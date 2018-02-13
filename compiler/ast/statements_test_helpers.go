//+build !release

package ast

import (
	"testing"
)

/*
 BaseNode
*/

// IsClassStmt fails the test and returns nil by default
func (b *BaseNode) IsClassStmt(t *testing.T) *TestableClassStatement {
	t.Fatalf(nodeFailureMsgFormat, "class statement", b)
	return nil
}

// IsModuleStmt fails the test and returns nil by default
func (b *BaseNode) IsModuleStmt(t *testing.T) *TestableModuleStatement {
	t.Fatalf(nodeFailureMsgFormat, "module statement", b)
	return nil
}

// IsReturnStmt fails the test and returns nil by default
func (b *BaseNode) IsReturnStmt(t *testing.T) *TestableReturnStatement {
	t.Fatalf(nodeFailureMsgFormat, "return statement", b)
	return nil
}

// IsDefStmt fails the test and returns nil by default
func (b *BaseNode) IsDefStmt(t *testing.T) *TestableDefStatement {
	t.Fatalf(nodeFailureMsgFormat, "method definition", b)
	return nil
}

// IsWhileStmt fails the test and returns nil by default
func (b *BaseNode) IsWhileStmt(t *testing.T) (ws *TestableWhileStatement) {
	t.Fatalf(nodeFailureMsgFormat, "while statement", b)
	return nil
}

/*
 ClassStatement
*/

func (cs *ClassStatement) IsClassStmt(t *testing.T) *TestableClassStatement {
	return &TestableClassStatement{t: t, ClassStatement: cs}
}

/*
 Module Statement
*/

// IsModuleStmt returns a pointer of the module statement
func (ms *ModuleStatement) IsModuleStmt(t *testing.T) *TestableModuleStatement {
	return &TestableModuleStatement{ModuleStatement: ms, t: t}
}

// IsDefStmt returns a pointer of the DefStatement
func (ds *DefStatement) IsDefStmt(t *testing.T) *TestableDefStatement {
	return &TestableDefStatement{DefStatement: ds, t: t}
}

// IsDefStmt returns a pointer of the ReturnStatement
func (rs *ReturnStatement) IsReturnStmt(t *testing.T) (trs *TestableReturnStatement) {
	return &TestableReturnStatement{t: t, ReturnStatement: rs}
}

// IsExpressionStmt returns ExpressionStatement itself
func (ts *ExpressionStatement) IsExpression(t *testing.T) TestableExpression {
	return ts.Expression.(TestableExpression)
}

// IsWhileStmt returns the pointer of current while statement
func (ws *WhileStatement) IsWhileStmt(t *testing.T) *TestableWhileStatement {
	return &TestableWhileStatement{WhileStatement: ws, t: t}
}
