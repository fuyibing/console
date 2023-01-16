// author: wsfuyibing <websearch@163.com>
// date: 2023-01-11

package main

import (
	"github.com/fuyibing/console/v3"
	"github.com/fuyibing/console/v3/managers"
)

var (
	manager managers.Manager
	err     error
)

func init() {
	if manager, err = console.Default(); err == nil {
		manager.SetDescription("About example")
	}
}

func main() {
	if err == nil {
		err = manager.RunTerminal()
	}
	if err != nil {
		println(err.Error())
	}
}
