# Project 2: A Regex Parser

This file is Project 2's technical specification. Everything in `README.md` still holds, while this file adds what's new for Project 2. The project handout on the course website covers the schedule, the grammar analysis, deliverables, reflections, and grading.

As in Project 1, the specification has two parts: this file, and the doc comments of the new contract declarations. Read `go doc ./regex Parse`, `go doc ./regex SyntaxError`, and `go doc ./engine Compile`. `Parse`'s doc comment gives the grammar and the exact definition of an error's offset.

## New files for Project 2

Project 2 arrives as a pull request that only adds files. Nothing that shipped with Project 1 changed.

| File | Ships | What it is |
|---|---|---|
| `regex/parse.go` | stub | `Parse`, phase 2.1 |
| `regex/syntaxerror.go` | complete | `SyntaxError`, the error `Parse` returns |
| `engine/compile.go` | stub | `Compile`, phase 2.2 |
| `tests/contract_p2_test.go` | complete | compile-time signature checks for the three |
| `tests/parse_test.go` | complete | a few examples, including a round trip through your `String()` |
| `README-P2.md` | complete | this file |

## Rules

**1. Leave these alone,** in addition to everything on `README.md`'s list:

- **`regex/syntaxerror.go`.** Grading replaces it with the shipped copy, like `corpus/` and `oracle/`.
- **`tests/contract_p2_test.go`.** Like `tests/contract_test.go`, it checks the new signatures at compile time.
- **The contract declarations in the new stub files.**

**2. Write the parser yourself:** outside your tests, you should not use a parser generator (such as goyacc), a parser-combinator library, a PEG library, or Go's `regexp` or `regexp/syntax` packages. Your tests may use Go's `regexp` as an oracle.

## Workflow

```sh
go test ./tests/ -run TestP2         # the shipped examples
go test ./... -run 'TestP2_1'        # phase 2.1, if you name your tests by phase
```

`Parse` and `Compile` start as stubs that panic with `unimplemented: phase 2.1 Parse` and `unimplemented: phase 2.2 Compile`. The shipped examples wrap each case in `oracle.Case`, so until you write `Parse` they fail without stopping the rest of `tests/`.

Name your own tests by phase, as in Project 1. The grading suites are packages of their own, so no name in your test files can collide with one of theirs, even one like `TestP2_1_Precedence` below. In package `tests`, just avoid the four names `tests/parse_test.go` declares.

## Contract conventions

The doc comments state what `Parse`, `SyntaxError`, and `Compile` each promise. What spans them:

- **Every string is a legal argument.** Project 1's functions took machines and `Node`s that had to come from the contract. `Parse` and `Compile` take any Go string, including the empty string, non-ASCII text, and invalid UTF-8. They answer every one and never panic - all strings either parse or don't.
- **No call depends on an earlier one:** nothing carries over from one `Parse` to the next.
- **Unspecified:** the tree's shape, `Msg`'s wording, the `Node` or `Regexp` returned with an error, and performance beyond finishing promptly.

## Design considerations

- **Check your `String()` against the grammar.** Phase 1.5 graded `String()` through Go's `regexp`, which accepts some things course syntax doesn't, such as `(?:a|b)`, `()`, and the empty alternative in `(a|)`. If your `String()` produces one, your own correct `Parse` will reject it, and phase 2.2's grading, which parses your renderings, will fail. `TestP2_1_RoundTrip` in `tests/parse_test.go` finds this as soon as `Parse` works. The fix belongs in `String()`.
- **Return a literal nil error.** A nil `*SyntaxError` stored in an `error` is not a nil `error`. A parser whose helpers return `*SyntaxError` gets this wrong if it passes their result straight through:

  ```go
  n, err := parseAll(pattern) // err is a *SyntaxError, and nil here
  return n, err            // returns a non-nil error holding a nil pointer
  ```

  Check for nil and return a literal `nil` instead.
- **Multi-byte tokens can be cut off.** A pattern can end partway through `ε` or `∅`, or hold their first byte followed by something else. Reading `pattern[i+1]` to finish an `ε` panics on a pattern that ends in its first byte; `strings.HasPrefix(pattern[i:], "ε")` doesn't. And `for i, r := range pattern` decodes runes, turning each invalid byte into U+FFFD, which is easy to mistake for a real character. Offsets count bytes, so `pattern[:e.Offset]` is exactly the part that could begin a valid pattern.
- **Scale.** The grading tests include deeply nested and long patterns, with the sizes and time limits in the handout. A correct parser takes milliseconds on each. Recursion as deep as the nesting is fine, because Go stacks grow as needed. Work that grows exponentially with the nesting depth isn't.

## Testing

**Prefer deciding to sampling, still.** A parse is right when it denotes the right language, and with Project 1 you can decide that. Build the AST you expect by hand, compile both over one alphabet, determinize, and compare:

```go
func TestP2_1_Precedence(t *testing.T) {
    sigma := []automata.Symbol{'a', 'b', 'c'}
    got := oracle.Must(regex.Parse("a|bc*"))
    want := regex.Or(regex.Lit('a'), regex.Seq(regex.Lit('b'), regex.Star(regex.Lit('c'))))
    oracle.AssertEquivalent(t,
        automata.NFAToDFA(regex.Compile(got, sigma)),
        automata.NFAToDFA(regex.Compile(want, sigma)))
}
```

Swap the hand-built AST for a second `Parse` and the same code decides whether two patterns denote the same language: a regex equivalence checker, with no new code.

Deciding runs your Project 1 `Equivalent`, so its speed is your product construction's. On machines of hundreds of states, the size of the handout's largest phase 2.1 patterns, that can take seconds or more. Phase 2.1's grading doesn't run your automata: it checks each parse with Go's `regexp` on chosen words, below. At those sizes, so can you.

Other checks worth having:

- **The round trip.** `Parse(n.String())` should denote the language of `n` for every `Node` `n`, as `TestP2_1_RoundTrip` checks on generated `Node`s.
- **Go's `regexp` as an oracle.** A valid course pattern `p` is also a Go pattern after `oracle.GoPattern(p)`, with the same precedence. So ``regexp.MustCompile(`\A(?:` + oracle.GoPattern(p) + `)\z`)`` decides the language of `p`, and you can compare `engine.Compile(p)` against it on generated words. Only valid patterns: `GoPattern` translates tokens and checks nothing.
- **Rejections and offsets.** A table of malformed patterns and their offsets, read with `errors.As`, as in `TestP2_1_Rejections`, which already has every example in `Parse`'s doc comment. Go beyond them: include patterns that are wrong at the end, in the middle, inside parentheses, after more than one factor or alternative, and in bytes that aren't tokens at all.
