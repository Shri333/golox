package interpreter

import (
	"fmt"

	"github.com/Shri333/golox/fault"
	"github.com/Shri333/golox/scanner"
)

type environment struct {
	values map[string]interface{}
}

func (e *environment) get(name *scanner.Token) interface{} {
	if value, ok := e.values[name.Lexeme]; ok {
		return value
	}

	message := fmt.Sprintf("undefined variable %s", name.Lexeme)
	panic(fault.NewFault(name.Line, message))
}

func (e *environment) assign(name *scanner.Token, value interface{}) {
	if _, ok := e.values[name.Lexeme]; ok {
		e.values[name.Lexeme] = value
	} else {
		message := fmt.Sprintf("undefined variable %s", name.Lexeme)
		panic(fault.NewFault(name.Line, message))
	}
}

func (e *environment) define(name string, value interface{}) {
	e.values[name] = value
}
