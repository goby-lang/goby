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
func (b *BaseNode) IsWhileStmt(t *testing.T) (ws *WhileStatement) {
	t.Fatalf(nodeFailureMsgFormat, "while statement", b)
	return nil
}

/*
 ClassStatement
*/

func (cs *ClassStatement) IsClassStmt(t *testing.T) *TestableClassStatement {
	return &TestableClassStatement{t: t, ClassStatement: cs}
}

// NameIs returns the compare result of current class name and target class name
func (cs *ClassStatement) NameIs(n string) bool {
	if cs.Name.Value == n {
		return true
	}

	return false
}

/*
 Module Statement
*/

// IsModuleStmt returns a pointer of the module statement
func (ms *ModuleStatement) IsModuleStmt(t *testing.T) *TestableModuleStatement {
	return &TestableModuleStatement{ModuleStatement: ms, t: t}
}

// NameIs returns the compare result of current module name and target module name
func (ms *ModuleStatement) NameIs(n string) bool {
	if ms.Name.Value == n {
		return true
	}

	return false
}

/*
 DefStatement
*/

// IsDefStmt returns a pointer of the DefStatement
func (ds *DefStatement) IsDefStmt(t *testing.T) *TestableDefStatement {
	return &TestableDefStatement{DefStatement: ds, t: t}
}

/*
ReturnStatement
*/

func (rs *ReturnStatement) IsReturnStmt(t *testing.T) (trs *TestableReturnStatement) {
	return &TestableReturnStatement{t: t, ReturnStatement: rs}
}

/*
ExpressionStatement
*/

// IsExpressionStmt returns ExpressionStatement itself
func (ts *ExpressionStatement) IsExpression(t *testing.T) TestingExpression {
	return ts.Expression.(TestingExpression)
}

/*
WhileStatement
*/

// Block returns while statement's code block as a set of TestingStatements
func (we *WhileStatement) CodeBlock() CodeBlock {
	var tss []TestingStatement

	for _, stmt := range we.Body.Statements {
		tss = append(tss, stmt.(TestingStatement))
	}

	return tss
}

// ConditionExpression returns while statement's condition as TestingExpression
func (we *WhileStatement) ConditionExpression() TestingExpression {
	return we.Condition.(TestingExpression)
}

// IsWhileStmt returns the pointer of current while statement
func (ws *WhileStatement) IsWhileStmt(t *testing.T) *WhileStatement {
	return ws
}
