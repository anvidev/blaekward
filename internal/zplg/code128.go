package zplg

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

func (d *ZPL) AddCode128(x, y float64, content string, opts ...code128OptionFunc) {
	defaultModuleWidth := 2
	defaultBarRatio := 3.0
	defaultBarHeight := 10

	cfg := &code128Option{
		moduleWidth: defaultModuleWidth,
		barRatio:    defaultBarRatio,
		barHeight:   defaultBarHeight,

		orientation:             d.orientation,
		barcodeHeight:           defaultBarHeight, // TODO: remove from function
		interpretationLineBelow: true,
		interpretationLineAbove: false,
		checkDigit:              false,
		mode:                    ModeNone,
	}

	for _, opt := range opts {
		opt(cfg)
	}

	d.FO(x, y)
	d.BY(
		cfg.barHeight,
		WithBYBarWidthRatio(cfg.barRatio),
		WithBYModuleWidth(cfg.moduleWidth),
	)
	d.BC(
		WithBCOrientation(cfg.orientation),
		WithBCTextBelow(cfg.interpretationLineBelow),
		WithBCTextAbove(cfg.interpretationLineAbove),
		WithBCCheckDigit(cfg.checkDigit),
		WithBCMode(BCMode(cfg.mode)),
	)
	d.FD(content)
	d.FS()
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
