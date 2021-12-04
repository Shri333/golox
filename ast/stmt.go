package ast

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
