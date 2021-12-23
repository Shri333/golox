package interpreter

import (
	"fmt"
	"strconv"

	"github.com/Shri333/golox/ast"
	"github.com/Shri333/golox/fault"
	"github.com/Shri333/golox/scanner"
)

type interpreter struct {
	global *environment
	env    *environment
}

func NewInterpreter() *interpreter {
	global := &environment{nil, make(map[string]interface{})}
	global.define("clock", &clock{})
	return &interpreter{global, global}
}

func (i *interpreter) Interpret(stmts []ast.Stmt) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = r.(error)
		}
	}()

	for _, stmt := range stmts {
		stmt.Accept(i)
	}

	return
}

func (i *interpreter) VisitExprStmt(e *ast.ExprStmt) interface{} {
	e.Expression.Accept(i)
	return nil
}

func (i *interpreter) VisitPrintStmt(p *ast.PrintStmt) interface{} {
	value := p.Expression.Accept(i)
	switch v := value.(type) {
	case float64:
		fmt.Println(strconv.FormatFloat(v, 'f', -1, 64))
	case bool:
		fmt.Println(strconv.FormatBool(v))
	case nil:
		fmt.Println("nil")
	default:
		fmt.Println(v)
	}

	return nil
}

func (i *interpreter) VisitVarStmt(v *ast.VarStmt) interface{} {
	var value interface{}
	if v.Initializer != nil {
		value = v.Initializer.Accept(i)
	}

	i.env.define(v.Name.Lexeme, value)
	return nil
}

func (i *interpreter) VisitBlockStmt(b *ast.BlockStmt) interface{} {
	prev := i.env
	defer func() {
		i.env = prev
	}()

	i.env = &environment{prev, make(map[string]interface{})}
	for _, stmt := range b.Statements {
		stmt.Accept(i)
	}

	return nil
}

func (i *interpreter) VisitIfStmt(i_ *ast.IfStmt) interface{} {
	value := i_.Condition.Accept(i)
	if isTruthy(value) {
		i_.ThenBranch.Accept(i)
	} else if i_.ElseBranch != nil {
		i_.ElseBranch.Accept(i)
	}

	return nil
}

func (i *interpreter) VisitWhileStmt(w *ast.WhileStmt) interface{} {
	for isTruthy(w.Condition.Accept(i)) {
		w.Body.Accept(i)
	}

	return nil
}

func (i *interpreter) VisitFunStmt(f *ast.FunStmt) interface{} {
	i.env.define(f.Name.Lexeme, &function{f})
	return nil
}

func (i *interpreter) VisitReturnStmt(v *ast.ReturnStmt) interface{} {
	var value interface{}
	if v.Value != nil {
		value = v.Value.Accept(i)
	}

	panic(value)
}

func (i *interpreter) VisitBinaryExpr(b *ast.BinaryExpr) interface{} {
	left := b.Left.Accept(i)
	right := b.Right.Accept(i)

	switch b.Operator.TokenType {
	case scanner.BANG_EQUAL:
		return left != right
	case scanner.EQUAL_EQUAL:
		return left == right
	case scanner.GREATER:
		leftValue, rightValue := i.checkNumberOperands(b.Operator, left, right)
		return leftValue > rightValue
	case scanner.GREATER_EQUAL:
		leftValue, rightValue := i.checkNumberOperands(b.Operator, left, right)
		return leftValue >= rightValue
	case scanner.LESS:
		leftValue, rightValue := i.checkNumberOperands(b.Operator, left, right)
		return leftValue < rightValue
	case scanner.LESS_EQUAL:
		leftValue, rightValue := i.checkNumberOperands(b.Operator, left, right)
		return leftValue <= rightValue
	case scanner.MINUS:
		leftValue, rightValue := i.checkNumberOperands(b.Operator, left, right)
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

		panic(fault.NewFault(b.Operator.Line, "operands must be two numbers or two strings"))
	case scanner.SLASH:
		leftValue, rightValue := i.checkNumberOperands(b.Operator, left, right)
		return leftValue / rightValue
	case scanner.STAR:
		leftValue, rightValue := i.checkNumberOperands(b.Operator, left, right)
		return leftValue * rightValue
	}

	return nil
}

func (i *interpreter) VisitGroupingExpr(g *ast.GroupingExpr) interface{} {
	return g.Expression.Accept(i)
}

func (i *interpreter) VisitLiteralExpr(l *ast.LiteralExpr) interface{} {
	return l.Value
}

func (i *interpreter) VisitUnaryExpr(u *ast.UnaryExpr) interface{} {
	right := u.Right.Accept(i)

	if u.Operator.TokenType == scanner.MINUS {
		if value, ok := right.(float64); ok {
			return -value
		}

		panic(fault.NewFault(u.Operator.Line, "operand must be a number"))
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

func (i *interpreter) VisitVariableExpr(v *ast.VariableExpr) interface{} {
	return i.env.get(v.Name)
}

func (i *interpreter) VisitAssignExpr(a *ast.AssignExpr) interface{} {
	value := a.Value.Accept(i)
	i.env.assign(a.Name, value)
	return value
}

func (i *interpreter) VisitLogicalExpr(l *ast.LogicalExpr) interface{} {
	left := l.Left.Accept(i)

	if (l.Operator.TokenType == scanner.OR && isTruthy(left)) || !isTruthy(left) {
		return left
	}

	return l.Right.Accept(i)
}

func (i *interpreter) VisitCallExpr(c *ast.CallExpr) interface{} {
	callee := c.Callee.Accept(i)
	args := []interface{}{}
	for _, arg := range c.Arguments {
		args = append(args, arg.Accept(i))
	}

	if f, ok := callee.(callable); ok {
		if len(args) != f.arity() {
			message := fmt.Sprintf("expected %d arguments but got %d.", f.arity(), len(args))
			panic(fault.NewFault(c.Paren.Line, message))
		}

		return f.call(i, args)
	}

	panic(fault.NewFault(c.Paren.Line, "can only call functions and classes"))
}

func (i *interpreter) checkNumberOperands(operator *scanner.Token, left interface{}, right interface{}) (float64, float64) {
	if leftValue, leftOk := left.(float64); leftOk {
		if rightValue, rightOk := right.(float64); rightOk {
			return leftValue, rightValue
		}
	}

	panic(fault.NewFault(operator.Line, "operands must be numbers"))
}

func isTruthy(value interface{}) bool {
	if value == nil {
		return false
	}

	if boolean, ok := value.(bool); ok {
		return boolean
	}

	return true
}
