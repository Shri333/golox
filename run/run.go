package run

import (
	"bufio"
	"fmt"
	"log"
	"os"

	"github.com/Shri333/golox/fault"
	"github.com/Shri333/golox/interpreter"
	"github.com/Shri333/golox/parser"
	"github.com/Shri333/golox/resolver"
	"github.com/Shri333/golox/scanner"
)

func RunFile(path string) {
	bytes, err := os.ReadFile(path)
	if err != nil {
		log.Fatal(err)
	}

	stmts, err := scanAndParse(string(bytes))
	if err != nil && fault.W == os.Stdout {
		os.Exit(65)
	}

	i := interpreter.NewInterpreter()
	r := resolver.NewResolver(i)
	err = r.Resolve(stmts)
	if err != nil && fault.W == os.Stdout {
		os.Exit(65)
	}

	err = i.Interpret(stmts)
	if err != nil && fault.W == os.Stdout {
		os.Exit(70)
	}
}

func RunPrompt() {
	s := bufio.NewScanner(os.Stdin)
	i := interpreter.NewInterpreter()
	fmt.Print("> ")
	for s.Scan() {
		stmts, err := scanAndParse(s.Text())
		if err == nil {
			r := resolver.NewResolver(i)
			err = r.Resolve(stmts)
		}
		if err == nil {
			i.Interpret(stmts)
		}
		fmt.Print("> ")
	}

	if err := s.Err(); err == nil {
		fmt.Println("bye")
		os.Exit(0)
	}
}

func scanAndParse(source string) ([]parser.Stmt, error) {
	s := scanner.NewScanner(source)
	err := s.ScanTokens()
	if err != nil {
		return nil, err
	}

	p := parser.NewParser(s.Tokens)
	stmts, err := p.Parse()
	if err != nil {
		return nil, err
	}

	return stmts, nil
}
