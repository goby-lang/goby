package vm

import (
	"fmt"
	"math"
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

func (i *IntegerObject) objectType() objectType {
	return integerObj
}

func (i *IntegerObject) Inspect() string {
	return fmt.Sprintf("%d", i.Value)
}

func (i *IntegerObject) returnClass() Class {
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

				err := checkArgumentLen(args, integerClass, "**")

				if err != nil {
					return err
				}

				leftValue := receiver.(*IntegerObject).Value
				right, ok := args[0].(*IntegerObject)

				if !ok {
					return wrongTypeError(integerClass)
				}

				rightValue := right.Value
				result := math.Pow(float64(leftValue), float64(rightValue))
				return &IntegerObject{Value: int(result), Class: integerClass}
			}
		},
		Name: "**",
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

				err := checkArgumentLen(args, integerClass, ">=")
				if err != nil {
					return err
				}

				leftValue := receiver.(*IntegerObject).Value
				right, ok := args[0].(*IntegerObject)

				if !ok {
					return wrongTypeError(integerClass)
				}

				rightValue := right.Value

				if leftValue >= rightValue {
					return TRUE
				}

				return FALSE
			}
		},
		Name: ">=",
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

				err := checkArgumentLen(args, integerClass, "<=")
				if err != nil {
					return err
				}

				leftValue := receiver.(*IntegerObject).Value
				right, ok := args[0].(*IntegerObject)

				if !ok {
					return wrongTypeError(integerClass)
				}

				rightValue := right.Value

				if leftValue <= rightValue {
					return TRUE
				}

				return FALSE
			}
		},
		Name: "<=",
	},
	{
		Fn: func(receiver Object) builtinMethodBody {
			return func(vm *VM, args []Object, blockFrame *callFrame) Object {

				err := checkArgumentLen(args, integerClass, "<=")
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
					return initilaizeInteger(-1)
				} else if leftValue > rightValue {
					return initilaizeInteger(1)
				} else {
					return initilaizeInteger(0)
				}
			}
		},
		Name: "<=>",
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
				int.Value++
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
				int.Value--
				return int
			}
		},
		Name: "--",
	},
	{
		Fn: func(receiver Object) builtinMethodBody {
			return func(vm *VM, args []Object, blockFrame *callFrame) Object {

				if len(args) > 0 {
					return &Error{Message: "Too many arguments for Integer#to_s"}
				}

				int := receiver.(*IntegerObject)
				return initializeString(fmt.Sprint(int.Value))
			}
		},
		Name: "to_s",
	},
	{
		Fn: func(receiver Object) builtinMethodBody {
			return func(vm *VM, args []Object, blockFrame *callFrame) Object {
				if len(args) > 0 {
					return &Error{Message: "Too many arguments for Integer#even"}
				}

				i := receiver.(*IntegerObject)
				even := i.Value%2 == 0

				if even {
					return TRUE
				}

				return FALSE
			}
		},
		Name: "even",
	},
	{
		Fn: func(receiver Object) builtinMethodBody {
			return func(vm *VM, args []Object, blockFrame *callFrame) Object {
				if len(args) > 0 {
					return &Error{Message: "Too many arguments for Integer#odd"}
				}

				i := receiver.(*IntegerObject)
				odd := i.Value%2 != 0
				if odd {
					return TRUE
				}

				return FALSE
			}
		},
		Name: "odd",
	},
}

func initInteger() {
	methods := newEnvironment()

	for _, m := range builtinIntegerMethods {
		methods.set(m.Name, m)
	}

	bc := &BaseClass{Name: "Integer", Methods: methods, ClassMethods: newEnvironment(), Class: classClass, SuperClass: objectClass}
	ic := &RInteger{BaseClass: bc}
	integerClass = ic
}
