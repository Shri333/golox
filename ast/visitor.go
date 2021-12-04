package ast

type ExprVisitor interface {
	VisitBinaryExpr(b *BinaryExpr) interface{}
	VisitGroupingExpr(g *GroupingExpr) interface{}
	VisitLiteralExpr(l *LiteralExpr) interface{}
	VisitUnaryExpr(u *UnaryExpr) interface{}
}

type StmtVisitor interface {
	VisitExprStmt(e *ExprStmt) interface{}
	VisitPrintStmt(p *PrintStmt) interface{}
}
