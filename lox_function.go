package main

type LoxFunction struct {
	declaration StmtFunction
}

func (l LoxFunction) call(i Interpreter, arguments []Object) Object {
	environment := NewGlobalEnvironment()

	for i := 0; i < len(l.declaration.params); i++ {
		environment.Define(l.declaration.params[i].lexeme, arguments[i])
	}

	i.executeBlock([]Stmt{l.declaration.body}, *environment)

	return nil
}

func (l LoxFunction) arity() int {
	return len(l.declaration.params)
}

func (l LoxFunction) toString() string {
	return "<fn " + l.declaration.name.lexeme + ">"
}
