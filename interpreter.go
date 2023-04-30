package main

import (
	"fmt"
)

type Interpreter struct {
	Lox *Lox
}

func (i Interpreter) Interpret(expr Expr) error {
	val, err := i.evaluate(expr)
	if err != nil {
		i.Lox.runtimeError(err.(RuntimeError))
		return err
	}

	fmt.Println(stringify(val))
	return nil
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
		leftStr, _ := expr.left.(any).(string)
		rightStr, _ := expr.right.(any).(string)

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
	case EQUAL:
		return isEqual(left, right), nil
	}

	return nil, nil
}

func checkStringOperands(left any, right any) bool {
	_, leftIsString := left.(string)
	_, rightIsString := right.(string)
	return leftIsString && rightIsString
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

func (e RuntimeError) Error() string {
	return e.msg
}
