package ast

import "testing"

type TestingExpression interface {
	Expression
	// Test Helpers
	IsArrayExpression(t *testing.T) *TestableArrayExpression
	IsAssignExpression(t *testing.T) *TestableAssignExpression
	IsCallExpression(t *testing.T) *TestableCallExpression
	IsConditionalExpression(t *testing.T) *ConditionalExpression
	IsConstant(t *testing.T) *TestableConstant
	IsHashExpression(t *testing.T) *TestableHashExpression
	IsIdentifier(t *testing.T) *TestableIdentifier
	IsIfExpression(t *testing.T) *IfExpression
	IsInfixExpression(t *testing.T) *TestableInfixExpression
	IsIntegerLiteral(t *testing.T) *IntegerLiteral
	IsStringLiteral(t *testing.T) *StringLiteral
	IsYieldExpression(t *testing.T) *YieldExpression
}

/*TestableArrayExpression*/

type TestableArrayExpression struct {
	*ArrayExpression
	t *testing.T
}

func (tae *TestableArrayExpression) TestableElements() (tes []TestingExpression) {
	for _, elem := range tae.Elements {
		tes = append(tes, elem.(TestingExpression))
	}

	return
}

/*TestableAssignExpression*/

type TestableAssignExpression struct {
	*AssignExpression
	t *testing.T
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

/*TestableHashExpression*/

type TestableHashExpression struct {
	*HashExpression
	t *testing.T
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

/*TestableInfixExpression*/

type TestableInfixExpression struct {
	*InfixExpression
	t *testing.T
}

// ShouldHasOperator checks if the infix expression has expected operator
func (tie *TestableInfixExpression) ShouldHasOperator(expectedOperator string) {
	if tie.Operator != expectedOperator {
		tie.t.Fatalf("Expect infix expression to have %s operator, got %s", expectedOperator, tie.Operator)
	}
}

// LeftExpression returns infix expression's left expression as TestingExpression
func (tie *TestableInfixExpression) TestableLeftExpression() TestingExpression {
	return tie.Left.(TestingExpression)
}

// RightExpression returns infix expression's right expression as TestingExpression
func (tie *TestableInfixExpression) TestableRightExpression() TestingExpression {
	return tie.Right.(TestingExpression)
}

/*TestableConstant*/

type TestableConstant struct {
	*Constant
	t *testing.T
}

func (tc *TestableConstant) ShouldHasName(expectedName string) {
	if tc.Value != expectedName {
		tc.t.Fatalf("expect current identifier to be '%s', got '%s'", expectedName, tc.Value)
	}
}
