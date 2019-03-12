package vm

import (
	"fmt"

	"github.com/goby-lang/goby/vm/classes"
	"github.com/goby-lang/goby/vm/errors"
)

// GoMap ...
type GoMap struct {
	*BaseObj
	data map[string]interface{}
}

// Class methods --------------------------------------------------------
var builtinGoMapClassMethods = []*BuiltinMethodObject{
	{
		// Initialize a new GoMap instance.
		// It can be called without any arguments, which will create an empty map.
		// Or you can pass a hash as argument, so the map will have same pairs.
		//
		// @return [GoMap]
		Name: "new",
		Fn: func(receiver Object, sourceLine int, t *Thread, args []Object, blockFrame *normalCallFrame) Object {
			m := make(map[string]interface{})

			if len(args) == 0 {
				return t.vm.initGoMap(m)
			}

			hash, ok := args[0].(*HashObject)

			if !ok {
				return t.vm.InitErrorObject(errors.TypeError, sourceLine, errors.WrongArgumentTypeFormat, classes.HashClass, args[0].Class().Name)
			}

			for k, v := range hash.Pairs {
				m[k] = v.Value()
			}

			return t.vm.initGoMap(m)

		},
	},
}

// Instance methods -----------------------------------------------------
var builtinGoMapInstanceMethods = []*BuiltinMethodObject{
	{
		Name: "get",
		Fn: func(receiver Object, sourceLine int, t *Thread, args []Object, blockFrame *normalCallFrame) Object {
			if len(args) != 1 {
				return t.vm.InitErrorObject(errors.ArgumentError, sourceLine, errors.WrongNumberOfArgument, 1, len(args))
			}

			key, ok := args[0].(*StringObject)

			if !ok {
				return t.vm.InitErrorObject(errors.TypeError, sourceLine, errors.WrongArgumentTypeFormat, classes.StringClass, args[0].Class().Name)
			}

			m := receiver.(*GoMap).data

			result, ok := m[key.value]

			if !ok {
				return NULL
			}

			obj, ok := result.(Object)

			if !ok {
				obj = t.vm.InitObjectFromGoType(result)
			}

			return obj

		},
	},
	{
		Name: "set",
		Fn: func(receiver Object, sourceLine int, t *Thread, args []Object, blockFrame *normalCallFrame) Object {
			if len(args) != 2 {
				return t.vm.InitErrorObject(errors.ArgumentError, sourceLine, errors.WrongNumberOfArgument, 2, len(args))
			}

			key, ok := args[0].(*StringObject)

			if !ok {
				return t.vm.InitErrorObject(errors.TypeError, sourceLine, errors.WrongArgumentTypeFormatNum, 1, classes.StringClass, args[0].Class().Name)
			}

			m := receiver.(*GoMap).data

			m[key.value] = args[1]

			return args[1]

		},
	},
	{
		Name: "to_hash",
		Fn: func(receiver Object, sourceLine int, t *Thread, args []Object, blockFrame *normalCallFrame) Object {
			if len(args) != 0 {
				return t.vm.InitErrorObject(errors.ArgumentError, sourceLine, errors.WrongNumberOfArgument, 0, len(args))
			}

			m := receiver.(*GoMap)

			pairs := map[string]Object{}

			for k, obj := range m.data {
				pairs[k] = t.vm.InitObjectFromGoType(obj)

			}

			return t.vm.InitHashObject(pairs)

		},
	},
}

// Internal functions ===================================================

// Functions for initialization -----------------------------------------

func (vm *VM) initGoMap(d map[string]interface{}) *GoMap {
	return &GoMap{data: d, BaseObj: &BaseObj{class: vm.TopLevelClass(classes.GoMapClass)}}
}

func (vm *VM) initGoMapClass() *RClass {
	sc := vm.initializeClass(classes.GoMapClass)
	sc.setBuiltinMethods(builtinGoMapClassMethods, true)
	sc.setBuiltinMethods(builtinGoMapInstanceMethods, false)
	vm.objectClass.setClassConstant(sc)
	return sc
}

// Polymorphic helper functions -----------------------------------------

// Value returns the object
func (m *GoMap) Value() interface{} {
	return m.data
}

// ToString returns the object's name as the string format
func (m *GoMap) ToString() string {
	return fmt.Sprintf("<GoMap: %p>", m)
}

// Inspect delegates to ToString
func (m *GoMap) Inspect() string {
	return m.ToString()
}

// ToJSON just delegates to ToString
func (m *GoMap) ToJSON(t *Thread) string {
	return m.ToString()
}
