// author: wsfuyibing <websearch@163.com>
// date: 2021-12-16

package docs

import (
    "bufio"
    "fmt"
    "os"
    "os/exec"
    "regexp"
    "sort"
    "strings"
    "sync"
    "time"
)

type (
    // Scanner
    // 扫描器接口.
    Scanner interface {
        // AddPayload
        // 添加出入参.
        //
        // .AddPayload("docs/app/logics/topic/CreateRequest")
        // .AddPayload("docs/app/logics/topic/CreateResponse")
        AddPayload(pkg string) error

        // GetBasePath
        // 读取基础(项目)目录.
        //
        // return "/data/sketch"
        GetBasePath() string

        // GetControllerPath
        // 读取控制器目录.
        //
        // return "/app/controllers"
        GetControllerPath() string

        // GetDescription
        // 读取文档描述.
        //
        // return "关于示例文档"
        GetDescription() string

        // GetDocsPath
        // 读取文档(Markdown)存储位置.
        //
        // return "/docs/api"
        GetDocsPath() string

        // GetDomain
        // 读取域名.
        //
        // return "example.com", "www"
        GetDomain() (prefix, domain string)

        // GetHost
        // 读取主机名.
        //
        // return "0.0.0.0", 8080
        GetHost() (host, port string)

        // GetModule
        // 读取模块名称.
        //
        // return "gsdoc"
        GetModule() string

        // GetPayload
        // 读取出入参.
        GetPayload(path string) Payload

        // GetTitle
        // 读取文档标题.
        //
        // return "示例文档"
        GetTitle() string

        // GetVersion
        // 读取版本号.
        //
        // return "0.0"
        GetVersion() string

        // IsRecursion
        // 是否递归目录.
        IsRecursion() bool

        // Markdown
        // 导出Markdown文件.
        Markdown() error

        // Save
        // 保存文件.
        //
        // Save("/data/sketch/docs/api/README.md", `# Title`)
        Save(path, text string) error

        // Template
        // 模板参数替换.
        Template(text string, args map[string]interface{}) string
    }

    // 扫描器结构体.
    scanner struct {
        directory Directory
        mu        *sync.RWMutex

        basePath, controllerPath, docsPath       string
        recursion                                bool
        module, host, port, domain, domainPrefix string
        description, title, version              string

        payloads map[string]Payload
    }
)

// 构造扫描实例.
func NewScan(basePath, controllerPath, docsPath string) (Scanner, error) {
    // 1. 准备实例.
    o := &scanner{}
    o.mu = new(sync.RWMutex)
    o.payloads = make(map[string]Payload)

    // 2. 基础数据.
    o.basePath = basePath
    o.controllerPath = controllerPath
    o.docsPath = docsPath
    o.recursion = true
    o.version = "0.0"
    o.host = "0.0.0.0"
    o.port = "8080"
    o.domain = "example.com"

    // 2.1 基础目录.
    //     计算基础目录的绝对路径.
    if o.basePath == "." || o.basePath == "./" {
        s, err := os.Getwd()
        if err != nil {
            return nil, err
        }
        o.basePath = s
    }

    // 3. 扫描过程.
    //    扫描控制器所在目录下的全部控制器文件, 基于控制器文件与
    //    语句、注释、注解计算文档内容.
    if err := o.run(); err != nil {
        return nil, err
    }

    // 4. 入口处理.
    //    a. 基于payload创建docs/api/main/main.go文件.
    //    c. 执行docs/api/main/main.go生成temp文件.
    if err := o.mainBuilder(); err != nil {
        return nil, err
    }
    if err := o.mainExecute(); err != nil {
        return nil, err
    }
    return o, nil
}

// AddPayload
// 添加出入参文件路径.
func (o *scanner) AddPayload(path string) (err error) {
    path = o.payloadKey(path)
    if m := regexp.MustCompile(`^(.+)\.([_a-zA-Z0-9]+)$`).FindStringSubmatch(path); len(m) == 3 {
        o.mu.Lock()
        defer o.mu.Unlock()

        // 1. 存在.
        if _, ok := o.payloads[path]; ok {
            return
        }

        // 2. 创建.
        var p Payload
        if p, err = NewPayload(o, path, m[1], m[2]); err == nil {
            o.payloads[path] = p
        }
    }
    return
}

// GetPayload
// 读取出入参文件操作实例.
func (o *scanner) GetPayload(path string) Payload {
    path = o.payloadKey(path)
    o.mu.RLock()
    defer o.mu.RUnlock()
    if p, ok := o.payloads[path]; ok {
        return p
    }
    return nil
}

func (o *scanner) GetBasePath() string                { return o.basePath }
func (o *scanner) GetControllerPath() string          { return o.controllerPath }
func (o *scanner) GetDescription() string             { return o.description }
func (o *scanner) GetDocsPath() string                { return o.docsPath }
func (o *scanner) GetDomain() (prefix, domain string) { return o.domainPrefix, o.domain }
func (o *scanner) GetHost() (host, port string)       { return o.host, o.port }
func (o *scanner) GetModule() string                  { return o.module }
func (o *scanner) GetTitle() string                   { return o.title }
func (o *scanner) GetVersion() string                 { return o.version }

