//+build !release

package ast

import (
	"testing"
)

/*
 BaseNode
*/

// IsClassStmt fails the test and returns nil by default
func (b *BaseNode) IsClassStmt(t *testing.T) *testableClassStatement {
	t.Helper()
	t.Fatalf(nodeFailureMsgFormat, "class statement", b)
	return nil
}

// IsModuleStmt fails the test and returns nil by default
func (b *BaseNode) IsModuleStmt(t *testing.T) *testableModuleStatement {
	t.Helper()
	t.Fatalf(nodeFailureMsgFormat, "module statement", b)
	return nil
}

// IsReturnStmt fails the test and returns nil by default
func (b *BaseNode) IsReturnStmt(t *testing.T) *testableReturnStatement {
	t.Helper()
	t.Fatalf(nodeFailureMsgFormat, "return statement", b)
	return nil
}

// IsDefStmt fails the test and returns nil by default
func (b *BaseNode) IsDefStmt(t *testing.T) *testableDefStatement {
	t.Helper()
	t.Fatalf(nodeFailureMsgFormat, "method definition", b)
	return nil
}

// IsWhileStmt fails the test and returns nil by default
func (b *BaseNode) IsWhileStmt(t *testing.T) (ws *testableWhileStatement) {
	t.Helper()
	t.Fatalf(nodeFailureMsgFormat, "while statement", b)
	return nil
}

// IsClassStmt fails the test and returns nil by default
func (cs *ClassStatement) IsClassStmt(t *testing.T) *testableClassStatement {
	return &testableClassStatement{t: t, ClassStatement: cs}
}

// IsModuleStmt returns a pointer of the module statement
func (ms *ModuleStatement) IsModuleStmt(t *testing.T) *testableModuleStatement {
	return &testableModuleStatement{ModuleStatement: ms, t: t}
}

// IsDefStmt returns a pointer of the DefStatement
func (ds *DefStatement) IsDefStmt(t *testing.T) *testableDefStatement {
	return &testableDefStatement{DefStatement: ds, t: t}
}

// IsReturnStmt returns a pointer of the ReturnStatement
func (rs *ReturnStatement) IsReturnStmt(t *testing.T) (trs *testableReturnStatement) {
	return &testableReturnStatement{t: t, ReturnStatement: rs}
}

// IsExpression returns ExpressionStatement itself
func (ts *ExpressionStatement) IsExpression(t *testing.T) testableExpression {
	return ts.Expression.(testableExpression)
}

// IsWhileStmt returns the pointer of current while statement
func (ws *WhileStatement) IsWhileStmt(t *testing.T) *testableWhileStatement {
	return &testableWhileStatement{WhileStatement: ws, t: t}
}
