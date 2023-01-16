// author: wsfuyibing <websearch@163.com>
// date: 2023-01-11

package managers

import (
	"fmt"
	"runtime"
	"strings"
)

type (
	// Command
	// operation interface.
	Command interface {
		AddOption(opts ...Option) error
		GetDescription() string
		GetHidden() bool
		GetName() string
		GetOption(key string) Option
		GetOptions() map[string]Option
		Run(manager Manager, arguments Arguments) error
		SetDescription(s string) Command
		SetHandler(handler CommandHandler) Command
		SetHidden(b bool) Command
	}

	// CommandHandler
	// callable handler on command.
	CommandHandler func(manager Manager, arguments Arguments) error

	command struct {
		Handler           CommandHandler
		Hidden            bool
		Name, Description string
		OptionKeys        map[string]string
		OptionMapper      map[string]Option
	}
)

func NewCommand(name string) Command {
	return (&command{
		Name:         name,
		OptionKeys:   make(map[string]string),
		OptionMapper: make(map[string]Option),
	}).initFields()
}

// /////////////////////////////////////////////////////////////
// Interface methods
// /////////////////////////////////////////////////////////////

func (o *command) AddOption(opts ...Option) error            { return o.addOption(opts) }
func (o *command) GetDescription() string                    { return o.Description }
func (o *command) GetHidden() bool                           { return o.Hidden }
func (o *command) GetName() string                           { return o.Name }
func (o *command) GetOption(key string) Option               { return o.getOption(key) }
func (o *command) GetOptions() map[string]Option             { return o.OptionMapper }
func (o *command) Run(m Manager, a Arguments) error          { return o.run(m, a) }
func (o *command) SetDescription(s string) Command           { o.Description = s; return o }
func (o *command) SetHandler(handler CommandHandler) Command { o.Handler = handler; return o }
func (o *command) SetHidden(b bool) Command                  { o.Hidden = b; return o }

// /////////////////////////////////////////////////////////////
// Access and constructor
// /////////////////////////////////////////////////////////////

func (o *command) addOption(opts []Option) error {
	for _, opt := range opts {
		if opt == nil {
			continue
		}

		// Set mapper.
		o.OptionMapper[opt.GetName()] = opt

		// Full name mapper.
		o.OptionKeys[opt.GetName()] = opt.GetName()

		// Short name mapper.
		if s := opt.GetShortName(); s != "" {
			o.OptionKeys[opt.GetShortName()] = opt.GetName()
		}
	}
	return nil
}

func (o *command) getOption(key string) Option {
	if k, exists := o.OptionKeys[key]; exists {
		if v, ok := o.OptionMapper[k]; ok {
			return v
		}
	}
	return nil
}

func (o *command) initFields() *command {
	return o
}

func (o *command) run(m Manager, a Arguments) (err error) {
	if o.Handler == nil {
		err = fmt.Errorf("command handler not defined: %s", o.Name)
		return
	}

	// Catch
	// command runner panic.
	defer func() {
		if r := recover(); r != nil {
			es := []string{
				fmt.Sprintf("command panic on %v: %v", o.Name, r),
			}

			for i := 0; ; i++ {
				if _, f, l, g := runtime.Caller(i); g {
					es = append(es, fmt.Sprintf("%s:%d", strings.TrimSpace(f), l))
					continue
				}
				break
			}

			err = fmt.Errorf("%s", strings.Join(es, "\n"))
		}
	}()

	// Call handler.
	err = o.Handler(m, a)
	return
}
