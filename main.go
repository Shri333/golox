package main

import (
	"bufio"
	"fmt"
	"log"
	"os"

	"github.com/Shri333/golox/ast"
	"github.com/Shri333/golox/interpreter"
	"github.com/Shri333/golox/scanner"
)

func main() {
	if len(os.Args) > 2 {
		log.Fatal("Usage golox [script]")
	} else if len(os.Args) == 2 {
		runFile(os.Args[1])
	} else {
		runPrompt()
	}
}

func runFile(path string) {
	bytes, err := os.ReadFile(path)
	if err != nil {
		log.Fatal(err)
	}

	stmts, err := scanAndParse(string(bytes))
	if err != nil {
		os.Exit(65)
	}

	interpreter := interpreter.NewInterpreter()
	err = interpreter.Interpret(stmts)
	if err != nil {
		os.Exit(70)
	}
}

func runPrompt() {
	scanner := bufio.NewScanner(os.Stdin)
	interpreter := interpreter.NewInterpreter()
	fmt.Print("> ")
	for scanner.Scan() {
		stmts, err := scanAndParse(scanner.Text())
		if err == nil {
			interpreter.Interpret(stmts)
		}
		fmt.Print("> ")
	}

	if err := scanner.Err(); err == nil {
		fmt.Println("bye")
		os.Exit(0)
	}
}

func scanAndParse(source string) ([]ast.Stmt, error) {
	scanner := scanner.NewScanner(source)
	err := scanner.ScanTokens()
	if err != nil {
		return nil, err
	}

	parser := ast.NewParser(scanner.Tokens)
	stmts, err := parser.Parse()
	if err != nil {
		return nil, err
	}

	return stmts, nil
}
