package engine

// Phase 2.2: the full pipeline, from a pattern string to a matcher.

// Compile parses pattern with regex.Parse and builds a Regexp from the result
// with FromAST. It is not regex.Compile, which is Thompson's construction
// from an AST (phase 1.6).
//
// Compile accepts exactly the patterns Parse accepts. For those, the error
// is nil and the Regexp is not, and the Regexp matches exactly what
// FromAST's would for the parsed Node: the pattern's language, whole-word.
// For the rest, errors.As finds the *regex.SyntaxError that Parse returned,
// with the same Offset, and nothing is promised about the Regexp. Like
// Parse, Compile never panics.
func Compile(pattern string) (*Regexp, error) {
	panic("unimplemented: phase 2.2 Compile")
}
