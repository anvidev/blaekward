package zpl

// The ^FS command:
// - Marks the end of a field definition
// - Resets field parameters to their defaults
// - No parameters are available for this command
// - Must appear after every ^FD (Field Data) command
type fs struct{}

// FS creates a new Field Separator command.
// This command should be used after each complete field definition.
//
// Example:
//
//	label := zpl.NewLabel()
//	label.AddCommand(zpl.FD("Sample Text")) // Field data
//	label.AddCommand(zpl.FS())              // Terminate field
//	label.AddCommand(zpl.FD("More Text"))   // New field
//	label.AddCommand(zpl.FS())              // Terminate field
func FS() *fs {
	return &fs{}
}

// ToZPL converts the FS command to its ZPL string representation.
// Returns the literal "^FS" command string.
// This implementation always returns the same constant value since
// the command has no parameters or variations.
func (fs *fs) ToZPL() string {
	return "^FS"
}
