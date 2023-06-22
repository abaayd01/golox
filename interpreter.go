package main

import (
	"fmt"
)

type Interpreter struct {
	Lox         *Lox
	Environment *Environment
}

func (i Interpreter) InterpretStatements(statements []Stmt) error {
	for _, stmt := range statements {
		err := i.execute(stmt)
		if err != nil {
			i.Lox.runtimeError(err.(RuntimeError))
			return err
		}
	}

	return nil
}

func (i Interpreter) VisitStmtExpression(stmt StmtExpression) (any, error) {
	return i.evaluate(stmt.expression)
}

func (i Interpreter) VisitStmtWhile(stmt StmtWhile) (any, error) {
	for {
		cond, err := i.evaluate(stmt.condition)
		if err != nil {
			return nil, err
		}

		if !isTruthy(cond) {
			break
		}

		err = i.execute(stmt.body)
		if err != nil {
			return nil, err
		}
	}
	return nil, nil
}

func (i Interpreter) VisitStmtPrint(stmt StmtPrint) (any, error) {
	value, err := i.evaluate(stmt.expression)
	if err != nil {
		return nil, err
	}
	fmt.Println(value)
	return nil, nil
}

func (i Interpreter) VisitStmtVar(stmt StmtVar) (any, error) {
	if stmt.initializer == nil {
		i.Environment.Define(stmt.name.lexeme, nil)
		return nil, nil
	}

	value, err := i.evaluate(stmt.initializer)
	if err != nil {
		return nil, err
	}

	i.Environment.Define(stmt.name.lexeme, value)
	return value, nil
}

func (i Interpreter) VisitStmtBlock(expr StmtBlock) (any, error) {
	err := i.executeBlock(expr.statements, *NewEnvironmentWithEnclosing(i.Environment))
	if err != nil {
		return nil, err
	}
	return nil, nil
}

func (i Interpreter) VisitStmtIf(expr StmtIf) (any, error) {
	cond, err := i.evaluate(expr.condition)
	if err != nil {
		return nil, err
	}

	if isTruthy(cond) {
		i.execute(expr.thenBranch)
	} else if expr.elseBranch != nil {
		i.execute(expr.elseBranch)
	}

	return nil, nil
}

func (i Interpreter) executeBlock(statements []Stmt, environment Environment) error {
	prevEnv := i.Environment

	for _, stmt := range statements {
		i.Environment = &environment
		err := i.execute(stmt)
		if err != nil {
			i.Environment = prevEnv
			return err
		}
	}

	i.Environment = prevEnv
	return nil
}

func (i Interpreter) InterpretExpression(expr Expr) (any, error) {
	val, err := i.evaluate(expr)
	if err != nil {
		i.Lox.runtimeError(err.(RuntimeError))
		return nil, err
	}

	fmt.Println(stringify(val))
	return val, nil
}

func (i Interpreter) VisitLiteral(expr Literal) (any, error) {
	return expr.value, nil
}

func (i Interpreter) VisitGrouping(expr Grouping) (any, error) {
	return i.evaluate(expr.expression)
}

func (i Interpreter) VisitUnary(expr Unary) (any, error) {
	right, _ := i.evaluate(expr.right)

	switch expr.operator.tokenType {
	case MINUS:
		err := checkNumberOperand(expr.operator, right)
		if err != nil {
			return nil, err
		}
		return -1.0 * right.(float64), nil // nb will panic if not a float
	case BANG:
		return !isTruthy(right), nil
	}

	return nil, nil
}

