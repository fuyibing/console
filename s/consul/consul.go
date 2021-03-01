// author: wsfuyibing <websearch@163.com>
// date: 2021-02-27

package consul

import (
	"errors"
	"fmt"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/hashicorp/consul/api"

	"github.com/fuyibing/console/v2/i"
)

var (
	regexpParseContent = regexp.MustCompile(`kv://([/_a-zA-Z0-9]+)`)
	regexpKvComment    = regexp.MustCompile(`^\s*#`)
	regexpKvEmptyLine  = regexp.MustCompile(`^\s*$`)
	regexpKvFile       = regexp.MustCompile(`^([_a-zA-Z0-9\-\.]+.yaml)\s*:\s*`)
)

type kv struct {
	command i.ICommand
	console i.IConsole
}

// New kv manager instance.
func Manager(console i.IConsole, command i.ICommand) *kv {
	o := &kv{command: command, console: console}
	return o
}

// Register service.
func (o *kv) Register() {}

// Remove service.
func (o *kv) Deregister() {}

// Download kv.
func (o *kv) Download() {
	// Base client.
	cli, err := o.client()
	if err != nil {
		o.console.PrintError(errors.New(fmt.Sprintf("Command %s: read consul fail: %v.", o.command.GetName(), err)))
		return
	}
	// Read content.
	str := ""
	content := ""
	for _, name := range strings.Split(o.command.GetOption("name").ToString(), ",") {
		str, err = o.readKvContents(cli, name)
		if err != nil {
			o.console.PrintError(err)
			return
		}
		content += "\n\n" + str
	}
	// Parse depth.
	if o.command.GetOption("parse").ToBool() {
		content, err = o.parseKvContents(cli, content)
		if err != nil {
			o.console.PrintError(err)
			return
		}
	}
	// Split content.
	var res = make(map[string][]string)
	if res, err = o.splitKvContents(content); err != nil {
		o.console.PrintError(err)
		return
	}
	// Check file exist if not override.
	if !o.command.GetOption("override").ToBool() {
		for name, _ := range res {
			if err = o.checkFileExist(name); err != nil {
				o.console.PrintError(err)
				return
			}
		}
	}
	// Write file.
	for name, lines := range res {
		if err = o.writeConfigFile(name, lines); err != nil {
			o.console.PrintError(err)
			return
		}
	}
}

// Upload kv.
func (o *kv) Upload() {

}

// Check file exists.
func (o *kv) checkFileExist(name string) error {
	path := o.command.GetOption("path").ToString()
	// Check path.
	if err := o.checkPath(path); err != nil {
		return err
	}
	// Check file.
	file := fmt.Sprintf("%s/%s", path, name)
	if f, err := os.OpenFile(file, os.O_RDONLY, os.ModePerm); err == nil {
		_ = f.Close()
		return errors.New(fmt.Sprintf("Command %s: file exist: %s", o.command.GetName(), file))
	}
	// File not found.
	return nil
}

// Check directory.
func (o *kv) checkPath(path string) error {
	if err := os.MkdirAll(path, os.ModePerm); err != nil {
		return errors.New(fmt.Sprintf("Command %s: check directory error: %s", o.command.GetName(), path))
	}
	return nil
}

// KV Client.
func (o *kv) client() (*api.Client, error) {
	return api.NewClient(&api.Config{
		Address:  o.command.GetOption("addr").ToString(),
		Scheme:   "http",
		WaitTime: time.Duration(5) * time.Second,
	})
}

// Parse kv content from consul kv and replace depth.
func (o *kv) parseKvContents(cli *api.Client, content string) (string, error) {
	var err error
	var str = ""
	content = regexpParseContent.ReplaceAllStringFunc(content, func(s string) string {
		if m := regexpParseContent.FindStringSubmatch(s); len(m) == 2 {
			if str, err = o.readKvContents(cli, m[1]); err == nil {
				return str
			}
		}
		return s
	})
	// Completed.
	return content, err
}

// Read kv content from consul kv.
func (o *kv) readKvContents(cli *api.Client, name string) (string, error) {
	p, _, err := cli.KV().Get(name, nil)
	if err != nil {
		return "", errors.New(fmt.Sprintf("Command %s: read kv content fail: %v.", o.command.GetName(), err))
	}
	if p == nil {
		return "", errors.New(fmt.Sprintf("Command %s: kv not found: %s.", o.command.GetName(), name))
	}
	o.console.Info("Command %s: read kv contents  : %s.", o.command.GetName(), name)
	return string(p.Value), nil
}

// Split config content from consul kv.
func (o *kv) splitKvContents(content string) (map[string][]string, error) {
	var name = ""
	var res = make(map[string][]string)
	for _, line := range strings.Split(content, "\n") {
		// ignore comment.
		if regexpKvComment.MatchString(line) || regexpKvEmptyLine.MatchString(line) {
			continue
		}
		// file name.
		if m := regexpKvFile.FindStringSubmatch(line); len(m) == 2 {
			name = m[1]
			// repeat: kv format.
			if _, ok := res[name]; ok {
				return nil, errors.New(fmt.Sprintf("Command %s: file name exist - kv format error", o.command.GetName()))
			}
			// init map.
			res[name] = make([]string, 0)
			continue
		}
		// append line.
		if chs := len(line); chs <= 2 {
			continue
		}
		// last file name.
		if name == "" {
			return nil, errors.New(fmt.Sprintf("Command %s: unknown file name for append", o.command.GetName()))
		}
		res[name] = append(res[name], line[2:])
	}
	return res, nil
}

// Write content to config file.
func (o *kv) writeConfigFile(name string, lines []string) error {
	var err error
	var handler *os.File
	// text format.
	var text = ""
	text += fmt.Sprintf("# {%s}\n", name)
	text += fmt.Sprintf("# date: %s.\n", time.Now().Format("2006-01-02 15:04:05"))
	for _, line := range lines {
		text += fmt.Sprintf("%s\n", line)
	}
	// file name.
	file := fmt.Sprintf("%s/%s", o.command.GetOption("path").ToString(), name)
	// open file.
	handler, err = os.OpenFile(file, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, os.ModePerm)
	if err != nil {
		return errors.New(fmt.Sprintf("Command %s: create file error: %s", o.command.GetName(), file))
	}
	// close when end.
	defer func() {
		_ = handler.Close()
	}()
	// write contents.
	if _, err = handler.WriteString(text); err != nil {
		return errors.New(fmt.Sprintf("Command %s: write file error: %s", o.command.GetName(), file))
	}
	o.console.Info("Command %s: write config file : %s.", o.command.GetName(), file)
	return nil
}
