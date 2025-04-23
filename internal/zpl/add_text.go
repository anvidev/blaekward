package zpl

import "fmt"

// AddTextCommand represents a composed text field command with position and content.
// This higher-order command combines FO (Field Origin), FD (Field Data) and FS (Field Separator)
// commands into a single operation for convenience.
type addTextCommand struct {
	x, y uint
	data string
}

// AddText creates a new positioned text field command.
// This helper combines:
//   - ^FO (Field Origin) for positioning
//   - ^FD (Field Data) for text content
//   - ^FS (Field Separator) for marking end of field definition
//
// Example:
//
//	// Add "Hello" at position (100,200)
//	label.AddCommand(zpl.AddText(100, 200, "Hello"))
func AddText(x, y uint, data string) *addTextCommand {
	return &addTextCommand{
		x:    x,
		y:    y,
		data: data,
	}
}

// ToZPL converts the command to its ZPL representation.
// Returns the combined ^FO and ^FD commands with termination (^FS).
// Example output: "^FO100,200^FDHello^FS"
func (at *addTextCommand) ToZPL() string {
	fo := FO(at.x, at.y)
	fd := FD(at.data)
	fs := FS()
	return fmt.Sprintf("%s%s%s", fo.ToZPL(), fd.ToZPL(), fs.ToZPL())
}
