package prepare

import (
	"bytes"
	"fmt"
	"sort"
	"strings"
	"text/template"

	"github.com/azhai/gozzo-db/construct"
	"github.com/azhai/gozzo-db/rewrite"
	"github.com/azhai/gozzo-db/schema"
	"github.com/azhai/gozzo-db/utils"
	"github.com/jinzhu/gorm"
)

func CreateModels(db *gorm.DB, app AppConfig) (names []string, err error) {
	var fname string
	s := schema.NewSchema(db.DB())
	for i, table := range s.GetTableNames("") {
		cols := s.GetColumnInfos(table, "")
		tbInfo := s.GetTableInfo(table, "")
		if app.TablePrefix != "" && strings.HasPrefix(table, app.TablePrefix) {
			table = table[len(app.TablePrefix):]
		}
		if app.PluralTable {
			table = utils.ToSingular(table)
		}
		name := utils.ToCamel(table)
		names = append(names, name)
		code, _ := GenModelCode(name, tbInfo, cols)

		fname = fmt.Sprintf("%s/%s.go", app.OutputDir, table)
		if i == 0 {
			utils.FileSize(fname, true)
		}
		cs := rewrite.NewCodeSource()
		err = cs.SetPackage("models")
		cs.AddImport("github.com/azhai/gozzo-db/construct", "base")
		cs.AddImport("github.com/jinzhu/gorm", "")
		cs.AddImport("time", "")
		err = cs.AddCode(code)
		cs.DelImport("github.com/jinzhu/gorm", "")
		cs.DelImport("time", "")
		err = cs.WriteTo(fname)
	}
	return
}

func GenModelCode(name string, table schema.TableInfo, columns []*schema.ColumnInfo) ([]byte, error) {
	funs := template.FuncMap{
		"genNameType": func(col *schema.ColumnInfo) string {
			name := utils.ToCamel(col.FieldName)
			if name == "DeletedAt" {
				return name + " *time.Time"
			}
			tpname := construct.GuessTypeName(col)
			return name + " " + tpname
		},
		"genTagComment": func(col *schema.ColumnInfo) string {
			var blank, comment string
			tags := construct.GuessStructTags(col)
			if tags != "" {
				blank = " "
			}
			if col.Comment != "" {
				comment = " // " + col.Comment
			}
			tpl := "`json:\"%s\"%s%s`%s"
			return fmt.Sprintf(tpl, col.FieldName, blank, tags, comment)
		},
	}
	data := map[string]interface{}{
		"Name": name, "Table": table, "Columns": columns,
	}
	var buf bytes.Buffer
	// New()名称要与ParseFiles()文件名一致 Funcs()必须在Parse之前
	t := template.New("gen_model.tmpl").Funcs(funs)
	t = template.Must(t.ParseFiles("gen_model.tmpl"))
	err := t.Execute(&buf, data)
	return buf.Bytes(), err
}

func GenInitFile(names []string, app AppConfig, fileName, diaName, drvName string) (err error) {
	var buf bytes.Buffer
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
	data := map[string]interface{}{
		"FileName": fileName, "DiaName": diaName,
		"Prefix": app.TablePrefix, "Plural": app.PluralTable,
		"Models": buf.String(),
	}

	buf.Reset()
	t := template.New("gen_init.tmpl")
	t = template.Must(t.ParseFiles("gen_init.tmpl"))
	if err = t.Execute(&buf, data); err == nil {
		fname := fmt.Sprintf("%s/init.go", app.OutputDir)
		cs := rewrite.NewCodeSource()
		err = cs.SetPackage("models")
		cs.AddImport("github.com/azhai/gozzo-db/prepare", "")
		cs.AddImport("github.com/jinzhu/gorm", "")
		cs.AddImport("github.com/jinzhu/gorm/dialects/"+drvName, "_")
		err = cs.AddCode(buf.Bytes())
		err = cs.WriteTo(fname)
	}
	return
}
