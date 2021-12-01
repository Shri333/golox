package ast

type visitor interface {
	visitBinaryExpr(expr *binary) interface{}
	visitGroupingExpr(expr *grouping) interface{}
	visitLiteralExpr(expr *literal) interface{}
	visitUnaryExpr(expr *unary) interface{}
}
