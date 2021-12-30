# golox

Tree-walk interpreter for the Lox programming language.

To build the interpreter (using a modern Go toolchain), run `go build` in the root directory of this repository.
From there, run `./golox` with the name of the Lox source file (or without a source file to start the REPL).

To run the tests, run `go test` inside of the `tests` directory.
All tests in `testdata` come from the [Crafting Interpreters repository](https://github.com/munificent/craftinginterpreters/tree/master/test).

Thank you Bob Nystrom for writing such an excellent book!
