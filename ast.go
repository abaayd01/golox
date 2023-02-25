package main

type Expr interface {
	Accept(visitor Visitor) (any, error)
}
type Unary struct {
	operator Token
	right    Expr
}

func (t Unary) Accept(visitor Visitor) (any, error) {
	return visitor.VisitUnary(t)
}

type Binary struct {
	left     Expr
	operator Token
	right    Expr
}

func (t Binary) Accept(visitor Visitor) (any, error) {
	return visitor.VisitBinary(t)
}

type Grouping struct {
	expression Expr
}

func (t Grouping) Accept(visitor Visitor) (any, error) {
	return visitor.VisitGrouping(t)
}

type Literal struct {
	value Object
}

func (t Literal) Accept(visitor Visitor) (any, error) {
	return visitor.VisitLiteral(t)
}

type Visitor interface {
	VisitUnary(expr Unary) (any, error)
	VisitBinary(expr Binary) (any, error)
	VisitGrouping(expr Grouping) (any, error)
	VisitLiteral(expr Literal) (any, error)
}
