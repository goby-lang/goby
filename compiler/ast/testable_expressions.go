package ast

import "testing"

type TestingExpression interface {
	Expression
	// Test Helpers
	IsArrayExpression(t *testing.T) *ArrayExpression
	IsAssignExpression(t *testing.T) *AssignExpression
	IsCallExpression(t *testing.T) *TestableCallExpression
	IsConditionalExpression(t *testing.T) *ConditionalExpression
	IsConstant(t *testing.T) *Constant
	IsHashExpression(t *testing.T) *HashExpression
	IsIdentifier(t *testing.T) *TestableIdentifier
	IsIfExpression(t *testing.T) *IfExpression
	IsInfixExpression(t *testing.T) *InfixExpression
	IsIntegerLiteral(t *testing.T) *IntegerLiteral
	IsStringLiteral(t *testing.T) *StringLiteral
	IsYieldExpression(t *testing.T) *YieldExpression
}

/*TestableCallExpression*/

type TestableCallExpression struct {
	*CallExpression
	t *testing.T
}

// NthArgument returns n-th argument of the call expression as TestingExpression
func (tce *TestableCallExpression) NthArgument(n int) TestingExpression {
	return tce.Arguments[n-1].(TestingExpression)
}

// ReceiverExpression returns call expression's receiver as TestingExpression
func (tce *TestableCallExpression) TestableReceiver() TestingExpression {
	return tce.Receiver.(TestingExpression)
}

// ShouldHasMethodName
func (tce *TestableCallExpression) ShouldHasMethodName(expectedName string) {
	if tce.Method != expectedName {
		tce.t.Fatalf("expect call expression's method name to be '%s', got '%s'", expectedName, tce.Method)
	}
}

/*TestableIdentifier*/

type TestableIdentifier struct {
	*Identifier
	t *testing.T
}

func (ti *TestableIdentifier) ShouldHasName(expectedName string) {
	if ti.Value != expectedName {
		ti.t.Fatalf("expect current identifier to be '%s', got '%s'", expectedName, ti.Value)
	}
}
