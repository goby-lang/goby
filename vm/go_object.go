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
func builtinGoClassMethods() []*BuiltinMethodObject {
	return []*BuiltinMethodObject{}
}

// Instance methods -----------------------------------------------------
func builtinGoInstanceMethods() []*BuiltinMethodObject {
	return []*BuiltinMethodObject{
		{
			Name: "go_func",
			Fn: func(receiver Object) builtinMethodBody {
				return func(t *thread, args []Object, blockFrame *callFrame) Object {
					s, ok := args[0].(*StringObject)

					if !ok {
						return t.vm.initErrorObject(errors.TypeError, errors.WrongArgumentTypeFormat, classes.StringClass, args[0].Class().Name)
					}

					funcName := s.value
					r := receiver.(*GoObject)

					funcArgs, err := convertToGoFuncArgs(args[1:])

					if err != nil {
						t.vm.initErrorObject(errors.TypeError, err.Error())
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
	sc.setBuiltinMethods(builtinGoClassMethods(), true)
	sc.setBuiltinMethods(builtinGoInstanceMethods(), false)
	vm.objectClass.setClassConstant(sc)
	return sc
}

// Polymorphic helper functions -----------------------------------------

func (s *GoObject) Value() interface{} {
	return s.data
}
func (s *GoObject) toString() string {
	return fmt.Sprintf("<GoObject: %p>", s)
}

func (s *GoObject) toJSON() string {
	return s.toString()
}

// Other helper functions -----------------------------------------------

func convertToGoFuncArgs(args []Object) ([]interface{}, error) {
	funcArgs := []interface{}{}

	for _, arg := range args {
		v, ok := arg.(builtinType)

		if ok {
			switch v := v.(type) {
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

			funcArgs = append(funcArgs, v.Value())
		} else {
			err := fmt.Errorf("Can't pass %s type object when calling go function", arg.Class().Name)
			return nil, err
		}
	}

	return funcArgs, nil
}
