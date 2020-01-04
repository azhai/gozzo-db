package construct

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/azhai/gozzo-db/rewrite"
)

/**
 * 时间相关的三个典型字段
 */
type TimeModel struct {
	CreatedAt time.Time  `json:"-" gorm:"comment:'创建时间'"`       // 创建时间
	UpdatedAt time.Time  `json:"-" gorm:"comment:'更新时间'"`       // 更新时间
	DeletedAt *time.Time `json:"-" gorm:"index;comment:'删除时间'"` // 删除时间
}

var ModelClasses = map[string]ClassSummary{
	"base.TimeModel": {
		Name:   "base.TimeModel",
		Import: "github.com/azhai/gozzo-db/construct",
		Alias:  "base",
		Fields: map[string]string{
			"CreatedAt": "time.Time",
			"UpdatedAt": "time.Time",
			"DeletedAt": "*time.Time",
		},
	},
	"base.NestedModel": {
		Name:   "base.NestedModel",
		Import: "github.com/azhai/gozzo-db/construct",
		Alias:  "base",
		Fields: map[string]string{
			"Lft":   "uint",
			"Rgt":   "uint",
			"Depth": "uint8",
		},
	},
}

type ClassSummary struct {
	Name          string
	Import, Alias string
	Fields        map[string]string
}

func (s ClassSummary) String() string {
	var buf bytes.Buffer
	sep := ""
	if len(s.Alias) > 0 {
		sep = " "
	}
	buf.WriteString(fmt.Sprintf("import %s%s\"%s\"\n", s.Alias, sep, s.Import))
	buf.WriteString(fmt.Sprintf("%s\n", s.Name))
	buf.WriteString(fmt.Sprintf("%#v\n", s.Fields))
	return buf.String()
}

func ScanModelDir(dir string) {
	files, _ := FindFiles(dir, "*.go")
	for fname, _ := range files {
		cp, err := rewrite.NewFileParser(fname)
		if err != nil {
			continue
		}
		for _, node := range cp.AllDeclNode("type.struct") {
			if len(node.Fields) == 0 {
				continue
			}
			name := node.GetName()
			ffcode := cp.GetNodeCode(node.Fields[0])
			first := strings.TrimSpace(ffcode)
			if !strings.HasSuffix(name, "Model") &&
				!strings.HasSuffix(first, "Model") {
				continue
			}
			summary := ClassSummary{Name:name}
			for _, f := range node.Fields {
				pieces := strings.SplitN(cp.GetNodeCode(f), " ", 3)
				if len(pieces) >= 2 {
					summary.Fields[pieces[0]] = pieces[1]
				}
			}
			ModelClasses[name] = summary
		}
	}
}

// 遍历目录下的文件
func FindFiles(dir, ext string) (map[string]os.FileInfo, error) {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	var result = make(map[string]os.FileInfo)
	for _, file := range files {
		fname := file.Name()
		if !strings.HasSuffix(fname, ext) {
			continue
		}
		fname = filepath.Join(dir, fname)
		result[fname] = file
	}
	return result, nil
}
