package zpl

import "fmt"

type by struct {
	moduleWidth   uint
	barRatio      float32
	barcodeHeight uint
}

func BY() *by {
	return &by{}
}

func (by *by) ToZPL() string {
	cmd := "^BY"

	if by.moduleWidth > 0 {
		cmd += fmt.Sprintf("%d,", by.moduleWidth)
	} else {
		cmd += ","
	}

	if by.barRatio > 0 {
		cmd += fmt.Sprintf("%.1f,", by.barRatio)
	} else {
		cmd += ","
	}

	if by.barcodeHeight > 0 {
		cmd += fmt.Sprintf("%d", by.barcodeHeight)
	}

	return cmd
}

func (by *by) SetModuleWidth(width uint) *by {
	by.moduleWidth = width
	return by
}

func (by *by) SetBarRatio(ratio float32) *by {
	by.barRatio = ratio
	return by
}

func (by *by) SetBarcodeHeight(height uint) *by {
	by.barcodeHeight = height
	return by
}
