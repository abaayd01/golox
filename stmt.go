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

type StmtBlock struct {
	statements []Stmt
}

func (t StmtBlock) Accept(visitor StmtVisitor) (any, error) {
	return visitor.VisitStmtBlock(t)
}

type StmtIf struct {
	condition  Expr
	thenBranch Stmt
	elseBranch Stmt
}

func (t StmtIf) Accept(visitor StmtVisitor) (any, error) {
	return visitor.VisitStmtIf(t)
}

type StmtWhile struct {
	condition Expr
	body      Stmt
}

func (t StmtWhile) Accept(visitor StmtVisitor) (any, error) {
	return visitor.VisitStmtWhile(t)
}

type StmtFunction struct {
	name   Token
	params []Token
	body   StmtBlock
}

func (t StmtFunction) Accept(visitor StmtVisitor) (any, error) {
	return visitor.VisitStmtFunction(t)
}

type StmtReturn struct {
	returnKeyword Token
	value         Expr
}

func (t StmtReturn) Accept(visitor StmtVisitor) (any, error) {
	return visitor.VisitStmtReturn(t)
}

type StmtVisitor interface {
	VisitStmtExpression(expr StmtExpression) (any, error)
	VisitStmtPrint(expr StmtPrint) (any, error)
	VisitStmtVar(expr StmtVar) (any, error)
	VisitStmtBlock(expr StmtBlock) (any, error)
	VisitStmtIf(expr StmtIf) (any, error)
	VisitStmtWhile(expr StmtWhile) (any, error)
	VisitStmtFunction(expr StmtFunction) (any, error)
	VisitStmtReturn(expr StmtReturn) (any, error)
}
