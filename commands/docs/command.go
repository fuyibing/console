// author: wsfuyibing <websearch@163.com>
// date: 2023-01-11

// Package docs
// generate application documents as markdown files or postman collection and so on.
package docs

import (
	"github.com/fuyibing/console/v3/managers"
)

type Command struct {
	Command managers.Command
	Err     error
	Name    string
}

// Handle
// callable registered on command manager interface.
func (o *Command) Handle(m managers.Manager, a managers.Arguments) error {
	return nil
}

// /////////////////////////////////////////////////////////////
// Access and constructor methods
// /////////////////////////////////////////////////////////////

func (o *Command) initField() *Command {
	o.Command = managers.NewCommand(o.Name)
	o.Command.SetHandler(o.Handle)
	o.Command.SetDescription("Generate application documents as markdown files or postman collection and so on")
	return o
}

func (o *Command) initOption() *Command {
	o.Err = o.Command.AddOption(
		managers.NewOption("adapter").SetShortName('a').SetDescription("Specify document formatter, accept: postman, markdown").SetDefault("markdown"),
		managers.NewOption("base").SetShortName('b').SetDescription("Specify your working base path").SetDefault("./"),
		managers.NewOption("controller").SetShortName('c').SetDescription("Specify your controller path").SetDefault("/app/controllers"),
		managers.NewOption("document").SetShortName('d').SetDescription("Built documents storage location").SetDefault("/docs/api"),
	)

	return o
}

// New
// function create and return instance.
func New() (managers.Command, error) {
	o := (&Command{Name: "docs"}).
		initField().
		initOption()

	return o.Command, o.Err
}
