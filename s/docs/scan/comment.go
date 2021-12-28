// author: wsfuyibing <websearch@163.com>
// date: 2021-12-26

package scan

import (
    "regexp"
    "strings"
)

type (
    // Comment
    // 注释接口.
    Comment interface {
        // Add
        // 添加注释行.
        //
        // .Add("Message")
        // .Add("@Version(1.0)")
        Add(line string)

        // Is
        // 校验注释类型.
        Is(ct CommentType) bool

        // Markdown
        // 导出Markdown文档.
        Markdown() string
    }

    // CommentBlock
    // 注释块接口.
    CommentBlock interface {
        // Comment
        // 注释接口.
        Comment

        // GetAnnotation
        // 读取注解项.
        //
        // .GetAnnotation("Version")
        GetAnnotation(key string) (count int, values []string, has bool)

        // GetAnnotations
        // 读取注解列表.
        GetAnnotations() map[string][]string

        // GetPrefix
        // 读取前缀.
        GetPrefix() string

        // GetTitle
        // 读取标题.
        GetTitle() string

        // SetPrefix
        // 设置前缀.
        //
        // .SetPrefix("Package controllers")
        // .SetPrefix("ExampleController")
        // .SetPrefix("PostList")
        SetPrefix(s string) CommentBlock
    }

    // CommentType
    // 注释类型.
    CommentType int
)

const (
    _ CommentType = iota
    CommentTypeBlock
    CommentTypeCode
    CommentTypeQuote
)

// 注释块结构体.
type commentBlock struct {
    pre         string
    ct          CommentType
    annotations map[string][]string
    comments    []Comment
    first, last Comment
}

// NewComment
// 注释块实例.
func NewComment() CommentBlock {
    o := &commentBlock{ct: CommentTypeBlock}
    o.annotations = make(map[string][]string)
    o.comments = make([]Comment, 0)
    return o
}

func (o *commentBlock) Add(line string)                                { o.add(line) }
func (o *commentBlock) GetAnnotation(key string) (int, []string, bool) { return o.getAnnotation(key) }
func (o *commentBlock) GetAnnotations() map[string][]string            { return o.annotations }
func (o *commentBlock) GetPrefix() string                              { return o.pre }
func (o *commentBlock) GetTitle() string                               { return o.getTitle() }
func (o *commentBlock) Is(ct CommentType) bool                         { return o.is(ct) }
func (o *commentBlock) Markdown() string                               { return o.markdown() }
func (o *commentBlock) SetPrefix(s string) CommentBlock                { o.pre = s; return o }

// 添加注释.
func (o *commentBlock) add(line string) {
    // 1. 注解.
    if o.isAnnotation(line) {
        if o.last != nil {
            o.last = nil
        }
        return
    }

    // 2. 空行.
    if s := strings.TrimSpace(line); s == "" {
        if o.last != nil {
            if o.last.Is(CommentTypeQuote) {
                o.last = nil
            } else if o.last.Is(CommentTypeCode) {
                o.last.Add("")
            }
        }
        return
    }

    // 3. 代码.
    if o.isCode(line) {
        return
    }

    // 4. 注释.
    o.isQuote(line)
}

// 读取注解.
func (o *commentBlock) getAnnotation(key string) (n int, vs []string, has bool) {
    k := strings.ToLower(key)
    if vs, has = o.annotations[k]; has {
        n = len(vs)
    }
    return
}

// 读取标题.
func (o *commentBlock) getTitle() string {
    if o.first != nil {
        return o.first.(*commentQuote).title()
    }
    return ""
}

// 类型比较.
func (o *commentBlock) is(ct CommentType) bool {
    return ct == o.ct
}

// 是否注解.
func (o *commentBlock) isAnnotation(line string) bool {
    b := false
    r := regexp.MustCompile(`^\s*@\s*([_a-zA-Z0-9]+)\s*\(([^)]*)\)`)
    if m := r.FindStringSubmatch(line); len(m) > 0 {
        b = true
        k, v := strings.ToLower(m[1]), strings.TrimSpace(m[2])

        if _, ok := o.annotations[k]; !ok {
            o.annotations[k] = make([]string, 0)
        }

        if v != "" {
            v = strings.TrimSpace(v)
            v = strings.TrimPrefix(v, "'")
            v = strings.TrimPrefix(v, `"`)
            v = strings.TrimSuffix(v, "'")
            v = strings.TrimSuffix(v, `"`)
            v = strings.TrimSpace(v)
        }

        o.annotations[k] = append(
            o.annotations[k],
            v,
        )
    }

    return b
}

// 是否代码.
func (o *commentBlock) isCode(line string) bool {
    b := false
    r := regexp.MustCompile(`^([\s]{2,})(.+)$`)
    if m := r.FindStringSubmatch(line); len(m) == 3 {
        b = true
        if o.last != nil && !o.last.Is(CommentTypeCode) {
            o.last = nil
        }
        if o.last == nil {
            o.last = NewCommentCode(o)
            o.comments = append(o.comments, o.last)
        }
        o.last.Add(line)
    }
    return b
}

// 是否注释.
func (o *commentBlock) isQuote(line string) {
    if o.last != nil && !o.last.Is(CommentTypeQuote) {
        o.last = nil
    }

    if o.last == nil {
        o.last = NewCommentQuote(o)
        if o.first == nil {
            o.first = o.last
        } else {
            o.comments = append(o.comments, o.last)
        }
    }

    line = strings.TrimSpace(line)
    o.last.Add(line)
    if regexp.MustCompile(`[.]+$`).MatchString(line) {
        o.last = nil
    }
}

// 导出文档.
func (o *commentBlock) markdown() string {
    cs := make([]string, 0)
    for _, c := range o.comments {
        cs = append(cs, c.Markdown())
    }

    return strings.Join(cs, "\n\n")
}
