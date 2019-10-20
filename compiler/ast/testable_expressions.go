//+build !release

package ast

import "testing"

// TestableExpression interface implements basic Expression's functions, and helper functions to assert node's type
type TestableExpression interface {
	Expression
	// Test Helpers
	IsArrayExpression(t *testing.T) *TestableArrayExpression
	IsAssignExpression(t *testing.T) *TestableAssignExpression
	IsBooleanExpression(t *testing.T) *TestableBooleanExpression
	IsCallExpression(t *testing.T) *TestableCallExpression
	IsConditionalExpression(t *testing.T) *TestableConditionalExpression
	IsConstant(t *testing.T) *TestableConstant
	IsHashExpression(t *testing.T) *TestableHashExpression
	IsIdentifier(t *testing.T) *TestableIdentifier
	IsIfExpression(t *testing.T) *TestableIfExpression
	IsInfixExpression(t *testing.T) *TestableInfixExpression
	IsInstanceVariable(t *testing.T) *TestableInstanceVariable
	IsIntegerLiteral(t *testing.T) *TestableIntegerLiteral
	IsSelfExpression(t *testing.T) *TestableSelfExpression
	IsStringLiteral(t *testing.T) *TestableStringLiteral
	IsYieldExpression(t *testing.T) *TestableYieldExpression
}

// TestableArrayExpression
type TestableArrayExpression struct {
	*ArrayExpression
	t *testing.T
}

// TestableElements returns array expression's element nodes and assert them as TestableExpression
func (tae *TestableArrayExpression) TestableElements() (tes []TestableExpression) {
	for _, elem := range tae.Elements {
		tes = append(tes, elem.(TestableExpression))
	}

	return
}

// TestableAssignExpression
type TestableAssignExpression struct {
	*AssignExpression
	t *testing.T
}

// NthVariable returns the nth variable of the assignment as a TestableExpression
func (tae *TestableAssignExpression) NthVariable(n int) TestableExpression {
	return tae.Variables[n-1].(TestableExpression)
}

// TestableValue returns the assignment's value as a TestableExpression
func (tae *TestableAssignExpression) TestableValue() TestableExpression {
	return tae.Value.(TestableExpression)
}

// TestableBooleanExpression
type TestableBooleanExpression struct {
	*BooleanExpression
	t *testing.T
}

// ShouldEqualTo compares if the boolean expression's value equals to the expected value
func (tbe *TestableBooleanExpression) ShouldEqualTo(expected bool) {
	if tbe.Value != expected {
		tbe.t.Helper()
		tbe.t.Fatalf("Expect boolean literal to be %t, got %t", expected, tbe.Value)
	}
}

//TestableCallExpression
type TestableCallExpression struct {
	*CallExpression
	t *testing.T
}

// NthArgument returns n-th argument of the call expression as TestingExpression
func (tce *TestableCallExpression) NthArgument(n int) TestableExpression {
	return tce.Arguments[n-1].(TestableExpression)
}

// ReceiverExpression returns call expression's receiver as TestingExpression
func (tce *TestableCallExpression) TestableReceiver() TestableExpression {
	return tce.Receiver.(TestableExpression)
}

// ShouldHaveMethodName checks if the method's name is same as we expected
func (tce *TestableCallExpression) ShouldHaveMethodName(expectedName string) {
	if tce.Method != expectedName {
		tce.t.Helper()
		tce.t.Fatalf("expect call expression's method name to be '%s', got '%s'", expectedName, tce.Method)
	}
}

// ShouldHaveNumbersOfArguments checks if the method call's argument number is same we expected
func (tce *TestableCallExpression) ShouldHaveNumbersOfArguments(n int) {
	if len(tce.Arguments) != n {
		tce.t.Helper()
		tce.t.Fatalf("expect call expression to have %d arguments, got %d", n, len(tce.Arguments))
	}
}

// TestableConditionalExpression
type TestableConditionalExpression struct {
	*ConditionalExpression
	t *testing.T
}

func (tce *TestableConditionalExpression) TestableCondition() TestableExpression {
	return tce.Condition.(TestableExpression)
}

func (tce *TestableConditionalExpression) TestableConsequence() CodeBlock {
	var tss []TestableStatement
	for _, stmt := range tce.Consequence.Statements {
		tss = append(tss, stmt.(TestableStatement))
	}

	return tss
}

// TestableConstant
type TestableConstant struct {
	*Constant
	t *testing.T
}

// ShouldHaveName checks if the constant's name is same as we expected
func (tc *TestableConstant) ShouldHaveName(expectedName string) {
	if tc.Value != expectedName {
		tc.t.Helper()
		tc.t.Fatalf("expect current identifier to be '%s', got '%s'", expectedName, tc.Value)
	}
}

