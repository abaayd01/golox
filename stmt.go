// this file is auto-generated with bin/generate_ast
// DO NOT EDIT
package main

type Stmt interface {
	Accept(visitor StmtVisitor) (any, error)
}
type Expression struct {
	expression Expr
}

func (t Expression) Accept(visitor StmtVisitor) (any, error) {
	return visitor.VisitExpression(t)
}

type Print struct {
	expression Expr
}

func (t Print) Accept(visitor StmtVisitor) (any, error) {
	return visitor.VisitPrint(t)
}

type StmtVisitor interface {
	VisitExpression(expr Expression) (any, error)
	VisitPrint(expr Print) (any, error)
}
