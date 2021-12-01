package ast

import "strconv"

type Printer struct{}

func (p *Printer) Print(e expr) string {
	return e.accept(p).(string)
}

func (p *Printer) visitBinaryExpr(expr *binary) interface{} {
	return p.parenthesize(expr.operator.Lexeme, expr.left, expr.right)
}

func (p *Printer) visitGroupingExpr(expr *grouping) interface{} {
	return p.parenthesize("group", expr.expression)
}

func (p *Printer) visitLiteralExpr(expr *literal) interface{} {
	switch value := expr.value.(type) {
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

func (p *Printer) visitUnaryExpr(expr *unary) interface{} {
	return p.parenthesize(expr.operator.Lexeme, expr.right)
}

func (p *Printer) parenthesize(name string, exprs ...expr) string {
	str := "(" + name

	for _, e := range exprs {
		str += " " + e.accept(p).(string)
	}

	str += ")"
	return str
}
