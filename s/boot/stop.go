// author: wsfuyibing <websearch@163.com>
// date: 2021-03-13

package boot

import (
	"fmt"
	"os"
	"syscall"

	"github.com/fuyibing/console/v2/base"
	"github.com/fuyibing/console/v2/i"
)

type stop struct {
	base.Command
}

func NewStop() i.ICommand {
	// 1. normal.
	o := &stop{}
	o.Initialize()
	o.SetDescription(StopDescription)
	o.SetName(StopName)
	// 2. before stop.
	o.SetHandlerBefore(o.beforeStop)
	// n. completed.
	return o
}

// Before stop executed.
func (o *stop) beforeStop(c i.IConsole) bool {
	pid, _, e0 := read()
	// Application stop.
	c.Info("App stop.")
	c.Info("    > pid = %d.", pid)
	c.Info("    > file = %s", PidFile)
	// Not started.
	if pid == 0 {
		if e0 != nil {
			c.PrintError(fmt.Errorf("Read pid error: %v.", e0))
		} else {
			c.PrintError(fmt.Errorf("No started pid found."))
		}
		return false
	}
	// Send signal error.
	if e1 := syscall.Kill(pid, syscall.SIGTERM); e1 != nil {
		c.PrintError(fmt.Errorf("Error: %v.", e1))
		return false
	}
	// Remove pid file error.
	if e2 := os.Remove(PidFile); e2 != nil {
		c.PrintError(fmt.Errorf("Error: %v.", e2))
		return false
	}

	return true
}
