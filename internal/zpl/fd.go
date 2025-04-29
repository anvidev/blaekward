package zpl

import "fmt"

// The ^FD command:
// - Contains the printable data for a text field
// - Must be preceded by field formatting commands (^A, ^FO, etc.)
// - Must be terminated with ^FS (Field Separator)
// - Maximum length for data is 3072 bytes
// - Character encoding depends on current ^CI (Character Set) setting
// - Field orientation affects rendering (^FW, ^FR commands)
// - Always requires terminating ^FS command
type fd struct {
	data string
}

// FD creates a new Field Data command with the specified content.
// The data string will be rendered according to current field settings.
//
// Example:
//
//	label.AddCommand(zpl.FO(100, 100))    // Set position
//	label.AddCommand(zpl.A(50, 50))       // Set font
//	label.AddCommand(zpl.FD("Hello World")) // Set text
//	label.AddCommand(zpl.FS())             // Terminate field
func FD(data string) *fd {
	return &fd{data}
}

// ToZPL converts the FD command to its ZPL string representation.
// Returns the formatted ^FD command with the contained data.
// Example output: "^FDHello World^FS"
func (fd *fd) ToZPL() string {
	return fmt.Sprintf("^FD%s", fd.data)
}
