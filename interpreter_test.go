package main

import (
	"github.com/matryer/is"
	"testing"
)

func TestInterpreter_Interpret(t *testing.T) {
	type test struct {
		description string
		input       Expr
		expected    any
	}

	tests := []test{
		{
			description: "it adds two numbers",
			input: Binary{
				operator: Token{
					tokenType: PLUS,
					lexeme:    "+",
					literal:   nil,
					line:      0,
				},
				left:  Literal{value: 1.0},
				right: Literal{value: 2.0},
			},
			expected: 3.0,
		},
		{
			description: "it subtracts two numbers",
			input: Binary{
				operator: Token{
					tokenType: MINUS,
					lexeme:    "-",
					literal:   nil,
					line:      0,
				},
				left:  Literal{value: 1.0},
				right: Literal{value: 2.0},
			},
			expected: -1.0,
		},
		{
			description: "concatenates two strings",
			input: Binary{
				operator: Token{
					tokenType: PLUS,
					lexeme:    "+",
					literal:   nil,
					line:      0,
				},
				left:  Literal{value: "hello "},
				right: Literal{value: "world!"},
			},
			expected: "hello world!",
		},
	}

	is := is.NewRelaxed(t)
	for _, tc := range tests {
		t.Run(tc.description, func(t *testing.T) {
			is := is.NewRelaxed(t)
			interpreter := Interpreter{
				Lox: &Lox{},
			}

			result, err := interpreter.Interpret(tc.input)

			is.NoErr(err)
			is.Equal(result, tc.expected)
		})
	}
}
