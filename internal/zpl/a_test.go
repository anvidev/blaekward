package zpl

import (
	"bytes"
	"testing"
)

func TestA(t *testing.T) {
	z := &zpl{buffer: bytes.Buffer{}}
	z.A("N")

	if string(z.buffer.Bytes()) != "^AN\n" {
		t.Errorf(
			"TestA() - command is not added to buffer correctly - want: %s, got: %s",
			"^AN\n",
			string(z.buffer.Bytes()),
		)
	}
}

func TestAWithFieldOrientation(t *testing.T) {
	z := &zpl{buffer: bytes.Buffer{}}
	z.A("N", WithFieldOrientation(OrientationInverted))

	if string(z.buffer.Bytes()) != "^AN,I\n" {
		t.Errorf(
			"TestA() - command is not added to buffer correctly - want: %s, got: %s",
			"^AN,I\n",
			string(z.buffer.Bytes()),
		)
	}
}

func TestAWithHeight(t *testing.T) {
	z := &zpl{buffer: bytes.Buffer{}}
	z.A("N", WithCharacterHeight(15))

	if string(z.buffer.Bytes()) != "^AN,15\n" {
		t.Errorf(
			"TestA() - command is not added to buffer correctly - want: %s, got: %s",
			"^AN,15\n",
			string(z.buffer.Bytes()),
		)
	}
}

func TestAWithWidth(t *testing.T) {
	z := &zpl{buffer: bytes.Buffer{}}
	z.A("N", WithCharacterWidth(20))

	if string(z.buffer.Bytes()) != "^AN\n" {
		t.Errorf(
			"TestA() - command is not added to buffer correctly - want: %s, got: %s",
			"^AN\n",
			string(z.buffer.Bytes()),
		)
	}
}

func TestAWithHeightAndWidth(t *testing.T) {
	z := &zpl{buffer: bytes.Buffer{}}
	z.A("N", WithCharacterHeight(15), WithCharacterWidth(20))

	if string(z.buffer.Bytes()) != "^AN,15,20\n" {
		t.Errorf(
			"TestA() - command is not added to buffer correctly - want: %s, got: %s",
			"^AN,15,20\n",
			string(z.buffer.Bytes()),
		)
	}
}
