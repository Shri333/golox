package interpreter

import (
	"fmt"

	"github.com/Shri333/golox/ast"
	"github.com/Shri333/golox/fault"
	"github.com/Shri333/golox/scanner"
)

type interpreter struct {
	Error bool
}

func NewInterpreter() *interpreter {
	return &interpreter{false}
}

func (i *interpreter) Interpret(expr ast.Expr) {
	value := expr.Accept(i)
	if !i.Error {
		fmt.Println(value)
	}
}

func (i *interpreter) VisitBinaryExpr(b *ast.Binary) interface{} {
	if i.Error {
		return nil
	}

	left := b.Left.Accept(i)
	right := b.Right.Accept(i)

	switch b.Operator.TokenType {
	case scanner.BANG_EQUAL:
		return left != right
	case scanner.EQUAL_EQUAL:
		return left == right
	case scanner.GREATER:
		leftValue, rightValue, err := i.checkNumberOperands(b.Operator, left, right)
		if err != nil {
			return nil
		}

		return leftValue > rightValue
	case scanner.GREATER_EQUAL:
		leftValue, rightValue, err := i.checkNumberOperands(b.Operator, left, right)
		if err != nil {
			return nil
		}

		return leftValue >= rightValue
	case scanner.LESS:
		leftValue, rightValue, err := i.checkNumberOperands(b.Operator, left, right)
		if err != nil {
			return nil
		}

		return leftValue < rightValue
	case scanner.LESS_EQUAL:
		leftValue, rightValue, err := i.checkNumberOperands(b.Operator, left, right)
		if err != nil {
			return nil
		}

		return leftValue <= rightValue
	case scanner.MINUS:
		leftValue, rightValue, err := i.checkNumberOperands(b.Operator, left, right)
		if err != nil {
			return nil
		}

		return leftValue - rightValue
	case scanner.PLUS:
		if leftValue, leftOk := left.(float64); leftOk {
			if rightValue, rightOk := right.(float64); rightOk {
				return leftValue + rightValue
			}
		}

		if leftValue, leftOk := left.(string); leftOk {
			if rightValue, rightOk := right.(string); rightOk {
				return leftValue + rightValue
			}
		}

		i.Error = true
		fault.NewFault(b.Operator.Line, "operands must be two numbers or two strings")
		return nil
	case scanner.SLASH:
		leftValue, rightValue, err := i.checkNumberOperands(b.Operator, left, right)
		if err != nil {
			return nil
		}

		return leftValue / rightValue
	case scanner.STAR:
		leftValue, rightValue, err := i.checkNumberOperands(b.Operator, left, right)
		if err != nil {
			return nil
		}

		return leftValue * rightValue
	}

	return nil
}

func (i *interpreter) VisitGroupingExpr(g *ast.Grouping) interface{} {
	if i.Error {
		return nil
	}

	return g.Expression.Accept(i)
}

func (i *interpreter) VisitLiteralExpr(l *ast.Literal) interface{} {
	if i.Error {
		return nil
	}

	return l.Value
}

func (i *interpreter) VisitUnaryExpr(u *ast.Unary) interface{} {
	if i.Error {
		return nil
	}

	right := u.Right.Accept(i)

	if u.Operator.TokenType == scanner.MINUS {
		if value, ok := right.(float64); ok {
			return -value
		}

		fault.NewFault(u.Operator.Line, "operand must be a number")
		i.Error = true
	}

	if u.Operator.TokenType == scanner.BANG {
		switch value := right.(type) {
		case bool:
			return !value
		case nil:
			return true
		default:
			return false
		}
	}

	return nil
}

func (i *interpreter) checkNumberOperands(operator *scanner.Token, left interface{}, right interface{}) (float64, float64, error) {
	if leftValue, leftOk := left.(float64); leftOk {
		if rightValue, rightOk := right.(float64); rightOk {
			return leftValue, rightValue, nil
		}
	}

	i.Error = true
	return 0.0, 0.0, fault.NewFault(operator.Line, "operands must be numbers")
}
