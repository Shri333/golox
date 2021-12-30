package tests

import (
	"bytes"
	"testing"

	"github.com/Shri333/golox/fault"
	"github.com/Shri333/golox/run"
)

func assert(t *testing.T, path, expected string) {
	buf := &bytes.Buffer{}
	fault.W = buf
	run.RunFile(path)
	actual := buf.String()
	if actual != expected {
		t.Errorf("Expected %s but got %s", expected, actual)
	}
}
