package zpl

// The ^XZ command:
//   - Terminates the current label format
//   - Instructs the printer to process the label
//   - Must be paired with a starting ^XA command
//   - Can only appear once per label format
//   - Always appears as the last command in a ZPL format
//   - No parameters are available for this command
type xz struct{}

// XZ creates a new End Format command.
// This command is typically added automatically when calling Label.String()
// but can be added manually if needed.
//
// Example:
//
//	label := zpl.NewLabel()  // Starts with ^XA
//	// ... add other commands ...
//	label.AddCommand(zpl.XZ())  // Explicitly end format (not recommended)
//	// OR
//	zplStr := label.String()    // Automatically adds ^XZ if missing
func XZ() *xz {
	return &xz{}
}

// ToZPL converts the XZ command to its ZPL string representation.
// Returns the literal "^XZ" command string.
// This implementation always returns the same constant value since
// the command has no parameters or variations.
func (xz *xz) ToZPL() string {
	return "^XZ"
}
