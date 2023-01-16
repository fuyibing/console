// author: wsfuyibing <websearch@163.com>
// date: 2023-01-11

package managers

import (
	"fmt"
	"os"
)

const (
	Version = "3.0.0"
)

type (
	// Manager
	// operation interface.
	Manager interface {
		AddCommand(c Command) error
		GetCommand(key string) Command
		GetCommands() map[string]Command
		GetDescription() string
		Run(a Arguments) error
		RunTerminal() error
		SetDescription(s string) Manager
	}

	manager struct {
		Commands    map[string]Command
		Description string
	}
)

func NewManager() Manager {
	return &manager{
		Commands: make(map[string]Command),
	}
}

// /////////////////////////////////////////////////////////////
// Interface methods
// /////////////////////////////////////////////////////////////

func (o *manager) AddCommand(c Command) error      { return o.addCommand(c) }
func (o *manager) GetCommand(key string) Command   { return o.getCommand(key) }
func (o *manager) GetCommands() map[string]Command { return o.Commands }
func (o *manager) GetDescription() string          { return o.Description }
func (o *manager) Run(a Arguments) error           { return o.run(a) }
func (o *manager) RunTerminal() error              { return o.runTerminal() }
func (o *manager) SetDescription(s string) Manager { o.Description = s; return o }

// /////////////////////////////////////////////////////////////
// Access and constructor
// /////////////////////////////////////////////////////////////

func (o *manager) addCommand(c Command) error {
	if c == nil {
		return nil
	}

	// Name required.
	if c.GetName() == "" {
		return fmt.Errorf("can not add unnamed command to manager")
	}

	// Add twice.
	if _, ok := o.Commands[c.GetName()]; ok {
		return fmt.Errorf("command exists in manager: %s", c.GetName())
	}

	// Set mapper.
	o.Commands[c.GetName()] = c
	return nil
}

func (o *manager) getCommand(key string) Command {
	if c, ok := o.Commands[key]; ok {
		return c
	}
	return nil
}

func (o *manager) run(a Arguments) error {
	var (
		cmd      Command
		exists   bool
		selector = a.GetSelector()
	)

	// Use default command
	// if selected is empty.
	if selector == "" {
		selector = ArgumentsHelp
	}

	// Read command from mapper.
	if cmd, exists = o.Commands[selector]; exists {
		// Return error
		// if arguments option not registered in command.
		for ak, av := range a.GetMapper() {
			if co := cmd.GetOption(ak); co != nil {
				if err := co.Assign(av); err != nil {
					return err
				}
				continue
			}
			return fmt.Errorf("option not recognized: %s", ak)
		}

		// Return error
		// if command option validate failed.
		for _, cv := range cmd.GetOptions() {
			if err := cv.Validate(); err != nil {
				return err
			}
		}

		// Run command.
		return cmd.Run(o, a)
	}

	// Return error
	// if command not registered.
	return fmt.Errorf("command not registered in manager: %s", selector)
}

func (o *manager) runTerminal() error {
	a := NewArguments()

	if err := a.Parse(os.Args...); err != nil {
		return err
	}

	return o.run(a)
}
