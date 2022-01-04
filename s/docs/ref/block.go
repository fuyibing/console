// author: wsfuyibing <websearch@163.com>
// date: 2021-12-27

package ref

import (
    "encoding/json"
    "fmt"
    "os"
    "reflect"
    "regexp"
    "sort"
    "strings"
)

// 块结构体.
type block struct {
    fields map[string]Element
    keys   []string
    level  int
}

// NewBlock
// 构造块实例.
//
//   var (
//       s = pkg.Struct{}
//       r = reflect.ValueOf(s)
//       b = ref.NewBlock(0)
//   )
//
//   if err := b.Run(r); err != nil {
//       println("reflection parse error:", err.Error())
//       return
//   }
//
//   println("reflection parse completed")
func NewBlock(level int) Block {
    return &block{
        fields: make(map[string]Element),
        keys:   make([]string, 0),
        level:  level,
    }
}

// Code
// 处理代码片段.
func (o *block) Code(cs map[string]interface{}) error {
    for _, f := range o.fields {
        if err := f.Code(cs); err != nil {
            return err
        }
    }
    return nil
}

// Markdown
// 导出Markdown文档.
func (o *block) Markdown(basePath, sourcePath string) error {
    // 1. 代码.
    cs := make(map[string]interface{})
    if err := o.Code(cs); err != nil {
        return err
    }

    // 2. 入参.
    req := []string{
        fmt.Sprintf(`| 字段名 | 类型 | 必须 | 校验 | 备注 |`),
        fmt.Sprintf(`| :---- | :---- | :----: | :---- | :---- |`),
    }
    if err := o.Request(&req); err != nil {
        return err
    }

    // 3. 出参.
    res := []string{
        fmt.Sprintf(`| 字段名 | 类型 | 备注 |`),
        fmt.Sprintf(`| :---- | :---- | :---- |`),
    }
    if err := o.Response(&res); err != nil {
        return err
    }

    // 4. 保存.
    if err := o.saveCode(basePath, sourcePath, "0", cs); err != nil {
        return err
    }
    if err := o.saveTable(basePath, sourcePath, "1", req); err != nil {
        return err
    }
    if err := o.saveTable(basePath, sourcePath, "2", res); err != nil {
        return err
    }
    return nil
}

// Request
// 处理入参.
func (o *block) Request(s *[]string) error {
    a := o.keys
    sort.Strings(a)

    for _, k := range a {
        if f, ok := o.fields[k]; ok {
            if err := f.Request(s); err != nil {
                return err
            }
        }
    }

    return nil
}

// Response
// 处理出参.
func (o *block) Response(s *[]string) error {
    a := o.keys
    sort.Strings(a)

    for _, k := range a {
        if f, ok := o.fields[k]; ok {
            if err := f.Response(s); err != nil {
                return err
            }
        }
    }

    return nil
}

// Run
// 执行解析.
func (o *block) Run(x reflect.Value) error {
    return o.each(x, func(v reflect.Value, sf reflect.StructField) error {
        f := NewField(o.level)
        if err := f.Run(v, sf); err != nil {
            return err
        }
        k := f.SortKey()
        o.fields[k] = f
        o.keys = append(o.keys, k)
        return nil
    })
}

// 遍历字段.
func (o *block) each(x reflect.Value, callback func(v reflect.Value, sf reflect.StructField) error) error {
    if callback != nil {
        r := regexp.MustCompile(`^[A-Z]`)
        for i := 0; i < x.NumField(); i++ {
            v := x.Field(i)
            sf := x.Type().Field(i)
            // 匿名.
            if sf.Anonymous {
                k := v.Kind()
                if k == reflect.Struct {
                    if err := o.each(v, callback); err != nil {
                        return err
                    }
                } else if k == reflect.Ptr {
                    if err := o.each(v.Elem(), callback); err != nil {
                        return err
                    }
                } else {
                    return fmt.Errorf("invalid anonymous: %s", k.String())
                }
            }
            // 忽略.
            if !r.MatchString(sf.Name) {
                continue
            }
            // 回调.
            if err := callback(v, sf); err != nil {
                return err
            }
        }
    }
    return nil
}

// 保存文档.
func (o *block) save(basePath, sourcePath, suffix, text string) error {
    var (
        name, path = strings.ToLower(
            strings.ReplaceAll(strings.TrimPrefix(sourcePath, "/"), "/", "-") + "." + suffix),
            fmt.Sprintf("%s/.md", basePath)
    )

    if err := os.MkdirAll(path, os.ModePerm); err != nil {
        return err
    }

    f, err := os.OpenFile(fmt.Sprintf("%s/%s", path, name), os.O_CREATE|os.O_RDWR|os.O_TRUNC, os.ModePerm)
    if err != nil {
        return err
    }
    defer func() { _ = f.Close() }()

    if _, err = f.WriteString(text); err != nil {
        return err
    }

    return nil
}

func (o *block) saveCode(basePath, sourcePath, suffix string, cs map[string]interface{}) error {
    buf, err := json.MarshalIndent(cs, "", "    ")
    if err != nil {
        return err
    }
    return o.save(basePath, sourcePath, suffix, string(buf))
}

func (o *block) saveTable(basePath, sourcePath, suffix string, ts []string) error {
    return o.save(basePath, sourcePath, suffix, strings.Join(ts, "\n"))
}

func (o *block) save2(table []string, code map[string]interface{}, basePath, sourcePath, suffix string) error {
    // 1. 名称.
    var (
        name, path, text = strings.ToLower(
            strings.ReplaceAll(
                strings.TrimPrefix(sourcePath, "/"),
                "/",
                "-",
            ) + "." + suffix,
        ), fmt.Sprintf(
            "%s/.md", basePath,
        ), ""
    )

    // 2. 内容.
    if suffix == "0" {
        // JSON.
        // 用于在Postman中的出入参.
        buf, err := json.MarshalIndent(code, "", "    ")
        if err != nil {
            return err
        }
        text = string(buf)
    } else {
        text += strings.Join(table, "\n")
        // 1.1 代码片段.
        if len(code) > 0 {
            buf, err := json.MarshalIndent(code, "", "    ")
            if err != nil {
                return err
            }
            text += "\n\n**Example**:\n\n"
            text += "```json\n"
            text += string(buf) + "\n"
            text += "```"
        }
    }

    // 2. 目录.
    //    创建基础目录.
    if err := os.MkdirAll(path, os.ModePerm); err != nil {
        return err
    }

    // 3. 打开.
    f, err := os.OpenFile(fmt.Sprintf("%s/%s", path, name), os.O_CREATE|os.O_RDWR|os.O_TRUNC, os.ModePerm)
    if err != nil {
        return err
    }
    defer func() { _ = f.Close() }()

    // 4. 写入.
    if _, err = f.WriteString(text); err != nil {
        return err
    }

    return nil
}
