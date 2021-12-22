// author: wsfuyibing <websearch@163.com>
// date: 2021-12-20

package docs

import (
    "encoding/json"
    "fmt"
    "os"
    "reflect"
    "regexp"
    "strconv"
    "strings"
)

var (
    xreExported = regexp.MustCompile(`^([A-Z][_a-zA-Z0-9]*)`)
    xreRequired = regexp.MustCompile(`required`)
)

// Parser.

type (
    // XP
    // X Parser 接口.
    XP interface {
        Markdown() error
    }

    // X Parser.
    xp struct {
        basePath, name string
        err            error
        s              XS
    }
)

// X
// 构造XParser实例.
//
//   X(Struct{})
//   X(&Struct{})
func NewX(basePath, name string, v interface{}) XP {
    o := &xp{basePath: basePath, name: name}
    o.with(v)
    return o
}

// Markdown
// 生成Markdown文档.
func (o *xp) Markdown() error {
    // 1. 检查有错.
    if o.err != nil {
        return o.err
    }

    // 2. 准备导出.
    code, level, table := make(map[string]interface{}), 0, []string{
        fmt.Sprintf("| Field | Type | Required | Validate | Comment |"),
        fmt.Sprintf("| :---- | :---- | :----: | :----: | :---- |"),
    }
    if o.err = o.s.Markdown(code, &table, level); o.err != nil {
        return o.err
    }

    // 3. 文档内容.
    var body []byte
    if body, o.err = json.MarshalIndent(code, "", "    "); o.err != nil {
        return o.err
    }
    return o.save(fmt.Sprintf(
        "%s\n\n**Code**: \n\n```json\n%s\n```",
        strings.Join(table, "\n"),
        body,
    ))
}

// 保存文件.
func (o *xp) save(text string) error {
    var (
        err  error
        f    *os.File
        path = fmt.Sprintf("%s/%s", o.basePath, o.name)
    )

    if f, err = os.OpenFile(path, os.O_CREATE|os.O_RDWR|os.O_TRUNC, os.ModePerm); err != nil {
        return err
    }
    defer func() { _ = f.Close() }()

    if _, err = f.WriteString(text); err != nil {
        return err
    }
    return nil
}

// 使用入参.
func (o *xp) with(v interface{}) {
    r := reflect.ValueOf(v)
    k := r.Kind()

    if k == reflect.Struct {
        o.s, o.err = NewXS(r)
        return
    }

    if k == reflect.Ptr {
        if x := r.Elem(); x.Kind() == reflect.Struct {
            o.s, o.err = NewXS(x)
            return
        }
    }

    o.err = fmt.Errorf("invalid parse kind: %s", k)
    return
}

// Struct.

type (
    // XS
    // X Struct 接口.
    XS interface {
        Markdown(code map[string]interface{}, table *[]string, level int) error
    }

    // X Struct.
    xs struct {
        pkg, key string
        kind     reflect.Kind

        xss []XS
        xfs []XF
    }
)

// NewXS
// 构造Struct实例.
//
//   NewXS(
//       reflect.ValueOf(
//           Struct{},
//       ),
//   )
func NewXS(v reflect.Value) (XS, error) {
    o := &xs{pkg: v.Type().PkgPath(), kind: v.Kind()}
    o.xss = make([]XS, 0)
    o.xfs = make([]XF, 0)
    err := o.collect(v)
    return o, err
}

// Markdown
// 导出Markdown文档.
func (o *xs) Markdown(code map[string]interface{}, table *[]string, level int) error {
    // 1. 匿名.
    for _, a := range o.xss {
        if err := a.Markdown(code, table, level); err != nil {
            return err
        }
    }

    // 2. 字段.
    for _, f := range o.xfs {
        if err := f.Markdown(code, table, level); err != nil {
            return err
        }
    }

    return nil
}

// 收集字段.
func (o *xs) collect(v reflect.Value) error {
    for i := 0; i < v.NumField(); i++ {
        f := v.Field(i)
        s := v.Type().Field(i)

        // 匿名.
        if s.Anonymous {
            if err := o.collectAnonymous(v, f, s); err != nil {
                return err
            }
            continue
        }

        // 忽略.
        if !xreExported.MatchString(s.Name) {
            continue
        }

        // 字段.
        if err := o.collectField(v, f, s); err != nil {
            return err
        }
    }
    return nil
}

// 收集匿名字段.
func (o *xs) collectAnonymous(v, f reflect.Value, s reflect.StructField) (err error) {
    var (
        k = f.Kind()
        x XS
    )

    // 匿名: 结构体.
    if k == reflect.Struct {
        if x, err = NewXS(f); err == nil {
            o.xss = append(o.xss, x)
        }
        return
    }

    // 匿名: 指针.
    if k == reflect.Ptr {
        if n := reflect.New(f.Type().Elem()).Elem(); n.Kind() == reflect.Struct {
            if x, err = NewXS(n); err == nil {
                o.xss = append(o.xss, x)
            }
            return
        }
    }

    // 匿名: 无效.
    return fmt.Errorf("invalid anonymous kind: %s", k)
}

// 收集结构字段.
func (o *xs) collectField(v, f reflect.Value, s reflect.StructField) (err error) {
    var x XF
    if x, err = NewXF(v, f, s); err != nil {
        return
    }
    o.xfs = append(o.xfs, x)
    return
}

// Field.

type (
    // XF
    // X Field 接口.
    XF interface {
        Markdown(code map[string]interface{}, table *[]string, level int) error
    }

    // X Field.
    xf struct {
        pkg, key                                string
        comment, description, mock, name, title string
        validate, required                      string

        kind  reflect.Kind
        child XS

        vk     reflect.Kind
        vd, vm interface{}
        um     bool
    }
)

