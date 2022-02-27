package builders

import (
	"fmt"
	"go/ast"
	"go/token"
	"strings"
)

var (
	// Zero equals 0 constant
	Zero = IntegerConstant(0).Expr()
	// Blank equals _ ident
	Blank = ast.NewIdent("_")
	// Nil equals nil ident
	Nil = ast.NewIdent("nil")
	// EmptyInterface equals empty interface
	EmptyInterface = &ast.InterfaceType{
		Methods:    &ast.FieldList{},
		Incomplete: true, // todo ? check it
	}

	// UInt represents the data type uint
	UInt = ast.NewIdent("uint")
	// UInt8 represents the data type uint8
	UInt8 = ast.NewIdent("uint8")
	// UInt16 represents the data type uint16
	UInt16 = ast.NewIdent("uint16")
	// UInt32 represents the data type uint32
	UInt32 = ast.NewIdent("uint32")
	// UInt64 represents the data type uint64
	UInt64 = ast.NewIdent("uint64")

	// Int represents the data type int
	Int = ast.NewIdent("int")
	// Int8 represents the data type int8
	Int8 = ast.NewIdent("int8")
	// Int16 represents the data type int16
	Int16 = ast.NewIdent("int16")
	// Int32 represents the data type int32
	Int32 = ast.NewIdent("int32")
	// Int64 represents the data type int64
	Int64 = ast.NewIdent("int64")

	// Float32 represents the data type float32
	Float32 = ast.NewIdent("float32")
	// Float64 represents the data type float64
	Float64 = ast.NewIdent("float64")

	// String represents the data type string
	String = ast.NewIdent("string")

	// ContextType represents the `context.Context` interface
	ContextType = SimpleSelector("context", "Context")

	// ErrorType represents the `error` interface
	ErrorType = ast.NewIdent("error")
)

// Import represents import declaration with token.IMPORT
func Import(imports map[string]string) ast.Decl {
	var impSpec []ast.Spec
	impSpec = makeImportSpec(imports)
	return &ast.GenDecl{
		Tok:   token.IMPORT,
		Specs: impSpec,
	}
}

func makeImportSpec(imports map[string]string) []ast.Spec {
	var impSpec = make([]ast.Spec, 0, len(imports))
	for packageKey, packagePath := range imports {
		pathSplit := strings.Split(packagePath, "/")
		impElm := ast.ImportSpec{
			Path: &ast.BasicLit{
				Kind:  token.STRING,
				Value: fmt.Sprintf("\"%s\"", packagePath),
			},
		}
		// fixme: in general - this is lie
		if pathSplit[len(pathSplit)-1] == packageKey {
			impElm.Name = ast.NewIdent(packageKey)
		}
		impSpec = append(impSpec, &impElm)
	}
	return impSpec
}

// CommentGroup wraps the lines in the ast.CommentGroup structure. Returns nil if arguments is omitted or empty
func CommentGroup(comments ...string) *ast.CommentGroup {
	if len(comments) == 0 {
		return nil
	} else {
		if len(comments) == 1 && strings.TrimSpace(comments[0]) == "" {
			return nil
		}
	}
	var prefChar = "\n// "
	var g ast.CommentGroup
	for _, line := range comments {
		g.List = append(g.List, &ast.Comment{Text: prefChar + line, Slash: 1})
		prefChar = "// "
	}
	return &g
}

// ast.Field constructor.
// docAndComments contains the first line as Docstring, all other lines turn into CommentGroup
func Field(name string, tag *ast.BasicLit, fieldType ast.Expr, docAndComments ...string) *ast.Field {
	if fieldType == nil {
		return nil
	}
	var (
		doc      = ""
		comments []string
		names    = make([]*ast.Ident, 0, 1)
	)
	if name != "" {
		names = []*ast.Ident{ast.NewIdent(name)}
	}
	if docAndComments = truncateEmpty(docAndComments); len(docAndComments) > 0 {
		doc = fmt.Sprintf("%s %s", name, docAndComments[0])
		comments = docAndComments[1:]
	}
	return &ast.Field{
		Doc:     CommentGroup(doc),
		Names:   names,
		Type:    fieldType,
		Tag:     tag,
		Comment: CommentGroup(comments...),
	}
}

// creates ast.FieldList, any nil values will be excluded from list
func FieldList(fields ...*ast.Field) *ast.FieldList {
	var list = ast.FieldList{
		List: make([]*ast.Field, 0, len(fields)),
	}
	for i, field := range fields {
		if field != nil {
			list.List = append(list.List, fields[i])
		}
	}
	return &list
}

// creates ast.TypeSpec with Type field
func TypeSpec(name string, varType ast.Expr, comment ...string) *ast.TypeSpec {
	return &ast.TypeSpec{
		Name: ast.NewIdent(name),
		Type: varType,
		Doc:  CommentGroup(comment...),
	}
}

// creates ast.ValueSpec with Type field
func VariableType(name string, varType ast.Expr, vals ...VarValue) *ast.ValueSpec {
	valSpec := ast.ValueSpec{
		Names: []*ast.Ident{ast.NewIdent(name)},
		Type:  varType,
	}
	for _, val := range vals {
		valSpec.Values = append(valSpec.Values, val.Expr())
	}
	return &valSpec
}

// creates ast.ValueSpec with Values field
func VariableValue(name string, vals ...VarValue) *ast.ValueSpec {
	valSpec := ast.ValueSpec{
		Names: []*ast.Ident{
			ast.NewIdent(name),
		},
		Values: []ast.Expr{},
	}
	for _, val := range vals {
		valSpec.Values = append(valSpec.Values, val.Expr())
	}
	return &valSpec
}

// creates ast.TypeSpec with Type field
func StructType(fields ...*ast.Field) *ast.StructType {
	return &ast.StructType{
		Fields: FieldList(fields...),
	}
}

type assignToken int

const (
	Assignment  assignToken = iota + 1 // =
	Incremental                        // +=
	Decremental                        // -=
	Definition                         // :=
)

func (t assignToken) token() token.Token {
	switch t {
	case Assignment:
		return token.ASSIGN
	case Incremental:
		return token.ADD_ASSIGN
	case Decremental:
		return token.SUB_ASSIGN
	case Definition:
		return token.DEFINE
	default:
		panic("unknown assignment token")
	}
}

type (
	// represents a list of variable names
	VarNames []string
)

func (c VarNames) expression() []ast.Expr {
	var varNames = make([]ast.Expr, 0, len(c))
	for _, varName := range c {
		varNames = append(varNames, ast.NewIdent(varName))
	}
	return varNames
}

// creates ast.AssignStmt which assigns a variable with a value
func Assign(varNames VarNames, tok assignToken, rhs ...ast.Expr) ast.Stmt {
	return &ast.AssignStmt{
		Lhs: varNames.expression(),
		Tok: tok.token(),
		Rhs: rhs,
	}
}

// todo move
func truncateEmpty(s []string) []string {
	var result []string
	for i := range s {
		if line := strings.TrimSpace(s[i]); len(result) > 0 || line != "" {
			result = append(result, line)
		}
	}
	return result
}
