// author: wsfuyibing <websearch@163.com>
// date: 2021-12-26

package scan

import (
    "bufio"
    "fmt"
    "os"
    "regexp"
    "strings"
)

type (
    Reader interface {
        Run() error
    }

    reader struct {
        scanner   Scanner
        directory Directory
        path      string
    }
)

// NewReader
// 构造文件读取实例.
//
//   NewReader(s, d, "/doc.go")
//   NewReader(s, d, "/topic/controller.go")
func NewReader(scanner Scanner, directory Directory, path string) Reader {
    o := &reader{scanner: scanner, directory: directory, path: path}
    return o
}

func (o *reader) Run() (err error) {
    var (
        file *os.File
        path = fmt.Sprintf(
            "%s%s%s",
            o.scanner.GetBasePath(),
            o.scanner.GetControllerPath(),
            o.path,
        )
    )

    // 1. 打开文件.
    //    结束后关闭文件.
    if file, err = os.OpenFile(path, os.O_RDONLY, os.ModePerm); err != nil {
        return err
    }
    defer func() { _ = file.Close() }()

    // 2. 解析准备.
    var (
        cb CommentBlock

        // 类型定义.
        r0 = regexp.MustCompile(`^type\s*\(`)

        // 注释语句.
        r1 = regexp.MustCompile(`^[/]{2,}[\s]?(.*)`)

        // 包名定义.
        //
        //   r2[0] = package controllers
        //   r2[1] = controllers
        r2 = regexp.MustCompile(`^package\s+([_a-zA-Z0-9]+)`)

        // 控制器定义.
        //
        //   r3[0] = ExampleController struct
        //   r3[1] = ExampleController
        r3 = regexp.MustCompile(`([_a-zA-Z0-9]*Controller)\s*struct\s*{`)

        // API定义.
        //
        //   r4[0] = func (x *Controller) Get(ctx iris.Context)
        //   r4[1] = x
        //   r4[2] = Controller
        //   r4[3] = Get
        //   r4[4] = ctx
        r4 = regexp.MustCompile(`func\s*\(\s*([_a-zA-Z0-9]*)\s*[*]*([_a-zA-Z0-9]*Controller)\)\s*([A-Z][_a-zA-Z0-9]*)\(\s*([_a-zA-Z0-9]*)\s*iris\.Context\s*\)`)

        reset = func() {
            if cb != nil {
                cb = nil
            }
        }
    )

    // 3. 逐行检查.
    num := 0
    buf := bufio.NewScanner(file)
    for buf.Scan() {
        num++
        str := strings.TrimSpace(buf.Text())

        // 忽略.
        // 空行或类型定义.
        if str == "" || r0.MatchString(str) {
            reset()
            continue
        }

        // 注释.
        if m := r1.FindStringSubmatch(str); len(m) == 2 {
            if cb == nil {
                cb = NewComment()
            }
            cb.Add(m[1])
            continue
        }

        // 语法.
        if m2 := r2.FindStringSubmatch(str); len(m2) > 0 {
            if cb != nil {
                cb.SetPrefix(fmt.Sprintf("Package %s", m2[1]))
                o.scanner.SetCommentBlock(cb)
            }
        } else if m3 := r3.FindStringSubmatch(str); len(m3) > 0 {
            if cb == nil {
                cb = NewComment()
            }
            o.runController(cb, m3[1], num)
        } else if m4 := r4.FindStringSubmatch(str); len(m4) > 0 {
            if cb == nil {
                cb = NewComment()
            }
            o.runAction(cb, m4[2], m4[3], num)
        }

        // 重置.
        reset()
    }
    return nil
}

func (o *reader) key(name string) string {
    return strings.TrimPrefix(fmt.Sprintf("%s%s", o.scanner.GetControllerPath(),
        regexp.MustCompile(`/([_a-zA-Z0-9-.]+)$`).ReplaceAllString(o.path, ""),
    ), "/") + "." + name
}

func (o *reader) runAction(cb CommentBlock, controllerName, name string, line int) {
    k := o.key(controllerName)
    c := o.scanner.GetController(k)
    c.Add(NewAction(c, name).SetCommentBlock(cb.SetPrefix(name)).SetSource(o.path, line))
}

func (o *reader) runController(cb CommentBlock, name string, line int) {
    k := o.key(name)
    c := o.scanner.GetController(k)
    c.SetScanner(o.scanner).SetCommentBlock(cb.SetPrefix(name)).SetName(name).SetSource(o.path, line)
}
