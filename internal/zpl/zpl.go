package zpl

// https://labelary.com/viewer.html
// https://www.servopack.de/support/zebra/ZPLII-Prog.pdf

type unit string

const (
	Inches     unit = "inches"
	Centimeter      = "cm"
	Millimeter      = "mm"
)

type ZPLLabel struct {
	dppm int // dpi

	unit   unit
	width  int
	height int
}

func (l *ZPLLabel) draw() []byte {
	return []byte("")
}
