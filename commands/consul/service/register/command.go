// author: wsfuyibing <websearch@163.com>
// date: 2023-01-11

// Package register
// build service and register to consul.
package register

import (
	"fmt"
	"github.com/fuyibing/console/v3/commands/consul"
	"github.com/fuyibing/console/v3/managers"
	"github.com/hashicorp/consul/api"
)

const (
	CmdDesc = "Register service to consul"
	CmdName = "service:register"
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
		cfg  = api.DefaultNonPooledConfig()
		keys map[string]interface{}
		port int64
		req  = &api.AgentServiceRegistration{}
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

	// Service address.
	if req.Address, err = o.Command.GetOption(consul.OptServiceAddr).ToString(); err != nil {
		return
	}

	// Service id.
	if req.ID, err = o.Command.GetOption(consul.OptServiceId).ToString(); err != nil {
		return
	}

	// Service name.
	if req.Name, err = o.Command.GetOption(consul.OptServiceName).ToString(); err != nil {
		return
	}

	// Service port.
	if port, err = o.Command.GetOption(consul.OptServicePort).ToInt(); err != nil {
		return
	} else {
		req.Port = int(port)
	}

	// Send
	// register request.
	keys, err = consul.Client.Register(cfg, req)
	managers.Output.Map(keys, fmt.Sprintf("Register service: %s", req.Name))
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
		managers.NewOption(consul.OptServiceAddr).SetDescription(consul.OptServiceAddrDesc).SetMode(managers.ModeRequired),
		managers.NewOption(consul.OptServiceId).SetDescription(consul.OptServiceIdDesc).SetMode(managers.ModeRequired),
		managers.NewOption(consul.OptServiceName).SetDescription(consul.OptServiceNameDesc).SetMode(managers.ModeRequired),
		managers.NewOption(consul.OptServicePort).SetDescription(consul.OptServicePortDesc).SetMode(managers.ModeRequired).SetValueType(managers.ValueTypeInteger),
	)
	return o
}

// New
// create and return instance.
//
//   go run main.go service:register \
//     --addr=consul.example.com \
//     --scheme=https \
//     --service-addr=127.0.0.1 \
//     --service-port=8080 \
//     --service-id=myapp-hash-string \
//     --service-name=myapp
func New() (managers.Command, error) {
	o := (&Command{Name: CmdName}).
		InitField().
		InitOption()

	return o.Command, o.Err
}