func (o *scanner) IsRecursion() bool { return o.recursion }

// Markdown
// 导出Markdown文件.
func (o *scanner) Markdown() error {
    var err error

    // 1. 准备参数.
    //    出参/入参.
    if err = o.prepare(); err != nil {
        return err
    }

    // 2. 导出README.md文件.
    if err = o.render(); err != nil {
        return err
    }

    // 3. 导出接口文件.
    o.directory.EachAction(func(a Action) {
        if err != nil {
            return
        }
        err = a.Markdown()
    })
    if err != nil {
        return err
    }

    // 4. 导出完成.
    o.clean()
    return nil
}

// Save
// 保存文件.
func (o *scanner) Save(path, text string) error {
    var (
        err error
        m   []string
    )

    // 1. 解析路径.
    //    路径名与文件名.
    if m = regexp.MustCompile(`^(.+)/([_a-zA-Z0-9.]+)$`).FindStringSubmatch(path); len(m) == 0 {
        return fmt.Errorf("invalid file path")
    }

    // 2. 创建路径.
    if err = os.MkdirAll(m[1], os.ModePerm); err != nil {
        return err
    }

    // 3. 打开文件.
    var f *os.File
    if f, err = os.OpenFile(path, os.O_CREATE|os.O_RDWR|os.O_TRUNC, os.ModePerm); err != nil {
        return err
    }
    defer func() {
        _ = f.Close()
    }()
    if _, err = f.WriteString(text); err != nil {
        return err
    }
    return nil
}

// Template
// 模板参数替换.
func (o *scanner) Template(text string, args map[string]interface{}) string {
    // 1. 追加变量.
    args["UPDATED"] = time.Now().Format("2006-01-02 15:04")

    // 2. 替换过程.
    r := regexp.MustCompile(`\{\{([_a-zA-Z0-9-]*)\}\}`)
    return strings.ReplaceAll(r.ReplaceAllStringFunc(text, func(s string) string {
        m := r.FindStringSubmatch(s)
        if v, ok := args[m[1]]; ok {
            return fmt.Sprintf("%v", v)
        }
        return ""
    }), "·", "`")
}

func (o *scanner) clean() {
    path := fmt.Sprintf("%s%s/main", o.GetBasePath(), o.GetDocsPath())
    for _, p := range o.payloads {
        os.Remove(fmt.Sprintf("%s/%s", path, p.GetPathKey()))
    }
    os.Remove(fmt.Sprintf("%s/main.go", path))
}

// 基于payload生成临时文件.
func (o *scanner) mainBuilder() error {
    // 基础目录.
    offset := 0
    imports, imported, payloads := make(map[string]string), make([]string, 0), make([]string, 0)
    mainPath := fmt.Sprintf("%s%s/main", o.basePath, o.docsPath)

    for _, x := range o.payloads {
        pkg := x.GetPkg()
        alias := ""
        ok := false

        if alias, ok = imports[pkg]; !ok {
            alias = fmt.Sprintf("t%d", offset)
            imports[pkg] = alias
            imported = append(imported, fmt.Sprintf(`    %s "%s"`, alias, pkg))
            offset++
        }

        payloads = append(
            payloads,
            fmt.Sprintf(`        "%s": %s.%s{},`, x.GetPathKey(), alias, x.GetName()),
        )
    }

    text := o.Template(templateMain, map[string]interface{}{
        "IMPORTS":   strings.Join(imported, "\n"),
        "PAYLOADS":  strings.Join(payloads, "\n"),
        "MAIN_PATH": mainPath,
    })

    path := fmt.Sprintf("%s/main.go", mainPath)
    return o.Save(path, text)
}

// fork子进程, 执行临时文件/docs/api/main/main.go, 生成
// 出入参.
func (o *scanner) mainExecute() error {
    // 1. 临时文件.
    //    docs/api/main/main.go.
    src := fmt.Sprintf("%s%s/main/main.go", o.basePath, o.docsPath)

    // 2. Shell脚本.
    cmd := exec.Command("go", "run", src)

    // 3. 执行命令
    if err := cmd.Run(); err != nil {
        return err
    }
    return nil
}

