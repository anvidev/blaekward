package zpl

func (z *ZPL) FS() {
	z.writeCmd("^FS")
}
