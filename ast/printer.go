package ast

import "strconv"

type Printer struct{}

func (p *Printer) Print(e expr) string {
	return e.accept(p).(string)
}

func (p *Printer) visitBinaryExpr(b *binary) interface{} {
	return p.parenthesize(b.operator.Lexeme, b.left, b.right)
}

func (p *Printer) visitGroupingExpr(g *grouping) interface{} {
	return p.parenthesize("group", g.expression)
}

func (p *Printer) visitLiteralExpr(l *literal) interface{} {
	switch value := l.value.(type) {
	case float64:
		return strconv.FormatFloat(value, 'f', -1, 64)
	case bool:
		return strconv.FormatBool(value)
	case nil:
		return "nil"
	default:
		return value
	}
}

func (p *Printer) visitUnaryExpr(u *unary) interface{} {
	return p.parenthesize(u.operator.Lexeme, u.right)
}

func (p *Printer) parenthesize(name string, exprs ...expr) string {
	str := "(" + name

	for _, e := range exprs {
		str += " " + e.accept(p).(string)
	}

	str += ")"
	return str
}
