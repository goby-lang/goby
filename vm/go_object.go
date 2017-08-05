package vm

import (
	"fmt"
	"github.com/st0012/metago"
)

func (vm *VM) initGoObject(d interface{}) *GoObject {
	return &GoObject{data: d, baseObj: &baseObj{class: vm.topLevelClass(goObjectClass)}}
}

func (vm *VM) initGoClass() *RClass {
	sc := vm.initializeClass(goObjectClass, false)
	sc.setBuiltInMethods(builtinGoClassMethods(), true)
	sc.setBuiltInMethods(builtinGoInstanceMethods(), false)
	vm.objectClass.setClassConstant(sc)
	return sc
}

// GoObject ...
type GoObject struct {
	*baseObj
	data interface{}
}

// Polymorphic helper functions -----------------------------------------
func (s *GoObject) toString() string {
	return fmt.Sprintf("<GoObject: %p>", s)
}

func (s *GoObject) toJSON() string {
	return s.toString()
}

func (s *GoObject) Value() interface{} {
	return s.data
}

func builtinGoClassMethods() []*BuiltInMethodObject {
	return []*BuiltInMethodObject{}
}

func builtinGoInstanceMethods() []*BuiltInMethodObject {
	return []*BuiltInMethodObject{
		{
			Name: "send",
			Fn: func(receiver Object) builtinMethodBody {
				return func(t *thread, args []Object, blockFrame *callFrame) Object {
					s, ok := args[0].(*StringObject)

					if !ok {
						return t.vm.initErrorObject(TypeError, WrongArgumentTypeFormat, stringClass, args[0].Class().Name)
					}

					funcName := s.value
					r := receiver.(*GoObject)

					funcArgs, err := convertToGoFuncArgs(args[1:])

					if err != nil {
						t.vm.initErrorObject(TypeError, err.Error())
					}

					result := metago.CallFunc(r.data, funcName, funcArgs...)
					return t.vm.initObjectFromGoType(result)
				}
			},
		},
	}
}

func convertToGoFuncArgs(args []Object) ([]interface{}, error) {
	funcArgs := []interface{}{}

	for _, arg := range args {
		v, ok := arg.(builtInType)

		if ok {
			if integer, ok := v.(*IntegerObject); ok {
				switch integer.flag {
				case integer64:
					funcArgs = append(funcArgs, int64(integer.value))
					continue
				case integer32:
					funcArgs = append(funcArgs, int32(integer.value))
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
