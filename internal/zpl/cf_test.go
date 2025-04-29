package zpl

import "testing"

func TestCF(t *testing.T) {
	tests := []struct {
		name string
		cmd  Commander
		want string
	}{
		{name: "Default", cmd: CF(), want: "^CF,,"},
		{name: "With font", cmd: CF().SetFont("N"), want: "^CFN,,"},
		{name: "With character height", cmd: CF().SetCharacterHeight(10), want: "^CF,10,"},
		{name: "With character width", cmd: CF().SetCharacterWidth(20), want: "^CF,,20"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.want != tt.cmd.ToZPL() {
				t.Errorf("cf command failed. want %s, got %s", tt.want, tt.cmd.ToZPL())
			}
		})
	}
}
