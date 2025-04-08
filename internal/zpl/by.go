package zpl

import "fmt"

type BYCommand struct {
	moduleWidth   *int
	barWidthRatio *float64
	barcodeHeight *int
}

func (z *zpl) BY(opts ...OptionFunc[BYCommand]) {
	cfg := &BYCommand{}

	for _, opt := range opts {
		opt(cfg)
	}

	cmd := "^BY"

	if cfg.moduleWidth != nil {
		cmd += fmt.Sprintf("%d,", *cfg.moduleWidth)
	} else {
		cmd += ","
	}

	if cfg.barWidthRatio != nil {
		cmd += fmt.Sprintf("%.1f,", *cfg.barWidthRatio)
	} else {
		cmd += ","
	}

	if cfg.barcodeHeight != nil {
		cmd += fmt.Sprintf("%d", *cfg.barcodeHeight)
	}

	z.writeCmd(cmd)
}

func WithModuleWidth(v int) OptionFunc[BYCommand] {
	return func(b *BYCommand) { b.moduleWidth = &v }
}

func WithBarWidthRatio(v float64) OptionFunc[BYCommand] {
	return func(b *BYCommand) { b.barWidthRatio = &v }
}

func WithBarcodeHeight(v int) OptionFunc[BYCommand] {
	return func(b *BYCommand) { b.barcodeHeight = &v }
}
