package interpreter

import (
	"fmt"
	"time"

	"github.com/Shri333/golox/parser"
)

type callable interface {
	arity() int
	call(i *Interpreter, args []interface{}) interface{}
}

type clock struct{}

func (c *clock) arity() int { return 0 }

func (c *clock) call(i *Interpreter, args []interface{}) interface{} {
	return float64(time.Now().UnixMilli() / 1000)
}

func (c clock) String() string {
	return "<native function clock>"
}

type function struct {
	declaration *parser.FunStmt
	closure     *environment
}

func (f *function) arity() int { return len(f.declaration.Params) }

func (f *function) call(i *Interpreter, args []interface{}) (value interface{}) {
	env := &environment{f.closure, make(map[string]interface{})}
	for i := 0; i < f.arity(); i++ {
		env.define(f.declaration.Params[i].Lexeme, args[i])
	}

	prev := i.current
	defer func() {
		i.current = prev
		r := recover()
		if err, ok := r.(error); ok {
			panic(err)
		} else {
			value = r
		}
	}()

	i.current = env
	for _, stmt := range f.declaration.Body.Statements {
		stmt.Accept(i)
	}

	return value
}

func (f function) String() string {
	return fmt.Sprintf("<function %s>", f.declaration.Name.Lexeme)
}
