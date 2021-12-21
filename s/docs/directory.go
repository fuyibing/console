// author: wsfuyibing <websearch@163.com>
// date: 2021-12-16

package docs

import (
    "fmt"
    "os"
    "regexp"
    "strings"
    "sync"
)

type (
    Directory interface {
        // EachAction
        // 遍历Action.
        EachAction(callback func(action Action))
        EachController(callback func(controller Controller))

        // GetController
        // 读取控制器实例.
        //
        // GetController("ExampleController")
        GetController(name string) Controller

        // GetControllerWithInitialize
        // 读取控制器实例, 若不存在则创建.
        //
        // GetControllerWithInitialize("ExampleController")
        GetControllerWithInitialize(name string) Controller

        // GetFolder
        // 读取目录文件夹名称.
        //
        // return "/topic/manage"
        GetFolder() string

        // Scanner
        // 返回Scanner实例.
        Scanner() Scanner
    }

    directory struct {
        folder  string
        mu      *sync.RWMutex
        scanner Scanner

        controllers map[string]Controller
        directories []Directory
    }
)

// 构造目录实例.
func NewDirectory(scanner Scanner, folder string) (Directory, error) {
    o := &directory{scanner: scanner, folder: folder}
    o.mu = new(sync.RWMutex)
    o.controllers = make(map[string]Controller, 0)
    o.directories = make([]Directory, 0)
    if err := o.run(); err != nil {
        return nil, err
    }
    return o, nil
}

// 遍历.
func (o *directory) EachAction(callback func(action Action)) {
    if o.scanner.IsRecursion() {
        for _, d := range o.directories {
            d.EachAction(callback)
        }
    }
    for _, c := range o.controllers {
        c.EachAction(callback)
    }
}

func (o *directory) EachController(callback func(controller Controller)) {
    if o.scanner.IsRecursion() {
        for _, d := range o.directories {
            d.EachController(callback)
        }
    }
    for _, c := range o.controllers {
        callback(c)
    }
}

// GetController
// 读取控制器实例.
func (o *directory) GetController(name string) Controller {
    o.mu.RLock()
    defer o.mu.RUnlock()
    if c, ok := o.controllers[name]; ok {
        return c
    }
    return nil
}

// GetControllerWithCreate
// 读取控制器实例, 若不存在则初始化.
func (o *directory) GetControllerWithInitialize(name string) Controller {
    if c := o.GetController(name); c != nil {
        return c
    }

    o.mu.Lock()
    defer o.mu.Unlock()

    c := NewController(o, name)
    o.controllers[name] = c
    return c
}

// GetFolder
// 读取目录文件夹名称.
func (o *directory) GetFolder() string { return o.folder }

func (o *directory) Scanner() Scanner { return o.scanner }

// 执行扫描.
func (o *directory) run() error {
    // 1. 当前目录.
    p := fmt.Sprintf("%s%s%s", o.scanner.GetBasePath(), o.scanner.GetControllerPath(), o.folder)

    // 2. 打开目录.
    ds, err := os.ReadDir(p)
    if err != nil {
        return err
    }

    // 3. 遍历文件.
    r1 := regexp.MustCompile(`^[a-zA-Z][_a-zA-Z0-9-]*$`)
    r2 := regexp.MustCompile(`^([_a-z0-9-]*controller\.go)$`)

    for _, d := range ds {
        // 子目录.
        if d.IsDir() {
            if r1.MatchString(d.Name()) && o.scanner.IsRecursion() {
                child, err2 := NewDirectory(o.scanner, o.folder+"/"+d.Name())
                if err2 != nil {
                    return err2
                }
                o.directories = append(o.directories, child)
            }
            continue
        }
        // 控制器文件.
        if r2.MatchString(strings.ToLower(d.Name())) {
            _, err3 := NewControllerReader(o, o.folder, d.Name())
            if err3 != nil {
                return err3
            }
        }
    }

    return nil
}
