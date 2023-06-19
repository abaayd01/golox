package main

import "fmt"

type Environment struct {
	Values map[string]any
}

func (e *Environment) Define(key string, value any) {
	e.Values[key] = value
}

func (e *Environment) Get(name Token) (any, error) {
	key := name.lexeme
	val, ok := e.Values[key]
	if !ok {
		return nil, NewRuntimeError(name, fmt.Sprintf("Undefined variable '%s'.", name.lexeme))
	}

	return val, nil
}

func (e *Environment) Assign(name Token, value any) error {
	_, err := e.Get(name)

	if err != nil {
		return err
	}

	e.Values[name.lexeme] = value
	return nil
}
