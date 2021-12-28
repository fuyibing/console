// author: wsfuyibing <websearch@163.com>
// date: 2021-12-16

package scan

var templateMain = `package main

import (
    "fmt"
    "reflect"

    "gsjob/tests/ref"
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
        p := ref.NewBlock(0)
        if err := p.Run(reflect.ValueOf(v)); err == nil {
            if err = p.Markdown(b, k); err != nil {
                break
            }
        }
    }
}
`
