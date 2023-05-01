package main

import "testing"

func TestLox_Run(t *testing.T) {
	l := Lox{}
	_ = l.run("1 + 2")
}
