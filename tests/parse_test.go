package tests

import (
	"errors"
	"fmt"
	"testing"

	"401/automata"
	"401/engine"
	"401/oracle"
	"401/regex"
)

// Examples for Project 2. They fail with "unimplemented: phase 2.1 Parse"
// until you write Parse. Run just these with go test ./tests/ -run TestP2.
//
// Like the smoke tests, they wrap each case in oracle.Case, so an
// unimplemented Parse doesn't stop the rest of this package's tests. Your own
// tests shouldn't (see README.md).
//
// This file declares nothing but its four TestP2_* functions, so those four
// names are the only ones your own test files in package tests must avoid.

// TestP2_1_Examples parses a few patterns and checks each against an AST
// built by hand. The check is exact rather than sampled: compile both ASTs
// over one alphabet, determinize, and decide equivalence. A parse is right
// when it denotes the right language, whatever shape its tree has.
func TestP2_1_Examples(t *testing.T) {
	sigma := []automata.Symbol{'a', 'b', 'c'}
	for _, c := range []struct {
		pattern string
		want    func() regex.Node // a function, so a panicking constructor stays inside oracle.Case
	}{
		{"a", func() regex.Node { return regex.Lit('a') }},
		{"a|bc", func() regex.Node {
			return regex.Or(regex.Lit('a'), regex.Seq(regex.Lit('b'), regex.Lit('c')))
		}},
		{"ab*", func() regex.Node { return regex.Seq(regex.Lit('a'), regex.Star(regex.Lit('b'))) }},
		{"(a|b)*abb", func() regex.Node {
			return regex.Seq(regex.Star(regex.Or(regex.Lit('a'), regex.Lit('b'))),
				regex.Lit('a'), regex.Lit('b'), regex.Lit('b'))
		}},
		{"ε|a", func() regex.Node { return regex.Or(regex.Eps(), regex.Lit('a')) }},
		{"a∅|b", func() regex.Node {
			return regex.Or(regex.Seq(regex.Lit('a'), regex.Emptyset()), regex.Lit('b'))
		}},
		{"(a*)*", func() regex.Node { return regex.Star(regex.Star(regex.Lit('a'))) }},
	} {
		oracle.Case(t, c.pattern, func(t *testing.T) {
			got, err := regex.Parse(c.pattern)
			if err != nil {
				t.Fatalf("Parse(%q) returned an error: %v", c.pattern, err)
			}
			oracle.AssertEquivalent(t,
				automata.NFAToDFA(regex.Compile(got, sigma)),
				automata.NFAToDFA(regex.Compile(c.want(), sigma)))
		})
	}
}

// TestP2_1_Rejections checks the examples in Parse's doc comment, which
// defines Offset: each must be rejected with a *SyntaxError at that offset.
func TestP2_1_Rejections(t *testing.T) {
	for _, c := range []struct {
		pattern string
		offset  int
	}{
		{"", 0},      // ends too early
		{"a|", 2},    // ends too early
		{"(a", 2},    // ends too early: the ) is missing
		{"*a", 0},    // no pattern starts with *
		{"()", 1},    // a group can't be empty
		{"a||b", 2},  // an alternative can't be empty
		{"a**", 2},   // a nested star needs parentheses
		{"a b", 1},   // a space isn't a token
		{"*a b", 0},  // the first error counts, not the first bad byte
		{"aεb)", 4},  // ε is two bytes
		{"a\xce", 1}, // half an ε isn't a whole token
	} {
		oracle.Case(t, fmt.Sprintf("%q", c.pattern), func(t *testing.T) {
			_, err := regex.Parse(c.pattern)
			if err == nil {
				t.Fatalf("Parse(%q) accepted it", c.pattern)
			}
			var se *regex.SyntaxError
			if !errors.As(err, &se) {
				t.Fatalf("Parse(%q) returned %q, which isn't a *regex.SyntaxError", c.pattern, err)
			}
			if se.Offset != c.offset {
				t.Errorf("Parse(%q): Offset is %d, want %d", c.pattern, se.Offset, c.offset)
			}
		})
	}
}

// TestP2_1_RoundTrip parses what your String() renders, for random ASTs, and
// checks that the parse denotes the AST's language.
//
// It tests String() as much as Parse. Phase 1.5 checked String() only through
// Go's regexp, which accepts some things course syntax doesn't, such as
// (?:a|b), (), and the empty alternative in (a|). If your String() emits one,
// your Parse is right to reject it, and the fix is in String().
func TestP2_1_RoundTrip(t *testing.T) {
	sigma := []automata.Symbol{'a', 'b'}
	g := oracle.NewGen(1, sigma)
	oracle.Case(t, "random ASTs", func(t *testing.T) {
		for range 200 {
			n := g.AST(5)
			rendered := n.String()
			back, err := regex.Parse(rendered)
			if err != nil {
				t.Fatalf("Parse rejected %q, which your String() rendered: %v", rendered, err)
			}
			oracle.AssertEquivalent(t,
				automata.NFAToDFA(regex.Compile(back, sigma)),
				automata.NFAToDFA(regex.Compile(n, sigma)))
			if t.Failed() {
				t.Fatalf("Parse(%q) doesn't denote the language of the AST it was rendered from", rendered)
			}
		}
	})
}

// TestP2_2_Examples runs patterns through the whole pipeline: engine.Compile,
// then Contains.
func TestP2_2_Examples(t *testing.T) {
	oracle.Case(t, "(a|b)*abb", func(t *testing.T) {
		r, err := engine.Compile("(a|b)*abb")
		if err != nil {
			t.Fatalf("Compile returned an error: %v", err)
		}
		oracle.Accepts(t, r, "abb")
		oracle.Accepts(t, r, "babb")
		oracle.Rejects(t, r, "abba")
		oracle.Rejects(t, r, "")
	})

	oracle.Case(t, "a||b", func(t *testing.T) {
		_, err := engine.Compile("a||b")
		var se *regex.SyntaxError
		if !errors.As(err, &se) {
			t.Fatalf("Compile(%q) returned %v; errors.As should find Parse's *regex.SyntaxError", "a||b", err)
		}
		if se.Offset != 2 {
			t.Errorf("Offset is %d, want 2", se.Offset)
		}
	})
}