func (i Interpreter) VisitBinary(expr Binary) (any, error) {
	left, _ := i.evaluate(expr.left)
	right, _ := i.evaluate(expr.right)

	operandsAreBothStrings := checkStringOperands(left, right)

	// support for string concatenation
	if operandsAreBothStrings {
		leftStr, _ := left.(string)
		rightStr, _ := right.(string)

		switch expr.operator.tokenType {
		case PLUS:
			return fmt.Sprintf("%s%s", leftStr, rightStr), nil
		}

		return nil, NewRuntimeError(expr.operator, fmt.Sprintf("Cannot use operator '%s' with string operands", expr.operator.lexeme))
	}

	err := checkNumberOperands(expr.operator, left, right)
	if err != nil {
		return nil, err
	}

	switch expr.operator.tokenType {
	case MINUS:
		return left.(float64) - right.(float64), nil
	case PLUS:
		return left.(float64) + right.(float64), nil
	case SLASH:
		err = checkDivideByZero(expr.operator, right)
		if err != nil {
			return nil, err
		}

		return left.(float64) / right.(float64), nil
	case STAR:
		return left.(float64) * right.(float64), nil
	case GREATER:
		return left.(float64) > right.(float64), nil
	case GREATER_EQUAL:
		return left.(float64) >= right.(float64), nil
	case LESS:
		return left.(float64) < right.(float64), nil
	case LESS_EQUAL:
		return left.(float64) <= right.(float64), nil
	case BANG_EQUAL:
		return !isEqual(left, right), nil
	case EQUAL_EQUAL:
		return isEqual(left, right), nil
	}

	return nil, nil
}

func (i Interpreter) VisitLogical(expr Logical) (any, error) {
	if expr.operator.tokenType == OR {
		lv, err := i.evaluate(expr.left)
		if err != nil {
			return nil, err
		}

		if isTruthy(lv) {
			return lv, nil
		}

		rv, err := i.evaluate(expr.right)
		if err != nil {
			return nil, err
		}

		return rv, nil
	}

	if expr.operator.tokenType == AND {
		lv, err := i.evaluate(expr.left)
		if err != nil {
			return nil, err
		}

		if !isTruthy(lv) {
			return lv, nil
		}

		rv, err := i.evaluate(expr.right)
		if err != nil {
			return nil, err
		}

		return rv, nil
	}

	return nil, NewRuntimeError(expr.operator, "Logical operator must be 'or' or 'and'.")
}

func (i Interpreter) VisitVar(expr Var) (any, error) {
	return i.Environment.Get(expr.name)
}

func (i Interpreter) VisitAssign(expr Assign) (any, error) {
	val, err := i.evaluate(expr.value)
	if err != nil {
		return nil, err
	}

	i.Environment.Assign(expr.name, val)

	return nil, nil
}

func (i Interpreter) execute(stmt Stmt) error {
	_, err := stmt.Accept(i)
	return err
}

func (i Interpreter) evaluate(expr Expr) (any, error) {
	return expr.Accept(i)
}

func stringify(val any) string {
	// could insert more stringifying logic here if needed
	return fmt.Sprintf("%v", val)
}

func checkNumberOperand(op Token, expr any) error {
	_, ok := expr.(float64)
	if !ok {
		return NewRuntimeError(op, "Operand must be a number.")
	}
	return nil
}

func checkNumberOperands(op Token, left any, right any) error {
	_, leftOk := left.(float64)
	_, rightOk := right.(float64)
	if !leftOk || !rightOk {
		return NewRuntimeError(op, "Operands must both be numbers.")
	}
	return nil
}

func checkStringOperands(left any, right any) bool {
	_, leftIsString := left.(string)
	_, rightIsString := right.(string)
	return leftIsString && rightIsString
}

func checkDivideByZero(op Token, right any) error {
	if right == 0.0 {
		return NewRuntimeError_DivideByZero(op)
	}

	return nil
}

// isTruthy checks if x is a truthy value or not.
// Truthy being defined as not falsey.
// Falsey being defined as:
//   - false
//   - nil
func isTruthy(x any) bool {
	if x == nil {
		return false
	}

	xBool, ok := x.(bool)

	// if x is not a boolean, (and not nil) return true
	if !ok {
		return true
	}

	return xBool
}

func isEqual(a any, b any) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil {
		return false
	}

	return a == b
}

type RuntimeError struct {
	Token Token
	msg   string
}

func NewRuntimeError(token Token, msg string) RuntimeError {
	return RuntimeError{
		Token: token,
		msg:   msg,
	}
}

func NewRuntimeError_DivideByZero(token Token) RuntimeError {
	return NewRuntimeError(token, "Cannot divide by zero")
}

func (e RuntimeError) Error() string {
	return e.msg
}
