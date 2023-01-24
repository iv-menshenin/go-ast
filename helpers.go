package asthlp

import "go/ast"

// MakeVarNames creates VarNames inline
func MakeVarNames(vars ...string) VarNames {
	var varNames = make([]ast.Expr, 0, len(vars))
	for _, varName := range vars {
		varNames = append(varNames, ast.NewIdent(varName))
	}
	return varNames
}

// ClearEmptyExpressions returns an []ast.Expr, any nil values will be excluded from this array
func ClearEmptyExpressions(first ast.Expr, next ...ast.Expr) []ast.Expr {
	var result = make([]ast.Expr, 0, len(next)+1)
	if first != nil {
		result = append(result, first)
	}
	for i, expr := range next {
		if expr != nil {
			result = append(result, next[i])
		}
	}
	return result
}

// IfKeyVal returns ast.KeyValueExpr or nil if the `value` attribute is nil. useful with E helper
func IfKeyVal(key string, value ast.Expr) ast.Expr {
	if value == nil {
		return nil
	}
	return &ast.KeyValueExpr{
		Key:   ast.NewIdent(key),
		Value: value,
	}
}

func safeExpr(expression Expression) ast.Expr {
	if expression == nil {
		return nil
	}
	return expression.Expr()
}
