package vm

import (
	"math"
	"strconv"
)

var (
	integerClass *RInteger
)

// RInteger is integer class
type RInteger struct {
	*BaseClass
}

// IntegerObject represents number objects which can bring into mathematical calculations.
//
// ```ruby
// 1 + 1 # => 2
// 2 * 2 # => 4
// ```
type IntegerObject struct {
	Class *RInteger
	Value int
}

func (i *IntegerObject) objectType() objectType {
	return integerObj
}

func (i *IntegerObject) Inspect() string {
	return strconv.Itoa(i.Value)
}

func (i *IntegerObject) returnClass() Class {
	return i.Class
}

func (i *IntegerObject) equal(e *IntegerObject) bool {
	return i.Value == e.Value
}

func initilaizeInteger(value int) *IntegerObject {
	return &IntegerObject{Value: value, Class: integerClass}
}

var builtinIntegerMethods = []*BuiltInMethod{
	{
		// Returns the sum of self and another Integer.
		Name: "+",
		Fn: func(receiver Object) builtinMethodBody {
			return func(vm *VM, args []Object, blockFrame *callFrame) Object {

				leftValue := receiver.(*IntegerObject).Value
				right, ok := args[0].(*IntegerObject)

				if !ok {
					return wrongTypeError(integerClass)
				}

				rightValue := right.Value
				return &IntegerObject{Value: leftValue + rightValue, Class: integerClass}
			}
		},
	},
	{
		// Returns the subtraction of another Integer from self.
		Name: "-",
		Fn: func(receiver Object) builtinMethodBody {
			return func(vm *VM, args []Object, blockFrame *callFrame) Object {

				leftValue := receiver.(*IntegerObject).Value
				right, ok := args[0].(*IntegerObject)

				if !ok {
					return wrongTypeError(integerClass)
				}

				rightValue := right.Value
				return &IntegerObject{Value: leftValue - rightValue, Class: integerClass}
			}
		},
	},
	{
		// Returns self multiplying another Integer.
		Name: "*",
		Fn: func(receiver Object) builtinMethodBody {
			return func(vm *VM, args []Object, blockFrame *callFrame) Object {

				leftValue := receiver.(*IntegerObject).Value
				right, ok := args[0].(*IntegerObject)

				if !ok {
					return wrongTypeError(integerClass)
				}

				rightValue := right.Value
				return &IntegerObject{Value: leftValue * rightValue, Class: integerClass}
			}
		},
	},
	{
		// Returns self squaring another Integer.
		Name: "**",
		Fn: func(receiver Object) builtinMethodBody {
			return func(vm *VM, args []Object, blockFrame *callFrame) Object {

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
	},
	{
		// Returns self divided by another Integer.
		Name: "/",
		Fn: func(receiver Object) builtinMethodBody {
			return func(vm *VM, args []Object, blockFrame *callFrame) Object {

				leftValue := receiver.(*IntegerObject).Value
				right, ok := args[0].(*IntegerObject)

				if !ok {
					return wrongTypeError(integerClass)
				}

				rightValue := right.Value
				return &IntegerObject{Value: leftValue / rightValue, Class: integerClass}
			}
		},
	},
	{
		// Returns if self is larger than another Integer.
		Name: ">",
		Fn: func(receiver Object) builtinMethodBody {
			return func(vm *VM, args []Object, blockFrame *callFrame) Object {

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
	},
	{
		// Returns if self is larger than or equals to another Integer.
		Name: ">=",
		Fn: func(receiver Object) builtinMethodBody {
			return func(vm *VM, args []Object, blockFrame *callFrame) Object {

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
	},
	{
		// Returns if self is smaller than another Integer.
		Name: "<",
		Fn: func(receiver Object) builtinMethodBody {
			return func(vm *VM, args []Object, blockFrame *callFrame) Object {

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
	},
	{
		// Returns if self is smaller than or equals to another Integer.
		Name: "<=",
		Fn: func(receiver Object) builtinMethodBody {
			return func(vm *VM, args []Object, blockFrame *callFrame) Object {

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
	},
	{
		// Returns 1 if self is larger than the incoming Integer, -1 if smaller. Otherwise 0.
		Name: "<=>",
		Fn: func(receiver Object) builtinMethodBody {
			return func(vm *VM, args []Object, blockFrame *callFrame) Object {

				leftValue := receiver.(*IntegerObject).Value
				right, ok := args[0].(*IntegerObject)

				if !ok {
					return wrongTypeError(integerClass)
				}

				rightValue := right.Value

				if leftValue < rightValue {
					return initilaizeInteger(-1)
				}
				if leftValue > rightValue {
					return initilaizeInteger(1)
				}

				return initilaizeInteger(0)
			}
		},
	},
	{
		// Returns if self is equal to another Integer.
		Name: "==",
		Fn: func(receiver Object) builtinMethodBody {
			return func(vm *VM, args []Object, blockFrame *callFrame) Object {

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
	},
	{
		// Returns if self is not equal to another Integer.
		Name: "!=",
		Fn: func(receiver Object) builtinMethodBody {
			return func(vm *VM, args []Object, blockFrame *callFrame) Object {

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
	},
	{
		// Adds 1 to self and returns.
		Name: "++",
		Fn: func(receiver Object) builtinMethodBody {
			return func(vm *VM, args []Object, blockFrame *callFrame) Object {

				int := receiver.(*IntegerObject)
				int.Value++
				return int
			}
		},
	},
	{
		// Substracts 1 from self and returns.
		Name: "--",
		Fn: func(receiver Object) builtinMethodBody {
			return func(vm *VM, args []Object, blockFrame *callFrame) Object {

				int := receiver.(*IntegerObject)
				int.Value--
				return int
			}
		},
	},
	{
		// Returns a `String` representation of self.
		Name: "to_s",
		Fn: func(receiver Object) builtinMethodBody {
			return func(vm *VM, args []Object, blockFrame *callFrame) Object {

				int := receiver.(*IntegerObject)

				return initializeString(strconv.Itoa(int.Value))
			}
		},
	},
	{
		// Returns self.
		Name: "to_i",
		Fn: func(receiver Object) builtinMethodBody {
			return func(vm *VM, args []Object, blockFrame *callFrame) Object {
				return receiver
			}
		},
	},
	{
		// Returns if self is even.
		Name: "even",
		Fn: func(receiver Object) builtinMethodBody {
			return func(vm *VM, args []Object, blockFrame *callFrame) Object {

				i := receiver.(*IntegerObject)
				even := i.Value%2 == 0

				if even {
					return TRUE
				}

				return FALSE
			}
		},
	},
	{
		// Returns if self is odd.
		//
		// ```ruby
		// 3.odd # => true
		// 4.odd # => false
		// ```
		Name: "odd",
		Fn: func(receiver Object) builtinMethodBody {
			return func(vm *VM, args []Object, blockFrame *callFrame) Object {

				i := receiver.(*IntegerObject)
				odd := i.Value%2 != 0
				if odd {
					return TRUE
				}

				return FALSE
			}
		},
	},
	{
		// Returns self + 1.
		Name: "next",
		Fn: func(receiver Object) builtinMethodBody {
			return func(vm *VM, args []Object, blockFrame *callFrame) Object {
				i := receiver.(*IntegerObject)
				return initilaizeInteger(i.Value + 1)
			}
		},
	},
	{
		// Returns self - 1.
		Name: "prev",
		Fn: func(receiver Object) builtinMethodBody {
			return func(vm *VM, args []Object, blockFrame *callFrame) Object {
				i := receiver.(*IntegerObject)
				return initilaizeInteger(i.Value - 1)
			}
		},
	},
	{
		// Yields a block a number of times equals to self.
		Name: "times",
		Fn: func(receiver Object) builtinMethodBody {
			return func(vm *VM, args []Object, blockFrame *callFrame) Object {
				n := receiver.(*IntegerObject)

				if n.Value < 0 {
					return newError("Expect paramentr to be greater 0. got=%d", n.Value)
				}

				if blockFrame == nil {
					return newError("Can't yield without a block")
				}

				for i := 0; i < n.Value; i++ {
					builtInMethodYield(vm, blockFrame, initilaizeInteger(i))
				}

				return n
			}
		},
	},
}

func initInteger() {
	methods := newEnvironment()

	for _, m := range builtinIntegerMethods {
		methods.set(m.Name, m)
	}

	bc := &BaseClass{Name: "Integer", Methods: methods, ClassMethods: newEnvironment(), Class: classClass, pseudoSuperClass: objectClass, superClass: objectClass}
	ic := &RInteger{BaseClass: bc}
	integerClass = ic
}
