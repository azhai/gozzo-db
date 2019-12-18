package prepare

import (
	"bytes"
	"database/sql"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/azhai/gozzo-db/rewrite"
	"github.com/azhai/gozzo-db/schema"
	"github.com/azhai/gozzo-db/utils"
)

var (
	sourceDir  = "models/"
	targetDir = "cmd/tmp/"
)

// 将代码中的表和字段注释，覆盖数据库的SQL语句注释
func AmendComments(db *sql.DB, prefix string, onlyTable, verbose bool) error {
	var buf bytes.Buffer
	tables, colDefs := FindTables(db, verbose)
	fname := utils.GetAbsFile(filepath.Join(targetDir, "models.go"))
	if fsize := utils.MkdirForFile(fname); fsize == 0 {
		err := CollectCode(sourceDir, fname, "", verbose)
		if err != nil {
			return err
		}
	}
	var fp *os.File
	outname := utils.GetAbsFile(filepath.Join(targetDir, "comments.sql"))
	flag := os.O_WRONLY | os.O_CREATE | os.O_TRUNC | os.O_APPEND
	fp, err := os.OpenFile(outname, flag, 0644)
	defer fp.Close()
	if err != nil {
		return err
	}

	tbComments, colComments := ParseModelComments(fname, prefix, verbose)
	tpl := "\tCHANGE `%s` `%s` %s COMMENT '%s',\n"
	for i, table := range tables {
		columns, comments := colDefs[i], colComments[table]
		buf.Reset()
		buf.WriteString(fmt.Sprintf("ALTER TABLE `%s`", table))
		if !onlyTable {
			buf.WriteString("\n")
			for col, def := range columns {
				if comment, ok := comments[col]; ok {
					buf.WriteString(fmt.Sprintf(tpl, col, col, def, comment))
				}
			}
		}
		buf.WriteString(fmt.Sprintf(" COMMENT = '%s';\n", tbComments[table]))
		tbSql := buf.String() + "\n"
		fp.WriteString(tbSql)
		if verbose {
			fmt.Print(tbSql)
		}
	}
	return nil
}

// 找出所有表名，并列出每张表的字段定义
func FindTables(db *sql.DB, verbose bool) ([]string, []map[string]string) {
	s := schema.NewSchema(db)
	tables := s.GetTableNames("")
	defines := make([]map[string]string, 0)
	if !utils.CheckError(s.Error) {
		return tables, defines
	}
	for _, table := range tables {
		cols := s.GetColumnInfos(table, "")
		if !utils.CheckError(s.Error) {
			return tables, defines
		}
		if verbose {
			fmt.Println(table)
		}
		tbDef := make(map[string]string)
		for _, ci := range cols {
			tbDef[ci.Name()] = ci.GetDefine()
			if verbose {
				fmt.Printf("Column: %s %s COMMENT %s\n",
					ci.Name(), ci.GetDefine(), ci.Comment)
			}
		}
		defines = append(defines, tbDef)
	}
	return tables, defines
}

// 在目录下所有的go代码中，找出baseModel子类的完整代码写入一个文件
func CollectCode(dirname, outname, baseModel string, verbose bool) (err error) {
	var data string = `package tmp

import (
	"time"

	base "github.com/azhai/gozzo-db/construct"
	"github.com/jinzhu/gorm"
)
`
	files, err := ioutil.ReadDir(dirname)
	if !utils.CheckError(err) {
		return
	}
	for _, file := range files {
		fname := file.Name()
		if !strings.HasSuffix(fname, ".go") {
			continue
		}
		fname = filepath.Join(dirname, fname)
		if verbose {
			fmt.Println(fname)
		}
		cp, err := rewrite.NewFileParser(fname)
		if !utils.CheckError(err) {
			continue
		}
		for _, node := range cp.AllDeclNode("type") {
			if len(node.Fields) == 0 {
				continue
			}
			if baseModel != "" { // 必须是某个Model的子类
				ffcode := cp.GetNodeCode(node.Fields[0])
				if strings.TrimSpace(ffcode) != baseModel {
					continue
				}
			}
			data += "\n\n"
			if node.Comment != nil {
				data += cp.GetComment(node.Comment, false) + "\n"
			}
			data += cp.GetNodeCode(node)
		}
	}
	err = ioutil.WriteFile(outname, []byte(data+"\n"), 0644)
	return
}

// 分析文件代码，找出所有Model注释和其中的字段注释
func ParseModelComments(filename, prefix string, verbose bool) (map[string]string, map[string](map[string]string)) {
	tbComments := make(map[string]string)
	colComments := make(map[string](map[string]string))
	cp, err := rewrite.NewFileParser(filename)
	if !utils.CheckError(err) {
		return tbComments, colComments
	}
	for _, node := range cp.AllDeclNode("type") {
		table := prefix + utils.ToSnake(strings.Join(node.Names, ", "))
		tbComments[table] = cp.GetComment(node.Comment, true)
		if verbose {
			fmt.Printf("%s %s\n", table, tbComments[table])
		}
		colComments[table] = make(map[string]string)
		for i, fd := range node.Fields {
			comment := cp.GetComment(fd.Comment, true)
			if comment == "" {
				continue
			}
			column := utils.ToSnake(strings.Join(fd.Names, ", "))
			colComments[table][column] = comment
			if verbose {
				fmt.Printf("%d: %s %s\n", i, column, comment)
			}
		}
		if verbose {
			fmt.Println("")
		}
	}
	return tbComments, colComments
}
