package ast

type visitor interface {
	visitBinaryExpr(b *binary) interface{}
	visitGroupingExpr(g *grouping) interface{}
	visitLiteralExpr(l *literal) interface{}
	visitUnaryExpr(u *unary) interface{}
}
