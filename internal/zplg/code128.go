package zplg

import "fmt"

type Code128Mode string

const (
	ModeNone    Code128Mode = "N" // no mode selected
	ModeUCCCase             = "U" // UCC Case mode
	ModeAuto                = "A" // automatic mode
	ModeUCCEAN              = "D" // UCC/EAN mode
)

type code128OptionFunc func(*code128Option)

type code128Option struct {
	moduleWidth             int         // ^BY
	barRatio                float64     // ^BY
	barHeight               int         // ^BY
	orientation             Orientation // ^BC
	barcodeHeight           int         // ^BC
	interpretationLineBelow bool        // ^BC
	interpretationLineAbove bool        // ^BC
	checkDigit              bool        // ^BC
	mode                    Code128Mode // ^BC
}

func (d *ZPLDrawer) AddCode128(
	x,
	y float64,
	content string,
	opts ...code128OptionFunc,
) {
	defaultModuleWidth := 2
	defaultBarRatio := 3.0
	defaultBarHeight := 10

	cfg := &code128Option{
		moduleWidth: defaultModuleWidth,
		barRatio:    defaultBarRatio,
		barHeight:   defaultBarHeight,

		orientation:             d.orientation,
		barcodeHeight:           defaultBarHeight,
		interpretationLineBelow: true,
		interpretationLineAbove: false,
		checkDigit:              false,
		mode:                    ModeNone,
	}

	for _, opt := range opts {
		opt(cfg)
	}

	dotx, doty := d.calculateDots(x), d.calculateDots(y)
	d.buffer.WriteString(fmt.Sprintf("^FO%d,%d", dotx, doty))

	// ^BY command
	if cfg.moduleWidth < 1 && cfg.moduleWidth > 10 {
		cfg.moduleWidth = defaultModuleWidth
	}
	if cfg.barRatio < 2.0 && cfg.barRatio > 3.0 {
		cfg.barRatio = defaultBarRatio
	}

	byCmd := fmt.Sprintf("^BY%d,%.1f,%d", cfg.moduleWidth, cfg.barRatio, cfg.barHeight)
	d.buffer.WriteString(byCmd)

	// ^BC command
	if cfg.barcodeHeight < 1 && cfg.barcodeHeight > 32000 {
		cfg.barcodeHeight = cfg.barHeight
	}
	lineBelow := "Y"
	if !cfg.interpretationLineBelow {
		lineBelow = "N"
	}
	lineAbove := "N"
	if cfg.interpretationLineAbove {
		lineAbove = "Y"
	}
	checkDigit := "N"
	if cfg.checkDigit {
		checkDigit = "Y"
	}

	bcCmd := fmt.Sprintf("^BC%s,%d,%s,%s,%s,%s", cfg.orientation, cfg.barcodeHeight, lineBelow, lineAbove, checkDigit, cfg.mode)
	d.buffer.WriteString(bcCmd)

	d.buffer.WriteString(fmt.Sprintf("^FD%s", content))
	d.buffer.WriteString("^FS\n")
}

func WithModuleWidth(v int) code128OptionFunc {
	return func(o *code128Option) {
		o.moduleWidth = v
	}
}

func WithBarRatio(v float64) code128OptionFunc {
	return func(o *code128Option) {
		o.barRatio = v
	}
}

func WithBarHeight(v int) code128OptionFunc {
	return func(o *code128Option) {
		o.barHeight = v
		o.barcodeHeight = v // default value is inherited from bar height if it not set
	}
}

func WithBarcodeHeight(v int) code128OptionFunc {
	return func(o *code128Option) {
		o.barcodeHeight = v
	}
}

func WithInterpretationLineBelow(v bool) code128OptionFunc {
	return func(o *code128Option) {
		o.interpretationLineBelow = v
	}
}

func WithInterpretationLineAbove(v bool) code128OptionFunc {
	return func(o *code128Option) {
		o.interpretationLineAbove = v
	}
}

