// author: wsfuyibing <websearch@163.com>
// date: 2021-02-28

package build

var (
	TypeMapping = map[string]string{
		"float":     "float64",
		"double":    "float64",
		"decimal":   "float64",
		"bigint":    "int64",
		"int":       "int",
		"tinyint":   "int",
		"smallint":  "int",
		"mediumint": "int",
		"time":      "time.Time:time",
		"timestamp": "time.Time:time",
		"datetime":  "time.Time:time",
		"char":      "string",
		"text":      "string",
		"enum":      "string",
		"varchar":   "string",
	}
)

type BeanColumn struct {
	Comment string
	Default string
	Field   string
	Key     string
	Null    string
	Type    string
}

type BeanTable struct {
	Name    string
	Created string `xorm:"Create_time"`
	Updated string `xorm:"Update_time"`
	Comment string
}
