// author: wsfuyibing <websearch@163.com>
// date: 2021-03-02

package docs

import (
	"github.com/fuyibing/console/v2/base"
	"github.com/fuyibing/console/v2/i"
	"github.com/fuyibing/gdoc/adapters/markdown"
	"github.com/fuyibing/gdoc/adapters/postman"
	gb "github.com/fuyibing/gdoc/base"
	"github.com/fuyibing/gdoc/conf"
	"github.com/fuyibing/gdoc/reflectors"
	"github.com/fuyibing/gdoc/scanners"
)

const (
	Description = "Export application document"
	Name        = "docs"
)

// 导出文档.
type command struct {
	base.Command

	basePath, controllerPath, docsPath string
}

// New
// 构造导出文档实例.
func New() *command {
	// 1. 构建实例.
	o := &command{}
	o.basePath = "./"
	o.controllerPath = "/app/controllers"
	o.docsPath = "/docs"

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

	// 7. 完成配置.
	return o
}

// 后置.
func (o *command) after(c i.IConsole) {
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

	// Init config
	conf.Path.SetBasePath(o.basePath)
	conf.Path.SetControllerPath(o.controllerPath)
	conf.Path.SetDocumentPath(o.docsPath)
	conf.Config.Load()

	return true
}

// 过程.
func (o *command) run(c i.IConsole) {
	scanners.Scanner.Scan()

	ref := reflectors.New(gb.Mapper)
	ref.Configure()
	if err := ref.Make(); err != nil {
		return
	}

	postman.New(gb.Mapper).Run()
	markdown.New(gb.Mapper).Run()
	ref.Clean()
}
