// author: wsfuyibing <websearch@163.com>
// date: 2021-02-25

// Package command for help.
package help

import (
	"errors"
	"fmt"
	"os"

	"github.com/fuyibing/console/v2/base"
	"github.com/fuyibing/console/v2/i"
)

const (
	Description = ""
	Name        = "help"
)

// Command struct.
type command struct {
	base.Command
}

// New help command instance.
func New() i.ICommand {
	o := &command{}
	o.Initialize()
	// 1. normal.
	o.SetHidden(true)
	o.SetDefaulter(true)
	o.SetDescription(Description)
	o.SetName(Name)
	// 3. prepare.
	return o
}

// Print usage.
func (o *command) Run(console i.IConsole) {

	// Call command usage.
	args := os.Args
	if len(args) > 2 {
		if m := i.RegexpName.FindStringSubmatch(args[2]); len(m) > 0 {
			if cmd := console.GetCommand(m[1]); cmd != nil {
				cmd.Usage(console)
			} else {
				console.PrintError(errors.New(fmt.Sprintf("Command %s - command name not recognize: %s", o.GetName(), m[1])))
			}
			return
		}
	}
	// Print usage.
	console.PrintUsage(o)
	// Print command info.
	n := 0
	for _, k := range console.GetNames() {
		if c := console.GetCommand(k); c != nil {
			if c.IsHidden() {
				continue
			}
			console.PrintCommandItem(n, c, false)
			n++
		}
	}
	console.PrintCommandItem(n, nil, true)
	// Print ended.
	console.PrintCommandMore(o)
}
