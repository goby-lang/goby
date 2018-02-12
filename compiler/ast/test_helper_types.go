//+build !release

package ast

const nodeFailureMsgFormat = "Node is not %s, is %T"

type TestingIdentifier string

type testingNode interface {
	// Belows are test helpers
	NameIs(name string) bool
}
