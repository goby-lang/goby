package vm

import (
	"fmt"

	"github.com/goby-lang/goby/vm/classes"
	"github.com/goby-lang/goby/vm/errors"
	"github.com/st0012/metago"
)

// GoObject ...
type GoObject struct {
	*BaseObj
	data interface{}
}

// Class methods --------------------------------------------------------
var builtinGoObjectClassMethods = []*BuiltinMethodObject{}

// Instance methods -----------------------------------------------------
var builtinGoObjectInstanceMethods = []*BuiltinMethodObject{
	{
		// An experimental method for loading plugins (written in Golang) dynamically.
		// Needs improvements.
		//
		/// ```ruby
		/// require "plugin"
		//
		//	p = Plugin.use "../test_fixtures/import_test/plugin/plugin.go"
		//	p.go_func("Foo", "!")
		//	p.go_func("Baz")
		// ```
		//
		// @param name [String]
		// @return [Object]
		Name: "go_func",
		Fn: func(receiver Object, sourceLine int, t *Thread, args []Object, blockFrame *normalCallFrame) Object {
			s, ok := args[0].(*StringObject)

			if !ok {
				return t.vm.InitErrorObject(errors.TypeError, sourceLine, errors.WrongArgumentTypeFormat, classes.StringClass, args[0].Class().Name)
			}

			funcName := s.value
			r := receiver.(*GoObject)

			funcArgs, err := ConvertToGoFuncArgs(args[1:])

			if err != nil {
				t.vm.InitErrorObject(errors.TypeError, sourceLine, err.Error())
			}

			result := metago.CallFunc(r.data, funcName, funcArgs...)
			return t.vm.InitObjectFromGoType(result)

		},
	},
}

// Internal functions ===================================================

// Functions for initialization -----------------------------------------

func (vm *VM) initGoObject(d interface{}) *GoObject {
	return &GoObject{data: d, BaseObj: NewBaseObject(vm.TopLevelClass(classes.GoObjectClass))}
}

func (vm *VM) initGoClass() *RClass {
	sc := vm.initializeClass(classes.GoObjectClass)
	sc.setBuiltinMethods(builtinGoObjectClassMethods, true)
	sc.setBuiltinMethods(builtinGoObjectInstanceMethods, false)
	vm.objectClass.setClassConstant(sc)
	return sc
}

// Polymorphic helper functions -----------------------------------------

// Value returns the object
func (s *GoObject) Value() interface{} {
	return s.data
}

// ToString returns the object's name as the string format
func (s *GoObject) ToString() string {
	return fmt.Sprintf("<GoObject: %p>", s)
}

// Inspect delegates to ToString
func (s *GoObject) Inspect() string {
	return s.ToString()
}

// ToJSON just delegates to ToString
func (s *GoObject) ToJSON(t *Thread) string {
	return s.ToString()
}

// Other helper functions -----------------------------------------------

// ConvertToGoFuncArgs converts Goby's args to Go func's args
func ConvertToGoFuncArgs(args []Object) ([]interface{}, error) {
	funcArgs := []interface{}{}

	for _, arg := range args {
		switch v := arg.(type) {
		case *IntegerObject:
			switch v.flag {
			case f64:
				funcArgs = append(funcArgs, float64(v.value))
				continue
			case f32:
				funcArgs = append(funcArgs, float32(v.value))
				continue
			case ui64:
				funcArgs = append(funcArgs, uint64(v.value))
				continue
			case ui32:
				funcArgs = append(funcArgs, uint32(v.value))
				continue
			case ui16:
				funcArgs = append(funcArgs, uint16(v.value))
				continue
			case ui8:
				funcArgs = append(funcArgs, uint8(v.value))
				continue
			case i64:
				funcArgs = append(funcArgs, int64(v.value))
				continue
			case i32:
				funcArgs = append(funcArgs, int32(v.value))
				continue
			case i16:
				funcArgs = append(funcArgs, int16(v.value))
				continue
			case i8:
				funcArgs = append(funcArgs, int8(v.value))
				continue
			}
		}

		funcArgs = append(funcArgs, arg.Value())
	}

	return funcArgs, nil
}
