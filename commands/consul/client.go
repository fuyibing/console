// author: wsfuyibing <websearch@163.com>
// date: 2023-01-16

package consul

import (
	"fmt"
	"github.com/hashicorp/consul/api"
	"os"
	"strings"
	"time"
)

var (
	// Client
	// instance for consul.
	Client *ClientManager
)

type (
	// ClientManager
	// for consul agent manager.
	ClientManager struct{}
)

// Deregister
// remove service of consul.
func (o *ClientManager) Deregister(cfg *api.Config, serviceName, serviceId string) (res map[string]interface{}, err error) {
	var (
		cli  *api.Client
		idx  = 0
		key  string
		list []*api.CatalogService
	)

	// Prepare
	// download results.
	res = make(map[string]interface{})

	// Build
	// consul api client.
	if cli, err = api.NewClient(cfg); err != nil {
		res[serviceName] = err
		return
	}

	// List service
	// by name.
	if list, _, err = cli.Catalog().Service(serviceName, "", nil); err != nil {
		res[serviceName] = err
		return
	}

	// Range service.
	for _, item := range list {
		if serviceId != "*" && serviceId != item.ServiceID {
			continue
		}

		// Build
		// result index.
		idx++
		key = fmt.Sprintf("index=%d, node=%v, service-id=%v", idx, item.Node, item.ServiceID)

		// Send deregister request.
		if _, de := cli.Catalog().Deregister(&api.CatalogDeregistration{
			Node: item.Node, ServiceID: item.ServiceID,
		}, nil); de != nil {
			res[key] = de
		} else {
			res[key] = "deleted"
		}
	}

	return
}

// Download
// remote configuration from consul and store as local files.
func (o *ClientManager) Download(cfg *api.Config, key, path string, override bool) (res map[string]interface{}, err error) {
	var (
		cli  *api.Client
		text string
	)

	// Prepare
	// download results.
	res = make(map[string]interface{})

	// Build
	// consul api client.
	if cli, err = api.NewClient(cfg); err != nil {
		return
	}

	// Read
	// key contents from consul.
	if text, err = o.keyReader(cli, res, key); err == nil {
		for _, rows := range o.keySplit(text) {
			if len(rows) > 1 {
				if err = o.keySave(res, override, path, rows[0], rows[1:]); err != nil {
					break
				}
			}
		}
	}
	return
}

// Register
// add new service to consul.
func (o *ClientManager) Register(cfg *api.Config, req *api.AgentServiceRegistration) (res map[string]interface{}, err error) {
	var (
		cli *api.Client
	)

	// Prepare
	// download results.
	res = make(map[string]interface{})

	// Catch
	// when process end.
	defer func() {
		if err != nil {
			res[req.Name] = err
		} else {
			res[req.Name] = "succeed"
		}
	}()

	// Build
	// consul api client.
	if cli, err = api.NewClient(cfg); err != nil {
		err = cli.Agent().ServiceRegister(req)
	}
	return
}

// Upload
// read local config files contents and put to consul.
func (o *ClientManager) Upload(cfg *api.Config, key, path string) (res map[string]interface{}, err error) {
	var (
		cli  *api.Client
		text string
	)

	// Prepare
	// uploaded results.
	res = make(map[string]interface{})

	// Read contents.
	if text, err = o.ymlReader(res, path); err != nil {
		return
	}

	// Build
	// consul api client.
	if cli, err = api.NewClient(cfg); err != nil {
		return
	}

	// Build params and send upload request.
	_, err = cli.KV().Put(&api.KVPair{
		Key:   key,
		Value: []byte(text),
	}, nil)

	return
}

// /////////////////////////////////////////////////////////////
// Access and construct methods
// /////////////////////////////////////////////////////////////

// Init
// client instance.
func (o *ClientManager) init() *ClientManager {
	return o
}

