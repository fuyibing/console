// author: wsfuyibing <websearch@163.com>
// date: 2023-01-12

// Package consul
// work for service register and find.
package consul

import (
	"regexp"
	"sync"
)

const (
	OptAddr     = "addr"
	OptAddrByte = 'a'
	OptAddrDesc = "Consul server address, such as: 127.0.0.1, consul.example.com"

	OptKey     = "name"
	OptKeyByte = 'n'
	OptKeyDesc = "Consul key name"

	OptScheme        = "scheme"
	OptSchemeByte    = 's'
	OptSchemeDefault = "http"
	OptSchemeDesc    = "Consul server scheme, accept http or https"

	OptOverride        = "override"
	OptOverrideByte    = 'o'
	OptOverrideDefault = false
	OptOverrideDesc    = "Override config files if exists"

	OptPath        = "path"
	OptPathByte    = 'p'
	OptPathDefault = "./config"
	OptPathDesc    = "Config file storage location"

	OptServiceAddr     = "service-addr"
	OptServiceAddrDesc = "Consul service address, such as: 172.16.0.100, app.example.com"

	OptServiceId     = "service-id"
	OptServiceIdDesc = "Consul service id, such as: myapp-hash"

	OptServiceName     = "service-name"
	OptServiceNameDesc = "Consul service name, such as: myapp"

	OptServicePort     = "service-port"
	OptServicePortDesc = "Consul service port, such as: 80, 8080"
)

var (
	RegexDepth          = regexp.MustCompile(`kv://([._a-zA-Z0-9-/]+)`)
	RegexFilename       = regexp.MustCompile(`^([a-zA-Z][._a-zA-Z0-9-]*\.ya?ml)$`)
	RegexFilenameRemote = regexp.MustCompile(`^([a-zA-Z][._a-zA-Z0-9-]*\.ya?ml):[^\n]*`)
)

func init() {
	new(sync.Once).Do(func() {
		Client = (&ClientManager{}).init()
	})
}
