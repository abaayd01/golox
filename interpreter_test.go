package main

import (
	"testing"
)

func TestInterpreter_Interpret(t *testing.T) {
	expr := Unary{
		operator: Token{
			tokenType: MINUS,
			lexeme:    "-",
			literal:   nil,
			line:      0,
		},
		right: Literal{value: "abc"},
	}

	l := Lox{}
	interpreter := Interpreter{
		Lox: &l,
	}
	_ = interpreter.Interpret(expr)
}
