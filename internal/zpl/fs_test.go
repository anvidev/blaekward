package zpl

import (
	"bytes"
	"testing"
)

func TestFS(t *testing.T) {
	z := &ZPL{buffer: bytes.Buffer{}}
	z.FS()

	if string(z.buffer.Bytes()) != "^FS\n" {
		t.Errorf("TestFO() - command is not added to buffer correctly - want: %s, got: %s", "^FS\n", string(z.buffer.Bytes()))
	}
}
