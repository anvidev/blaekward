package zpl

import "fmt"

type gb struct {
	width, height             uint
	borderWidth, borderRadius uint
	lineColor                 string
}

func GB(width, height uint) *gb {
	return &gb{
		width:  width,
		height: height,
	}
}

func (gb *gb) ToZPL() string {
	cmd := fmt.Sprintf("^GB%d,%d,", gb.width, gb.height)

	if gb.borderWidth > 0 {
		cmd += fmt.Sprintf("%d,", gb.borderWidth)
	} else {
		cmd += ","
	}

	if gb.lineColor != "" {
		cmd += fmt.Sprintf("%s,", gb.lineColor)
	} else {
		cmd += ","
	}

	if gb.borderRadius > 0 {
		cmd += fmt.Sprintf("%d", gb.borderRadius)
	}
	return cmd
}

func (gb *gb) SetBorderWidth(width uint) *gb {
	gb.borderWidth = width
	return gb
}

func (gb *gb) SetBorderRadius(radius uint) *gb {
	gb.borderRadius = radius
	return gb
}

func (gb *gb) SetLineColor(color string) *gb {
	gb.lineColor = color
	return gb
}
