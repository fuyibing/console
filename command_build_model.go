// author: wsfuyibing <websearch@163.com>
// date: 2021-03-23

package console

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/fuyibing/db"
)

// Build model command struct.
type buildModelCommand struct {
	command    *Command
	pathName   string
	fileName   string
	modelName  string
	tableName  string
	structName string
	jsonMode   string
	imports    map[string]string
}

// Handle after.
func (o *buildModelCommand) after(cs *Console) error { return nil }

// Handle before.
func (o *buildModelCommand) before(cs *Console) error {
	o.imports = make(map[string]string)
	// normal name.
	o.modelName = o.command.GetOption("name").String()
	o.structName = o.toExportName(o.modelName)
	o.jsonMode = strings.ToLower(o.command.GetOption("json").String())
	o.pathName = fmt.Sprintf("%s", o.command.GetOption("path").String())
	o.fileName = fmt.Sprintf("%s.go", o.command.GetOption("name").String())
	// table name.
	if o.tableName = o.command.GetOption("table").String(); o.tableName == "" {
		o.tableName = o.modelName
	}
	// Allow override model file.
	if o.command.GetOption("override").Bool() {
		return nil
	}
	// Return nil if file not exist.
	if _, err := os.Stat(o.pathName + "/" + o.fileName); err != nil {
		if os.IsNotExist(err) {
			return nil
		}
	}
	// Return error if file exist.
	return fmt.Errorf("model file exists: %s/%s", o.pathName, o.fileName)
}

// Handle build command.
func (o *buildModelCommand) handler(cs *Console) error {
	// Read columns
	columns, err := o.listColumns(cs)
	if err != nil {
		return err
	}
	// Render struct.
	body, err2 := o.renderStruct(cs, columns)
	if err2 != nil {
		return err2
	}
	// Result.
	text := o.renderCopyright(cs)
	text += o.renderImports(cs)
	text += body
	return o.write(cs, text)
}

// List columns.
func (o *buildModelCommand) listColumns(cs *Console) ([]*Column, error) {
	columns := make([]*Column, 0)
	if err := db.Slave().SQL(fmt.Sprintf("SHOW FULL COLUMNS FROM `%s`", o.tableName)).Find(&columns); err != nil {
		return nil, err
	}
	return columns, nil
}

func (o *buildModelCommand) listTable(cs *Console) (*Table, error) {
	table := &Table{}
	if _, err := db.Slave().SQL(fmt.Sprintf("SHOW TABLE STATUS LIKE '%s'", o.tableName)).Get(table); err != nil {
		return nil, err
	}
	return table, nil
}

func (o *buildModelCommand) renderCopyright(cs *Console) string {
	// copyright.
	str := fmt.Sprintf("// author: console\n")
	str += fmt.Sprintf("// date: %s\n", time.Now().Format("2006-01-02 15:04:05"))
	// append command line.
	str += fmt.Sprintf("// command: %s", o.command.name)
	if len(os.Args) > 0 {
		for n, arg := range os.Args {
			if n > 1 {
				str += fmt.Sprintf(" %s", arg)
			}
		}
	}
	str += fmt.Sprintf("\n\n")
	// header.
	str += fmt.Sprintf("package models\n\n")
	return str
}

func (o *buildModelCommand) renderImports(cs *Console) string {
	str := ""
	if len(o.imports) > 0 {
		str += fmt.Sprintf("import (\n")
		for _, v := range o.imports {
			str += fmt.Sprintf("    \"%s\"\n", v)
		}
		str += fmt.Sprintf(")\n")
		str += "\n"
	}
	return str
}

func (o *buildModelCommand) renderStruct(cs *Console, columns []*Column) (string, error) {
	// Table info.
	str := ""
	table, err := o.listTable(cs)
	if err != nil {
		return str, err
	}
	// Comment.
	if table.Comment == "" {
		str += fmt.Sprintf("// Model struct.\n")
	} else {
		str += fmt.Sprintf("// %s\n", table.Comment)
	}
	// standard.
	str += fmt.Sprintf("// name: %s\n", o.modelName)
	str += fmt.Sprintf("// table: %s\n", table.Name)
	str += fmt.Sprintf("// create: %s\n", table.Created.Format("2006-01-02 15:04"))
	str += fmt.Sprintf("type %s struct {\n", o.structName)
	// Columns.
	for n, column := range columns {
		// Separator.
		if n > 0 {
			str += fmt.Sprintf("\n")
		}
		// Column comment.
		if column.Comment != "" {
			str += fmt.Sprintf("    // %s.\n", column.Comment)
		}
		// Column type.
		if column.Type != "" {
			str += fmt.Sprintf("    // type: %s\n", column.Type)
		}
		// Column definition.
		str += fmt.Sprintf(
			"    %s %s `xorm:\"%s\" json:\"%s\"`\n",
			o.toExportName(column.Name),
			o.toExportType(column.Type),
			o.toExportOrm(column.Name, column.Key),
			o.toTargetName(column.Name),
		)
	}
	str += fmt.Sprintf("}\n")
	// table name.
	str += fmt.Sprintf("\n")
	str += fmt.Sprintf("// Return table name.\n")
	str += fmt.Sprintf("func (o *%s) TableName() string {\n", o.structName)
	str += fmt.Sprintf("    return \"%s\"\n", o.tableName)
	str += fmt.Sprintf("}\n")
	// result.
	return str, nil
}

