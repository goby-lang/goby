//+build !release

package ast

import "testing"

type TestingIdentifier string

type testingNode interface {
	// Belows are test helpers
	NameIs(name string) bool
}

type TestingStatement interface {
	Statement
	// Test Helpers
	IsClassStmt(t *testing.T, className string) *ClassStatement
	IsDefStmt(t *testing.T) *DefStatement
	IsExpression(t *testing.T) TestingExpression
	IsModuleStmt(t *testing.T, className string) *ModuleStatement
	IsReturnStmt(t *testing.T) *ReturnStatement
	IsWhileStmt(t *testing.T) *WhileStatement
}

type TestingExpression interface {
	Expression
	// Test Helpers
	IsAssignExpression(t *testing.T) *AssignExpression
	IsCallExpression(t *testing.T) *CallExpression
	IsConstant(t *testing.T) *Constant
	IsIdentifier(t *testing.T) *Identifier
	IsInfixExpression(t *testing.T) *InfixExpression
	IsYieldExpression(t *testing.T) *YieldExpression
}

type CodeBlock []TestingStatement

func (cb CodeBlock) NthStmt(n int) TestingStatement {
	return cb[n-1]
}
