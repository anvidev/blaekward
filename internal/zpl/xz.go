package zpl

type xz struct{}

func XZ() *xz {
	return &xz{}
}

func (xz *xz) ToZPL() string {
	return "^XZ"
}
