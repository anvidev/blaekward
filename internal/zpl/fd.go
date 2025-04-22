package zpl

import "fmt"

type fd struct {
	data string
}

func FD(data string) *fd {
	return &fd{data}
}

func (fd *fd) ToZPL() string {
	return fmt.Sprintf("^FD%s", fd.data)
}
