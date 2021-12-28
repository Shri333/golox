package interpreter

import (
	"fmt"
	"strconv"

	"github.com/Shri333/golox/fault"
	"github.com/Shri333/golox/parser"
	"github.com/Shri333/golox/scanner"
)

type Interpreter struct {
	global  *environment
	current *environment
	locals  map[parser.Expr]int
}

func NewInterpreter() *Interpreter {
	global := &environment{nil, make(map[string]interface{})}
	global.define("clock", &clock{})
	return &Interpreter{global, global, make(map[parser.Expr]int)}
}

func (i *Interpreter) Interpret(stmts []parser.Stmt) (err error) {
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

func (i *Interpreter) Resolve(expr parser.Expr, depth int) {
	i.locals[expr] = depth
}

func (i *Interpreter) VisitExprStmt(e *parser.ExprStmt) interface{} {
	e.Expression.Accept(i)
	return nil
}

func (i *Interpreter) VisitPrintStmt(p *parser.PrintStmt) interface{} {
	value := p.Expression.Accept(i)
	switch v := value.(type) {
	case float64:
		fmt.Println(strconv.FormatFloat(v, 'f', -1, 64))
	case bool:
		fmt.Println(strconv.FormatBool(v))
	default:
		fmt.Println(v)
	}

	return nil
}

func (i *Interpreter) VisitVarStmt(v *parser.VarStmt) interface{} {
	var value interface{}
	if v.Initializer != nil {
		value = v.Initializer.Accept(i)
	}

	i.current.define(v.Name.Lexeme, value)
	return nil
}

func (i *Interpreter) VisitBlockStmt(b *parser.BlockStmt) interface{} {
	prev := i.current
	defer func() { i.current = prev }()

	i.current = &environment{prev, make(map[string]interface{})}
	for _, stmt := range b.Statements {
		stmt.Accept(i)
	}

	return nil
}

func (i *Interpreter) VisitIfStmt(i_ *parser.IfStmt) interface{} {
	value := i_.Condition.Accept(i)
	if isTruthy(value) {
		i_.ThenBranch.Accept(i)
	} else if i_.ElseBranch != nil {
		i_.ElseBranch.Accept(i)
	}

	return nil
}

func (i *Interpreter) VisitWhileStmt(w *parser.WhileStmt) interface{} {
	for isTruthy(w.Condition.Accept(i)) {
		w.Body.Accept(i)
	}

	return nil
}

func (i *Interpreter) VisitFunStmt(f *parser.FunStmt) interface{} {
	fn := &function{f, i.current, false}
	i.current.define(f.Name.Lexeme, fn)
	return nil
}

func (i *Interpreter) VisitReturnStmt(v *parser.ReturnStmt) interface{} {
	var value interface{}
	if v.Value != nil {
		value = v.Value.Accept(i)
	}

	panic(value)
}

func (i *Interpreter) VisitClassStmt(c *parser.ClassStmt) interface{} {
	i.current.define(c.Name.Lexeme, nil)
	methods := make(map[string]*function)
	for _, method := range c.Methods {
		if method.Name.Lexeme == "init" {
			methods[method.Name.Lexeme] = &function{method, i.current, true}
		} else {
			methods[method.Name.Lexeme] = &function{method, i.current, false}
		}
	}

	i.current.assign(c.Name, &class{c.Name.Lexeme, methods})
	return nil
}

func (i *Interpreter) VisitBinaryExpr(b *parser.BinaryExpr) interface{} {
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

func (i *Interpreter) VisitGroupingExpr(g *parser.GroupingExpr) interface{} {
	return g.Expression.Accept(i)
}

func (i *Interpreter) VisitLiteralExpr(l *parser.LiteralExpr) interface{} {
	return l.Value
}

func (i *Interpreter) VisitUnaryExpr(u *parser.UnaryExpr) interface{} {
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

func (i *Interpreter) VisitVariableExpr(v *parser.VariableExpr) interface{} {
	if dist, ok := i.locals[v]; ok {
		return i.current.getAt(v.Name.Lexeme, dist)
	}

	return i.global.get(v.Name)
}

func (i *Interpreter) VisitAssignExpr(a *parser.AssignExpr) interface{} {
	value := a.Value.Accept(i)
	if dist, ok := i.locals[a]; ok {
		i.current.assignAt(a.Name.Lexeme, value, dist)
	} else {
		i.global.assign(a.Name, value)
	}

	return value
}

func (i *Interpreter) VisitLogicalExpr(l *parser.LogicalExpr) interface{} {
	left := l.Left.Accept(i)
	if (l.Operator.TokenType == scanner.OR && isTruthy(left)) || !isTruthy(left) {
		return left
	}

	return l.Right.Accept(i)
}

func (i *Interpreter) VisitCallExpr(c *parser.CallExpr) interface{} {
	callee := c.Callee.Accept(i)
	args := []interface{}{}
	for _, arg := range c.Arguments {
		args = append(args, arg.Accept(i))
	}

	if f, ok := callee.(callable); ok {
		if len(args) != f.arity() {
			message := fmt.Sprintf("expected %d arguments but got %d", f.arity(), len(args))
			panic(fault.NewFault(c.Paren.Line, message))
		}

		return f.call(i, args)
	}

	panic(fault.NewFault(c.Paren.Line, "can only call functions and classes"))
}

func (i *Interpreter) VisitGetExpr(g *parser.GetExpr) interface{} {
	object := g.Object.Accept(i)
	if o, ok := object.(*instance); ok {
		return o.get(g.Name)
	}

	panic(fault.NewFault(g.Name.Line, "only instances have properties"))
}

func (i *Interpreter) VisitSetExpr(s *parser.SetExpr) interface{} {
	object := s.Object.Accept(i)
	if o, ok := object.(*instance); ok {
		value := s.Value.Accept(i)
		o.set(s.Name, value)
		return value
	}

	panic(fault.NewFault(s.Name.Line, "only instances have fields"))
}

func (i *Interpreter) VisitThisExpr(t *parser.ThisExpr) interface{} {
	if dist, ok := i.locals[t]; ok {
		return i.current.getAt(t.Keyword.Lexeme, dist)
	}

	return i.global.get(t.Keyword)
}

func (i *Interpreter) checkNumberOperands(operator *scanner.Token, left interface{}, right interface{}) (float64, float64) {
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
