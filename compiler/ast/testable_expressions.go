//+build !release

package ast

import "testing"

// testableExpression interface implements basic Expression's functions, and helper functions to assert node's type
type testableExpression interface {
	Expression
	// Test Helpers
	IsArrayExpression(t *testing.T) *testableArrayExpression
	IsAssignExpression(t *testing.T) *testableAssignExpression
	IsBooleanExpression(t *testing.T) *testableBooleanExpression
	IsCallExpression(t *testing.T) *testableCallExpression
	IsConditionalExpression(t *testing.T) *testableConditionalExpression
	IsConstant(t *testing.T) *testableConstant
	IsHashExpression(t *testing.T) *testableHashExpression
	IsIdentifier(t *testing.T) *testableIdentifier
	IsIfExpression(t *testing.T) *testableIfExpression
	IsInfixExpression(t *testing.T) *testableInfixExpression
	IsInstanceVariable(t *testing.T) *testableInstanceVariable
	IsIntegerLiteral(t *testing.T) *testableIntegerLiteral
	IsSelfExpression(t *testing.T) *testableSelfExpression
	IsStringLiteral(t *testing.T) *testableStringLiteral
	IsYieldExpression(t *testing.T) *testableYieldExpression
}

type testableArrayExpression struct {
	*ArrayExpression
	t *testing.T
}

// TestableElements returns array expression's element nodes and assert them as testableExpression
func (tae *testableArrayExpression) TestableElements() (tes []testableExpression) {
	for _, elem := range tae.Elements {
		tes = append(tes, elem.(testableExpression))
	}

	return
}

type testableAssignExpression struct {
	*AssignExpression
	t *testing.T
}

// NthVariable returns the nth variable of the assignment as a testableExpression
func (tae *testableAssignExpression) NthVariable(n int) testableExpression {
	return tae.Variables[n-1].(testableExpression)
}

// TestableValue returns the assignment's value as a testableExpression
func (tae *testableAssignExpression) TestableValue() testableExpression {
	return tae.Value.(testableExpression)
}

type testableBooleanExpression struct {
	*BooleanExpression
	t *testing.T
}

// ShouldEqualTo compares if the boolean expression's value equals to the expected value
func (tbe *testableBooleanExpression) ShouldEqualTo(expected bool) {
	if tbe.Value != expected {
		tbe.t.Helper()
		tbe.t.Fatalf("Expect boolean literal to be %t, got %t", expected, tbe.Value)
	}
}

type testableCallExpression struct {
	*CallExpression
	t *testing.T
}

// NthArgument returns n-th argument of the call expression as TestingExpression
func (tce *testableCallExpression) NthArgument(n int) testableExpression {
	return tce.Arguments[n-1].(testableExpression)
}

// ReceiverExpression returns call expression's receiver as TestingExpression
func (tce *testableCallExpression) TestableReceiver() testableExpression {
	return tce.Receiver.(testableExpression)
}

// ShouldHaveMethodName checks if the method's name is same as we expected
func (tce *testableCallExpression) ShouldHaveMethodName(expectedName string) {
	if tce.Method != expectedName {
		tce.t.Helper()
		tce.t.Fatalf("expect call expression's method name to be '%s', got '%s'", expectedName, tce.Method)
	}
}

// ShouldHaveNumbersOfArguments checks if the method call's argument number is same we expected
func (tce *testableCallExpression) ShouldHaveNumbersOfArguments(n int) {
	if len(tce.Arguments) != n {
		tce.t.Helper()
		tce.t.Fatalf("expect call expression to have %d arguments, got %d", n, len(tce.Arguments))
	}
}

type testableConditionalExpression struct {
	*ConditionalExpression
	t *testing.T
}

func (tce *testableConditionalExpression) TestableCondition() testableExpression {
	return tce.Condition.(testableExpression)
}

func (tce *testableConditionalExpression) TestableConsequence() CodeBlock {
	var tss []TestableStatement
	for _, stmt := range tce.Consequence.Statements {
		tss = append(tss, stmt.(TestableStatement))
	}

	return tss
}

type testableConstant struct {
	*Constant
	t *testing.T
}

// ShouldHaveName checks if the constant's name is same as we expected
func (tc *testableConstant) ShouldHaveName(expectedName string) {
	if tc.Value != expectedName {
		tc.t.Helper()
		tc.t.Fatalf("expect current identifier to be '%s', got '%s'", expectedName, tc.Value)
	}
}

