package zpl

type fs struct{}

func FS() *fs {
	return &fs{}
}

func (fs *fs) ToZPL() string {
	return "^FS"
}
