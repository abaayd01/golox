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

func TestInterpreter_AddNumbers(t *testing.T) {
	expr := Binary{
		operator: Token{
			tokenType: PLUS,
			lexeme:    "+",
			literal:   nil,
			line:      0,
		},
		left:  Literal{value: 1.0},
		right: Literal{value: 2.0},
	}

	l := Lox{}
	interpreter := Interpreter{
		Lox: &l,
	}
	_ = interpreter.Interpret(expr)
}

func TestInterpreter_ConcatStrings(t *testing.T) {
	expr := Binary{
		operator: Token{
			tokenType: PLUS,
			lexeme:    "+",
			literal:   nil,
			line:      0,
		},
		left:  Literal{value: "a"},
		right: Literal{value: "b"},
	}

	l := Lox{}
	interpreter := Interpreter{
		Lox: &l,
	}
	_ = interpreter.Interpret(expr)
}

func TestInterpreter_ConcatStringsError(t *testing.T) {
	expr := Binary{
		operator: Token{
			tokenType: MINUS,
			lexeme:    "-",
			literal:   nil,
			line:      0,
		},
		left:  Literal{value: "a"},
		right: Literal{value: "b"},
	}

	l := Lox{}
	interpreter := Interpreter{
		Lox: &l,
	}
	_ = interpreter.Interpret(expr)
}
