package main

import "fmt"

type Environment struct {
	EnclosingEnv *Environment
	Values       map[string]any
}

func NewEnvironment() *Environment {
	return &Environment{Values: map[string]any{}}
}

func NewEnvironmentWithEnclosing(env *Environment) *Environment {
	return &Environment{
		EnclosingEnv: env,
		Values:       map[string]any{},
	}
}

func (e *Environment) Define(key string, value any) {
	e.Values[key] = value
}

func (e *Environment) Get(name Token) (any, error) {
	key := name.lexeme
	val, ok := e.Values[key]
	if ok {
		return val, nil
	}

	if e.EnclosingEnv != nil {
		return e.EnclosingEnv.Get(name)
	}

	return nil, NewRuntimeError(name, fmt.Sprintf("Undefined variable '%s'.", name.lexeme))
}

func (e *Environment) Assign(name Token, value any) error {
	_, err := e.Get(name)

	if err != nil {
		return err
	}

	e.Values[name.lexeme] = value
	return nil
}
