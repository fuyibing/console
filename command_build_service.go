// author: wsfuyibing <websearch@163.com>
// date: 2021-03-23

package console

import (
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/fuyibing/db"
)

// Build service command struct.
type buildServiceCommand struct {
	command         *Command
	imports         map[string]string
	pathName        string
	fileName        string
	serviceName     string
	structName      string
	modelStructName string
	moduleName      string
	tableName       string
}

// Handle after.
func (o *buildServiceCommand) after(cs *Console) error { return nil }

// Handle before.
func (o *buildServiceCommand) before(cs *Console) error {
	// detect module path.
	if err := o.renderModule(cs); err != nil {
		return err
	}
	// append path.
	o.imports = make(map[string]string)
	o.imports["models"] = fmt.Sprintf("%s/app/models", o.moduleName)
	o.imports["db"] = "github.com/fuyibing/db"
	o.imports["xorm"] = "xorm.io/xorm"
	// normal name.
	o.serviceName = o.command.GetOption("name").String()
	o.modelStructName = o.toExportName(o.serviceName)
	o.structName = o.modelStructName + "Service"
	o.pathName = fmt.Sprintf("%s", o.command.GetOption("path").String())
	o.fileName = fmt.Sprintf("%s_service.go", o.command.GetOption("name").String())
	// table name.
	if o.tableName = o.command.GetOption("table").String(); o.tableName == "" {
		o.tableName = o.serviceName
	}
	// Allow override service file.
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
	return fmt.Errorf("service file exists: %s/%s", o.pathName, o.fileName)
}

// Handle build command.
func (o *buildServiceCommand) handler(cs *Console) error {

	columns, err := o.listColumns(cs)
	if err != nil {
		return err
	}

	str := ""

	str += o.renderCopyright(cs)
	str += o.renderImports(cs)
	str += o.renderStruct(cs)
	str += o.renderAdd(cs, columns)
	str += o.renderGetByPk(cs, columns)

	// println(str)
	return o.write(cs, str)
}

// List columns.
func (o *buildServiceCommand) listColumns(cs *Console) ([]*Column, error) {
	columns := make([]*Column, 0)
	if err := db.Slave().SQL(fmt.Sprintf("SHOW FULL COLUMNS FROM `%s`", o.tableName)).Find(&columns); err != nil {
		return nil, err
	}
	return columns, nil
}

// Primary key.
func (o *buildServiceCommand) renderAdd(cs *Console, columns []*Column) string {
	var column *Column
	str := fmt.Sprintf("// Add model by request.\n")
	str += fmt.Sprintf("func (o *%s) Add(req *models.%s) (*models.%s, error) {\n", o.structName, o.modelStructName, o.modelStructName)
	str += fmt.Sprintf("    // todo: created by console.\n")
	str += fmt.Sprintf("    // define bean.\n")
	str += fmt.Sprintf("    bean := &models.%s{}\n", o.modelStructName)
	str += fmt.Sprintf("    // assign values.\n")
	for _, c := range columns {
		if c.Key == "PRI" {
			column = c
			continue
		}
		str += fmt.Sprintf("    bean.%s = req.%s\n", o.toExportName(c.Name), o.toExportName(c.Name))
	}
	str += fmt.Sprintf("    // execute db query.\n")
	str += fmt.Sprintf("    if _, err := o.Master().Insert(bean); err != nil {\n")
	str += fmt.Sprintf("        return nil, err\n")
	str += fmt.Sprintf("    }\n")
	if column != nil {
		str += fmt.Sprintf("    // check primary key.\n")
		str += fmt.Sprintf("    if bean.%s > 0 {\n", o.toExportName(column.Name))
		str += fmt.Sprintf("        return bean, nil\n")
		str += fmt.Sprintf("    }\n")
		str += fmt.Sprintf("    // return not found.\n")
		str += fmt.Sprintf("    return nil, nil\n")
	} else {
		str += fmt.Sprintf("    // return result.\n")
		str += fmt.Sprintf("    return bean, nil\n")
	}
	str += fmt.Sprintf("}\n\n")
	return str
}

// Render copyright.
func (o *buildServiceCommand) renderCopyright(cs *Console) string {
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
	str += fmt.Sprintf("package services\n\n")
	return str
}

