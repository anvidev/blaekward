package label

import "github.com/anvidev/blaekward/internal/zpl"

func NewTest(name, street, city, country, barcode string) []byte {
	z := zpl.New(76, 51)

	z.CF(zpl.WithCFFont("A"), zpl.WithCFCharHeight(30))
	z.FO(30, 30)
	z.FD(name)
	z.FS()
	z.FO(30, 65)
	z.FD(street)
	z.FS()
	z.FO(30, 100)
	z.FD(city)
	z.FS()
	z.FO(30, 135)
	z.FD(country)
	z.FS()
	z.BY(zpl.WithModuleWidth(3), zpl.WithBarcodeHeight(100))
	z.FO(30, 200)
	z.BC()
	z.FD(barcode)
	z.FS()

	return z.Bytes()
}
