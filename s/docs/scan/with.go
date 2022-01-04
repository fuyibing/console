// author: wsfuyibing <websearch@163.com>
// date: 2022-01-04

package scan

var With WithInterface

type (
    // ErrCode
    // 错误编码.
    ErrCode int

    // WithInterface
    // API模式下标准化的返回数据结构.
    WithInterface interface {
        // Data
        // 返回Object数据结构.
        //
        // return app.With.Data(data)
        //
        // return {
        //     "errno": 0,
        //     "error": "",
        //     "data": {
        //         "key": "value"
        //     },
        //     "dataType": "OBJECT"
        // }
        Data(data interface{}) WithResponse

        // Error
        // 返回Error数据结构.
        //
        // return app.With.Error(err)
        //
        // return {
        //     "errno": 404,
        //     "error": "HTTP 404 NOT FOUND",
        //     "data": {},
        //     "dataType": "ERROR"
        // }
        Error(err error) WithResponse

        // ErrorCode
        // 返回Error数据结构.
        ErrorCode(err error, code ErrCode) WithResponse

        // Html
        // 返回HTML数据.
        Html(text string) string

        // List
        // 返回List数据结构.
        //
        // return app.With.List(list)
        //
        // return {
        //     "errno": 0,
        //     "error": "",
        //     "data": [
        //          {
        //              "id": 1,
        //              "key": "value",
        //          },
        //          {
        //              "id": 2,
        //              "key": "value - 2",
        //          }
        //     ],
        //     "dataType": "LIST"
        // }
        List(data interface{}) WithResponse

        // Paging
        // 返回Paging数据结构.
        //
        // return {
        //     "errno": 0,
        //     "error": "",
        //     "data": {
        //         "body": [
        //              {
        //                  "id": 1,
        //                  "key": "value",
        //              },
        //              {
        //                  "id": 2,
        //                  "key": "value - 2",
        //              }
        //         ],
        //         "paging": {
        //         }
        //     },
        //     "dataType": "PAGING"
        // }
        Paging(data interface{}, total int64, page, size int) WithResponse

        // String
        // 字符串数据类型.
        //
        // return "example"
        String(str string) string

        // Success
        // 返回Object数据类型.
        //
        // return app.With.Success()
        Success() WithResponse
    }

    // WithPaging
    // 分页数据结构.
    WithPaging struct {
        First      int   `json:"first" mock:"1" label:"首页"`
        Before     int   `json:"before" mock:"1" label:"前页"`
        Current    int   `json:"current" mock:"1" label:"本页"`
        Next       int   `json:"next" mock:"1" label:"后页"`
        Last       int   `json:"last" mock:"1" label:"尾页"`
        Limit      int   `json:"limit" mock:"1" label:"每页"`
        TotalPages int   `json:"totalPages" mock:"1" label:"总页数"`
        TotalItems int64 `json:"totalItems" mock:"1" label:"总条数"`
    }

    // WithResponse
    // 数据主体结构.
    WithResponse struct {
        Data     interface{} `json:"data" label:"主数据结构" desc:"匹配Object, List, Paging, Error等数据结构"`
        Errno    int         `json:"errno" label:"错误编码" desc:"0: 成功<br />n: 错误码"`
        Error    string      `json:"error" label:"错误原因" desc:"当errno非0时, 此处错误描述"`
        DataType string      `json:"dataType" label:"数据类型" desc:"指明data字段的数据结构"`
    }
)

// 错误码定义.
const (
    ErrCodePlus           = 1000 // 最小Code码标准值.
    SystemErrCode ErrCode = iota // 系统错误/当未指定Code码时, 使用本编码.
)

// Int
// 错误转成整型.
//
// 当转码时, 200-999保留给http状态码, 不进行转换, 当编码
// 大0-199之间时, 自动加上最小Code码标准值.
//
//   // 用法 1
//   err := fmt.Errorf("my error")
//   app.With.ErrorCode(err, app.SystemErrCode)
//
//   // 用法 2
//   println("code1:", app.ErrCode(404).Int())  // 404
//   println("code2:", app.ErrCode(2).Int())    // 1002
//   println("code2:", app.SystemErrCode.Int()) // 1000
func (c ErrCode) Int() int {
    n := int(c)
    if c < 200 {
        return ErrCodePlus + n
    }
    return n
}

// 输出结构.
type with struct{}

// Data
// 返回Object数据结构.
func (o *with) Data(data interface{}) WithResponse {
    return WithResponse{
        Data: data, DataType: "OBJECT",
        Errno: 0, Error: "",
    }
}

// Error
// 返回Error数据结构.
func (o *with) Error(err error) WithResponse { return o.ErrorCode(err, SystemErrCode) }

// ErrorCode
// 返回Error数据结构.
func (o *with) ErrorCode(err error, code ErrCode) WithResponse {
    return WithResponse{
        Data: map[string]interface{}{}, DataType: "ERROR",
        Errno: code.Int(), Error: err.Error(),
    }
}

// Html
// 返回HTML内容.
func (o *with) Html(text string) string { return text }

// List
// 返回列表数据结构.
func (o *with) List(data interface{}) WithResponse {
    return WithResponse{
        Data: map[string]interface{}{"body": data}, DataType: "LIST",
        Errno: 0, Error: "",
    }
}

// Paging
// 返回分页数据结构.
func (o *with) Paging(data interface{}, total int64, page, size int) WithResponse {
    return WithResponse{
        Data: map[string]interface{}{"body": data, "paging": o.paging(total, page, size)}, DataType: "PAGING",
        Errno: 0, Error: "",
    }
}

// String
// 返回文本.
func (o *with) String(str string) string { return str }

// Success
// 返回成功数据结构.
func (o *with) Success() WithResponse { return o.Data(map[string]interface{}{}) }

// 初始化.
// 在包的init()方法中调用.
func (o *with) init() *with { return o }

// 分页器.
// 处理分页数据时, 各字段值.
func (o *with) paging(total int64, limit, page int) WithPaging {
    // 1. 默认值.
    //    查询结果只有1页时的参数.
    p := WithPaging{
        First: 1, Before: 1, Current: page, Next: 1, Last: 1,
        Limit: limit, TotalPages: 1, TotalItems: total,
    }
    // 1.1 最小页码.
    if page < 1 {
        p.Current = 1
    }
    // 1.1 最小数量.
    if limit < 1 {
        p.Limit = 1
    }

    // 2. 计算分页.
    t := int(total)
    if t > p.Limit {
        // 总页数.
        p.TotalPages = t / p.Limit
        if t%p.Limit > 0 {
            p.TotalPages += 1
        }

        // 最大页码.
        p.Last = p.TotalPages

        // 前页锚点.
        if p.Current > 2 {
            p.Before = p.Current - 1
        }

        // 下页锚点.
        if p.Current < p.Last {
            p.Next = p.Current + 1
        } else {
            p.Next = p.Last
        }
    }
    return p
}
