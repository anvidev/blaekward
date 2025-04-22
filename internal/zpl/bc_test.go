package zpl

import "testing"

func TestBC(t *testing.T) {
	tests := []struct {
		name string
		cmd  Commander
		want string
	}{
		{name: "Default", cmd: BC(), want: "^BC,,,,,"},
		{name: "With orientation", cmd: BC().SetOrientation("R"), want: "^BCR,,,,"},
		{name: "With barcode height", cmd: BC().SetBarcodeHeight(30), want: "^BC,30,,,"},
		{name: "With text below (off)", cmd: BC().SetTextBelow(false), want: "^BC,,N,,"},
		{name: "With text above (on)", cmd: BC().SetTextAbove(true), want: "^BC,,,Y,"},
		{name: "With ucc check digit (on)", cmd: BC().SetUCCCheckDigit(true), want: "^BC,,,Y,"},
		{name: "With mode", cmd: BC().SetMode("U"), want: "^BC,,,,U"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.want != tt.cmd.ToZPL() {
				t.Errorf("by command failed. want %s, got %s", tt.want, tt.cmd.ToZPL())
			}
		})
	}
}
