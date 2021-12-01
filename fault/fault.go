package fault

import "fmt"

type Fault struct {
	line    int
	message string
}

func (f *Fault) Error() string {
	return fmt.Sprintf("Error (line %d): %s.", f.line, f.message)
}

func NewFault(line int, message string) *Fault {
	fault := &Fault{line, message}
	fmt.Println(fault)
	return fault
}
