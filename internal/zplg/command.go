package zplg

import "fmt"

type BQErrorCorrection string
type BCMode string

const (
	ULTRA_HIGH_RELIABILITY BQErrorCorrection = "H"
	HIGH_RELIABILITY                         = "Q" // default if empty
	STANDARD                                 = "M" // default if invalid value
	HIGH_DENSITY                             = "L"

	NONE_MODE      BCMode = "N" // default
	UCC_CASE_MODE         = "U"
	AUTOMATIC_MODE        = "A"
	UCC_EAN_MODE          = "D"
)

func (z *ZPL) FS() {
	z.buffer.WriteString("^FS\n")
}

func (z *ZPL) FO(x, y float64) {
	dotx, doty := z.calculateDots(x), z.calculateDots(y)
	z.buffer.WriteString(fmt.Sprintf("^FO%d,%d", dotx, doty))
}

func (z *ZPL) FD(v string) {
	z.buffer.WriteString(fmt.Sprintf("^FD%s", v))
}

type BYOptions struct {
	moduleWidth   int
	barWidthRatio float64
	barcodeHeight int
}

type BYOptionsFunc func(*BYOptions)

func WithBYModuleWidth(v int) BYOptionsFunc {
	return func(b *BYOptions) { b.moduleWidth = v }
}

func WithBYBarWidthRatio(v float64) BYOptionsFunc {
	return func(b *BYOptions) { b.barWidthRatio = v }
}

func (z *ZPL) BY(barcodeHeight int, opts ...BYOptionsFunc) {
	defaultModuleWidth := 2
	defaultBarRatio := 3.0

	cfg := &BYOptions{
		moduleWidth:   defaultModuleWidth,
		barWidthRatio: defaultBarRatio,
		barcodeHeight: barcodeHeight,
	}

	for _, opt := range opts {
		opt(cfg)
	}

	if cfg.moduleWidth < 1 || cfg.moduleWidth > 10 {
		cfg.moduleWidth = defaultModuleWidth
	}
	if cfg.barWidthRatio < 2.0 || cfg.barWidthRatio > 3.0 {
		cfg.barWidthRatio = defaultBarRatio
	}

	byCmd := fmt.Sprintf(
		"^BY%d,%.1f,%d",
		cfg.moduleWidth,
		cfg.barWidthRatio,
		cfg.barcodeHeight,
	)
	z.buffer.WriteString(byCmd)
}

type BQOptions struct {
	orientation     Orientation
	model           int
	magnification   int
	errorCorrection BQErrorCorrection
	mask            int
}

type BQOptionsFunc func(*BQOptions)

func WithBQModel(v int) BQOptionsFunc {
	return func(b *BQOptions) { b.model = v }
}

func WithBQMagnification(v int) BQOptionsFunc {
	return func(b *BQOptions) { b.magnification = v }
}

func WithBQErrorCorrection(v BQErrorCorrection) BQOptionsFunc {
	return func(b *BQOptions) { b.errorCorrection = v }
}

func WithBQMask(v int) BQOptionsFunc {
	return func(b *BQOptions) { b.mask = v }
}

func (z *ZPL) BQ(v string, opts ...BQOptionsFunc) {
	defaultOrientation := OrientationNormal
	defaultModel := 2
	defaultMagnification := 2
	defaultErrorCorrection := HIGH_RELIABILITY
	defaultMask := 7

	if z.dpi == 150 {
		defaultMagnification = 1
	} else if z.dpi == 200 {
		defaultMagnification = 2
	} else if z.dpi == 300 {
		defaultMagnification = 3
	} else if z.dpi == 600 {
		defaultMagnification = 6
	}

	cfg := &BQOptions{
		orientation:     defaultOrientation,
		model:           defaultModel,
		magnification:   defaultMagnification,
		errorCorrection: BQErrorCorrection(defaultErrorCorrection),
		mask:            defaultMask,
	}

	for _, opt := range opts {
		opt(cfg)
	}

	if cfg.model < 1 || cfg.model > 2 {
		cfg.model = defaultModel
	}

	if cfg.magnification < 1 || cfg.magnification > 100 {
		cfg.magnification = defaultMagnification
	}

	if cfg.mask < 0 || cfg.mask > 7 {
		cfg.mask = defaultMask
	}

	bqCmd := fmt.Sprintf(
		"^BQ%s,%d,%d,%s,%d",
		cfg.orientation,
		cfg.model,
		cfg.magnification,
		cfg.errorCorrection,
		cfg.mask,
	)
	z.buffer.WriteString(bqCmd)
}

type BCOptions struct {
	orientation      Orientation
	height           *int
	showTextBelow    bool
	showTextAbove    bool
	addUCCCheckDigit bool
	mode             BCMode
}

type BCOptionsFunc func(*BCOptions)

func WithBCOrientation(v Orientation) BCOptionsFunc {
	return func(b *BCOptions) { b.orientation = v }
}

func WithBCHeight(v int) BCOptionsFunc {
	return func(b *BCOptions) { b.height = &v }
}

func WithBCTextBelow(v bool) BCOptionsFunc {
	return func(b *BCOptions) { b.showTextBelow = v }
}

func WithBCTextAbove(v bool) BCOptionsFunc {
	return func(b *BCOptions) { b.showTextAbove = v }
}

func WithBCCheckDigit(v bool) BCOptionsFunc {
	return func(b *BCOptions) { b.addUCCCheckDigit = v }
}

func WithBCMode(v BCMode) BCOptionsFunc {
	return func(b *BCOptions) { b.mode = v }
}

func (z *ZPL) BC(opts ...BCOptionsFunc) {
	defaultOritation := z.orientation
	defaultShowTextBelow := true
	defaultShowTextAbove := false
	defaultAddUCCCheckDigit := false
	defaultMode := NONE_MODE

	cfg := &BCOptions{
		orientation:      defaultOritation,
		height:           nil,
		showTextBelow:    defaultShowTextBelow,
		showTextAbove:    defaultShowTextAbove,
		addUCCCheckDigit: defaultAddUCCCheckDigit,
		mode:             defaultMode,
	}

	bcCmd := fmt.Sprintf("^BC")

	for _, opt := range opts {
		opt(cfg)
	}

	if cfg.orientation != defaultOritation {
		fmt.Printf("cfg.orientation was not the same as default. is %s - defualt %s", cfg.orientation, defaultOritation)
		bcCmd += fmt.Sprintf("%s,", cfg.orientation)
	} else {
		bcCmd += ","
	}

	if cfg.height != nil {
		bcCmd += fmt.Sprintf("%d,", *cfg.height)
	} else {
		bcCmd += ","
	}

	if !cfg.showTextBelow {
		bcCmd += fmt.Sprintf("N,")
	} else {
		bcCmd += ","
	}

	if cfg.showTextAbove {
		bcCmd += fmt.Sprintf("Y,")
	} else {
		bcCmd += ","
	}

	if cfg.addUCCCheckDigit {
		bcCmd += fmt.Sprintf("Y,")
	} else {
		bcCmd += ","
	}

	if cfg.mode != defaultMode {
		bcCmd += fmt.Sprintf("%s,", cfg.mode)
	}

	z.buffer.WriteString(bcCmd)
}
