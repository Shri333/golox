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

func (b *binary) accept(v visitor) interface{} {
	return v.visitBinaryExpr(b)
}

type grouping struct {
	expression expr
}

func (g *grouping) accept(v visitor) interface{} {
	return v.visitGroupingExpr(g)
}

type literal struct {
	value interface{}
}

func (l *literal) accept(v visitor) interface{} {
	return v.visitLiteralExpr(l)
}

type unary struct {
	operator *scanner.Token
	right    expr
}

func (u *unary) accept(v visitor) interface{} {
	return v.visitUnaryExpr(u)
}
