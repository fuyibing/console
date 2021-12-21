// author: wsfuyibing <websearch@163.com>
// date: 2021-12-16

package docs

type (
    Controller interface {
        // Add
        // 添加Action.
        Add(a Action) Controller

        // Each
        // 遍历Action.
        EachAction(callback func(action Action))

        // GetDescription
        // 读取文件描述.
        GetDescription() string

        // GetFile
        // 读取文件.
        GetFile() (path string, line int)

        // GetName
        // 读取控制器名称.
        //
        // return "ExampleController"
        GetName() string

        // GetRoutePrefix
        // 读取路由前缀.
        //
        // return "/topic"
        GetRoutePrefix() string

        // GetTitle
        // 读取控制器标题.
        GetTitle() string

        Scanner() Scanner

        // SetDescription
        // 设置标题.
        //
        // SetDescription("About controller description messages")
        SetDescription(s string) Controller

        SetRoutePrefix(s string) Controller

        // SetTitle
        // 设置标题.
        //
        // SetTitle("Example controller")
        SetTitle(s string) Controller

        // SetVersion
        // 设置版本号.
        //
        // SetVersion("1.2.3")
        SetVersion(s string) Controller

        // With
        // 绑定文件名与行号.
        //
        // With("example_controller.go", 123)
        With(file string, line int)
    }

    controller struct {
        actions                 []Action
        directory               Directory
        file, name, routePrefix string
        line                    int

        description, title, version string
    }
)

// NewController
// 构造控制器实例.
func NewController(directory Directory, name string) Controller {
    o := &controller{directory: directory, name: name}
    o.actions = make([]Action, 0)
    o.routePrefix = directory.GetFolder()

    o.title = name
    o.version = "0.0"
    return o
}

func (o *controller) Add(a Action) Controller { o.actions = append(o.actions, a); return o }

func (o *controller) EachAction(callback func(action Action)) {
    for _, a := range o.actions {
        callback(a)
    }
}

func (o *controller) GetDescription() string           { return o.description }
func (o *controller) GetFile() (path string, line int) { return o.file, o.line }
func (o *controller) GetName() string                  { return o.name }
func (o *controller) GetRoutePrefix() string           { return o.routePrefix }
func (o *controller) GetTitle() string                 { return o.title }

func (o *controller) Scanner() Scanner { return o.directory.Scanner() }

func (o *controller) SetDescription(s string) Controller { o.description = s; return o }
func (o *controller) SetRoutePrefix(s string) Controller { o.routePrefix = s; return o }
func (o *controller) SetTitle(s string) Controller       { o.title = s; return o }
func (o *controller) SetVersion(s string) Controller     { o.version = s; return o }

func (o *controller) With(file string, line int) {
    o.file = file
    o.line = line
}
