// author: wsfuyibing <websearch@163.com>
// date: 2021-03-23

package console

import (
	"fmt"
	"sync"
)

// Command handler.
type CommandHandler func(console *Console) error

// Command struct.
type Command struct {
	description   string
	isDefault     bool
	isHidden      bool
	mu            *sync.RWMutex
	name          string
	optionList    map[string]*Option
	optionName    []string
	handler       CommandHandler
	handlerAfter  CommandHandler
	handlerBefore CommandHandler
}

// New command instance.
func NewCommand(name string) *Command {
	return &Command{
		mu:         new(sync.RWMutex),
		name:       name,
		optionList: make(map[string]*Option),
		optionName: make([]string, 0),
	}
}

// Add option.
func (o *Command) Add(options ...*Option) *Command {
	// lock & unlock.
	o.mu.Lock()
	defer o.mu.Unlock()
	// add option to command.
	for _, option := range options {
		if _, ok := o.optionList[option.name]; !ok {
			o.optionName = append(o.optionName, option.name)
		}
		o.optionList[option.name] = option
	}
	return o
}

// Set default status.
func (o *Command) Default(is bool) *Command {
	o.isDefault = is
	return o
}

// Set hidden status.
func (o *Command) Hidden(is bool) *Command {
	o.isHidden = is
	return o
}

// Get option by name.
func (o *Command) GetOption(name string) *Option {
	o.mu.RLock()
	defer o.mu.RUnlock()
	if option, ok := o.optionList[name]; ok {
		return option
	}
	return nil
}

// Get option list.
func (o *Command) GetOptions() []*Option {
	o.mu.RLock()
	defer o.mu.RUnlock()
	ls := make([]*Option, 0)
	for _, name := range o.optionName {
		if option, ok := o.optionList[name]; ok {
			ls = append(ls, option)
		}
	}
	return ls
}

// Set command description.
func (o *Command) SetDescription(description string) *Command {
	o.description = description
	return o
}

// Set handler.
// Handler callback when command execute.
func (o *Command) SetHandler(handler CommandHandler) *Command {
	o.handler = handler
	return o
}

// Set after handler.
// Execute after command executed.
func (o *Command) SetHandlerAfter(handler CommandHandler) *Command {
	o.handlerAfter = handler
	return o
}

// Set before handler.
// Execute before handler executed.
func (o *Command) SetHandlerBefore(handler CommandHandler) *Command {
	o.handlerBefore = handler
	return o
}

// Run command.
func (o *Command) run(console *Console, args ...string) error {
	// Handler defined check.
	if o.handler == nil {
		return fmt.Errorf("handler not defined: %s", o.name)
	}
	// Validate options.
	o.mu.RLock()
	for _, option := range o.optionList {
		if err := option.validate(args...); err != nil {
			o.mu.RUnlock()
			return err
		}
	}
	o.mu.RUnlock()
	// Call before handler.
	if o.handlerBefore != nil {
		if err := o.handlerBefore(console); err != nil {
			return err
		}
	}
	// Call main handler.
	if err := o.handler(console); err != nil {
		return err
	}
	// Call after.
	if o.handlerAfter != nil {
		if err := o.handlerAfter(console); err != nil {
			return err
		}
	}
	// Succeed.
	return nil
}
