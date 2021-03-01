// author: wsfuyibing <websearch@163.com>
// date: 2021-02-26

package path

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/fuyibing/console/v2/base"
	"github.com/fuyibing/console/v2/i"
)

const (
	Description = "Build path for iris application"
	Name        = "build:path"
)

// Command struct.
type command struct {
	base.Command
}

// New build path instance.
func New() i.ICommand {
	// normal.
	o := &command{}
	o.Initialize()
	o.SetDescription(Description)
	o.SetName(Name)
	// app path.
	o.Add(
		base.NewOption(i.OptionalMode, i.StrValue).
			SetName("path").SetShortName("p").
			SetDefaultValue("./app").
			SetDescription("Application working path"),
	)
	return o
}

// Run command.
func (o *command) Run(console i.IConsole) {
	console.Info("Command %s: begin", o.GetName())
	defer console.Info("Command %s: completed", o.GetName())
	p := o.GetOption("path").ToString()
	for _, name := range []string{
		"commands",
		"controllers",
		"models",
		"logics",
		"services",
	} {
		path := fmt.Sprintf("%s/%s", p, name)
		if err := o.makeDir(console, path); err != nil {
			console.PrintError(err)
			return
		}
	}
}

// Make directory.
func (o *command) makeDir(console i.IConsole, path string) error {
	if err := os.MkdirAll(path, os.ModePerm); err != nil {
		return errors.New(fmt.Sprintf("Command %s: make directory fail: %v", o.GetName(), err))
	}
	console.Info("        make: %s", path)
	// File: open and close.
	file := path + "/.gitKeep"
	handler, err := os.OpenFile(file, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, os.ModePerm)
	if err != nil {
		return errors.New(fmt.Sprintf("Command %s: build file fail: %v", o.GetName(), err))
	}
	defer func() {
		_ = handler.Close()
	}()
	if _, err = handler.WriteString("# date: " + time.Now().Format("2006-01-02 15:04:05")); err != nil {
		return errors.New(fmt.Sprintf("Command %s: write file fail: %v", o.GetName(), err))
	}
	// Succeed.
	console.Info("        open: %s", file)
	return nil
}
