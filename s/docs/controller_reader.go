// author: wsfuyibing <websearch@163.com>
// date: 2021-12-16

package docs

import (
    "bufio"
    "fmt"
    "os"
    "regexp"
    "strconv"
    "strings"
)

type (
    // ControllerReader
    // 读取接口.
    ControllerReader interface {
    }

    // 读取结构体.
    controllerReader struct {
        directory      Directory
        folder, name   string
        controllerFile string
    }
)

// NewControllerReader
// 构造控制器文件读取过程.
//
// NewControllerReader(directory, "/topic", "controller.go")
func NewControllerReader(directory Directory, folder, name string) (ControllerReader, error) {
    o := &controllerReader{directory: directory, folder: folder, name: name}
    o.controllerFile = fmt.Sprintf("%s%s/%s", directory.Scanner().GetControllerPath(), folder, name)

    if err := o.run(); err != nil {
        return nil, err
    }
    return o, nil
}

// 解析: 方法/Action.
func (o *controllerReader) parseAction(line int, name, controlerName string, an map[string][]string, cs []string) error {
    c := o.directory.GetControllerWithInitialize(controlerName)
    a, err := NewAction(c, name, line)
    if err != nil {
        return err
    }

    defer c.Add(a)

    // 1. 注解.
    for k, vs := range an {

        switch k {
        case "ignore":
            if len(vs) > 0 {
                if b, be := strconv.ParseBool(vs[0]); be == nil {
                    a.SetIgnore(b)
                }
            } else {
                a.SetIgnore(true)
            }
        case "request", "input":
            if len(vs) > 0 {
                a.SetRequest(vs[0])
            }
        case "response", "output":
            if len(vs) > 0 {
                if len(vs) > 0 {
                    a.SetResponse(vs[0])
                }
            }
        case "sdk":
            if len(vs) > 0 {
                a.SetSdk(vs[0])
            }
        case "version":
            if len(vs) > 0 {
                a.SetVersion(vs[0])
            }
        }
    }

    // 2. 注释.
    title, description := o.parseRemark(name, cs)
    if title != "" {
        a.SetTitle(title)
    }
    if description != "" {
        a.SetDescription(description)
    }

    return nil
}

// 解析: 控制器/Controller.
func (o *controllerReader) parseController(line int, name string, an map[string][]string, cs []string) error {
    c := o.directory.GetControllerWithInitialize(name)
    c.With(o.controllerFile, line)

    // 1. 注解.
    for k, vs := range an {
        switch k {
        case "routeprefix":
            if len(vs) > 0 {
                c.SetRoutePrefix(vs[0])
            }
        }
    }

    // 2. 注释.
    title, description := o.parseRemark(name, cs)
    if title != "" {
        c.SetTitle(title)
    }
    if description != "" {
        c.SetDescription(description)
    }

    return nil
}

// 解析: 注解内容.
func (o *controllerReader) parseAnnatationValue(str string) string {
    for _, s := range []string{`"`, `''`} {
        str = strings.TrimPrefix(str, s)
        str = strings.TrimSuffix(str, s)
        str = strings.TrimSpace(str)
    }
    return str
}

func (o *controllerReader) parseRemark(prefix string, cs []string) (title, description string) {
    n := 0
    r := regexp.MustCompile(`\s*[.]+$`)
    for _, s := range cs {
        if n == 0 {
            if s = strings.TrimSpace(strings.TrimPrefix(s, prefix)); s == "" {
                continue
            }
            if s = strings.TrimSpace(r.ReplaceAllString(s, "")); s == "" {
                continue
            }
        }
        if n == 0 {
            title = s
        } else {
            description = fmt.Sprintf("%s%s", description, s)
        }
        n++
    }
    return
}

// 执行过程.
func (o *controllerReader) run() error {
    var (
        err  error
        file *os.File
        path = fmt.Sprintf("%s%s", o.directory.Scanner().GetBasePath(), o.controllerFile)
    )

    // 1. 打开文件.
    //    结束后关闭文件.
    if file, err = os.OpenFile(path, os.O_RDONLY, os.ModePerm); err != nil {
        return err
    }
    defer func() { _ = file.Close() }()

    // 2. 准备数据.
    var (
        an = make(map[string][]string)
        cs = make([]string, 0)
        ig = false
        r0 = regexp.MustCompile(`^type\s*\(`)
        r1 = regexp.MustCompile(`^[/]+\s*(.*)`)
        r2 = regexp.MustCompile(`@([a-zA-Z0-9]+)\s*\(([^)]*)\)`)
        r3 = regexp.MustCompile(`([_a-zA-Z0-9]*Controller)\s*struct\s*\{`)
        r4 = regexp.MustCompile(`func\s*\(\s*([_a-zA-Z0-9]*)\s*[*]*([_a-zA-Z0-9]*Controller)\)\s*([A-Z][_a-zA-Z0-9]*)\(\s*([_a-zA-Z0-9]*)\s*iris\.Context\s*\)`)

        reset = func() {
            an = make(map[string][]string)
            cs = make([]string, 0)
            ig = false
        }
    )

    // 3. 逐行检查.
    num := 0
    buf := bufio.NewScanner(file)
    for buf.Scan() {
        num++
        str := strings.TrimSpace(buf.Text())

        // 空行.
        //
        // 1. 清空备注.
        if str == "" || r0.MatchString(str) {
            reset()
            continue
        }

        // 注释.
        if m1 := r1.FindStringSubmatch(str); len(m1) == 2 {
            // 注解.
            if m2 := r2.FindStringSubmatch(m1[1]); len(m2) == 3 {
                ig = true
                mk := strings.ToLower(m2[1])

                if _, ok := an[mk]; !ok {
                    an[mk] = make([]string, 0)
                }

                if mv := o.parseAnnatationValue(m2[2]); mv != "" {
                    an[mk] = append(an[mk], mv)
                }

                continue
            }

            // 备注.
            if !ig {
                cs = append(cs, m1[1])
            }

            continue
        }

        // 语法行.
        if m3 := r3.FindStringSubmatch(str); len(m3) == 2 {
            // 控制器.
            // m[0] = type ExampleController struct
            // m[1] = ExampleController
            if err = o.parseController(num, m3[1], an, cs); err != nil {
                return err
            }
        } else if m4 := r4.FindStringSubmatch(str); len(m4) == 5 {
            // 方法名.
            // m[0] = func (x *Controller) Get(ctx iris.Context)
            // m[1] = x
            // m[2] = Controller
            // m[3] = Get
            // m[4] = ctx
            if err = o.parseAction(num, m4[3], m4[2], an, cs); err != nil {
                return err
            }
        }
        reset()
    }

    return nil
}
