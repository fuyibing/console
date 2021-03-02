// author: wsfuyibing <websearch@163.com>
// date: 2021-02-25

package console

import (
	"github.com/fuyibing/console/v2/base"
	"github.com/fuyibing/console/v2/i"
	"github.com/fuyibing/console/v2/s/build/model"
	"github.com/fuyibing/console/v2/s/build/path"
	"github.com/fuyibing/console/v2/s/build/service"
	"github.com/fuyibing/console/v2/s/consul/deregister"
	"github.com/fuyibing/console/v2/s/consul/download"
	"github.com/fuyibing/console/v2/s/consul/register"
	"github.com/fuyibing/console/v2/s/consul/upload"
	"github.com/fuyibing/console/v2/s/docs"
	"github.com/fuyibing/console/v2/s/help"
)

// Return default console.
func Default() i.IConsole {
	c := New()
	c.Add(help.New(), docs.New())
	c.Add(path.New(), model.New(), service.New())
	c.Add(download.New(), upload.New(), register.New(), deregister.New())
	return c
}

// Return new console.
func New() i.IConsole {
	return base.NewConsole()
}
