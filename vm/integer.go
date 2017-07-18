package vm

import (
	"math"
	"strconv"
)

// IntegerObject represents number objects which can bring into mathematical calculations.
//
// ```ruby
// 1 + 1 # => 2
// 2 * 2 # => 4
// ```
//
// - `Integer.new` is not supported.
type IntegerObject struct {
	*baseObj
	Value int
}

func (vm *VM) initIntegerObject(value int) *IntegerObject {
	return &IntegerObject{
		baseObj: &baseObj{class: vm.topLevelClass(integerClass)},
		Value:   value,
	}
}

func (vm *VM) initIntegerClass() *RClass {
	ic := vm.initializeClass(integerClass, false)
	ic.setBuiltInMethods(builtinIntegerInstanceMethods(), false)
	ic.setBuiltInMethods(builtInIntegerClassMethods(), true)
	return ic
}

func builtInIntegerClassMethods() []*BuiltInMethodObject {
	return []*BuiltInMethodObject{
		{
			Name: "new",
			Fn: func(receiver Object) builtinMethodBody {
				return func(t *thread, args []Object, blockFrame *callFrame) Object {
					return t.UnsupportedMethodError("#new", receiver)
				}
			},
		},
	}
}

func builtinIntegerInstanceMethods() []*BuiltInMethodObject {
	return []*BuiltInMethodObject{
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
						err := t.vm.initErrorObject(TypeError, WrongArgumentTypeFormat, integerClass, args[0].Class().Name)
						return err
					}

					rightValue := right.Value
					return t.vm.initIntegerObject(leftValue + rightValue)
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
						err := t.vm.initErrorObject(TypeError, WrongArgumentTypeFormat, integerClass, args[0].Class().Name)
						return err
					}

					rightValue := right.Value
					return t.vm.initIntegerObject(leftValue % rightValue)
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
						err := t.vm.initErrorObject(TypeError, WrongArgumentTypeFormat, integerClass, args[0].Class().Name)
						return err
					}

					rightValue := right.Value
					return t.vm.initIntegerObject(leftValue - rightValue)
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
						err := t.vm.initErrorObject(TypeError, WrongArgumentTypeFormat, integerClass, args[0].Class().Name)
						return err
					}

					rightValue := right.Value
					return t.vm.initIntegerObject(leftValue * rightValue)
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
						err := t.vm.initErrorObject(TypeError, WrongArgumentTypeFormat, integerClass, args[0].Class().Name)
						return err
					}

					rightValue := right.Value
					result := math.Pow(float64(leftValue), float64(rightValue))
					return t.vm.initIntegerObject(int(result))
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
						err := t.vm.initErrorObject(TypeError, WrongArgumentTypeFormat, integerClass, args[0].Class().Name)
						return err
					}

					rightValue := right.Value
					return t.vm.initIntegerObject(leftValue / rightValue)
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
						err := t.vm.initErrorObject(TypeError, WrongArgumentTypeFormat, integerClass, args[0].Class().Name)
						return err
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
						err := t.vm.initErrorObject(TypeError, WrongArgumentTypeFormat, integerClass, args[0].Class().Name)
						return err
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
						err := t.vm.initErrorObject(TypeError, WrongArgumentTypeFormat, integerClass, args[0].Class().Name)
						return err
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
						err := t.vm.initErrorObject(TypeError, WrongArgumentTypeFormat, integerClass, args[0].Class().Name)
						return err
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
						err := t.vm.initErrorObject(TypeError, WrongArgumentTypeFormat, integerClass, args[0].Class().Name)
						return err
					}

					rightValue := right.Value

					if leftValue < rightValue {
						return t.vm.initIntegerObject(-1)
					}
					if leftValue > rightValue {
						return t.vm.initIntegerObject(1)
					}

					return t.vm.initIntegerObject(0)
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
						err := t.vm.initErrorObject(TypeError, WrongArgumentTypeFormat, integerClass, args[0].Class().Name)
						return err
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
						err := t.vm.initErrorObject(TypeError, WrongArgumentTypeFormat, integerClass, args[0].Class().Name)
						return err
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

					return t.vm.initStringObject(strconv.Itoa(int.Value))
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
					return t.vm.initIntegerObject(i.Value + 1)
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
					return t.vm.initIntegerObject(i.Value - 1)
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
						t.builtInMethodYield(blockFrame, t.vm.initIntegerObject(i))
					}

					return n
				}
			},
		},
	}
}

// Polymorphic helper functions -----------------------------------------

// toString converts the receiver into string.
func (i *IntegerObject) toString() string {
	return strconv.Itoa(i.Value)
}

// toJSON converts the receiver into JSON string.
func (i *IntegerObject) toJSON() string {
	return i.toString()
}

func (i *IntegerObject) equal(e *IntegerObject) bool {
	return i.Value == e.Value
}
