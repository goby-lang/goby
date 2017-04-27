package vm

import (
	"testing"
)

func getBuiltInMethod(t *testing.T, receiver BaseObject, methodName string) builtinMethodBody {
	m := receiver.returnClass().LookupInstanceMethod(methodName)

	if m == nil {
		t.Fatalf("Undefined built in method %s for %s", methodName, receiver.returnClass().ReturnName())
	}

	method := m.(*BuiltInMethod)
	return method.Fn(receiver)
}
