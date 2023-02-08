package main

import "strings"

type AstPrinter struct{}

func (a AstPrinter) parenthesize(name string, exprs ...Expr[string]) string {
	var sb strings.Builder
	sb.WriteString("(")
	sb.WriteString(name)

	for _, e := range exprs {
		s := e.Accept(a)
		sb.WriteString(s)
	}

	return sb.String()
}

func (a AstPrinter) VisitUnary(expr Unary[string]) string {
	//TODO implement me
	panic("implement me")
}

func (a AstPrinter) VisitBinary(expr Binary[string]) string {
	//TODO implement me
	panic("implement me")
}

func (a AstPrinter) VisitGrouping(expr Grouping[string]) string {
	//TODO implement me
	panic("implement me")
}

func (a AstPrinter) VisitLiteral(expr Literal[string]) string {
	//TODO implement me
	panic("implement me")
}
