package ast

import (
	"fmt"

	"github.com/Shri333/golox/fault"
	"github.com/Shri333/golox/scanner"
)

type parser struct {
	tokens  []scanner.Token
	current int
	err     error
}

func NewParser(tokens []scanner.Token) *parser {
	return &parser{tokens, 0, nil}
}

func (p *parser) Parse() ([]Stmt, error) {
	stmts := []Stmt{}

	for p.tokens[p.current].TokenType != scanner.EOF {
		stmts = append(stmts, p.declaration())
	}

	return stmts, p.err
}

func (p *parser) declaration() Stmt {
	defer p.synchronize()

	if p.match(scanner.VAR) {
		return p.varDeclaration()
	}

	return p.statement()
}

func (p *parser) varDeclaration() Stmt {
	if !p.match(scanner.IDENTIFIER) {
		panic(fault.NewFault(p.tokens[p.current].Line, "expected variable name"))
	}

	name := p.tokens[p.current-1]
	var initializer Expr
	if p.match(scanner.EQUAL) {
		initializer = p.expression()
	}

	if !p.match(scanner.SEMICOLON) {
		panic(fault.NewFault(p.tokens[p.current].Line, "expected ';' after variable declaration"))
	}

	return &VarStmt{&name, initializer}
}

func (p *parser) statement() Stmt {
	if p.match(scanner.PRINT) {
		return p.printStatement()
	}

	if p.match(scanner.LEFT_BRACE) {
		return &BlockStmt{p.block()}
	}

	return p.exprStatement()
}

func (p *parser) printStatement() Stmt {
	expr := p.expression()

	if !p.match(scanner.SEMICOLON) {
		panic(fault.NewFault(p.tokens[p.current].Line, "expected ';' after print statement"))
	}

	return &PrintStmt{expr}
}

func (p *parser) exprStatement() Stmt {
	expr := p.expression()

	if !p.match(scanner.SEMICOLON) {
		panic(fault.NewFault(p.tokens[p.current].Line, "expected ';' after expression statement"))
	}

	return &ExprStmt{expr}
}

func (p *parser) block() []Stmt {
	stmts := []Stmt{}

	for p.tokens[p.current].TokenType != scanner.EOF && p.tokens[p.current].TokenType != scanner.RIGHT_BRACE {
		stmts = append(stmts, p.declaration())
	}

	if !p.match(scanner.RIGHT_BRACE) {
		panic(fault.NewFault(p.tokens[p.current].Line, "expected '}' after block"))
	}

	return stmts
}

func (p *parser) expression() Expr {
	return p.assignment()
}

func (p *parser) assignment() Expr {
	expr := p.equality()

	if p.match(scanner.EQUAL) {
		equals := p.tokens[p.current-1]
		value := p.assignment()

		if variable, ok := expr.(*VariableExpr); ok {
			return &AssignExpr{variable.Name, value}
		}

		fault.NewFault(equals.Line, "invalid assignment target")
	}

	return expr
}

func (p *parser) equality() Expr {
	left := p.comparison()

	for p.match(scanner.BANG_EQUAL, scanner.EQUAL_EQUAL) {
		operator := p.tokens[p.current-1]
		right := p.comparison()
		left = &BinaryExpr{left, &operator, right}
	}

	return left
}

func (p *parser) comparison() Expr {
	left := p.term()

	for p.match(scanner.GREATER, scanner.GREATER_EQUAL, scanner.LESS, scanner.LESS_EQUAL) {
		operator := p.tokens[p.current-1]
		right := p.term()
		left = &BinaryExpr{left, &operator, right}
	}

	return left
}

func (p *parser) term() Expr {
	left := p.factor()

	for p.match(scanner.MINUS, scanner.PLUS) {
		operator := p.tokens[p.current-1]
		right := p.factor()
		left = &BinaryExpr{left, &operator, right}
	}

	return left
}

func (p *parser) factor() Expr {
	left := p.unary()

	for p.match(scanner.SLASH, scanner.STAR) {
		operator := p.tokens[p.current-1]
		right := p.unary()
		left = &BinaryExpr{left, &operator, right}
	}

	return left
}

func (p *parser) unary() Expr {
	if p.match(scanner.BANG, scanner.MINUS) {
		operator := p.tokens[p.current-1]
		right := p.unary()
		return &UnaryExpr{&operator, right}
	}

	return p.primary()
}

func (p *parser) primary() Expr {
	if p.match(scanner.FALSE) {
		return &LiteralExpr{false}
	}

	if p.match(scanner.TRUE) {
		return &LiteralExpr{true}
	}

	if p.match(scanner.NIL) {
		return &LiteralExpr{nil}
	}

	if p.match(scanner.NUMBER, scanner.STRING) {
		value := p.tokens[p.current-1].Literal
		return &LiteralExpr{value}
	}

	if p.match(scanner.IDENTIFIER) {
		previous := &p.tokens[p.current-1]
		return &VariableExpr{previous}
	}

	if p.match(scanner.LEFT_PAREN) {
		e := p.expression()
		if p.tokens[p.current].TokenType != scanner.RIGHT_PAREN {
			message := fmt.Sprintf("expected ')' after \"%s\"", p.tokens[p.current].Lexeme)
			panic(fault.NewFault(p.tokens[p.current].Line, message))
		}
		p.current++
		return &GroupingExpr{e}
	}

	message := fmt.Sprintf("expected expression at \"%s\"", p.tokens[p.current].Lexeme)
	panic(fault.NewFault(p.tokens[p.current].Line, message))
}

func (p *parser) match(types ...int) bool {
	if p.tokens[p.current].TokenType == scanner.EOF {
		return false
	}

	actualType := p.tokens[p.current].TokenType
	for _, tokenType := range types {
		if actualType == tokenType {
			p.current++
			return true
		}
	}

	return false
}

func (p *parser) synchronize() {
	if r := recover(); r != nil {
		defer func() {
			p.err = r.(error)
		}()

		if p.tokens[p.current].TokenType != scanner.EOF {
			p.current++
		}

		for p.tokens[p.current].TokenType != scanner.EOF {
			if p.tokens[p.current-1].TokenType == scanner.SEMICOLON {
				return
			}

			switch p.tokens[p.current].TokenType {
			case scanner.CLASS:
				return
			case scanner.FUN:
				return
			case scanner.VAR:
				return
			case scanner.FOR:
				return
			case scanner.IF:
				return
			case scanner.WHILE:
				return
			case scanner.PRINT:
				return
			case scanner.RETURN:
				return
			}

			p.current++
		}
	}
}
