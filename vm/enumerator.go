package vm

import (
	"bytes"
	"strings"
)

func (vm *VM) initEnumeratorObject(elements []Object) *EnumeratorObject {
	return &EnumeratorObject{
		baseObj:  &baseObj{class: vm.topLevelClass(enumeratorClass)},
		Elements: elements,
		Index: 0,
	}
}

func (vm *VM) initEnumeratorClass() *RClass {
	ac := vm.initializeClass(enumeratorClass, false)
	ac.setBuiltInMethods(builtinEnumeratorInstanceMethods(), false)
	ac.setBuiltInMethods(builtInEnumeratorClassMethods(), true)
	return ac
}

// ArrayObject represents instance from Array class.
// An array is a collection of different objects that are ordered and indexed.
// Elements in an array can belong to any class.
type EnumeratorObject struct {
	*baseObj
	Elements []Object
	Index    int
}

func (a *EnumeratorObject) Value() interface{} {
	return a.Elements
}

// Polymorphic helper functions -----------------------------------------
func (a *EnumeratorObject) toString() string {
	return "Enumerator String Object"
}

func (a *EnumeratorObject) toJSON() string {
	return "Enumerator JSON Object"
}

func builtInEnumeratorClassMethods() []*BuiltInMethodObject {
	return []*BuiltInMethodObject{
		{
			Name: "new",
			Fn: func(receiver Object) builtinMethodBody {
				return func(t *thread, args []Object, blockFrame *callFrame) Object {
					return t.unsupportedMethodError("#new", receiver)
				}
			},
		},
	}
}

func builtinEnumeratorInstanceMethods() []*BuiltInMethodObject {
	return []*BuiltInMethodObject{}
}
