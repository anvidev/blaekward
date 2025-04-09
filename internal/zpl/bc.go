package zpl

import "fmt"

type BCCommand struct {
	orientation      *Orientation
	barcodeHeight    *int
	showTextBelow    *bool
	showTextAbove    *bool
	addUccCheckDigit *bool
	mode             *Mode
}

func (z *ZPL) BC(opts ...OptionFunc[BCCommand]) {
	cfg := &BCCommand{}

	for _, opt := range opts {
		opt(cfg)
	}

	cmd := "^BC"

	if cfg.orientation != nil {
		cmd += fmt.Sprintf("%s,", *cfg.orientation)
	} else {
		cmd += ","
	}

	if cfg.barcodeHeight != nil {
		cmd += fmt.Sprintf("%d,", *cfg.barcodeHeight)
	} else {
		cmd += ","
	}

	if cfg.showTextBelow != nil {
		value := "Y"
		if !*cfg.showTextBelow {
			value = "N"
		}
		cmd += fmt.Sprintf("%s,", value)
	} else {
		cmd += ","
	}

	if cfg.showTextAbove != nil {
		value := "N"
		if *cfg.showTextAbove {
			value = "Y"
		}
		cmd += fmt.Sprintf("%s,", value)
	} else {
		cmd += ","
	}

	if cfg.addUccCheckDigit != nil {
		value := "N"
		if *cfg.addUccCheckDigit {
			value = "Y"
		}
		cmd += fmt.Sprintf("%s,", value)
	} else {
		cmd += ","
	}

	if cfg.mode != nil {
		cmd += fmt.Sprintf("%s", *cfg.mode)
	}

	z.writeCmd(cmd)
}

func WithBCOrientation(v Orientation) OptionFunc[BCCommand] {
	return func(b *BCCommand) { b.orientation = &v }
}

func WithBCBarcodeHeight(v int) OptionFunc[BCCommand] {
	return func(b *BCCommand) { b.barcodeHeight = &v }
}

func WithBCTextBelow(v bool) OptionFunc[BCCommand] {
	return func(b *BCCommand) { b.showTextBelow = &v }
}

func WithBCTextAbove(v bool) OptionFunc[BCCommand] {
	return func(b *BCCommand) { b.showTextAbove = &v }
}

func WithBCCheckDigit(v bool) OptionFunc[BCCommand] {
	return func(b *BCCommand) { b.addUccCheckDigit = &v }
}

func WithBCMode(v Mode) OptionFunc[BCCommand] {
	return func(b *BCCommand) { b.mode = &v }
}
