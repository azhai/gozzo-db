package construct

import (
	"bytes"
	"fmt"
	"sort"
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
	Name           string
	Import, Alias  string
	Features       []string
	sortedFeatures []string
	FieldLines     []string
	IsChanged      bool
}

func NewClassSummary(name string) ClassSummary {
	return ClassSummary{Name: name}
}

func (s ClassSummary) GetInnerCode() string {
	var buf bytes.Buffer
	for _, line := range s.FieldLines {
		buf.WriteString(fmt.Sprintf("\t%s\n", line))
	}
	return buf.String()
}

func (s ClassSummary) GetSortedFeatures() []string {
	if len(s.sortedFeatures) == 0 {
		s.sortedFeatures = append([]string{}, s.Features...)
		sort.Strings(s.sortedFeatures)
	}
	return s.sortedFeatures
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
			s.Features[i] = ps[0]
		} else {
			s.Features[i] = ps[0] + " " + ps[1]
		}
		if cm := cp.GetComment(f.Comment, true); len(cm) > 0 {
			code += " //" + cm
		}
		s.FieldLines[i] = code
	}
	return size
}

func ReplaceModel(summary, sub ClassSummary) ClassSummary {
	var features, lines []string
	find := false
	sted := sub.GetSortedFeatures()
	for i, ft := range summary.Features {
		if !InStringList(ft, sted, CMP_STRING_EQUAL) {
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
			size := summary.ParseFields(cp, node)
			for n, sub := range ModelClasses {
				if n == summary.Name {
					continue
				}
				sted := sub.GetSortedFeatures()
				sorted := summary.GetSortedFeatures()
				if IsSubsetList(sted, sorted) {
					summary = ReplaceModel(summary, sub)
				} else if strings.HasPrefix(n, "base.") || n == summary.Name {
					continue
				} else if IsSubsetList(sorted, sted) {
					ModelClasses[n] = ReplaceModel(sub, summary)
				}
			}
			ModelClasses[name] = summary
			if summary.IsChanged && size > 0 {
				// cp.ReplaceCode(node.Fields[0], node.Fields[size-1], summary.GetInnerCode())
			}
		}
		// cp.WriteSource(fname)
	}
}
