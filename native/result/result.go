package result

//go:generate binder -in result.go -type Result

import (
	"github.com/goby-lang/goby/vm"
)

type Object = vm.Object
type Thread = vm.Thread

// Result is a variant return type
type Result struct {
	*vm.BaseObj
	empty bool
	used  bool
	name  Object
	value Object
}

func (r *Result) Value() interface{} {
	return r.value
}

func (r *Result) ToJSON(*vm.Thread) string {
	return ""
}

func (r *Result) ToString() string {
	return ""
}

// New creates and returns a new isntance of a Result
func (Result) New(t *Thread, name, value Object) Object {
	r := &Result{
		BaseObj: vm.NewBaseObj(t.VM(), "Result"),
		name:    name,
		value:   value,
	}
	if name == vm.NULL {
		r.empty = true
	}

	return r
}

// MethodMissing will be called for all methods other than 'or'
func (r *Result) MethodMissing(t *Thread, name Object, args Object) Object {
	if name == r.name && !r.used {
		r.used = true
		if t.BlockGiven() {
			t.Yield(r.value)
		}
	}

	return r
}

func (r *Result) Or(t *Thread) Object {
	if r.used || r.empty {
		return vm.NULL
	}

	if t.BlockGiven() {
		t.Yield(r.name, r.value)
	}

	return vm.NULL
}
