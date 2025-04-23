package zpl

import "fmt"

// BC represents the ZPL Code 128 Barcode command (^BC).
// This command creates a Code 128 barcode field with customizable parameters.
//
// The ^BC command:
//   - Generates a Code 128 barcode (UCC/EAN-128, Auto, or Standard)
//   - Supports human-readable text above/below the barcode
//   - Includes optional UCC check digit calculation
//   - Must be paired with ^FD (Field Data) and terminated with ^FS
//
// Technical Reference:
// Format: ^BCo,h,f,g,m
//
// Parameters:
//
//   - o = Orientation (N-normal, R-rotated 90°, I-inverted 180°, B-bottom-up 270°)
//   - h = Barcode height in dots (1-32000)
//   - f = Print interpretation line below barcode (Y/N)
//   - g = Print interpretation line above barcode (Y/N)
//   - m = Mode (A-auto, D-UCC/EAN, U-standard)
type bc struct {
	orientation   string
	barcodeHeight uint
	textBelow     *bool
	textAbove     *bool
	uccCheckDigit *bool
	mode          string
}

// BC creates a new Code 128 Barcode command.
//
// Example:
//
//	// Basic barcode with default settings
//	zpl.BC().SetBarcodeHeight(50)
//
//	// UCC/EAN-128 barcode with check digit
//	zpl.BC().SetMode("D").SetUCCCheckDigit(true)
func BC() *bc {
	return &bc{}
}

// ToZPL converts the BC command to its ZPL string representation.
// Returns the properly formatted ^BC command with specified parameters.
// Example outputs:
//
//	"^BCN,50,,,,"      // Normal orientation, 50-dot height
//	"^BC,,,Y,,A"       // Defaults with text above
//	"^BCR,30,Y,N,Y,D"  // Rotated UCC barcode with check digit
func (bc *bc) ToZPL() string {
	cmd := "^BC"

	if bc.orientation != "" {
		cmd += fmt.Sprintf("%s,", bc.orientation)
	} else {
		cmd += ","
	}

	if bc.barcodeHeight > 0 {
		cmd += fmt.Sprintf("%d,", bc.barcodeHeight)
	} else {
		cmd += ","
	}

	if bc.textBelow != nil {
		cmd += fmt.Sprintf("%s,", boolToZPL(*bc.textBelow))
	} else {
		cmd += ","
	}

	if bc.textAbove != nil {
		cmd += fmt.Sprintf("%s,", boolToZPL(*bc.textAbove))
	} else {
		cmd += ","
	}

	if bc.uccCheckDigit != nil {
		cmd += fmt.Sprintf("%s,", boolToZPL(*bc.uccCheckDigit))
	} else {
		cmd += ","
	}

	if bc.mode != "" {
		cmd += fmt.Sprintf("%s", bc.mode)
	}

	return cmd
}

// SetOrientation sets the barcode orientation.
// Valid values: N (normal), R (90°), I (180°), B (270°)
func (bc *bc) SetOrientation(orientation string) *bc {
	bc.orientation = orientation
	return bc
}

// SetBarcodeHeight sets the barcode height in dots (1-32000).
func (bc *bc) SetBarcodeHeight(height uint) *bc {
	bc.barcodeHeight = height
	return bc
}

// SetTextBelow controls printing of human-readable text below barcode.
func (bc *bc) SetTextBelow(on bool) *bc {
	bc.textBelow = &on
	return bc
}

// SetTextAbove controls printing of human-readable text above barcode.
func (bc *bc) SetTextAbove(on bool) *bc {
	bc.textAbove = &on
	return bc
}

// SetUCCCheckDigit enables/disables UCC check digit calculation.
func (bc *bc) SetUCCCheckDigit(on bool) *bc {
	bc.uccCheckDigit = &on
	return bc
}

// SetMode sets the Code 128 encoding mode.
// Valid values: A (auto), D (UCC/EAN-128), U (standard)
func (bc *bc) SetMode(mode string) *bc {
	bc.mode = mode
	return bc
}
