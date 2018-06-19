package result

import vm "github.com/goby-lang/goby/vm"

var _staticResult = new(Result)

func _binding_Result_New(receiver vm.Object, line int) vm.Method {
	return func(t *vm.Thread, args []vm.Object) vm.Object {
		r := _staticResult
		if len(args) != 2 {
			panic("NOT OK")
		}
		arg0, ok := args[0].(Object)
		if !ok {
			panic("NOT OK")
		}

		arg1, ok := args[1].(Object)
		if !ok {
			panic("NOT OK")
		}

		return r.New(t, arg0, arg1)
	}
}

func _binding_Result_MethodMissing(receiver vm.Object, line int) vm.Method {
	return func(t *vm.Thread, args []vm.Object) vm.Object {
		r, ok := receiver.(*Result)
		if !ok {
			panic("NOT OK")
		}
		if len(args) != 2 {
			panic("NOT OK")
		}
		arg0, ok := args[0].(Object)
		if !ok {
			panic("NOT OK")
		}

		arg1, ok := args[1].(Object)
		if !ok {
			panic("NOT OK")
		}

		return r.MethodMissing(t, arg0, arg1)
	}
}

func _binding_Result_Or(receiver vm.Object, line int) vm.Method {
	return func(t *vm.Thread, args []vm.Object) vm.Object {
		r, ok := receiver.(*Result)
		if !ok {
			panic("NOT OK")
		}
		if len(args) != 0 {
			panic("NOT OK")
		}
		return r.Or(t)
	}
}
