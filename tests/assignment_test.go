package tests

import "testing"

func TestAssociativity(t *testing.T) {
	assert(t, "../testdata/assignment/associativity.lox", "c\nc\nc\n")
}

func TestGlobal(t *testing.T) {
	assert(t, "../testdata/assignment/global.lox", "before\nafter\narg\narg\n")
}

func TestGrouping(t *testing.T) {
	assert(t, "../testdata/assignment/grouping.lox", "Error (line 2): invalid assignment target\n")
}

func TestInfixOperator(t *testing.T) {
	assert(t, "../testdata/assignment/infix_operator.lox", "Error (line 3): invalid assignment target\n")
}

func TestLocal(t *testing.T) {
	assert(t, "../testdata/assignment/local.lox", "before\nafter\narg\narg\n")
}

func TestPrefixOperator(t *testing.T) {
	assert(t, "../testdata/assignment/prefix_operator.lox", "Error (line 2): invalid assignment target\n")
}

func TestSyntax(t *testing.T) {
	assert(t, "../testdata/assignment/syntax.lox", "var\nvar\n")
}

func TestToThis(t *testing.T) {
	assert(t, "../testdata/assignment/to_this.lox", "Error (line 3): invalid assignment target\n")
}

func TestUndefined(t *testing.T) {
	assert(t, "../testdata/assignment/undefined.lox", "Error (line 1): undefined variable 'unknown'\n")
}
