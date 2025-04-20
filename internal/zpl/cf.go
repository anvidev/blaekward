package zpl

import "fmt"

type cfCommand struct {
	font       *string
	charHeight *int
	charWidth  *int
}

func (z *ZPL) CF(opts ...OptionFunc[cfCommand]) {
	cfg := &cfCommand{}

	for _, opt := range opts {
		opt(cfg)
	}

	cmd := "^CF"

	if cfg.font != nil {
		cmd += fmt.Sprintf("%s,", *cfg.font)
	} else {
		cmd += ","
	}

	if cfg.charHeight != nil {
		cmd += fmt.Sprintf("%d,", *cfg.charHeight)
	} else {
		cmd += ","
	}

	if cfg.charWidth != nil {
		cmd += fmt.Sprintf("%d", *cfg.charWidth)
	}

	z.writeCmd(cmd)
}

func WithCFFont(v string) OptionFunc[cfCommand] {
	return func(cc *cfCommand) { cc.font = &v }
}

func WithCFCharHeight(v int) OptionFunc[cfCommand] {
	return func(cc *cfCommand) { cc.charHeight = &v }
}

func WithCFCharWidth(v int) OptionFunc[cfCommand] {
	return func(cc *cfCommand) { cc.charWidth = &v }
}
