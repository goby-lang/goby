//+build !release

package ast

import (
	"testing"
)

/*
 BaseNode
*/

// IsAssignExpression fails the test and returns nil by default
func (b *BaseNode) IsArrayExpression(t *testing.T) (ae *ArrayExpression) {
	t.Fatalf(nodeFailureMsgFormat, "array expression", b)
	return
}

// IsAssignExpression fails the test and returns nil by default
func (b *BaseNode) IsAssignExpression(t *testing.T) (ae *AssignExpression) {
	t.Fatalf(nodeFailureMsgFormat, "assign expression", b)
	return
}

// IsCallExpression fails the test and returns nil by default
func (b *BaseNode) IsCallExpression(t *testing.T) (ce *TestableCallExpression) {
	t.Fatalf(nodeFailureMsgFormat, "call expression", b)
	return
}

// IsConditionalExpression fails the test and returns nil by default
func (b *BaseNode) IsConditionalExpression(t *testing.T) (ce *ConditionalExpression) {
	t.Fatalf(nodeFailureMsgFormat, "conditional expression", b)
	return
}

// IsConstant fails the test and returns nil by default
func (b *BaseNode) IsConstant(t *testing.T) (c *Constant) {
	t.Fatalf(nodeFailureMsgFormat, "constant", b)
	return
}

// IsExpression fails the test and returns nil by default
func (b *BaseNode) IsExpression(t *testing.T) (te TestingExpression) {
	t.Helper()
	t.Fatalf(nodeFailureMsgFormat, "expression", b)
	return
}

// HashExpression fails the test and returns nil by default
func (b *BaseNode) IsHashExpression(t *testing.T) (he *HashExpression) {
	t.Fatalf(nodeFailureMsgFormat, "hash expression", b)
	return
}

// IsIdentifier fails the test and returns nil by default
func (b *BaseNode) IsIdentifier(t *testing.T) (i *Identifier) {
	t.Fatalf(nodeFailureMsgFormat, "identifier", b)
	return
}

// IsIfExpression fails the test and returns nil by default
func (b *BaseNode) IsIfExpression(t *testing.T) (i *IfExpression) {
	t.Fatalf(nodeFailureMsgFormat, "if expression", b)
	return
}

// IsInfixExpression fails the test and returns nil by default
func (b *BaseNode) IsInfixExpression(t *testing.T) (ie *InfixExpression) {
	t.Fatalf(nodeFailureMsgFormat, "infix expression", b)
	return
}

// IsIntegerLiteral fails the test and returns nil by default
func (b *BaseNode) IsIntegerLiteral(t *testing.T) (il *IntegerLiteral) {
	t.Fatalf(nodeFailureMsgFormat, "integer literal", b)
	return
}

// IsStringLiteral fails the test and returns nil by default
func (b *BaseNode) IsStringLiteral(t *testing.T) (sl *StringLiteral) {
	t.Fatalf(nodeFailureMsgFormat, "string literal", b)
	return
}

// IsYieldExpression returns pointer of the receiver yield expression
func (b *BaseNode) IsYieldExpression(t *testing.T) (ye *YieldExpression) {
	t.Fatalf(nodeFailureMsgFormat, "yield expression", b)
	return
}

/*
ArrayExpression
*/

// IsArrayExpression returns pointer of the receiver array expression
func (ae *ArrayExpression) IsArrayExpression(t *testing.T) *ArrayExpression {
	return ae
}

func (ae *ArrayExpression) TestableElements() (tes []TestingExpression) {
	for _, elem := range ae.Elements {
		tes = append(tes, elem.(TestingExpression))
	}

	return
}

/*
AssignExpression
*/

// IsAssignExpression returns pointer of the receiver assign expression
func (ae *AssignExpression) IsAssignExpression(t *testing.T) *AssignExpression {
	return ae
}

// NameIs compares the assignment's variable name and expected name
func (ae *AssignExpression) NameIs(n string) bool {
	return ae.Variables[0].(testingNode).NameIs(n)
}

/*
CallExpression
*/

// IsCallExpression returns pointer of the receiver call expression
func (ce *CallExpression) IsCallExpression(t *testing.T) *TestableCallExpression {
	return &TestableCallExpression{CallExpression: ce, t: t}
}

