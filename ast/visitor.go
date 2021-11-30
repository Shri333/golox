package ast

type visitor interface {
	visitBinaryExpr(expr *Binary) interface{}
	visitGroupingExpr(expr *Grouping) interface{}
	visitLiteralExpr(expr *Literal) interface{}
	visitUnaryExpr(expr *Unary) interface{}
}