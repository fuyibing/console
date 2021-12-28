// author: wsfuyibing <websearch@163.com>
// date: 2021-12-27

package ano

type Request struct {
    AnonymousId  int   `json:"anonymous_id" validate:"required,min=5" label:"匿名ID" desc:"关于匿名ID的描述信息"`
    AnonymousIds []int `json:"anonymous_ids" validate:"required" label:"匿名ID列表" desc:"关于匿名ID列表的描述信息"`
}
