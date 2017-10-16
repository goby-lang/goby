package vm

import (
	"fmt"

	"github.com/goby-lang/goby/vm/classes"
	"github.com/goby-lang/goby/vm/errors"
)

// GoMap ...
type GoMap struct {
	*baseObj
	data map[string]interface{}
}

// Class methods --------------------------------------------------------
func builtinGoMapClassMethods() []*BuiltinMethodObject {
	return []*BuiltinMethodObject{
		{
			// Initialize a new GoMap instance.
			// It can be called without any arguments, which will create an empty map.
			// Or you can pass a hash as argument, so the map will have same pairs.
			//
			// @return [GoMap]
			Name: "new",
			Fn: func(receiver Object) builtinMethodBody {
				return func(t *thread, args []Object, blockFrame *normalCallFrame) Object {
					m := make(map[string]interface{})

					if len(args) == 0 {
						return t.vm.initGoMap(m)
					}

					hash, ok := args[0].(*HashObject)

					if !ok {
						return t.vm.initErrorObject(errors.TypeError, errors.WrongArgumentTypeFormat, classes.HashClass, args[0].Class().Name)
					}

					for k, v := range hash.Pairs {
						m[k] = v.Value()
					}

					return t.vm.initGoMap(m)
				}
			},
		},
	}
}

// Instance methods -----------------------------------------------------
func builtinGoMapInstanceMethods() []*BuiltinMethodObject {
	return []*BuiltinMethodObject{
		{
			Name: "to_hash",
			Fn: func(receiver Object) builtinMethodBody {
				return func(t *thread, args []Object, blockFrame *normalCallFrame) Object {
					if len(args) != 0 {
						return t.vm.initErrorObject(errors.ArgumentError, "Expect 0 argument. got: %d", len(args))
					}

					m := receiver.(*GoMap)

					pairs := map[string]Object{}

					for k, obj := range m.data {
						pairs[k] = t.vm.initObjectFromGoType(obj)

					}

					return t.vm.initHashObject(pairs)
				}
			},
		},
		{
			Name: "get",
			Fn: func(receiver Object) builtinMethodBody {
				return func(t *thread, args []Object, blockFrame *normalCallFrame) Object {
					if len(args) != 1 {
						return t.vm.initErrorObject(errors.ArgumentError, "Expect 1 argument. got: %d", len(args))
					}

					key, ok := args[0].(*StringObject)

					if !ok {
						return t.vm.initErrorObject(errors.TypeError, errors.WrongArgumentTypeFormat, classes.StringClass, args[0].Class().Name)
					}

					m := receiver.(*GoMap).data

					result, ok := m[key.value]

					if !ok {
						return NULL
					}

					obj, ok := result.(Object)

					if !ok {
						obj = t.vm.initObjectFromGoType(result)
					}

					return obj
				}
			},
		},
		{
			Name: "set",
			Fn: func(receiver Object) builtinMethodBody {
				return func(t *thread, args []Object, blockFrame *normalCallFrame) Object {
					if len(args) != 2 {
						return t.vm.initErrorObject(errors.ArgumentError, "Expect 2 argument. got: %d", len(args))
					}

					key, ok := args[0].(*StringObject)

					if !ok {
						return t.vm.initErrorObject(errors.TypeError, errors.WrongArgumentTypeFormat, classes.StringClass, args[0].Class().Name)
					}

					m := receiver.(*GoMap).data

					m[key.value] = args[1]

					return args[1]
				}
			},
		},
	}
}

// Internal functions ===================================================

// Functions for initialization -----------------------------------------

func (vm *VM) initGoMap(d map[string]interface{}) *GoMap {
	return &GoMap{data: d, baseObj: &baseObj{class: vm.topLevelClass(classes.GoMapClass)}}
}

func (vm *VM) initGoMapClass() *RClass {
	sc := vm.initializeClass(classes.GoMapClass, false)
	sc.setBuiltinMethods(builtinGoMapClassMethods(), true)
	sc.setBuiltinMethods(builtinGoMapInstanceMethods(), false)
	vm.objectClass.setClassConstant(sc)
	return sc
}

// Polymorphic helper functions -----------------------------------------

// Value returns the object
func (m *GoMap) Value() interface{} {
	return m.data
}

// toString returns the object's name as the string format
func (m *GoMap) toString() string {
	return fmt.Sprintf("<GoMap: %p>", m)
}

// toJSON just delegates to toString
func (m *GoMap) toJSON() string {
	return m.toString()
}
