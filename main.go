package main

import (
	"bufio"
	"fmt"
	"log"
	"os"

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

	run(string(bytes))
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

func run(source string) {
	scanner := scanner.NewScanner(source)
	scanner.ScanTokens()

	for _, token := range scanner.Tokens {
		fmt.Println(token)
	}
}
