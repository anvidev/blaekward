package zpl

func (z *zpl) XA() {
	z.writeCmd("^XA")
}