// TestableHashExpression
type TestableHashExpression struct {
	*HashExpression
	t *testing.T
}

// TestableDataPairs returns a map of hash expression's element and assert them as TestableExpression
func (the *TestableHashExpression) TestableDataPairs() (pairs map[string]TestableExpression) {
	pairs = make(map[string]TestableExpression)
	for k, v := range the.Data {
		pairs[k] = v.(TestableExpression)
	}

	return
}

// TestableIdentifier
type TestableIdentifier struct {
	*Identifier
	t *testing.T
}

// ShouldHaveName checks if the identifier's name is same as we expected
func (ti *TestableIdentifier) ShouldHaveName(expectedName string) {
	if ti.Value != expectedName {
		ti.t.Helper()
		ti.t.Fatalf("expect current identifier to be '%s', got '%s'", expectedName, ti.Value)
	}
}

// TestableTernaryExpression
type TestableTernaryExpression struct {
	*TernaryExpression
	t *testing.T
}

// TODO: Test helpers

// TestableIfExpression
type TestableIfExpression struct {
	*IfExpression
	t *testing.T
}

// ShouldHaveNumberOfConditionals checks if the number of condition matches the specified one.
func (tie *TestableIfExpression) ShouldHaveNumberOfConditionals(n int) {
	if len(tie.Conditionals) != n {
		tie.t.Helper()
		tie.t.Fatalf("Expect if expression to have %d conditionals, got %d", n, len(tie.Conditionals))
	}
}

// TestableConditionals returns if expression's conditionals and assert them as TestableExpression
func (tie *TestableIfExpression) TestableConditionals() (tes []TestableExpression) {
	for _, cond := range tie.Conditionals {
		tes = append(tes, cond)
	}

	return
}

// TestableAlternative returns if expression's alternative code block as TestableExpression
func (tie *TestableIfExpression) TestableAlternative() CodeBlock {
	var tss []TestableStatement
	for _, stmt := range tie.Alternative.Statements {
		tss = append(tss, stmt.(TestableStatement))
	}

	return tss
}

// TestableInfixExpression
type TestableInfixExpression struct {
	*InfixExpression
	t *testing.T
}

// ShouldHaveOperator checks if the infix expression has expected operator
func (tie *TestableInfixExpression) ShouldHaveOperator(expectedOperator string) {
	if tie.Operator != expectedOperator {
		tie.t.Helper()
		tie.t.Fatalf("Expect infix expression to have %s operator, got %s", expectedOperator, tie.Operator)
	}
}

// LeftExpression returns infix expression's left expression as TestingExpression
func (tie *TestableInfixExpression) TestableLeftExpression() TestableExpression {
	return tie.Left.(TestableExpression)
}

// RightExpression returns infix expression's right expression as TestingExpression
func (tie *TestableInfixExpression) TestableRightExpression() TestableExpression {
	return tie.Right.(TestableExpression)
}

// TestableInstanceVariable
type TestableInstanceVariable struct {
	*InstanceVariable
	t *testing.T
}

// ShouldHaveName checks if the instance variable's name is same as we expected
func (tiv *TestableInstanceVariable) ShouldHaveName(expectedName string) {
	if tiv.Value != expectedName {
		tiv.t.Helper()
		tiv.t.Fatalf("expect current instance variable to be '%s', got '%s'", expectedName, tiv.Value)
	}
}

// TestableIntegerLiteral
type TestableIntegerLiteral struct {
	*IntegerLiteral
	t *testing.T
}

// ShouldEqualTo compares if the integer literal's value equals to the expected value
func (til *TestableIntegerLiteral) ShouldEqualTo(expectedInt int) {
	if til.Value != expectedInt {
		til.t.Helper()
		til.t.Fatalf("Expect integer literal to be %d, got %d", expectedInt, til.Value)
	}
}

// TestableSelfExpression
type TestableSelfExpression struct {
	*SelfExpression
	t *testing.T
}

// TestableStringLiteral
type TestableStringLiteral struct {
	*StringLiteral
	t *testing.T
}

// ShouldEqualTo compares if the string literal's value equals to the expected value
func (tsl *TestableStringLiteral) ShouldEqualTo(expected string) {
	if tsl.Value != expected {
		tsl.t.Helper()
		tsl.t.Fatalf("Expect string literal to be %s, got %s", expected, tsl.Value)
	}
}

// TestableYieldExpression
type TestableYieldExpression struct {
	*YieldExpression
	t *testing.T
}

// NthArgument returns n-th argument of the yield expression as TestingExpression
func (tye *TestableYieldExpression) NthArgument(n int) TestableExpression {
	return tye.Arguments[n-1].(TestableExpression)
}
