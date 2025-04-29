package zpl

// XA represents the ZPL Start Format command (^XA).
// This command must begin every ZPL label format and is automatically added
// when creating a new label with NewLabel().
//
// The ^XA command:
//   - Initializes the printer for a new label format
//   - Clears any previous formatting information
//   - Must be paired with a terminating ^XZ command
//   - No parameters are available for this command
type xa struct{}

// XA creates a new Start Format command.
// This command is used internally by the NewLabel() constructor
// rather than being called directly.
//
// Example:
//
//	label := zpl.NewLabel()  // Automatically includes ^XA
//	// ... add other commands ...
func XA() *xa {
	return &xa{}
}

// ToZPL converts the XA command to its ZPL string representation.
// Returns the literal "^XA" command string.
// This implementation always returns the same constant value since
// the command has no parameters or variations.
func (xa *xa) ToZPL() string {
	return "^XA"
}
