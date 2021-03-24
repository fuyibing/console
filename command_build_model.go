// author: wsfuyibing <websearch@163.com>
// date: 2021-03-23

package console

// Build model command struct.
type buildModelCommand struct {
	command *Command
}

// Handle after.
func (o *buildModelCommand) after(cs *Console) error { return nil }

// Handle before.
func (o *buildModelCommand) before(cs *Console) error { return nil }

// Handle build command.
func (o *buildModelCommand) handler(cs *Console) error {
	return nil
}

// New build model command.
func newBuildModelCommand() *Command {
	// base.
	c := NewCommand("bm")
	c.SetDescription("Build model for application")
	// options.
	c.Add(
		NewOption("json").
			SetMode(OptionalMode).SetDefaultValue("same").
			SetDescription("Export as json string format, accept: auto,snake,camel, default: auto"),
		NewOption("name").SetTag('n').
			SetDescription(""),
		NewOption("override").SetTag('o').
			SetMode(OptionalMode).SetValue(NullValue),
		NewOption("table").SetTag('t').
			SetMode(OptionalMode),
	)
	// register handler.
	o := &buildModelCommand{command: c}
	c.SetHandlerBefore(o.before).SetHandler(o.handler).SetHandlerAfter(o.after)
	return c
}
