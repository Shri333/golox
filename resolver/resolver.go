package resolver

import (
	"github.com/Shri333/golox/fault"
	"github.com/Shri333/golox/interpreter"
	"github.com/Shri333/golox/parser"
	"github.com/Shri333/golox/scanner"
)

type Resolver struct {
	i      *interpreter.Interpreter
	scopes []map[string]bool
	ftype  int
	ctype  int
}

func NewResolver(i *interpreter.Interpreter) *Resolver {
	return &Resolver{i, []map[string]bool{}, 0, 0}
}

func (r *Resolver) Resolve(stmts []parser.Stmt) (err error) {
	defer func() {
		if r_ := recover(); r_ != nil {
			err = r_.(error)
		}
	}()

	for _, stmt := range stmts {
		stmt.Accept(r)
	}

	return err
}

func (r *Resolver) VisitExprStmt(e *parser.ExprStmt) interface{} {
	e.Expression.Accept(r)
	return nil
}

func (r *Resolver) VisitPrintStmt(p *parser.PrintStmt) interface{} {
	p.Expression.Accept(r)
	return nil
}

func (r *Resolver) VisitVarStmt(v *parser.VarStmt) interface{} {
	r.declare(v.Name)
	if v.Initializer != nil {
		v.Initializer.Accept(r)
	}
	r.define(v.Name)

	return nil
}

func (r *Resolver) VisitBlockStmt(b *parser.BlockStmt) interface{} {
	r.scopes = append(r.scopes, make(map[string]bool))
	for _, stmt := range b.Statements {
		stmt.Accept(r)
	}
	r.scopes = r.scopes[:len(r.scopes)-1]

	return nil
}

func (r *Resolver) VisitIfStmt(i *parser.IfStmt) interface{} {
	i.Condition.Accept(r)
	i.ThenBranch.Accept(r)
	if i.ElseBranch != nil {
		i.ElseBranch.Accept(r)
	}

	return nil
}

func (r *Resolver) VisitWhileStmt(w *parser.WhileStmt) interface{} {
	w.Condition.Accept(r)
	w.Body.Accept(r)
	return nil
}

func (r *Resolver) VisitFunStmt(f *parser.FunStmt) interface{} {
	r.declare(f.Name)
	r.define(f.Name)
	r.resolveFunction(f, 1)

	return nil
}

func (r *Resolver) VisitReturnStmt(r_ *parser.ReturnStmt) interface{} {
	if r.ftype == 0 {
		panic(fault.NewFault(r_.Keyword.Line, "cannot return outside of a function"))
	}

	if r_.Value != nil {
		if r.ftype == 3 {
			panic(fault.NewFault(r_.Keyword.Line, "cannot return a value from an initializer"))
		}

		r_.Value.Accept(r)
	}

	return nil
}

func (r *Resolver) VisitClassStmt(c *parser.ClassStmt) interface{} {
	enclosing := r.ctype
	r.ctype = 1

	r.declare(c.Name)
	r.define(c.Name)

	r.scopes = append(r.scopes, make(map[string]bool))
	scope := r.scopes[len(r.scopes)-1]
	scope["this"] = true

	for _, method := range c.Methods {
		if method.Name.Lexeme == "init" {
			r.resolveFunction(method, 3)
		} else {
			r.resolveFunction(method, 2)
		}
	}

	r.scopes = r.scopes[:len(r.scopes)-1]
	r.ctype = enclosing
	return nil
}

func (r *Resolver) VisitBinaryExpr(b *parser.BinaryExpr) interface{} {
	b.Left.Accept(r)
	b.Right.Accept(r)
	return nil
}

func (r *Resolver) VisitGroupingExpr(g *parser.GroupingExpr) interface{} {
	g.Expression.Accept(r)
	return nil
}

func (r *Resolver) VisitLiteralExpr(l *parser.LiteralExpr) interface{} {
	return nil
}

func (r *Resolver) VisitUnaryExpr(u *parser.UnaryExpr) interface{} {
	u.Right.Accept(r)
	return nil
}

func (r *Resolver) VisitVariableExpr(v *parser.VariableExpr) interface{} {
	if len(r.scopes) > 0 {
		scope := r.scopes[len(r.scopes)-1]
		if value, ok := scope[v.Name.Lexeme]; ok && !value {
			panic(fault.NewFault(v.Name.Line, "cannot read local variable in its own initializer"))
		}
	}

	r.resolveLocal(v, v.Name)
	return nil
}

func (r *Resolver) VisitAssignExpr(a *parser.AssignExpr) interface{} {
	a.Value.Accept(r)
	r.resolveLocal(a, a.Name)
	return nil
}

func (r *Resolver) VisitLogicalExpr(l *parser.LogicalExpr) interface{} {
	l.Left.Accept(r)
	l.Right.Accept(r)
	return nil
}

func (r *Resolver) VisitCallExpr(c *parser.CallExpr) interface{} {
	c.Callee.Accept(r)
	for _, arg := range c.Arguments {
		arg.Accept(r)
	}

	return nil
}

func (r *Resolver) VisitGetExpr(g *parser.GetExpr) interface{} {
	g.Object.Accept(r)
	return nil
}

func (r *Resolver) VisitSetExpr(s *parser.SetExpr) interface{} {
	s.Value.Accept(r)
	s.Object.Accept(r)
	return nil
}

func (r *Resolver) VisitThisExpr(t *parser.ThisExpr) interface{} {
	if r.ctype == 0 {
		panic(fault.NewFault(t.Keyword.Line, "cannot use 'this' outside of a class"))
	}

	r.resolveLocal(t, t.Keyword)
	return nil
}

func (r *Resolver) declare(name *scanner.Token) {
	if len(r.scopes) > 0 {
		scope := r.scopes[len(r.scopes)-1]
		if _, ok := scope[name.Lexeme]; ok {
			panic(fault.NewFault(name.Line, "variable cannot be redeclared in local scope"))
		}
		scope[name.Lexeme] = false
	}
}

func (r *Resolver) define(name *scanner.Token) {
	if len(r.scopes) > 0 {
		scope := r.scopes[len(r.scopes)-1]
		scope[name.Lexeme] = true
	}
}

func (r *Resolver) resolveLocal(expr parser.Expr, name *scanner.Token) {
	for i := len(r.scopes) - 1; i >= 0; i-- {
		if _, ok := r.scopes[i][name.Lexeme]; ok {
			r.i.Resolve(expr, len(r.scopes)-i-1)
			return
		}
	}
}

func (r *Resolver) resolveFunction(function *parser.FunStmt, ftype int) {
	enclosing := r.ftype
	r.ftype = ftype
	r.scopes = append(r.scopes, make(map[string]bool))

	for _, param := range function.Params {
		r.declare(param)
		r.define(param)
	}

	for _, stmt := range function.Body.Statements {
		stmt.Accept(r)
	}

	r.scopes = r.scopes[:len(r.scopes)-1]
	r.ftype = enclosing
}
