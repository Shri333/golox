package interpreter

import (
	"fmt"

	"github.com/Shri333/golox/fault"
	"github.com/Shri333/golox/scanner"
)

type environment struct {
	enclosing *environment
	values    map[string]interface{}
}

func (e *environment) get(name *scanner.Token) interface{} {
	if value, ok := e.values[name.Lexeme]; ok {
		return value
	}

	if e.enclosing != nil {
		return e.enclosing.get(name)
	}

	message := fmt.Sprintf("undefined variable %s", name.Lexeme)
	panic(fault.NewFault(name.Line, message))
}

func (e *environment) getAt(name string, dist int) interface{} {
	ancestor := e
	for i := 0; i < dist; i++ {
		ancestor = ancestor.enclosing
	}

	return ancestor.values[name]
}

func (e *environment) assign(name *scanner.Token, value interface{}) {
	if _, ok := e.values[name.Lexeme]; ok {
		e.values[name.Lexeme] = value
	} else if e.enclosing != nil {
		e.enclosing.assign(name, value)
	} else {
		message := fmt.Sprintf("undefined variable %s", name.Lexeme)
		panic(fault.NewFault(name.Line, message))
	}
}

func (e *environment) assignAt(name string, value interface{}, dist int) {
	ancestor := e
	for i := 0; i < dist; i++ {
		ancestor = ancestor.enclosing
	}

	ancestor.values[name] = value
}

func (e *environment) define(name string, value interface{}) {
	e.values[name] = value
}
