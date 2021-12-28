// author: wsfuyibing <websearch@163.com>
// date: 2021-12-26

package scan

import (
    "fmt"
    "regexp"
    "strings"
)

// 文本注释片段.
type commentQuote struct {
    cb    CommentBlock
    ct    CommentType
    lines []string
}

func NewCommentQuote(cb CommentBlock) Comment {
    o := &commentQuote{cb: cb, ct: CommentTypeQuote}
    o.lines = make([]string, 0)
    return o
}

func (o *commentQuote) Add(line string)        { o.add(line) }
func (o *commentQuote) Is(ct CommentType) bool { return o.is(ct) }
func (o *commentQuote) Markdown() string       { return o.markdown() }

// 添加注释.
//
//   .add("Message")
func (o *commentQuote) add(line string) (nl bool) {
    if line = strings.TrimSpace(line); line == "" {
        return true
    }

    re := regexp.MustCompile(`[.]+$`)
    if re.MatchString(line) {
        nl = true
        line = re.ReplaceAllString(line, "")
    }

    o.lines = append(o.lines, line)
    return nl
}

// 类型比例.
//
//   return true
//   return false
func (o *commentQuote) is(ct CommentType) bool {
    return ct == o.ct
}

// 注释内容.
//
//   return "> Message"
func (o *commentQuote) markdown() string {
    if len(o.lines) > 0 {
        return fmt.Sprintf("> %s.", strings.Join(o.lines, ""))
    }
    return ""
}

func (o *commentQuote) title() string {
    if len(o.lines) > 0 {
        s := strings.Join(o.lines, "")
        if p := o.cb.GetPrefix(); p != "" {
            s = strings.TrimPrefix(s, p)
        }
        return strings.TrimSpace(s)
    }
    return ""
}
