package main

type Expr interface {
	Accept(visitor Visitor) any
}
type Unary struct {
	operator Token
	right    Expr
}

func (t Unary) Accept(visitor Visitor) any {
	return visitor.VisitUnary(t)
}

type Binary struct {
	left     Expr
	operator Token
	right    Expr
}

func (t Binary) Accept(visitor Visitor) any {
	return visitor.VisitBinary(t)
}

type Grouping struct {
	expression Expr
}

func (t Grouping) Accept(visitor Visitor) any {
	return visitor.VisitGrouping(t)
}

type Literal struct {
	value Object
}

func (t Literal) Accept(visitor Visitor) any {
	return visitor.VisitLiteral(t)
}

type Visitor interface {
	VisitUnary(expr Unary) any
	VisitBinary(expr Binary) any
	VisitGrouping(expr Grouping) any
	VisitLiteral(expr Literal) any
}
