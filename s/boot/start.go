// author: wsfuyibing <websearch@163.com>
// date: 2021-03-13

package boot

import (
	"fmt"
	"os"

	"github.com/fuyibing/console/v2/base"
	"github.com/fuyibing/console/v2/i"
)

type start struct {
	base.Command
}

func NewStart() i.ICommand {
	// 1. normal.
	o := &start{}
	o.Initialize()
	o.SetDescription(StartDescription)
	o.SetName(StartName)
	// 2. before start.
	o.SetHandlerBefore(o.beforeStart)
	// n. completed
	return o
}

// Before start.
func (o *start) beforeStart(c i.IConsole) bool {
	// current pid.
	pid := os.Getpid()
	// Application start.
	c.Info("App start.")
	c.Info("    > pid = %d.", pid)
	c.Info("    > file = %s", PidFile)
	// Running error.
	if value, _, _ := read(); value > 0 {
		c.PrintError(fmt.Errorf("Another pid %d is running.", value))
		return false
	}
	// Can not create pid file.
	if err := write(pid); err != nil {
		c.PrintError(fmt.Errorf("Can not create pid file: %v", err))
		return false
	}
	// Allow call handler.
	return true
}
