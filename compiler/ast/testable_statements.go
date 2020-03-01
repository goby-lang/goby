//+build !release

package ast

import (
	"testing"
)

// TestableStatement holds predicate functions for checking statements
type TestableStatement interface {
	Statement
	// Test Helpers
	IsClassStmt(t *testing.T) *testableClassStatement
	IsDefStmt(t *testing.T) *testableDefStatement
	IsExpression(t *testing.T) testableExpression
	IsModuleStmt(t *testing.T) *testableModuleStatement
	IsReturnStmt(t *testing.T) *testableReturnStatement
	IsWhileStmt(t *testing.T) *testableWhileStatement
}

// CodeBlock is a list of TestableStatement
type CodeBlock []TestableStatement

// NthStmt returns the trailing TestableStatement
func (cb CodeBlock) NthStmt(n int) TestableStatement {
	return cb[n-1]
}

type testableClassStatement struct {
	*ClassStatement
	t *testing.T
}

// HasMethod checks if current class statement has target method, and returns the method if it has
func (tcs *testableClassStatement) HasMethod(methodName string) (ds *testableDefStatement) {
	for _, stmt := range tcs.Body.Statements {
		s, ok := stmt.(*DefStatement)

		if ok && s.Name.Value == methodName {
			ds = &testableDefStatement{
				DefStatement: s,
				t:            tcs.t,
			}
			return
		}
	}

	tcs.t.Fatalf("Can't find method '%s' in class '%s'", methodName, tcs.Name)
	return
}

// ShouldHaveName checks if current class's name matches the specified name.
func (tcs *testableClassStatement) ShouldHaveName(name string) {
	if tcs.Name.Value != name {
		tcs.t.Helper()
		tcs.t.Fatalf("Wrong class, this class is %s", tcs.Name.Value)
	}
}

// ShouldInherit checks if current class statement inherits the target class
func (tcs *testableClassStatement) ShouldInherit(className string) {
	if tcs.SuperClassName != className {
		tcs.t.Helper()
		tcs.t.Fatalf("Expect class %s to inherit class %s. got %s", tcs.Name, className, tcs.SuperClassName)
	}
}

type testableDefStatement struct {
	*DefStatement
	t *testing.T
}

// MethodBody returns method body's statements and assert them as TestingStatements
func (tds *testableDefStatement) MethodBody() CodeBlock {
	var tss []TestableStatement

	for _, stmt := range tds.BlockStatement.Statements {
		tss = append(tss, stmt.(TestableStatement))
	}

	return tss
}

// ShouldHaveName checks if the method's name is what we expected
func (tds *testableDefStatement) ShouldHaveName(expectedName string) {
	if tds.Name.Value != expectedName {
		tds.t.Helper()
		tds.t.Fatalf("It's method %s, not %s", tds.Name.Value, expectedName)
	}
}

// ShouldHaveNoParam checks if the method has no param
func (tds *testableDefStatement) ShouldHaveNoParam() {
	if len(tds.Parameters) != 0 {
		tds.t.Helper()
		tds.t.Fatalf("Expect method %s not to have any params, got: %d", tds.Name.Value, len(tds.Parameters))
	}
}

// ShouldHaveNormalParam checks if the method has expected normal argument
func (tds *testableDefStatement) ShouldHaveNormalParam(expectedName string) {
	for _, param := range tds.Parameters {
		p, ok := param.(*Identifier)

		if ok && p.Value == expectedName {
			return
		}
	}

	tds.t.Helper()
	tds.t.Fatalf("Can't find normal param '%s' in method '%s'", expectedName, tds.Name.Value)
}

// ShouldHaveOptionalParam checks if the method has expected optional argument
func (tds *testableDefStatement) ShouldHaveOptionalParam(expectedName string) {
	for _, param := range tds.Parameters {
		p, ok := param.(*AssignExpression)

		if ok {
			paramName := p.Variables[0].(*Identifier).Value
			if paramName == expectedName {
				return
			}
		}
	}

	tds.t.Helper()
	tds.t.Fatalf("Can't find optional param '%s' in method '%s'", expectedName, tds.Name.Value)
}

