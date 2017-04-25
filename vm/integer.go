package vm

import (
	"fmt"
)

var (
	integerClass *RInteger
)

// RInteger is integer class
type RInteger struct {
	*BaseClass
}

// IntegerObject represents integer instances
type IntegerObject struct {
	Class *RInteger
	Value int
}

func (i *IntegerObject) Type() objectType {
	return integerObj
}

func (i *IntegerObject) Inspect() string {
	return fmt.Sprintf("%d", i.Value)
}

func (i *IntegerObject) ReturnClass() Class {
	return i.Class
}

func initilaizeInteger(value int) *IntegerObject {
	return &IntegerObject{Value: value, Class: integerClass}
}

var builtinIntegerMethods = []*BuiltInMethod{
	{
		Fn: func(receiver Object) builtinMethodBody {
			return func(vm *VM, args []Object, blockFrame *callFrame) Object {

				err := checkArgumentLen(args, integerClass, "+")

				if err != nil {
					return err
				}

				leftValue := receiver.(*IntegerObject).Value
				right, ok := args[0].(*IntegerObject)

				if !ok {
					return wrongTypeError(integerClass)
				}

				rightValue := right.Value
				return &IntegerObject{Value: leftValue + rightValue, Class: integerClass}
			}
		},
		Name: "+",
	},
	{
		Fn: func(receiver Object) builtinMethodBody {
			return func(vm *VM, args []Object, blockFrame *callFrame) Object {

				err := checkArgumentLen(args, integerClass, "-")

				if err != nil {
					return err
				}

				leftValue := receiver.(*IntegerObject).Value
				right, ok := args[0].(*IntegerObject)

				if !ok {
					return wrongTypeError(integerClass)
				}

				rightValue := right.Value
				return &IntegerObject{Value: leftValue - rightValue, Class: integerClass}
			}
		},
		Name: "-",
	},
	{
		Fn: func(receiver Object) builtinMethodBody {
			return func(vm *VM, args []Object, blockFrame *callFrame) Object {

				err := checkArgumentLen(args, integerClass, "+")

				if err != nil {
					return err
				}

				leftValue := receiver.(*IntegerObject).Value
				right, ok := args[0].(*IntegerObject)

				if !ok {
					return wrongTypeError(integerClass)
				}

				rightValue := right.Value
				return &IntegerObject{Value: leftValue * rightValue, Class: integerClass}
			}
		},
		Name: "*",
	},
	{
		Fn: func(receiver Object) builtinMethodBody {
			return func(vm *VM, args []Object, blockFrame *callFrame) Object {

				err := checkArgumentLen(args, integerClass, "+")

				if err != nil {
					return err
				}

				leftValue := receiver.(*IntegerObject).Value
				right, ok := args[0].(*IntegerObject)

				if !ok {
					return wrongTypeError(integerClass)
				}

				rightValue := right.Value
				return &IntegerObject{Value: leftValue / rightValue, Class: integerClass}
			}
		},
		Name: "/",
	},
	{
		Fn: func(receiver Object) builtinMethodBody {
			return func(vm *VM, args []Object, blockFrame *callFrame) Object {

				err := checkArgumentLen(args, integerClass, ">")
				if err != nil {
					return err
				}

				leftValue := receiver.(*IntegerObject).Value
				right, ok := args[0].(*IntegerObject)

				if !ok {
					return wrongTypeError(integerClass)
				}

				rightValue := right.Value

				if leftValue > rightValue {
					return TRUE
				}

				return FALSE
			}
		},
		Name: ">",
	},
	{
		Fn: func(receiver Object) builtinMethodBody {
			return func(vm *VM, args []Object, blockFrame *callFrame) Object {

				err := checkArgumentLen(args, integerClass, "<")
				if err != nil {
					return err
				}

				leftValue := receiver.(*IntegerObject).Value
				right, ok := args[0].(*IntegerObject)

				if !ok {
					return wrongTypeError(integerClass)
				}

				rightValue := right.Value

				if leftValue < rightValue {
					return TRUE
				}

				return FALSE
			}
		},
		Name: "<",
	},
	{
		Fn: func(receiver Object) builtinMethodBody {
			return func(vm *VM, args []Object, blockFrame *callFrame) Object {

				err := checkArgumentLen(args, integerClass, "==")

				if err != nil {
					return err
				}

				leftValue := receiver.(*IntegerObject).Value
				right, ok := args[0].(*IntegerObject)

				if !ok {
					return wrongTypeError(integerClass)
				}

				rightValue := right.Value

				if leftValue == rightValue {
					return TRUE
				}

				return FALSE
			}
		},
		Name: "==",
	},
	{
		Fn: func(receiver Object) builtinMethodBody {
			return func(vm *VM, args []Object, blockFrame *callFrame) Object {

				err := checkArgumentLen(args, integerClass, "!=")

				if err != nil {
					return err
				}

				leftValue := receiver.(*IntegerObject).Value
				right, ok := args[0].(*IntegerObject)

				if !ok {
					return wrongTypeError(integerClass)
				}

				rightValue := right.Value

				if leftValue != rightValue {
					return TRUE
				}

				return FALSE
			}
		},
		Name: "!=",
	},
	{
		Fn: func(receiver Object) builtinMethodBody {
			return func(vm *VM, args []Object, blockFrame *callFrame) Object {

				if len(args) > 0 {
					return &Error{Message: "Too many arguments for Integer#++"}
				}

				int := receiver.(*IntegerObject)
				int.Value += 1
				return int
			}
		},
		Name: "++",
	},
	{
		Fn: func(receiver Object) builtinMethodBody {
			return func(vm *VM, args []Object, blockFrame *callFrame) Object {

				if len(args) > 0 {
					return &Error{Message: "Too many arguments for Integer#--"}
				}

				int := receiver.(*IntegerObject)
				int.Value -= 1
				return int
			}
		},
		Name: "--",
	},
	{
		Fn: func(receiver Object) builtinMethodBody {
			return func(vm *VM, args []Object, blockFrame *callFrame) Object {

				if len(args) > 0 {
					return &Error{Message: "Too many arguments for Integer#--"}
				}

				int := receiver.(*IntegerObject)
				return initializeString(fmt.Sprint(int.Value))
			}
		},
		Name: "to_s",
	},
}

func initInteger() {
	methods := NewEnvironment()

	for _, m := range builtinIntegerMethods {
		methods.Set(m.Name, m)
	}

	bc := &BaseClass{Name: "Integer", Methods: methods, ClassMethods: NewEnvironment(), Class: classClass, SuperClass: objectClass}
	ic := &RInteger{BaseClass: bc}
	integerClass = ic
}
