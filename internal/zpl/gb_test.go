package zpl

import "testing"

func TestGB(t *testing.T) {
	tests := []struct {
		name string
		cmd  Commander
		want string
	}{
		{name: "Default with x and y", cmd: GB(380, 200), want: "^GB380,200,,,"},
		{name: "With border width", cmd: GB(380, 200).SetBorderWidth(20), want: "^GB380,200,20,,"},
		{name: "With line color", cmd: GB(380, 200).SetLineColor("W"), want: "^GB380,200,,W,"},
		{name: "With border radius", cmd: GB(380, 200).SetBorderRadius(8), want: "^GB380,200,,,8"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.want != tt.cmd.ToZPL() {
				t.Errorf("gb command failed. want %s, got %s", tt.want, tt.cmd.ToZPL())
			}
		})
	}
}
