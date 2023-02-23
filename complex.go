package asthlp

import (
	"fmt"
	"go/ast"
	"go/token"
	"sort"
	"strings"
)

// MakeTagsForField with tags like map[tag]values, string `tag1:"values1" tag2:"values2"` is created
func MakeTagsForField(tags map[string][]string) *ast.BasicLit {
	if len(tags) == 0 {
		return nil
	}
	arrTags := make([]string, 0, len(tags))
	for key, val := range tags {
		if len(val) > 0 {
			arrTags = append(arrTags, fmt.Sprintf("%s:\"%s\"", key, strings.Join(val, ",")))
		}
	}
	sort.Strings(arrTags)
	return &ast.BasicLit{
		ValuePos: 1,
		Kind:     token.STRING,
		Value:    "`" + strings.Join(arrTags, " ") + "`",
	}
}

// MakeCallWithErrChecking creates a function call statement with error checking branch
//
//	if <varName>, err = callExpr(); err != nil {
//	    <body>
//	}
//
// varName can be omitted
func MakeCallWithErrChecking(varName string, callExpr *ast.CallExpr, body ...ast.Stmt) ast.Stmt {
	if len(body) == 0 {
		body = []ast.Stmt{ReturnEmpty()}
	}
	if varName != "" {
		return IfInit(
			Assign(MakeVarNames(varName, "err"), Assignment, callExpr),
			NotEqual(ast.NewIdent("err"), Nil),
			body...,
		)
	} else {
		return IfInit(
			Assign(MakeVarNames("err"), Assignment, callExpr),
			NotEqual(ast.NewIdent("err"), Nil),
			body...,
		)
	}
}

// MakeCallReturnIfError creates a function call statement with error checking branch contained `return err`
//
//	if <varName>, err = callExpr(); err != nil {
//	    return err
//	}
//
// varName can be omitted
func MakeCallReturnIfError(varName ast.Expr, callExpr *ast.CallExpr) ast.Stmt {
	var errVar = ast.NewIdent("err")
	if varName != nil {
		return IfInit(
			Assign(VarNames{varName, errVar}, Assignment, callExpr),
			NotEqual(errVar, Nil),
			Return(errVar),
		)
	} else {
		return IfInit(
			Assign(VarNames{errVar}, Assignment, callExpr),
			NotEqual(errVar, Nil),
			Return(errVar),
		)
	}
}

func MakeTypeSwitch(assign ast.Stmt, cases ...SwitchCase) ast.Stmt {
	return &ast.TypeSwitchStmt{
		Assign: assign,
		Body:   &ast.BlockStmt{List: casesToStatements(cases)},
	}
}

func MakeSwitch(init ast.Stmt, tag ast.Expr, cases ...SwitchCase) ast.Stmt {
	return &ast.SwitchStmt{
		Init: init,
		Tag:  tag,
		Body: &ast.BlockStmt{List: casesToStatements(cases)},
	}
}

type (
	SwitchCase struct {
		clause []ast.Expr
		body   []ast.Stmt
	}
)

func MakeSwitchCase(clause ...ast.Expr) SwitchCase {
	return SwitchCase{
		clause: clause,
	}
}

func (c SwitchCase) Body(statements ...ast.Stmt) SwitchCase {
	c.body = statements
	return c
}

func casesToStatements(cases []SwitchCase) []ast.Stmt {
	var result = make([]ast.Stmt, 0, len(cases))
	for _, oneCase := range cases {
		result = append(result, &ast.CaseClause{
			List: oneCase.clause,
			Body: oneCase.body,
		})
	}
	return result
}
