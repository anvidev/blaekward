package zpl

import "bytes"

type OptionFunc[T any] func(*T)

type Unit string
type Orientation string
type Justification int
type Mode string

const (
	Inches     Unit = "inches"
	Millimeter      = "mm"
	Centimeter      = "cm"

	OrientationNormal   Orientation = "N"
	OrientationRotated              = "R"
	OrientationInverted             = "I"
	OrientationBottomUp             = "B"

	JustifyLeft  Justification = 0
	JustifyRight               = 1
	JustifyAuto                = 2

	ModeNone Mode = "N"
	ModeUCC       = "U"
	ModeAuto      = "A"
	ModeEAN       = "D"
)

type zpl struct {
	// Internal buffer for writing zpl commands.
	buffer bytes.Buffer

	// DPI of label. Default is 203 dpi
	dpi int

	// DPMM of label (used for easy conversion of mm/cm to dots).
	dpmm int

	// Unit of label. Defualt is mm.
	unit Unit

	// Width of label in dots.
	width int

	// Height of label in dots.
	height int

	// Orientation of label. Default: N (normal)
	orientation Orientation
}

func New(width, height int, opts ...OptionFunc[zpl]) *zpl {
	z := &zpl{
		buffer:      bytes.Buffer{},
		dpi:         203,
		dpmm:        8,
		unit:        Millimeter,
		width:       width,
		height:      height,
		orientation: OrientationNormal,
	}

	z.XA()
	return z
}

func (z *zpl) Bytes() []byte {
	z.XZ()
	return z.buffer.Bytes()
}

func WithOrientation(o Orientation) OptionFunc[zpl] {
	return func(z *zpl) { z.orientation = o }
}

func WithUnit(unit Unit) OptionFunc[zpl] {
	return func(z *zpl) { z.unit = unit }
}

func WithDpi(v int) OptionFunc[zpl] {
	return func(z *zpl) {
		z.dpi = v
		z.dpmm = dpiToDpmm(v)
	}
}

func InchesToDots(inches float64, dpi int) int {
	return int(inches * float64(dpi))
}

func MmToDots(mm float64, dpi int) int {
	dpmm := float64(dpi) / 25.4
	return int(mm * dpmm)
}

func CmToDots(cm float64, dpi int) int {
	dpcm := float64(dpi) / 2.54
	return int(cm * dpcm)
}

func dpiToDpmm(dpi int) int {
	return int(float64(dpi) / 25.4)
}

func (z *zpl) writeCmd(cmd string) {
	z.buffer.WriteString(cmd)
	z.buffer.WriteByte('\n')
}
