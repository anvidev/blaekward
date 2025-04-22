package zpl

import "testing"

func TestFD(t *testing.T) {
	tests := []struct {
		name string
		cmd  Commander
		want string
	}{
		{name: "Empty string", cmd: FD(""), want: "^FD"},
		{name: "Foobar", cmd: FD("Foobar"), want: "^FDFoobar"},
		{name: "Numbers", cmd: FD("123321"), want: "^FD123321"},
		{name: "Symbols", cmd: FD("()[]{}?=+-"), want: "^FD()[]{}?=+-"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.cmd.ToZPL() != tt.want {
				t.Errorf("fd command failed. want %s, got %s", tt.want, tt.cmd.ToZPL())
			}
		})
	}
}
