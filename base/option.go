// author: wsfuyibing <websearch@163.com>
// date: 2021-02-26

package base

import (
	"fmt"
	"strconv"

	"github.com/fuyibing/console/v2/i"
)

type option struct {
	mode         i.Mode
	valueMode    i.ValueMode
	name         string
	shortName    string
	description  string
	value        string
	defaultValue string
}

func (o *option) IsIntValue() bool        { return o.valueMode == i.IntValue }
func (o *option) IsBoolValue() bool       { return o.valueMode == i.BoolValue }
func (o *option) IsNoneValue() bool       { return o.valueMode == i.NoneValue }
func (o *option) IsStrValue() bool        { return o.valueMode == i.StrValue }
func (o *option) IsOptional() bool        { return o.mode == i.OptionalMode }
func (o *option) IsRequired() bool        { return o.mode == i.RequiredMode }
func (o *option) GetDefaultValue() string { return o.defaultValue }
func (o *option) GetDescription() string  { return o.description }
func (o *option) GetName() string         { return o.name }
func (o *option) GetShortName() string    { return o.shortName }

// Set option description.
// Print when Usage() method call.
func (o *option) SetDefaultValue(defaultValue interface{}) i.IOption {
	o.defaultValue = fmt.Sprintf("%v", defaultValue)
	return o
}

// Set option description.
// Print when Usage() method call.
func (o *option) SetDescription(description string) i.IOption {
	o.description = description
	return o
}

// Set option short name.
//   `--name=value`
func (o *option) SetName(name string) i.IOption {
	o.name = name
	return o
}

// Set option short name.
//   `-n value`
func (o *option) SetShortName(name string) i.IOption {
	o.shortName = name
	return o
}

// Set option short name.
func (o *option) SetValue(value string) i.IOption {
	o.value = value
	return o
}

// Convert to boolean value.
func (o *option) ToBool() bool {
	var v = o.value
	if v == "" {
		if v = o.defaultValue; v == "" {
			return false
		}
	}
	b, _ := strconv.ParseBool(v)
	return b
}

// Convert to integer.
func (o *option) ToInt() int {
	v := o.value
	if v == "" {
		v = o.defaultValue
	}
	if v == "" {
		return 0
	}
	n, _ := strconv.ParseInt(v, 0, 32)
	return int(n)
}

// Convert to string.
func (o *option) ToString() string {
	if v := o.value; v != "" {
		return v
	}
	return o.defaultValue
}

// New option.
func NewOption(m i.Mode, vm i.ValueMode) i.IOption {
	o := &option{mode: m, valueMode: vm}
	if o.IsNoneValue() || o.IsBoolValue() {
		o.defaultValue = "false"
	}
	return o
}
