// author: wsfuyibing <websearch@163.com>
// date: 2021-03-02

package docs

import (
	"github.com/fuyibing/console/v2/base"
	"github.com/fuyibing/console/v2/i"
)

const (
	Description = "Export application document"
	Name        = "docs"
)

// Command struct.
type command struct {
	base.Command
}

// Return export document instance.
func New() *command {
	// normal.
	o := &command{}
	o.Initialize()
	o.SetName(Name)
	o.SetDescription(Description)
	return o
}

func Run(console i.IConsole) {

}
