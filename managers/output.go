// author: wsfuyibing <websearch@163.com>
// date: 2023-01-16

package managers

import (
	"fmt"
	"sort"
	"strings"
)

var (
	// Output
	// manager instance.
	Output OutputManager
)

type (
	// OutputManager
	// manager interface.
	OutputManager interface {
		Map(keys map[string]interface{}, desc string)
	}

	output struct{}
)

// Map
// format print.
func (o *output) Map(keys map[string]interface{}, desc string) {
	var (
		format       string
		index, width = 0, 0
		list         = make([]string, 0)
	)

	// Range
	// key to list and execute maximum width.
	for k, _ := range keys {
		list = append(list, k)

		// Generate
		// key maximum characters width.
		if w := len(k); width < w {
			width = w
		}
	}

	// Build
	// map format and sorts by string.
	format = fmt.Sprintf("%%-%ds  - %%v", width)
	sort.Strings(list)

	// Range key.
	for _, k := range list {
		if v, ok := keys[k]; ok {
			// Print description.
			if index++; index == 1 {
				o.println(desc)
				o.println(strings.Repeat("-", 80))
			}

			// Print key.
			o.println(format, k, v)
		}
	}
}

// Init output instance.
func (o *output) init() *output {
	return o
}

// Print contents.
func (o *output) println(text string, args ...interface{}) {
	println(fmt.Sprintf(text, args...))
}
