package ast

import (
	"fmt"

	"github.com/Shri333/golox/fault"
	"github.com/Shri333/golox/scanner"
)

type parser struct {
	tokens  []scanner.Token
	current int
	panic   bool
	Error   bool
}

func NewParser(tokens []scanner.Token) *parser {
	return &parser{tokens, 0, false, false}
}

func (p *parser) Parse() expr {
	root := p.expression()
	if p.panic {
		return nil
	}

	return root
}

func (p *parser) expression() expr {
	if p.panic {
		p.synchronize()
		return nil
	}

	return p.equality()
}

func (p *parser) equality() expr {
	if p.panic {
		p.synchronize()
		return nil
	}

	left := p.comparison()

	for p.match(scanner.BANG_EQUAL, scanner.EQUAL_EQUAL) {
		operator := p.tokens[p.current-1]
		right := p.comparison()
		left = &binary{left, &operator, right}
	}

	return left
}

func (p *parser) comparison() expr {
	if p.panic {
		p.synchronize()
		return nil
	}

	left := p.term()

	for p.match(scanner.GREATER, scanner.GREATER_EQUAL, scanner.LESS, scanner.LESS_EQUAL) {
		operator := p.tokens[p.current-1]
		right := p.term()
		left = &binary{left, &operator, right}
	}

	return left
}

func (p *parser) term() expr {
	if p.panic {
		p.synchronize()
		return nil
	}

	left := p.factor()

	for p.match(scanner.MINUS, scanner.PLUS) {
		operator := p.tokens[p.current-1]
		right := p.factor()
		left = &binary{left, &operator, right}
	}

	return left
}

func (p *parser) factor() expr {
	if p.panic {
		p.synchronize()
		return nil
	}

	left := p.unary()

	for p.match(scanner.SLASH, scanner.STAR) {
		operator := p.tokens[p.current-1]
		right := p.unary()
		left = &binary{left, &operator, right}
	}

	return left
}

func (p *parser) unary() expr {
	if p.panic {
		p.synchronize()
		return nil
	}

	if p.match(scanner.BANG, scanner.MINUS) {
		operator := p.tokens[p.current-1]
		right := p.unary()
		return &unary{&operator, right}
	}

	return p.primary()
}

func (p *parser) primary() expr {
	if p.panic {
		p.synchronize()
		return nil
	}

	if p.match(scanner.FALSE) {
		return &literal{false}
	}

	if p.match(scanner.TRUE) {
		return &literal{true}
	}

	if p.match(scanner.NIL) {
		return &literal{nil}
	}

	if p.match(scanner.NUMBER, scanner.STRING) {
		value := p.tokens[p.current-1].Literal
		return &literal{value}
	}

	if p.match(scanner.LEFT_PAREN) {
		e := p.expression()
		if p.tokens[p.current].TokenType != scanner.RIGHT_PAREN {
			fault.NewFault(p.tokens[p.current].Line, "expected ')' after expression")
			p.panic = true
			p.Error = true
			return nil
		}
		p.current++
		return &grouping{e}
	}

	message := fmt.Sprintf("expected expression at \"%s\"", p.tokens[p.current].Lexeme)
	fault.NewFault(p.tokens[p.current].Line, message)
	p.panic = true
	p.Error = true
	return nil
}

func (p *parser) match(types ...int) bool {
	if p.current == len(p.tokens) {
		return false
	}

	actualType := p.tokens[p.current].TokenType
	for _, tokenType := range types {
		if actualType == scanner.EOF {
			return false
		}
		if actualType == tokenType {
			p.current++
			return true
		}
	}

	return false
}

func (p *parser) synchronize() {
	p.current++
	p.panic = false

	for ; p.current < len(p.tokens); p.current++ {
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
	}
}
