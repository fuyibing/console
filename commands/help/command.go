// author: wsfuyibing <websearch@163.com>
// date: 2023-01-11

// Package help print manager and command information on
// terminal.
//
// # Source code mode
//   go run main.go
//   go run main.go help
//   go run main.go help docs
//
// # Built file mode
//   ./demo
//   ./demo help
//   ./demo help docs
package help

import (
	"fmt"
	"github.com/fuyibing/console/v3/managers"
	"os"
	"sort"
	"strings"
)

const (
	commandName  = "help"
	commandWidth = 100
)

type Command struct {
	Command managers.Command
	Err     error
	Name    string
}

// Handle
// callable registered on command manager interface.
func (o *Command) Handle(m managers.Manager, a managers.Arguments) error {
	// Handle command.
	if key := a.GetHelpSelector(); key != "" {
		if c := m.GetCommand(key); c != nil {
			return o.HandleCommand(a, c)
		}

		// Return error if command not recognize.
		return fmt.Errorf("command not recognized: %s", key)
	}

	// Handle manager.
	return o.HandleManager(m, a)
}

// HandleCommand
// generate command information and print.
func (o *Command) HandleCommand(a managers.Arguments, c managers.Command) error {
	o.RenderVersion()
	o.RenderUsage(a.GetScript(), c.GetName())
	o.RenderDescription(c.GetDescription())

	o.RenderOption(c)
	return nil
}

// HandleManager
// generate manager information and print.
func (o *Command) HandleManager(m managers.Manager, a managers.Arguments) error {
	o.RenderVersion()
	o.RenderUsage(a.GetScript(), "COMMAND")
	o.RenderDescription(m.GetDescription())

	o.RenderOption(o.Command)
	o.RenderCommands(m)
	o.RenderGuider(a.GetScript())
	return nil
}

// RenderCommands
// print manager command list.
//
//   Commands:
//     docs         Build application document files
//     kv:upload    Collect local config files and upload to consul
func (o *Command) RenderCommands(m managers.Manager) {
	var (
		c            managers.Command
		index, width = 0, 0
		keys         = make([]string, 0)
	)

	// Range commands.
	for _, c = range m.GetCommands() {
		if c.GetHidden() {
			continue
		}

		// Set maximum width of command.
		if w := len(c.GetName()); width < w {
			width = w
		}

		keys = append(keys, c.GetName())
	}

	// Sort
	// by command name.
	sort.Strings(keys)

	// Make formatters.
	var (
		format = fmt.Sprintf("  %%-%ds    %%s", width)
		holder = fmt.Sprintf("  %s    %%s", strings.Repeat(" ", width))
	)

	// Range commands.
	for _, key := range keys {
		if c = m.GetCommand(key); c == nil {
			continue
		}

		// Print header
		// if command index is zero.
		if index++; index == 1 {
			o.println("")
			o.println("Commands:")
		}

		// Print command.
		if ss := o.SplitWords(width, c.GetDescription()); len(ss) > 0 {
			// Multi-rows description.
			for i, s := range ss {
				if i == 0 {
					// First row.
					o.println(format, key, s)
				} else {
					// Not first row.
					o.println(holder, s)
				}
			}
		} else {
			// No description command.
			o.println(format, key, "")
		}
	}
}

// RenderDescription
// print command description information.
func (o *Command) RenderDescription(str string) {
	for i, s := range o.SplitWords(0, str) {
		if i == 0 {
			o.println("")
		}
		o.println("%s", s)
	}
}

// RenderGuider
// print guide information.
func (o *Command) RenderGuider(script string) {
	o.println("")
	o.println("Run '%s help COMMAND' for more information on a command", script)
	o.println("")
	o.println("To get more help with console, check out our guides at https://github.com/fuyibing/console/tree/v3")
}

// RenderOption
// print option information.
//
//   Options:
//     -b, --base=string      Specify working base path
//         --config=string    Specify config path
func (o *Command) RenderOption(c managers.Command) {
	var (
		index, width = 0, 0
		keys         = make([]string, 0)
		opt          managers.Option
		opts         = make(map[string]managers.Option)
	)

	// Range
	// command options.
	for _, opt = range c.GetOptions() {
		keys = append(keys, opt.GetName())
		opts[opt.GetName()] = opt

		// Set maximum width of label.
		if n := len(opt.GetLabel()); width < n {
			width = n
		}
	}

	// Sort
	// by option name.
	sort.Strings(keys)

	// Make formatters.
	var (
		format = fmt.Sprintf("  %%-%ds    %%s", width)
		holder = fmt.Sprintf("  %s    %%s", strings.Repeat(" ", width))
	)

	// Range
	// option names.
	for _, key := range keys {
		if opt = opts[key]; opt == nil {
			continue
		}

		// Print header
		// if option index is zero.
		if index++; index == 1 {
			o.println("")
			o.println("Options:")
		}

		// Print options.
		if cs := o.SplitWords(width, opt.GetDescription()); len(cs) > 0 {
			// Multi-rows description on option.
			for i, s := range cs {
				if i == 0 {
					// First row
					// of multi-rows.
					o.println(format, opt.GetLabel(), s)
				} else {
					// Not first row
					// of multi-rows.
					o.println(holder, s)
				}
			}
		} else {
			// No description.
			o.println(format, opt.GetLabel(), "")
		}
	}
}

// RenderUsage
// print usage information.
//
//   Usage: ./app COMMAND [OPTION]
//   Usage: ./app help COMMAND
func (o *Command) RenderUsage(script, name string) {
	o.println("Usage: %s %s [OPTIONS]", script, name)
}

// RenderVersion
// print version information.
//
//   Version: 3.0.0
func (o *Command) RenderVersion() {
	o.println("")
	o.println("Version: %s", managers.Version)
}

// SplitWords
// convert long-text string as multi-rows slice with specified
// width.
func (o *Command) SplitWords(w int, str string) []string {
	var (
		ln, n     = 0, 0
		tmp, rows = make([]string, 0), make([]string, 0)
		width     = commandWidth - w
	)

	// Range words by empty space.
	for _, word := range strings.Split(str, " ") {
		if word = strings.TrimSpace(word); word == "" {
			continue
		}

		// Append to tmp
		// if length is less than limit.
		if n = len(word); (ln + n + 1) < width {
			ln += n + 1
			tmp = append(tmp, word)
			continue
		}

		// Collect mid-words
		// to rows when tmp words width is greater than limit.
		rows = append(rows, strings.Join(tmp, " "))

		// Reset
		// tmp slice and length.
		tmp = []string{word}
		ln = n + 1
	}

	// Collect end-words
	// to rows if tmp is not empty.
	if len(tmp) > 0 {
		rows = append(rows, strings.Join(tmp, " "))
	}

	// Return
	// multi-rows slice.
	return rows
}

// /////////////////////////////////////////////////////////////
// Access and constructor methods
// /////////////////////////////////////////////////////////////

func (o *Command) initField() *Command {
	o.Command = managers.NewCommand(o.Name)
	o.Command.SetHidden(true).SetHandler(o.Handle)
	return o
}

func (o *Command) initOption() *Command {
	o.Err = o.Command.AddOption()
	return o
}

func (o *Command) println(format string, args ...interface{}) {
	_, _ = fmt.Fprintf(os.Stdout, "%s\n", fmt.Sprintf(format, args...))
}

// New
// function create and return instance.
func New() (managers.Command, error) {
	o := (&Command{Name: commandName}).
		initField().
		initOption()

	return o.Command, o.Err
}
