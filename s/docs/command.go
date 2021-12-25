// author: wsfuyibing <websearch@163.com>
// date: 2021-03-02

package docs

import (
    "fmt"
    "os"

    "github.com/fuyibing/console/v2/base"
    "github.com/fuyibing/console/v2/i"
)

const (
    Description = "Export application document"
    Name        = "docs"
)

// 导出文档.
type command struct {
    base.Command

    basePath, controllerPath, docsPath string
    scanner                            Scanner
}

// 构造导出文档实例.
func New() *command {
    o := &command{basePath: ".", controllerPath: "/app/controllers", docsPath: "/docs/api"}
    o.Initialize()
    o.SetName(Name)
    o.SetDescription(Description)
    // 执行过程.
    o.SetHandlerBefore(o.before)
    o.SetHandler(o.run)
    o.SetHandlerAfter(o.after)

    // 上传地址.
    // 生成的Markdown文件上传到哪里.
    //   -u gs-docs.turboradio.cn
    //   --upload=gs-docs.turboradio.cn
    o.Add(base.NewOption(i.OptionalMode, i.StrValue).SetName("upload").SetShortName("u").SetDescription("upload markdown to specified server"))

    // 保存文件.
    //   --parse
    o.Add(base.NewOption(i.OptionalMode, i.BoolValue).SetName("save").SetDefaultValue("true").SetDescription("save markdown to local"))
    return o
}

// 后置.
func (o *command) after(c i.IConsole) {}

// 前置.
func (o *command) before(c i.IConsole) bool {
    // 1. 根目录.
    if g := o.GetOption("base-path"); g != nil {
        if s := g.ToString(); s != "" {
            o.basePath = s
        }
    }

    // 2. 控制器文件目录.
    if g := o.GetOption("controller-path"); g != nil {
        if s := g.ToString(); s != "" {
            o.controllerPath = s
        }
    }

    // 3. 文档(Markdown)存储目录.
    if g := o.GetOption("docs-path"); g != nil {
        if s := g.ToString(); s != "" {
            o.docsPath = s
        }
    }

    // 4. 校验目录.
    stat, err := os.Stat(o.basePath + o.controllerPath)
    if err != nil {
        c.PrintError(err)
        return false
    }

    // 5. 合法目录.
    if stat.IsDir() {
        if o.scanner, err = NewScan(o.basePath, o.controllerPath, o.docsPath); err != nil {
            c.PrintError(err)
            return false
        }

        // 上传位置.
        if g := o.GetOption("upload"); g != nil {
            if s := g.ToString(); s != "" {
                o.scanner.SetUploadUrl(s)
            }
        }

        // 本地存储.
        sl := true
        if g := o.GetOption("save"); g != nil {
            sl = g.ToBool()
        }
        o.scanner.SetSaveLocal(sl)

        return true
    }

    // 6. 无效目录.
    c.PrintError(fmt.Errorf("invalid controller path: %s", o.controllerPath))
    return false
}

// 过程.
func (o *command) run(c i.IConsole) {
    c.Info("[docs] export markdown documents.")
    c.Info("       module: %s", o.scanner.GetModule())
    c.Info("       ---- ---- ---- ---- ---- ---- ---- ----")

    // 导出过程.
    if err := o.scanner.Markdown(); err != nil {
        c.PrintError(err)
        return
    }

    // 结束导出.
    c.Info("       ---- ---- ---- ---- ---- ---- ---- ----")
    c.Info("[docs] end export.")
}
