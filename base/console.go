// author: wsfuyibing <websearch@163.com>
// date: 2021-02-25

// Package for base.
package base

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"sync"

	"github.com/fuyibing/log/v2"

	"github.com/fuyibing/console/v2/i"
)

// Console struct.
type console struct {
	dc       i.ICommand            // default command.
	cs       map[string]i.ICommand // command list.
	ns       []string              // command name list.
	mu       *sync.RWMutex         // mutex
	cmdWidth int                   // command with
}

// Add command to console.
func (o *console) Add(cs ...i.ICommand) {
	o.mu.Lock()
	defer o.mu.Unlock()
	// command list.
	for _, c := range cs {
		// unique control.
		if _, ok := o.cs[c.GetName()]; ok {
			continue
		}
		// set default.
		if c.IsDefault() {
			o.dc = c
		}
		// append command list.
		o.cs[c.GetName()] = c
		o.ns = append(o.ns, c.GetName())
		// width
		w := len(c.GetName())
		// command with
		if o.cmdWidth < w {
			o.cmdWidth = w
		}
	}
}

// Delete command.
func (o *console) Del(cs ...i.ICommand) {
	o.mu.Lock()
	defer o.mu.Unlock()
	for _, c := range cs {
		// not added.
		if _, ok := o.cs[c.GetName()]; !ok {
			continue
		}
		// remove default command if specified command is
		// marked as default.
		if c.IsDefault() && o.dc != nil && c.GetName() == o.dc.GetName() {
			o.dc = nil
		}
		// remove specified command.
		delete(o.cs, c.GetName())
	}
}

// Return command.
func (o *console) GetCommand(name string) i.ICommand {
	o.mu.RLock()
	defer o.mu.RUnlock()
	if c, ok := o.cs[name]; ok {
		return c
	}
	return nil
}

// Return all commands.
func (o *console) GetCommands() map[string]i.ICommand {
	o.mu.RLock()
	defer o.mu.RUnlock()
	return o.cs
}

// Return default command.
func (o *console) GetDefaultCommand() i.ICommand { return o.dc }

// Return name list.
func (o *console) GetNames() []string {
	o.mu.RLock()
	defer o.mu.RUnlock()
	return o.ns
}

// Print info.
func (o *console) Info(text string, args ...interface{}) {
	println(fmt.Sprintf(text, args...))
}

// Print command item.
func (o *console) PrintCommandItem(n int, cmd i.ICommand, end bool) {
	// end command.
	if end {
		if n > 0 {
			println("└─────────────┴────────────────────────────────────────────────────────────────────────────────────┘")
		} else {
			println("└──────────────────────────────────────────────────────────────────────────────────────────────────┘")
		}
		return
	}
	// begin command, first item.
	if n == 0 {
		println("├─────────────┬────────────────────────────────────────────────────────────────────────────────────┤")
	}
	// normal
	var format, prefix, script, desc = "", "", "", ""
	// prefix.
	if prefix = fmt.Sprintf("Commands %2d", n+1); n > 0 {
		prefix = fmt.Sprintf("         %2d", n+1)
	}
	// format.
	format = fmt.Sprintf("%%-%ds", o.cmdWidth+4)
	script = fmt.Sprintf(format, cmd.GetName())
	if desc = cmd.GetDescription(); len(desc) > 55 {
		desc = desc[0:55] + "..."
	}
	// print.
	fmt.Printf("│ %-11s │ %-82s │\n", prefix, strings.TrimSpace(script+desc))
}

// Print more guide on command list.
func (o *console) PrintCommandMore(cmd i.ICommand) {
	fmt.Printf(
		"Run '%s help %s' for more information on a command.\n",
		i.Script,
		i.UsageDefaultCommand,
	)
}

