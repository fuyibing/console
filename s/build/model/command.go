// author: wsfuyibing <websearch@163.com>
// date: 2021-02-25

// Package command for build application model.
package model

import (
	"errors"
	"fmt"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/fuyibing/db"

	"github.com/fuyibing/console/v2/base"
	"github.com/fuyibing/console/v2/i"
	"github.com/fuyibing/console/v2/s/build"
)

const (
	Description = "Build model file for iris application, model struct dependent on table columns, " +
		"you need ./config/db.yaml config file settings, you can got config from " +
		"http://udsdk.turboradio.cn/ui/testing/kv/go/db/edit"
	Name = "build:model"
)

var (
	regexpFirstLetter = regexp.MustCompile(`^([a-zA-Z0-9])`)
	regexpResetSnake  = regexp.MustCompile(`[_]+([a-zA-Z0-9])`)
	regexpTypeName    = regexp.MustCompile(`^([_a-zA-Z0-9\-]+)`)
)

// Command struct.
type command struct {
	base.Command
	packages map[string]int
}

// New build model instance.
func New() i.ICommand {
	// normal.
	o := &command{packages: make(map[string]int)}
	o.Initialize()
	o.SetDescription(Description)
	o.SetName(Name)
	// model name.
	o.Add(
		base.NewOption(i.RequiredMode, i.StrValue).
			SetName("name").SetShortName("n").
			SetDescription("Model name, equal to file file."),
	)
	// table name.
	o.Add(
		base.NewOption(i.OptionalMode, i.StrValue).
			SetName("table").SetShortName("t").
			SetDescription("Table name, read columns from the table."),
	)
	// table prefix.
	o.Add(
		base.NewOption(i.OptionalMode, i.StrValue).
			SetName("prefix").
			SetDescription("Table prefix prefix."),
	)
	// application path.
	o.Add(
		base.NewOption(i.OptionalMode, i.StrValue).
			SetName("path").SetShortName("p").
			SetDefaultValue("./app").
			SetDescription("Application path."),
	)
	// list tables and columns.
	//   -l
	//   --list
	o.Add(
		base.NewOption(i.OptionalMode, i.BoolValue).
			SetName("list").SetShortName("l").
			SetDescription("List all tables and columns."),
	)
	// override if file exist.
	//   -o
	//   --override
	o.Add(
		base.NewOption(i.OptionalMode, i.BoolValue).
			SetName("override").SetShortName("o").
			SetDescription("Override if file exist"),
	)
	// prepared.
	return o
}

// Run command.
func (o *command) Run(console i.IConsole) {
	// variables.
	name := o.GetOption("name").ToString()
	exportName := o.toExportName(name)
	path := o.GetOption("path").ToString() + "/models"
	file := path + "/" + name + ".go"
	// logger.
	console.Info("Command %s: begin.", o.GetName())
	console.Info("        name: %s.", exportName)
	console.Info("        file: %s.", file)
	defer console.Info("Command %s: completed.", o.GetName())
	// file exist for not override.
	if ok, _ := o.fileExist(file); ok && !o.GetOption("override").ToBool() {
		console.PrintError(errors.New(fmt.Sprintf("Command %s: file exist", o.GetName())))
		return
	}
	table := o.GetOption("table").ToString()
	if table == "" {
		table = name
	}
	// List tables and columns.
	if o.GetOption("list").ToBool() {
		if err := o.listTables(); err != nil {
			console.PrintError(err)
			return
		}
		if err := o.listColumn(table); err != nil {
			console.PrintError(err)
			return
		}
	}
	// Prepare.
	var err error
	var dumpColumns, dumpHead, dumpTableName string
	// Dump progress.
	if dumpColumns, err = o.dumpColumns(exportName, table); err != nil {
		console.PrintError(err)
		return
	}
	dumpHead = o.dumpHead()
	dumpTableName = o.dumpTableName(exportName, table)
	// Save model.
	if err = o.write(path, file, dumpHead+"\n"+dumpColumns+"\n"+dumpTableName); err != nil {
		console.PrintError(err)
		return
	}
}

