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
}

func New() i.ICommand {
	o := &command{}
	o.Initialize()
	o.SetDescription(Description)
	o.SetName(Name)
	return o
}
