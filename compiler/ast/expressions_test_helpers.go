//+build !release

package ast

import "testing"

/*
 BaseNode
*/

// IsAssignExpression fails the test and returns nil by default
func (b *BaseNode) IsAssignExpression(t *testing.T) (ae *AssignExpression) {
	t.Fatalf("Node is not an assign expression, is %v", b)
	return
}

// IsCallExpression fails the test and returns nil by default
func (b *BaseNode) IsCallExpression(t *testing.T) (ce *CallExpression) {
	t.Fatalf("Node is not a call expression, is %v", b)
	return
}

// IsExpression fails the test and returns nil by default
func (b *BaseNode) IsExpression(t *testing.T) (te TestingExpression) {
	t.Fatalf("Node is not an expression, is %v", b)
	return
}

// IsIdentifier fails the test and returns nil by default
func (b *BaseNode) IsIdentifier(t *testing.T) (i *Identifier) {
	t.Fatalf("Node is not an identifier, is %v", b)
	return
}

// IsInfixExpression fails the test and returns nil by default
func (b *BaseNode) IsInfixExpression(t *testing.T) (ie *InfixExpression) {
	t.Fatalf("Node is not an infix expression, is %v", b)
	return
}

// IsYieldExpression returns pointer of the receiver yield expression
func (b *BaseNode) IsYieldExpression(t *testing.T) (ye *YieldExpression) {
	t.Fatalf("Node is not an yield expression, is %v", b)
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
func (ce *CallExpression) IsCallExpression(t *testing.T) *CallExpression {
	return ce
}

// ShouldHasMethodName
func (ce *CallExpression) ShouldHasMethodName(t *testing.T, expectedName string) {
	if ce.Method != expectedName {
		t.Fatalf("expect call expression's method name to be '%s', got '%s'", expectedName, ce.Method)
	}
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
YieldExpression
*/

// IsYieldExpression returns pointer of the receiver yield expression
func (ye *YieldExpression) IsYieldExpression(t *testing.T) *YieldExpression {
	return ye
}
