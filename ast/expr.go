package ast

import "github.com/Shri333/golox/scanner"

type expr interface {
	accept(v visitor) interface{}
}

type binary struct {
	left     expr
	operator *scanner.Token
	right    expr
}

func (expr *binary) accept(v visitor) interface{} {
	return v.visitBinaryExpr(expr)
}

type grouping struct {
	expression expr
}

func (expr *grouping) accept(v visitor) interface{} {
	return v.visitGroupingExpr(expr)
}

type literal struct {
	value interface{}
}

func (expr *literal) accept(v visitor) interface{} {
	return v.visitLiteralExpr(expr)
}

type unary struct {
	operator *scanner.Token
	right    expr
}

func (expr *unary) accept(v visitor) interface{} {
	return v.visitUnaryExpr(expr)
}
