package zpl

import (
	"bytes"
	"testing"
)

func TestBYCommand(t *testing.T) {
	tests := []struct {
		name    string
		options []OptionFunc[BYCommand]
		want    string
	}{
		{name: "Default", options: nil, want: "^BY,,\n"},
		{name: "With Module Width", options: []OptionFunc[BYCommand]{WithModuleWidth(3)}, want: "^BY3,,\n"},
		{name: "With Bar Width Ratio", options: []OptionFunc[BYCommand]{WithBarWidthRatio(2.58)}, want: "^BY,2.6,\n"},
		{name: "With Barcode Height", options: []OptionFunc[BYCommand]{WithBarcodeHeight(135)}, want: "^BY,,135\n"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			z := &zpl{buffer: bytes.Buffer{}}
			z.BY(tt.options...)
			got := string(z.buffer.Bytes())

			if got != tt.want {
				t.Errorf("want %q, got %q", tt.want, got)
			}
		})
	}
}
