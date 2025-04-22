package zpl

import "fmt"

type fo struct {
	x, y uint
}

func FO(x, y uint) *fo {
	return &fo{x, y}
}

func (fo *fo) ToZPL() string {
	return fmt.Sprintf("^FO%d,%d", fo.x, fo.y)
}
