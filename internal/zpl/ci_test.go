package zpl

import "testing"

func TestCI(t *testing.T) {
	tests := []struct {
		name string
		cmd  Commander
		want string
	}{
		{name: "UTF-8", cmd: CI(28), want: "^CI28"},
		{name: "Double byte asian encoding", cmd: CI(14), want: "^CI14"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.cmd.ToZPL() != tt.want {
				t.Errorf("ci command failed. want %s, got %s", tt.want, tt.cmd.ToZPL())
			}
		})
	}
}
