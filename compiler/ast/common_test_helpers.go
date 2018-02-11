//+build !release

package ast

import (
	"testing"
)

/*
Program
*/

// FirstStmt returns program's first statement as a TestStatement
func (p *Program) FirstStmt() TestingStatement {
	return p.Statements[0].(TestingStatement)
}

func (p *Program) NthStmt(nth int) TestingStatement {
	return p.Statements[nth-1].(TestingStatement)
}

/*
 BaseNode
*/

// IsClassStmt fails the test and returns nil by default
func (b *BaseNode) IsClassStmt(t *testing.T, className string) (cs *ClassStatement) {
	t.Fatalf("Node is not a class statement, is %v", b)
	return
}

// IsModuleStmt fails the test and returns nil by default
func (b *BaseNode) IsModuleStmt(t *testing.T, moduleName string) (cs *ModuleStatement) {
	t.Fatalf("Node is not a module statement, is %v", b)
	return
}

// IsReturnStmt fails the test and returns nil by default
func (b *BaseNode) IsReturnStmt(t *testing.T) (rs *ReturnStatement) {
	t.Fatalf("Node is not a return statement, is %v", b)
	return
}

// IsDefStmt fails the test and returns nil by default
func (b *BaseNode) IsDefStmt(t *testing.T) (rs *DefStatement) {
	t.Fatalf("Node is not a method definition, is %v", b)
	return
}

// IsExpression fails the test and returns nil by default
func (b *BaseNode) IsExpression(t *testing.T) (te TestingExpression) {
	t.Fatalf("Node is not an expression, is %v", b)
	return
}

// NameIs returns false by default
func (b *BaseNode) NameIs(n string) bool {
	return false
}

/*
interal helpers
*/

func compareInt(t *testing.T, exp Expression, value int) {
	il, ok := exp.(*IntegerLiteral)
	if !ok {
		t.Fatalf("expect exp to be IntegerLiteral. got=%T", exp)
	}
	if il.Value != value {
		t.Fatalf("il.Value is not %d. got=%d", value, il.Value)
	}
}

func compareString(t *testing.T, exp Expression, value string) {
	sl, ok := exp.(*StringLiteral)
	if !ok {
		t.Fatalf("expect exp to be StringLiteral. got=%T", exp)
	}
	if sl.Value != value {
		t.Fatalf("il.Value is not %s. got=%s", value, sl.Value)
	}
}

func compareIdentifier(t *testing.T, exp Expression, value TestingIdentifier) {
	sl, ok := exp.(*Identifier)
	if !ok {
		t.Fatalf("expect exp to be StringLiteral. got=%T", exp)
	}
	if sl.Value != string(value) {
		t.Fatalf("il.Value is not %s. got=%s", value, sl.Value)
	}
}

func compareBool(t *testing.T, exp Expression, value bool) {
	b, ok := exp.(*BooleanExpression)
	if !ok {
		t.Fatalf("expect exp to be IntegerLiteral. got=%T", exp)
	}
	if b.Value != value {
		t.Fatalf("il.Value is not %d. got=%d", value, b.Value)
	}
}
