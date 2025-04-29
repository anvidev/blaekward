package zpl

import "fmt"

// The ^GB command parameters:
//   - width: Box width in dots (required)
//   - height: Box height in dots (required)
//   - borderWidth: Line thickness in dots (optional, default: 1)
//   - lineColor: "B" for black or "W" for white (optional, default: "B")
//   - borderRadius: Corner radius in dots (optional, default: 0 - square corners)
//   - Parameters must be specified in order: width, height, borderWidth, lineColor, borderRadius
//   - Empty commas maintain parameter positions for optional values
type gb struct {
	width        uint   // Box width in dots (required)
	height       uint   // Box height in dots (required)
	borderWidth  uint   // Line thickness in dots (0 = default thickness)
	borderRadius uint   // Corner radius in dots (0 = square corners)
	lineColor    string // "B" (black) or "W" (white)
}

// GB creates a new Graphic Box command with required width and height.
// Default values will be used for unspecified parameters:
//   - borderWidth: 1 dot
//   - lineColor: "B" (black)
//   - borderRadius: 0 (square corners)
//
// Example:
//
//	// Create a 300x200 dot box with default settings
//	box := zpl.GB(300, 200)
//
//	// Create a customized box
//	box := zpl.GB(300, 200).
//	    SetBorderWidth(3).
//	    SetLineColor("W").
//	    SetBorderRadius(10)
func GB(width, height uint) *gb {
	return &gb{
		width:  width,
		height: height,
	}
}

// ToZPL converts the GB command to its ZPL string representation.
// Returns the properly formatted ^GB command with all specified parameters.
// Example outputs:
//
//	"^GB300,200,,,"      // Basic box
//	"^GB300,200,2,B,10"  // Customized box
func (gb *gb) ToZPL() string {
	cmd := fmt.Sprintf("^GB%d,%d,", gb.width, gb.height)

	if gb.borderWidth > 0 {
		cmd += fmt.Sprintf("%d,", gb.borderWidth)
	} else {
		cmd += ","
	}

	if gb.lineColor != "" {
		cmd += fmt.Sprintf("%s,", gb.lineColor)
	} else {
		cmd += ","
	}

	if gb.borderRadius > 0 {
		cmd += fmt.Sprintf("%d", gb.borderRadius)
	}
	return cmd
}

// SetBorderWidth sets the line thickness for the box border.
// A value of 0 will use the printer's default thickness (1 dot).
func (gb *gb) SetBorderWidth(width uint) *gb {
	gb.borderWidth = width
	return gb
}

// SetBorderRadius sets the corner radius for rounded corners.
// A value of 0 produces square corners.
func (gb *gb) SetBorderRadius(radius uint) *gb {
	gb.borderRadius = radius
	return gb
}

// SetLineColor sets the line color for the box.
// Valid values: "B" (black) or "W" (white).
func (gb *gb) SetLineColor(color string) *gb {
	gb.lineColor = color
	return gb
}