type testableHashExpression struct {
	*HashExpression
	t *testing.T
}

// TestableDataPairs returns a map of hash expression's element and assert them as testableExpression
func (the *testableHashExpression) TestableDataPairs() (pairs map[string]testableExpression) {
	pairs = make(map[string]testableExpression)
	for k, v := range the.Data {
		pairs[k] = v.(testableExpression)
	}

	return
}

type testableIdentifier struct {
	*Identifier
	t *testing.T
}

// ShouldHaveName checks if the identifier's name is same as we expected
func (ti *testableIdentifier) ShouldHaveName(expectedName string) {
	if ti.Value != expectedName {
		ti.t.Helper()
		ti.t.Fatalf("expect current identifier to be '%s', got '%s'", expectedName, ti.Value)
	}
}

type testableIfExpression struct {
	*IfExpression
	t *testing.T
}

// ShouldHaveNumberOfConditionals checks if the number of condition matches the specified one.
func (tie *testableIfExpression) ShouldHaveNumberOfConditionals(n int) {
	if len(tie.Conditionals) != n {
		tie.t.Helper()
		tie.t.Fatalf("Expect if expression to have %d conditionals, got %d", n, len(tie.Conditionals))
	}
}

// TestableConditionals returns if expression's conditionals and assert them as testableExpression
func (tie *testableIfExpression) TestableConditionals() (tes []testableExpression) {
	for _, cond := range tie.Conditionals {
		tes = append(tes, cond)
	}

	return
}

// TestableAlternative returns if expression's alternative code block as testableExpression
func (tie *testableIfExpression) TestableAlternative() CodeBlock {
	var tss []TestableStatement
	for _, stmt := range tie.Alternative.Statements {
		tss = append(tss, stmt.(TestableStatement))
	}

	return tss
}

type testableInfixExpression struct {
	*InfixExpression
	t *testing.T
}

// ShouldHaveOperator checks if the infix expression has expected operator
func (tie *testableInfixExpression) ShouldHaveOperator(expectedOperator string) {
	if tie.Operator != expectedOperator {
		tie.t.Helper()
		tie.t.Fatalf("Expect infix expression to have %s operator, got %s", expectedOperator, tie.Operator)
	}
}

// LeftExpression returns infix expression's left expression as TestingExpression
func (tie *testableInfixExpression) TestableLeftExpression() testableExpression {
	return tie.Left.(testableExpression)
}

// RightExpression returns infix expression's right expression as TestingExpression
func (tie *testableInfixExpression) TestableRightExpression() testableExpression {
	return tie.Right.(testableExpression)
}

type testableInstanceVariable struct {
	*InstanceVariable
	t *testing.T
}

// ShouldHaveName checks if the instance variable's name is same as we expected
func (tiv *testableInstanceVariable) ShouldHaveName(expectedName string) {
	if tiv.Value != expectedName {
		tiv.t.Helper()
		tiv.t.Fatalf("expect current instance variable to be '%s', got '%s'", expectedName, tiv.Value)
	}
}

type testableIntegerLiteral struct {
	*IntegerLiteral
	t *testing.T
}

// ShouldEqualTo compares if the integer literal's value equals to the expected value
func (til *testableIntegerLiteral) ShouldEqualTo(expectedInt int) {
	if til.Value != expectedInt {
		til.t.Helper()
		til.t.Fatalf("Expect integer literal to be %d, got %d", expectedInt, til.Value)
	}
}

type testableSelfExpression struct {
	*SelfExpression
	t *testing.T
}

type testableStringLiteral struct {
	*StringLiteral
	t *testing.T
}

// ShouldEqualTo compares if the string literal's value equals to the expected value
func (tsl *testableStringLiteral) ShouldEqualTo(expected string) {
	if tsl.Value != expected {
		tsl.t.Helper()
		tsl.t.Fatalf("Expect string literal to be %s, got %s", expected, tsl.Value)
	}
}

type testableYieldExpression struct {
	*YieldExpression
	t *testing.T
}

// NthArgument returns n-th argument of the yield expression as TestingExpression
func (tye *testableYieldExpression) NthArgument(n int) testableExpression {
	return tye.Arguments[n-1].(testableExpression)
}
