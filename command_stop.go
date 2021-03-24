// author: wsfuyibing <websearch@163.com>
// date: 2021-03-23

package console

import (
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"syscall"
)

// Stop command struct.
type stopCommand struct {
	command *Command
}

// Handle after stop command.
func (o *stopCommand) after(cs *Console) error {
	if err := os.Remove(PidFile); err != nil {
		return fmt.Errorf("remove pid error: %v", err)
	}
	return nil
}

// Handle before start command.
func (o *stopCommand) before(*Console) error {
	// return if read pid error.
	pid, e1 := o.readPid()
	if e1 != nil {
		return e1
	}
	// return if pid is zero, means not running.
	if pid == 0 {
		return fmt.Errorf("not running")
	}
	// return error if find process error.
	proc, e2 := os.FindProcess(pid)
	if e2 != nil {
		return fmt.Errorf("find process error: %v", e2)
	}
	// send sigterm signal.
	_ = proc.Signal(syscall.SIGTERM)
	return nil
}

// Read pid.
// Read process id from pid file, default pid file
// is ./pid (in working directory).
func (o *stopCommand) readPid() (int, error) {
	// read pid file.
	body, e1 := ioutil.ReadFile(PidFile)
	if e1 != nil {
		// return zero if file not exist.
		if os.IsNotExist(e1) {
			return 0, nil
		}
		// return error if read file.
		return 0, fmt.Errorf("%v", e1)
	}
	// parse file content.
	text := strings.TrimSpace(string(body))
	if text == "" {
		return 0, fmt.Errorf("empty pid file")
	}
	// parse content format.
	n, e2 := strconv.ParseInt(text, 0, 32)
	if e2 != nil {
		return 0, fmt.Errorf("invalid pid: %s", text)
	}
	// read succeed.
	return int(n), nil
}

// New start command.
//
//   cs := console.Default()
//   cs.Add(console.NewStop().SetHandler(func(cs *console.Console)error{
//       return nil
//   }))
func NewStop() *Command {
	c := NewCommand("stop")
	c.SetDescription("Stop application")
	o := &stopCommand{command: c}
	c.SetHandlerAfter(o.after).SetHandlerBefore(o.before)
	return c
}
