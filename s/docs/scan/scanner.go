// author: wsfuyibing <websearch@163.com>
// date: 2021-12-26

package scan

import (
    "encoding/json"
    "fmt"
    "io/ioutil"
    "net/http"
    "os"
    "os/exec"
    "regexp"
    "sort"
    "strings"
    "time"
)

type (
    // Scanner
    // 扫描接口.
    Scanner interface {
        AddPayload(path string)
        Clean() error

        GetBasePath() string
        GetController(key string) Controller
        GetControllerPath() string
        GetDocsPath() string
        GetModule() string
        GetPayload(path string) Payload

        IsRecursion() bool
        IsSaveEnable() bool
        IsUploadEnable() bool

        Markdown() error
        Run() error

        // Save
        // 保存文件.
        Save(path, text string) error

        SetBasePath(basePath string) Scanner
        SetCommentBlock(cb CommentBlock) Scanner
        SetControllerPath(controllerPath string) Scanner
        SetDocsPath(docsPath string) Scanner
        SetRecursion(b bool) Scanner
        SetSaveEnable(b bool) Scanner
        SetUploadEnable(b bool) Scanner
        SetUploadUrl(s string) Scanner
        Upload(name, text string) error
    }

    // 扫描结构体.
    scanner struct {
        basePath, controllerPath, docsPath string
        module                             string
        controllers                        map[string]Controller

        payloads map[string]Payload

        recursion    bool
        saveEnable   bool
        uploadEnable bool
        uploadUrl    string

        docVersion                                   string
        docTitle, docDescription                     string
        docHost, docPort, docDomain, docDomainPrefix string
    }
)

// NewScanner
// 构造Scanner接口.
func NewScanner() Scanner {
    o := &scanner{}
    o.saveEnable = true
    o.uploadEnable = false
    o.controllers = make(map[string]Controller)
    o.payloads = make(map[string]Payload)
    o.docTitle = o.module
    o.docDescription = ""
    o.docHost, o.docPort, o.docDomain, o.docDomainPrefix = "0.0.0.0", "8080", o.module, "example.com"
    o.docVersion = "0.0"
    return o
}

func (o *scanner) AddPayload(path string) {
    path = strings.TrimPrefix(path, "/")
    if !strings.HasSuffix(path, fmt.Sprintf("%s/", o.module)) {
        path = fmt.Sprintf("%s/%s", o.module, path)
    }

    if _, ok := o.payloads[path]; ok {
        return
    }

    o.payloads[path] = NewPayload(o, path)
}
func (o *scanner) Clean() error { return o.clean() }
func (o *scanner) GetPayload(path string) Payload {
    path = strings.TrimPrefix(path, "/")
    if !strings.HasSuffix(path, fmt.Sprintf("%s/", o.module)) {
        path = fmt.Sprintf("%s/%s", o.module, path)
    }
    if v, ok := o.payloads[path]; ok {
        return v
    }
    return nil
}
func (o *scanner) GetBasePath() string { return o.basePath }
func (o *scanner) GetController(key string) Controller {
    if c, ok := o.controllers[key]; ok {
        return c
    }
    c := NewController()
    o.controllers[key] = c
    return c
}
func (o *scanner) GetControllerPath() string { return o.controllerPath }
func (o *scanner) GetDocsPath() string       { return o.docsPath }
func (o *scanner) GetModule() string         { return o.module }
func (o *scanner) IsRecursion() bool         { return o.recursion }
func (o *scanner) IsSaveEnable() bool        { return o.saveEnable }
func (o *scanner) IsUploadEnable() bool      { return o.uploadEnable }
func (o *scanner) Markdown() error           { return o.markdown() }
func (o *scanner) Run() error                { return o.run() }

