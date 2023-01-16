// author: wsfuyibing <websearch@163.com>
// date: 2023-01-11

// Package upload
// scan and read local config files and put to consul kv
// storage.
package upload

import (
	"github.com/fuyibing/console/v3/commands/consul"
	"github.com/fuyibing/console/v3/managers"
	"github.com/hashicorp/consul/api"
)

const (
	CmdDesc = "Upload local config file to consul kv"
	CmdName = "kv:upload"
)

type Command struct {
	Command managers.Command
	Err     error
	Name    string
}

// Handle
// send upload request.
func (o *Command) Handle(_ managers.Manager, _ managers.Arguments) (err error) {
	var (
		cfg       = api.DefaultNonPooledConfig()
		key, path = "", ""
		keys      map[string]interface{}
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

	// Consul key name.
	if key, err = o.Command.GetOption(consul.OptKey).ToString(); err != nil {
		return
	}

	// Config storage path.
	if path, err = o.Command.GetOption(consul.OptPath).ToString(); err != nil {
		return
	}

	// Send upload request.
	keys, err = consul.Client.Upload(cfg, key, path)
	managers.Output.Map(keys, "Consul key uploaded results")
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
