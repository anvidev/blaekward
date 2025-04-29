package zpl

import "testing"

func TestAddText(t *testing.T) {
	tests := []struct {
		name string
		cmd  Commander
		want string
	}{
		{name: "Top left empty string", cmd: AddText(0, 0, ""), want: "^FO0,0^FD^FS"},
		{name: "Elsewhere foobar", cmd: AddText(45, 60, "foobar"), want: "^FO45,60^FDfoobar^FS"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.want != tt.cmd.ToZPL() {
				t.Errorf("add text command failed. want %s, got %s", tt.want, tt.cmd.ToZPL())
			}
		})
	}
}
