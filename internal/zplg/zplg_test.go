package zplg

import (
	"testing"
)

func TestAddBardcodeBytes(t *testing.T) {
	want := []byte("^XA\n^FWN\n^FO40,81^BY3,3.0,60^BCN,60,Y,N,N,N^FD1234ABC^FS\n^XZ")
	zpl := New(208, 51, 26, Millimeter, OrientationNormal)
	zpl.AddCode128(5, 10, "1234ABC", WithModuleWidth(3), WithBarHeight(60))
	got := zpl.Bytes()

	if string(want) != string(got) {
		t.Errorf("TestAddBarcodeBytes() failed - want: %s, got: %s", string(want), string(got))
	}
}

func TestAddTextBytes(t *testing.T) {
	want := []byte("^XA\n^FWN\n^FO40,40^A0,40^FDLabels styrer!^FS\n^XZ")
	zpl := New(208, 51, 26, Millimeter, OrientationNormal)
	zpl.AddText(5, 5, "Labels styrer!", 40, nil)
	got := zpl.Bytes()

	if string(want) != string(got) {
		t.Errorf("TestAddTextBytes() failed - want: %s, got: %s", string(want), string(got))
	}
}

func TestAddTextString(t *testing.T) {
	want := "^XA\n^FWN\n^FO40,40^A0,40^FDLabels styrer!^FS\n^XZ"
	zpl := New(208, 51, 26, Millimeter, OrientationNormal)
	zpl.AddText(5, 5, "Labels styrer!", 40, nil)
	got := zpl.String()

	if want != got {
		t.Errorf("TestAddTextString() failed - want: %s, got: %s", want, got)
	}
}

func TestEmptyLabelBytes(t *testing.T) {
	want := []byte("^XA\n^FWN\n^XZ")
	zpl := New(208, 51, 26, Millimeter, OrientationNormal)
	got := zpl.Bytes()

	if string(want) != string(got) {
		t.Errorf("TestEmptyLabelBytes() failed - want: %s, got: %s", string(want), string(got))
	}
}

func TestEmptyLabelString(t *testing.T) {
	want := "^XA\n^FWN\n^XZ"
	zpl := New(208, 51, 26, Millimeter, OrientationNormal)
	got := zpl.String()

	if want != got {
		t.Errorf("TestEmptyLabelString() failed - want: %s, got: %s", want, got)
	}
}

func TestCalculateDots(t *testing.T) {
	tests := []struct {
		dpi   int
		unit  Unit
		want  int
		input float64
	}{
		{dpi: 203, unit: Inches, input: 4.3, want: 872},
		{dpi: 152, unit: Centimeter, input: 15, want: 897},
		{dpi: 300, unit: Millimeter, input: 52, want: 614},
	}

	for _, tt := range tests {
		var got int
		switch tt.unit {
		case Inches:
			got = inchesToDots(tt.input, tt.dpi)
		case Centimeter:
			got = centimetersToDots(tt.input, tt.dpi)
		case Millimeter:
			got = millimetersToDots(tt.input, tt.dpi)
		}

		if got != tt.want {
			t.Errorf("TestCalculateDots() failed - want: %d, got: %d", tt.want, got)
		}
	}
}
