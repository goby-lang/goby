//+build !release

package ast

import "testing"

type TestingStatement interface {
	Statement
	// Test Helpers
	IsClassStmt(t *testing.T) *TestableClassStatement
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

/*TestableClassStatement*/

type TestableClassStatement struct {
	*ClassStatement
	t *testing.T
}

// HasMethod checks if current class statement has target method, and returns the method if it has
func (tcs *TestableClassStatement) HasMethod(methodName string) (ds *DefStatement) {
	for _, stmt := range tcs.Body.Statements {
		s, ok := stmt.(*DefStatement)

		if ok && s.Name.Value == methodName {
			ds = s
			return
		}
	}

	tcs.t.Fatalf("Can't find method '%s' in class '%s'", methodName, tcs.Name)
	return
}

func (tcs *TestableClassStatement) ShouldHasName(name string) {
	if tcs.Name.Value != name {
		tcs.t.Fatalf("Wrong class, this class is %s", tcs.Name.Value)
	}
}

// ShouldInherits checks if current class statement inherits the target class
func (tcs *TestableClassStatement) ShouldInherits(className string) {
	if tcs.SuperClassName != className {
		tcs.t.Fatalf("Expect class %s to inherit class %s. got %s", tcs.Name, className, tcs.SuperClassName)
	}
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
