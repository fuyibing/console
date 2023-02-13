// author: wsfuyibing <websearch@163.com>
// date: 2023-01-11

package managers

import (
	"fmt"
	"regexp"
	"strings"
)

const (
	ArgumentsHelp   = "help"
	ArgumentsScript = "go run main.go"
)

var (
	ArgumentsRegexOption        = regexp.MustCompile(`^-`)
	ArgumentsRegexOptionName    = regexp.MustCompile(`^(-+)([a-zA-Z0-9][_a-zA-Z0-9-]*)$`)
	ArgumentsRegexOptionPairs   = regexp.MustCompile(`^-+([a-zA-Z0-9][_a-zA-Z0-9-]*)=(.*)$`)
	ArgumentsRegexScriptBinary  = regexp.MustCompile(`^\./[_a-zA-Z0-9-]+$`)
	ArgumentsRegexScriptWorking = regexp.MustCompile(`^[_a-zA-Z0-9-]+$`)
)

type (
	// Arguments
	// operation interface.
	Arguments interface {
		Get(key string) string
		GetHelpSelector() string
		GetMapper() map[string]string
		GetScript() string
		GetSelector() string
		Has(key string) bool
		Parse(ss ...string) error
	}

	arguments struct {
		Mapper                         map[string]string
		Selector, HelpSelector, Script string
	}
)

func NewArguments() Arguments {
	return &arguments{
		Mapper: make(map[string]string),
	}
}

// /////////////////////////////////////////////////////////////
// Interface methods
// /////////////////////////////////////////////////////////////

func (o *arguments) Get(key string) string        { return o.get(key) }
func (o *arguments) GetHelpSelector() string      { return o.HelpSelector }
func (o *arguments) GetMapper() map[string]string { return o.Mapper }
func (o *arguments) GetScript() string            { return o.Script }
func (o *arguments) GetSelector() string          { return o.Selector }
func (o *arguments) Has(key string) bool          { return o.has(key) }
func (o *arguments) Parse(ss ...string) error     { return o.parse(ss) }

// /////////////////////////////////////////////////////////////
// Access and constructor
// /////////////////////////////////////////////////////////////

func (o *arguments) get(key string) string {
	if s, ok := o.Mapper[key]; ok {
		return s
	}
	return ""
}

func (o *arguments) has(key string) bool {
	if _, ok := o.Mapper[key]; ok {
		return true
	}
	return false
}

func (o *arguments) parse(ss []string) error {
	var (
		keys   = make([]string, 0)
		values = make([]string, 0)
	)

	// Range args.
	for i, s := range ss {
		// Find
		// arguments value.
		if !ArgumentsRegexOption.MatchString(s) {
			switch i {
			case 0:
				o.parseScript(s)
			case 1:
				o.Selector = s
			default:
				if i == 2 && o.Selector == ArgumentsHelp {
					o.HelpSelector = s
				} else {
					values = append(values, s)
				}
			}
			continue
		}

		// Collect tmp
		// to mapper and reset.
		if len(keys) > 0 || len(values) > 0 {
			// Set mapper
			// when option found.
			if err := o.setter(keys, values); err != nil {
				return err
			}

			// Reset tmp.
			keys = make([]string, 0)
			values = make([]string, 0)
		}

		// Find
		// key/value pairs.
		if m := ArgumentsRegexOptionPairs.FindStringSubmatch(s); len(m) == 3 {
			if err := o.setter([]string{m[1]}, []string{m[2]}); err != nil {
				return err
			}
			continue
		}

		// Find
		// argument option.
		if m := ArgumentsRegexOptionName.FindStringSubmatch(s); len(m) == 3 {
			if m[1] == "-" {
				// Short name.
				for _, c := range m[2] {
					keys = append(keys, string(c))
				}
			} else {
				// Full name.
				keys = []string{m[2]}
			}
		}
	}

	// Collect tmp
	// to mapper if not empty.
	if len(keys) > 0 || len(values) > 0 {
		if err := o.setter(keys, values); err != nil {
			return err
		}
	}

	return nil
}

func (o *arguments) parseScript(s string) {
	if ArgumentsRegexScriptBinary.MatchString(s) || ArgumentsRegexScriptWorking.MatchString(s) {
		o.Script = s
	} else {
		o.Script = ArgumentsScript
	}
}

func (o *arguments) setter(keys, values []string) error {
	var (
		n  = len(keys) - 1
		vs = strings.Join(values, " ")
	)

	// Range key
	// and assign value on last key.
	for i, key := range keys {
		if _, ok := o.Mapper[key]; ok {
			return fmt.Errorf("option can not specify twice: %s", key)
		}

		// Set mapper.
		if i == n {
			// Assign on last.
			o.Mapper[key] = vs
		} else {
			// Assign empty on not last.
			o.Mapper[key] = ""
		}
	}

	return nil
}
