package tests

import (
	"401/engine"
	"401/regex"
)

// Compile-time assertions for the declarations Project 2 adds. The shipped
// contract_test.go can't change, so they get a file of their own.

var (
	_ error = (*regex.SyntaxError)(nil)

	_ func(pattern string) (regex.Node, error)     = regex.Parse
	_ func(pattern string) (*engine.Regexp, error) = engine.Compile
)
