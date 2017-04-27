package vm

import (
	"testing"
)

func getBuiltInMethod(t *testing.T, receiver BaseObject, method_name string) builtinMethodBody {
	m := receiver.returnClass().LookupInstanceMethod(method_name)

	if m == nil {
		t.Fatalf("Undefined built in method %s for %s", method_name, receiver.returnClass().ReturnName())
	}

	method := m.(*BuiltInMethod)
	return method.Fn(receiver)
}
