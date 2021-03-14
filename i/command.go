// author: wsfuyibing <websearch@163.com>
// date: 2021-02-25

package i

// Command interface.
type ICommand interface {
	GetDescription() string
	GetName() string
	GetOption(name string) IOption
	GetOptions() map[string]IOption
	IsDefault() bool
	IsHidden() bool
	Run(console IConsole)
	SetHandler(func(IConsole))
	Usage(console IConsole)
	Validate(args []string) error
}
