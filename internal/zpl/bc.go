package zpl

import "fmt"

type bc struct {
	orientation   string
	barcodeHeight uint
	textBelow     *bool
	textAbove     *bool
	uccCheckDigit *bool
	mode          string
}

func BC() *bc {
	return &bc{}
}

func (bc *bc) ToZPL() string {
	cmd := "^BC"

	if bc.orientation != "" {
		cmd += fmt.Sprintf("%s,", bc.orientation)
	} else {
		cmd += ","
	}

	if bc.barcodeHeight > 0 {
		cmd += fmt.Sprintf("%d,", bc.barcodeHeight)
	} else {
		cmd += ","
	}

	if bc.textBelow != nil {
		cmd += fmt.Sprintf("%s,", boolToZPL(*bc.textBelow))
	} else {
		cmd += ","
	}

	if bc.textAbove != nil {
		cmd += fmt.Sprintf("%s,", boolToZPL(*bc.textAbove))
	} else {
		cmd += ","
	}

	if bc.uccCheckDigit != nil {
		cmd += fmt.Sprintf("%s,", boolToZPL(*bc.uccCheckDigit))
	} else {
		cmd += ","
	}

	if bc.mode != "" {
		cmd += fmt.Sprintf("%s", bc.mode)
	}

	return cmd
}

func (bc *bc) SetOrientation(orientation string) *bc {
	bc.orientation = orientation
	return bc
}

func (bc *bc) SetBarcodeHeight(height uint) *bc {
	bc.barcodeHeight = height
	return bc
}

func (bc *bc) SetTextBelow(on bool) *bc {
	bc.textBelow = &on
	return bc
}

func (bc *bc) SetTextAbove(on bool) *bc {
	bc.textAbove = &on
	return bc
}

func (bc *bc) SetUCCCheckDigit(on bool) *bc {
	bc.uccCheckDigit = &on
	return bc
}

func (bc *bc) SetMode(mode string) *bc {
	bc.mode = mode
	return bc
}
