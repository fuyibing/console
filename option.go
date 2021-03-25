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
	FloatValue
	BooleanValue
	NullValue
)

// Option struct.
type Option struct {
	accepts      []string
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

// Set accepts.
func (o *Option) SetAccepts(accepts ...string) *Option {
	o.accepts = accepts
	return o
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
	if value == NullValue {
		o.mode = OptionalMode
		if o.defaultValue == nil {
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

// To float64.
func (o *Option) Float() float64 {
	s := o.String()
	if n, err := strconv.ParseFloat(s, 64); err == nil {
		return n
	}
	return 0
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
	// value parse.
	o.validateArguments(args...)
	// option not specified on command line.
	if !o.userFound {
		if o.mode == RequiredMode {
			return fmt.Errorf("option '%s' not specified", o.name)
		}
		return nil
	}
	// empty value.
	if o.userValue == "" {
		if o.mode == RequiredMode {
			return fmt.Errorf("value for '%s' option can not empty", o.name)
		}
		return nil
	}
	// integer value.
	if o.value == IntegerValue {
		_, err := strconv.ParseInt(o.userValue, 0, 64)
		if err != nil {
			return fmt.Errorf("value for '%s' option must be integer", o.name)
		}
		return nil
	}
	// float value.
	if o.value == FloatValue {
		_, err := strconv.ParseFloat(o.userValue, 64)
		if err != nil {
			return fmt.Errorf("value for '%s' option must be float", o.name)
		}
		return nil
	}
	// boolean value.
	if o.value == BooleanValue {
		_, err := strconv.ParseBool(o.userValue)
		if err != nil {
			return fmt.Errorf("value for '%s' option must be boolean", o.name)
		}
		return nil
	}
	// accepts check.
	if o.accepts != nil && len(o.accepts) > 0 {
		denied := true
		for _, acc := range o.accepts {
			if acc == o.userValue {
				denied = false
			}
		}
		if denied {
			return fmt.Errorf("option '%s' not accept '%s' value", o.name, o.userValue)
		}
	}

	return nil
}

// Validate arguments.
func (o *Option) validateArguments(args ...string) {
	// prepare.
	argc := len(args)
	index := -1
	// execute before return.
	defer func() {
		// not found.
		if !o.userFound || o.value == NullValue || o.userValue != "" {
			return
		}
		// match value.
		if (index + 1) < argc {
			if !RegexpOptionPrefix.MatchString(args[index+1]) {
				o.userValue = args[index+1]
			}
		}
	}()
	// Loop arguments and check option status.
	for i, arg := range args {
		// Single mode:
		//   -t
		//   -t value
		//   --tag value
		if m := RegexpOptionSingle.FindStringSubmatch(arg); len(m) == 3 {
			if m[1] == "-" && o.tag > 0 {
				for _, n := range m[2] {
					if n == int32(o.tag) {
						o.userFound = true
						index = i
						return
					}
				}
			} else if m[1] == "--" {
				if m[2] == o.name {
					o.userFound = true
					index = i
					return
				}
			}
		}
		// Double mode.
		//   --key="value"
		if m := RegexpOptionStandard.FindStringSubmatch(arg); len(m) == 3 {
			if m[1] == o.name {
				o.userFound = true
				o.userValue = m[2]
				return
			}
		}
	}
	return
}