// ShouldHaveRequiredKeywordParam checks if the method has expected keyword argument
func (tds *testableDefStatement) ShouldHaveRequiredKeywordParam(expectedName string) {
	for _, param := range tds.Parameters {
		p, ok := param.(*ArgumentPairExpression)

		if ok {
			paramName := p.Key.(*Identifier).Value
			if expectedName == paramName && p.Value == nil {
				return
			}

		}
	}

	tds.t.Helper()
	tds.t.Fatalf("Can't find required keyword param '%s' in method '%s'", expectedName, tds.Name.Value)
}

// ShouldHaveOptionalKeywordParam checks if the method has expected optional keyword argument
func (tds *testableDefStatement) ShouldHaveOptionalKeywordParam(expectedName string) {
	for _, param := range tds.Parameters {
		p, ok := param.(*ArgumentPairExpression)

		if ok {
			paramName := p.Key.(*Identifier).Value
			if expectedName == paramName && p.Value != nil {
				return
			}

		}
	}

	tds.t.Helper()
	tds.t.Fatalf("Can't find optional keyword param '%s' in method '%s'", expectedName, tds.Name.Value)
}

// ShouldHaveSplatParam checks if the method has expected splat argument
func (tds *testableDefStatement) ShouldHaveSplatParam(expectedName string) {
	for _, param := range tds.Parameters {
		p, ok := param.(*PrefixExpression)

		if ok {
			paramName := p.Right.(*Identifier).Value
			if expectedName == paramName {
				return
			}
		}
	}

	tds.t.Helper()
	tds.t.Fatalf("Can't find splat param '%s' in method '%s'", expectedName, tds.Name.Value)
}

type testableModuleStatement struct {
	*ModuleStatement
	t *testing.T
}

// HasMethod checks if current class statement has target method, and returns the method if it has
func (tms *testableModuleStatement) HasMethod(t *testing.T, methodName string) (ds *testableDefStatement) {
	for _, stmt := range tms.Body.Statements {
		s, ok := stmt.(*DefStatement)

		if ok && s.Name.Value == methodName {
			ds = &testableDefStatement{
				DefStatement: s,
				t:            tms.t,
			}
			return
		}
	}

	t.Helper()
	t.Fatalf("Can't find method '%s' in module '%s'", methodName, tms.Name)
	return
}

// ShouldHaveName checks if current class's name matches the specified name.
func (tms *testableModuleStatement) ShouldHaveName(name string) {
	if tms.Name.Value != name {
		tms.t.Helper()
		tms.t.Fatalf("Wrong class, this class is %s", tms.Name.Value)
	}
}

type testableReturnStatement struct {
	*ReturnStatement
	t *testing.T
}

// ShouldHaveValue checks if the current value matches the specified name.
func (trs *testableReturnStatement) ShouldHaveValue(value interface{}) {
	t := trs.t
	t.Helper()
	rs := trs.ReturnStatement
	switch v := value.(type) {
	case int:
		compareInt(t, rs.ReturnValue, v)
	case string:
		compareString(t, rs.ReturnValue, v)
	case bool:
		compareBool(t, rs.ReturnValue, v)
	case TestableIdentifierValue:
		compareIdentifier(t, rs.ReturnValue, v)
	}
}

type testableWhileStatement struct {
	*WhileStatement
	t *testing.T
}

// CodeBlock returns while statement's code block as a set of TestingStatements
func (tws *testableWhileStatement) CodeBlock() CodeBlock {
	var tss []TestableStatement

	for _, stmt := range tws.Body.Statements {
		tss = append(tss, stmt.(TestableStatement))
	}

	return tss
}

// ConditionExpression returns while statement's condition as TestingExpression
func (tws *testableWhileStatement) ConditionExpression() testableExpression {
	return tws.Condition.(testableExpression)
}
