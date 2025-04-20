package zpl

import "fmt"

type ACommand struct {
	font             string
	fieldOrientation *Orientation
	charHeight       *int
	charWidth        *int
}

func (z *ZPL) A(font string, opts ...OptionFunc[ACommand]) {
	cfg := &ACommand{font: font}

	for _, opt := range opts {
		opt(cfg)
	}

	cmd := "^A"

	if isValidFont(cfg.font) {
		cmd += fmt.Sprintf("%s", cfg.font)
	} else {
		cmd += fmt.Sprintf("%s", "A")
	}

	if cfg.fieldOrientation != nil {
		cmd += fmt.Sprintf(",%s", *cfg.fieldOrientation)
	}

	if cfg.charHeight != nil {
		cmd += fmt.Sprintf(",%d", *cfg.charHeight)

		// width can only be applied if height is set
		if cfg.charWidth != nil {
			cmd += fmt.Sprintf(",%d", *cfg.charWidth)
		}
	}

	z.writeCmd(cmd)
}

func WithFieldOrientation(v Orientation) OptionFunc[ACommand] {
	return func(a *ACommand) { a.fieldOrientation = &v }
}

func WithCharacterWidth(v int) OptionFunc[ACommand] {
	return func(a *ACommand) { a.charWidth = &v }
}

func WithCharacterHeight(v int) OptionFunc[ACommand] {
	return func(a *ACommand) { a.charHeight = &v }
}

func isValidFont(font string) bool {
	if len(font) != 1 {
		return false
	}
	r := rune(font[0])
	return (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9')
}
