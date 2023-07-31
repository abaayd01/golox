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

func TestParser_ParseFunctionDeclaration(t *testing.T) {
	is := is2.New(t)

	lox := Lox{}
	scanner := Scanner{
		lox:     &lox,
		source:  "fun myFun(a, b) { a + b; }",
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
		StmtFunction{
			name: Token{
				tokenType: IDENTIFIER,
				lexeme:    "myFun",
				literal:   nil,
				line:      0,
			},
			params: []Token{
				{
					tokenType: IDENTIFIER,
					lexeme:    "a",
					literal:   nil,
					line:      0,
				},
				{
					tokenType: IDENTIFIER,
					lexeme:    "b",
					literal:   nil,
					line:      0,
				},
			},
			body: StmtBlock{
				statements: []Stmt{
					StmtExpression{
						expression: Binary{
							operator: Token{
								tokenType: PLUS,
								lexeme:    "+",
								literal:   nil,
								line:      0,
							},
							left: Var{
								name: Token{
									tokenType: IDENTIFIER,
									lexeme:    "a",
									literal:   nil,
									line:      0,
								},
							},
							right: Var{
								name: Token{
									tokenType: IDENTIFIER,
									lexeme:    "b",
									literal:   nil,
									line:      0,
								},
							},
						},
					},
				},
			},
		}})
}
