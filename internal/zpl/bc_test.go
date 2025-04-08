package zpl

import (
	"bytes"
	"testing"
)

func TestBCCommand(t *testing.T) {
	tests := []struct {
		name    string
		options []OptionFunc[BCCommand]
		want    string
	}{
		{"Default", nil, "^BC,,,,,\n"},
		{"With Orientation", []OptionFunc[BCCommand]{WithBCOrientation(OrientationInverted)}, "^BCI,,,,,\n"},
		{"With Barcode Height", []OptionFunc[BCCommand]{WithBCBarcodeHeight(60)}, "^BC,60,,,,\n"},
		{"With Text Below", []OptionFunc[BCCommand]{WithBCTextBelow(false)}, "^BC,,N,,,\n"},
		{"With Text Above", []OptionFunc[BCCommand]{WithBCTextAbove(true)}, "^BC,,,Y,,\n"},
		{"With Check Digit", []OptionFunc[BCCommand]{WithBCCheckDigit(true)}, "^BC,,,,Y,\n"},
		{"With Mode", []OptionFunc[BCCommand]{WithBCMode(ModeEAN)}, "^BC,,,,,D\n"},
		{
			"Full",
			[]OptionFunc[BCCommand]{
				WithBCOrientation(OrientationBottomUp),
				WithBCBarcodeHeight(135),
				WithBCTextBelow(false),
				WithBCTextAbove(true),
				WithBCCheckDigit(true),
				WithBCMode(ModeUCC),
			},
			"^BCB,135,N,Y,Y,U\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			z := &zpl{buffer: bytes.Buffer{}}
			z.BC(tt.options...)
			got := string(z.buffer.Bytes())

			if got != tt.want {
				t.Errorf("want %q, got %q", tt.want, got)
			}
		})
	}
}
