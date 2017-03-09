package evaluator

import (
	"testing"
)

func getBuiltInMethod(t *testing.T, receiver BaseObject, method_name string) BuiltinMethodBody {
	m := receiver.ReturnClass().LookupInstanceMethod(method_name)

	if m == nil {
		t.Fatalf("Undefined built in method %s for %s", method_name, receiver.ReturnClass().ReturnName())
	}

	method := m.(*BuiltInMethod)
	return method.Fn(receiver)
}
