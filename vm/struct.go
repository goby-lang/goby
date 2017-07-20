package vm

import (
	"fmt"
	"reflect"
)

// StructObject ...
type StructObject struct {
	*baseObj
	data interface{}
}

func (vm *VM) initStructObject(d interface{}) *StructObject {
	return &StructObject{data: d, baseObj: &baseObj{class: vm.topLevelClass(structClass)}}
}

func (vm *VM) initStructClass() *RClass {
	sc := vm.initializeClass(structClass, false)
	sc.setBuiltInMethods(builtinStructClassMethods(), true)
	sc.setBuiltInMethods(builtinStructInstanceMethods(), false)
	vm.objectClass.setClassConstant(sc)
	return sc
}

// Only initialize file related methods after it's being required.
func builtinStructClassMethods() []*BuiltInMethodObject {
	return []*BuiltInMethodObject{}
}

// Only initialize file related methods after it's being required.
func builtinStructInstanceMethods() []*BuiltInMethodObject {
	return []*BuiltInMethodObject{
		{
			Name: "send",
			Fn: func(receiver Object) builtinMethodBody {
				return func(t *thread, args []Object, blockFrame *callFrame) Object {
					s, ok := args[0].(*StringObject)

					if !ok {
						return t.vm.initErrorObject(TypeError, WrongArgumentTypeFormat, stringClass, args[0].Class().Name)
					}

					funcName := s.Value
					r := receiver.(*StructObject)

					funcArgs := make([]reflect.Value, len(args)-1)
					for i := range args[1:] {
						funcArgs[i] = reflect.ValueOf(args[i])
					}

					result := reflect.ValueOf(reflect.ValueOf(r.data).MethodByName(funcName).Call(funcArgs))

					fmt.Println(result)

					return NULL
				}
			},
		},
	}
}

// Polymorphic helper functions -----------------------------------------

func (s *StructObject) toString() string {
	return fmt.Sprintf("<Strcut: %p>", s)
}

func (s *StructObject) toJSON() string {
	return s.toString()
}
