package regex

// Phase 2.1: the parser.

// Parse reads a pattern in course regex syntax, the syntax String renders,
// and returns an AST for it. The syntax is this grammar, with Union as the
// start symbol:
//
//	Union   → Union '|' Concat | Concat
//	Concat  → Concat Star | Star
//	Star    → Atom '*' | Atom
//	Atom    → LetterOrDigit | 'ε' | '∅' | '(' Union ')'
//
// The tokens are an ASCII letter or digit (case-sensitive), ε (U+03B5, two
// bytes in UTF-8), ∅ (U+2205, three bytes), and the four one-byte tokens
// | * ( and ). Anything else is an error, whitespace and other non-ASCII
// characters included. So are the empty pattern, an empty group "()", an
// empty alternative "a|", and a doubled star "a**": a nested star is
// written (a*)*.
//
// Parse accepts every string the grammar generates and rejects every other
// string. It never panics, whatever the input, invalid UTF-8 included.
//
// For an accepted pattern, the error is nil and the Node is not. The Node
// denotes the pattern's language, with star binding tighter than
// concatenation and concatenation tighter than union. Its shape is yours:
// a|b|c may come out as Or(a, b, c), Or(Or(a, b), c), or Or(a, Or(b, c)).
// It mentions no literal the pattern doesn't, and like any Node it can be
// passed to String, Compile, and engine.FromAST.
//
// For a rejected pattern, the error is a *SyntaxError, and nothing is
// promised about the Node. Offset is the length of the longest prefix of
// pattern that is a sequence of whole tokens and that some valid pattern
// begins with. That is the byte offset where the first token that can't
// continue a valid pattern begins, or where a byte sits that doesn't begin a
// whole token, or len(pattern) if the pattern ends too early:
//
//	pattern  Offset
//	""       0       ends too early
//	"a|"     2       ends too early
//	"(a"     2       ends too early
//	"*a"     0       no pattern starts with *
//	"()"     1       a group can't be empty
//	"a||b"   2       an alternative can't be empty
//	"a**"    2       a nested star needs parentheses
//	"a b"    1       a space isn't a token
//	"*a b"   0       the first error counts, not the first bad byte
//	"aεb)"   4       ε is two bytes
//	"a\xce"  1       half an ε isn't a whole token
//
// The wording of Msg is yours.
//
// Write the parser yourself, top-down or bottom-up, as your phase 2.0
// writeup plans it. README-P2.md says what that rules out.
func Parse(pattern string) (Node, error) {
	panic("unimplemented: phase 2.1 Parse")
}
