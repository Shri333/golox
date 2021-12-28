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
	init        bool
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
			if f.init {
				value = f.closure.getAt("this", 0)
			} else {
				value = r
			}
		}
	}()

	i.current = env
	for _, stmt := range f.declaration.Body.Statements {
		stmt.Accept(i)
	}

	if f.init {
		return f.closure.getAt("this", 0)
	}

	return value
}

func (f *function) bind(i *instance) *function {
	env := &environment{f.closure, make(map[string]interface{})}
	env.define("this", i)
	return &function{f.declaration, env, f.init}
}

func (f function) String() string {
	return fmt.Sprintf("<function %s>", f.declaration.Name.Lexeme)
}

type class struct {
	name    string
	methods map[string]*function
}

func (c *class) arity() int {
	initializer := c.findMethod("init")
	if initializer != nil {
		return initializer.arity()
	}

	return 0
}

func (c *class) call(i *Interpreter, args []interface{}) interface{} {
	inst := &instance{c, make(map[string]interface{})}
	initializer := c.findMethod("init")
	if initializer != nil {
		initializer.bind(inst).call(i, args)
	}

	return inst
}

func (c *class) findMethod(name string) *function {
	if fn, ok := c.methods[name]; ok {
		return fn
	}

	return nil
}

func (c class) String() string {
	return fmt.Sprintf("<class %s>", c.name)
}