// Save
// 写入文件.
func (o *scanner) Save(path, text string) error {
    var (
        err error
        src = path
    )

    defer func() {
        if err != nil {
            println("       saveto error:", err.Error())
        } else {
            println("       saveto [", src, "]")
        }
    }()

    path = fmt.Sprintf("%s%s", o.basePath, path)

    // 1. 检查.
    //    解析路径出错.
    m := regexp.MustCompile(`^(.+)/([_a-zA-Z0-9-.]+)$`).FindStringSubmatch(path)
    if len(m) == 0 {
        return fmt.Errorf("parse saved fail")
    }

    // 2. 目录.
    //    创建基础目录.
    if err = os.MkdirAll(m[1], os.ModePerm); err != nil {
        return err
    }

    // 3. 打开.
    f, err := os.OpenFile(path, os.O_CREATE|os.O_RDWR|os.O_TRUNC, os.ModePerm)
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

func (o *scanner) SetBasePath(basePath string) Scanner {
    o.basePath = basePath
    return o
}

func (o *scanner) SetCommentBlock(cb CommentBlock) Scanner {
    o.setCommentBlock(cb)
    return o
}

func (o *scanner) SetControllerPath(controllerPath string) Scanner {
    o.controllerPath = controllerPath
    return o
}

func (o *scanner) SetDocsPath(docsPath string) Scanner {
    o.docsPath = docsPath
    return o
}

func (o *scanner) SetRecursion(b bool) Scanner {
    o.recursion = b
    return o
}

func (o *scanner) SetSaveEnable(b bool) Scanner {
    o.saveEnable = b
    return o
}

func (o *scanner) SetUploadEnable(b bool) Scanner {
    o.uploadEnable = b
    return o
}

func (o *scanner) SetUploadUrl(s string) Scanner {
    if !regexp.MustCompile(`^https*://`).MatchString(s) {
        s = fmt.Sprintf("http://%s", s)
    }
    o.uploadUrl = s
    return o
}

// Upload
// 上传文件.
func (o *scanner) Upload(name, content string) error {
    var (
        buf       []byte
        data      = map[string]string{"key": o.module, "name": strings.TrimPrefix(name, "/"), "content": content}
        body, err = json.Marshal(data)
        req       *http.Request
        res       *http.Response
    )

    // 1. 记录结果.
    defer func() {
        if err != nil {
            println("       upload error:", err.Error())
        } else {
            route := regexp.MustCompile(`/([_a-zA-Z0-9]+).md$`).ReplaceAllString(data["name"], "")
            target := strings.TrimSuffix(o.uploadUrl, "/") + "/" + data["key"] + "/" + data["name"]
            println("       upload [", route, "] to [", target, "]")
        }
    }()

    // 2. 入参校验.
    if err != nil {
        return err
    }

    // 3. 上传请求.
    if req, err = http.NewRequest(http.MethodPut, o.uploadUrl, strings.NewReader(string(body))); err != nil {
        return err
    }

    // 4. 发送内容.
    req.Header.Set("Content-Type", "application/json")
    if res, err = (&http.Client{}).Do(req); err != nil {
        return err
    }
    defer func() { _ = res.Body.Close() }()
    if buf, err = ioutil.ReadAll(res.Body); err != nil {
        return err
    }

    // 5. 发送结果.
    v := &struct {
        Errno int    `json:"errno"`
        Error string `json:"error"`
    }{}
    if err = json.Unmarshal(buf, v); err != nil {
        return err
    }

    // 6. 上传失败.
    if v.Errno != 0 {
        err = fmt.Errorf(v.Error)
        return err
    }

    // 7. 完成上传.
    return nil
}

func (o *scanner) clean() error { return nil }

func (o *scanner) eachAction(fn func(a Action)) {
    if fn != nil {
        for _, c := range o.controllers {
            c.Each(fn)
        }
    }
}

func (o *scanner) eachController(fn func(c Controller)) {
    if fn != nil {
        for _, c := range o.controllers {
            fn(c)
        }
    }
}

// 执行脚本.
//
// 执行由runScript()方法生成的main.go文件.
func (o *scanner) execution() error {
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

// 导出文档.
func (o *scanner) markdown() error {
    // 1. 临时文件.
    if err := o.execution(); err != nil {
        return err
    }

    // 2. 导航菜单.
    if err := o.renderMenu(); err != nil {
        return err
    }

    // 3. 文档主页.
    if err := o.renderReadme(); err != nil {
        return err
    }

    // 4. 文档明细.
    return o.renderActions()
    // return nil
}

func (o *scanner) renderActions() (err error) {
    o.eachAction(func(a Action) {
        if err == nil {
            err = a.Markdown()
        }
    })
    return
}

// 渲染菜单.
func (o *scanner) renderMenu() error {
    list := make([]string, 0)

    // 1. 控制器列表.
    cs := make([]string, 0)
    cm := make(map[string]Controller)
    o.eachController(func(c Controller) {
        ck := c.GetName()
        cs = append(cs, ck)
        cm[ck] = c
    })

    // 2. 遍历控制器.
    //    按名称顺序排列.
    sort.Strings(cs)
    for _, ck := range cs {
        if c, ok := cm[ck]; ok {
            // 2.1 API列表
            as := make([]string, 0)
            am := make(map[string]Action)
            c.Each(func(a Action) {
                if !a.Ignore() {
                    ak := a.GetRouteLink()
                    as = append(as, ak)
                    am[ak] = a
                }
            })
            if len(as) == 0 {
                continue
            }

            // 2.2 遍历API.
            //     按URI顺序排列.
            list = append(list, fmt.Sprintf(`1. %s`, c.GetTitle()))
            for _, ak := range as {
                if a, ao := am[ak]; ao {
                    list = append(list, fmt.Sprintf(
                        `    1. [%s](.%s)`,
                        a.GetTitle(),
                        a.GetRouteLink(),
                    ))
                }
            }
        }
    }

    text := strings.ReplaceAll(strings.Join(list, "\n"), "·", "`")

    if o.saveEnable {
        if err := o.Save(
            fmt.Sprintf(
                "%s/menu.md",
                o.GetDocsPath(),
            ), text,
        ); err != nil {
            return err
        }
    }

    if o.uploadEnable {
        if err := o.Upload("menu.md", text); err != nil {
            return err
        }
    }

    return nil
}

// 渲染入口.
func (o *scanner) renderReadme() error {
    text := ""

    // 1. 标题.
    text += fmt.Sprintf(`# %s`, strings.TrimSpace(o.docTitle)) + "\n"
    text += "\n"

    // 2. 参数.
    text += fmt.Sprintf("**域名** : ·%s.%s·", o.docDomainPrefix, o.docDomain) + "<br />\n"
    text += fmt.Sprintf("**部署** : ·%s:%s·", o.docHost, o.docPort) + "<br />\n"
    text += fmt.Sprintf("**版本** : ·%s·", o.docVersion) + "<br />\n"
    text += fmt.Sprintf("**更新** : ·%s·", time.Now().Format("2006-01-02 15:04")) + "<br />\n"
    text += "\n"

    // 3. 描述.
    if o.docDescription != "" {
        text += o.docDescription + "\n"
        text += "\n"
    }

    // n. 标记.
    text = strings.ReplaceAll(text, "·", "`")

    if o.saveEnable {
        if err := o.Save(
            fmt.Sprintf(
                "%s/README.md",
                o.GetDocsPath(),
            ), text,
        ); err != nil {
            return err
        }
    }

    if o.uploadEnable {
        if err := o.Upload("README.md", text); err != nil {
            return err
        }
    }

    return nil
}

// 执行扫描.
func (o *scanner) run() error {
    // 1. 读取模块.
    buf, err := os.ReadFile(fmt.Sprintf("%s/go.mod", o.basePath))
    if err != nil {
        return err
    }

    // 2. 解析模块.
    get := regexp.MustCompile(`module\s+([_a-zA-Z0-9-.]+)`).FindStringSubmatch(string(buf))
    if len(get) == 0 {
        return fmt.Errorf("parse module fail")
    }
    o.module = get[1]

    // 3. 目录扫描.
    if err = NewDirectory(o, "").Run(); err != nil {
        return err
    }

    // 4. 创建脚本.
    return o.runScript()
}

// 生成脚本.
//
// 可成可执行的main.go脚本, 用于通过反射解析结构体的出入参数据.
func (o *scanner) runScript() error {
    var (
        importKey  = 0
        importKeys = make(map[string]string)
        importList = make([]string, 0)
        payloads   = make([]string, 0)
    )

    for k, x := range o.payloads {
        pkg := x.GetPkg()
        ali, ok := importKeys[pkg]

        if !ok {
            ali = fmt.Sprintf("t%d", importKey)
            importKey++
            importKeys[pkg] = ali
            importList = append(importList, fmt.Sprintf(`    %s "%s"`, ali, pkg))
        }

        payloads = append(
            payloads,
            fmt.Sprintf(`        "%s" : %s.%s{}, `, k, ali, x.GetName()),
        )
    }

    // m. 变更.
    res := map[string]string{
        "MAIN_PATH": fmt.Sprintf("%s%s/main", o.basePath, o.docsPath),
        "IMPORTS":   strings.Join(importList, "\n"),
        "PAYLOADS":  strings.Join(payloads, "\n"),
    }

    // n. 模板.
    //    替换模板变量, 写入目标文件.
    reg := regexp.MustCompile(`{{([_a-zA-Z0-9]+)}}`)
    str := reg.ReplaceAllStringFunc(templateMain, func(s string) string {
        if m := reg.FindStringSubmatch(s); len(m) > 0 {
            if v, ok := res[m[1]]; ok {
                return v
            }
        }
        return ""
    })
    return o.Save(o.docsPath+"/main/main.go", str)
}

func (o *scanner) setCommentBlock(cb CommentBlock) {
    if s := cb.GetTitle(); s != "" {
        o.docTitle = s
    }
    if s := cb.Markdown(); s != "" {
        o.docDescription = s
    }

    o.docHost, o.docPort, o.docDomain, o.docDomainPrefix = "0.0.0.0", "8080", "example.com", o.module
    o.docVersion = "0.0"

    if n, vs, has := cb.GetAnnotation("host"); has && n > 0 && vs[0] != "" {
        o.docHost = vs[0]
    }
    if n, vs, has := cb.GetAnnotation("port"); has && n > 0 && vs[0] != "" {
        o.docPort = vs[0]
    }
    if n, vs, has := cb.GetAnnotation("domain"); has && n > 0 && vs[0] != "" {
        o.docDomain = vs[0]
    }
    if n, vs, has := cb.GetAnnotation("domainprefix"); has && n > 0 && vs[0] != "" {
        o.docDomainPrefix = vs[0]
    }
    if n, vs, has := cb.GetAnnotation("version"); has && n > 0 && vs[0] != "" {
        o.docVersion = vs[0]
    }
}
