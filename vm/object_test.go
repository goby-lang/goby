package vm

import (
	"testing"
)

func getBuiltInMethod(t *testing.T, receiver Object, methodName string) builtinMethodBody {
	m := receiver.returnClass().lookupInstanceMethod(methodName)

	if m == nil {
		t.Fatalf("Undefined built in method %s for %s", methodName, receiver.returnClass().ReturnName())
	}

	method := m.(*BuiltInMethodObject)
	return method.Fn(receiver)
}
