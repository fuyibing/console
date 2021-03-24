// author: wsfuyibing <websearch@163.com>
// date: 2021-02-25

// Package console.
package console

import (
	"regexp"
)

const (
	PidFile = "pid"
)

var (
	RegexpArgumentCommandName = regexp.MustCompile(`^([a-zA-Z][_a-zA-Z0-9-]*)`)
)
