package asthlp

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
	// NewFn is a construction of the `new` function
	NewFn = makeFunc(ast.NewIdent("new"), 1, true)
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
	// StrconvAtoiFn is a construction of the `strconv.Atoi` function
	StrconvAtoiFn = makeFunc(SimpleSelector("strconv", "Atoi"), 1, false)
	// StrconvParseIntFn is a construction of the `strconv.ParseInt` function
	StrconvParseIntFn = makeFunc(SimpleSelector("strconv", "ParseInt"), 3, false)
	// StrconvParseUintFn is a construction of the `strconv.ParseUint` function
	StrconvParseUintFn = makeFunc(SimpleSelector("strconv", "ParseUint"), 3, false)
	// StrconvParseFloatFn is a construction of the `strconv.ParseFloat` function
	StrconvParseFloatFn = makeFunc(SimpleSelector("strconv", "ParseFloat"), 2, false)
	// StrconvParseBoolFn is a construction of the `strconv.ParseBool` function
	StrconvParseBoolFn = makeFunc(SimpleSelector("strconv", "ParseBool"), 1, false)

	// StrconvFormatIntFn is a construction of the `strconv.FormatInt` function
	StrconvFormatIntFn = makeFunc(SimpleSelector("strconv", "FormatInt"), 2, false)
	// StrconvFormatFloatFn is a construction of the `strconv.FormatFloat` function
	StrconvFormatFloatFn = makeFunc(SimpleSelector("strconv", "FormatFloat"), 4, false)
	// StrconvFormatBoolFn is a construction of the `strconv.FormatBool` function
	StrconvFormatBoolFn = makeFunc(SimpleSelector("strconv", "FormatBool"), 1, false)

	// StringsEqualFoldFn is a construction of the `strings.EqualFold` function
	StringsEqualFoldFn = makeFunc(SimpleSelector("strings", "EqualFold"), 2, false)
	// StringsToLowerFn is a construction of the `strings.ToLower` function
	StringsToLowerFn = makeFunc(SimpleSelector("strings", "ToLower"), 1, false)
	// StringsJoinFn is a construction of the `strings.Join` function
	StringsJoinFn = makeFunc(SimpleSelector("strings", "Join"), 2, false)

	// BytesEqualFoldFn is a construction of the `bytes.EqualFold` function
	BytesEqualFoldFn = makeFunc(SimpleSelector("bytes", "EqualFold"), 2, false)
	// BytesEqualFn is a construction of the `bytes.EqualFold` function
	BytesEqualFn = makeFunc(SimpleSelector("bytes", "Equal"), 2, false)
	// BytesNewBufferFn is a construction of the `bytes.NewBuffer` function
	BytesNewBufferFn = makeFunc(SimpleSelector("bytes", "NewBuffer"), 1, false)

	// FmtSprintfFn is a construction of the `fmt.Sprintf` function
	FmtSprintfFn = makeFunc(SimpleSelector("fmt", "Sprintf"), 1, true)
	// FmtFscanfFn is a construction of the `fmt.Fscanf` function
	FmtFscanfFn = makeFunc(SimpleSelector("fmt", "Fscanf"), 1, true)
	// FmtErrorfFn is a construction of the `fmt.Errorf` function
	FmtErrorfFn = makeFunc(SimpleSelector("fmt", "Errorf"), 1, true)

	// JsonUnmarshal is a construction of the `json.Unmarshall` function
	JsonUnmarshal = makeFunc(SimpleSelector("json", "Unmarshal"), 2, false)
	// JsonMarshal is a construction of the `json.Marshall` function
	JsonMarshal = makeFunc(SimpleSelector("json", "Marshal"), 1, false)
	// JsonNewEncoder is a construction of the `json.NewEncoder` function
	JsonNewEncoder = makeFunc(SimpleSelector("json", "NewEncoder"), 1, false)
	// JsonNewDecoder is a construction of the `json.NewDecoder` function
	JsonNewDecoder = makeFunc(SimpleSelector("json", "NewDecoder"), 1, false)

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

	// BytesToIntFn represents utils.BytesToInt function
	BytesToIntFn = makeFunc(SimpleSelector("utils", "BytesToInt"), 1, false)
	// BytesToUintFn represents utils.BytesToUint function
	BytesToUintFn = makeFunc(SimpleSelector("utils", "BytesToUint"), 1, false)
	// BytesToInt64Fn represents utils.BytesToInt64 function
	BytesToInt64Fn = makeFunc(SimpleSelector("utils", "BytesToInt64"), 1, false)
	// BytesToUint64Fn represents utils.BytesToUint64 function
	BytesToUint64Fn = makeFunc(SimpleSelector("utils", "BytesToUint64"), 1, false)
	// BytesToFloat64Fn represents utils.BytesToFloat64 function
	BytesToFloat64Fn = makeFunc(SimpleSelector("utils", "BytesToFloat64"), 1, false)
)

func makeFunc(f ast.Expr, m int, e bool) CallFunctionDescriber {
	return CallFunctionDescriber{
		FunctionName:                f,
		MinimumNumberOfArguments:    m,
		ExtensibleNumberOfArguments: e,
	}
}

func InlineFunc(f ast.Expr) CallFunctionDescriber {
	return CallFunctionDescriber{
		FunctionName:                f,
		MinimumNumberOfArguments:    0,
		ExtensibleNumberOfArguments: true,
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

func CallStmt(x *ast.CallExpr) ast.Stmt {
	return &ast.ExprStmt{X: x}
}
