package prepare

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"path"
	"path/filepath"
	"sort"
	"strings"
	"text/template"

	base "github.com/azhai/gozzo-db/construct"
	"github.com/azhai/gozzo-db/rewrite"
	"github.com/azhai/gozzo-db/schema"
	"github.com/azhai/gozzo-utils/filesystem"
	"github.com/jinzhu/gorm"
)

const (
	MODE_SAFE uint = iota
	MODE_INIT
	MODE_TWO_FILES
	MODE_THREE_FILES
	MODE_QUERY_DISTRI
	MODE_TABLE_DISTRI
)

var (
	TemplateNotFound   = errors.New("template not found")
	TemplateModeIntros = map[uint]string{
		MODE_SAFE:         "安全模式，与 5 类似，但会在每个文件名前面加一个下划线",
		MODE_INIT:         "只生成 init.go 文件",
		MODE_TWO_FILES:    "除了 init.go 文件， table 和 query 都放入 tables.go 中",
		MODE_THREE_FILES:  "除了 init.go 文件， table 都放入 tables.go 中， query 都放入 queries.go 中",
		MODE_QUERY_DISTRI: "除了 init.go 文件， table 都放入 tables.go 中， query 分开放入对应模型名文件中",
		MODE_TABLE_DISTRI: "除了 init.go 文件， table 和 query 一起放入对应模型名文件中",
	}
)

func GetTemplate(funcMap template.FuncMap, name, fname string, others ...string) *template.Template {
	var buf bytes.Buffer
	// Funcs()必须在Parse之前
	t := template.New(name).Funcs(funcMap)
	for i := -1; i < len(others); i++ {
		if i >= 0 {
			fname = others[i]
		}
		fpath := filesystem.GetAbsFile(fname)
		if fsize, _ := filesystem.FileSize(fpath); fsize > 0 {
			if ftext, err := ioutil.ReadFile(fpath); err == nil {
				buf.WriteString("\n\n")
				buf.Write(ftext)
			}
		} else if text, ok := templates[path.Base(fname)]; ok {
			buf.WriteString(text)
		}
	}
	return template.Must(t.Parse(buf.String()))
}

func CreateModels(conf *Config, db *gorm.DB, mode uint) (names []string, err error) {
	// 模板
	var funcMap = template.FuncMap{
		"GetRule":       GetRule,
		"GenTagComment": GenTagComment,
		"GenNameType": func(col *schema.ColumnInfo, rule RuleConfig) string {
			return GenNameType(col, rule, conf.NullPointers)
		},
	}
	var tpl, tpl2 *template.Template
	if mode != MODE_INIT {
		if mode == MODE_THREE_FILES || mode == MODE_QUERY_DISTRI {
			tpl = GetTemplate(funcMap, "gen_table", "gen_table.tmpl")
			tpl2 = GetTemplate(funcMap, "gen_query", "gen_query.tmpl")
			if tpl2 == nil {
				return nil, TemplateNotFound
			}
		} else {
			tpl = GetTemplate(funcMap, "gen_model", "gen_table.tmpl", "gen_query.tmpl")
		}
		if tpl == nil {
			return nil, TemplateNotFound
		}
	}
	// 参数
	filePre := ""
	if mode == MODE_SAFE {
		filePre = "_"
	}
	outDir := conf.Application.OutputDir
	isPlural := conf.Application.PluralTable
	tablePrefix := conf.GetTablePrefix(conf.ConnName)

	tableNames := make(map[string]string)
	s := schema.NewSchema(db.DB())
	for _, tableName := range s.GetTableNames("") {
		table := tableName
		if tablePrefix != "" && strings.HasPrefix(table, tablePrefix) {
			table = table[len(tablePrefix):]
		}
		if isPlural {
			table = ToSingular(table)
		}
		name := ToCamel(table)
		names = append(names, name)
		tableNames[name] = tableName
	}
	sort.Strings(names)
	if mode == MODE_INIT { // 此模式只生成init.go文件
		return
	}

	var buf, buf2 bytes.Buffer
	for _, name := range names {
		tableName := tableNames[name]
		// 收集数据，渲染模板
		tbInfo := s.GetTableInfo(tableName, "")
		data := map[string]interface{}{
			"Name": name, "Table": tbInfo,
			"Columns": s.GetColumnInfos(tableName, ""),
			"Rules":   conf.GetRules(tbInfo.TableName),
		}
		if mode == MODE_SAFE || mode == MODE_TABLE_DISTRI {
			buf.Reset()
		}
		if err = tpl.Execute(&buf, data); err != nil {
			continue
		}
		if mode == MODE_THREE_FILES || mode == MODE_QUERY_DISTRI {
			if mode == MODE_QUERY_DISTRI {
				buf2.Reset()
			}
			if err = tpl2.Execute(&buf2, data); err != nil {
				continue
			}
		}
		// 写入文件
		fname := fmt.Sprintf("%s/%s%s.go", outDir, filePre, strings.ToLower(name))
		if mode == MODE_SAFE || mode == MODE_TABLE_DISTRI {
			err = WriteModelFile(buf, fname)
		} else {
			if mode == MODE_QUERY_DISTRI {
				err = WriteModelFile(buf2, fname)
			}
			if mode == MODE_THREE_FILES {
				fname = fmt.Sprintf("%s/queries.go", outDir)
				err = WriteModelFile(buf2, fname)
			}
			fname = fmt.Sprintf("%s/tables.go", outDir)
			err = WriteModelFile(buf, fname)
		}
	}
	return
}

