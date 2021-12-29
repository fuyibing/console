// author: wsfuyibing <websearch@163.com>
// date: 2021-12-27

package ref

import (
    "fmt"
    "reflect"
    "regexp"
    "strconv"
)

// 字段结构体.
type field struct {
    level int

    // 名称.
    //
    // name=Id
    // fieldName=id
    // fieldType=int
    name, fieldName, fieldType string

    // 控制.
    required, executable bool

    // 配置.
    //
    // title=标题
    // description=关于Field描述
    // mock=message
    // validate=required,min=1
    title, description, mock, validate string

    arr      int         // 是否数组
    child    Block       // 子级
    defaults interface{} // 默认值
}

// NewField
// 构造字段实例.
//
//   f := ref.NewField(0)
//
//   if err := f.Run(v, sf); err != nil {
//       println("field parse error:", err.Error())
//       return
//   }
//
//   println("field parse completed")
func NewField(level int) Field {
    return &field{level: level}
}

// Code
// 追加代码片段.
func (o *field) Code(cs map[string]interface{}) error {
    if o.child != nil {
        v := make(map[string]interface{})
        if err := o.child.Code(v); err != nil {
            return err
        }
        if o.arr > 0 {
            var b []interface{}
            for i := 0; i < o.arr; i++ {
                if i == 0 {
                    b = []interface{}{v}
                } else {
                    b = []interface{}{b}
                }
            }
            cs[o.fieldName] = b
        } else {
            cs[o.fieldName] = v
        }
    } else {
        if o.arr > 0 {
            var tmp []interface{}
            for i := 0; i < o.arr; i++ {
                if i == 0 {
                    tmp = []interface{}{o.defaults}
                } else {
                    tmp = []interface{}{tmp}
                }
            }
            cs[o.fieldName] = tmp
        } else {
            cs[o.fieldName] = o.defaults
        }
    }
    return nil
}

// Request
// 入参数据.
func (o *field) Request(s *[]string) (err error) {
    // 1. 行数据.
    *s = append(*s,
        fmt.Sprintf(
            `| %s | %s | %s | %s | %s |`,
            o.buildMarkdownFieldName(),
            o.buildMarkdownType(),
            o.buildMarkdownRequired(),
            o.validate,
            o.buildMarkdownRemark(),
        ),
    )

    // 2. 子数据.
    if o.child != nil {
        err = o.child.Request(s)
    }

    return
}

// Response
// 出参数据.
func (o *field) Response(s *[]string) (err error) {
    // 1. 行数据.
    *s = append(*s,
        fmt.Sprintf(
            `| %s | %s | %s |`,
            o.buildMarkdownFieldName(),
            o.buildMarkdownType(),
            o.buildMarkdownRemark(),
        ),
    )

    // 2. 子数据.
    if o.child != nil {
        err = o.child.Response(s)
    }

    return
}

// Run
// 执行解析.
func (o *field) Run(v reflect.Value, sf reflect.StructField) (err error) {
    if err = o.runStructField(sf); err == nil {
        err = o.runReflectValue(v)
    }
    return err
}

// SortKey
// 排序键名.
func (o *field) SortKey() string {
    a := "1"
    if o.required {
        a = "0"
    }
    return fmt.Sprintf("%s-%s", a, o.fieldName)
}

// MD: 字段名.
func (o *field) buildMarkdownFieldName() string {
    s := ""
    if o.level > 0 {
        for i := 0; i < (o.level - 1); i++ {
            s += "　　"
        }
        s += "　└ "
    }
    return s + o.fieldName
}

// MD: 字段描述.
func (o *field) buildMarkdownRemark() string {
    str := ""
    if o.title != "" {
        str += fmt.Sprintf(`%s`, o.title)
    }
    if o.description != "" {
        if str != "" {
            str += "<br />"
        }
        str += regexp.MustCompile(`[\n]+`).ReplaceAllString(o.description, "<br />")
    }
    return str
}

func (o *field) buildMarkdownRequired() string {
    if o.required {
        return "√"
    }
    return " "
}

