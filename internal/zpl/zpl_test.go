package zpl

import "testing"

func TestNewZPL(t *testing.T) {
	zpl := New(10, 20)

	if zpl.width != 10 {
		t.Errorf("TestNewZPL() - zpl.width is incorrect - want %d, got: %d", 76, zpl.width)
	}

	if zpl.height != 20 {
		t.Errorf("TestNewZPL() - zpl.height is incorrect - want %d, got: %d", 51, zpl.width)
	}

	if zpl.orientation != OrientationNormal {
		t.Errorf("TestNewZPL() - zpl.orientation is incorrect - want %s, got: %s", OrientationNormal, zpl.orientation)
	}

	if zpl.dpi != 203 {
		t.Errorf("TestNewZPL() - zpl.dpi is incorrect - want %d, got: %d", 203, zpl.dpi)
	}

	if zpl.dpmm != 8 {
		t.Errorf("TestNewZPL() - zpl.dpmm is incorrect - want %d, got: %d", 8, zpl.dpmm)
	}

	if zpl.unit != Millimeter {
		t.Errorf("TestNewZPL() - zpl.unit is incorrect - want %s, got: %s", Millimeter, zpl.unit)
	}
}
