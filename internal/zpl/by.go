package zpl

import "fmt"

// BY represents the ZPL Barcode Field Default command (^BY).
// This command sets the default parameters for all subsequent barcode fields.
//
// The ^BY command:
//   - Defines default module width, width ratio, and height for barcodes
//   - Affects all barcodes until changed by another ^BY command
//   - Must be specified before each ^B (barcode) command
//
// Technical Reference:
// Format: ^BYw,r,h
// Parameters:
//
//   - w = Module width in dots (1-10), default 2
//   - r = Wide/narrow ratio (2.0-3.0), default 3.0
//   - h = Barcode height in dots (1-32000), default 10
type by struct {
	moduleWidth   uint
	barRatio      float32
	barcodeHeight uint
}

// BY creates a new Barcode Field Default command with default values (2,3.0,10).
// Use builder methods to customize the barcode parameters.
//
// Example:
//
//	// Set narrow barcode with 1:2 ratio
//	zpl.BY().SetModuleWidth(1).SetBarRatio(2.0)
//
//	// Complete barcode specification
//	zpl.BY().SetModuleWidth(3).SetBarRatio(2.5).SetBarcodeHeight(50)
func BY() *by {
	return &by{}
}

// ToZPL converts the BY command to its ZPL string representation.
// Returns the properly formatted ^BY command with specified parameters.
// Example outputs:
//
//	"^BY2,,"       // Default module width only
//	"^BY,2.5,"     // Custom ratio only
//	"^BY3,2.5,50"  // Complete specification
func (by *by) ToZPL() string {
	cmd := "^BY"

	if by.moduleWidth > 0 {
		cmd += fmt.Sprintf("%d,", by.moduleWidth)
	} else {
		cmd += ","
	}

	if by.barRatio > 0 {
		cmd += fmt.Sprintf("%.1f,", by.barRatio)
	} else {
		cmd += ","
	}

	if by.barcodeHeight > 0 {
		cmd += fmt.Sprintf("%d", by.barcodeHeight)
	}

	return cmd
}

// SetModuleWidth sets the barcode module width in dots (1-10).
func (by *by) SetModuleWidth(width uint) *by {
	by.moduleWidth = width
	return by
}

// SetBarRatio sets the wide/narrow bar ratio (2.0-3.0).
func (by *by) SetBarRatio(ratio float32) *by {
	by.barRatio = ratio
	return by
}

// SetBarcodeHeight sets the barcode height in dots (1-32000).
func (by *by) SetBarcodeHeight(height uint) *by {
	by.barcodeHeight = height
	return by
}
