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

	result := run(string(bytes))
	if result == 1 {
		os.Exit(65)
	}
	if result == 2 {
		os.Exit(70)
	}
}

func runPrompt() {
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Print("> ")
	for scanner.Scan() {
		text := scanner.Text()
		run(text)
		fmt.Print("> ")
	}

	if err := scanner.Err(); err == nil {
		fmt.Println("bye")
		os.Exit(0)
	}
}

func run(source string) int {
	scanner := scanner.NewScanner(source)
	scanner.ScanTokens()
	if scanner.Error {
		return 1
	}

	parser := ast.NewParser(scanner.Tokens)
	expr := parser.Parse()
	if parser.Error {
		return 1
	}

	interpreter := interpreter.NewInterpreter()
	interpreter.Interpret(expr)
	if interpreter.Error {
		return 2
	}

	return 0
}
