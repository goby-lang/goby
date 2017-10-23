package vm

import (
	"fmt"

	"github.com/goby-lang/goby/vm/classes"
	"github.com/goby-lang/goby/vm/errors"
	"github.com/st0012/metago"
)

// GoObject ...
type GoObject struct {
	*baseObj
	data interface{}
}

// Class methods --------------------------------------------------------
func builtinGoObjectClassMethods() []*BuiltinMethodObject {
	return []*BuiltinMethodObject{}
}

// Instance methods -----------------------------------------------------
func builtinGoObjectInstanceMethods() []*BuiltinMethodObject {
	return []*BuiltinMethodObject{
		{
			Name: "go_func",
			Fn: func(receiver Object, instruction *instruction) builtinMethodBody {
				return func(t *thread, args []Object, blockFrame *normalCallFrame) Object {
					s, ok := args[0].(*StringObject)

					if !ok {
						return t.vm.initErrorObject(errors.TypeError, instruction, errors.WrongArgumentTypeFormat, classes.StringClass, args[0].Class().Name)
					}

					funcName := s.value
					r := receiver.(*GoObject)

					funcArgs, err := convertToGoFuncArgs(args[1:])

					if err != nil {
						t.vm.initErrorObject(errors.TypeError, instruction, err.Error())
					}

					result := metago.CallFunc(r.data, funcName, funcArgs...)
					return t.vm.initObjectFromGoType(result)
				}
			},
		},
	}
}

// Internal functions ===================================================

// Functions for initialization -----------------------------------------

func (vm *VM) initGoObject(d interface{}) *GoObject {
	return &GoObject{data: d, baseObj: &baseObj{class: vm.topLevelClass(classes.GoObjectClass)}}
}

func (vm *VM) initGoClass() *RClass {
	sc := vm.initializeClass(classes.GoObjectClass, false)
	sc.setBuiltinMethods(builtinGoObjectClassMethods(), true)
	sc.setBuiltinMethods(builtinGoObjectInstanceMethods(), false)
	vm.objectClass.setClassConstant(sc)
	return sc
}

// Polymorphic helper functions -----------------------------------------

// Value returns the object
func (s *GoObject) Value() interface{} {
	return s.data
}

// toString returns the object's name as the string format
func (s *GoObject) toString() string {
	return fmt.Sprintf("<GoObject: %p>", s)
}

// toJSON just delegates to toString
func (s *GoObject) toJSON() string {
	return s.toString()
}

// Other helper functions -----------------------------------------------

func convertToGoFuncArgs(args []Object) ([]interface{}, error) {
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
