package explorer

import (
	"fmt"
	"go/ast"
	"go/token"
	"sort"
	"strings"
)

type (
	Discoverer struct {
		imports map[string]UsedPackage
	}
	UsedPackage struct {
		Package Package
		Alias   string
	}
	Package struct {
		Path string
		Kind PkgKind
	}
	PkgKind int8
)

const (
	PkgKindSystem PkgKind = iota
	PkgKindExternal
	PkgKindInternal
)

var (
	knownPackages = map[string]Package{
		"tar":       {Path: "archive/tar", Kind: PkgKindSystem},
		"zip":       {Path: "archive/zip", Kind: PkgKindSystem},
		"bufio":     {Path: "bufio", Kind: PkgKindSystem},
		"bytes":     {Path: "bytes", Kind: PkgKindSystem},
		"context":   {Path: "context", Kind: PkgKindSystem},
		"aes":       {Path: "crypto/aes", Kind: PkgKindSystem},
		"des":       {Path: "crypto/des", Kind: PkgKindSystem},
		"md5":       {Path: "crypto/md5", Kind: PkgKindSystem},
		"crand":     {Path: "crypto/rand", Kind: PkgKindSystem},
		"sha512":    {Path: "crypto/sha512", Kind: PkgKindSystem},
		"x509":      {Path: "crypto/x509", Kind: PkgKindSystem},
		"sql":       {Path: "database/sql", Kind: PkgKindSystem},
		"hex":       {Path: "encoding/hex", Kind: PkgKindSystem},
		"json":      {Path: "encoding/json", Kind: PkgKindSystem},
		"xml":       {Path: "encoding/xml", Kind: PkgKindSystem},
		"errors":    {Path: "errors", Kind: PkgKindSystem},
		"fmt":       {Path: "fmt", Kind: PkgKindSystem},
		"color":     {Path: "image/color", Kind: PkgKindSystem},
		"gif":       {Path: "image/gif", Kind: PkgKindSystem},
		"jpeg":      {Path: "image/jpeg", Kind: PkgKindSystem},
		"png":       {Path: "image/png", Kind: PkgKindSystem},
		"io":        {Path: "io", Kind: PkgKindSystem},
		"log":       {Path: "log", Kind: PkgKindSystem},
		"math":      {Path: "math", Kind: PkgKindSystem},
		"big":       {Path: "math/big", Kind: PkgKindSystem},
		"rand":      {Path: "math/rand", Kind: PkgKindSystem},
		"mime":      {Path: "mime", Kind: PkgKindSystem},
		"multipart": {Path: "mime/multipart", Kind: PkgKindSystem},
		"net":       {Path: "net", Kind: PkgKindSystem},
		"http":      {Path: "net/http", Kind: PkgKindSystem},
		"url":       {Path: "net/url", Kind: PkgKindSystem},
		"os":        {Path: "os", Kind: PkgKindSystem},
		"path":      {Path: "path", Kind: PkgKindSystem},
		"reflect":   {Path: "reflect", Kind: PkgKindSystem},
		"regexp":    {Path: "regexp", Kind: PkgKindSystem},
		"sort":      {Path: "sort", Kind: PkgKindSystem},
		"strconv":   {Path: "strconv", Kind: PkgKindSystem},
		"sync":      {Path: "sync", Kind: PkgKindSystem},
		"time":      {Path: "time", Kind: PkgKindSystem},
		"unicode":   {Path: "unicode", Kind: PkgKindSystem},
		"utf8":      {Path: "unicode/utf8", Kind: PkgKindSystem},
		"utf16":     {Path: "unicode/utf16", Kind: PkgKindSystem},
		"unsafe":    {Path: "unsafe", Kind: PkgKindSystem},
		"fasthttp":  {Path: "github.com/valyala/fasthttp", Kind: PkgKindExternal},
		"fastjson":  {Path: "github.com/valyala/fastjson", Kind: PkgKindExternal},
		"router":    {Path: "github.com/fasthttp/router", Kind: PkgKindExternal},
		"uuid":      {Path: "github.com/google/uuid", Kind: PkgKindExternal},
	}
)

func RegisterPackage(packName string, pkg Package) {
	knownPackages[packName] = pkg
}

func New() *Discoverer {
	return &Discoverer{
		imports: make(map[string]UsedPackage),
	}
}

func (i *Discoverer) Explore(node ast.Node) {
	ast.Walk(i, node)
}

func (i *Discoverer) Visit(node ast.Node) (w ast.Visitor) {
	sel, ok := node.(*ast.SelectorExpr)
	if !ok {
		return i
	}
	x, ok := sel.X.(*ast.Ident)
	if !ok {
		return i
	}
	pack, ok := knownPackages[x.String()]
	if ok {
		i.imports[pack.Path] = UsedPackage{
			Package: pack,
			Alias:   x.String(),
		}
	}
	return i
}

func (i *Discoverer) ImportSpec() []ast.Spec {
	var imports []UsedPackage
	for _, pkg := range i.imports {
		imports = append(imports, pkg)
	}
	sort.SliceStable(imports, func(i, j int) bool {
		if imports[i].Package.Kind == imports[j].Package.Kind {
			return imports[i].Package.Path < imports[j].Package.Path
		}
		return imports[i].Package.Kind < imports[j].Package.Kind
	})

	var (
		currT PkgKind = -1
		specs []ast.Spec
	)
	for _, imp := range imports {
		var addLine string
		var alias string
		split := strings.Split(imp.Package.Path, "/")
		if split[len(split)-1] != imp.Alias {
			alias = imp.Alias + " "
		}
		if currT != imp.Package.Kind {
			currT = imp.Package.Kind
			addLine = "\n\t"
		}
		specs = append(specs, &ast.ImportSpec{Path: &ast.BasicLit{
			Kind:  token.STRING,
			Value: fmt.Sprintf("%s%s\"%s\"", addLine, alias, imp.Package.Path),
		}})
	}
	return specs
}
