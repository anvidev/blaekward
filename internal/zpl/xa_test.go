package zpl

import (
	"bytes"
	"testing"
)

func TestXA(t *testing.T) {
	z := &ZPL{buffer: bytes.Buffer{}}
	z.XA()

	if string(z.buffer.Bytes()) != "^XA\n" {
		t.Errorf("TestXA() - command is not added to buffer correctly - want: %s, got: %s", "^XA\n", string(z.buffer.Bytes()))
	}
}
