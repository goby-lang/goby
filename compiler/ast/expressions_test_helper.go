//+build !release

package ast

import "testing"


// IsYieldExpression returns pointer of the receiver yield expression
func (b *BaseNode) IsYieldExpression(t *testing.T) (ye *YieldExpression) {
	t.Fatalf("Node is not an yield expression, is %v", b)
	return
}

/*
AssignExpression
*/

// NameIs compares the assignment's variable name and expected name
func (ae *AssignExpression) NameIs(n string) bool {
	return ae.Variables[0].(testingNode).NameIs(n)
}

/*
Identifier
*/

// NameIs compares the identifier's name and expected name
func (i *Identifier) NameIs(n string) bool {
	if i.Value == n {
		return true
	}

	return false
}

/*
YieldExpression
*/

// IsYieldExpression returns pointer of the receiver yield expression
func (ye *YieldExpression) IsYieldExpression(t *testing.T) *YieldExpression {
	return ye
}