// Print option item.
func (o *console) PrintOptionItem(n int, opt i.IOption, end bool) {
	// end command.
	if end {
		if n > 0 {
			println("└─────────────┴────────────────────────────────────────────────────────────────────────────────────┘")
		} else {
			println("└──────────────────────────────────────────────────────────────────────────────────────────────────┘")
		}
		return
	}
	// normal
	var prefix, script, desc, value = "", "", "", ""
	// prefix.
	if prefix = fmt.Sprintf("Options  %2d", n+1); n > 0 {
		prefix = fmt.Sprintf("         %2d", n+1)
	}
	// option:name:short.
	if sn := opt.GetShortName(); sn != "" {
		script = fmt.Sprintf("-%s,", sn)
	} else {
		script = "   "
	}
	// option:name.
	script += fmt.Sprintf("--%s", opt.GetName())
	// option:value
	if !opt.IsNoneValue() {
		if value = opt.GetDefaultValue(); value == "" {
			if opt.IsIntValue() {
				value = "int"
			} else if opt.IsStrValue() {
				value = "str"
			} else {
				value = "value"
			}
		}
		if opt.IsRequired() {
			script += fmt.Sprintf("=<%s>", value)
		} else {
			script += fmt.Sprintf("[=%s]", value)
		}
	}
	script = fmt.Sprintf("%-28s", script)
	// desc.
	if desc = opt.GetDescription(); len(desc) > 44 {
		desc = desc[0:44] + "..."
	}
	// print.
	// begin command, first item.
	if n == 0 {
		println("├─────────────┬────────────────────────────────────────────────────────────────────────────────────┤")
	}
	fmt.Printf("│ %-11s │ %-82s │\n", prefix, script+strings.TrimSpace(desc))
}

// Print usage.
func (o *console) PrintUsage(cmd i.ICommand) {
	// print usage.
	println("┌──────────────────────────────────────────────────────────────────────────────────────────────────┐")
	fmt.Printf("│ %-96s │\n", fmt.Sprintf("Usage: %s %s %s", i.Script, cmd.GetName(), i.UsageDefaultOption))
	// print description.
	if desc := cmd.GetDescription(); desc != "" {
		println("│ ________________________________________________________________________________________________ │")
		buf := ""
		num := 0
		for _, str := range strings.Split(desc, " ") {
			str = strings.TrimSpace(str)
			if str == "" {
				continue
			}
			size := len(str)
			if (size + num) >= 82 {
				fmt.Printf("│ %-96s │\n", strings.TrimSpace(buf))
				num = size
				buf = str
			} else {
				buf += " " + str
				num += size
			}
		}
		if buf != "" {
			fmt.Printf("│ %-96s │\n", strings.TrimSpace(buf))
		}
	}
}

// Print error.
func (o *console) PrintError(err error) {
	println(fmt.Sprintf("%c[%d;%dm%5s%c[0m", 0x1B, 0, 31, err.Error(), 0x1B))
}

// Run console.
func (o *console) Run(args ...string) {
	// Use args from command line.
	if len(args) == 0 {
		args = os.Args
	}
	// Arg 0: reset as script.
	if len(args) == 0 {
		args = append(args, i.Script)
	} else {
		args[0] = i.Script
	}
	// Arg 1: command name.
	if len(args) == 1 {
		if o.dc == nil {
			args = append(args, "")
		} else {
			args = append(args, o.dc.GetName())
		}
	}
	// Print error if command is empty.
	if args[1] == "" {
		o.PrintError(errors.New(fmt.Sprintf("Command - command name not specified")))
		return
	}
	// Validate command name.
	if i.RegexpName.MatchString(args[1]) {
		if cmd := o.GetCommand(args[1]); cmd != nil {
			// Call command.
			if err := cmd.Validate(args[2:]); err != nil {
				o.PrintError(err)
				return
			}
			log.Config.SetLevel("off")
			cmd.Run(o)
		} else {
			// Print error if command name invalid.
			o.PrintError(errors.New(fmt.Sprintf("Command - command name not registed: %s", args[1])))
		}
	} else {
		// Print error if command name not validated.
		o.PrintError(errors.New(fmt.Sprintf("Command - command name not recognize: %s", args[1])))
	}
}

// Return new console.
func NewConsole() i.IConsole {
	o := &console{}
	o.cs = make(map[string]i.ICommand)
	o.mu = new(sync.RWMutex)
	o.ns = make([]string, 0)
	return o
}
