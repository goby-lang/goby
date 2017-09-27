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
				return func(t *thread, args []Object, blockFrame *callFrame) Object {
					m := make(map[string]interface{})

					if len(args) == 0 {
						return t.vm.initGoMap(m)
					}

					hash, ok := args[0].(*HashObject)

					if !ok {
						return t.vm.initErrorObject(errors.TypeError, errors.WrongArgumentTypeFormat, classes.HashClass, args[0].Class().Name)
					}

					for k, v := range hash.Pairs {
						m[k] = v
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
				return func(t *thread, args []Object, blockFrame *callFrame) Object {
					if blockFrame == nil {
						return t.vm.initErrorObject(errors.InternalError, errors.CantYieldWithoutBlockFormat)
					}

					m := receiver.(*GoMap)

					pairs := map[string]Object{}

					for k, v := range m.data {
						pairs[k] = t.vm.initObjectFromGoType(v)

					}

					return t.vm.initHashObject(pairs)
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
