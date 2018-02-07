//+build !release

package ast

import (
	"testing"
)

func (p *Program) FirstStmt() TestStatement {
	return p.Statements[0].(TestStatement)
}

// BaseNode

func (b *BaseNode) IsClassStmt(t *testing.T, className string) (cs *ClassStatement) {
	t.Fatalf("Node is not a class statement, is %v", b)
	return
}

func (b *BaseNode) IsModuleStmt(t *testing.T, moduleName string) (cs *ModuleStatement) {
	t.Fatalf("Node is not a module statement, is %v", b)
	return
}

func (b *BaseNode) NameIs(n string) bool {
	return false
}

// ClassStatement

func (cs *ClassStatement) IsClassStmt(t *testing.T, className string) *ClassStatement {
	return cs
}

func (cs *ClassStatement) ShouldInherits(t *testing.T, className string) {
	if cs.SuperClassName != className {
		t.Fatalf("Expect class %s to inherit class %s. got %s", cs.Name, className, cs.SuperClassName)
	}
}

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

func (cs *ClassStatement) NameIs(n string) bool {
	if cs.Name.Value == n {
		return true
	}

	return false
}

// Module Statement

func (ms *ModuleStatement) IsModuleStmt(t *testing.T, moduleName string) *ModuleStatement {
	return ms
}

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

func (ms *ModuleStatement) NameIs(n string) bool {
	if ms.Name.Value == n {
		return true
	}

	return false
}

// DefStatement

func (ds *DefStatement) ShouldHasName(t *testing.T, expectedName string) {
	if ds.Name.Value != expectedName {
		t.Fatalf("It's method %s, not %s", ds.Name.Value, expectedName)
	}
}

func (ds *DefStatement) ShouldHasNoParam(t *testing.T) {
	if len(ds.Parameters) != 0 {
		t.Fatalf("Expect method %s not to have any params, got: %d", ds.Name.Value, len(ds.Parameters))
	}
}

func (ds *DefStatement) HasNormalParam(t *testing.T, paramName string) {
	for _, param := range ds.Parameters {
		p, ok := param.(*Identifier)

		if ok && p.NameIs(paramName) {
			return
		}
	}

	t.Fatalf("Can't find normal param '%s' in method '%s'", paramName, ds.Name.Value)
}

func (ds *DefStatement) HasOptionalParam(t *testing.T, paramName string) {
	for _, param := range ds.Parameters {
		p, ok := param.(*AssignExpression)

		if ok && p.NameIs(paramName) {
			return
		}
	}

	t.Fatalf("Can't find optional param '%s' in method '%s'", paramName, ds.Name.Value)
}

func (ds *DefStatement) HasRequiredKeywordParam(t *testing.T, paramName string) {
	for _, param := range ds.Parameters {
		p, ok := param.(*ArgumentPairExpression)

		if ok && p.NameIs(paramName) && p.Value == nil {
			return
		}
	}

	t.Fatalf("Can't find required keyword param '%s' in method '%s'", paramName, ds.Name.Value)
}

func (ds *DefStatement) HasOptionalKeywordParam(t *testing.T, paramName string) {
	for _, param := range ds.Parameters {
		p, ok := param.(*ArgumentPairExpression)

		if ok && p.NameIs(paramName) && p.Value != nil {
			return
		}
	}

	t.Fatalf("Can't find optional keyword param '%s' in method '%s'", paramName, ds.Name.Value)
}

func (ds *DefStatement) HasSplatParam(t *testing.T, paramName string) {
	for _, param := range ds.Parameters {
		p, ok := param.(*PrefixExpression)

		if ok && p.NameIs(paramName) && p.Operator == "*" {
			return
		}
	}

	t.Fatalf("Can't find splat param '%s' in method '%s'", paramName, ds.Name.Value)
}

// AssignExpression

func (ae *AssignExpression) NameIs(n string) bool {
	return ae.Variables[0].(testNode).NameIs(n)
}

// Identifier

func (i *Identifier) NameIs(n string) bool {
	if i.Value == n {
		return true
	}

	return false
}
