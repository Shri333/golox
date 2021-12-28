package interpreter

import (
	"fmt"

	"github.com/Shri333/golox/fault"
	"github.com/Shri333/golox/scanner"
)

type instance struct {
	c      *class
	fields map[string]interface{}
}

func (i *instance) get(name *scanner.Token) interface{} {
	if value, ok := i.fields[name.Lexeme]; ok {
		return value
	}

	method := i.c.findMethod(name.Lexeme)
	if method != nil {
		return method.bind(i)
	}

	message := fmt.Sprintf("undefined property %s", name.Lexeme)
	panic(fault.NewFault(name.Line, message))
}

func (i *instance) set(name *scanner.Token, value interface{}) {
	i.fields[name.Lexeme] = value
}

func (i instance) String() string {
	return fmt.Sprintf("%s instance", i.c.name)
}
