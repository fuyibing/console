// author: wsfuyibing <websearch@163.com>
// date: 2021-03-23

package console

import (
	"fmt"
	"os"
	"sync"
)

// Console struct.
type Console struct {
	commandList map[string]*Command // Not sort.
	commandName []string            // Sorted able.
	defaultName string
	mu          *sync.RWMutex
	script      string
}

// Default console instance.
func Default() *Console {
	return New().Add(
		NewStart(),
		NewStop(),
		newBuildModelCommand(),
		newBuildServiceCommand(),
		newKvDownloadCommand(),
	)
}

// New console instance.
func New() *Console {
	return NewConsole().Add(
		newHelpCommand(),
	)
}

// New console instance.
func NewConsole() *Console {
	return &Console{
		mu:          new(sync.RWMutex),
		commandList: make(map[string]*Command),
		commandName: make([]string, 0),
	}
}

// Add commands to console.
func (o *Console) Add(commands ...*Command) *Console {
	// lock & unlock.
	o.mu.Lock()
	defer o.mu.Unlock()
	// loop commands.
	for _, command := range commands {
		// append name.
		if _, ok := o.commandList[command.name]; !ok {
			o.commandName = append(o.commandName, command.name)
		}
		// set or reset list.
		o.commandList[command.name] = command
	}
	// reset default command.
	defaultName := ""
	for _, command := range o.commandList {
		if command.isDefault {
			defaultName = command.name
			break
		}
	}
	o.defaultName = defaultName
	return o
}

// Delete commands from added list.
func (o *Console) Delete(commands ...*Command) *Console {
	for _, command := range commands {
		o.DeleteByName(command.name)
	}
	return o
}

// Delete command by name.
func (o *Console) DeleteByName(names ...string) *Console {
	// lock & unlock.
	o.mu.Lock()
	defer o.mu.Unlock()
	// loop names.
	for _, name := range names {
		// name not registered.
		if _, ok := o.commandList[name]; !ok {
			continue
		}
		// delete if exist.
		delete(o.commandList, name)
	}
	// remove from names.
	defaultName := ""
	commandName := make([]string, 0)
	for _, name := range o.commandName {
		if command, ok := o.commandList[name]; ok {
			commandName = append(commandName, name)
			if command.isDefault {
				defaultName = name
			}
		}
	}
	o.commandName = commandName
	o.defaultName = defaultName
	return o
}

// Get added command instance.
func (o *Console) GetCommand(name string) *Command {
	o.mu.RLock()
	defer o.mu.RUnlock()
	if command, ok := o.commandList[name]; ok {
		return command
	}
	return nil
}

// Get added command instance.
func (o *Console) GetCommands() []*Command {
	o.mu.RLock()
	defer o.mu.RUnlock()
	cs := make([]*Command, 0)
	for _, name := range o.commandName {
		if command, ok := o.commandList[name]; ok {
			cs = append(cs, command)
		}
	}
	return cs
}

// Run console.
func (o *Console) Run(args ...string) error {
	// reset arguments use command line.
	if args == nil || len(args) == 0 {
		args = os.Args
	}
	// arguments.
	argc := len(args)
	if argc == 0 {
		return fmt.Errorf("invalid arguments")
	}
	if o.script == "" {
		o.script = args[0]
	}
	// command name.
	name := o.defaultName
	if argc > 1 && RegexpArgumentCommandName.MatchString(args[1]) {
		name = args[1]
	}
	// read command and execute.
	if command, ok := o.commandList[name]; ok {
		return command.run(o, args...)
	}
	// error returned.
	if name == "" {
		return fmt.Errorf("command name not specified")
	}
	return fmt.Errorf("command name not defined: %s", name)
}

// Set console script.
func (o *Console) SetScript(script string) *Console {
	o.script = script
	return o
}
