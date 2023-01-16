// author: wsfuyibing <websearch@163.com>
// date: 2023-01-11

// Package console define quick guider functions.
//
// > Run with source code
//   go run main.go
//   go run main.go help
//   go run main.go help docs
//
// > Build binary file named as demo
//   go build -o demo
//
// > Run with binary file
//   ./demo
//   ./demo help
//   ./demo help docs
package console

import (
	"github.com/fuyibing/console/v3/commands/consul/kv/download"
	"github.com/fuyibing/console/v3/commands/consul/kv/upload"
	"github.com/fuyibing/console/v3/commands/consul/service/deregister"
	"github.com/fuyibing/console/v3/commands/consul/service/register"
	"github.com/fuyibing/console/v3/commands/docs"
	"github.com/fuyibing/console/v3/commands/help"
	"github.com/fuyibing/console/v3/managers"
)

// Default
// function create and return default manager with built-in commands.
func Default() (mng managers.Manager, err error) {
	var (
		cmd managers.Command

		// Built-in command definitions.
		list = []func() (managers.Command, error){
			docs.New,
			download.New,
			upload.New,
			deregister.New,
			register.New,
		}
	)

	// Create and add
	// built-in commands to manager.
	if mng, err = New(); err == nil {
		for _, f := range list {
			// Return error
			// if create failed reason.
			if cmd, err = f(); err != nil {
				return
			}

			// Add to manager.
			if err = mng.AddCommand(cmd); err != nil {
				return
			}
		}
	}
	return
}

// Latest
// function create and return latest.
func Latest() (mng managers.Manager, err error) {
	return Default()
}

// New
// function create and return manager instance with help command.
func New() (managers.Manager, error) {
	mng := managers.NewManager()
	c, err := help.New()
	if err == nil {
		err = mng.AddCommand(c)
	}
	return mng, err
}
