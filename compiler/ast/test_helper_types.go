//+build !release

package ast

import "testing"

const nodeFailureMsgFormat = "Node is not %s, is %v"

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
	IsArrayExpression(t *testing.T) *ArrayExpression
	IsAssignExpression(t *testing.T) *AssignExpression
	IsCallExpression(t *testing.T) *CallExpression
	IsConstant(t *testing.T) *Constant
	IsHashExpression(t *testing.T) *HashExpression
	IsIdentifier(t *testing.T) *Identifier
	IsInfixExpression(t *testing.T) *InfixExpression
	IsStringLiteral(t *testing.T) *StringLiteral
	IsYieldExpression(t *testing.T) *YieldExpression
}

type CodeBlock []TestingStatement

func (cb CodeBlock) NthStmt(n int) TestingStatement {
	return cb[n-1]
}
