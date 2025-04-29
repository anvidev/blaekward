package zpl

import "testing"

func TestBY(t *testing.T) {
	tests := []struct {
		name string
		cmd  Commander
		want string
	}{
		{name: "Default", cmd: BY(), want: "^BY,,"},
		{name: "With module width", cmd: BY().SetModuleWidth(3), want: "^BY3,,"},
		{name: "With bar ratio", cmd: BY().SetBarRatio(2.5), want: "^BY,2.5,"},
		{name: "With bar ratio (multiple decimals)", cmd: BY().SetBarRatio(2.584), want: "^BY,2.6,"},
		{name: "With barcode height", cmd: BY().SetBarcodeHeight(30), want: "^BY,,30"},
		{name: "With module with + barcode height", cmd: BY().SetModuleWidth(3).SetBarcodeHeight(30), want: "^BY3,,30"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.want != tt.cmd.ToZPL() {
				t.Errorf("by command failed. want %s, got %s", tt.want, tt.cmd.ToZPL())
			}
		})
	}
}
