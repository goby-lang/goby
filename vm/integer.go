package vm

import (
	"math"
	"strconv"

	"github.com/goby-lang/goby/vm/classes"
	"github.com/goby-lang/goby/vm/errors"
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
	value int
	flag  int
}

/*
This is enum defined for integer's flag
*/
const (
	_ int = iota
	ui
	ui8
	ui16
	ui32
	ui64

	i
	i8
	i16
	i32
	i64

	f32
	f64
)

// Class methods --------------------------------------------------------
func builtinIntegerClassMethods() []*BuiltinMethodObject {
	return []*BuiltinMethodObject{
		{
			Name: "new",
			Fn: func(receiver Object) builtinMethodBody {
				return func(t *thread, args []Object, blockFrame *callFrame) Object {
					return t.unsupportedMethodError("#new", receiver)
				}
			},
		},
	}
}

// Instance methods -----------------------------------------------------
func builtinIntegerInstanceMethods() []*BuiltinMethodObject {
	return []*BuiltinMethodObject{
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

					leftValue := receiver.(*IntegerObject).value
					right, ok := args[0].(*IntegerObject)

					if !ok {
						err := t.vm.initErrorObject(errors.TypeError, errors.WrongArgumentTypeFormat, classes.IntegerClass, args[0].Class().Name)
						return err
					}

					rightValue := right.value
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

					leftValue := receiver.(*IntegerObject).value
					right, ok := args[0].(*IntegerObject)

					if !ok {
						err := t.vm.initErrorObject(errors.TypeError, errors.WrongArgumentTypeFormat, classes.IntegerClass, args[0].Class().Name)
						return err
					}

					rightValue := right.value
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

					leftValue := receiver.(*IntegerObject).value
					right, ok := args[0].(*IntegerObject)

					if !ok {
						err := t.vm.initErrorObject(errors.TypeError, errors.WrongArgumentTypeFormat, classes.IntegerClass, args[0].Class().Name)
						return err
					}

					rightValue := right.value
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

					leftValue := receiver.(*IntegerObject).value
					right, ok := args[0].(*IntegerObject)

					if !ok {
						err := t.vm.initErrorObject(errors.TypeError, errors.WrongArgumentTypeFormat, classes.IntegerClass, args[0].Class().Name)
						return err
					}

					rightValue := right.value
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

					leftValue := receiver.(*IntegerObject).value
					right, ok := args[0].(*IntegerObject)

					if !ok {
						err := t.vm.initErrorObject(errors.TypeError, errors.WrongArgumentTypeFormat, classes.IntegerClass, args[0].Class().Name)
						return err
					}

					rightValue := right.value
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

					leftValue := receiver.(*IntegerObject).value
					right, ok := args[0].(*IntegerObject)

					if !ok {
						err := t.vm.initErrorObject(errors.TypeError, errors.WrongArgumentTypeFormat, classes.IntegerClass, args[0].Class().Name)
						return err
					}

					rightValue := right.value
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

					leftValue := receiver.(*IntegerObject).value
					right, ok := args[0].(*IntegerObject)

					if !ok {
						err := t.vm.initErrorObject(errors.TypeError, errors.WrongArgumentTypeFormat, classes.IntegerClass, args[0].Class().Name)
						return err
					}

					rightValue := right.value

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

					leftValue := receiver.(*IntegerObject).value
					right, ok := args[0].(*IntegerObject)

					if !ok {
						err := t.vm.initErrorObject(errors.TypeError, errors.WrongArgumentTypeFormat, classes.IntegerClass, args[0].Class().Name)
						return err
					}

					rightValue := right.value

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

					leftValue := receiver.(*IntegerObject).value
					right, ok := args[0].(*IntegerObject)

					if !ok {
						err := t.vm.initErrorObject(errors.TypeError, errors.WrongArgumentTypeFormat, classes.IntegerClass, args[0].Class().Name)
						return err
					}

					rightValue := right.value

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

					leftValue := receiver.(*IntegerObject).value
					right, ok := args[0].(*IntegerObject)

					if !ok {
						err := t.vm.initErrorObject(errors.TypeError, errors.WrongArgumentTypeFormat, classes.IntegerClass, args[0].Class().Name)
						return err
					}

					rightValue := right.value

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

					leftValue := receiver.(*IntegerObject).value
					right, ok := args[0].(*IntegerObject)

					if !ok {
						err := t.vm.initErrorObject(errors.TypeError, errors.WrongArgumentTypeFormat, classes.IntegerClass, args[0].Class().Name)
						return err
					}

					rightValue := right.value

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

					leftValue := receiver.(*IntegerObject).value
					right, ok := args[0].(*IntegerObject)

					if !ok {
						return FALSE
					}

					rightValue := right.value

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

					leftValue := receiver.(*IntegerObject).value
					right, ok := args[0].(*IntegerObject)

					if !ok {
						return TRUE
					}

					rightValue := right.value

					if leftValue != rightValue {
						return TRUE
					}

					return FALSE
				}
			},
		},
		{
			// Returns if self is even.
			//
			// ```Ruby
			// 1.even? # => false
			// 2.even? # => true
			// ```
			// @return [Boolean]
			Name: "even?",
			Fn: func(receiver Object) builtinMethodBody {
				return func(t *thread, args []Object, blockFrame *callFrame) Object {

					i := receiver.(*IntegerObject)
					even := i.value%2 == 0

					if even {
						return TRUE
					}

					return FALSE
				}
			},
		},
		// Returns the `Float` conversion of self.
		//
		// ```Ruby
		// 100.to_f # => '100.0'.to_f
		// ```
		// @return [Float]
		{
			Name: "to_f",
			Fn: func(receiver Object) builtinMethodBody {
				return func(t *thread, args []Object, blockFrame *callFrame) Object {
					r := receiver.(*IntegerObject)
					newFloat := t.vm.initFloatObject(float64(r.value))
					return newFloat
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

					return t.vm.initStringObject(strconv.Itoa(int.value))
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
					return t.vm.initIntegerObject(i.value + 1)
				}
			},
		},
		{
			// Returns if self is odd.
			//
			// ```ruby
			// 3.odd? # => true
			// 4.odd? # => false
			// ```
			// @return [Boolean]
			Name: "odd?",
			Fn: func(receiver Object) builtinMethodBody {
				return func(t *thread, args []Object, blockFrame *callFrame) Object {

					i := receiver.(*IntegerObject)
					odd := i.value%2 != 0
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
					return t.vm.initIntegerObject(i.value - 1)
				}
			},
		},
		{
			// Yields a block a number of times equals to self.
			//
			// ```Ruby
			// a = 0
			// 3.times do
			//    a += 1
			// end
			// a # => 3
			// ```
			Name: "times",
			Fn: func(receiver Object) builtinMethodBody {
				return func(t *thread, args []Object, blockFrame *callFrame) Object {
					n := receiver.(*IntegerObject)

					if n.value < 0 {
						return t.vm.initErrorObject(errors.InternalError, "Expect integer greater than or equal 0. got: %d", n.value)
					}

					if blockFrame == nil {
						return t.vm.initErrorObject(errors.InternalError, errors.CantYieldWithoutBlockFormat)
					}

					for i := 0; i < n.value; i++ {
						t.builtinMethodYield(blockFrame, t.vm.initIntegerObject(i))
					}

					return n
				}
			},
		},
		{
			Name: "to_int",
			Fn: func(receiver Object) builtinMethodBody {
				return func(t *thread, args []Object, blockFrame *callFrame) Object {
					r := receiver.(*IntegerObject)
					newInt := t.vm.initIntegerObject(r.value)
					newInt.flag = i
					return newInt
				}
			},
		},
		{
			Name: "to_int8",
			Fn: func(receiver Object) builtinMethodBody {
				return func(t *thread, args []Object, blockFrame *callFrame) Object {
					r := receiver.(*IntegerObject)
					newInt := t.vm.initIntegerObject(r.value)
					newInt.flag = i8
					return newInt
				}
			},
		},
		{
			Name: "to_int16",
			Fn: func(receiver Object) builtinMethodBody {
				return func(t *thread, args []Object, blockFrame *callFrame) Object {
					r := receiver.(*IntegerObject)
					newInt := t.vm.initIntegerObject(r.value)
					newInt.flag = i16
					return newInt
				}
			},
		},
		{
			Name: "to_int32",
			Fn: func(receiver Object) builtinMethodBody {
				return func(t *thread, args []Object, blockFrame *callFrame) Object {
					r := receiver.(*IntegerObject)
					newInt := t.vm.initIntegerObject(r.value)
					newInt.flag = i32
					return newInt
				}
			},
		},
		{
			Name: "to_int64",
			Fn: func(receiver Object) builtinMethodBody {
				return func(t *thread, args []Object, blockFrame *callFrame) Object {
					r := receiver.(*IntegerObject)
					newInt := t.vm.initIntegerObject(r.value)
					newInt.flag = i64
					return newInt
				}
			},
		},
		{
			Name: "to_uint",
			Fn: func(receiver Object) builtinMethodBody {
				return func(t *thread, args []Object, blockFrame *callFrame) Object {
					r := receiver.(*IntegerObject)
					newInt := t.vm.initIntegerObject(r.value)
					newInt.flag = ui
					return newInt
				}
			},
		},
		{
			Name: "to_uint8",
			Fn: func(receiver Object) builtinMethodBody {
				return func(t *thread, args []Object, blockFrame *callFrame) Object {
					r := receiver.(*IntegerObject)
					newInt := t.vm.initIntegerObject(r.value)
					newInt.flag = ui8
					return newInt
				}
			},
		},
		{
			Name: "to_uint16",
			Fn: func(receiver Object) builtinMethodBody {
				return func(t *thread, args []Object, blockFrame *callFrame) Object {
					r := receiver.(*IntegerObject)
					newInt := t.vm.initIntegerObject(r.value)
					newInt.flag = ui16
					return newInt
				}
			},
		},
		{
			Name: "to_uint32",
			Fn: func(receiver Object) builtinMethodBody {
				return func(t *thread, args []Object, blockFrame *callFrame) Object {
					r := receiver.(*IntegerObject)
					newInt := t.vm.initIntegerObject(r.value)
					newInt.flag = ui32
					return newInt
				}
			},
		},
		{
			Name: "to_uint64",
			Fn: func(receiver Object) builtinMethodBody {
				return func(t *thread, args []Object, blockFrame *callFrame) Object {
					r := receiver.(*IntegerObject)
					newInt := t.vm.initIntegerObject(r.value)
					newInt.flag = ui64
					return newInt
				}
			},
		},
		{
			Name: "to_float32",
			Fn: func(receiver Object) builtinMethodBody {
				return func(t *thread, args []Object, blockFrame *callFrame) Object {
					r := receiver.(*IntegerObject)
					newInt := t.vm.initIntegerObject(r.value)
					newInt.flag = f32
					return newInt
				}
			},
		},
		{
			Name: "to_float64",
			Fn: func(receiver Object) builtinMethodBody {
				return func(t *thread, args []Object, blockFrame *callFrame) Object {
					r := receiver.(*IntegerObject)
					newInt := t.vm.initIntegerObject(r.value)
					newInt.flag = f64
					return newInt
				}
			},
		},
		{
			Name: "ptr",
			Fn: func(receiver Object) builtinMethodBody {
				return func(t *thread, args []Object, blockFrame *callFrame) Object {
					r := receiver.(*IntegerObject)
					return t.vm.initGoObject(&r.value)
				}
			},
		},
	}
}

// Internal functions ===================================================

// Functions for initialization -----------------------------------------

func (vm *VM) initIntegerObject(value int) *IntegerObject {
	return &IntegerObject{
		baseObj: &baseObj{class: vm.topLevelClass(classes.IntegerClass)},
		value:   value,
		flag:    i,
	}
}

func (vm *VM) initIntegerClass() *RClass {
	ic := vm.initializeClass(classes.IntegerClass, false)
	ic.setBuiltinMethods(builtinIntegerInstanceMethods(), false)
	ic.setBuiltinMethods(builtinIntegerClassMethods(), true)
	return ic
}

// Polymorphic helper functions -----------------------------------------

// Value returns the object
func (i *IntegerObject) Value() interface{} {
	return i.value
}

// toString returns the object's name as the string format
func (i *IntegerObject) toString() string {
	return strconv.Itoa(i.value)
}

// toJSON just delegates to toString
func (i *IntegerObject) toJSON() string {
	return i.toString()
}

// equal checks if the integer values between receiver and argument are equal
func (i *IntegerObject) equal(e *IntegerObject) bool {
	return i.value == e.value
}
