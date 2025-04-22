package zpl

import (
	"fmt"
)

type cf struct {
	font          string
	height, width uint
}

func CF() *cf {
	return &cf{}
}

func (cf *cf) ToZPL() string {
	cmd := "^CF"

	if cf.font != "" {
		cmd += fmt.Sprintf("%s,", cf.font)
	} else {
		cmd += ","
	}

	if cf.height > 0 {
		cmd += fmt.Sprintf("%d,", cf.height)
	} else {
		cmd += ","
	}

	if cf.width > 0 {
		cmd += fmt.Sprintf("%d", cf.width)
	}

	return cmd
}

func (cf *cf) SetFont(font string) *cf {
	cf.font = font
	return cf
}

func (cf *cf) SetCharacterHeight(height uint) *cf {
	cf.height = height
	return cf
}

func (cf *cf) SetCharacterWidth(width uint) *cf {
	cf.width = width
	return cf
}