func WithCheckDigit(v bool) code128OptionFunc {
	return func(o *code128Option) {
		o.checkDigit = v
	}
}

func WithMode(v Code128Mode) code128OptionFunc {
	return func(o *code128Option) {
		o.mode = v
	}
}

func WithOrientation(v Orientation) code128OptionFunc {
	return func(o *code128Option) {
		o.orientation = v
	}
}

// func (d *ZPLDrawer) AddBarcodeCode128(
// 	x,
// 	y float64,
// 	moduleWidth *int,
// 	barRatio *float64,
// 	barHeight *int,
// 	orientation BarcodeOrientation,
// 	barcodeHeight *int,
// 	interpretationLineBelow *bool,
// 	interpretationLineAbove *bool,
// 	checkDigit *bool,
// 	mode BarcodeMode,
// 	content string,
// ) {
// 	dotx, doty := d.calculateDots(x), d.calculateDots(y)
// 	d.buffer.WriteString(fmt.Sprintf("^FO%d,%d", dotx, doty))
//
// 	// base barcode config
//
// 	defualtModuleWidth := 2
// 	if moduleWidth == nil {
// 		moduleWidth = &defualtModuleWidth
// 	} else if *moduleWidth < 1 {
// 		*moduleWidth = defualtModuleWidth
// 	} else if *moduleWidth > 10 {
// 		*moduleWidth = defualtModuleWidth
// 	}
//
// 	defualtBarRatio := 3.0
// 	if barRatio == nil {
// 		barRatio = &defualtBarRatio
// 	} else if *barRatio < 2.0 {
// 		*barRatio = defualtBarRatio
// 	} else if *barRatio > 3.0 {
// 		*barRatio = defualtBarRatio
// 	} else if !checkDecimalPlaces(1, *barRatio) {
// 		*barRatio = defualtBarRatio
// 	}
//
// 	defualtBarHeight := 10
// 	if barHeight == nil {
// 		barHeight = &defualtBarHeight
// 	}
//
// 	byCmd := fmt.Sprintf("^BY%d,%.1f,%d", *moduleWidth, *barRatio, *barHeight)
//
// 	d.buffer.WriteString(byCmd)
//
// 	// code 128 command
//
// 	defaultBarcodeHeight := *barHeight
// 	if barcodeHeight == nil {
// 		barcodeHeight = &defaultBarcodeHeight
// 	} else if *barcodeHeight < 1 {
// 		*barcodeHeight = defaultBarcodeHeight
// 	} else if *barcodeHeight > 32000 {
// 		*barcodeHeight = defaultBarcodeHeight
// 	}
//
// 	defaultInterpretationLineBelow := true
// 	if interpretationLineBelow == nil {
// 		interpretationLineBelow = &defaultInterpretationLineBelow
// 	}
// 	interpretationLineBelowValue := "Y"
// 	if !*interpretationLineBelow {
// 		interpretationLineBelowValue = "N"
// 	}
//
// 	defaultInterpretationLineAbove := false
// 	if interpretationLineAbove == nil {
// 		interpretationLineAbove = &defaultInterpretationLineAbove
// 	}
// 	interpretationLineAboveValue := "Y"
// 	if !*interpretationLineAbove {
// 		interpretationLineAboveValue = "N"
// 	}
//
// 	defaultCheckDigit := false
// 	if checkDigit == nil {
// 		checkDigit = &defaultCheckDigit
// 	}
// 	checkDigitValue := "Y"
// 	if !*checkDigit {
// 		checkDigitValue = "N"
// 	}
//
// 	bcCmd := fmt.Sprintf(
// 		"^BC%s,%d,%s,%s,%s,%s",
// 		orientation,
// 		*barcodeHeight,
// 		interpretationLineBelowValue,
// 		interpretationLineAboveValue,
// 		checkDigitValue,
// 		mode,
// 	)
// 	d.buffer.WriteString(bcCmd)
//
// 	d.buffer.WriteString(fmt.Sprintf("^FD%s", content))
//
// 	d.buffer.WriteString("^FS\n")
// }
