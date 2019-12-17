package prepare

import (
	"bytes"
	"errors"
	"fmt"
	"path"
	"sort"
	"strings"
	"text/template"

	"github.com/azhai/gozzo-db/construct"
	"github.com/azhai/gozzo-db/rewrite"
	"github.com/azhai/gozzo-db/schema"
	"github.com/azhai/gozzo-db/utils"
	"github.com/jinzhu/gorm"
)

var TemplateNotFound = errors.New("template not found")

func GetTemplate(name string, funcMap template.FuncMap) *template.Template {
	if fsize, _ := utils.FileSize(utils.GetAbsFile(name)); fsize > 0 {
		// New()名称要与ParseFiles()文件名一致 Funcs()必须在Parse之前
		t := template.New(name).Funcs(funcMap)
		return template.Must(t.ParseFiles(name))
	}
	if content, ok := templates[path.Base(name)]; ok {
		// Funcs()必须在Parse之前
		t := template.New(name).Funcs(funcMap)
		return template.Must(t.Parse(content))
	}
	return nil
}

func CreateModels(conf *Config, db *gorm.DB) (names []string, err error) {
	var funcMap = template.FuncMap{
		"GetRule":       GetRule,
		"GenNameType":   GenNameType,
		"GenTagComment": GenTagComment,
	}
	t := GetTemplate("gen_model.tmpl", funcMap)
	if t == nil {
		return nil, TemplateNotFound
	}

	var buf bytes.Buffer
	app := conf.Application
	tablePrefix := conf.GetTablePrefix(conf.ConnName)
	s := schema.NewSchema(db.DB())
	for i, table := range s.GetTableNames("") {
		// 收集数据
		cols := s.GetColumnInfos(table, "")
		tbInfo := s.GetTableInfo(table, "")
		if tablePrefix != "" && strings.HasPrefix(table, tablePrefix) {
			table = table[len(tablePrefix):]
		}
		if app.PluralTable {
			table = utils.ToSingular(table)
		}
		name := utils.ToCamel(table)
		names = append(names, name)

		// 渲染模板
		data := map[string]interface{}{
			"Name": name, "Table": tbInfo, "Columns": cols,
			"Rules": conf.GetRules(tbInfo.TableName),
		}
		buf.Reset()
		if err = t.Execute(&buf, data); err != nil {
			continue
		}

		// 写入文件
		fname := fmt.Sprintf("%s/%s.go", app.OutputDir, table)
		if i == 0 {
			utils.MkdirForFile(fname)
		}
		cs := rewrite.NewCodeSource()
		err = cs.SetPackage("models")
		cs.AddImport("github.com/azhai/gozzo-db/construct", "base")
		cs.AddImport("github.com/jinzhu/gorm", "")
		cs.AddImport("time", "")
		err = cs.AddCode(buf.Bytes())
		cs.DelImport("github.com/jinzhu/gorm", "")
		cs.DelImport("time", "")
		err = cs.WriteTo(fname)
	}
	return
}

func GenInitFile(conf *Config, names []string) error {
	var buf bytes.Buffer
	app := conf.Application
	sort.Strings(names)
	for i, name := range names {
		if i == 0 {
			buf.WriteString(fmt.Sprintf("&%s{}", name))
		} else if i%3 == 0 {
			buf.WriteString(fmt.Sprintf(",\n\t\t&%s{}", name))
		} else {
			buf.WriteString(fmt.Sprintf(", &%s{}", name))
		}
	}
	models := buf.String()

	// 渲染模板
	driverName := conf.GetDriverName(conf.ConnName)
	tablePrefix := conf.GetTablePrefix(conf.ConnName)
	data := map[string]interface{}{
		"FileName": conf.FileName, "ConnName": conf.ConnName,
		"Prefix": tablePrefix, "Plural": app.PluralTable,
		"Models": models,
	}
	t := GetTemplate("gen_init.tmpl", nil)
	if t == nil {
		return TemplateNotFound
	}
	buf.Reset()
	if err := t.Execute(&buf, data); err != nil {
		return err
	}

	// 写入文件
	fname := fmt.Sprintf("%s/init.go", app.OutputDir)
	cs := rewrite.NewCodeSource()
	err := cs.SetPackage("models")
	cs.AddImport("log", "")
	cs.AddImport("os", "")
	cs.AddImport("github.com/azhai/gozzo-db/construct", "base")
	cs.AddImport("github.com/azhai/gozzo-db/prepare", "")
	cs.AddImport("github.com/jinzhu/gorm", "")
	cs.AddImport("github.com/jinzhu/gorm/dialects/"+driverName, "_")
	err = cs.AddCode(buf.Bytes())
	err = cs.WriteTo(fname)
	return err
}

func GenNameType(col *schema.ColumnInfo, rule RuleConfig) string {
	if rule.Name == "" {
		rule.Name = utils.ToCamel(col.FieldName)
	}
	if rule.Type == "" {
		rule.Type = construct.GuessTypeName(col)
	}
	return rule.Name + " " + rule.Type
}

func GenTagComment(col *schema.ColumnInfo, rule RuleConfig) string {
	var blank, comment string
	if rule.Json == "" {
		rule.Json = col.FieldName
	}
	if rule.Tags == "" {
		rule.Tags = construct.GuessStructTags(col)
	}
	if rule.Tags != "" {
		blank = " "
	}
	if col.Comment != "" {
		comment = " // " + col.Comment
	} else if rule.Comment != "" {
		comment = " // " + rule.Comment
	}
	tpl := "`json:\"%s\"%s%s`%s"
	return fmt.Sprintf(tpl, rule.Json, blank, rule.Tags, comment)
}
