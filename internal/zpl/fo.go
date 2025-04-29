package zpl

import "fmt"

// The ^FO command:
//   - x: Horizontal position in dots from left edge (required)
//   - y: Vertical position in dots from top edge (required)
//   - Sets the reference point for subsequent field elements
//   - Coordinates are absolute positions (0,0 = top-left corner)
//   - Affects all following fields until another ^FO is encountered
//   - Must be specified before any field data or graphic elements
//
// Note: The position is affected by any previous ^LH (Label Home) settings.
type fo struct {
	x uint // Horizontal position in dots (left edge = 0)
	y uint // Vertical position in dots (top edge = 0)
}

// FO creates a new Field Origin command with specified coordinates.
// Both x and y positions are required parameters.
//
// Example:
//
//	// Position text at 100 dots from left, 200 dots from top
//	label.AddCommand(zpl.FO(100, 200))
//	label.AddCommand(zpl.FD("Positioned Text"))
//	label.AddCommand(zpl.FS())
//
//	// Position a box at different coordinates
//	label.AddCommand(zpl.FO(50, 50))
//	label.AddCommand(zpl.GB(100, 100))
func FO(x, y uint) *fo {
	return &fo{x, y}
}

// ToZPL converts the FO command to its ZPL string representation.
// Returns the formatted ^FO command with x,y coordinates.
// Example output: "^FO100,200" for position (100,200)
func (fo *fo) ToZPL() string {
	return fmt.Sprintf("^FO%d,%d", fo.x, fo.y)
}
