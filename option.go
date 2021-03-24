// author: wsfuyibing <websearch@163.com>
// date: 2021-03-23

package console

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
	userValue    string
}

// New option instance.
func NewOption(name string) *Option {
	o := &Option{name: name}
	o.SetMode(RequiredMode).SetValue(StringValue)
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
	return o
}

// Validate option values of command line.
func (o *Option) validate(args ...string) error {
	return nil
}
