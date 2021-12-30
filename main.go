package main

import (
	"log"
	"os"

	"github.com/Shri333/golox/run"
)

func main() {
	if len(os.Args) > 2 {
		log.Fatal("Usage golox [script]")
	} else if len(os.Args) == 2 {
		run.RunFile(os.Args[1])
	} else {
		run.RunPrompt()
	}
}
