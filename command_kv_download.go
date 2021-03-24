// author: wsfuyibing <websearch@163.com>
// date: 2021-03-23

package console

// Download kv command struct.
type kvDownloadCommand struct {
	command *Command
}

// Handle download kv.
func (o *kvDownloadCommand) handler(cs *Console) error {
	return nil
}

// New download kv.
func newKvDownloadCommand() *Command {
	// 1. normal.
	c := NewCommand("kvd")
	c.SetDescription("Download kv from consul server")
	// 2. register option.
	c.Add(
		NewOption("addr").SetTag('a').
			SetDescription("Consul server address").
			SetMode(RequiredMode).SetValue(StringValue),
		NewOption("name").SetTag('n').
			SetDescription("Registered name on consul kv").
			SetMode(RequiredMode).SetValue(StringValue),
		NewOption("path").
			SetDescription("Location for downloaded yaml file save").
			SetMode(RequiredMode).SetValue(StringValue),
		NewOption("recursion").SetTag('r').
			SetDescription("Match prefix 'kv://' in yaml").
			SetMode(RequiredMode).SetValue(StringValue),
		NewOption("override").SetTag('o').
			SetDescription("Override if file exists").
			SetMode(RequiredMode).SetValue(StringValue),
	)
	// 3. register handler.
	c.SetHandler((&kvDownloadCommand{command: c}).handler)
	return c
}
