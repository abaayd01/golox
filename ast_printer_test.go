package main

import (
	"fmt"
	"testing"
)

func TestAstPrinter_Print(t *testing.T) {
	astPrinter := AstPrinter{}

	l := Literal{value: 1}
	r := Literal{value: 2}
	plusToken := Token{
		tokenType: PLUS,
		lexeme:    "+",
		literal:   nil,
		line:      0,
	}
	b := Binary{
		left:     l,
		operator: plusToken,
		right:    r,
	}

	res := astPrinter.Print(b)

	expected := "(+ 1 2)"

	if res != expected {
		t.Errorf("result: %s, expected: %s", res, expected)
	}
	fmt.Println(res)
}
