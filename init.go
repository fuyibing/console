// author: wsfuyibing <websearch@163.com>
// date: 2021-02-25

// Package console.
package console

import (
	"regexp"
	"time"
)

const (
	PidFile = "pid"
)

var (
	RegexpArgumentCommandName = regexp.MustCompile(`^([a-zA-Z][_a-zA-Z0-9-]*)`)
	RegexpOptionPrefix        = regexp.MustCompile(`^-`)
	RegexpOptionSingle        = regexp.MustCompile(`^([-]{1,2})([a-zA-Z0-9][_a-zA-Z0-9-]*)$`)
	RegexpOptionStandard      = regexp.MustCompile(`^--([a-zA-Z0-9][_a-zA-Z0-9-]*)=(.*)$`)
	RegexpUnderlinePrefix     = regexp.MustCompile(`^[_]+`)
	RegexpUnderlineMiddle     = regexp.MustCompile(`[_]+([a-zA-Z0-9])`)
	RegexpFirstLetter         = regexp.MustCompile(`^([a-zA-Z])`)
	RegexpUpperLetter         = regexp.MustCompile(`([A-Z])`)
	RegexpColumnType          = regexp.MustCompile(`^([a-zA-Z0-9]+)`)
)

// Column.
type Column struct {
	Comment string `xorm:"Comment"`
	Key     string `xorm:"Comment"`
	Name    string `xorm:"Field"`
	Type    string `xorm:"Comment"`
}

// Table meta.
type Table struct {
	Comment string    `xorm:"Comment"`
	Created time.Time `xorm:"Create_time"`
	Name    string    `xorm:"Name"`
}
