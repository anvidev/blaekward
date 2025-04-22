package zpl

import "testing"

func TestNewLabel(t *testing.T) {
	l := NewLabel()

	if l == nil {
		t.Errorf("new label should not be nil")
	}

	if len(l.commands) != 1 {
		t.Errorf("new label should only have 1 command at start")
	}

	_, ok := l.commands[0].(*xa)
	if !ok {
		t.Errorf("first command in new label should be ^XA")
	}

	zpl := l.String()
	if zpl != "^XA^XZ" {
		t.Errorf("first and last command should be ^XA and ^XZ in that order")
	}

	l.AddCommand(XA())
	if len(l.commands) != 3 {
		t.Errorf("label add command does not add command correctly to command slice")
	}
}
