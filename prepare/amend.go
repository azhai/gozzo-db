package prepare

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"

	"github.com/azhai/gozzo-db/rewrite"
	"github.com/azhai/gozzo-db/schema"
	"github.com/azhai/gozzo-db/utils"
)

var (
	sourceDir  = "models/"
	targetFile = "amend/models.go"
)

// 将代码中的表和字段注释，覆盖数据库的SQL语句注释
func AmendComments(db *sql.DB, execute, verbose bool) {
	tables, colDefs := FindTables(db, verbose)
	fname := utils.GetAbsFile(targetFile)
	if fsize, _ := utils.FileSize(fname, true); fsize == 0 {
		_ = CollectCode(sourceDir, fname, "gorm.Model", verbose)
	}
	tbComments, colComments := ParseModelComments(fname, verbose)
	atSql, tpl := "", "CHANGE `%s` `%s` %s COMMENT '%s',\n"
	for i, table := range tables {
		columns, comments := colDefs[i], colComments[table]
		atSql = fmt.Sprintf("ALTER TABLE `%s`\n", table)
		for col, def := range columns {
			if comment, ok := comments[col]; ok {
				atSql += fmt.Sprintf(tpl, col, col, def, comment)
			}
		}
		atSql += fmt.Sprintf("COMMENT = '%s';", tbComments[table])
		fmt.Println(atSql, "\n\n")
		if execute {
			db.Exec(atSql)
		}
	}
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
	var data string = `package amend

import (
	"time"
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
func ParseModelComments(filename string, verbose bool) (map[string]string, map[string](map[string]string)) {
	tbComments := make(map[string]string)
	colComments := make(map[string](map[string]string))
	cp, err := rewrite.NewFileParser(filename)
	if !utils.CheckError(err) {
		return tbComments, colComments
	}
	for _, node := range cp.AllDeclNode("type") {
		table := utils.ToSnake(strings.Join(node.Names, ", "))
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
