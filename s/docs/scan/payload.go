// author: wsfuyibing <websearch@163.com>
// date: 2021-12-26

package scan

import (
    "fmt"
    "os"
    "regexp"
    "strings"
)

type (
    Payload interface {
        GetName() string
        GetPath() string
        GetPkg() string

        Markdown(resp bool) string
        Postman() string
    }

    payload struct {
        scanner         Scanner
        path, pkg, name string
        message         string
    }
)

func NewPayload(scanner Scanner, path string) Payload {
    o := &payload{scanner: scanner, path: path}
    if m := regexp.MustCompile(`^(\S+)\.([A-Z][_a-zA-Z0-9]*)$`).FindStringSubmatch(path); len(m) == 3 {
        o.pkg = m[1]
        o.name = m[2]
    }
    return o
}

func (o *payload) GetName() string { return o.name }
func (o *payload) GetPath() string { return o.path }
func (o *payload) GetPkg() string  { return o.pkg }

func (o *payload) Markdown(resp bool) string {
    var suffix = "1"
    if resp {
        suffix = "2"
    }

    var (
        name, path = strings.ToLower(
            strings.ReplaceAll(
                strings.TrimPrefix(o.path, "/"),
                "/",
                "-",
            ) + "." + suffix,
        ), fmt.Sprintf("%s%s/main/.md",
            o.scanner.GetBasePath(),
            o.scanner.GetDocsPath(),
        )
    )

    src := fmt.Sprintf("%s/%s", path, name)
    if buf, err := os.ReadFile(src); err == nil {
        return strings.TrimSpace(string(buf))
    }
    return ""
}

func (o *payload) Postman() string {
    var (
        name, path = strings.ToLower(
            strings.ReplaceAll(
                strings.TrimPrefix(o.path, "/"),
                "/",
                "-",
            ) + ".0",
        ), fmt.Sprintf("%s%s/main/.md",
            o.scanner.GetBasePath(),
            o.scanner.GetDocsPath(),
        )
    )

    src := fmt.Sprintf("%s/%s", path, name)
    if buf, err := os.ReadFile(src); err == nil {
        return strings.TrimSpace(string(buf))
    }
    return ""
}
