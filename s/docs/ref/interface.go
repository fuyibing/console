// author: wsfuyibing <websearch@163.com>
// date: 2021-12-27

package ref

import "reflect"

type (
    // Element
    // 用于Markdown的元素接口.
    Element interface {
        // Code
        // 加入代码片段.
        Code(cs map[string]interface{}) error

        // Request
        // 处理入参数据.
        Request(s *[]string) error

        // Response
        // 处理出参数据.
        Response(s *[]string) error
    }

    // Block
    // 块实例接口.
    Block interface {
        Element

        // Markdown
        // 导出Markdown文档.
        Markdown(basePath, sourcePath string) error

        // Run
        // 执行块结构解析.
        Run(x reflect.Value) error
    }

    // Field
    // 字段实例接口.
    Field interface {
        Element

        // Run
        // 执行字段解析.
        Run(v reflect.Value, sf reflect.StructField) error

        // SortKey
        // 返回唯一键值.
        SortKey() string
    }
)
