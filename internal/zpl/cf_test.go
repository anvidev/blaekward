package zpl

import (
	"bytes"
	"testing"
)

func TestCFCommand(t *testing.T) {
	tests := []struct {
		name    string
		options []OptionFunc[cfCommand]
		want    string
	}{
		{name: "Default", options: nil, want: "^CF,,\n"},
		{name: "With Font", options: []OptionFunc[cfCommand]{WithCFFont("Z")}, want: "^CFZ,,\n"},
		{name: "With Character Height", options: []OptionFunc[cfCommand]{WithCFCharHeight(30)}, want: "^CF,30,\n"},
		{name: "With Character Width", options: []OptionFunc[cfCommand]{WithCFCharWidth(30)}, want: "^CF,,30\n"},
		{name: "Full", options: []OptionFunc[cfCommand]{
			WithCFFont("H"),
			WithCFCharHeight(30),
			WithCFCharWidth(20),
		}, want: "^CFH,30,20\n"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			z := &ZPL{buffer: bytes.Buffer{}}
			z.CF(tt.options...)
			got := string(z.buffer.Bytes())

			if got != tt.want {
				t.Errorf("want %q, got %q", tt.want, got)
			}
		})
	}
}
