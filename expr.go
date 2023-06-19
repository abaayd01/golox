// this file is auto-generated with bin/generate_ast
// DO NOT EDIT
package main

type Expr interface {
	Accept(visitor ExprVisitor) (any, error)
}
type Unary struct {
	operator Token
	right    Expr
}

func (t Unary) Accept(visitor ExprVisitor) (any, error) {
	return visitor.VisitUnary(t)
}

type Binary struct {
	left     Expr
	operator Token
	right    Expr
}

func (t Binary) Accept(visitor ExprVisitor) (any, error) {
	return visitor.VisitBinary(t)
}

type Grouping struct {
	expression Expr
}

func (t Grouping) Accept(visitor ExprVisitor) (any, error) {
	return visitor.VisitGrouping(t)
}

type Literal struct {
	value Object
}

func (t Literal) Accept(visitor ExprVisitor) (any, error) {
	return visitor.VisitLiteral(t)
}

type Var struct {
	name Token
}

func (t Var) Accept(visitor ExprVisitor) (any, error) {
	return visitor.VisitVar(t)
}

type Assign struct {
	name  Token
	value Expr
}

func (t Assign) Accept(visitor ExprVisitor) (any, error) {
	return visitor.VisitAssign(t)
}

type ExprVisitor interface {
	VisitUnary(expr Unary) (any, error)
	VisitBinary(expr Binary) (any, error)
	VisitGrouping(expr Grouping) (any, error)
	VisitLiteral(expr Literal) (any, error)
	VisitVar(expr Var) (any, error)
	VisitAssign(expr Assign) (any, error)
}
