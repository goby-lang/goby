//+build !release

package ast

import (
	"testing"
)

/*
 BaseNode
*/

// IsClassStmt fails the test and returns nil by default
func (b *BaseNode) IsClassStmt(t *testing.T, className string) (cs *ClassStatement) {
	t.Fatalf(nodeFailureMsgFormat, "class statement", b)
	return
}

// IsModuleStmt fails the test and returns nil by default
func (b *BaseNode) IsModuleStmt(t *testing.T, moduleName string) (cs *ModuleStatement) {
	t.Fatalf(nodeFailureMsgFormat, "module statement", b)
	return
}

// IsReturnStmt fails the test and returns nil by default
func (b *BaseNode) IsReturnStmt(t *testing.T) (trs *TestableReturnStatement) {
	t.Fatalf(nodeFailureMsgFormat, "return statement", b)
	return
}

// IsDefStmt fails the test and returns nil by default
func (b *BaseNode) IsDefStmt(t *testing.T) (rs *DefStatement) {
	t.Fatalf(nodeFailureMsgFormat, "method definition", b)
	return
}

// IsWhileStmt fails the test and returns nil by default
func (b *BaseNode) IsWhileStmt(t *testing.T) (ws *WhileStatement) {
	t.Fatalf(nodeFailureMsgFormat, "while statement", b)
	return
}

/*
 ClassStatement
*/

// IsClassStmt returns a pointer of the class statement
func (cs *ClassStatement) IsClassStmt(t *testing.T, className string) *ClassStatement {
	return cs
}

// ShouldInherits checks if current class statement inherits the target class
func (cs *ClassStatement) ShouldInherits(t *testing.T, className string) {
	if cs.SuperClassName != className {
		t.Fatalf("Expect class %s to inherit class %s. got %s", cs.Name, className, cs.SuperClassName)
	}
}

// HasMethod checks if current class statement has target method, and returns the method if it has
func (cs *ClassStatement) HasMethod(t *testing.T, methodName string) (ds *DefStatement) {
	for _, stmt := range cs.Body.Statements {
		s, ok := stmt.(*DefStatement)

		if ok && s.Name.Value == methodName {
			ds = s
			return
		}
	}

	t.Fatalf("Can't find method '%s' in class '%s'", methodName, cs.Name)
	return
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
func (ms *ModuleStatement) IsModuleStmt(t *testing.T, moduleName string) *ModuleStatement {
	return ms
}

// HasMethod checks if current class statement has target method, and returns the method if it has
func (ms *ModuleStatement) HasMethod(t *testing.T, methodName string) (ds *DefStatement) {
	for _, stmt := range ms.Body.Statements {
		s, ok := stmt.(*DefStatement)

		if ok && s.Name.Value == methodName {
			ds = s
			return
		}
	}

	t.Fatalf("Can't find method '%s' in module '%s'", methodName, ms.Name)
	return
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
func (ds *DefStatement) IsDefStmt(t *testing.T) *DefStatement {
	return ds
}

// MethodBody returns method body's statements and assert them as TestingStatements
func (ds *DefStatement) MethodBody() CodeBlock {
	var tss []TestingStatement

	for _, stmt := range ds.BlockStatement.Statements {
		tss = append(tss, stmt.(TestingStatement))
	}

	return tss
}

// ShouldHasName checks if the method's name is what we expected
func (ds *DefStatement) ShouldHasName(t *testing.T, expectedName string) {
	if ds.Name.Value != expectedName {
		t.Fatalf("It's method %s, not %s", ds.Name.Value, expectedName)
	}
}

// ShouldHasNoParam checks if the method has no param
func (ds *DefStatement) ShouldHasNoParam(t *testing.T) {
	if len(ds.Parameters) != 0 {
		t.Fatalf("Expect method %s not to have any params, got: %d", ds.Name.Value, len(ds.Parameters))
	}
}

// ShouldHasNormalParam checks if the method has expected normal argument
func (ds *DefStatement) ShouldHasNormalParam(t *testing.T, paramName string) {
	for _, param := range ds.Parameters {
		p, ok := param.(*Identifier)

		if ok && p.NameIs(paramName) {
			return
		}
	}

	t.Fatalf("Can't find normal param '%s' in method '%s'", paramName, ds.Name.Value)
}

// ShouldHasOptionalParam checks if the method has expected optional argument
func (ds *DefStatement) ShouldHasOptionalParam(t *testing.T, paramName string) {
	for _, param := range ds.Parameters {
		p, ok := param.(*AssignExpression)

		if ok && p.NameIs(paramName) {
			return
		}
	}

	t.Fatalf("Can't find optional param '%s' in method '%s'", paramName, ds.Name.Value)
}

// ShouldHasRequiredKeywordParam checks if the method has expected keyword argument
func (ds *DefStatement) ShouldHasRequiredKeywordParam(t *testing.T, paramName string) {
	for _, param := range ds.Parameters {
		p, ok := param.(*ArgumentPairExpression)

		if ok && p.NameIs(paramName) && p.Value == nil {
			return
		}
	}

	t.Fatalf("Can't find required keyword param '%s' in method '%s'", paramName, ds.Name.Value)
}

// ShouldHasOptionalKeywordParam checks if the method has expected optional keyword argument
func (ds *DefStatement) ShouldHasOptionalKeywordParam(t *testing.T, paramName string) {
	for _, param := range ds.Parameters {
		p, ok := param.(*ArgumentPairExpression)

		if ok && p.NameIs(paramName) && p.Value != nil {
			return
		}
	}

	t.Fatalf("Can't find optional keyword param '%s' in method '%s'", paramName, ds.Name.Value)
}

// ShouldHasSplatParam checks if the method has expected splat argument
func (ds *DefStatement) ShouldHasSplatParam(t *testing.T, paramName string) {
	for _, param := range ds.Parameters {
		p, ok := param.(*PrefixExpression)

		if ok && p.NameIs(paramName) && p.Operator == "*" {
			return
		}
	}

	t.Fatalf("Can't find splat param '%s' in method '%s'", paramName, ds.Name.Value)
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
