//+build !release

package ast

import (
	"testing"
)

const nodeFailureMsgFormat = "Node is not %s, is %T"

// TestableIdentifierValue for marking a string as an identifier's value for test
type TestableIdentifierValue string

/*
 BaseNode
*/

// IsArrayExpression fails the test and returns nil by default
func (b *BaseNode) IsArrayExpression(t *testing.T) *testableArrayExpression {
	t.Helper()
	t.Fatalf(nodeFailureMsgFormat, "array expression", b)
	return nil
}

// IsAssignExpression fails the test and returns nil by default
func (b *BaseNode) IsAssignExpression(t *testing.T) (ae *testableAssignExpression) {
	t.Helper()
	t.Fatalf(nodeFailureMsgFormat, "assign expression", b)
	return
}

// IsBooleanExpression fails the test and returns nil by default
func (b *BaseNode) IsBooleanExpression(t *testing.T) (ae *testableBooleanExpression) {
	t.Helper()
	t.Fatalf(nodeFailureMsgFormat, "boolean expression", b)
	return
}

// IsCallExpression fails the test and returns nil by default
func (b *BaseNode) IsCallExpression(t *testing.T) (ce *testableCallExpression) {
	t.Helper()
	t.Fatalf(nodeFailureMsgFormat, "call expression", b)
	return
}

// IsConditionalExpression fails the test and returns nil by default
func (b *BaseNode) IsConditionalExpression(t *testing.T) *testableConditionalExpression {
	t.Helper()
	t.Fatalf(nodeFailureMsgFormat, "conditional expression", b)
	return nil
}

// IsConstant fails the test and returns nil by default
func (b *BaseNode) IsConstant(t *testing.T) (c *testableConstant) {
	t.Helper()
	t.Fatalf(nodeFailureMsgFormat, "constant", b)
	return
}

// IsExpression fails the test and returns nil by default
func (b *BaseNode) IsExpression(t *testing.T) (te testableExpression) {
	t.Helper()
	t.Fatalf(nodeFailureMsgFormat, "expression", b)
	return
}

// IsHashExpression fails the test and returns nil by default
func (b *BaseNode) IsHashExpression(t *testing.T) *testableHashExpression {
	t.Helper()
	t.Fatalf(nodeFailureMsgFormat, "hash expression", b)
	return nil
}

// IsIdentifier fails the test and returns nil by default
func (b *BaseNode) IsIdentifier(t *testing.T) *testableIdentifier {
	t.Helper()
	t.Fatalf(nodeFailureMsgFormat, "identifier", b)
	return nil
}

// IsIfExpression fails the test and returns nil by default
func (b *BaseNode) IsIfExpression(t *testing.T) *testableIfExpression {
	t.Helper()
	t.Fatalf(nodeFailureMsgFormat, "if expression", b)
	return nil
}

// IsInfixExpression fails the test and returns nil by default
func (b *BaseNode) IsInfixExpression(t *testing.T) (ie *testableInfixExpression) {
	t.Helper()
	t.Fatalf(nodeFailureMsgFormat, "infix expression", b)
	return
}

// IsInstanceVariable fails the test and returns nil by default
func (b *BaseNode) IsInstanceVariable(t *testing.T) (ie *testableInstanceVariable) {
	t.Helper()
	t.Fatalf(nodeFailureMsgFormat, "instance variable", b)
	return
}

// IsIntegerLiteral fails the test and returns nil by default
func (b *BaseNode) IsIntegerLiteral(t *testing.T) (il *testableIntegerLiteral) {
	t.Helper()
	t.Fatalf(nodeFailureMsgFormat, "integer literal", b)
	return
}

// IsSelfExpression fails the test and returns nil by default
func (b *BaseNode) IsSelfExpression(t *testing.T) (sl *testableSelfExpression) {
	t.Helper()
	t.Fatalf(nodeFailureMsgFormat, "self expression", b)
	return
}

