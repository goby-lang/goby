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
//
// - `Integer.new` is not supported.
type IntegerObject struct {
	Class *RInteger
	Value int
}

func (i *IntegerObject) toString() string {
	return strconv.Itoa(i.Value)
}

func (i *IntegerObject) toJSON() string {
	return i.toString()
}

func (i *IntegerObject) returnClass() Class {
	return i.Class
}

func (i *IntegerObject) equal(e *IntegerObject) bool {
	return i.Value == e.Value
}

func initIntegerObject(value int) *IntegerObject {
	return &IntegerObject{Value: value, Class: integerClass}
}

var builtinIntegerInstanceMethods = []*BuiltInMethodObject{
	{
		// Returns the sum of self and another Integer.
		//
		// ```Ruby
		// 1 + 2 # => 3
		// ```
		// @return [Integer]
		Name: "+",
		Fn: func(receiver Object) builtinMethodBody {
			return func(t *thread, args []Object, blockFrame *callFrame) Object {

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
		// Divides left hand operand by right hand operand and returns remainder.
		//
		// ```Ruby
		// 5 % 2 # => 1
		// ```
		// @return [Integer]
		Name: "%",
		Fn: func(receiver Object) builtinMethodBody {
			return func(t *thread, args []Object, blockFrame *callFrame) Object {

				leftValue := receiver.(*IntegerObject).Value
				right, ok := args[0].(*IntegerObject)

				if !ok {
					return wrongTypeError(integerClass)
				}

				rightValue := right.Value
				return &IntegerObject{Value: leftValue % rightValue, Class: integerClass}
			}
		},
	},
	{
		// Returns the subtraction of another Integer from self.
		//
		// ```Ruby
		// 1 - 1 # => 0
		// ```
		// @return [Integer]
		Name: "-",
		Fn: func(receiver Object) builtinMethodBody {
			return func(t *thread, args []Object, blockFrame *callFrame) Object {

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
		//
		// ```Ruby
		// 2 * 10 # => 20
		// ```
		// @return [Integer]
		Name: "*",
		Fn: func(receiver Object) builtinMethodBody {
			return func(t *thread, args []Object, blockFrame *callFrame) Object {

				leftValue := receiver.(*IntegerObject).Value
				right, ok := args[0].(*IntegerObject)

				if !ok {
					return initErrorObject(TypeErrorClass, "Expect Integer. got=%T (%+v)", args[0], args[0])
				}

				rightValue := right.Value
				return &IntegerObject{Value: leftValue * rightValue, Class: integerClass}
			}
		},
	},
	{
		// Returns self squaring another Integer.
		//
		// ```Ruby
		// 2 ** 8 # => 256
		// ```
		// @return [Integer]
		Name: "**",
		Fn: func(receiver Object) builtinMethodBody {
			return func(t *thread, args []Object, blockFrame *callFrame) Object {

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
		//
		// ```Ruby
		// 6 / 3 # => 2
		// ```
		// @return [Integer]
		Name: "/",
		Fn: func(receiver Object) builtinMethodBody {
			return func(t *thread, args []Object, blockFrame *callFrame) Object {

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
		//
		// ```Ruby
		// 10 > -1 # => true
		// 3 > 3 # => false
		// ```
		// @return [Boolean]
		Name: ">",
		Fn: func(receiver Object) builtinMethodBody {
			return func(t *thread, args []Object, blockFrame *callFrame) Object {

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
		//
		// ```Ruby
		// 2 >= 1 # => true
		// 1 >= 1 # => true
		// ```
		// @return [Boolean]
		Name: ">=",
		Fn: func(receiver Object) builtinMethodBody {
			return func(t *thread, args []Object, blockFrame *callFrame) Object {

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
		//
		// ```Ruby
		// 1 < 3 # => true
		// 1 < 1 # => false
		// ```
		// @return [Boolean]
		Name: "<",
		Fn: func(receiver Object) builtinMethodBody {
			return func(t *thread, args []Object, blockFrame *callFrame) Object {

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
		//
		// ```Ruby
		// 1 <= 3 # => true
		// 1 <= 1 # => true
		// ```
		// @return [Boolean]
		Name: "<=",
		Fn: func(receiver Object) builtinMethodBody {
			return func(t *thread, args []Object, blockFrame *callFrame) Object {

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
		//
		// ```Ruby
		// 1 <=> 3 # => -1
		// 1 <=> 1 # => 0
		// 3 <=> 1 # => 1
		// ```
		// @return [Integer]
		Name: "<=>",
		Fn: func(receiver Object) builtinMethodBody {
			return func(t *thread, args []Object, blockFrame *callFrame) Object {

				leftValue := receiver.(*IntegerObject).Value
				right, ok := args[0].(*IntegerObject)

				if !ok {
					return wrongTypeError(integerClass)
				}

				rightValue := right.Value

				if leftValue < rightValue {
					return initIntegerObject(-1)
				}
				if leftValue > rightValue {
					return initIntegerObject(1)
				}

				return initIntegerObject(0)
			}
		},
	},
	{
		// Returns if self is equal to another Integer.
		//
		// ```Ruby
		// 1 == 3 # => false
		// 1 == 1 # => true
		// ```
		// @return [Boolean]
		Name: "==",
		Fn: func(receiver Object) builtinMethodBody {
			return func(t *thread, args []Object, blockFrame *callFrame) Object {

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
		//
		// ```Ruby
		// 1 != 3 # => true
		// 1 != 1 # => false
		// ```
		// @return [Boolean]
		Name: "!=",
		Fn: func(receiver Object) builtinMethodBody {
			return func(t *thread, args []Object, blockFrame *callFrame) Object {

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
		//
		// ```Ruby
		// 1++ # => 2
		// ```
		// @return [Integer]
		Name: "++",
		Fn: func(receiver Object) builtinMethodBody {
			return func(t *thread, args []Object, blockFrame *callFrame) Object {
				int := receiver.(*IntegerObject)

				t.vm.Lock()
				defer t.vm.Unlock()

				int.Value++
				return int
			}
		},
	},
	{
		// Substracts 1 from self and returns.
		//
		// ```Ruby
		// 0-- # => -1
		// ```
		// @return [Integer]
		Name: "--",
		Fn: func(receiver Object) builtinMethodBody {
			return func(t *thread, args []Object, blockFrame *callFrame) Object {

				int := receiver.(*IntegerObject)

				t.vm.Lock()
				defer t.vm.Unlock()

				int.Value--
				return int
			}
		},
	},
	{
		// Returns a `String` representation of self.
		//
		// ```Ruby
		// 100.to_s # => "100"
		// ```
		// @return [String]
		Name: "to_s",
		Fn: func(receiver Object) builtinMethodBody {
			return func(t *thread, args []Object, blockFrame *callFrame) Object {

				int := receiver.(*IntegerObject)

				return initStringObject(strconv.Itoa(int.Value))
			}
		},
	},
	{
		// Returns self.
		//
		// ```Ruby
		// 100.to_i # => 100
		// ```
		// @return [Integer]
		Name: "to_i",
		Fn: func(receiver Object) builtinMethodBody {
			return func(t *thread, args []Object, blockFrame *callFrame) Object {
				return receiver
			}
		},
	},
	{
		// Returns if self is even.
		//
		// ```Ruby
		// 1.even # => false
		// 2.even # => true
		// ```
		// @return [Boolean]
		Name: "even",
		Fn: func(receiver Object) builtinMethodBody {
			return func(t *thread, args []Object, blockFrame *callFrame) Object {

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
		// @return [Boolean]
		Name: "odd",
		Fn: func(receiver Object) builtinMethodBody {
			return func(t *thread, args []Object, blockFrame *callFrame) Object {

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
		//
		// ```ruby
		// 100.next # => 101
		// ```
		// @return [Integer]
		Name: "next",
		Fn: func(receiver Object) builtinMethodBody {
			return func(t *thread, args []Object, blockFrame *callFrame) Object {
				i := receiver.(*IntegerObject)
				return initIntegerObject(i.Value + 1)
			}
		},
	},
	{
		// Returns self - 1.
		//
		// ```ruby
		// 40.pred # => 39
		// ```
		// @return [Integer]
		Name: "pred",
		Fn: func(receiver Object) builtinMethodBody {
			return func(t *thread, args []Object, blockFrame *callFrame) Object {
				i := receiver.(*IntegerObject)
				return initIntegerObject(i.Value - 1)
			}
		},
	},
	{
		// Yields a block a number of times equals to self.
		//
		// ```Ruby
		// a = 0
		// 3.times do
		//    a++
		// end
		// a # => 3
		// ```
		Name: "times",
		Fn: func(receiver Object) builtinMethodBody {
			return func(t *thread, args []Object, blockFrame *callFrame) Object {
				n := receiver.(*IntegerObject)

				if n.Value < 0 {
					return newError("Expect paramentr to be greater 0. got=%d", n.Value)
				}

				if blockFrame == nil {
					return newError("Can't yield without a block")
				}

				for i := 0; i < n.Value; i++ {
					t.builtInMethodYield(blockFrame, initIntegerObject(i))
				}

				return n
			}
		},
	},
}

func initIntegerClass() {
	bc := &BaseClass{Name: "Integer", Methods: newEnvironment(), ClassMethods: newEnvironment(), Class: classClass, pseudoSuperClass: objectClass, superClass: objectClass}
	ic := &RInteger{BaseClass: bc}
	ic.setBuiltInMethods(builtinIntegerInstanceMethods, false)
	integerClass = ic
}
