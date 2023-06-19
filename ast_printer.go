package main

import (
	"fmt"
	"strings"
)

type AstPrinter struct{}

func (a AstPrinter) Print(expr Expr) string {
	s, _ := expr.Accept(a)
	return s.(string)
}

func (a AstPrinter) parenthesize(name string, exprs ...Expr) string {
	var sb strings.Builder
	sb.WriteString("(")
	sb.WriteString(name)
	sb.WriteString(" ")

	for i, e := range exprs {
		s, _ := e.Accept(a)
		str := s.(string)
		sb.WriteString(str)

		if i < len(exprs)-1 {
			sb.WriteString(" ")
		}
	}

	sb.WriteString(")")
	return sb.String()
}

func (a AstPrinter) VisitUnary(expr Unary) (any, error) {
	return a.parenthesize(expr.operator.lexeme, expr.right), nil
}

func (a AstPrinter) VisitBinary(expr Binary) (any, error) {
	return a.parenthesize(expr.operator.lexeme, expr.left, expr.right), nil
}

func (a AstPrinter) VisitGrouping(expr Grouping) (any, error) {
	return a.parenthesize("group", expr.expression), nil
}

func (a AstPrinter) VisitLiteral(expr Literal) (any, error) {
	return fmt.Sprintf("%v", expr.value), nil
}

func (a AstPrinter) VisitVar(expr Var) (any, error) {
	return fmt.Sprintf("%v", expr.name.lexeme), nil
}

func (a AstPrinter) VisitAssign(expr Assign) (any, error) {
	return fmt.Sprintf("%v = %v", expr.name.lexeme, expr.value), nil
}
