package main

import (
	"fmt"
	"strings"
)

type AstPrinter struct{}

func (a AstPrinter) Print(expr Expr) string {
	return expr.Accept(a).(string)
}

func (a AstPrinter) parenthesize(name string, exprs ...Expr) string {
	var sb strings.Builder
	sb.WriteString("(")
	sb.WriteString(name)
	sb.WriteString(" ")

	for i, e := range exprs {
		s := e.Accept(a).(string)
		sb.WriteString(s)

		if i < len(exprs)-1 {
			sb.WriteString(" ")
		}
	}

	sb.WriteString(")")
	return sb.String()
}

func (a AstPrinter) VisitUnary(expr Unary) any {
	return a.parenthesize(expr.operator.lexeme, expr.right)
}

func (a AstPrinter) VisitBinary(expr Binary) any {
	return a.parenthesize(expr.operator.lexeme, expr.left, expr.right)
}

func (a AstPrinter) VisitGrouping(expr Grouping) any {
	return a.parenthesize("group", expr.expression)
}

func (a AstPrinter) VisitLiteral(expr Literal) any {
	return fmt.Sprintf("%v", expr.value)
}
