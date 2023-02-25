package main

import "testing"

// TestParser_Parse is just a smoke test of the parser for now.
// It just checks the parser can parse something without panicking.
func TestParser_Parse(t *testing.T) {
	l := Lox{}
	_ = l.run("1+2")
}

func TestLox_Run(t *testing.T) {
	l := Lox{}
	_ = l.run("1 + 2")
}
