package parser

import "github.com/Shri333/golox/scanner"

type Stmt interface {
	Accept(v StmtVisitor) interface{}
}

type ExprStmt struct {
	Expression Expr
}

func (e *ExprStmt) Accept(v StmtVisitor) interface{} {
	return v.VisitExprStmt(e)
}

type PrintStmt struct {
	Expression Expr
}

func (p *PrintStmt) Accept(v StmtVisitor) interface{} {
	return v.VisitPrintStmt(p)
}

type VarStmt struct {
	Name        *scanner.Token
	Initializer Expr
}

func (v *VarStmt) Accept(v_ StmtVisitor) interface{} {
	return v_.VisitVarStmt(v)
}

type BlockStmt struct {
	Statements []Stmt
}

func (b *BlockStmt) Accept(v StmtVisitor) interface{} {
	return v.VisitBlockStmt(b)
}

type IfStmt struct {
	Condition  Expr
	ThenBranch Stmt
	ElseBranch Stmt
}

func (i *IfStmt) Accept(v StmtVisitor) interface{} {
	return v.VisitIfStmt(i)
}

type WhileStmt struct {
	Condition Expr
	Body      Stmt
}

func (w *WhileStmt) Accept(v StmtVisitor) interface{} {
	return v.VisitWhileStmt(w)
}

type FunStmt struct {
	Name   *scanner.Token
	Params []*scanner.Token
	Body   *BlockStmt
}

func (f *FunStmt) Accept(v StmtVisitor) interface{} {
	return v.VisitFunStmt(f)
}

type ReturnStmt struct {
	Keyword *scanner.Token
	Value   Expr
}

func (r *ReturnStmt) Accept(v StmtVisitor) interface{} {
	return v.VisitReturnStmt(r)
}
