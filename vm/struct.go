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

					for i, arg := range args[1:] {
						v, ok := arg.(builtInType)

						if ok {
							funcArgs[i] = reflect.ValueOf(v.value())
						} else {
							return t.vm.initErrorObject(InternalError, "Can't pass %s type object when calling go function", arg.Class().Name)
						}
					}

					return t.vm.initObjectFromGoType(callMethod(r.data, funcName, funcArgs))
				}
			},
		},
	}
}

func callMethod(i interface{}, methodName string, args []reflect.Value) interface{} {
	var ptr reflect.Value
	var value reflect.Value
	var finalMethod reflect.Value

	fmt.Println(i)
	value = reflect.ValueOf(i)

	// if we start with a pointer, we need to get value pointed to
	// if we start with a value, we need to get a pointer to that value
	if value.Type().Kind() == reflect.Ptr {
		fmt.Println("Ptr")
		ptr = value
		value = ptr.Elem()
	} else {
		fmt.Println("Obj")
		ptr = reflect.New(reflect.TypeOf(i))
		temp := ptr.Elem()
		temp.Set(value)
	}

	// check for method on value
	method := value.MethodByName(methodName)
	if method.IsValid() {
		finalMethod = method
	}

	// check for method on pointer
	method = ptr.MethodByName(methodName)

	if method.IsValid() {
		finalMethod = method
	}

	if finalMethod.IsValid() {
		return finalMethod.Call(args)[0].Interface()
	}

	// return or panic, method not found of either type
	panic("!!!!!!!!!!!!!!!!!!")
}

// Polymorphic helper functions -----------------------------------------

func (s *StructObject) toString() string {
	return fmt.Sprintf("<Strcut: %p>", s)
}

func (s *StructObject) toJSON() string {
	return s.toString()
}
