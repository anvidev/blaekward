package zpl

import (
	"fmt"
)

// CF represents the ZPL Change Font command (^CF).
// This command sets the default font and character dimensions for subsequent text fields.
//
// Format: ^CFf,h,w
// Parameters:
//   - f: Font identifier (A-Z, 0-9)
//   - h: Character height in dots (0-32000)
//   - w: Character width in dots (0-32000)
//
// Technical Reference:
//   - Default values at power-up: Font A, height 9, width 5
//   - Affects all subsequent text fields until changed
//   - Proportional scaling occurs when only height or width is specified
//   - Invalid font values default to font A
//
// Official ZPL Documentation:
// The ^CF command simplifies programs by setting default font characteristics.
// If no alphanumeric field command (^AF) is used or an invalid font is specified,
// all data prints using font A by default.
type cf struct {
	font   string // Font identifier (A-Z, 0-9)
	height uint   // Character height in dots (0-32000)
	width  uint   // Character width in dots (0-32000)
}

// CF creates a new Change Font command with the specified content.
// Use the builder methods to customize the font settings.
//
// Example:
//
//	// Set default font to D with 24-dot height (auto-width)
//	label.AddCommand(zpl.CF().
//	    SetFont("D").
//	    SetCharacterHeight(24))
//
//	// Set complete font specification
//	label.AddCommand(zpl.CF().
//	    SetFont("B").
//	    SetCharacterHeight(18).
//	    SetCharacterWidth(12))
func CF() *cf {
	return &cf{}
}

// ToZPL converts the FD command to its ZPL string representation.
// Returns the formatted ^FD command with the contained data.
// Example output: "^FDHello World^FS"
func (cf *cf) ToZPL() string {
	cmd := "^CF"

	if cf.font != "" {
		cmd += fmt.Sprintf("%s,", cf.font)
	} else {
		cmd += ","
	}

	if cf.height > 0 {
		cmd += fmt.Sprintf("%d,", cf.height)
	} else {
		cmd += ","
	}

	if cf.width > 0 {
		cmd += fmt.Sprintf("%d", cf.width)
	}

	return cmd
}

// SetFont specifies the default font (A-Z, 0-9).
// Invalid values default to "A".
func (cf *cf) SetFont(font string) *cf {
	cf.font = font
	return cf
}

// SetCharacterHeight sets the character height in dots (0-32000).
// Value 0 uses printer default (9).
func (cf *cf) SetCharacterHeight(height uint) *cf {
	cf.height = height
	return cf
}

// SetCharacterWidth sets the character width in dots (0-32000).
// Value 0 maintains proportional scaling.
func (cf *cf) SetCharacterWidth(width uint) *cf {
	cf.width = width
	return cf
}
