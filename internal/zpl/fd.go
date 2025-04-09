package zpl

import (
	"fmt"
)

func (z *ZPL) FD(content string) {
	cmd := fmt.Sprintf("^FD%s", content)
	z.writeCmd(cmd)
}
