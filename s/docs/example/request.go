// author: wsfuyibing <websearch@163.com>
// date: 2021-12-27

package example

import (
    "gsjob/tests/example/ano"
    "gsjob/tests/example/child"
)

type Request struct {
    ano.Request

    UserId int `json:"user_id" label:"用户ID"`

    Child child.Child

    // Arr *[]*[]*int `json:"arr" mock:"1001"`

    // ReqInt11 int
    // ReqInt12 *int
    // ReqInt13 []int
    // ReqInt14 []*int
    // ReqInt15 *[]int
    // ReqInt16 *[]*int

    // ChildBlockStruct         child.Child
    // ChildBlockPointer        *child.Child
    // ChildPointerSlice *[]child.Child `json:"child_pointer_slice"`
    // ChildPointerSlicePointer *[]*child.Child
    // ChildSliceStruct         []child.Child
    // ChildSlicePointer        []*child.Child
}
