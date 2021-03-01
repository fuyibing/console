// author: wsfuyibing <websearch@163.com>
// date: 2021-02-25

package i

// Console interface.
type IConsole interface {
	Add(cs ...ICommand)
	Del(cs ...ICommand)
	GetCommand(name string) ICommand
	GetCommands() map[string]ICommand
	GetNames() []string
	Info(text string, args ... interface{})
	PrintCommandItem(n int, cmd ICommand, end bool)
	PrintCommandMore(cmd ICommand)
	PrintError(err error)
	PrintOptionItem(n int, opt IOption, end bool)
	PrintUsage(cmd ICommand)
	Run(args ...string)
}
