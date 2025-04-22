package zpl

import "testing"

func TestFO(t *testing.T) {
	tests := []struct {
		name string
		cmd  Commander
		want string
	}{
		{name: "Top left", cmd: FO(0, 0), want: "^FO0,0"},
		{name: "Elsewhere", cmd: FO(230, 43), want: "^FO230,43"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.cmd.ToZPL() != tt.want {
				t.Errorf("fo command failed. want %s, got %s", tt.want, tt.cmd.ToZPL())
			}
		})
	}
}