// To export name.
// Large camel format.
func (o *buildModelCommand) toExportName(name string) string {
	// clear prefix with underline.
	name = RegexpUnderlinePrefix.ReplaceAllString(name, "")
	// change first letter as upper.
	name = RegexpFirstLetter.ReplaceAllStringFunc(name, func(s string) string {
		return strings.ToUpper(s)
	})
	// match underline of middle, remove underline and change upper.
	name = RegexpUnderlineMiddle.ReplaceAllStringFunc(name, func(s string) string {
		if m := RegexpUnderlineMiddle.FindStringSubmatch(s); len(m) == 2 {
			return strings.ToUpper(m[1])
		}
		return s
	})
	// return result string.
	return name
}

// To export type.
func (o *buildModelCommand) toExportType(name string) (s string) {
	s = "interface{}"
	if m := RegexpColumnType.FindStringSubmatch(name); len(m) == 2 {
		switch strings.ToLower(m[1]) {
		// timeline.
		case "timestamp", "datetime":
			s = "time.Time"
			o.imports[s] = "time"
		case "bigint":
			s = "int64"
		case "tinyint", "smallint", "mediumint", "int":
			s = "int"
		case "double", "float", "decimal":
			s = "float64"
		case "char", "varchar", "text", "enum", "date", "time":
			s = "string"
		}
	}
	return
}

// Export orm tag.
func (o *buildModelCommand) toExportOrm(name, key string) string {
	str := name
	if key == "PRI" {
		str += " pk autoincr"
	}
	return str
}

// Export json tag.
// Convert export name to target name.
func (o *buildModelCommand) toTargetName(name string) string {
	// Same as export name.
	// convert any name as Large camel.
	if o.jsonMode == "same" {
		return o.toExportName(name)
	}
	// Camel model, with lower.
	if o.jsonMode == "camel" {
		// 1. remove underline prefix.
		name = RegexpUnderlinePrefix.ReplaceAllString(name, "")
		// 2. change first letter as lower.
		name = RegexpFirstLetter.ReplaceAllStringFunc(name, func(s string) string {
			return strings.ToLower(s)
		})
		// 3. change underline with lower letter as upper letter.
		name = RegexpUnderlineMiddle.ReplaceAllStringFunc(name, func(s string) string {
			m := RegexpUnderlineMiddle.FindStringSubmatch(s)
			return strings.ToUpper(m[1])
		})
		// 4. completed.
		return name
	}
	// Snake mode.
	if o.jsonMode == "snake" {
		// 1. remove underline prefix.
		name = RegexpUnderlinePrefix.ReplaceAllString(name, "")
		// 2. change first letter as lower.
		name = RegexpFirstLetter.ReplaceAllStringFunc(name, func(s string) string {
			return strings.ToLower(s)
		})
		// 3. change upper letter as underline with lower letter.
		name = RegexpUpperLetter.ReplaceAllStringFunc(name, func(s string) string {
			return "_" + strings.ToLower(s)
		})
		// 4. completed.
		return name
	}
	// auto mode.
	return name
}

// Write to file.
func (o *buildModelCommand) write(cs *Console, text string) error {
	// Make directory.
	if err := os.MkdirAll(o.pathName, os.ModePerm); err != nil {
		return err
	}
	// Open file, create if not exist.
	fp, err2 := os.OpenFile(o.pathName+"/"+o.fileName, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, os.ModePerm)
	if err2 != nil {
		return err2
	}
	// Close file handler when ended.
	defer func() {
		_ = fp.Close()
	}()
	// Write contents
	_, err3 := fp.WriteString(text)
	if err3 != nil {
		return err3
	}
	// Succeed.
	return nil
}

// New build model command.
func newBuildModelCommand() *Command {
	// base.
	c := NewCommand("bm")
	c.SetDescription("Build model for application")
	// options.
	c.Add(
		NewOption("json").
			SetMode(OptionalMode).
			SetAccepts("auto", "same", "snake", "camel").SetDefaultValue("auto").
			SetDescription("Export as json string format, accept <auto>, <same>, <snake>, <camel>, default is <auto>"),
		NewOption("name").SetTag('n').
			SetDescription("Model name"),
		NewOption("override").SetTag('o').
			SetMode(OptionalMode).SetValue(NullValue).
			SetDescription("Override if model exist, default is <false>"),
		NewOption("path").SetTag('p').
			SetMode(OptionalMode).SetValue(StringValue).
			SetDefaultValue("app/models").
			SetDescription("Created model file save to, default is <app/models>"),
		NewOption("table").SetTag('t').
			SetMode(OptionalMode).
			SetDescription("Specify table name, default is name option value if not specified"),
	)
	// register handler.
	o := &buildModelCommand{command: c}
	c.SetHandlerBefore(o.before).SetHandler(o.handler).SetHandlerAfter(o.after)
	return c
}
