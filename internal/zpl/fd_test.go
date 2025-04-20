package zpl

import (
	"bytes"
	"testing"
)

func TestFD(t *testing.T) {
	z := &ZPL{buffer: bytes.Buffer{}}
	z.FD("ZPL commands are awesome!")

	if string(z.buffer.Bytes()) != "^FDZPL commands are awesome!\n" {
		t.Errorf(
			"TestFD() - command is not added to buffer correctly - want: %s, got: %s",
			"^FDZPL commands are awesome!\n",
			string(z.buffer.Bytes()),
		)
	}
}
