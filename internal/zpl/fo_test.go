package zpl

import (
	"bytes"
	"testing"
)

func TestFO(t *testing.T) {
	z := &ZPL{buffer: bytes.Buffer{}}
	z.FO(51, 26)

	if string(z.buffer.Bytes()) != "^FO51,26\n" {
		t.Errorf("TestFO() - command is not added to buffer correctly - want: %s, got: %s", "^FO51,26\n", string(z.buffer.Bytes()))
	}
}

func TestFOWithJustify(t *testing.T) {
	z := &ZPL{buffer: bytes.Buffer{}}
	z.FO(51, 26, WithJustify(JustifyRight))

	if string(z.buffer.Bytes()) != "^FO51,26,1\n" {
		t.Errorf("TestFO() - command is not added to buffer correctly - want: %s, got: %s", "^FO51,261\n", string(z.buffer.Bytes()))
	}
}
