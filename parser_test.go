package main

import (
	is2 "github.com/matryer/is"
	"testing"
)

func TestParser_Parse(t *testing.T) {
	is := is2.New(t)

	lox := Lox{}
	scanner := Scanner{
		lox:     &lox,
		source:  "var foo = 20;",
		tokens:  nil,
		start:   0,
		current: 0,
		line:    0,
	}

	tokens, err := scanner.scanTokens()
	is.NoErr(err)

	parser := Parser{
		Lox:     &lox,
		Tokens:  tokens,
		current: 0,
	}

	statements, err := parser.Parse()
	is.NoErr(err)

	is.Equal(statements, []Stmt{
		StmtVar{
			name: Token{
				tokenType: IDENTIFIER,
				lexeme:    "foo",
				literal:   nil,
				line:      0,
			},
			initializer: Literal{value: 20},
		}})
}
