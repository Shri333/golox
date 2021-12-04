package ast

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
