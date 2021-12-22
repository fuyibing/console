// author: wsfuyibing <websearch@163.com>
// date: 2021-12-16

package docs

var templateMain = `package main

import (
    "fmt"

    "github.com/fuyibing/console/v2/s/docs"
{{IMPORTS}}
)

func main() {
    defer func() {
        if r := recover(); r != nil {
            println(fmt.Errorf("%v", r))
        }
    }()

    b := "{{MAIN_PATH}}"

    m := map[string]interface{}{
{{PAYLOADS}}
    }

    for k, v := range m {
        if err := docs.NewX(b, k, v).Markdown(); err != nil {
            println(err.Error())
            break
        }
    }
}
`
