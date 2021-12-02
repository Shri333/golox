package ast

import "github.com/Shri333/golox/scanner"

type Expr interface {
	Accept(v Visitor) interface{}
}

type Binary struct {
	Left     Expr
	Operator *scanner.Token
	Right    Expr
}

func (b *Binary) Accept(v Visitor) interface{} {
	return v.VisitBinaryExpr(b)
}

type Grouping struct {
	Expression Expr
}

func (g *Grouping) Accept(v Visitor) interface{} {
	return v.VisitGroupingExpr(g)
}

type Literal struct {
	Value interface{}
}

func (l *Literal) Accept(v Visitor) interface{} {
	return v.VisitLiteralExpr(l)
}

type Unary struct {
	Operator *scanner.Token
	Right    Expr
}

func (u *Unary) Accept(v Visitor) interface{} {
	return v.VisitUnaryExpr(u)
}
