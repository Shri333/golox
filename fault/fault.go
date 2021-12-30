package fault

import (
	"fmt"
	"io"
	"os"
)

var W io.Writer = os.Stdout

type Fault struct {
	line    int
	message string
}

func (f *Fault) Error() string {
	return fmt.Sprintf("Error (line %d): %s", f.line, f.message)
}

func NewFault(line int, message string) *Fault {
	fault := &Fault{line, message}
	fmt.Fprintln(W, fault)
	return fault
}
