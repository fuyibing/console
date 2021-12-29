// author: wsfuyibing <websearch@163.com>
// date: 2021-12-27

package child

type Child2 struct {
    Id2 int
    Child3 Child3
}

type Child3 struct {
    Id3 int
}

type Child struct {
    ChildId int `json:"child_id" label:"子集ID"`
    Id int
    Child2 Child2
}
