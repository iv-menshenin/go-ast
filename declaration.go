package asthlp

import (
	"go/ast"
	"go/token"
)

type (
	funcDecl struct {
		name *ast.Ident
		comm []*ast.Comment
		recv *ast.Field
		parm *ast.FieldList
		resl *ast.FieldList
		stmt []ast.Stmt
	}
	FuncDecl interface {
		Comments(...string) FuncDecl
		Receiver(*ast.Field) FuncDecl
		Params(...*ast.Field) FuncDecl
		Results(...*ast.Field) FuncDecl
		AppendStmt(...ast.Stmt) FuncDecl
		Decl() ast.Decl
		Lit() ast.Expr
	}
)

func DeclareFunction(name *ast.Ident) FuncDecl {
	return &funcDecl{
		name: name,
	}
}

func (f *funcDecl) Comments(comments ...string) FuncDecl {
	for _, comment := range comments {
		f.comm = append(f.comm, &ast.Comment{Text: comment})
	}
	return f
}

func (f *funcDecl) Receiver(recv *ast.Field) FuncDecl {
	f.recv = recv
	return f
}

func (f *funcDecl) Params(params ...*ast.Field) FuncDecl {
	if f.parm == nil {
		f.parm = &ast.FieldList{}
	}
	f.parm.List = append(f.parm.List, params...)
	return f
}

func (f *funcDecl) Results(results ...*ast.Field) FuncDecl {
	if f.resl == nil {
		f.resl = &ast.FieldList{}
	}
	f.resl.List = append(f.resl.List, results...)
	return f
}

func (f *funcDecl) AppendStmt(stmt ...ast.Stmt) FuncDecl {
	f.stmt = append(f.stmt, stmt...)
	return f
}

func (f *funcDecl) Decl() ast.Decl {
	var recv *ast.FieldList
	if f.recv != nil {
		recv = &ast.FieldList{List: []*ast.Field{f.recv}}
	}
	return &ast.FuncDecl{
		Doc:  &ast.CommentGroup{List: f.comm},
		Recv: recv,
		Name: f.name,
		Type: &ast.FuncType{
			Params:  f.parm,
			Results: f.resl,
		},
		Body: &ast.BlockStmt{List: f.stmt},
	}
}

func (f *funcDecl) Lit() ast.Expr {
	if f.recv != nil {
		panic("can't use a literal on methods (the receiver presents)")
	}
	return &ast.FuncLit{
		Type: &ast.FuncType{
			Params:  f.parm,
			Results: f.resl,
		},
		Body: &ast.BlockStmt{List: f.stmt},
	}
}

type (
	varDecl struct {
		comm []*ast.Comment
		spec []ast.Spec
	}
	VarDecl interface {
		Comments(comments ...string) VarDecl
		AppendSpec(spec ...ast.Spec) VarDecl
		Decl() ast.Decl
		Stmt() ast.Stmt
	}
)

func DeclareVariable() VarDecl {
	return &varDecl{}
}

func (v *varDecl) Comments(comments ...string) VarDecl {
	for _, comment := range comments {
		v.comm = append(v.comm, &ast.Comment{Text: comment})
	}
	return v
}

func (v *varDecl) AppendSpec(spec ...ast.Spec) VarDecl {
	v.spec = append(v.spec, spec...)
	return v
}

func (v *varDecl) Decl() ast.Decl {
	var comm *ast.CommentGroup
	if len(v.comm) > 0 {
		comm = &ast.CommentGroup{List: v.comm}
	}
	return &ast.GenDecl{
		Doc:   comm,
		Tok:   token.VAR,
		Specs: v.spec,
	}
}

func (v *varDecl) Stmt() ast.Stmt {
	return &ast.DeclStmt{Decl: v.Decl()}
}
