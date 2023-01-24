package asthlp

import (
	"fmt"
	"go/ast"
	"go/token"
	"strconv"
	"strings"
)

type (
	Expression interface {
		Expr() ast.Expr
	}
	StructFiller interface {
		Expression
		FillKeyValue(key string, value ast.Expr) StructFiller
	}
	BoolConstant     bool
	StringConstant   string  // string constant e.g. "abc"
	RuneConstant     rune    // rune constant e.g. 'r'
	IntegerConstant  int64   // integer constant e.g. 123
	UnsignedConstant uint64  // unsigned integer constant e.g. 123
	FloatConstant    float64 // float constant e.g. 123.45
	SliceByteLiteral []byte  // []byte{'f', 'i', 'l', 't', 'e', 'r'}
	VariableName     string  // any variable name
	freeExpression   struct {
		expr ast.Expr
	}

	structLiteral struct {
		name ast.Expr
		exps []ast.Expr
	}
)

func (b BoolConstant) Expr() ast.Expr {
	return &ast.Ident{Name: strconv.FormatBool(bool(b))}
}

// Expr creates ast.BasicLit with token.STRING
func (c StringConstant) Expr() ast.Expr {
	if strings.Contains(string(c), "\"") || strings.Contains(string(c), "\n") {
		return &ast.BasicLit{
			Kind:  token.STRING,
			Value: fmt.Sprintf("`%s`", c),
		}
	} else {
		return &ast.BasicLit{
			Kind:  token.STRING,
			Value: fmt.Sprintf("\"%s\"", c),
		}
	}
}

// Expr creates ast.BasicLit with token.STRING
func (c RuneConstant) Expr() ast.Expr {
	return &ast.BasicLit{
		Kind:  token.STRING,
		Value: "'" + string(c) + "'",
	}
}

// Expr creates ast.BasicLit with token.INT
func (c IntegerConstant) Expr() ast.Expr {
	return &ast.BasicLit{
		ValuePos: 1,
		Kind:     token.INT,
		Value:    fmt.Sprintf("%d", c),
	}
}

// Expr creates ast.BasicLit with token.INT
func (c UnsignedConstant) Expr() ast.Expr {
	return &ast.BasicLit{
		ValuePos: 1,
		Kind:     token.INT,
		Value:    fmt.Sprintf("%d", c),
	}
}

// Expr creates ast.BasicLit with token.FLOAT
func (c FloatConstant) Expr() ast.Expr {
	return &ast.BasicLit{
		ValuePos: 1,
		Kind:     token.FLOAT,
		Value:    fmt.Sprintf("%f", c),
	}
}

func (s SliceByteLiteral) Expr() ast.Expr {
	var elts []ast.Expr
	for _, char := range s {
		var val = "'" + string(char) + "'"
		var kind = token.CHAR
		if char < 32 {
			val = strconv.Itoa(int(char))
			kind = token.INT
		}
		elts = append(elts, &ast.BasicLit{
			Kind:  kind,
			Value: val,
		})
	}
	return &ast.CompositeLit{
		Type: &ast.ArrayType{Elt: Byte},
		Elts: elts,
	}
}

// Expr creates ast.Ident with variable name
func (c VariableName) Expr() ast.Expr {
	return ast.NewIdent(string(c))
}

// Expr creates ast.Expr
func (c freeExpression) Expr() ast.Expr {
	return c.expr
}

// FreeExpression creates Expression from ast.Expr
func FreeExpression(e ast.Expr) Expression {
	return freeExpression{expr: e}
}

func StructLiteral(name ast.Expr) StructFiller {
	return &structLiteral{
		name: name,
	}
}

func (c *structLiteral) Expr() ast.Expr {
	return &ast.CompositeLit{
		Type: c.name,
		Elts: c.exps,
	}
}

func (c *structLiteral) FillKeyValue(key string, value ast.Expr) StructFiller {
	c.exps = append(c.exps, &ast.KeyValueExpr{
		Key:   ast.NewIdent(key),
		Value: value,
	})
	return c
}

// Index creates the array element picker expression
//   someArr[1]
func Index(x ast.Expr, index Expression) ast.Expr {
	return &ast.IndexExpr{
		X:      x,
		Lbrack: 1,
		Index:  safeExpr(index),
		Rbrack: 2,
	}
}

// SimpleSelector represents a dot notation expression like "pack.object" from string arguments
func SimpleSelector(pack, object string) ast.Expr {
	return Selector(ast.NewIdent(pack), object)
}

// Selector represents a dot notation expression like "pack.object"
// <x>.<object>
func Selector(x ast.Expr, object string) ast.Expr {
	return &ast.SelectorExpr{
		X:   x,
		Sel: ast.NewIdent(object),
	}
}

// Unary represents unary expression
//   <tok><expr> e.g. !expr
// you can use this constant as `tok` attribute:
//   token.ADD     // +
//   token.SUB     // -
//   token.MUL     // *
//   token.QUO     // /
//   token.REM     // %
//   token.AND     // &
//   token.OR      // |
//   token.XOR     // ^
//   token.SHL     // <<
//   token.SHR     // >>
//   token.AND_NOT // &^
func Unary(expr ast.Expr, tok token.Token) ast.Expr {
	if tok == token.MUL {
		return Star(expr)
	}
	return &ast.UnaryExpr{
		OpPos: 1,
		Op:    tok,
		X:     expr,
	}
}

// Star represents star expression
//   *<expr>
func Star(expr ast.Expr) ast.Expr {
	return &ast.StarExpr{
		Star: 1,
		X:    expr,
	}
}

