package parser

import (
	"fmt"

	"github.com/Shri333/golox/fault"
	"github.com/Shri333/golox/scanner"
)

type Parser struct {
	tokens  []scanner.Token
	current int
	err     error
}

func NewParser(tokens []scanner.Token) *Parser {
	return &Parser{tokens, 0, nil}
}

func (p *Parser) Parse() ([]Stmt, error) {
	stmts := []Stmt{}
	for p.tokens[p.current].TokenType != scanner.EOF {
		stmts = append(stmts, p.declaration())
	}

	return stmts, p.err
}

func (p *Parser) declaration() Stmt {
	defer p.synchronize()

	if p.match(scanner.VAR) {
		return p.varDeclaration()
	}

	if p.match(scanner.FUN) {
		return p.funDeclaration("function")
	}

	if p.match(scanner.CLASS) {
		return p.classDeclaration()
	}

	return p.statement()
}

func (p *Parser) varDeclaration() *VarStmt {
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

func (p *Parser) funDeclaration(kind string) *FunStmt {
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
	if p.tokens[p.current].TokenType != scanner.RIGHT_PAREN && p.tokens[p.current].TokenType != scanner.EOF {
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

func (p *Parser) classDeclaration() *ClassStmt {
	if !p.match(scanner.IDENTIFIER) {
		panic(fault.NewFault(p.tokens[p.current].Line, "expected class name"))
	}
	name := p.tokens[p.current-1]

	var super *VariableExpr
	if p.match(scanner.LESS) {
		if !p.match(scanner.IDENTIFIER) {
			panic(fault.NewFault(p.tokens[p.current].Line, "expected superclass name after '<'"))
		}
		superName := p.tokens[p.current-1]
		super = &VariableExpr{&superName}
	}

	if !p.match(scanner.LEFT_BRACE) {
		panic(fault.NewFault(p.tokens[p.current].Line, "expected '{' before class body"))
	}

	methods := []*FunStmt{}
	for p.tokens[p.current].TokenType != scanner.RIGHT_BRACE && p.tokens[p.current].TokenType != scanner.EOF {
		methods = append(methods, p.funDeclaration("method"))
	}

	if !p.match(scanner.RIGHT_BRACE) {
		panic(fault.NewFault(p.tokens[p.current].Line, "expected '}' after class body"))
	}

	return &ClassStmt{&name, super, methods}
}

func (p *Parser) statement() Stmt {
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

func (p *Parser) printStatement() *PrintStmt {
	expr := p.expression()
	if !p.match(scanner.SEMICOLON) {
		panic(fault.NewFault(p.tokens[p.current].Line, "expected ';' after print statement"))
	}

	return &PrintStmt{expr}
}

func (p *Parser) ifStatement() *IfStmt {
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

func (p *Parser) forStatement() Stmt {
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
	if p.tokens[p.current].TokenType != scanner.SEMICOLON && p.tokens[p.current].TokenType != scanner.EOF {
		condition = p.expression()
	}
	if !p.match(scanner.SEMICOLON) {
		panic(fault.NewFault(p.tokens[p.current].Line, "expected ';' after conditional expression"))
	}

	var increment Expr
	if p.tokens[p.current].TokenType != scanner.RIGHT_PAREN && p.tokens[p.current].TokenType != scanner.EOF {
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

func (p *Parser) whileStatement() *WhileStmt {
	if !p.match(scanner.LEFT_PAREN) {
		panic(fault.NewFault(p.tokens[p.current].Line, "expected '(' after while"))
	}

	condition := p.expression()
	if !p.match(scanner.RIGHT_PAREN) {
		panic(fault.NewFault(p.tokens[p.current].Line, "expected ')' after conditional expression"))
	}

	return &WhileStmt{condition, p.statement()}
}

func (p *Parser) blockStatement() *BlockStmt {
	stmts := []Stmt{}
	for p.tokens[p.current].TokenType != scanner.RIGHT_BRACE && p.tokens[p.current].TokenType != scanner.EOF {
		stmts = append(stmts, p.declaration())
	}

	if !p.match(scanner.RIGHT_BRACE) {
		panic(fault.NewFault(p.tokens[p.current].Line, "expected '}' after block"))
	}

	return &BlockStmt{stmts}
}

func (p *Parser) exprStatement() *ExprStmt {
	expr := p.expression()
	if !p.match(scanner.SEMICOLON) {
		panic(fault.NewFault(p.tokens[p.current].Line, "expected ';' after expression statement"))
	}

	return &ExprStmt{expr}
}

func (p *Parser) returnStatement() *ReturnStmt {
	keyword := p.tokens[p.current-1]
	var value Expr
	if p.tokens[p.current].TokenType != scanner.SEMICOLON && p.tokens[p.current].TokenType != scanner.EOF {
		value = p.expression()
	}

	if !p.match(scanner.SEMICOLON) {
		panic(fault.NewFault(p.tokens[p.current].Line, "expected ';' after return statement"))
	}

	return &ReturnStmt{&keyword, value}
}

func (p *Parser) expression() Expr {
	return p.assignment()
}

func (p *Parser) assignment() Expr {
	expr := p.or()
	if p.match(scanner.EQUAL) {
		equals := p.tokens[p.current-1]
		value := p.assignment()

		if variable, ok := expr.(*VariableExpr); ok {
			return &AssignExpr{variable.Name, value}
		}

		if get, ok := expr.(*GetExpr); ok {
			return &SetExpr{get.Object, get.Name, value}
		}

		fault.NewFault(equals.Line, "invalid assignment target")
	}

	return expr
}

func (p *Parser) or() Expr {
	left := p.and()
	for p.match(scanner.OR) {
		operator := p.tokens[p.current-1]
		right := p.and()
		left = &LogicalExpr{left, &operator, right}
	}

	return left
}

func (p *Parser) and() Expr {
	left := p.equality()
	for p.match(scanner.AND) {
		operator := p.tokens[p.current-1]
		right := p.equality()
		left = &LogicalExpr{left, &operator, right}
	}

	return left
}

func (p *Parser) equality() Expr {
	left := p.comparison()
	for p.match(scanner.BANG_EQUAL, scanner.EQUAL_EQUAL) {
		operator := p.tokens[p.current-1]
		right := p.comparison()
		left = &BinaryExpr{left, &operator, right}
	}

	return left
}

func (p *Parser) comparison() Expr {
	left := p.term()
	for p.match(scanner.GREATER, scanner.GREATER_EQUAL, scanner.LESS, scanner.LESS_EQUAL) {
		operator := p.tokens[p.current-1]
		right := p.term()
		left = &BinaryExpr{left, &operator, right}
	}

	return left
}

func (p *Parser) term() Expr {
	left := p.factor()
	for p.match(scanner.MINUS, scanner.PLUS) {
		operator := p.tokens[p.current-1]
		right := p.factor()
		left = &BinaryExpr{left, &operator, right}
	}

	return left
}

func (p *Parser) factor() Expr {
	left := p.unary()
	for p.match(scanner.SLASH, scanner.STAR) {
		operator := p.tokens[p.current-1]
		right := p.unary()
		left = &BinaryExpr{left, &operator, right}
	}

	return left
}

func (p *Parser) unary() Expr {
	if p.match(scanner.BANG, scanner.MINUS) {
		operator := p.tokens[p.current-1]
		right := p.unary()
		return &UnaryExpr{&operator, right}
	}

	return p.call()
}

func (p *Parser) call() Expr {
	expr := p.primary()
	for {
		if p.match(scanner.LEFT_PAREN) {
			args, paren := p.arguments()
			expr = &CallExpr{expr, paren, args}
		} else if p.match(scanner.DOT) {
			if !p.match(scanner.IDENTIFIER) {
				panic(fault.NewFault(p.tokens[p.current].Line, "expected property name after '.'"))
			}
			name := p.tokens[p.current-1]
			expr = &GetExpr{expr, &name}
		} else {
			break
		}
	}

	return expr
}

func (p *Parser) arguments() ([]Expr, scanner.Token) {
	args := []Expr{}
	if p.tokens[p.current].TokenType != scanner.RIGHT_PAREN && p.tokens[p.current].TokenType != scanner.EOF {
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

func (p *Parser) primary() Expr {
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

	if p.match(scanner.THIS) {
		previous := &p.tokens[p.current-1]
		return &ThisExpr{previous}
	}

	if p.match(scanner.SUPER) {
		keyword := p.tokens[p.current-1]
		if !p.match(scanner.DOT) || !p.match(scanner.IDENTIFIER) {
			panic(fault.NewFault(p.tokens[p.current].Line, "expected property access after 'super'"))
		}
		method := p.tokens[p.current-1]
		return &SuperExpr{&keyword, &method}
	}

	if p.match(scanner.LEFT_PAREN) {
		e := p.expression()
		if !p.match(scanner.RIGHT_PAREN) {
			message := fmt.Sprintf("expected ')' after '%s'", p.tokens[p.current-1].Lexeme)
			panic(fault.NewFault(p.tokens[p.current].Line, message))
		}
		return &GroupingExpr{e}
	}

	message := fmt.Sprintf("expected expression at '%s'", p.tokens[p.current].Lexeme)
	panic(fault.NewFault(p.tokens[p.current].Line, message))
}

func (p *Parser) match(types ...int) bool {
	currentType := p.tokens[p.current].TokenType
	if currentType == scanner.EOF {
		return false
	}

	for _, tokenType := range types {
		if currentType == tokenType {
			p.current++
			return true
		}
	}

	return false
}

func (p *Parser) synchronize() {
	if r := recover(); r != nil {
		defer func() { p.err = r.(error) }()

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
