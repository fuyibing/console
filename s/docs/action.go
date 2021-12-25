// author: wsfuyibing <websearch@163.com>
// date: 2021-12-16

package docs

import (
    "fmt"
    "regexp"
    "strings"
)

type (
    Action interface {
        GetDescription() string
        // GetFile
        // 读取文件名与行号.
        //
        // return "/app/controllers/index_controller.go", 100
        GetFile() (string, int)

        // GetMethod
        // 读取请求方式.
        GetMethod() string

        // GetRoute
        // 读取路由地址.
        //
        // return "/topic/batch"
        GetRoute() string
        GetRouteLink() string
        GetSdk() string
        GetTitle() string

        // Ignored
        // 是否被忽略.
        Ignored() bool

        // Markdown
        // 导出Markdown文件.
        Markdown() error

        // SetDescription
        // 设置描述.
        SetDescription(s string) Action
        SetIgnore(b bool) Action
        SetRequest(s string) Action
        SetResponse(s string) Action
        SetSdk(s string) Action

        // SetTitle
        // 设置标题.
        SetTitle(s string) Action

        // SetVersion
        // 设置版本号.
        SetVersion(s string) Action
    }

    action struct {
        controller                  Controller
        name                        string
        line                        int
        description, title, version string
        method, uri                 string
        request, response, sdk      string
        ignore                      bool
    }
)

// NewAction
// 构造Action实例.
//
//   NewAction(controller, "PostBatch")
func NewAction(controller Controller, name string, line int) (Action, error) {
    o := &action{controller: controller, name: name, line: line}
    o.title = name
    o.version = "0.0"
    o.ignore = false
    o.run()
    return o, nil
}

// GetFile
// 读取文件名与行号.
func (o *action) GetFile() (string, int) {
    file, _ := o.controller.GetFile()
    return file, o.line
}
func (o *action) GetDescription() string { return o.description }
func (o *action) GetMethod() string      { return o.method }
func (o *action) GetRoute() string {
    if s := fmt.Sprintf("%s%s", o.controller.GetRoutePrefix(), o.uri); s != "" {
        return s
    }
    return "/"
}
func (o *action) GetRouteLink() string {
    s := o.GetRoute()

    if !strings.HasSuffix(s, "/") {
        s += "/"
    }

    s += strings.ToLower(o.method) + ".md"
    return s
}
func (o *action) GetSdk() string   { return o.sdk }
func (o *action) GetTitle() string { return o.title }

func (o *action) Ignored() bool { return o.ignore }

func (o *action) Markdown() error {
    if err := o.render(); err != nil {
        return err
    }
    return nil
}

func (o *action) SetDescription(s string) Action { o.description = s; return o }
func (o *action) SetIgnore(b bool) Action        { o.ignore = b; return o }
func (o *action) SetRequest(s string) Action {
    o.controller.Scanner().AddPayload(s)
    o.request = s
    return o
}
func (o *action) SetResponse(s string) Action {
    o.controller.Scanner().AddPayload(s)
    o.response = s
    return o
}
func (o *action) SetSdk(s string) Action     { o.sdk = s; return o }
func (o *action) SetTitle(s string) Action   { o.title = s; return o }
func (o *action) SetVersion(s string) Action { o.version = s; return o }

// 渲染模板.
func (o *action) render() error {
    // 源码.
    // SourceFile & SourceLine.
    sf, sl := o.GetFile()

    // 变量.
    args := map[string]interface{}{
        "TITLE":         o.GetTitle(),
        "DESCRIPTION":   o.GetDescription(),
        "METHOD":        o.GetMethod(),
        "ROUTE":         o.GetRoute(),
        "REQUEST":       "",
        "RESPONSE":      "",
        "CALLABLE_NAME": o.controller.GetName(), "CALLABLE_FUNC": o.name,
        "SOURCE_FILE": sf, "SOURCE_LINE": sl,
        "VERSION": o.version,
    }

    // 入参.
    if o.request != "" {
        if p := o.controller.Scanner().GetPayload(o.request); p != nil {
            args["REQUEST"] = p.Markdown(false)
        }
    }

    // 出参.
    if o.response != "" {
        if p := o.controller.Scanner().GetPayload(o.response); p != nil {
            args["RESPONSE"] = p.Markdown(true)
        }
    }

    // 替换.
    text := o.controller.Scanner().Template(templateAction, args)

    // 存储.
    if o.controller.Scanner().IsSaveLocal() {
        path := fmt.Sprintf(
            "%s%s%s",
            o.controller.Scanner().GetBasePath(),
            o.controller.Scanner().GetDocsPath(),
            o.GetRouteLink(),
        )
        if err := o.controller.Scanner().Save(path, text); err != nil {
            return err
        }
    }

    // 上传.
    if uploadUrl := o.controller.Scanner().GetUploadUrl(); uploadUrl != "" {
        if err := o.controller.Scanner().Upload(o.GetRouteLink(), text); err != nil {
            return err
        }
    }

    return nil
}

func (o *action) run() {
    // 1. 请求方式.
    //    GET/POST/PUT/PATCH/HEAD/OPTIONS/DELETE
    str := regexp.MustCompile(`^([A-Z][a-z]+)`).ReplaceAllStringFunc(o.name, func(s string) string {
        o.method = strings.ToUpper(s)
        return ""
    })

    // 2. 请求路由.
    if str != "" {
        o.uri = regexp.MustCompile(`[A-Z]`).ReplaceAllStringFunc(str, func(s string) string {
            return fmt.Sprintf("/%s", strings.ToLower(s))
        })
    }
}