// Read contents
// from consul.
func (o *ClientManager) keyReader(c *api.Client, res map[string]interface{}, key string) (text string, err error) {
	var (
		k  = fmt.Sprintf("%v", key)
		kp *api.KVPair
	)

	defer func() {
		if err != nil {
			res[k] = err
		} else {
			res[k] = "succeed"
		}
	}()

	// Get contents by key.
	if kp, _, err = c.KV().Get(key, nil); err != nil {
		return
	}

	// Return
	// if not found.
	if kp == nil {
		err = fmt.Errorf("not found")
		return
	}

	// Replace variables like `kv://name`
	text = RegexDepth.ReplaceAllStringFunc(string(kp.Value), func(s string) string {
		if m := RegexDepth.FindStringSubmatch(s); len(m) == 2 {
			if sr, se := o.keyReader(c, res, m[1]); se == nil {
				return sr
			}
		}
		return s
	})

	return
}

// Split contents
// by line separator.
func (o *ClientManager) keySplit(text string) (res [][]string) {
	var (
		lk string
		ls = make([]string, 0)
	)

	res = make([][]string, 0)

	// Range lines.
	for _, s := range strings.Split(text, "\n") {
		if strings.TrimSpace(s) == "" {
			continue
		}

		// Find file name.
		if m := RegexFilenameRemote.FindStringSubmatch(s); len(m) == 2 {
			if lk != "" && len(ls) > 1 {
				res = append(res, ls)
			}

			// Reset file pointers.
			lk = m[1]
			ls = []string{lk}
			continue
		}

		// Find line.
		if len(s) > 2 {
			ls = append(ls, s[2:])
		}
	}

	// Collect end lines.
	if lk != "" && len(ls) > 1 {
		res = append(res, ls)
	}

	return
}

// Save contents
// to local yaml files.
func (o *ClientManager) keySave(res map[string]interface{}, override bool, path, name string, line []string) (err error) {
	var (
		fp       *os.File
		fullPath = fmt.Sprintf("%s/%s", path, name)
		k        = fmt.Sprintf("%s", fullPath)
	)

	// Not override.
	if !override {
		if s, se := os.Stat(path); se == nil && !s.IsDir() {
			res[k] = "ignored"
			return
		}
	}

	// Close when end.
	defer func() {
		if fp != nil {
			err = fp.Close()
		}
		if err != nil {
			res[k] = err
		} else {
			res[k] = "succeed"
		}
	}()

	// Open config file.
	if fp, err = os.OpenFile(fullPath, os.O_TRUNC|os.O_WRONLY|os.O_CREATE, os.ModePerm); err != nil {
		return
	}

	// Write string.
	_, err = fp.WriteString(strings.Join(line, "\n"))
	return
}

// Read contents
// from local yaml files.
func (o *ClientManager) ymlReader(res map[string]interface{}, path string) (text string, err error) {
	var (
		buf      []byte
		ds       []os.DirEntry
		fullPath string
		now      = time.Now().Format("2006-01-02 15:04:05.999")
	)

	// Return error
	// if read directory failed.
	if ds, err = os.ReadDir(path); err != nil {
		res[path] = err
		return
	}

	// Range file.
	for _, d := range ds {
		// Ignore directory or not yaml file.
		if d.IsDir() || !RegexFilename.MatchString(d.Name()) {
			continue
		}

		// if read failed.
		// Return error
		fullPath = fmt.Sprintf("%s/%s", path, d.Name())
		if buf, err = os.ReadFile(fullPath); err != nil {
			res[fullPath] = err
			return
		}

		// Add file name.
		res[fullPath] = "succeed"
		text += fmt.Sprintf("%s: # uploaded: %v\n", d.Name(), now)

		// Range lines.
		for _, s := range strings.Split(string(buf), "\n") {
			if strings.TrimSpace(s) == "" {
				continue
			}

			// Append line.
			text += fmt.Sprintf("  %s\n", s)
		}
	}

	return
}
