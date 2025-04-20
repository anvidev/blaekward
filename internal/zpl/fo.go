package zpl

import "fmt"

type FOCommand struct {
	x             int
	y             int
	justification *Justification // only supported on firmware versions v60.14.x, v50.14.x or later
}

func (z *ZPL) FO(x, y int, opts ...OptionFunc[FOCommand]) {
	cfg := &FOCommand{
		x: x,
		y: y,
	}

	for _, opt := range opts {
		opt(cfg)
	}

	cmd := fmt.Sprintf("^FO%d,%d", x, y)

	if cfg.justification != nil {
		cmd += fmt.Sprintf(",%d", *cfg.justification)
	}

	z.writeCmd(cmd)
}

func WithJustify(v Justification) OptionFunc[FOCommand] {
	return func(f *FOCommand) { f.justification = &v }
}
