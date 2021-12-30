package tests

import "testing"

func TestEmptyFile(t *testing.T) {
	assert(t, "../testdata/empty_file.lox", "")
}

func TestPrecedence(t *testing.T) {
	assert(t, "../testdata/precedence.lox", "14\n8\n4\n0\ntrue\ntrue\ntrue\ntrue\n0\n0\n0\n0\n4\n")
}

func TestUnexpectedCharacter(t *testing.T) {
	assert(t, "../testdata/unexpected_character.lox", "Error (line 3): unknown character '|'\n")
}
