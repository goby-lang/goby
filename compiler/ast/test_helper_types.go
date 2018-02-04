//+build !release

package ast

import "testing"

type testNode interface {
	// Belows are test helpers
	NameIs(name string) bool
}

type TestStatement interface {
	Statement
	// Test Helpers
	IsClassStmt(t *testing.T, className string) *ClassStatement
	IsModuleStmt(t *testing.T, className string) *ModuleStatement
}