// MD: 字段类型.
func (o *field) buildMarkdownType() string {
    str := o.fieldType
    for i := 0; i < o.arr; i++ {
        str = fmt.Sprintf(`[]%s`, str)
    }
    return str
}

// 解析: 字段类型.
func (o *field) runReflectValue(v reflect.Value) (err error) {
    k := v.Kind()

    // 1. 任意值.
    //
    //    {
    //        X interface{}
    //    }
    if k == reflect.Interface {
        o.fieldType = "ANY"
        o.defaults = "*"
        return
    }

    // 2. 集合值.
    if k == reflect.Map {
        o.fieldType = "JSON"
        o.defaults = make([]interface{}, 0)
        return
    }

    // 3. 结构体.
    //
    //    struct {
    //        X1 Struct
    //    }
    if k == reflect.Struct {
        err = o.runStruct(v)
        return
    }

    // 4. 指针值.
    //
    //    struct {
    //        X1 *Struct
    //        X2 *[]Struct
    //        X3 *[]*Struct
    //    }
    if k == reflect.Ptr {
        return o.runReflectValue(reflect.New(v.Type().Elem()).Elem())
    }

    // 5. 切片值.
    if k == reflect.Slice {
        o.arr++
        return o.runReflectValue(reflect.New(v.Type().Elem()).Elem())
    }

    // 6. 系统值.
    switch k {
    // String/字符串.
    case reflect.String:
        {
            o.fieldType = k.String()
            o.defaults = o.mock
            return
        }

    // Boolean类型.
    case reflect.Bool:
        {
            o.fieldType = k.String()
            o.defaults = false
            if o.mock != "" {
                if p, pe := strconv.ParseBool(o.mock); pe == nil {
                    o.defaults = p
                }
            }
            return
        }

    // Integer/整形.
    case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
        {
            o.fieldType = k.String()
            o.defaults = 0
            if o.mock != "" {
                if p, pe := strconv.ParseUint(o.mock, 10, 64); pe == nil {
                    o.defaults = p
                }
            }
            return
        }

    // Float/浮点.
    case reflect.Float32, reflect.Float64, reflect.Complex64, reflect.Complex128:
        {
            o.fieldType = k.String()
            o.defaults = 0.0
            if o.mock != "" {
                if p, pe := strconv.ParseFloat(o.mock, 64); pe == nil {
                    o.defaults = p
                }
            }
            return
        }
    }

    // n. 忽略类型.
    //
    //   reflect.Uintptr
    //   reflect.Array
    //   reflect.Chan
    //   reflect.Func
    //   reflect.UnsafePointer
    return fmt.Errorf("disabled field type: %s", k.String())
}

// 解析: 结构体字段.
func (o *field) runStruct(v reflect.Value) error {
    o.fieldType = "JSON"
    o.child = NewBlock(o.level + 1)
    return o.child.Run(v)
}

// 解析: 字段属性.
func (o *field) runStructField(sf reflect.StructField) error {
    // 1. 准备.
    o.name = sf.Name
    o.fieldName = sf.Name
    o.title = sf.Name

    // 2. 名称.
    for _, k := range []string{"json", "yaml", "xml", "form", "url"} {
        if v := sf.Tag.Get(k); v != "" {
            o.fieldName = v
        }
    }

    // 3. 标题.
    for _, k := range []string{"label", "title"} {
        if v := sf.Tag.Get(k); v != "" {
            o.title = v
        }
    }

    // 4. 描述.
    for _, k := range []string{"desc", "description", "remark"} {
        if v := sf.Tag.Get(k); v != "" {
            o.description = v
        }
    }

    // 5. 校验.
    for _, k := range []string{"validate"} {
        if v := sf.Tag.Get(k); v != "" {
            o.validate = v
            o.required = regexp.MustCompile(`required`).MatchString(v)
        }
    }

    // 6. MOCK.
    for _, k := range []string{"mock"} {
        if v := sf.Tag.Get(k); v != "" {
            o.mock = v
        }
    }

    // 7. 计算属性.
    if v := sf.Tag.Get("exec"); v != "" {
        if b, be := strconv.ParseBool(v); be == nil {
            o.executable = b
        }
    }

    return nil
}
