package rewrite

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/printer"
	"go/token"
	"io/ioutil"

	"golang.org/x/tools/go/ast/astutil"
	"golang.org/x/tools/imports"
)

func PrintNodes(objs ...interface{}) {
	for _, obj := range objs {
		if node, ok := obj.(ast.Node); ok {
			fmt.Printf("pos=[%d %d] %#v\n\n", node.Pos(), node.End(), node)
		} else {
			fmt.Printf("%#v\n\n", obj)
		}
	}
}

type CodeSource struct {
	Fileast *ast.File
	Fileset *token.FileSet
	Source  []byte
	*printer.Config
}

func NewCodeSource() *CodeSource {
	return &CodeSource{
		Fileset: token.NewFileSet(),
		Config: &printer.Config{
			Mode:     printer.TabIndent,
			Tabwidth: 4,
		},
	}
}

func (cs *CodeSource) SetSource(source []byte) (err error) {
	cs.Source = source
	cs.Fileast, err = parser.ParseFile(cs.Fileset, "", source, parser.ParseComments)
	return
}

func (cs *CodeSource) GetContent() ([]byte, error) {
	var buf bytes.Buffer
	err := cs.Config.Fprint(&buf, cs.Fileset, cs.Fileast)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (cs *CodeSource) AddCode(code []byte) error {
	content, err := cs.GetContent()
	if err != nil {
		return err
	}
	return cs.SetSource(append(content, code...))
}

func (cs *CodeSource) AddStringCode(code string) error {
	return cs.AddCode([]byte(code))
}

func (cs *CodeSource) WriteTo(filename string) error {
	code, err := cs.GetContent()
	if err != nil {
		return err
	}
	code, err = format.Source(code) // 格式化代码
	if err != nil {
		fmt.Println(err)
	}
	// imports分组排序
	var dst []byte
	err = ioutil.WriteFile(filename, code, 0644)
	if err != nil {
		return err
	}
	dst, err = imports.Process(filename, code, nil)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(filename, dst, 0644)
}

func (cs *CodeSource) GetPackage() string {
	if cs.Fileast != nil {
		return cs.Fileast.Name.Name
	}
	return ""
}

func (cs *CodeSource) GetNodeCode(node ast.Node) string {
	// 请先保证node不是nil
	pos := cs.Fileset.PositionFor(node.Pos(), false)
	end := cs.Fileset.PositionFor(node.End(), false)
	return string(cs.Source[pos.Offset:end.Offset])
}

func (cs *CodeSource) GetComment(c *ast.CommentGroup, trim bool) string {
	if c == nil {
		return ""
	}
	comment := cs.GetNodeCode(c)
	if trim {
		comment = TrimComment(comment)
	}
	return comment
}

func (cs *CodeSource) SetPackage(name string) (err error) {
	if cs.Fileast == nil {
		code := fmt.Sprintf("package %s", name)
		err = cs.SetSource([]byte(code))
	} else {
		cs.Fileast.Name.Name = name
	}
	return
}

func (cs *CodeSource) AddImport(path, alias string) bool {
	return astutil.AddNamedImport(cs.Fileset, cs.Fileast, alias, path)
}

func (cs *CodeSource) DelImport(path, alias string) bool {
	if astutil.UsesImport(cs.Fileast, path) {
		return false
	}
	return astutil.DeleteNamedImport(cs.Fileset, cs.Fileast, alias, path)
}
