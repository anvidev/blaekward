package zpl

import (
	"fmt"
)

func (z *zpl) FD(content string) {
	cmd := fmt.Sprintf("^FD%s", content)
	z.writeCmd(cmd)
}
