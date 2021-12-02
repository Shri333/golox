package ast

type Visitor interface {
	VisitBinaryExpr(b *Binary) interface{}
	VisitGroupingExpr(g *Grouping) interface{}
	VisitLiteralExpr(l *Literal) interface{}
	VisitUnaryExpr(u *Unary) interface{}
}
