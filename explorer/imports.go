package explorer

import (
	"fmt"
	"go/ast"
	"go/token"
	"sort"
)

type ImportExplorer struct {
	imports map[string]struct{}
}

var (
	knownPackages = map[string]string{
		"context": "context",
		"errors":  "errors",
		"http":    "net/http",
		"json":    "encoding/json",
		"reflect": "reflect",
		"time":    "time",
	}
)

func (i *ImportExplorer) Visit(node ast.Node) (w ast.Visitor) {
	switch v := node.(type) {
	case *ast.SelectorExpr:
		switch x := v.X.(type) {
		case *ast.Ident:
			pack, ok := knownPackages[x.String()]
			if ok {
				i.imports[pack] = struct{}{}
			}
		}
	}
	return i
}

func NewImportExplorer() *ImportExplorer {
	return &ImportExplorer{
		imports: map[string]struct{}{},
	}
}

func (i *ImportExplorer) Explore(node ast.Node) {
	ast.Walk(i, node)
}

func (i *ImportExplorer) Spec() []ast.Spec {
	var imports []string
	for s := range i.imports {
		imports = append(imports, s)
	}
	sort.Strings(imports)
	var specs []ast.Spec
	for _, imp := range imports {
		specs = append(specs, &ast.ImportSpec{Path: &ast.BasicLit{
			Kind:  token.STRING,
			Value: fmt.Sprintf("\"%s\"", imp),
		}})
	}
	return specs
}