func (o *scanner) openDocFile() error {
    var (
        err  error
        file *os.File
        m    []string
        p    = fmt.Sprintf("%s%s", o.basePath, o.controllerPath)
        pkg  string
    )

    // 1. 解析包名。
    if m = regexp.MustCompile(`/([_a-zA-Z0-9-]+)$`).FindStringSubmatch(p); len(m) != 2 {
        return fmt.Errorf("can not parse controller package")
    }
    pkg = m[1]

    // 2. 读取文档.
    if file, err = os.OpenFile(fmt.Sprintf("%s/doc.go", p), os.O_RDONLY, os.ModePerm); err != nil {
        return err
    }
    defer func() { _ = file.Close() }()

    // 3. 按行遍历.
    buf := bufio.NewScanner(file)
    comments := make([]string, 0)
    num := 0
    r1 := regexp.MustCompile(`^package\s+`)
    r2 := regexp.MustCompile(`^[/]+\s*`)
    r3 := regexp.MustCompile(`^@([_a-zA-Z0-9]+)\s*\(([^)]*)\)`)
    r4 := regexp.MustCompile(`\s*[.]+$`)
    r5 := regexp.MustCompile(fmt.Sprintf(`package\s+%s`, pkg))
    for buf.Scan() {
        s := strings.TrimSpace(buf.Text())

        // 空行.
        if s == "" {
            comments = make([]string, 0)
            num = 0
            continue
        }

        // 包名.
        if r1.MatchString(s) {
            break
        }

        // 注释.
        if s = strings.TrimSpace(r2.ReplaceAllString(s, "")); s == "" {
            continue
        }

        // 注解.
        if m = r3.FindStringSubmatch(s); len(m) == 3 {
            mk := strings.ToLower(m[1])
            mv := strings.TrimSpace(strings.ReplaceAll(strings.ReplaceAll(m[2], `"`, ""), `'`, ""))

            switch mk {
            case "version":
                if mv != "" {
                    o.version = mv
                }

            case "host":
                if mv != "" {
                    o.host = mv
                }

            case "port":
                if mv != "" {
                    o.port = mv
                }

            case "domain":
                if mv != "" {
                    o.domain = mv
                }

            case "domainprefix":
                if mv != "" {
                    o.domainPrefix = mv
                }
            }

            continue
        }

        // 首行.
        s = r4.ReplaceAllString(s, "")
        if num == 0 {
            if s = strings.TrimSpace(r5.ReplaceAllString(s, "")); s == "" {
                continue
            }
        }
        comments = append(comments, s)
        num++
    }

    // 4. 处理注释.
    for i, s := range comments {
        if i == 0 {
            o.title = s
        } else {
            o.description = fmt.Sprintf("%s%s", o.description, s)
        }
    }
    return nil
}

func (o *scanner) openModFile() error {
    body, err := os.ReadFile(o.basePath + "/go.mod")
    if err != nil {
        return err
    }

    if m := regexp.MustCompile(`module\s+([_a-zA-Z0-9-]+)[^\n]*\n`).FindStringSubmatch(string(body)); len(m) == 2 {
        o.module = m[1]
        o.domainPrefix = o.module
        return nil
    }

    return fmt.Errorf("invalid go mod file")
}

func (o *scanner) payloadKey(s string) string {
    if !strings.HasPrefix(s, o.module) {
        s = fmt.Sprintf("%s/%s", o.module, s)
    }
    return s
}

func (o *scanner) prepare() error {
    return nil
}

func (o *scanner) render() error {
    // 1. 基础变量.
    args := map[string]interface{}{
        "TITLE":         o.title,
        "DESCRIPTION":   o.description,
        "MODULE":        o.module,
        "HOST":          o.host,
        "PORT":          o.port,
        "DOMAIN":        o.domain,
        "DOMAIN_PREFIX": o.domainPrefix,
        "VERSION":       o.version,
    }

    // 2. 菜单列表.
    // cn := make(map[string]string)
    cks := make([]string, 0)
    cms := make(map[string]Controller)

    o.directory.EachController(func(controller Controller) {
        f, n := controller.GetFile()
        k := fmt.Sprintf("%s:%d", f, n)
        cks = append(cks, k)
        cms[k] = controller
    })

    // 3. 菜单处理.
    sort.Strings(cks)
    menu, comma := "", ""

    for _, ck := range cks {
        if c, ok := cms[ck]; ok {
            aks := make([]string, 0)
            ams := make(map[string]Action)

            c.EachAction(func(a Action) {
                if !a.Ignored() {
                    k := a.GetRouteLink()
                    aks = append(aks, k)
                    ams[k] = a
                }
            })

            if len(ams) == 0 {
                continue
            }

            // 连接符.
            menu += comma
            comma = "\n"

            // 控制器.
            menu += fmt.Sprintf("1. %s", c.GetTitle())

            if desc := c.GetDescription(); desc != "" {
                menu += " - " + desc
            }

            // 方法名.
            sort.Strings(aks)
            for _, ak := range aks {
                if a, got := ams[ak]; got {
                    str := fmt.Sprintf(
                        "    1. [%s](.%s) - ·%s·",
                        a.GetTitle(),
                        a.GetRouteLink(),
                        a.GetMethod(),
                    )
                    if sdk := a.GetSdk(); sdk != "" {
                        str = fmt.Sprintf("%s ·SDK·", str)
                    }
                    menu += comma + str
                }
            }
        }
    }

    args["MENU"] = menu

    // 4. 更新模板.
    path := fmt.Sprintf("%s%s/README.md", o.basePath, o.docsPath)
    text := o.Template(templateReadme, args)
    return o.Save(path, text)
}

func (o *scanner) run() error {
    var err error

    // 1. 读取go.mod文件.
    if err = o.openModFile(); err != nil {
        return err
    }

    // 2. 读取app/controllers/doc.go文件.
    if err = o.openDocFile(); err != nil {
        return err
    }

    // 3. 递归控制器文件.
    if o.directory, err = NewDirectory(o, ""); err != nil {
        return err
    }

    return nil
}
