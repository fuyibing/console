// author: wsfuyibing <websearch@163.com>
// date: 2021-03-23

package console

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"
)

type kvItem struct {
	Key   string `json:"Key"`
	Value string `json:"Value"`
}

// Convert base64 value to normal.
func (o *kvItem) Decode() ([]byte, error) {
	str, err := base64.StdEncoding.DecodeString(o.Value)
	if err != nil {
		return nil, err
	}
	return str, nil
}

func (o *kvItem) Yaml() string {
	// strbytes := []byte(o.Value)
	// encoded := base64.StdEncoding.EncodeToString(strbytes)

	decoded, _ := base64.StdEncoding.DecodeString(o.Value)
	return string(decoded)
	// decodestr := string(decoded)
	//
	//
	// base64.NewDecoder().Read()
	//
	// base64.StdEncoding.EncodeToString(strbytes)
	//
	//
	// return o.Value
}

// Download kv command struct.
type kvDownloadCommand struct {
	command    *Command
	consulAddr string
	kvName     string
	kvPath     string
	override   bool
	recursion  bool
}

// Handle after download kv.
func (o *kvDownloadCommand) after(cs *Console) error { return nil }

// Handle before download kv.
func (o *kvDownloadCommand) before(cs *Console) error {
	o.kvName = o.command.GetOption("name").String()
	o.kvPath = o.command.GetOption("path").String()
	o.override = o.command.GetOption("override").Bool()
	o.recursion = o.command.GetOption("recursion").Bool()
	o.consulAddr = o.command.GetOption("addr").String()
	if !regexp.MustCompile(`^https?://`).MatchString(o.consulAddr) {
		o.consulAddr = "http://" + o.consulAddr
	}
	return nil
}

// Handle download kv.
func (o *kvDownloadCommand) handler(cs *Console) error {
	// normal.
	text, err := o.getKvContent(cs, o.kvName)
	if err != nil {
		return err
	}
	// recursion.
	if o.recursion {
		if text, err = o.parseKvContent(cs, text); err != nil {
			return nil
		}
	}
	// split yaml.
	list, err2 := o.splitContent(cs, text)
	if err2 != nil {
		return err2
	}
	// write one by one.
	if err3 := o.writeDir(cs); err3 != nil {
		return err3
	}
	for listName, listText := range list {
		if err4 := o.writeYaml(cs, listName, listText); err4 != nil {
			return err4
		}
	}
	return nil
}

// Get remote kv content.
func (o *kvDownloadCommand) getKvContent(cs *Console, name string) (string, error) {
	text := ""
	// build request.
	req, e1 := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/v1/kv/%s", o.consulAddr, name), nil)
	if e1 != nil {
		return text, e1
	}
	// read response and close body when ended.
	cli := &http.Client{Timeout: time.Duration(2) * time.Second}
	res, e2 := cli.Do(req)
	if e2 != nil {
		return text, e2
	}
	defer func() {
		_ = res.Body.Close()
	}()
	// parse response.
	if res.StatusCode != http.StatusOK {
		return text, fmt.Errorf("%s, KV %s", res.Status, name)
	}
	// read body error.
	body, err3 := ioutil.ReadAll(res.Body)
	if err3 != nil {
		return text, err3
	}
	// parse items struct.
	items := make([]*kvItem, 0)
	if err4 := json.Unmarshal(body, &items); err4 != nil {
		return text, err4
	}
	// parse item string.
	for _, item := range items {
		dec, err5 := item.Decode()
		if err5 != nil {
			return "", err3
		}
		text += strings.TrimSpace(string(dec))
	}
	return text, nil
}

// Parse kv content use recursion.
func (o *kvDownloadCommand) parseKvContent(cs *Console, origin string) (text string, err error) {
	text = RegexpKvRecursion.ReplaceAllStringFunc(origin, func(s string) string {
		m := RegexpKvRecursion.FindStringSubmatch(s)
		sub, se := o.getKvContent(cs, m[1])
		if se != nil {
			err = se
			return s
		}
		return sub
	})
	return
}

// Split with new line.
func (o *kvDownloadCommand) splitContent(cs *Console, text string) (map[string]string, error) {
	key := ""
	res := make(map[string]string)
	regEmptyLine := regexp.MustCompile(`^\s+$`)
	regYamlFile := regexp.MustCompile(`^\s*([_a-zA-Z0-9-]+.yaml)\s*:\s*$`)
	regYamlTable := regexp.MustCompile(`^s+`)
	for _, line := range strings.Split(text, "\n") {
		// ignore empty line.
		if regEmptyLine.MatchString(line) {
			continue
		}
		// define file repeat.
		if m := regYamlFile.FindStringSubmatch(line); len(m) > 0 {
			key = m[1]
			if _, ok := res[key]; ok {
				return nil, fmt.Errorf("yaml define repeated: %s", key)
			}
			str := fmt.Sprintf("# name: %s.\n", o.kvName)
			str += fmt.Sprintf("# path: %s.\n", o.kvPath)
			str += fmt.Sprintf("# date: %s.\n", time.Now().Format("2006-01-02 15:03:04"))
			res[key] = str
			continue
		}
		// not defined.
		if key == "" {
			return nil, fmt.Errorf("yaml file not defined: %s", key)
		}
		// append define.
		line = regYamlTable.ReplaceAllStringFunc(line, func(s string) string {
			return strings.ReplaceAll(s, "\t", "  ")
		})
		if size := len(line); size > 2 {
			res[key] += line[2:] + "\n"
		}
	}
	return res, nil
}

// Write directory.
func (o *kvDownloadCommand) writeDir(cs *Console) error {
	f, err := os.Stat(o.kvPath)
	if err != nil {
		if os.IsExist(err) {
			return err
		}
	}
	if f.IsDir() {
		return nil
	}
	return os.Mkdir(o.kvPath, os.ModePerm)
}

// Write content.
func (o *kvDownloadCommand) writeYaml(cs *Console, name, text string) error {
	fp, err := os.OpenFile(o.kvPath+"/"+name, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, os.ModePerm)
	if err != nil {
		return err
	}
	defer func() {
		_ = fp.Close()
	}()
	if _, err = fp.WriteString(text); err != nil {
		return err
	}
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
		NewOption("override").SetTag('o').
			SetDescription("Override if file exists").
			SetMode(OptionalMode).SetValue(NullValue),
		NewOption("path").
			SetDescription("Location for downloaded yaml file save").
			SetDefaultValue("tmp").
			SetMode(OptionalMode).SetValue(StringValue),
		NewOption("recursion").SetTag('r').
			SetDescription("Match prefix 'kv://' in yaml").
			SetMode(OptionalMode).SetValue(NullValue),
	)
	// 3. register handler.
	o := &kvDownloadCommand{command: c}
	c.SetHandlerBefore(o.before).SetHandler(o.handler).SetHandlerAfter(o.after)
	return c
}