// Check model file exist.
func (o *command) fileExist(file string) (bool, error) {
	_, err := os.Stat(file)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

// Dump model fields.
// All fields read from table columns.
func (o *command) dumpColumns(exportName, table string) (string, error) {
	str := ""
	// Table info.
	tbl := &build.BeanTable{}
	if _, err := db.Slave().SQL(fmt.Sprintf("SHOW TABLE STATUS LIKE '%s'", table)).Get(tbl); err != nil {
		return "", errors.New(fmt.Sprintf("Command %s: show table status error: %v", o.GetName(), err))
	}
	if tbl.Name == "" {
		return "", errors.New(fmt.Sprintf("Command %s: show table status failed", o.GetName()))
	}
	// Table comment.
	if tbl.Comment == "" {
		str += fmt.Sprintf("// %s.\n", exportName)
	} else {
		str += fmt.Sprintf("// %s.\n", tbl.Comment)
	}
	// Model struct.
	str += fmt.Sprintf("type %s struct {\n", exportName)
	// Read columns.
	cols := make([]*build.BeanColumn, 0)
	if err := db.Slave().SQL(fmt.Sprintf("SHOW FULL COLUMNS FROM %s", table)).Find(&cols); err != nil {
		return "", errors.New(fmt.Sprintf("Command %s: show table columns error: %v", o.GetName(), err))
	}
	// Loop columns.
	for n, col := range cols {
		// not first column.
		if n > 0 {
			str += fmt.Sprintf("\n")
		}
		// comment.
		if col.Comment != "" {
			str += fmt.Sprintf("    // %s.\n", col.Comment)
		}
		// column.
		str += fmt.Sprintf("    // name: %s.\n", col.Field)
		str += fmt.Sprintf("    // type: %s.\n", col.Type)
		str += fmt.Sprintf("    %s %s %s\n", o.toExportName(col.Field), o.toExportType(col), o.toExportTag(col))
	}
	str += fmt.Sprintf("}\n")
	// Completed.
	return str, nil
}

// Dump model header.
func (o *command) dumpHead() (str string) {
	// base.
	str += fmt.Sprintf("// date: %s\n", time.Now().Format("2006-01-02 15:04"))
	str += fmt.Sprintf("// author: %s\n", o.GetName())
	str += fmt.Sprintf("// command: %s %s\n", i.Script, o.GetName())
	// package.
	str += fmt.Sprintf("\n")
	str += fmt.Sprintf("package models\n")
	// import.
	str += fmt.Sprintf("\n")
	if len(o.packages) > 0 {
		str += fmt.Sprintf("import (\n")
		for pkg, _ := range o.packages {
			str += fmt.Sprintf("    \"%s\"\n", pkg)
		}
		str += fmt.Sprintf(")\n")
	}
	return
}

// Dump table name method.
func (o *command) dumpTableName(exportName, table string) string {
	str := fmt.Sprintf("// Return table name.\n")
	str += fmt.Sprintf("func (*%s) TableName() string {\n", exportName)
	str += fmt.Sprintf("    return \"%s\"\n", table)
	str += fmt.Sprintf("}\n")
	return str
}

// List all columns.
func (o *command) listColumn(name string) error {
	cols := make([]*build.BeanColumn, 0)
	if err := db.Slave().SQL(fmt.Sprintf("SHOW FULL COLUMNS FROM %s", name)).Find(&cols); err != nil {
		return err
	}
	// prepare table.
	table := base.NewTable()
	// append head.
	table.Head().Add(
		base.NewCell("ID"),
		base.NewCell("Name").SetColor(base.ColorBlue),
		base.NewCell("Type").SetColor(base.ColorBlue),
		base.NewCell("Field Name"),
		base.NewCell("Field Comment"),
	)
	// append body.
	for n, col := range cols {
		// append row.
		table.Body().Add(base.NewRow().Add(
			base.NewCell(fmt.Sprintf("%d", n+1)).SetAlign(base.AlignRight),
			base.NewCell(o.toExportName(col.Field)).SetColor(base.ColorBlue),
			base.NewCell(o.toExportType(col)).SetColor(base.ColorBlue),
			base.NewCell(col.Field),
			base.NewCell(col.Comment),
		))
	}
	// print table.
	table.SetPrefix("        ").Print()
	return nil
}

// List all tables.
func (o *command) listTables() error {
	// generate query.
	var query = "SHOW TABLE STATUS"
	if prefix := o.GetOption("prefix").ToString(); prefix != "" {
		query += fmt.Sprintf(" LIKE '%s%%'", prefix)
	}
	// show columns.
	tbl := make([]*build.BeanTable, 0)
	if err := db.Slave().SQL(query).Find(&tbl); err != nil {
		return errors.New(fmt.Sprintf(
			"Command %s: show table status error: %v",
			o.GetName(),
			err,
		))
	}
	// prepare table.
	table := base.NewTable()
	// append head.
	table.Head().Add(
		base.NewCell("ID"),
		base.NewCell("Name").SetColor(base.ColorBlue),
		base.NewCell("Created Time"),
		base.NewCell("Table comment"),
	)
	// append body.
	for n, t := range tbl {
		// append row.
		table.Body().Add(base.NewRow().Add(
			base.NewCell(fmt.Sprintf("%d", n+1)).SetAlign(base.AlignRight),
			base.NewCell(t.Name).SetColor(base.ColorBlue),
			base.NewCell(t.Created),
			base.NewCell(t.Comment),
		))
	}
	// print table.
	table.SetPrefix("        ").Print()
	return nil
}

// Convert to model field name.
func (o *command) toExportName(name string) string {
	return regexpFirstLetter.ReplaceAllStringFunc(
		regexpResetSnake.ReplaceAllStringFunc(name, func(s string) string {
			m := regexpResetSnake.FindStringSubmatch(s)
			return strings.ToUpper(m[1])
		}), func(s string) string {
			m := regexpFirstLetter.FindStringSubmatch(s)
			return strings.ToUpper(m[1])
		},
	)
}

// Convert to xorm tag.
func (o *command) toExportTag(col *build.BeanColumn) string {
	str := ""
	if col.Key == "PRI" {
		str += "pk autoincr "
	}
	str += col.Field
	return fmt.Sprintf("`xorm:\"%s\"`", str)
}

// Convert to model field type.
func (o *command) toExportType(col *build.BeanColumn) string {
	if m := regexpTypeName.FindStringSubmatch(col.Type); len(m) == 2 {
		if t, ok := build.TypeMapping[m[1]]; ok {
			if s := strings.Split(t, ":"); len(s) == 2 {
				if s[1] == "" {
					s[1] = s[0]
				}
				o.packages[s[1]] = 1
				return s[0]
			}
			return t
		}
		return m[1]
	}
	return "interface{}"
}

// Write content to specified file.
func (o *command) write(path, file, content string) error {
	// make directory.
	if err := os.MkdirAll(path, os.ModePerm); err != nil {
		return err
	}
	// prepare.
	var err error
	var fp *os.File
	// open file.
	if fp, err = os.OpenFile(file, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, os.ModePerm); err != nil {
		return errors.New(fmt.Sprintf("Command %s: create file error: %v", o.GetName(), err))
	}
	// close when end.
	defer func() {
		_ = fp.Close()
	}()
	// write content.
	if _, err = fp.WriteString(content); err != nil {
		return errors.New(fmt.Sprintf("Command %s: create file error: %v", o.GetName(), err))
	}
	// completed.
	return nil
}
