// Package zpl2 provides a builder-style API for generating ZPL (Zebra Programming Language) labels.
// It enables type-safe, composable construction of ZPL commands that can be sent to Zebra printers.
package zpl

import (
	"strings"
)

// Commander is the core interface that all ZPL commands must implement.
// It represents any command that can be converted to ZPL format.
type Commander interface {
	// ToZPL converts the command to its ZPL string representation
	ToZPL() string
	// Validate() error // future validation support
}

// Label represents a complete ZPL label with multiple commands.
// It manages the sequence of commands and compiles them to valid ZPL.
type Label struct {
	commands []Commander
	// errors   []error // future error tracking for validation fails
	// unit string // future unit convertion
	// dpi uint // future unit convertion
}

// NewLabel creates a new Label with default initialization.
// The label is pre-configured with a starting XA (Start Format) command,
// which is required at the beginning of every ZPL label.
func NewLabel() *Label {
	l := &Label{commands: []Commander{}}
	l.AddCommand(XA())
	return l
}

// AddCommand appends a new command to the label's command sequence.
// Commands are executed by the printer in the order they are added.
func (l *Label) AddCommand(cmd Commander) {
	l.commands = append(l.commands, cmd)
}

// String compiles all commands into a complete ZPL label string.
// It automatically adds the XZ (End Format) command as the last command.
// The returned string is ready to be sent to a Zebra printer.
func (l *Label) String() string {
	var sb strings.Builder
	l.AddCommand(XZ())
	for _, cmd := range l.commands {
		sb.WriteString(cmd.ToZPL())
	}
	return sb.String()
}

// boolToZPL converts a boolean value to ZPL's boolean representation.
// Returns "Y" for true and "N" for false, which is the standard in ZPL.
// This is an internal helper function not exposed outside the package.
func boolToZPL(b bool) string {
	return map[bool]string{true: "Y", false: "N"}[b]
}