// Ref represents reference
//   &<expr>
func Ref(expr ast.Expr) ast.Expr {
	return Unary(expr, token.AND)
}

// Not represents inversion
//   !<expr>
func Not(expr ast.Expr) ast.Expr {
	return Unary(expr, token.NOT)
}

// Binary represents binary expression. Use token.* constants as `tok` attribute
//   <left> <tok> <right> e.g. left == right
func Binary(left, right ast.Expr, tok token.Token) ast.Expr {
	if left == nil || right == nil {
		panic("unsupported")
	}
	return &ast.BinaryExpr{
		X:     left,
		OpPos: 1,
		Op:    tok,
		Y:     right,
	}
}

// ArrayType represents array expression, use `l` attribute if you want to specify array length, else omit
//   [<l>]<expr>
func ArrayType(expr ast.Expr, l ...ast.Expr) ast.Expr {
	var lenExpr ast.Expr = nil
	if len(l) > 0 {
		lenExpr = l[0]
		if len(l) > 1 {
			panic("allowed only one value")
		}
	}
	return &ast.ArrayType{
		Lbrack: 1,
		Len:    lenExpr,
		Elt:    expr,
	}
}

// MapType represents map expression
//   map[<T>]<expr>
func MapType(key, expr ast.Expr) ast.Expr {
	return &ast.MapType{
		Map:   1,
		Key:   key,
		Value: expr,
	}
}

// NotEqual represents comparison operation
//   <left> != <right>
func NotEqual(left, right ast.Expr) ast.Expr {
	return Binary(left, right, token.NEQ)
}

// Equal represents comparison operation
//   <left> == <right>
func Equal(left, right ast.Expr) ast.Expr {
	return Binary(left, right, token.EQL)
}

// Great represents comparison operation
//   <left> > <right>
func Great(left, right ast.Expr) ast.Expr {
	return Binary(left, right, token.GTR)
}

// Add represents an addition operation
//   <expr1> + <expr2> + <expr3>
func Add(exps ...ast.Expr) ast.Expr {
	var acc ast.Expr = nil
	for _, expr := range exps {
		if acc == nil {
			acc = expr
		} else {
			acc = Binary(acc, expr, token.ADD)
		}
	}
	return acc
}

// Sub represents a subtraction operation
//   <expr1> - <expr2> - <expr3>
func Sub(exps ...ast.Expr) ast.Expr {
	var acc ast.Expr = nil
	for _, expr := range exps {
		if acc == nil {
			acc = expr
		} else {
			acc = Binary(acc, expr, token.SUB)
		}
	}
	return acc
}

// NotNil represents non-nil-comparison operation
//   <expr> != nil
func NotNil(expr ast.Expr) ast.Expr {
	return Binary(expr, Nil, token.NEQ)
}

// IsNil represents nil-comparison operation
//   <expr> == nil
func IsNil(expr ast.Expr) ast.Expr {
	return Binary(expr, Nil, token.EQL)
}

// And represents `&&` in comparison operation
//   <expr> && <expr> && <expr>
func And(left ast.Expr, expr ...ast.Expr) ast.Expr {
	if len(expr) == 0 {
		return left
	}
	return Binary(left, And(expr[0], expr[1:]...), token.LAND)
}

// Or represents `||` in comparison operation
//   <expr> || <expr> || <expr>
func Or(left ast.Expr, expr ...ast.Expr) ast.Expr {
	if len(expr) == 0 {
		return left
	}
	return Binary(left, And(expr[0], expr[1:]...), token.LOR)
}

// VariableTypeAssert represents variable type assertion expression
//   <varName>.(<t>) e.g. varName.(string)
func VariableTypeAssert(varName string, t ast.Expr) ast.Expr {
	return &ast.TypeAssertExpr{
		X:    ast.NewIdent(varName),
		Type: t,
	}
}

// ExpressionTypeAssert represents expression type assertion
//   <expr>.(<t>) e.g. varName.(string)
func ExpressionTypeAssert(expr, t ast.Expr) ast.Expr {
	return &ast.TypeAssertExpr{
		X:    expr,
		Type: t,
	}
}

// VariableTypeConvert represents variable type conversion expression
//   <t>(<varName>) e.g. string(varName)
func VariableTypeConvert(varName string, t ast.Expr) ast.Expr {
	return Call(
		CallFunctionDescriber{
			FunctionName:                t,
			MinimumNumberOfArguments:    1,
			ExtensibleNumberOfArguments: false,
		},
		ast.NewIdent(varName),
	)
}

// ExpressionTypeConvert represents the expression type conversion expression
//   <t>(<expr>) e.g. string(varName)
func ExpressionTypeConvert(expr ast.Expr, t ast.Expr) ast.Expr {
	return Call(
		CallFunctionDescriber{
			FunctionName:                t,
			MinimumNumberOfArguments:    1,
			ExtensibleNumberOfArguments: false,
		},
		expr,
	)
}

// MakeLenGreatThanZero makes len() > 0 expression
//   len(<arrayName>) > 0
func MakeLenGreatThanZero(arrayName string) ast.Expr {
	return &ast.BinaryExpr{
		X:  Call(LengthFn, ast.NewIdent(arrayName)),
		Op: token.GTR,
		Y:  Zero,
	}
}

func Slice(varName string, lo, hi Expression) ast.Expr {
	return &ast.SliceExpr{
		X:    ast.NewIdent(varName),
		High: safeExpr(hi),
		Low:  safeExpr(lo),
	}
}
