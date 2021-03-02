// author: wsfuyibing <websearch@163.com>
// date: 2021-02-25

// Package command for build application service.
package service

import (
	"github.com/fuyibing/console/v2/base"
	"github.com/fuyibing/console/v2/i"
)

const (
	Description = "Build server file for iris application"
	Name        = "build:service"
)

// Command struct.
type command struct {
	base.Command
	packages map[string]int
}

// New build service instance.
func New() i.ICommand {
	o := &command{packages: make(map[string]int)}
	o.Initialize()
	o.SetDescription(Description)
	o.SetName(Name)
	// service name.
	o.Add(base.NewOption(i.RequiredMode, i.StrValue).SetName("name").SetShortName("n").SetDescription("Service name, no suffix, equal to model name."))
	// application path.
	o.Add(base.NewOption(i.OptionalMode, i.StrValue).SetName("path").SetShortName("p").SetDefaultValue("./app").SetDescription("Application path."))
	// override if file exist.
	//   -o
	//   --override
	o.Add(base.NewOption(i.OptionalMode, i.BoolValue).SetName("override").SetShortName("o").SetDescription("Override if file exist"))
	// with
	o.Add(base.NewOption(i.OptionalMode, i.BoolValue).SetDefaultValue(true).SetName("no-add").SetDescription("Export Add(req *Model) method"))
	o.Add(base.NewOption(i.OptionalMode, i.BoolValue).SetDefaultValue(true).SetName("no-get").SetDescription("Export Get(req *Model) method"))
	o.Add(base.NewOption(i.OptionalMode, i.BoolValue).SetDefaultValue(true).SetName("no-get-by-id").SetDescription("Export GetById(id int) method"))
	// prepared.
	return o
}

// Run command.
func (o *command) Run(console i.IConsole) {
}

func (o *command) dumpHead()          {}
func (o *command) dumpBody()          {}
func (o *command) dumpMethodAdd()     {}
func (o *command) dumpMethodGet()     {}
func (o *command) dumpMethodGetById() {}
