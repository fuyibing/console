# Console

> Package `github.com/fuyibing/console`, Console for golang application. 


```text
// author: wsfuyibing <websearch@163.com>
// date: 2021-03-16

package b

import (
	"fmt"
	"strconv"
)

type Mode int

const (
	UnknownMode Mode = iota
	OptionalMode
	RequiredMode
)

type Value int

const (
	UnknownValue Value = iota
	NullValue    Value = iota
	BooleanValue
	IntegerValue
	StringValue
)

// 选项结构体.
type Option struct {
	Mode         Mode
	Value        Value
	_default     interface{}
	_description string
	_name        string
	_tag         int32
	_found       bool
	_value       string
}

// 创建选项.
func NewOption() *Option { return &Option{Mode: UnknownMode, Value: UnknownValue} }

func (o *Option) AsMode(Mode Mode) *Option         { o.Mode = Mode; return o }
func (o *Option) AsValue(Value Value) *Option      { o.Value = Value; return o }
func (o *Option) GetDescription() string           { return o._description }
func (o *Option) GetName() string                  { return o._name }
func (o *Option) GetTag() byte                     { return byte(o._tag) }
func (o *Option) IsOptional() bool                 { return o.Mode == OptionalMode }
func (o *Option) IsRequired() bool                 { return o.Mode == RequiredMode }
func (o *Option) IsNullValue() bool                { return o.Value == NullValue }
func (o *Option) IsBooleanValue() bool             { return o.Value == BooleanValue }
func (o *Option) IsIntegerValue() bool             { return o.Value == IntegerValue }
func (o *Option) IsStringValue() bool              { return o.Value == StringValue }
func (o *Option) SetDefault(v interface{}) *Option { o._default = v; return o }
func (o *Option) SetDescription(v string) *Option  { o._description = v; return o }
func (o *Option) SetName(v string) *Option         { o._name = v; return o }
func (o *Option) SetTag(v byte) *Option            { o._tag = int32(v); return o }

// 填充选项值.
func (o *Option) Fill() string {
	// 空值.
	if o.Value == NullValue {
		return "[=false]"
	}
	// 填充.
	v := ""
	if o._default != nil {
		v = fmt.Sprintf("%v", o._default)
		return fmt.Sprintf("=\"%s\"", v)
	}
	switch o.Value {
	case BooleanValue:
		v = "boolean"
	case IntegerValue:
		v = "integer"
	case StringValue:
		v = "string"
	}
	if v != "" {
		if o.Mode == RequiredMode {
			return "=<" + v + ">"
		}
		return "[=" + v + "]"
	}
	return ""
}

// 校验入参.
func (o *Option) Validate(args ...string) error {
	// 解析项选.
	o.queryArgument(args...)
	// 必须项: 未定义时报错.
	if o.Mode == RequiredMode && !o._found {
		return fmt.Errorf("required option not speicified: %s", o._name)
	}
	// 空值项.
	if o.Value == NullValue {
		return nil
	}
	// 布尔项.
	if o.Value == BooleanValue && o._value != "" {
		res, err := strconv.ParseBool(o._value)
		if err != nil {
			return fmt.Errorf("invalid boolean value option: %s", o._name)
		}
		if res {
			o._value = "true"
		} else {
			o._value = "false"
		}
	} else if o.Value == IntegerValue && o._value != "" {
		_, err := strconv.ParseInt(o._value, 0, 64)
		if err != nil {
			return fmt.Errorf("invalid integer value option: %s", o._name)
		}
	}
	return nil
}

// 读取运行布尔值.
func (o *Option) Bool() bool {
	// 空值布尔.
	if o.Value == NullValue {
		return o._found
	}
	// 布尔类型.
	if o.Value == BooleanValue {
		return o._value == "true"
	}
	// 其它类型.
	return false
}

// 读取运行整型值.
func (o *Option) Int() int64 {
	if v := o.String(); v != "" {
		if n, e := strconv.ParseInt(v, 0, 64); e != nil {
			return n
		}
	}
	return 0
}

// 读取运行字符串.
func (o *Option) String() string {
	if o._value != "" {
		return o._value
	}
	if o._default != nil {
		return fmt.Sprintf("%v", o._default)
	}
	return ""
}

// 读取参数.
//   script command -a
//   script command -a value
//   script command --arg
//   script command --arg value
//   script command --arg=value
func (o *Option) queryArgument(args ...string) {
	for k, arg := range args {
		// Single.
		if m := RegexpEqSingle.FindStringSubmatch(arg); len(m) == 3 {
			if m[1] == "-" {
				if o._tag == 0 {
					return
				}
				for _, n := range m[2] {
					if n == o._tag {
						o.queryFound(k, "", args...)
						return
					}
				}
			} else if m[1] == "--" {
				if m[2] == o._name {
					o.queryFound(k, "", args...)
					return
				}
			}
		}
		// Double.
		if m := RegexpEqDouble.FindStringSubmatch(arg); len(m) == 4 {
			if m[2] == o._name {
				o.queryFound(k, m[3])
				return
			}
		}
	}
}

// 发现选项.
func (o *Option) queryFound(k int, v string, args ...string) {
	o._found = true
	// 空值.
	if o.Value == NullValue {
		return
	}
	// 赋值.
	if v == "" {
		if l := len(args) - 1; k < l {
			if v = args[k+1]; !RegexpIsOption.MatchString(v) {
				o._value = v
			}
		}
	} else {
		o._value = v
	}
}

```