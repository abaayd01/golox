package main

type Expr[R any] interface {
	Accept(visitor Visitor[R]) R
}
type Unary[R any] struct {
	operator Token
	right    Expr[R]
}

func (t Unary[R]) Accept(visitor Visitor[R]) R {
	return visitor.VisitUnary(t)
}

type Binary[R any] struct {
	left     Expr[R]
	operator Token
	right    Expr[R]
}

func (t Binary[R]) Accept(visitor Visitor[R]) R {
	return visitor.VisitBinary(t)
}

type Grouping[R any] struct {
	expression Expr[R]
}

func (t Grouping[R]) Accept(visitor Visitor[R]) R {
	return visitor.VisitGrouping(t)
}

type Literal[R any] struct {
	value Object
}

func (t Literal[R]) Accept(visitor Visitor[R]) R {
	return visitor.VisitLiteral(t)
}

type Visitor[R any] interface {
	VisitUnary(expr Unary[R]) R
	VisitBinary(expr Binary[R]) R
	VisitGrouping(expr Grouping[R]) R
	VisitLiteral(expr Literal[R]) R
}
