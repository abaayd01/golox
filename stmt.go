// this file is auto-generated with bin/generate_ast
// DO NOT EDIT
package main

type Stmt interface {
	Accept(visitor StmtVisitor) (any, error)
}
type StmtExpression struct {
	expression Expr
}

func (t StmtExpression) Accept(visitor StmtVisitor) (any, error) {
	return visitor.VisitStmtExpression(t)
}

type StmtPrint struct {
	expression Expr
}

func (t StmtPrint) Accept(visitor StmtVisitor) (any, error) {
	return visitor.VisitStmtPrint(t)
}

type StmtVar struct {
	name        Token
	initializer Expr
}

func (t StmtVar) Accept(visitor StmtVisitor) (any, error) {
	return visitor.VisitStmtVar(t)
}

type StmtVisitor interface {
	VisitStmtExpression(expr StmtExpression) (any, error)
	VisitStmtPrint(expr StmtPrint) (any, error)
	VisitStmtVar(expr StmtVar) (any, error)
}
