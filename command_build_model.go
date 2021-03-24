// author: wsfuyibing <websearch@163.com>
// date: 2021-03-23

package console

// Build model command struct.
type buildModelCommand struct {
	command *Command
}

// Handle build model command.
func (o *buildModelCommand) handler(cs *Console) error {
	return nil
}

// New build model command.
func newBuildModelCommand() *Command {
	c := NewCommand("bm")
	c.SetDescription("Build model for application")
	o := &helpCommand{command: c}
	c.SetHandler(o.handler)
	return c
}
