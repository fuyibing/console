// author: wsfuyibing <websearch@163.com>
// date: 2021-03-13

package boot

import (
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"strconv"
	"strings"
)

const (
	PidFile          = "./pid"
	StartName        = "start"
	StartDescription = "start application"
	StopName         = "stop"
	StopDescription  = "stop application"
)

// Read pid file.
func read() (value int, found bool, err error) {
	// read content.
	var body []byte
	if body, err = ioutil.ReadFile(PidFile); err != nil {
		if regexp.MustCompile(`no\s+such\s+file\s+or\s+directory`).MatchString(err.Error()) {
			err = nil
		}
		return
	}
	// empty value.
	found = true
	text := strings.TrimSpace(string(body))
	if text == "" {
		return
	}
	// invalid value.
	n := int64(0)
	if n, err = strconv.ParseInt(text, 0, 32); err != nil {
		err = nil
		return
	}
	value = int(n)
	return
}

// Write pid file.
func write(pid int) error {
	fp, err := os.OpenFile(PidFile, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, os.ModePerm)
	if err != nil {
		return err
	}
	defer func() {
		_ = fp.Close()
	}()
	if _, err = fp.WriteString(fmt.Sprintf("%d", pid)); err != nil {
		return err
	}
	return nil
}
