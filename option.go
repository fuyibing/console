// author: wsfuyibing <websearch@163.com>
// date: 2021-03-23

package console

import (
	"fmt"
	"strconv"
)

type Mode int

const (
	RequiredMode Mode = iota
	OptionalMode
)

type Value int

const (
	StringValue Value = iota
	IntegerValue
	BooleanValue
	NullValue
)

// Option struct.
type Option struct {
	name         string
	description  string
	mode         Mode
	value        Value
	tag          byte
	defaultValue interface{}
	userFound    bool
	userValue    string
}

// New option instance.
func NewOption(name string) *Option {
	return &Option{
		name:  name,
		mode:  RequiredMode,
		value: StringValue,
	}
}

// Set default value.
func (o *Option) SetDefaultValue(v interface{}) *Option {
	o.defaultValue = v
	return o
}

// Set option description.
func (o *Option) SetDescription(description string) *Option {
	o.description = description
	return o
}

// Set option mode.
func (o *Option) SetMode(mode Mode) *Option {
	o.mode = mode
	return o
}

// Set tag.
func (o *Option) SetTag(tag byte) *Option {
	o.tag = tag
	return o
}

// Set option value mode.
func (o *Option) SetValue(value Value) *Option {
	o.value = value
	if o.defaultValue == nil {
		if value == NullValue {
			o.defaultValue = "false"
		}
	}
	return o
}

// To boolean.
func (o *Option) Bool() bool {
	// Null value.
	// Return true if specified on command line, else false.
	if o.value == NullValue {
		return o.userFound
	}
	// parse string value to boolean.
	if s := o.String(); s != "" {
		if b, err := strconv.ParseBool(s); err == nil {
			return b
		}
	}
	return false
}

// To boolean.
func (o *Option) Integer() int {
	s := o.String()
	if n, err := strconv.ParseInt(s, 0, 32); err == nil {
		return int(n)
	}
	return 0
}

// To string.
func (o *Option) String() string {
	if o.userValue != "" {
		return o.userValue
	}
	if o.defaultValue != nil {
		return fmt.Sprintf("%v", o.defaultValue)
	}
	return ""
}

// Validate option values of command line.
func (o *Option) validate(args ...string) error {
	return nil
}
