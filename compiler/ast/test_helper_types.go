//+build !release

package ast

import "testing"

const nodeFailureMsgFormat = "Node is not %s, is %T"

type TestingIdentifier string

type testingNode interface {
	// Belows are test helpers
	NameIs(name string) bool
}

type TestingExpression interface {
	Expression
	// Test Helpers
	IsArrayExpression(t *testing.T) *ArrayExpression
	IsAssignExpression(t *testing.T) *AssignExpression
	IsCallExpression(t *testing.T) *CallExpression
	IsConditionalExpression(t *testing.T) *ConditionalExpression
	IsConstant(t *testing.T) *Constant
	IsHashExpression(t *testing.T) *HashExpression
	IsIdentifier(t *testing.T) *Identifier
	IsIfExpression(t *testing.T) *IfExpression
	IsInfixExpression(t *testing.T) *InfixExpression
	IsIntegerLiteral(t *testing.T) *IntegerLiteral
	IsStringLiteral(t *testing.T) *StringLiteral
	IsYieldExpression(t *testing.T) *YieldExpression
}