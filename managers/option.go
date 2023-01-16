// author: wsfuyibing <websearch@163.com>
// date: 2023-01-11

package managers

import (
	"fmt"
	"strconv"
	"strings"
)

type (
	Mode      int
	ValueType int
)

const (
	ModeOptional Mode = iota
	ModeRequired
)

const (
	ValueTypeString ValueType = iota
	ValueTypeBoolean
	ValueTypeFloat
	ValueTypeInteger
	ValueTypeNull
)

var (
	ValueTypeText = map[ValueType]string{
		ValueTypeBoolean: "boolean",
		ValueTypeFloat:   "float",
		ValueTypeInteger: "integer",
		ValueTypeString:  "string",
	}
)

type (
	// Option
	// operation interface.
	Option interface {
		Assign(s string) error
		Assigned() bool
		GetDescription() string
		GetLabel() string
		GetName() string
		GetShortName() string
		SetDefault(v interface{}) Option
		SetDescription(ss ...string) Option
		SetMode(m Mode) Option
		SetShortName(b byte) Option
		SetValueType(vt ValueType) Option
		ToBool() (bool, error)
		ToFloat() (float64, error)
		ToInt() (int64, error)
		ToString() (string, error)
		Validate() error
	}

	option struct {
		Default         interface{}
		Descriptions    []string
		Label           string
		Mode            Mode
		Name, ShortName string
		Value           string
		ValueAssigned   bool
		ValueType       ValueType
	}
)

func NewOption(name string) Option {
	return (&option{
		Descriptions: make([]string, 0),
		Name:         name,
		Mode:         ModeOptional, ValueType: ValueTypeString,
	}).initLabel()
}

// /////////////////////////////////////////////////////////////
// Interface methods
// /////////////////////////////////////////////////////////////

func (o *option) Assign(s string) error              { return o.assign(s) }
func (o *option) Assigned() bool                     { return o.ValueAssigned }
func (o *option) GetDescription() string             { return o.getDescription() }
func (o *option) GetLabel() string                   { return o.Label }
func (o *option) GetName() string                    { return o.Name }
func (o *option) GetShortName() string               { return o.ShortName }
func (o *option) SetDefault(v interface{}) Option    { o.Default = v; return o }
func (o *option) SetDescription(ss ...string) Option { o.setDescription(ss...); return o }
func (o *option) SetMode(m Mode) Option              { o.Mode = m; return o.initLabel() }
func (o *option) SetShortName(b byte) Option         { o.ShortName = string(b); return o.initLabel() }
func (o *option) SetValueType(vt ValueType) Option   { o.ValueType = vt; return o.initLabel() }
func (o *option) ToBool() (bool, error)              { return o.toBool() }
func (o *option) ToFloat() (float64, error)          { return o.toFloat() }
func (o *option) ToInt() (int64, error)              { return o.toInt() }
func (o *option) ToString() (string, error)          { return o.toString() }
func (o *option) Validate() error                    { return o.validate() }

// /////////////////////////////////////////////////////////////
// Access and constructor
// /////////////////////////////////////////////////////////////

func (o *option) assign(s string) error {
	if o.Value = s; o.Value != "" && o.ValueType == ValueTypeNull {
		return fmt.Errorf("option not accept any value: %s", o.Name)
	}

	o.ValueAssigned = true
	return nil
}

func (o *option) getDescription() string {
	ss := o.Descriptions

	if o.Default != nil {
		ss = append(ss, fmt.Sprintf("(default: %v)", o.Default))
	}

	return strings.Join(ss, " ")
}

func (o *option) initLabel() *option {
	// Short name.
	if o.ShortName == "" {
		// No short name.
		o.Label = "    "
	} else {
		// Add short name.
		o.Label = fmt.Sprintf("-%s, ", o.ShortName)
	}

	// Add full name.
	o.Label = fmt.Sprintf("%s--%s", o.Label, o.Name)

	// Option type.
	if o.ValueType != ValueTypeNull {
		if s, ok := ValueTypeText[o.ValueType]; ok {
			if o.Mode == ModeOptional {
				o.Label = fmt.Sprintf("%s[=%s]", o.Label, s)
			} else {
				o.Label = fmt.Sprintf("%s=<%s>", o.Label, s)
			}
		}
	}

	return o
}

func (o *option) setDescription(ss ...string) {
	ds := make([]string, 0)
	for _, s := range ss {
		if s = strings.TrimSpace(s); s != "" {
			ds = append(ds, s)
		}
	}
	o.Descriptions = ds
}

func (o *option) toBool() (bool, error) {
	if o.ValueType != ValueTypeBoolean {
		return false, fmt.Errorf("option type not matched on boolean: %s", o.Name)
	}

	// Return
	// user value.
	if o.Value != "" {
		if v, err := strconv.ParseBool(o.Value); err == nil {
			return v, nil
		}
		return false, fmt.Errorf("option value convert to boolean failed: %s", o.Name)
	}

	// Return
	// default value.
	if o.Default != nil {
		if v, ok := o.Default.(bool); ok {
			return v, nil
		}
		return false, fmt.Errorf("option value convert to boolean failed: %s", o.Name)
	}

	// Return
	// system value.
	return false, nil
}

func (o *option) toFloat() (float64, error) {
	if o.ValueType != ValueTypeFloat {
		return 0, fmt.Errorf("option type not matched on float: %s", o.Name)
	}

	// Return
	// user value.
	if o.Value != "" {
		if v, err := strconv.ParseFloat(o.Value, 64); err == nil {
			return v, nil
		}
		return 0, fmt.Errorf("option value convert to float failed: %s", o.Name)
	}

	// Return
	// default value.
	if o.Default != nil {
		if v, ok := o.Default.(float64); ok {
			return v, nil
		}
		return 0, fmt.Errorf("option value convert to float failed: %s", o.Name)
	}

	// Return
	// system value.
	return 0, nil
}

func (o *option) toInt() (int64, error) {
	if o.ValueType != ValueTypeInteger {
		return 0, fmt.Errorf("option type not matched on integer: %s", o.Name)
	}

	// Return
	// user value.
	if o.Value != "" {
		if v, err := strconv.ParseInt(o.Value, 10, 64); err == nil {
			return v, nil
		}
		return 0, fmt.Errorf("option value convert to integer failed: %s", o.Name)
	}

	// Return
	// default value.
	if o.Default != nil {
		if v, ok := o.Default.(int64); ok {
			return v, nil
		}
		return 0, fmt.Errorf("option value convert to integer failed: %s", o.Name)
	}

	// Return
	// system value.
	return 0, nil
}

func (o *option) toString() (string, error) {
	if o.ValueType != ValueTypeString {
		return "", fmt.Errorf("option type not matched on string: %s", o.Name)
	}

	// Return
	// user value.
	if o.Value != "" {
		return o.Value, nil
	}

	// Return
	// default value.
	if o.Default != nil {
		if v, ok := o.Default.(string); ok {
			return v, nil
		}
		return "", fmt.Errorf("option value convert to string failed: %s", o.Name)
	}

	// Return
	// system value.
	return "", nil
}

func (o *option) validate() error {
	if o.Value == "" && o.Mode == ModeRequired {
		return fmt.Errorf("option is required: %s", o.Name)
	}
	return nil
}
