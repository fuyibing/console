// author: wsfuyibing <websearch@163.com>
// date: 2021-03-02

package docs

import (
    "path/filepath"

    "github.com/fuyibing/console/v2/s/docs/scan"

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
    clean, saveEnable, uploadEnable    bool
    scanner                            scan.Scanner
    uploadUrl                          string
}

// 构造导出文档实例.
func New() *command {
    // 1. 构建实例.
    o := &command{clean: true}
    o.basePath = "./"
    o.controllerPath = "/app/controllers"
    o.docsPath = "/docs/api"
    o.saveEnable = true
    o.uploadEnable = false
    o.uploadUrl = "http://gs-docs.turboradio.cn"

    // 2. 初始化.
    o.Initialize()
    o.SetName(Name)
    o.SetDescription(Description)

    // 3. 执行过程.
    o.SetHandlerBefore(o.before)
    o.SetHandler(o.run)
    o.SetHandlerAfter(o.after)

    // 4. 路径定义
    o.Add(base.NewOption(i.OptionalMode, i.StrValue).
        SetName("base-path").
        SetDefaultValue("./").
        SetDescription("application base path"),
    )
    o.Add(base.NewOption(i.OptionalMode, i.StrValue).
        SetName("controller-path").
        SetDefaultValue("/app/controllers").
        SetDescription("controller path of application"))
    o.Add(base.NewOption(i.OptionalMode, i.StrValue).
        SetName("docs-path").
        SetDefaultValue("/docs/api").
        SetDescription("documents path of application"))

    // 5. 存储状态.
    o.Add(base.NewOption(i.OptionalMode, i.BoolValue).
        SetName("save").
        SetShortName("s").
        SetDefaultValue("true").
        SetDescription("save to documents path or not, default: true"),
    )

    // 6. 上传状态.
    o.Add(base.NewOption(i.OptionalMode, i.BoolValue).
        SetName("upload").
        SetShortName("u").
        SetDescription("upload to server or not, default: false"),
    )

    // 8. 上传位置.
    o.Add(base.NewOption(i.OptionalMode, i.StrValue).
        SetName("upload-url").
        SetDefaultValue("gs-docs.turboradio.cn").
        SetDescription("where documents storage"),
    )

    // 7. 完成配置.
    return o
}

// 后置.
func (o *command) after(c i.IConsole) {
    if o.clean {
        o.scanner.Clean()
    }
}

// 前置.
func (o *command) before(c i.IConsole) bool {
    // 1. 项目目录.
    if g := o.GetOption("base-path"); g != nil {
        if s := g.ToString(); s != "" {
            o.basePath = s
        }
    }

    // 2. 控制器目录.
    if g := o.GetOption("controller-path"); g != nil {
        if s := g.ToString(); s != "" {
            o.controllerPath = s
        }
    }

    // 3. 文档目录.
    if g := o.GetOption("docs-path"); g != nil {
        if s := g.ToString(); s != "" {
            o.docsPath = s
        }
    }

    // 4. 是否存储.
    if g := o.GetOption("save"); g != nil {
        o.saveEnable = g.ToBool()
    }

    // 5. 是否上传.
    if g := o.GetOption("upload"); g != nil {
        o.uploadEnable = g.ToBool()
    }

    // 6. 上传位置.
    if g := o.GetOption("upload-url"); g != nil {
        if s := g.ToString(); s != "" {
            o.uploadUrl = s
        }
    }

    // 7. 完成.
    if o.basePath == "" || o.basePath == "." || o.basePath == "./" {
        if s, err := filepath.Abs("./"); err == nil {
            o.basePath = s
        }
    }
    return true
}

// 过程.
func (o *command) run(c i.IConsole) {
    // 1. 准备执行.
    o.scanner = scan.NewScanner()
    o.scanner.SetBasePath(o.basePath).SetControllerPath(o.controllerPath).SetDocsPath(o.docsPath)
    o.scanner.SetSaveEnable(o.saveEnable).SetUploadEnable(o.uploadEnable).SetUploadUrl(o.uploadUrl)

    // 2. 扫描目录.
    if err := o.scanner.Run(); err != nil {
        o.clean = false
        c.PrintError(err)
        return
    }

    // 3. 导出过程.
    c.Info("[docs] export markdown documents.")
    c.Info("       module: %s", o.scanner.GetModule())
    c.Info("       ---- ---- ---- ---- ---- ---- ---- ----")
    if err := o.scanner.Markdown(); err != nil {
        o.clean = false
        c.PrintError(err)
        return
    }

    // 4. 结束导出.
    c.Info("       ---- ---- ---- ---- ---- ---- ---- ----")
    c.Info("[docs] end export.")
}
