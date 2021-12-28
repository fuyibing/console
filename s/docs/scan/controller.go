// author: wsfuyibing <websearch@163.com>
// date: 2021-12-26

package scan

type (
    // Controller
    // 控制器接口.
    Controller interface {
        Add(a Action)

        // Each
        // 遍历API.
        Each(fn func(action Action))

        GetName() string

        // GetRoutePrefix
        // 读取路由前缀.
        //
        // return ""
        // return "/user"
        GetRoutePrefix() string

        GetScanner() Scanner
        GetSource() (string, int)

        GetTitle() string

        SetCommentBlock(cb CommentBlock) Controller
        SetName(name string) Controller
        SetScanner(scanner Scanner) Controller
        SetSource(path string, line int) Controller
    }

    // 控制器结构体.
    controller struct {
        scanner Scanner
        actions []Action

        name, title string

        routePrefix string
        sourcePath  string
        sourceLine  int
    }
)

// NewController
// 构造控制器接口实例.
func NewController() Controller {
    o := &controller{}
    o.actions = make([]Action, 0)
    return o
}

func (o *controller) Add(a Action) {
    o.actions = append(o.actions, a)
}

func (o *controller) Each(fn func(action Action)) {
    if fn != nil {
        for _, a := range o.actions {
            fn(a)
        }
    }
}

func (o *controller) GetName() string { return o.name }

// GetRoutePrefix
// 读取路由前缀.
func (o *controller) GetRoutePrefix() string { return o.routePrefix }

func (o *controller) GetScanner() Scanner { return o.scanner }

func (o *controller) GetSource() (string, int) { return o.sourcePath, o.sourceLine }

func (o *controller) GetTitle() string {
    if o.title != "" {
        return o.title
    }
    return o.name
}

func (o *controller) SetScanner(scanner Scanner) Controller {
    o.scanner = scanner
    return o
}

func (o *controller) SetCommentBlock(cb CommentBlock) Controller {
    // 1. 标题.
    if ti := cb.GetTitle(); ti != "" {
        o.title = ti
    }

    // 2. 注解.
    for k, vs := range cb.GetAnnotations() {
        switch k {
        case "routeprefix":
            {
                if len(vs) > 0 && vs[0] != "" {
                    o.routePrefix = vs[0]
                }
            }
        }
    }

    return o
}

func (o *controller) SetName(name string) Controller {
    o.name = name
    return o
}

func (o *controller) SetSource(path string, line int) Controller {
    o.sourcePath = path
    o.sourceLine = line
    return o
}
