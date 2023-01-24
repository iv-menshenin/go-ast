package asthlp

import (
	"go/ast"
	"go/token"
)

// Var creates ast.DeclStmt with VAR token, nil values will be excluded from List
func Var(spec ...ast.Spec) ast.Stmt {
	var decl = ast.GenDecl{
		Tok:   token.VAR,
		Specs: make([]ast.Spec, 0, len(spec)),
	}
	for i, s := range spec {
		if s != nil {
			decl.Specs = append(decl.Specs, spec[i])
		}
	}
	return &ast.DeclStmt{
		Decl: &decl,
	}
}

// Return represents return statement
//   return a, b, c, ...
// nil values will be excluded
func Return(results ...ast.Expr) *ast.ReturnStmt {
	var ret = ast.ReturnStmt{
		Results: make([]ast.Expr, 0, len(results)),
	}
	for i, result := range results {
		if result != nil {
			ret.Results = append(ret.Results, results[i])
		}
	}
	return &ret
}

// ReturnEmpty represents empty return statement
//   return
func ReturnEmpty() ast.Stmt {
	return Return()
}

// Block represents block of statement
//   {
//      ... // statements
//   }
// nil values will be excluded from List
func Block(statements ...ast.Stmt) *ast.BlockStmt {
	var block = ast.BlockStmt{
		List: make([]ast.Stmt, 0, len(statements)),
	}
	for i, stmt := range statements {
		if stmt != nil {
			block.List = append(block.List, statements[i])
		}
	}
	return &block
}

// If represents `if` statement
//   if <condition> { <body> }
// nil values will be excluded from Body.List
func If(condition ast.Expr, body ...ast.Stmt) ast.Stmt {
	return &ast.IfStmt{
		If:   1,
		Cond: condition,
		Body: Block(body...),
	}
}

// IfElse represents `if` statement
//   if <condition> { <body> } else { <alternative> }
// nil values will be excluded from Body.List
func IfElse(condition ast.Expr, body *ast.BlockStmt, alternative *ast.BlockStmt) ast.Stmt {
	return &ast.IfStmt{
		If:   1,
		Cond: condition,
		Body: body,
		Else: alternative,
	}
}

// IfInit represents `if` statement with initialization
//   if <init>; <condition> { <body> }
// nil values will be excluded from Body.List
func IfInit(initiation ast.Stmt, condition ast.Expr, body ...ast.Stmt) ast.Stmt {
	return &ast.IfStmt{
		If:   1,
		Init: initiation,
		Cond: condition,
		Body: Block(body...),
	}
}

// Range represents `for` statement with range expression
//   for <key>, <value> := range <x> { <body> }
// ":=" replaced by "=" if define is FALSE
func Range(define bool, key, value string, x ast.Expr, body ...ast.Stmt) ast.Stmt {
	var (
		tok           = token.ASSIGN
		k, v ast.Expr = nil, nil
	)
	if key != "" {
		k = ast.NewIdent(key)
	}
	if value != "" {
		v = ast.NewIdent(value)
	}
	if define {
		tok = token.DEFINE
	}
	return &ast.RangeStmt{
		For:    1,
		Key:    k,
		Value:  v,
		TokPos: 2,
		Tok:    tok,
		X:      x,
		Body:   Block(body...),
	}
}

func EmptyStmt() ast.Stmt {
	return &ast.EmptyStmt{}
}
