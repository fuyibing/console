// author: wsfuyibing <websearch@163.com>
// date: 2021-02-25

package base

import (
	"errors"
	"fmt"
	"strconv"
	"sync"

	"github.com/fuyibing/console/v2/i"
)

// Command struct.
type Command struct {
	defaulter   bool                      // default status.
	description string                    // command description
	handler     func(iConsole i.IConsole) // command handler.
	hidden      bool                      // hidden status, do not list in console
	mu          *sync.RWMutex             // Mutex
	name        string                    // command name
	keys        []string                  //
	opts        map[string]i.IOption      // command options
	width       int
}

// Add option to command.
func (o *Command) Add(opts ...i.IOption) {
	o.mu.Lock()
	defer o.mu.Unlock()
	// command list.
	for _, opt := range opts {
		// unique control.
		if _, ok := o.opts[opt.GetName()]; ok {
			continue
		}
		// append option list.
		if size := len(opt.GetName()); o.width < size {
			o.width = size
		}
		o.opts[opt.GetName()] = opt
		o.keys = append(o.keys, opt.GetName())
	}
}

// Initialize struct fields.
func (o *Command) Initialize() {
	o.mu = new(sync.RWMutex)
	o.keys = make([]string, 0)
	o.opts = make(map[string]i.IOption)
}

// Is default command.
// Call when name not specified int command line.
func (o *Command) IsDefault() bool { return o.defaulter }

// Is hidden status.
// Do not list on console.
func (o *Command) IsHidden() bool { return o.hidden }

// Return command description.
func (o *Command) GetDescription() string { return o.description }

// Return command name.
func (o *Command) GetName() string { return o.name }

// Return specified option.
func (o *Command) GetOption(name string) i.IOption {
	o.mu.RLock()
	defer o.mu.RUnlock()
	if opt, ok := o.opts[name]; ok {
		return opt
	}
	return nil
}

// Return option map.
func (o *Command) GetOptions() map[string]i.IOption {
	o.mu.RLock()
	defer o.mu.RUnlock()
	return o.opts
}

// Set command description.
// Call in command new method, not export to interface.
func (o *Command) SetDefaulter(defaulter bool) { o.defaulter = defaulter }

// Set command description.
func (o *Command) SetDescription(description string) { o.description = description }

// Set command name.
func (o *Command) SetName(name string) { o.name = name }

// Set command hidden status.
// Call in command new method, not export to interface.
func (o *Command) SetHidden(hidden bool) { o.hidden = hidden }

// Run command.
func (o *Command) Run(console i.IConsole) {
	if o.handler != nil {
		o.handler(console)
		return
	}
	console.PrintError(errors.New(fmt.Sprintf("Command %s - Run() method not override", o.GetName())))
}

// Print command usage in console.
func (o *Command) Usage(console i.IConsole) {
	// Print usage.
	console.PrintUsage(o)
	// Print options.
	o.mu.RLock()
	defer o.mu.RUnlock()
	n := 0
	for _, key := range o.keys {
		if opt, ok := o.opts[key]; ok {
			console.PrintOptionItem(n, opt, false)
			n++
		}
	}
	console.PrintOptionItem(n, nil, true)
}

// Validate options.
func (o *Command) Validate(args []string) error {
	// Lock.
	o.mu.RLock()
	defer o.mu.RUnlock()
	for _, opt := range o.opts {
		if err := o.validateOption(opt, args); err != nil {
			return err
		}
	}
	return nil
}

// validate option.
func (o *Command) validateOption(opt i.IOption, args []string) error {
	// prepare.
	found := false
	name := opt.GetName()
	shortName := opt.GetShortName()
	// scan argument.
	for offset, arg := range args {
		if a := i.RegexpOptionShortName.FindStringSubmatch(arg); len(a) == 2 {
			for n := 0; n < len(a[1]); n++ {
				if s := string(a[1][n]); s == shortName {
					v := ""
					if offset < (len(args) - 1) {
						v = args[offset+1]
					}
					found = true
					if err := o.validateValue(opt, v); err != nil {
						return err
					}
				}
			}
		} else if b := i.RegexpOptionName.FindStringSubmatch(arg); len(b) == 3 {
			if b[1] == name {
				found = true
				if err := o.validateValue(opt, b[2]); err != nil {
					return err
				}
			}
		}
	}
	// found.
	if !found && opt.IsRequired() {
		return errors.New(fmt.Sprintf("Option '--%s' value not specified", opt.GetName()))
	}
	// completed.
	return nil
}

// validate value.
func (o *Command) validateValue(opt i.IOption, value string) error {
	// none value.
	if opt.IsNoneValue() {
		opt.SetValue("true")
		return nil
	}
	// bool value.
	if opt.IsBoolValue() {
		if value == "" {
			value = "true"
			if opt.IsRequired() {
				return errors.New(fmt.Sprintf("Option '--%s' value can not be empty", opt.GetName()))
			}
		}
		b, err := strconv.ParseBool(value)
		if err != nil {
			return errors.New(fmt.Sprintf("Option '--%s' value can not convert to boolean", opt.GetName()))
		}
		if b {
			opt.SetValue("true")
		} else {
			opt.SetValue("false")
		}
		return nil
	}
	// integer value.
	if opt.IsIntValue() {
		if value == "" {
			value = "0"
			if opt.IsRequired() {
				return errors.New(fmt.Sprintf("Option '--%s' value can not be empty", opt.GetName()))
			}
		}
		n, err := strconv.ParseInt(value, 0, 32)
		if err != nil {
			return errors.New(fmt.Sprintf("Option '--%s' value can not convert to integer", opt.GetName()))
		}
		opt.SetValue(fmt.Sprintf("%d", n))
		return nil
	}
	// string value
	if opt.IsStrValue() {
		if value == "" {
			if opt.IsRequired() {
				return errors.New(fmt.Sprintf("Option '--%s' value can not be empty", opt.GetName()))
			}
		}
		opt.SetValue(value)
		return nil
	}
	return nil
}