// Header imports.
func (o *buildServiceCommand) renderImports(cs *Console) string {
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

// Render module.
func (o *buildServiceCommand) renderModule(cs *Console) error {
	reg := regexp.MustCompile(`module\s+([^\n]+)`)
	for _, file := range []string{"go.mod", "../go.mod"} {
		if body, err := ioutil.ReadFile(file); err == nil {
			if m := reg.FindStringSubmatch(string(body)); len(m) == 2 {
				m[1] = strings.TrimSpace(m[1])
				if m[1] != "" {
					o.moduleName = m[1]
					return nil
				}
			}
		}
	}
	return fmt.Errorf("can not find go.md file")
}

// Render struct.
func (o *buildServiceCommand) renderStruct(cs *Console) string {
	str := ""
	// definition.
	str += fmt.Sprintf("// service struct.\n")
	str += fmt.Sprintf("type %s struct{\n", o.structName)
	str += fmt.Sprintf("    db.Service\n")
	str += fmt.Sprintf("}\n\n")
	// new services.
	str += fmt.Sprintf("// create instance.\n")
	str += fmt.Sprintf("func New%s(sess ...*xorm.Session) *%s {\n", o.structName, o.structName)
	str += fmt.Sprintf("    o := &%s{}\n", o.structName)
	str += fmt.Sprintf("    o.Use(sess...)\n")
	str += fmt.Sprintf("    return o\n")
	str += fmt.Sprintf("}\n\n")
	return str
}

// Primary key.
func (o *buildServiceCommand) renderGetByPk(cs *Console, columns []*Column) string {
	var str = ""
	var column *Column
	for _, c := range columns {
		if c.Key == "PRI" {
			column = c
			break
		}
	}
	if column == nil {
		return str
	}
	str += fmt.Sprintf("// Get model by primary key.\n")
	str += fmt.Sprintf("func (o *%s) GetByPk(id int64) (*models.%s, error) {\n", o.structName, o.modelStructName)
	str += fmt.Sprintf("    return o.GetBy%s(id)\n", o.toExportName(column.Name))
	str += fmt.Sprintf("}\n\n")
	str += fmt.Sprintf("// Get by column name.\n")
	str += fmt.Sprintf("func (o *%s) GetBy%s(id int64) (*models.%s, error) {\n", o.structName, o.toExportName(column.Name), o.modelStructName)
	str += fmt.Sprintf("    // todo: created by console.\n")
	str += fmt.Sprintf("    // define model bean.\n")
	str += fmt.Sprintf("    bean := &models.%s{}\n", o.modelStructName)
	str += fmt.Sprintf("    // get one and assign values.\n")
	str += fmt.Sprintf("    if _, err := o.Slave().Where(\"%s = ?\", id).Get(bean); err != nil {\n", column.Name)
	str += fmt.Sprintf("        return nil, err\n")
	str += fmt.Sprintf("    }\n")
	str += fmt.Sprintf("    // check assigned value state.\n")
	str += fmt.Sprintf("    if bean.%s > 0 {\n", o.toExportName(column.Name))
	str += fmt.Sprintf("        return bean, nil\n")
	str += fmt.Sprintf("    }\n")
	str += fmt.Sprintf("    // return nil if not found.\n")
	str += fmt.Sprintf("    return nil, nil\n")
	str += fmt.Sprintf("}\n\n")
	return str
}

// To export name.
// Large camel format.
func (o *buildServiceCommand) toExportName(name string) string {
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

// Write to file.
func (o *buildServiceCommand) write(cs *Console, text string) error {
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
func newBuildServiceCommand() *Command {
	// base.
	c := NewCommand("bs")
	c.SetDescription("Build service for application")
	// options.
	c.Add(
		NewOption("name").SetTag('n').
			SetDescription("Service name, equal to model name and suffix is <_service>"),
		NewOption("override").SetTag('o').
			SetMode(OptionalMode).SetValue(NullValue).
			SetDescription("Override if service exist, default is <false>"),
		NewOption("path").SetTag('p').
			SetMode(OptionalMode).SetValue(StringValue).
			SetDefaultValue("app/services").
			SetDescription("Created service file save to, default is <app/services>"),
		NewOption("table").SetTag('t').
			SetMode(OptionalMode).
			SetDescription("Specify table name, default is name option value if not specified"),
	)
	// register handler.
	o := &buildServiceCommand{command: c}
	c.SetHandlerBefore(o.before).SetHandler(o.handler).SetHandlerAfter(o.after)
	return c
}
