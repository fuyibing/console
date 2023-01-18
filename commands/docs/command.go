// author: wsfuyibing <websearch@163.com>
// date: 2023-01-11

// Package docs
// generate application documents as markdown files or postman collection and so on.
package docs

import (
	"fmt"
	"github.com/fuyibing/console/v3/managers"
	"github.com/fuyibing/gdoc/adapters/markdown"
	"github.com/fuyibing/gdoc/adapters/postman"
	"github.com/fuyibing/gdoc/base"
	"github.com/fuyibing/gdoc/conf"
	"github.com/fuyibing/gdoc/reflectors"
	"github.com/fuyibing/gdoc/scanners"
)

const (
	CmdDesc = "Generate application documents as markdown files or postman collection and so on"
	CmdName = "docs"

	OptAdapter        = "adapter"
	OptAdapterByte    = 'a'
	OptAdapterDesc    = "Specify document formatter, accept: postman, markdown"
	OptAdapterDefault = "markdown"

	OptBase        = "base"
	OptBaseByte    = 'b'
	OptBaseDesc    = "Specify your working base path"
	OptBaseDefault = "./"

	OptController        = "controller"
	OptControllerByte    = 'c'
	OptControllerDesc    = "Specify your controller path"
	OptControllerDefault = "/app/controllers"

	OptDocument        = "document"
	OptDocumentByte    = 'd'
	OptDocumentDesc    = "Built documents storage location"
	OptDocumentDefault = "/docs/api"
)

type Command struct {
	Command managers.Command
	Err     error
	Name    string
}

// Handle
// callable registered on command manager interface.
func (o *Command) Handle(_ managers.Manager, _ managers.Arguments) (err error) {
	var (
		s1, s2, s3, s4 string
	)

	// Read
	// options value.
	if s1, err = o.Command.GetOption(OptAdapter).ToString(); err != nil {
		return
	}
	if s2, err = o.Command.GetOption(OptBase).ToString(); err != nil {
		return
	}
	if s3, err = o.Command.GetOption(OptController).ToString(); err != nil {
		return
	}
	if s4, err = o.Command.GetOption(OptDocument).ToString(); err != nil {
		return
	}

	// Use
	// option value.
	conf.Path.SetBasePath(s2)
	conf.Path.SetControllerPath(s3)
	conf.Path.SetDocumentPath(s4)

	// Scan
	// controller files.
	scanners.Scanner.Scan()

	// Reflect.
	ref := reflectors.New(base.Mapper)
	ref.Configure()

	if err = ref.Make(); err != nil {
		return
	}

	switch s1 {
	case "postman":
		postman.New(base.Mapper).Run()
	case "markdown":
		markdown.New(base.Mapper).Run()
	default:
		err = fmt.Errorf("unknown adapter")
	}

	if err != nil {
		ref.Clean()
	}

	return
}

// /////////////////////////////////////////////////////////////
// Access and constructor methods
// /////////////////////////////////////////////////////////////

func (o *Command) InitField() *Command {
	o.Command = managers.NewCommand(o.Name)
	o.Command.SetDescription(CmdDesc).SetHandler(o.Handle)
	return o
}

func (o *Command) InitOption() *Command {
	o.Err = o.Command.AddOption(
		managers.NewOption(OptAdapter).SetShortName(OptAdapterByte).SetDescription(OptAdapterDesc).SetDefault(OptAdapterDefault),
		managers.NewOption(OptBase).SetShortName(OptBaseByte).SetDescription(OptBaseDesc).SetDefault(OptBaseDefault),
		managers.NewOption(OptController).SetShortName(OptControllerByte).SetDescription(OptControllerDesc).SetDefault(OptControllerDefault),
		managers.NewOption(OptDocument).SetShortName(OptDocumentByte).SetDescription(OptDocumentDesc).SetDefault(OptDocumentDefault),
	)

	return o
}

// New
// function create and return instance.
//
//   go run main.go docs \
//     --adapter=markdown \
//     --base=./ \
//     --controller=/app/controllers \
//     --document=/docs/api
func New() (managers.Command, error) {
	o := (&Command{Name: CmdName}).
		InitField().
		InitOption()

	return o.Command, o.Err
}