// IsStringLiteral fails the test and returns nil by default
func (b *BaseNode) IsStringLiteral(t *testing.T) *testableStringLiteral {
	t.Helper()
	t.Fatalf(nodeFailureMsgFormat, "string literal", b)
	return nil
}

// IsYieldExpression returns pointer of the receiver yield expression
func (b *BaseNode) IsYieldExpression(t *testing.T) *testableYieldExpression {
	t.Helper()
	t.Fatalf(nodeFailureMsgFormat, "yield expression", b)
	return nil
}

/*
AST node's helpers
*/

// IsArrayExpression returns pointer of the receiver array expression
func (ae *ArrayExpression) IsArrayExpression(t *testing.T) *testableArrayExpression {
	return &testableArrayExpression{ArrayExpression: ae, t: t}
}

// IsAssignExpression returns pointer of the receiver assign expression
func (ae *AssignExpression) IsAssignExpression(t *testing.T) *testableAssignExpression {
	return &testableAssignExpression{AssignExpression: ae, t: t}
}

// IsBooleanExpression returns pointer of the receiver boolean expression
func (be *BooleanExpression) IsBooleanExpression(t *testing.T) *testableBooleanExpression {
	return &testableBooleanExpression{BooleanExpression: be, t: t}
}

// IsCallExpression returns pointer of the receiver call expression
func (ce *CallExpression) IsCallExpression(t *testing.T) *testableCallExpression {
	return &testableCallExpression{CallExpression: ce, t: t}
}

// IsConditionalExpression returns pointer of the receiver conditional expression
func (ce *ConditionalExpression) IsConditionalExpression(t *testing.T) *testableConditionalExpression {
	return &testableConditionalExpression{ConditionalExpression: ce, t: t}
}

// IsConstant returns pointer of the current receiver constant
func (c *Constant) IsConstant(t *testing.T) *testableConstant {
	return &testableConstant{Constant: c, t: t}
}

// IsHashExpression returns pointer of the receiver hash expression
func (he *HashExpression) IsHashExpression(t *testing.T) *testableHashExpression {
	return &testableHashExpression{HashExpression: he, t: t}
}

// IsIdentifier returns pointer of the receiver identifier
func (i *Identifier) IsIdentifier(t *testing.T) *testableIdentifier {
	return &testableIdentifier{Identifier: i, t: t}
}

// IsIfExpression returns pointer of the receiver if expression
func (ie *IfExpression) IsIfExpression(t *testing.T) *testableIfExpression {
	return &testableIfExpression{IfExpression: ie, t: t}
}

// IsInfixExpression returns pointer of the receiver infix expression
func (ie *InfixExpression) IsInfixExpression(t *testing.T) *testableInfixExpression {
	return &testableInfixExpression{InfixExpression: ie, t: t}
}

// IsInstanceVariable returns pointer of the receiver instance variable
func (iv *InstanceVariable) IsInstanceVariable(t *testing.T) *testableInstanceVariable {
	return &testableInstanceVariable{InstanceVariable: iv, t: t}
}

// IsIntegerLiteral returns pointer of the receiver integer literal
func (il *IntegerLiteral) IsIntegerLiteral(t *testing.T) *testableIntegerLiteral {
	return &testableIntegerLiteral{IntegerLiteral: il, t: t}
}

// IsSelfExpression returns pointer of the receiver self expression
func (se *SelfExpression) IsSelfExpression(t *testing.T) *testableSelfExpression {
	return &testableSelfExpression{SelfExpression: se, t: t}
}

// IsStringLiteral returns pointer of the receiver string literal
func (sl *StringLiteral) IsStringLiteral(t *testing.T) *testableStringLiteral {
	return &testableStringLiteral{StringLiteral: sl, t: t}
}

// IsYieldExpression returns pointer of the receiver yield expression
func (ye *YieldExpression) IsYieldExpression(t *testing.T) *testableYieldExpression {
	return &testableYieldExpression{YieldExpression: ye, t: t}
}
