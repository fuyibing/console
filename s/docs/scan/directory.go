// author: wsfuyibing <websearch@163.com>
// date: 2021-12-26

package scan

import (
    "fmt"
    "os"
    "regexp"
)

type (
    // Directory
    // 目录接口.
    Directory interface {
        Run() error
    }

    // 目标结构体.
    directory struct {
        scanner Scanner
        folder  string
    }
)

// NewDirectory
// 构造目标接口实例.
func NewDirectory(scanner Scanner, folder string) Directory {
    o := &directory{scanner: scanner, folder: folder}
    return o
}

func (o *directory) Run() (err error) {
    var (
        ds   []os.DirEntry
        path = fmt.Sprintf(
            "%s%s%s",
            o.scanner.GetBasePath(),
            o.scanner.GetControllerPath(),
            o.folder,
        )
    )

    // 1. 读取目录.
    if ds, err = os.ReadDir(path); err != nil {
        return err
    }

    // 2. 遍历列表.
    r1 := regexp.MustCompile(`^[a-zA-Z]`)
    r2 := regexp.MustCompile(`\.go$`)

    for _, d := range ds {
        // 目录.
        if d.IsDir() {
            if o.scanner.IsRecursion() && r1.MatchString(d.Name()) {
                if err = NewDirectory(o.scanner, fmt.Sprintf("%s/%s", o.folder, d.Name())).Run(); err != nil {
                    return err
                }
            }
            continue
        }

        // 文件.
        if r2.MatchString(d.Name()) {
            if err = NewReader(o.scanner, o, fmt.Sprintf("%s/%s", o.folder, d.Name())).Run(); err != nil {
                return err
            }
        }
    }

    return nil
}
