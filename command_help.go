// author: wsfuyibing <websearch@163.com>
// date: 2021-03-23

package console

import (
	"fmt"
	"os"
)

// Help command struct.
type helpCommand struct {
	command *Command
}

// Handle help command.
func (o *helpCommand) handler(cs *Console) error {
	args := os.Args
	if argc := len(args); argc > 2 && RegexpArgumentCommandName.MatchString(args[2]) {
		if c := cs.GetCommand(args[2]); c != nil {
			return o.listOption(cs, c)
		}
		return fmt.Errorf("command name not defined: %s", args[2])
	}
	return o.listCommand(cs)
}

// List commands.
func (o *helpCommand) listCommand(cs *Console) error {
	num := 0
	str := fmt.Sprintf("Usage: %s <COMMAND> [OPTIONS]\n", cs.script)
	for _, command := range cs.GetCommands() {
		if command.isHidden {
			continue
		}
		if num == 0 {
			str += fmt.Sprintf("Commands: \n")
		}
		num++
		str += "    " + fmt.Sprintf("%-22s  %s", command.name, command.description) + "\n"
	}
	str += fmt.Sprintf("Run `%s help <COMMAND>` for more information.\n", cs.script)
	print(str)
	return nil
}

// List option of command.
func (o *helpCommand) listOption(cs *Console, c *Command) error {
	num := 0
	str := fmt.Sprintf("Usage: %s %s [OPTIONS]\n", cs.script, c.name)
	for _, option := range c.GetOptions() {
		if num == 0 {
			str += fmt.Sprintf("Options: \n")
		}
		num++
		s := "    "
		if option.tag > 0 {
			s = fmt.Sprintf("-%s, ", string(option.tag))
		}
		s += fmt.Sprintf("--%s", option.name)
		str += "    " + fmt.Sprintf("%-32s  %s", s, option.description) + "\n"
	}
	print(str)
	return nil
}

// New help command.
func newHelpCommand() *Command {
	c := NewCommand("help")
	c.Default(true).Hidden(true)
	o := &helpCommand{command: c}
	c.SetHandler(o.handler)
	return c
}
