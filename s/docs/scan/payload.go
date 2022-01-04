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
    if resp {
        return o.fileRead("2")
    }
    return o.fileRead("1")
}

func (o *payload) Postman() string {
    return o.fileRead("0")
}

func (o *payload) fileName(suffix string) string {
    name := strings.ToLower(strings.ReplaceAll(strings.TrimPrefix(o.path, "/"), "/", "-") + "." + suffix)
    path := fmt.Sprintf("%s%s/main/.md", o.scanner.GetBasePath(), o.scanner.GetDocsPath())
    return fmt.Sprintf("%s/%s", path, name)
}

func (o *payload) fileRead(suffix string) string {
    if buf, err := os.ReadFile(o.fileName(suffix)); err == nil {
        return strings.TrimSpace(string(buf))
    }
    return ""
}