// NewXF
// 构造Field实例.
func NewXF(v, f reflect.Value, s reflect.StructField) (XF, error) {
    o := &xf{key: s.Name, name: s.Name, pkg: v.Type().PkgPath()}
    o.kind = f.Kind()
    o.initField(s)
    e := o.initType(f)
    if e == nil {
        o.initMock()
    }
    return o, e
}

func (o *xf) GetKey() string  { return o.key }
func (o *xf) GetName() string { return o.name }
func (o *xf) GetPkg() string  { return o.pkg }

// Markdown
// 合成Markdown结果.
func (o *xf) Markdown(res map[string]interface{}, table *[]string, level int) error {
    // 1. 加入表格.
    *table = append(
        *table, fmt.Sprintf(
            "| %s | %s | %s | %s | %s |",
            o.name,
            o.vk.String(),
            o.required,
            o.validate,
            o.comment,
        ),
    )

    if o.child != nil {
        c := make(map[string]interface{})
        if err := o.child.Markdown(c, table, level+1); err != nil {
            return err
        }
        if o.kind == reflect.Slice {
            res[o.name] = []interface{}{c}
        } else {
            res[o.name] = c
        }
    } else {
        v := o.vd
        if o.um {
            v = o.vm
        }
        if o.kind == reflect.Slice {
            res[o.name] = []interface{}{v}
        } else {
            res[o.name] = v
        }
    }
    return nil
}

// 初始: 字段参数.
func (o *xf) initField(s reflect.StructField) *xf {
    cs := make([]string, 0)

    // 1. tag:name
    for _, k := range []string{"url", "form", "json"} {
        if v := s.Tag.Get(k); v != "" {
            o.name = v
        }
    }

    // 2. tag:title
    for _, k := range []string{"label", "title"} {
        if v := s.Tag.Get(k); v != "" {
            o.title = v
            cs = append(cs, v)
        }
    }

    // 3. tag:description
    for _, k := range []string{"desc", "description"} {
        if v := s.Tag.Get(k); v != "" {
            o.description = v
            cs = append(cs, v)
        }
    }

    // 4. tag:mock
    for _, k := range []string{"mock"} {
        if v := s.Tag.Get(k); v != "" {
            o.mock = v
        }
    }

    // 5. tag:validate
    for _, k := range []string{"validate"} {
        if v := s.Tag.Get(k); v != "" {
            o.validate = v
            if xreRequired.MatchString(v) {
                o.required = "Y"
            }
        }
    }

    // 6. comment.
    if len(cs) > 0 {
        o.comment = strings.Join(cs, "<br />")
    }

    return o
}

// 初始: Mock数据.
func (o *xf) initMock() {
    if o.mock == "" {
        return
    }
    switch o.vk {
    case reflect.Bool:
        if b, e := strconv.ParseBool(o.mock); e == nil {
            o.um = true
            o.vm = b
        }
    case reflect.Int,
        reflect.Int8,
        reflect.Int16,
        reflect.Int32,
        reflect.Int64:
        if n, e := strconv.ParseInt(o.mock, 10, 64); e == nil {
            o.um = true
            o.vm = int64(n)
        }
    case reflect.Uint,
        reflect.Uint8,
        reflect.Uint16,
        reflect.Uint32,
        reflect.Uint64,
        reflect.Uintptr:
        if n, e := strconv.ParseUint(o.mock, 10, 64); e == nil {
            o.um = true
            o.vm = uint64(n)
        }
    case reflect.Float32,
        reflect.Float64,
        reflect.Complex64,
        reflect.Complex128:
        if n, e := strconv.ParseFloat(o.mock, 64); e == nil {
            o.um = true
            o.vm = float64(n)
        }
    case reflect.Interface,
        reflect.String:
        o.um = true
        o.vm = o.mock
    }
}

// 初始: 字段类型.
func (o *xf) initType(f reflect.Value) error {
    k := f.Kind()

    // struct {
    //     Key1 Struct
    // }
    if k == reflect.Struct {
        return o.withStruct(f)
    }

    // struct {
    //     Key1 *int
    //     Key2 *Struct
    //     Key3 *map[string]interface{}
    //     Key4 *[]Struct
    //     Key4 *[]*Struct
    // }
    if k == reflect.Ptr {
        return o.initType(reflect.New(f.Type().Elem()).Elem())
    }

    // struct {
    //     Key1 []int
    //     Key2 []*int
    //     Key3 []Struct
    //     Key4 []*Struct
    // }
    if k == reflect.Slice {
        e := f.Type().Elem()
        n := reflect.New(e).Elem()
        if e.Kind() == reflect.Struct {
            return o.withStruct(n)
        }
        return o.initType(n)
    }

    if k == reflect.Map {
        return nil
    }

    switch k {
    case reflect.Bool:
        o.vd = false
    case reflect.Int,
        reflect.Int8,
        reflect.Int16,
        reflect.Int32,
        reflect.Int64,
        reflect.Uint,
        reflect.Uint8,
        reflect.Uint16,
        reflect.Uint32,
        reflect.Uint64,
        reflect.Uintptr:
        o.vd = 0
    case reflect.Float32,
        reflect.Float64,
        reflect.Complex64,
        reflect.Complex128:
        o.vd = 0.0
    case reflect.Interface,
        reflect.String:
        o.vd = ""
    default:
        return fmt.Errorf("invalid field kind: %s", o.kind)
    }

    o.vk = k
    return nil
}

// 基于: 结构体.
func (o *xf) withStruct(f reflect.Value) (err error) {
    o.child, err = NewXS(f)
    return
}
