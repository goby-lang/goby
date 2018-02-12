package ast

import "testing"

type TestingExpression interface {
	Expression
	// Test Helpers
	IsArrayExpression(t *testing.T) *TestableArrayExpression
	IsAssignExpression(t *testing.T) *TestableAssignExpression
	IsCallExpression(t *testing.T) *TestableCallExpression
	IsConditionalExpression(t *testing.T) *TestableConditionalExpression
	IsConstant(t *testing.T) *TestableConstant
	IsHashExpression(t *testing.T) *TestableHashExpression
	IsIdentifier(t *testing.T) *TestableIdentifier
	IsIfExpression(t *testing.T) *TestableIfExpression
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

/*TestableConditionalExpression*/

type TestableConditionalExpression struct {
	*ConditionalExpression
	t *testing.T
}

func (tce *TestableConditionalExpression) TestableConsequence() CodeBlock {
	var tss []TestingStatement
	for _, stmt := range tce.Consequence.Statements {
		tss = append(tss, stmt.(TestingStatement))
	}

	return tss
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

/*TestableIfExpression*/

type TestableIfExpression struct {
	*IfExpression
	t *testing.T
}

func (tie *TestableIfExpression) ShouldHasNumberOfConditionals(n int) {
	if len(tie.Conditionals) != n {
		tie.t.Fatalf("Expect if expression to have %d conditionals, got %d", n, len(tie.Conditionals))
	}
}

func (tie *TestableIfExpression) TestableConditionals() (tes []TestingExpression) {
	for _, cond := range tie.Conditionals {
		tes = append(tes, cond)
	}

	return
}

func (tie *TestableIfExpression) TestableAlternative() CodeBlock {
	var tss []TestingStatement
	for _, stmt := range tie.Alternative.Statements {
		tss = append(tss, stmt.(TestingStatement))
	}

	return tss
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
