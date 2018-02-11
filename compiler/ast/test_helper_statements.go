//+build !release

package ast

import "testing"

type TestingStatement interface {
	Statement
	// Test Helpers
	IsClassStmt(t *testing.T, className string) *ClassStatement
	IsDefStmt(t *testing.T) *DefStatement
	IsExpression(t *testing.T) TestingExpression
	IsModuleStmt(t *testing.T, className string) *ModuleStatement
	IsReturnStmt(t *testing.T) *TestableReturnStatement
	IsWhileStmt(t *testing.T) *WhileStatement
}

type CodeBlock []TestingStatement

func (cb CodeBlock) NthStmt(n int) TestingStatement {
	return cb[n-1]
}

/*TestableReturnStatement*/

type TestableReturnStatement struct {
	*ReturnStatement
	t *testing.T
}

func (trs *TestableReturnStatement) ShouldHasValue(value interface{}) {
	t := trs.t
	rs := trs.ReturnStatement
	switch v := value.(type) {
	case int:
		compareInt(t, rs.ReturnValue, v)
	case string:
		compareString(t, rs.ReturnValue, v)
	case bool:
		compareBool(t, rs.ReturnValue, v)
	case TestingIdentifier:
		compareIdentifier(t, rs.ReturnValue, v)
	}
}