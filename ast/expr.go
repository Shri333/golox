package ast

import "github.com/Shri333/golox/scanner"

type Expr interface {
	accept(v visitor) interface{}
}

type Binary struct {
	Left     Expr
	Operator *scanner.Token
	Right    Expr
}

func (expr *Binary) accept(v visitor) interface{} {
	return v.visitBinaryExpr(expr)
}

type Grouping struct {
	Expr Expr
}

func (expr *Grouping) accept(v visitor) interface{} {
	return v.visitGroupingExpr(expr)
}

type Literal struct {
	Value interface{}
}

func (expr *Literal) accept(v visitor) interface{} {
	return v.visitLiteralExpr(expr)
}

type Unary struct {
	Operator *scanner.Token
	Right    Expr
}

func (expr *Unary) accept(v visitor) interface{} {
	return v.visitUnaryExpr(expr)
}
