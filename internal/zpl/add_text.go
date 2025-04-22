package zpl

import "fmt"

type addTextCommand struct {
	x, y uint
	data string
}

func AddText(x, y uint, data string) *addTextCommand {
	return &addTextCommand{
		x:    x,
		y:    y,
		data: data,
	}
}

func (at *addTextCommand) ToZPL() string {
	fo := FO(at.x, at.y)
	fd := FD(at.data)
	return fmt.Sprintf("%s%s", fo.ToZPL(), fd.ToZPL())
}
