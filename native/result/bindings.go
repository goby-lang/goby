package result

import (
	"fmt"
	vm "github.com/goby-lang/goby/vm"
	errors "github.com/goby-lang/goby/vm/errors"
)

func init() {
	vm.RegisterExternalClass(
		"result", vm.ExternalClass(
			"Result",
			"result.gb",
			map[string]vm.Method{
				"empty": bindingResultEmpty,
				"new":   bindingResultNew,
			},
			map[string]vm.Method{
				"method_missing": bindingResultMethodMissing,
				"or":             bindingResultOr,
			}))
}

var staticResult = new(Result)

func bindingResultNew(receiver vm.Object, line int, t *vm.Thread, args []vm.Object) vm.Object {
	r := staticResult
	if len(args) != 2 {
		return t.VM().InitErrorObject(errors.ArgumentError, line, errors.WrongNumberOfArgumentFormat, 2, len(args))
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

func bindingResultEmpty(receiver vm.Object, line int, t *vm.Thread, args []vm.Object) vm.Object {
	r := staticResult
	if len(args) != 0 {
		return t.VM().InitErrorObject(errors.ArgumentError, line, errors.WrongNumberOfArgumentFormat, 0, len(args))
	}
	return r.Empty(t)
}

func bindingResultMethodMissing(receiver vm.Object, line int, t *vm.Thread, args []vm.Object) vm.Object {
	r, ok := receiver.(*Result)
	if !ok {
		panic(fmt.Sprintf("Impossible receiver type. Wanted Result got %s", receiver))
	}
	if len(args) != 1 {
		return t.VM().InitErrorObject(errors.ArgumentError, line, errors.WrongNumberOfArgumentFormat, 1, len(args))
	}
	arg0, ok := args[0].(Object)
	if !ok {
		panic("Argument 0 must be Object")
	}

	return r.MethodMissing(t, arg0)
}

func bindingResultOr(receiver vm.Object, line int, t *vm.Thread, args []vm.Object) vm.Object {
	r, ok := receiver.(*Result)
	if !ok {
		panic(fmt.Sprintf("Impossible receiver type. Wanted Result got %s", receiver))
	}
	if len(args) != 0 {
		return t.VM().InitErrorObject(errors.ArgumentError, line, errors.WrongNumberOfArgumentFormat, 0, len(args))
	}
	return r.Or(t)
}