/*
ConditionalExpression
*/

func (ce *ConditionalExpression) IsConditionalExpression(t *testing.T) *ConditionalExpression {
	return ce
}

func (ce *ConditionalExpression) TestableConsequence() CodeBlock {
	var tss []TestingStatement
	for _, stmt := range ce.Consequence.Statements {
		tss = append(tss, stmt.(TestingStatement))
	}

	return tss
}

/*
Constant
*/

// IsConstant returns pointer of the current receiver constant
func (c *Constant) IsConstant(t *testing.T) *Constant {
	return c
}

func (c *Constant) ShouldHasName(t *testing.T, expectedName string) {
	if c.Value != expectedName {
		t.Fatalf("expect current identifier to be '%s', got '%s'", expectedName, c.Value)
	}
}

/*
HashExpression
*/

// IsHashExpression returns pointer of the receiver hash expression
func (he *HashExpression) IsHashExpression(t *testing.T) *HashExpression {
	return he
}

/*
Identifier
*/

// IsIdentifier returns pointer of the receiver identifier
func (i *Identifier) IsIdentifier(t *testing.T) *Identifier {
	return i
}

func (i *Identifier) ShouldHasName(t *testing.T, expectedName string) {
	if i.Value != expectedName {
		t.Fatalf("expect current identifier to be '%s', got '%s'", expectedName, i.Value)
	}
}

// NameIs compares the identifier's name and expected name
func (i *Identifier) NameIs(n string) bool {
	if i.Value == n {
		return true
	}

	return false
}

/*
IfExpression
*/

// IsIfExpression returns pointer of the receiver if expression
func (ie *IfExpression) IsIfExpression(t *testing.T) *IfExpression {
	return ie
}

func (ie *IfExpression) ShouldHasNumberOfConditionals(t *testing.T, n int) {
	if len(ie.Conditionals) != n {
		t.Fatalf("Expect if expression to have %d conditionals, got %d", n, len(ie.Conditionals))
	}
}

func (ie *IfExpression) TestableConditionals() (tes []TestingExpression) {
	for _, cond := range ie.Conditionals {
		tes = append(tes, cond)
	}

	return
}

func (ie *IfExpression) TestableAlternative() CodeBlock {
	var tss []TestingStatement
	for _, stmt := range ie.Alternative.Statements {
		tss = append(tss, stmt.(TestingStatement))
	}

	return tss
}

/*
InfixExpression
*/

// IsInfixExpression returns pointer of the receiver infix expression
func (ie *InfixExpression) IsInfixExpression(t *testing.T) *InfixExpression {
	return ie
}

// ShouldHasOperator checks if the infix expression has expected operator
func (ie *InfixExpression) ShouldHasOperator(t *testing.T, expectedOperator string) {
	if ie.Operator != expectedOperator {
		t.Fatalf("Expect infix expression to have %s operator, got %s", expectedOperator, ie.Operator)
	}
}

// LeftExpression returns infix expression's left expression as TestingExpression
func (ie *InfixExpression) LeftExpression() TestingExpression {
	return ie.Left.(TestingExpression)
}

// RightExpression returns infix expression's right expression as TestingExpression
func (ie *InfixExpression) RightExpression() TestingExpression {
	return ie.Right.(TestingExpression)
}

/*
IntegerLiteral
*/

// IsIntegerLiteral returns pointer of the receiver string literal
func (il *IntegerLiteral) IsIntegerLiteral(t *testing.T) *IntegerLiteral {
	return il
}

func (il *IntegerLiteral) ShouldEqualTo(t *testing.T, expectedInt int) {
	if il.Value != expectedInt {
		t.Fatalf("Expect integer literal to be %d, got %d", expectedInt, il.Value)
	}
}

/*
StringLiteral
*/

// IsStringLiteral returns pointer of the receiver string literal
func (sl *StringLiteral) IsStringLiteral(t *testing.T) *StringLiteral {
	return sl
}

/*
YieldExpression
*/

// IsYieldExpression returns pointer of the receiver yield expression
func (ye *YieldExpression) IsYieldExpression(t *testing.T) *YieldExpression {
	return ye
}
