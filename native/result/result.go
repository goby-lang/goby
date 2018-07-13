package result

//go:generate binder -in result.go -type Result

import "github.com/goby-lang/goby/vm"

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

func (r *Result) ToJSON(*Thread) string {
	return ""
}

func (r *Result) ToString() string {
	return ""
}

func (r *Result) Value() interface{} {
	return r.value
}

// New creates and returns a new isntance of a Result
func (Result) New(t *Thread, name Object, value Object) Object {
	r := &Result{
		name:    name,
		value:   value,
		BaseObj: vm.NewBaseObject(t.VM(), "Result"),
	}
	if name == vm.NULL {
		r.empty = true
	}

	return r
}

// Empty creats a new empty Result
func (Result) Empty(t *Thread) Object {
	return &Result{empty: true,
		BaseObj: vm.NewBaseObject(t.VM(), "Result"),
	}
}

// MethodMissing will be called for all methods other than 'or'
func (r *Result) MethodMissing(t *Thread, name Object) Object {

	if name.Value() == r.name.Value() && !r.used {
		r.used = true
		if t.BlockGiven() {
			t.Yield(r.value)
		}
	}

	return r
}

// Or should be the final catch all for a result call chain
func (r *Result) Or(t *Thread) Object {
	if r.used || r.empty {
		return vm.NULL
	}

	if t.BlockGiven() {
		t.Yield(r.name, r.value)
	}

	return vm.NULL
}
