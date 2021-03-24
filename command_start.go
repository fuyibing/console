// author: wsfuyibing <websearch@163.com>
// date: 2021-03-23

package console

import (
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

// Start command struct.
type startCommand struct {
	command *Command
}

// Handle after start command.
func (o *startCommand) after(cs *Console) error {
	if err := os.Remove(PidFile); err != nil {
		return fmt.Errorf("remove pid error: %v", err)
	}
	return nil
}

// Handle before start command.
func (o *startCommand) before(cs *Console) error {
	// Read pid value.
	pid, e1 := o.readPid()
	if e1 != nil {
		return e1
	}
	// Running pid.
	if pid > 0 {
		return fmt.Errorf("running pid: %d", pid)
	}
	// Write pid.
	if err := o.writePid(); err != nil {
		return fmt.Errorf("create pid error: %v", err)
	}
	return nil
}

// Read pid value.
func (o *startCommand) readPid() (int, error) {
	body, e1 := ioutil.ReadFile(PidFile)
	if e1 != nil {
		if os.IsNotExist(e1) {
			return 0, nil
		}
		return 0, fmt.Errorf("%v", e1)
	}
	text := strings.TrimSpace(string(body))
	if text == "" {
		return 0, nil
	}
	n, e2 := strconv.ParseInt(text, 0, 32)
	if e2 != nil {
		return 0, nil
	}
	return int(n), nil
}

// Write pid value.
func (o *startCommand) writePid() error {
	fp, err := os.OpenFile(PidFile, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, os.ModePerm)
	if err != nil {
		return err
	}
	defer func() {
		_ = fp.Close()
	}()
	if _, err = fp.WriteString(fmt.Sprintf("%d", os.Getpid())); err != nil {
		return err
	}
	return nil
}

// New start command.
//
//   cs := console.Default()
//   cs.Add(console.NewStart().SetHandler(func(cs *console.Console)error{
//       return nil
//   }))
func NewStart() *Command {
	c := NewCommand("start")
	c.SetDescription("Start application")
	o := &startCommand{command: c}
	c.SetHandlerAfter(o.after).SetHandlerBefore(o.before)
	return c
}
