// author: wsfuyibing <websearch@163.com>
// date: 2023-01-12

// Package managers
// implements support for command management.
package managers

import (
	"sync"
)

func init() {
	new(sync.Once).Do(func() {
		Output = (&output{}).init()
	})
}
