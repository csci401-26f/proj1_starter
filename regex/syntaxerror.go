package regex

import "fmt"

// Phase 2.1: the error Parse returns for a rejected pattern. This file ships
// complete, and grading replaces it with the shipped copy, so don't edit it.
// Anything of your own goes in another file.

// SyntaxError reports a pattern that isn't valid course regex syntax. Parse
// and engine.Compile return one for every pattern they reject, as is or
// wrapped with %w; find it with errors.As.
type SyntaxError struct {
	Offset int    // byte offset of the error, as Parse's documentation defines it
	Msg    string // what was expected at Offset
}

// Error returns the offset and the message.
func (e *SyntaxError) Error() string { return fmt.Sprintf("offset %d: %s", e.Offset, e.Msg) }
