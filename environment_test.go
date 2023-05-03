package main

import (
	is2 "github.com/matryer/is"
	"testing"
)

func TestEnvironment_DefineAndGet(t *testing.T) {
	is := is2.New(t)
	env := Environment{
		Values: map[string]any{},
	}

	env.Define("foo", 123)

	result, err := env.Get(Token{
		tokenType: IDENTIFIER,
		lexeme:    "foo",
		literal:   nil,
		line:      0,
	})

	is.NoErr(err)
	is.Equal(result, 123)
}

func TestEnvironment_GetRuntimeError(t *testing.T) {
	is := is2.New(t)
	env := Environment{
		Values: map[string]any{},
	}

	token := Token{
		tokenType: IDENTIFIER,
		lexeme:    "foo",
		literal:   nil,
		line:      0,
	}

	result, err := env.Get(Token{
		tokenType: IDENTIFIER,
		lexeme:    "foo",
		literal:   nil,
		line:      0,
	})

	is.Equal(err, RuntimeError{
		Token: token,
		msg:   "Undefined variable 'foo'.",
	})

	is.Equal(result, nil)
}
