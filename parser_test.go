package main

import (
	"testing"

	is2 "github.com/matryer/is"
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

func TestParser_ParseReturnStatementWithExpression(t *testing.T) {
	is := is2.New(t)

	lox := Lox{}
	scanner := Scanner{
		lox:     &lox,
		source:  "fun myFun() { return 1 + 2;}",
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

	expected := []Stmt{
		StmtFunction{
			name: Token{
				tokenType: IDENTIFIER,
				lexeme:    "myFun",
				literal:   nil,
				line:      0,
			},
			params: nil,
			body: StmtBlock{
				statements: []Stmt{
					StmtReturn{
						returnKeyword: Token{
							tokenType: RETURN,
							lexeme:    "return",
							literal:   nil,
							line:      0,
						},
						value: Binary{
							left: Literal{
								value: 1.0,
							},
							operator: Token{
								tokenType: PLUS,
								lexeme:    "+",
								literal:   nil,
								line:      0,
							},
							right: Literal{
								value: 2.0,
							},
						},
					},
				},
			},
		}}

	is.Equal(statements, expected)
}

func TestParser_ParseReturnStatementWithNoReturnValue(t *testing.T) {
	is := is2.New(t)

	lox := Lox{}
	scanner := Scanner{
		lox:     &lox,
		source:  "fun myFun() { return;}",
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

	expected := []Stmt{
		StmtFunction{
			name: Token{
				tokenType: IDENTIFIER,
				lexeme:    "myFun",
				literal:   nil,
				line:      0,
			},
			params: nil,
			body: StmtBlock{
				statements: []Stmt{
					StmtReturn{
						returnKeyword: Token{
							tokenType: RETURN,
							lexeme:    "return",
							literal:   nil,
							line:      0,
						},
						value: nil,
					},
				},
			},
		}}

	is.Equal(statements, expected)
}
