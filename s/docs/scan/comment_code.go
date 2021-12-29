// author: wsfuyibing <websearch@163.com>
// date: 2021-12-26

package scan

import (
    "fmt"
    "strings"
)

// 代码片段注释.
type commentCode struct {
    cb    CommentBlock
    ct    CommentType
    lines []string
}

func NewCommentCode(cb CommentBlock) Comment {
    o := &commentCode{cb: cb, ct: CommentTypeCode}
    o.lines = make([]string, 0)
    return o
}

func (o *commentCode) Add(line string)        { o.add(line) }
func (o *commentCode) Is(ct CommentType) bool { return o.is(ct) }
func (o *commentCode) Markdown() string       { return o.markdown() }

// 添加代码.
//
//   .add("n := 1")
func (o *commentCode) add(line string) bool {
    if len(line) > 2 {
        o.lines = append(o.lines, line[2:])
    } else {
        o.lines = append(o.lines, "")
    }
    return false
}

// 类型比例.
//
//   return true
//   return false
func (o *commentCode) is(ct CommentType) bool {
    return ct == o.ct
}

// 代码片段.
//
//   return ```
//   n := 1
//   ```
func (o *commentCode) markdown() string {
    if n := len(o.lines); n > 0 {
        x := n

        for i := n - 1; i >= 0; i-- {
            if s := o.lines[i]; s == "" {
                x--
                continue
            }
            break
        }

        if x > 0 {
            return fmt.Sprintf(
                "```\n%s\n```",
                strings.Join(o.lines[0:x], "\n"),
            )
        }
    }
    return ""
}
