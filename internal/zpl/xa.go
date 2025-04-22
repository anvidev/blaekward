package zpl

type xa struct{}

func XA() *xa {
	return &xa{}
}

func (xa *xa) ToZPL() string {
	return "^XA"
}
