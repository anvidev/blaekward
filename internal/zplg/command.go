package zplg

import (
	"fmt"
)

type BQErrorCorrection string
type BCMode string

const (
	ULTRA_HIGH_RELIABILITY BQErrorCorrection = "H"
	HIGH_RELIABILITY                         = "Q"
	STANDARD                                 = "M"
	HIGH_DENSITY                             = "L"

	NONE_MODE BCMode = "N"
	UCC_MODE         = "U"
	AUTO_MODE        = "A"
	EAN_MODE         = "D"
)

// FS wrties the ZPL command (^FS) to [*ZPL] internal bytes.Buffer and denotes the end of a field definition.
// It is recommended to call FS efter every command that creates a printable line.
func (z *ZPL) FS() {
	z.buffer.WriteString("^FS\n")
}

// FO writes the ZPL command (^FO) to [*ZPL] internal bytes.Buffer and sets the field origin relative to the label
// dimensions in the specified unit which is converted into dots.
func (z *ZPL) FO(x, y float64) {
	dotx, doty := z.calculateDots(x), z.calculateDots(y)
	z.buffer.WriteString(fmt.Sprintf("^FO%d,%d", dotx, doty))
}

// FD writes the ZPL command (^FD) to [*ZPL] internal bytes.Buffer and defines the data string for a field.
// The field data [v] can be any printable character except those used as command prefixes (^ and ~) and can
// be any data string up to 3072 bytes.
func (z *ZPL) FD(v string) {
	z.buffer.WriteString(fmt.Sprintf("^FD%s", v))
}

type CFOptions struct {
	font       string
	charHeight int
	charWidth  int
}

type CFOptionsFunc func(*CFOptions)

func WithCFFont(v string) CFOptionsFunc {
	return func(c *CFOptions) { c.font = v }
}

func WithCFCharHeight(v int) CFOptionsFunc {
	return func(c *CFOptions) { c.charHeight = v }
}

func WithCFCharWidth(v int) CFOptionsFunc {
	return func(c *CFOptions) { c.charWidth = v }
}

func (z *ZPL) CF(opts ...CFOptionsFunc) {
	defaultFont := "A"
	defaultCharHeight := 9
	defaultCharWidth := 5 // TODO implement last permanent saved value into ZPL

	cfg := &CFOptions{
		font:       defaultFont,
		charHeight: defaultCharHeight,
		charWidth:  defaultCharWidth,
	}

	for _, opt := range opts {
		opt(cfg)
	}

	if len(cfg.font) > 1 {
		cfg.font = defaultFont
	}

	if cfg.charHeight < 0 || cfg.charHeight > 32000 {
		cfg.charHeight = defaultCharHeight
	}

	if cfg.charWidth < 0 || cfg.charWidth > 32000 {
		cfg.charWidth = defaultCharWidth
	}

	cfCmd := fmt.Sprintf("^CF")

	if cfg.font != defaultFont {
		cfCmd += fmt.Sprintf("%s,", cfg.font)
	} else {
		cfCmd += ","
	}

	if cfg.charHeight != defaultCharHeight {
		cfCmd += fmt.Sprintf("%d,", cfg.charHeight)
	} else {
		cfCmd += ","
	}

	if cfg.charWidth != defaultCharWidth {
		cfCmd += fmt.Sprintf("%d,", cfg.charWidth)
	} else {
		cfCmd += ","
	}

	z.buffer.WriteString(cfCmd)
}

// BYOptions are parameters used to configure the BY function call
type BYOptions struct {
	// Accepted values are 1 to 10. Default value is 2.
	moduleWidth int

	// Accepted values are 2.0 to 3.0, in 0.1 increments. Default is 3.0.
	barWidthRatio float64

	// Default value is 10.
	barcodeHeight int
}

type BYOptionsFunc func(*BYOptions)

func WithBYModuleWidth(v int) BYOptionsFunc {
	return func(b *BYOptions) { b.moduleWidth = v }
}

func WithBYBarWidthRatio(v float64) BYOptionsFunc {
	return func(b *BYOptions) { b.barWidthRatio = v }
}

// BY writes the ZPL command (^BY) to [*ZPL] internal bytes.Buffer and is used to change the default values for module
// width (in dots), the wide bar to narrow bar width ratio (in dots) and the barcode height (in dots).
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

// BQOptions are parameters used to configure the BQ function call
type BQOptions struct {
	// Default value is normal [OrientationNormal] (N). Other values has no effect on this command.
	orientation Orientation

	// Accepted values are 1 (original) and 2 (enhanced). Default is 2.
	model int

	// Accepted values are 1 to 100. Default value depends in DPI
	//
	// Defaults:
	// 1 on 150 DPI printers
	// 2 on 200 DPI printers
	// 3 on 300 DPI printers
	// 6 on 600 DPI printers
	magnification int

	// Default value is [HIGH_RELIABILITY] (Q)
	errorCorrection BQErrorCorrection

	// Accepted values are 0 to 7. Default is 7.
	mask int
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

// BQ writes the ZPL command (^BQ) to [*ZPL] internal bytes.Buffer and produces a matrix symbology of an array of
// nominally square modules arranged in an overall square pattern (QR Code).
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

// BCOptions are parameters used to configure the BC function call
type BCOptions struct {
	// orientation sets the field orientation for the barcode.
	// Accepted values:
	//   - OrientationNormal (N)
	//   - OrientationRotated (R)
	//   - OrientationInverted (I)
	//   - OrientationBottomUp (B)
	// Default is the current ^FW setting or the printer's default orientation.
	orientation Orientation

	// height defines the height of the barcode in dots.
	// Valid range: 1 to 32000.
	// If nil, the height is determined by the ^BY setting or printer default.
	height *int

	// showTextBelow controls whether to print the interpretation line (human-readable text) below the barcode.
	// Set to true to show the text below.
	// Default: true.
	showTextBelow bool

	// showTextAbove controls whether to print the interpretation line above the barcode.
	// Set to true to show the text above.
	// Default: false.
	showTextAbove bool

	// addUCCCheckDigit enables or disables the UCC Mod 10 check digit.
	// When true, both Mod 10 and Mod 103 check digits are included.
	// Default: false.
	addUCCCheckDigit bool

	// mode sets the encoding mode for the barcode.
	// Accepted values (BCMode):
	//   - NONE_MODE (N): No special mode
	//   - UCC_MODE (U): UCC Case Mode (pads/truncates ^FD to 19 digits)
	//   - AUTO_MODE (A): Automatic mode for optimal subset packing
	//   - EAN_MODE (D): UCC/EAN Mode (adds FNC1 and encodes using application identifiers)
	// Default: NONE_MODE.
	mode BCMode
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

// BC writes the ZPL command (^BC) to [*ZPL] internal bytes.Buffer and produces a code 128 barcode.
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
