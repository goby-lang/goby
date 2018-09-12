//+build !release

package ast

import (
	"testing"
)

type TestableStatement interface {
	Statement
	// Test Helpers
	IsClassStmt(t *testing.T) *TestableClassStatement
	IsDefStmt(t *testing.T) *TestableDefStatement
	IsExpression(t *testing.T) TestableExpression
	IsModuleStmt(t *testing.T) *TestableModuleStatement
	IsReturnStmt(t *testing.T) *TestableReturnStatement
	IsWhileStmt(t *testing.T) *TestableWhileStatement
}

type CodeBlock []TestableStatement

func (cb CodeBlock) NthStmt(n int) TestableStatement {
	return cb[n-1]
}

/*TestableClassStatement*/

type TestableClassStatement struct {
	*ClassStatement
	t *testing.T
}

// HasMethod checks if current class statement has target method, and returns the method if it has
func (tcs *TestableClassStatement) HasMethod(methodName string) (ds *TestableDefStatement) {
	for _, stmt := range tcs.Body.Statements {
		s, ok := stmt.(*DefStatement)

		if ok && s.Name.Value == methodName {
			ds = &TestableDefStatement{
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
func (tcs *TestableClassStatement) ShouldHaveName(name string) {
	if tcs.Name.Value != name {
		tcs.t.Helper()
		tcs.t.Fatalf("Wrong class, this class is %s", tcs.Name.Value)
	}
}

// ShouldInherit checks if current class statement inherits the target class
func (tcs *TestableClassStatement) ShouldInherit(className string) {
	if tcs.SuperClassName != className {
		tcs.t.Helper()
		tcs.t.Fatalf("Expect class %s to inherit class %s. got %s", tcs.Name, className, tcs.SuperClassName)
	}
}

/*TestableDefStatement*/

type TestableDefStatement struct {
	*DefStatement
	t *testing.T
}

// MethodBody returns method body's statements and assert them as TestingStatements
func (tds *TestableDefStatement) MethodBody() CodeBlock {
	var tss []TestableStatement

	for _, stmt := range tds.BlockStatement.Statements {
		tss = append(tss, stmt.(TestableStatement))
	}

	return tss
}

// ShouldHaveName checks if the method's name is what we expected
func (tds *TestableDefStatement) ShouldHaveName(expectedName string) {
	if tds.Name.Value != expectedName {
		tds.t.Helper()
		tds.t.Fatalf("It's method %s, not %s", tds.Name.Value, expectedName)
	}
}

// ShouldHaveNoParam checks if the method has no param
func (tds *TestableDefStatement) ShouldHaveNoParam() {
	if len(tds.Parameters) != 0 {
		tds.t.Helper()
		tds.t.Fatalf("Expect method %s not to have any params, got: %d", tds.Name.Value, len(tds.Parameters))
	}
}

// ShouldHaveNormalParam checks if the method has expected normal argument
func (tds *TestableDefStatement) ShouldHaveNormalParam(expectedName string) {
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
func (tds *TestableDefStatement) ShouldHaveOptionalParam(expectedName string) {
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
func (tds *TestableDefStatement) ShouldHaveRequiredKeywordParam(expectedName string) {
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
func (tds *TestableDefStatement) ShouldHaveOptionalKeywordParam(expectedName string) {
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
func (tds *TestableDefStatement) ShouldHaveSplatParam(expectedName string) {
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

/*TestableModuleStatement*/

type TestableModuleStatement struct {
	*ModuleStatement
	t *testing.T
}

// HasMethod checks if current class statement has target method, and returns the method if it has
func (tms *TestableModuleStatement) HasMethod(t *testing.T, methodName string) (ds *TestableDefStatement) {
	for _, stmt := range tms.Body.Statements {
		s, ok := stmt.(*DefStatement)

		if ok && s.Name.Value == methodName {
			ds = &TestableDefStatement{
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
func (tms *TestableModuleStatement) ShouldHaveName(name string) {
	if tms.Name.Value != name {
		tms.t.Helper()
		tms.t.Fatalf("Wrong class, this class is %s", tms.Name.Value)
	}
}

/*TestableReturnStatement*/

type TestableReturnStatement struct {
	*ReturnStatement
	t *testing.T
}

// ShouldHaveValue checks if the current value matches the specified name.
func (trs *TestableReturnStatement) ShouldHaveValue(value interface{}) {
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

/*TestableWhileStatement*/

type TestableWhileStatement struct {
	*WhileStatement
	t *testing.T
}

// Block returns while statement's code block as a set of TestingStatements
func (tws *TestableWhileStatement) CodeBlock() CodeBlock {
	var tss []TestableStatement

	for _, stmt := range tws.Body.Statements {
		tss = append(tss, stmt.(TestableStatement))
	}

	return tss
}

// ConditionExpression returns while statement's condition as TestingExpression
func (tws *TestableWhileStatement) ConditionExpression() TestableExpression {
	return tws.Condition.(TestableExpression)
}
