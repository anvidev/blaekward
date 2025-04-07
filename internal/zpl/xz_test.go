package zpl

import (
	"bytes"
	"testing"
)

func TestXZ(t *testing.T) {
	z := &zpl{buffer: bytes.Buffer{}}
	z.XZ()

	if string(z.buffer.Bytes()) != "^XZ\n" {
		t.Errorf("TestXZ() - command is not added to buffer correctly - want: %s, got: %s", "^XZ\n", string(z.buffer.Bytes()))
	}
}
