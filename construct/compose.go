package construct

import (
	"bytes"
	"fmt"
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
		Features: []string{
			"CreatedAt time.Time",
			"UpdatedAt time.Time",
			"DeletedAt *time.Time",
		},
		FieldLines: []string{
			"CreatedAt time.Time  `json:\"-\" gorm:\"comment:'创建时间'\"`       // 创建时间",
			"UpdatedAt time.Time  `json:\"-\" gorm:\"comment:'更新时间'\"`       // 更新时间",
			"DeletedAt *time.Time `json:\"-\" gorm:\"index;comment:'删除时间'\"` // 删除时间",
		},
	},
	"base.NestedModel": {
		Name:   "base.NestedModel",
		Import: "github.com/azhai/gozzo-db/construct",
		Alias:  "base",
		Features: []string{
			"Lft uint",
			"Rgt uint",
			"Depth uint8",
		},
		FieldLines: []string{
			"Lft   uint  `json:\"-\" gorm:\"not null;default:0;comment:'左边界'\"`          // 左边界",
			"Rgt   uint  `json:\"-\" gorm:\"not null;index;default:0;comment:'右边界'\"`    // 右边界",
			"Depth uint8 `json:\"depth\" gorm:\"not null;index;default:1;comment:'高度'\"` // 高度",
		},
	},
}

type ClassSummary struct {
	Name          string
	Import, Alias string
	Features        []string
	FieldLines      []string
	IsChanged      bool
}

func NewClassSummary(name string) ClassSummary {
	return ClassSummary{Name: name}
}

func (s ClassSummary) String() string {
	var buf bytes.Buffer
	sep := ""
	if len(s.Alias) > 0 {
		sep = " "
	}
	buf.WriteString(fmt.Sprintf("import %s%s\"%s\"\n", s.Alias, sep, s.Import))
	buf.WriteString(fmt.Sprintf("%s %v\n", s.Name, s.IsChanged))
	buf.WriteString(fmt.Sprintf("%#v\n", s.Features))
	buf.WriteString(fmt.Sprintf("%#v\n", s.FieldLines))
	return buf.String()
}

func (s *ClassSummary) ParseFields(cp *rewrite.CodeParser, node *rewrite.DeclNode) int {
	size := len(node.Fields)
	s.Features = make([]string, size)
	s.FieldLines = make([]string, size)
	for i, f := range node.Fields {
		code := cp.GetNodeCode(f)
		ps := strings.Fields(code)
		if len(ps) == 0 {
			continue
		}
		if len(ps) == 1 {
			s.Features[i]  = ps[0]
		} else {
			s.Features[i]  = ps[0] + " " + ps[1]
		}
		comment := cp.GetComment(f.Comment, true)
		s.FieldLines[i] = code + " " + comment
	}
	return size
}

func ReplaceModel(summary, sub ClassSummary) ClassSummary {
	var features, lines []string
	find := false
	for i, ft := range summary.Features {
		if !InStringList(ft, sub.Features, CMP_STRING_EQUAL) {
			features = append(features, ft)
			lines = append(lines, summary.FieldLines[i])
		} else if !find {
			features = append(features, sub.Name)
			lines = append(lines, sub.Name)
			find = true
			summary.IsChanged = true
		}
	}
	summary.Features, summary.FieldLines = features, lines
	return summary
}

func ScanModelDir(dir string) {
	files, _ := FindFiles(dir, ".go")
	for fname, _ := range files {
		cp, err := rewrite.NewFileParser(fname)
		if err != nil {
			fmt.Println(fname, " error: ", err)
			continue
		}
		for _, node := range cp.AllDeclNode("type") {
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
			summary := NewClassSummary(name)
			summary.ParseFields(cp, node)
			for n, s := range ModelClasses {
				if IsSubsetList(s.Features, summary.Features) {
					summary = ReplaceModel(summary, s)
				} else if !strings.HasPrefix(n, "base.") &&
						IsSubsetList(summary.Features, s.Features) {
					ModelClasses[n] = ReplaceModel(s, summary)
				}
			}
			ModelClasses[name] = summary
		}
	}
}
