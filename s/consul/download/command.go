// author: wsfuyibing <websearch@163.com>
// date: 2021-02-25

package download

import (
	"github.com/fuyibing/console/v2/base"
	"github.com/fuyibing/console/v2/i"
	"github.com/fuyibing/console/v2/s/consul"
)

const (
	Description = "Read configuration from consul kv and split with specified line, " +
		"then extract content to local yaml file. Application read extracted yaml file when start."
	Name = "kv:download"
)

// Command struct.
type command struct {
	base.Command
}

func New() i.ICommand {
	// normal.
	o := &command{}
	o.Initialize()
	o.SetDescription(Description)
	o.SetName(Name)
	// consul addr.
	//   -a 127.0.0.1:8500
	//   --addr=127.0.0.1:8500
	o.Add(base.NewOption(i.OptionalMode, i.StrValue).SetDefaultValue("127.0.0.1:8500").SetName("addr").SetShortName("a").SetDescription("consul server address"))
	// consul name.
	//   -n goapps/demo
	//   --name=goapps/demo
	o.Add(base.NewOption(i.RequiredMode, i.StrValue).SetName("name").SetShortName("n").SetDescription("name of consul consul"))
	// extract to specified directory.
	//   -p ./tmp
	//   --path=./tmp
	o.Add(base.NewOption(i.OptionalMode, i.StrValue).SetDefaultValue("./tmp").SetName("path").SetShortName("p").SetDescription("extract yaml file to specified directory"))
	// parse content.
	//   --parse
	o.Add(base.NewOption(i.OptionalMode, i.BoolValue).SetName("parse").SetDefaultValue("true").SetDescription("Parse depth content"))
	// override if file exist.
	//   -o
	//   --override
	o.Add(base.NewOption(i.OptionalMode, i.BoolValue).SetName("override").SetShortName("o").SetDescription("override if file exist"))
	// n. prepared.
	return o
}

// Run download.
func (o *command) Run(console i.IConsole) {
	consul.Manager(console, o).Download()
}
