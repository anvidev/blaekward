package zplg

import "fmt"

func (d *ZPL) AddText(x, y float64, content string, height int, width *int) {
	dotx, doty := d.calculateDots(x), d.calculateDots(y)
	d.buffer.WriteString(fmt.Sprintf("^FO%d,%d", dotx, doty))

	fontCmd := fmt.Sprintf("^A%d,%d", 0, pointsToDots(float64(height), d.dpi))
	if width != nil {
		fontCmd += fmt.Sprintf(",%d", *width)
	}
	d.buffer.WriteString(fontCmd)

	d.buffer.WriteString(fmt.Sprintf("^FD%s", content))

	d.buffer.WriteString("^FS\n")
}
