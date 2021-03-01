// author: wsfuyibing <websearch@163.com>
// date: 2021-02-25

package i

// Option interface.
type IOption interface {
	IsIntValue() bool
	IsBoolValue() bool
	IsNoneValue() bool
	IsStrValue() bool
	IsOptional() bool
	IsRequired() bool
	GetDefaultValue() string
	GetDescription() string
	GetName() string
	GetShortName() string
	SetDescription(description string) IOption
	SetDefaultValue(defaultValue interface{}) IOption
	SetValue(value string) IOption
	SetName(name string) IOption
	SetShortName(name string) IOption
	ToBool() bool
	ToInt() int
	ToString() string
}
