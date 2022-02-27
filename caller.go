package builders

import (
	"go/ast"
	"go/token"
)

type (
	// CallFunctionDescriber describes a function so that we can do minimal checks
	CallFunctionDescriber struct {
		FunctionName ast.Expr
		// MinimumNumberOfArguments limits the number of arguments, unless indicated that it can expand
		MinimumNumberOfArguments int
		// ExtensibleNumberOfArguments shows that the number of arguments can be increased (notation ...)
		ExtensibleNumberOfArguments bool
	}
)

var (
	// MakeFn is a construction of the `make` function
	MakeFn = makeFunc(ast.NewIdent("make"), 1, true)
	// LengthFn is a construction of the `len` function
	LengthFn = makeFunc(ast.NewIdent("len"), 1, false)
	// CapFn is a construction of the `cap` function
	CapFn = makeFunc(ast.NewIdent("cap"), 1, false)
	// AppendFn is a construction of the `append` function
	AppendFn = makeFunc(ast.NewIdent("append"), 1, true)

	// StrconvItoaFn is a construction of the `strconv.Itoa` function
	StrconvItoaFn = makeFunc(SimpleSelector("strconv", "Itoa"), 1, false)
	// StringsEqualFoldFn is a construction of the `strings.EqualFold` function
	StringsEqualFoldFn = makeFunc(SimpleSelector("strings", "EqualFold"), 2, false)
	// StringsToLowerFn is a construction of the `strings.ToLower` function
	StringsToLowerFn = makeFunc(SimpleSelector("strings", "ToLower"), 1, false)
	// StringsJoinFn is a construction of the `strings.Join` function
	StringsJoinFn = makeFunc(SimpleSelector("strings", "Join"), 2, false)

	// FmtSprintfFn is a construction of the `fmt.Sprintf` function
	FmtSprintfFn = makeFunc(SimpleSelector("fmt", "Sprintf"), 1, true)
	// FmtFscanfFn is a construction of the `fmt.Fscanf` function
	FmtFscanfFn = makeFunc(SimpleSelector("fmt", "Fscanf"), 1, true)

	// JsonUnmarshal is a construction of the `json.Unmarshall` function
	JsonUnmarshal = makeFunc(SimpleSelector("json", "Unmarshal"), 2, false)
	// JsonMarshal is a construction of the `json.Marshall` function
	JsonMarshal = makeFunc(SimpleSelector("json", "Marshal"), 1, false)

	// TimeNowFn is a construction of the `time.Now` function
	TimeNowFn = makeFunc(SimpleSelector("time", "Now"), 0, false)

	// DbQueryFn is a construction of the `db.Query` function
	DbQueryFn = makeFunc(SimpleSelector("db", "Query"), 1, true)
	// RowsNextFn is a construction of the `rows.Next` function
	RowsNextFn = makeFunc(SimpleSelector("rows", "Next"), 0, false)
	// RowsErrFn is a construction of the `rows.Err` function
	RowsErrFn = makeFunc(SimpleSelector("rows", "Err"), 0, false)
	// RowsScanFn is a construction of the `rows.Scan` function
	RowsScanFn = makeFunc(SimpleSelector("rows", "Scan"), 1, true)
)

func makeFunc(f ast.Expr, m int, e bool) CallFunctionDescriber {
	return CallFunctionDescriber{
		FunctionName:                f,
		MinimumNumberOfArguments:    m,
		ExtensibleNumberOfArguments: e,
	}
}

func (c CallFunctionDescriber) checkArgsCount(a int) {
	if c.MinimumNumberOfArguments > a {
		panic("the minimum number of arguments has not been reached")
	}
	if !c.ExtensibleNumberOfArguments && a > c.MinimumNumberOfArguments {
		panic("the maximum number of arguments exceeded")
	}
}

// DeferCall represents a deferred function call statement
func DeferCall(fn CallFunctionDescriber, args ...ast.Expr) ast.Stmt {
	fn.checkArgsCount(len(args))
	return &ast.DeferStmt{
		Call: &ast.CallExpr{
			Fun:  fn.FunctionName,
			Args: args,
		},
	}
}

// Call represents a function call expression
func Call(fn CallFunctionDescriber, args ...ast.Expr) *ast.CallExpr {
	fn.checkArgsCount(len(args))
	return &ast.CallExpr{
		Fun:      fn.FunctionName,
		Args:     args,
		Ellipsis: token.NoPos,
	}
}

// CallEllipsis represents a function call expression with ellipsis after the last argument
func CallEllipsis(fn CallFunctionDescriber, args ...ast.Expr) *ast.CallExpr {
	fn.checkArgsCount(len(args))
	return &ast.CallExpr{
		Fun:      fn.FunctionName,
		Args:     args,
		Ellipsis: 1,
	}
}
