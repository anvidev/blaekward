package zplg

import (
	"bytes"
	"fmt"
	"math"
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
	Orientation90Deg                = "R" // 90 degrees (clockwise)
	Orientation180Deg               = "I" // inverted 180 degress
	OrientationBottomUp             = "B" // read from bottom up (270 degrees)
)

type ZPLDrawer struct {
	dpi         int
	unit        Unit
	width       int
	height      int
	orientation Orientation

	buffer bytes.Buffer
}

func New(dpi, width, height int, unit Unit, orientation Orientation) *ZPLDrawer {
	zpld := &ZPLDrawer{
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

func (d *ZPLDrawer) Bytes() []byte {
	d.buffer.WriteString("^XZ")
	return d.buffer.Bytes()
}

func (d *ZPLDrawer) String() string {
	d.buffer.WriteString("^XZ")
	return d.buffer.String()
}

func (d *ZPLDrawer) calculateDots(v float64) int {
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

func checkDecimalPlaces(i int, value float64) bool {
	valuef := value * float64(math.Pow(10.0, float64(i)))
	extra := valuef - float64(int(valuef))
	return extra == 0
}
