package result

import (
	"fmt"
	vm "github.com/goby-lang/goby/vm"
)

func init() {
	vm.RegisterExternalClass(
		"result", vm.ExternalClass(
			"Result",
			"result.gb",
			map[string]vm.Method{"new": _binding_Result_New},
			map[string]vm.Method{
				"method_missing": _binding_Result_MethodMissing,
				"or":             _binding_Result_Or,
			}))
}

var _staticResult = new(Result)

func _binding_Result_New(receiver vm.Object, line int, t *vm.Thread, args []vm.Object) vm.Object {
	r := _staticResult
	if len(args) != 2 {
		panic(fmt.Sprintf("Wrong NArgs. Wanted: 2 got: %d", len(args)))

	}
	arg0, ok := args[0].(Object)
	if !ok {
		panic("Argument 0 must be Object")
	}

	arg1, ok := args[1].(Object)
	if !ok {
		panic("Argument 1 must be Object")
	}

	return r.New(t, arg0, arg1)
}

func _binding_Result_MethodMissing(receiver vm.Object, line int, t *vm.Thread, args []vm.Object) vm.Object {
	r, ok := receiver.(*Result)
	if !ok {
		panic(fmt.Sprintf("Impossible receiver type. Wanted Result got %s", receiver))
	}
	if len(args) != 2 {
		panic(fmt.Sprintf("Wrong NArgs. Wanted: 2 got: %d", len(args)))

	}
	arg0, ok := args[0].(Object)
	if !ok {
		panic("Argument 0 must be Object")
	}

	arg1, ok := args[1].(Object)
	if !ok {
		panic("Argument 1 must be Object")
	}

	return r.MethodMissing(t, arg0, arg1)
}

func _binding_Result_Or(receiver vm.Object, line int, t *vm.Thread, args []vm.Object) vm.Object {
	r, ok := receiver.(*Result)
	if !ok {
		panic(fmt.Sprintf("Impossible receiver type. Wanted Result got %s", receiver))
	}
	if len(args) != 0 {
		panic(fmt.Sprintf("Wrong NArgs. Wanted: 0 got: %d", len(args)))

	}
	return r.Or(t)
}
