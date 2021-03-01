// author: wsfuyibing <websearch@163.com>
// date: 2021-02-25

package i

import (
	"regexp"
	"sync"
)

type Mode int

const (
	OptionalMode Mode = iota
	RequiredMode
)

type ValueMode int

const (
	NoneValue ValueMode = iota
	BoolValue
	IntValue
	StrValue
)

const (
	UsageDefaultCommand = "COMMAND"
	UsageDefaultOption  = "[OPTIONS]"
)

var (
	RegexpName            = regexp.MustCompile(`^([a-zA-Z][:a-zA-Z0-9\-]*)$`)
	RegexpOptionName      = regexp.MustCompile(`^--([a-zA-Z][_a-zA-Z0-9\-]*)[=]?(.*)$`)
	RegexpOptionShortName = regexp.MustCompile(`^-([a-zA-Z0-9]+)$`)
	Script                = "go run main.go"
)

func init() {
	new(sync.Once).Do(func() {})
}
