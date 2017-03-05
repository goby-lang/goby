package initializer_test

import (
	"github.com/st0012/Rooby/initializer"
	"github.com/st0012/Rooby/object"
	"testing"
)

func init() {
	initializer.InitializeProgram()
}

func getBuiltInMethod(t *testing.T, receiver object.BaseObject, method_name string) object.BuiltinMethodBody {
	m := receiver.ReturnClass().LookupInstanceMethod(method_name)

	if m == nil {
		t.Fatalf("Undefined built in method %s for %s", method_name, receiver.ReturnClass().ReturnName())
	}

	method := m.(*object.BuiltInMethod)
	return method.Fn(receiver)
}
