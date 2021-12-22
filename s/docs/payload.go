// author: wsfuyibing <websearch@163.com>
// date: 2021-12-16

package docs

import (
    "fmt"
    "os"
    "regexp"
)

type (
    // Payload
    // 入出参接口.
    Payload interface {
        GetName() string
        GetPath() string
        GetPathKey() string
        GetPkg() string
        Markdown(response bool) string
    }

    // 出入参结构体.
    payload struct {
        scanner         Scanner
        pkg, name, path string
    }
)

// NewPayload
// 构造出入参实例.
func NewPayload(scanner Scanner, path, pkg, name string) (Payload, error) {
    o := &payload{scanner: scanner, path: path, pkg: pkg, name: name}
    if err := o.run(); err != nil {
        return nil, err
    }
    return o, nil
}

func (o *payload) GetName() string { return o.name }
func (o *payload) GetPath() string { return o.path }
func (o *payload) GetPathKey() string {
    return regexp.MustCompile(`[_/.]`).ReplaceAllString(o.path, "_")
}
func (o *payload) GetPkg() string { return o.pkg }

func (o *payload) Markdown(response bool) string {
    path := fmt.Sprintf("%s%s/main/%s", o.scanner.GetBasePath(), o.scanner.GetDocsPath(), o.GetPathKey())
    body, err := os.ReadFile(path)
    if err == nil {
        return string(body)
    }
    return "payload markdown"
}

func (o *payload) run() error { return nil }
