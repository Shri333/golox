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

	if p.match(scanner.FUN) {
		return p.funDeclaration("function")
	}

	return p.statement()
}

func (p *parser) varDeclaration() *VarStmt {
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

func (p *parser) funDeclaration(kind string) *FunStmt {
	if !p.match(scanner.IDENTIFIER) {
		message := fmt.Sprintf("expected %s name", kind)
		panic(fault.NewFault(p.tokens[p.current].Line, message))
	}
	name := p.tokens[p.current-1]

	if !p.match(scanner.LEFT_PAREN) {
		message := fmt.Sprintf("expected '(' after %s name", kind)
		panic(fault.NewFault(p.tokens[p.current].Line, message))
	}

	params := []*scanner.Token{}
	if p.tokens[p.current].TokenType != scanner.RIGHT_PAREN {
		if !p.match(scanner.IDENTIFIER) {
			message := fmt.Sprintf("expected parameter name at %s", p.tokens[p.current].Lexeme)
			panic(fault.NewFault(p.tokens[p.current].Line, message))
		}
		params = append(params, &p.tokens[p.current-1])
		for p.match(scanner.COMMA) {
			if !p.match(scanner.IDENTIFIER) {
				message := fmt.Sprintf("expected parameter name at %s", p.tokens[p.current].Lexeme)
				panic(fault.NewFault(p.tokens[p.current].Line, message))
			}
			params = append(params, &p.tokens[p.current-1])
			if len(params) > 255 {
				panic(fault.NewFault(p.tokens[p.current].Line, "cannot have more than 255 parameters"))
			}
		}
	}

	if !p.match(scanner.RIGHT_PAREN) {
		panic(fault.NewFault(p.tokens[p.current].Line, "expected ')' after parameter list"))
	}

	if !p.match(scanner.LEFT_BRACE) {
		message := fmt.Sprintf("expected '{' before %s body", kind)
		panic(fault.NewFault(p.tokens[p.current].Line, message))
	}

	return &FunStmt{&name, params, p.blockStatement()}
}

func (p *parser) statement() Stmt {
	if p.match(scanner.PRINT) {
		return p.printStatement()
	}

	if p.match(scanner.IF) {
		return p.ifStatement()
	}

	if p.match(scanner.FOR) {
		return p.forStatement()
	}

	if p.match(scanner.WHILE) {
		return p.whileStatement()
	}

	if p.match(scanner.LEFT_BRACE) {
		return p.blockStatement()
	}

	if p.match(scanner.RETURN) {
		return p.returnStatement()
	}

	return p.exprStatement()
}

func (p *parser) printStatement() *PrintStmt {
	expr := p.expression()

	if !p.match(scanner.SEMICOLON) {
		panic(fault.NewFault(p.tokens[p.current].Line, "expected ';' after print statement"))
	}

	return &PrintStmt{expr}
}

func (p *parser) ifStatement() *IfStmt {
	if !p.match(scanner.LEFT_PAREN) {
		panic(fault.NewFault(p.tokens[p.current].Line, "expected '(' after if"))
	}

	condition := p.expression()
	if !p.match(scanner.RIGHT_PAREN) {
		panic(fault.NewFault(p.tokens[p.current].Line, "expected ')' after conditional expression"))
	}

	thenBranch := p.statement()
	var elseBranch Stmt
	if p.match(scanner.ELSE) {
		elseBranch = p.statement()
	}

	return &IfStmt{condition, thenBranch, elseBranch}
}

func (p *parser) forStatement() Stmt {
	if !p.match(scanner.LEFT_PAREN) {
		panic(fault.NewFault(p.tokens[p.current].Line, "expected '(' after for"))
	}

	var initializer Stmt
	if p.match(scanner.SEMICOLON) {
		initializer = nil
	} else if p.match(scanner.VAR) {
		initializer = p.varDeclaration()
	} else {
		initializer = p.exprStatement()
	}

	var condition Expr
	if p.tokens[p.current].TokenType != scanner.SEMICOLON {
		condition = p.expression()
	}
	if !p.match(scanner.SEMICOLON) {
		panic(fault.NewFault(p.tokens[p.current].Line, "expected ';' after conditional expression"))
	}

	var increment Expr
	if p.tokens[p.current].TokenType != scanner.RIGHT_PAREN {
		increment = p.expression()
	}
	if !p.match(scanner.RIGHT_PAREN) {
		panic(fault.NewFault(p.tokens[p.current].Line, "expected ')' after for clause"))
	}

	body := p.statement()
	if increment != nil {
		body = &BlockStmt{[]Stmt{body, &ExprStmt{increment}}}
	}

	if condition == nil {
		condition = &LiteralExpr{true}
	}

	body = &WhileStmt{condition, body}

	if initializer != nil {
		body = &BlockStmt{[]Stmt{initializer, body}}
	}

	return body
}

func (p *parser) whileStatement() *WhileStmt {
	if !p.match(scanner.LEFT_PAREN) {
		panic(fault.NewFault(p.tokens[p.current].Line, "expected '(' after while"))
	}

	condition := p.expression()
	if !p.match(scanner.RIGHT_PAREN) {
		panic(fault.NewFault(p.tokens[p.current].Line, "expected ')' after conditional expression"))
	}

	return &WhileStmt{condition, p.statement()}
}

func (p *parser) blockStatement() *BlockStmt {
	stmts := []Stmt{}

	for p.tokens[p.current].TokenType != scanner.RIGHT_BRACE {
		stmts = append(stmts, p.declaration())
	}

	if !p.match(scanner.RIGHT_BRACE) {
		panic(fault.NewFault(p.tokens[p.current].Line, "expected '}' after block"))
	}

	return &BlockStmt{stmts}
}

func (p *parser) exprStatement() *ExprStmt {
	expr := p.expression()

	if !p.match(scanner.SEMICOLON) {
		panic(fault.NewFault(p.tokens[p.current].Line, "expected ';' after expression statement"))
	}

	return &ExprStmt{expr}
}

func (p *parser) returnStatement() *ReturnStmt {
	keyword := p.tokens[p.current-1]
	var value Expr
	if p.tokens[p.current].TokenType != scanner.SEMICOLON {
		value = p.expression()
	}

	if !p.match(scanner.SEMICOLON) {
		panic(fault.NewFault(p.tokens[p.current].Line, "expected ';' after return statement"))
	}

	return &ReturnStmt{&keyword, value}
}

func (p *parser) expression() Expr {
	return p.assignment()
}

func (p *parser) assignment() Expr {
	expr := p.or()

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

func (p *parser) or() Expr {
	left := p.and()

	for p.match(scanner.OR) {
		operator := p.tokens[p.current-1]
		right := p.and()
		left = &LogicalExpr{left, &operator, right}
	}

	return left
}

func (p *parser) and() Expr {
	left := p.equality()

	for p.match(scanner.AND) {
		operator := p.tokens[p.current-1]
		right := p.equality()
		left = &LogicalExpr{left, &operator, right}
	}

	return left
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

	return p.call()
}

func (p *parser) call() Expr {
	expr := p.primary()

	for p.match(scanner.LEFT_PAREN) {
		args, paren := p.arguments()
		expr = &CallExpr{expr, paren, args}
	}

	return expr
}

func (p *parser) arguments() ([]Expr, scanner.Token) {
	args := []Expr{}
	if p.tokens[p.current].TokenType != scanner.RIGHT_PAREN {
		args = append(args, p.expression())
		for p.match(scanner.COMMA) {
			args = append(args, p.expression())
			if len(args) > 255 {
				panic(fault.NewFault(p.tokens[p.current].Line, "cannot have more than 255 arguments"))
			}
		}
	}

	if !p.match(scanner.RIGHT_PAREN) {
		panic(fault.NewFault(p.tokens[p.current].Line, "expected ')' after argument list"))
	}

	return args, p.tokens[p.current-1]
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
