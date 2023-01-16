// author: wsfuyibing <websearch@163.com>
// date: 2023-01-11

// Package download
// read configurations from consul kv storage
// and store to local yaml files.
package download

import (
	"github.com/fuyibing/console/v3/commands/consul"
	"github.com/fuyibing/console/v3/managers"
	"github.com/hashicorp/consul/api"
)

const (
	CmdDesc = "Download config from consul and save to local"
	CmdName = "kv:download"
)

// Command
// for consul kv download.
type Command struct {
	Command managers.Command
	Err     error
	Name    string
}

// Handle
// send download request.
func (o *Command) Handle(_ managers.Manager, _ managers.Arguments) (err error) {
	var (
		cfg       = api.DefaultNonPooledConfig()
		key, path = "", ""
		keys      map[string]interface{}
		override  bool
	)

	// Read address option.
	//
	//   -a consul.example.com
	//   --addr consul.example.com
	//   --addr="consul.example.com"
	if cfg.Address, err = o.Command.GetOption(consul.OptAddr).ToString(); err != nil {
		return
	}

	// Read scheme option.
	//
	//   -s https
	//   --scheme https
	//   --scheme="https"
	if cfg.Scheme, err = o.Command.GetOption(consul.OptScheme).ToString(); err != nil {
		return
	}

	// Read key name option.
	//
	//   -n app/myapp
	//   --name app/myapp
	//   --name="app/myapp"
	if key, err = o.Command.GetOption(consul.OptKey).ToString(); err != nil {
		return
	}

	// Config storage path.
	if path, err = o.Command.GetOption(consul.OptPath).ToString(); err != nil {
		return
	}

	// Config override or not.
	if override, err = o.Command.GetOption(consul.OptOverride).ToBool(); err != nil {
		return
	}

	// Send download request.
	keys, err = consul.Client.Download(cfg, key, path, override)
	managers.Output.Map(keys, "Consul key downloaded results")
	return
}

// InitField
// initialize command fields.
func (o *Command) InitField() *Command {
	o.Command = managers.NewCommand(o.Name)
	o.Command.SetDescription(CmdDesc).SetHandler(o.Handle)
	return o
}

// InitOption
// initialize command option.
func (o *Command) InitOption() *Command {
	o.Err = o.Command.AddOption(
		managers.NewOption(consul.OptAddr).SetShortName(consul.OptAddrByte).SetDescription(consul.OptAddrDesc).SetMode(managers.ModeRequired),
		managers.NewOption(consul.OptScheme).SetShortName(consul.OptSchemeByte).SetDescription(consul.OptSchemeDesc).SetDefault(consul.OptSchemeDefault),
		managers.NewOption(consul.OptKey).SetShortName(consul.OptKeyByte).SetDescription(consul.OptKeyDesc).SetMode(managers.ModeRequired),
		managers.NewOption(consul.OptOverride).SetShortName(consul.OptOverrideByte).SetDescription(consul.OptOverrideDesc).SetDefault(consul.OptOverrideDefault).SetValueType(managers.ValueTypeBoolean),
		managers.NewOption(consul.OptPath).SetShortName(consul.OptPathByte).SetDescription(consul.OptPathDesc).SetDefault(consul.OptPathDefault),
	)
	return o
}

// New
// create and return instance.
func New() (managers.Command, error) {
	o := (&Command{Name: CmdName}).
		InitField().
		InitOption()

	return o.Command, o.Err
}
