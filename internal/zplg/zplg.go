package zplg

import (
	"bytes"
	"fmt"
)

// https://labelary.com/viewer.html
// https://www.servopack.de/support/zebra/ZPLII-Prog.pdf

type Unit string
type Orientation string

const (
	Inches     Unit = "inches"
	Centimeter      = "cm"
	Millimeter      = "mm"

	OrientationNormal   Orientation = "N" // normal
	OrientationRotated              = "R" // 90 degrees (clockwise)
	OrientationInverted             = "I" // inverted 180 degress
	OrientationBottomUp             = "B" // read from bottom up (270 degrees)
)

type ZPL struct {
	dpi         int
	unit        Unit
	width       int
	height      int
	orientation Orientation

	buffer bytes.Buffer
}

func New(dpi, width, height int, unit Unit, orientation Orientation) *ZPL {
	zpld := &ZPL{
		dpi:         dpi,
		unit:        unit,
		width:       width,
		height:      height,
		orientation: orientation,
		buffer:      bytes.Buffer{},
	}
	zpld.buffer.WriteString("^XA\n")
	zpld.buffer.WriteString(fmt.Sprintf("^FW%s\n", zpld.orientation))
	return zpld
}

func (d *ZPL) Bytes() []byte {
	d.buffer.WriteString("^XZ")
	return d.buffer.Bytes()
}

func (d *ZPL) String() string {
	d.buffer.WriteString("^XZ")
	return d.buffer.String()
}

func (z *ZPL) AddQRCode(x, y float64, height int, data string) {
	z.FO(x, y)
	z.BY(height)
	z.BQ(data)
}

func (d *ZPL) calculateDots(v float64) int {
	switch d.unit {
	case Inches:
		return inchesToDots(v, d.dpi)
	case Centimeter:
		return centimetersToDots(v, d.dpi)
	case Millimeter:
		return millimetersToDots(v, d.dpi)
	default:
		return -1
	}
}

func inchesToDots(inches float64, dpi int) int {
	return int(inches * float64(dpi))
}

func centimetersToDots(cm float64, dpi int) int {
	return int((cm / 2.54) * float64(dpi))
}

func millimetersToDots(mm float64, dpi int) int {
	return int((mm / 25.4) * float64(dpi))
}

func pointsToDots(points float64, dpi int) int {
	return int((points / 72) * float64(dpi))
}

func isAlphaNumeric(c byte) bool {
	return (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || (c >= '0' && c <= '9')
}
