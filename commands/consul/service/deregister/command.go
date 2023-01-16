// author: wsfuyibing <websearch@163.com>
// date: 2023-01-11

// Package deregister
// remove registered service on consul.
package deregister

import (
	"fmt"
	"github.com/fuyibing/console/v3/commands/consul"
	"github.com/fuyibing/console/v3/managers"
	"github.com/hashicorp/consul/api"
)

const (
	CmdDesc = "Remove service from consul"
	CmdName = "service:deregister"
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
		cfg                    = api.DefaultNonPooledConfig()
		keys                   map[string]interface{}
		serviceId, serviceName string
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

	// Service id.
	if serviceId, err = o.Command.GetOption(consul.OptServiceId).ToString(); err != nil {
		return
	}

	// Service name.
	if serviceName, err = o.Command.GetOption(consul.OptServiceName).ToString(); err != nil {
		return
	}

	// Send
	// deregister request.
	keys, err = consul.Client.Deregister(cfg, serviceName, serviceId)
	managers.Output.Map(keys, fmt.Sprintf("Remove service: %v", serviceName))
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
		managers.NewOption(consul.OptServiceId).SetDescription(consul.OptServiceIdDesc).SetMode(managers.ModeRequired),
		managers.NewOption(consul.OptServiceName).SetDescription(consul.OptServiceNameDesc).SetMode(managers.ModeRequired),
	)
	return o
}

// New
// create and return instance.
//
//   go run main.go service:deregister \
//     --addr=consul.example.com \
//     --scheme=https \
//     --service-id=myapp-hash-string \
//     --service-name=myapp
func New() (managers.Command, error) {
	o := (&Command{Name: CmdName}).
		InitField().
		InitOption()

	return o.Command, o.Err
}
