//+build !release

package ast

import "testing"

type TestingIdentifier string

type testingNode interface {
	// Belows are test helpers
	NameIs(name string) bool
}

type TestingStatement interface {
	Statement
	// Test Helpers
	IsClassStmt(t *testing.T, className string) *ClassStatement
	IsExpressionStmt(t *testing.T) *ExpressionStatement
	IsModuleStmt(t *testing.T, className string) *ModuleStatement
	IsReturnStmt(t *testing.T) *ReturnStatement
	IsDefStmt(t *testing.T) *DefStatement
}

type MethodBody []TestingStatement

func (mb MethodBody) NthStmt(n int) TestingStatement {
	return mb[n-1]
}