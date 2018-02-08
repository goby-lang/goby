//+build !release

package ast

/*
AssignExpression
*/

// NameIs compares the assignment's variable name and expected name
func (ae *AssignExpression) NameIs(n string) bool {
	return ae.Variables[0].(testNode).NameIs(n)
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