func GenInitFile(conf *Config, names []string, mode uint) (err error) {
	// 模板
	tpl := GetTemplate(nil, "gen_init", "gen_init.tmpl")
	if tpl == nil {
		return TemplateNotFound
	}
	// 参数
	filePre := ""
	if mode == MODE_SAFE {
		filePre = "_"
	}
	outDir := conf.Application.OutputDir
	isPlural := conf.Application.PluralTable
	tablePrefix := conf.GetTablePrefix(conf.ConnName)
	driverName := conf.GetDriverName(conf.ConnName)
	fname := fmt.Sprintf("%s/%sinit.go", outDir, filePre)

	// 收集数据，渲染模板
	var buf bytes.Buffer
	sort.Strings(names)
	for i, name := range names {
		if i > 0 && i%3 == 0 {
			buf.WriteString("\n\t\t")
		}
		buf.WriteString(fmt.Sprintf("&%s{}, ", name))
	}
	models := buf.String()
	data := map[string]interface{}{
		"FileName": conf.FileName, "ConnName": conf.ConnName,
		"Prefix": tablePrefix, "Plural": isPlural, "Models": models,
	}
	buf.Reset()
	if err := tpl.Execute(&buf, data); err != nil {
		return err
	}
	// 写入文件
	return WriteInitFile(buf, fname, driverName)
}

func WriteModelFile(buf bytes.Buffer, fname string) (err error) {
	// 写入文件
	cs := rewrite.NewCodeSource()
	ns := filepath.Base(filepath.Dir(fname))
	if err = cs.SetPackage(ns); err != nil {
		return
	}
	// 添加可能引用的包，后面再尝试删除不一定会用的包
	cs.AddImport("database/sql", "")
	cs.AddImport("time", "")
	cs.AddImport("github.com/azhai/gozzo-db/construct", "base")
	cs.AddImport("github.com/jinzhu/gorm", "")
	if err = cs.AddCode(buf.Bytes()); err != nil {
		return
	}
	// 尝试删除，已用到的包不会被删除
	cs.DelImport("database/sql", "")
	cs.DelImport("time", "")
	cs.DelImport("github.com/jinzhu/gorm", "")
	err = cs.WriteTo(fname)
	return
}

func WriteInitFile(buf bytes.Buffer, fname, driverName string) (err error) {
	// 写入文件
	cs := rewrite.NewCodeSource()
	ns := filepath.Base(filepath.Dir(fname))
	if err = cs.SetPackage(ns); err != nil {
		return
	}
	// 以下包在默认模板都会引用
	cs.AddImport("log", "")
	cs.AddImport("os", "")
	cs.AddImport("github.com/azhai/gozzo-db/construct", "base")
	cs.AddImport("github.com/azhai/gozzo-db/cache", "")
	cs.AddImport("github.com/azhai/gozzo-db/export", "")
	cs.AddImport("github.com/azhai/gozzo-db/prepare", "")
	cs.AddImport("github.com/jinzhu/gorm", "")
	cs.AddImport("github.com/jinzhu/gorm/dialects/"+driverName, "_")
	if err = cs.AddCode(buf.Bytes()); err != nil {
		return
	}
	// 尝试删除，已用到的包不会被删除
	cs.DelImport("github.com/jinzhu/gorm", "")
	cs.DelImport("github.com/azhai/gozzo-db/cache", "")
	cs.DelImport("github.com/azhai/gozzo-db/export", "")
	cs.DelImport("github.com/azhai/gozzo-db/prepare", "")
	err = cs.WriteTo(fname)
	return err
}

func GenNameType(col *schema.ColumnInfo, rule RuleConfig, nps map[string]NullPointer) string {
	if rule.Name == "" {
		rule.Name = ToCamel(col.FieldName)
	}
	if rule.Type == "" {
		rule.Type = base.GuessTypeName(col)
		if !col.IsNotNull() && len(nps) > 0 {
			if NullPointerMatch(nps, rule, col) {
				rule.Type = "*" + rule.Type // 字段可为NULL时，使用对应的指针类型
			}
		}
	}
	return rule.Name + " " + rule.Type
}

func GenTagComment(col *schema.ColumnInfo, rule RuleConfig) string {
	var blank, comment string
	if rule.Json == "" {
		rule.Json = col.FieldName
	}
	if rule.Tags == "" {
		tag := base.GuessStructTags(col)
		if rule.Name != "" { // 非常规属性名,需要设置字段名
			tag.Set("column", col.FieldName)
		}
		rule.Tags = tag.String("gorm")
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
