package main

type Expr struct{}
type Unary struct {
	operator Token
	right    Expr
}
type Binary struct {
	left     Expr
	operator Token
	right    Expr
}
type Grouping struct {
	expression Expr
}
type Literal struct {
	value Object
}
