package vm

import (
	"math"
	"strconv"

	"github.com/goby-lang/goby/vm/classes"
	"github.com/goby-lang/goby/vm/errors"
)

// FloatObject represents an inexact real number using the native architecture's double-precision floating point
// representation.
//
// ```ruby
// '1.1'.to_f + '1.1'.to_f # => 2.2
// '2.1'.to_f * '2.1'.to_f # => 4.41
// ```
//
// - `Float.new` is not supported.
type FloatObject struct {
	*baseObj
	value float64
}

// Class methods --------------------------------------------------------
func builtinFloatClassMethods() []*BuiltinMethodObject {
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
func builtinFloatInstanceMethods() []*BuiltinMethodObject {
	return []*BuiltinMethodObject{
		{
			// Returns the sum of self and another Float.
			//
			// ```Ruby
			// 1 + 2 # => 3
			// ```
			// @return [Float]
			Name: "+",
			Fn: func(receiver Object) builtinMethodBody {
				return func(t *thread, args []Object, blockFrame *callFrame) Object {

					leftValue := receiver.(*FloatObject).value
					right, ok := args[0].(*FloatObject)

					if !ok {
						err := t.vm.initErrorObject(errors.TypeError, errors.WrongArgumentTypeFormat, classes.FloatClass, args[0].Class().Name)
						return err
					}

					rightValue := right.value
					return t.vm.initFloatObject(leftValue + rightValue)
				}
			},
		},
		{
			// Divides left hand operand by right hand operand and returns remainder.
			//
			// ```Ruby
			// 5 % 2 # => 1
			// ```
			// @return [Float]
			Name: "%",
			Fn: func(receiver Object) builtinMethodBody {
				return func(t *thread, args []Object, blockFrame *callFrame) Object {

					leftValue := receiver.(*FloatObject).value
					right, ok := args[0].(*FloatObject)

					if !ok {
						err := t.vm.initErrorObject(errors.TypeError, errors.WrongArgumentTypeFormat, classes.FloatClass, args[0].Class().Name)
						return err
					}

					rightValue := right.value
					result := math.Mod(leftValue, rightValue)
					return t.vm.initFloatObject(result)
				}
			},
		},
		{
			// Returns the subtraction of another Float from self.
			//
			// ```Ruby
			// 1 - 1 # => 0
			// ```
			// @return [Float]
			Name: "-",
			Fn: func(receiver Object) builtinMethodBody {
				return func(t *thread, args []Object, blockFrame *callFrame) Object {

					leftValue := receiver.(*FloatObject).value
					right, ok := args[0].(*FloatObject)

					if !ok {
						err := t.vm.initErrorObject(errors.TypeError, errors.WrongArgumentTypeFormat, classes.FloatClass, args[0].Class().Name)
						return err
					}

					rightValue := right.value
					return t.vm.initFloatObject(leftValue - rightValue)
				}
			},
		},
		{
			// Returns self multiplying another Float.
			//
			// ```Ruby
			// 2 * 10 # => 20
			// ```
			// @return [Float]
			Name: "*",
			Fn: func(receiver Object) builtinMethodBody {
				return func(t *thread, args []Object, blockFrame *callFrame) Object {

					leftValue := receiver.(*FloatObject).value
					right, ok := args[0].(*FloatObject)

					if !ok {
						err := t.vm.initErrorObject(errors.TypeError, errors.WrongArgumentTypeFormat, classes.FloatClass, args[0].Class().Name)
						return err
					}

					rightValue := right.value
					return t.vm.initFloatObject(leftValue * rightValue)
				}
			},
		},
		{
			// Returns self squaring another Float.
			//
			// ```Ruby
			// 2 ** 8 # => 256
			// ```
			// @return [Float]
			Name: "**",
			Fn: func(receiver Object) builtinMethodBody {
				return func(t *thread, args []Object, blockFrame *callFrame) Object {

					leftValue := receiver.(*FloatObject).value
					right, ok := args[0].(*FloatObject)

					if !ok {
						err := t.vm.initErrorObject(errors.TypeError, errors.WrongArgumentTypeFormat, classes.FloatClass, args[0].Class().Name)
						return err
					}

					rightValue := right.value
					result := math.Pow(leftValue, rightValue)
					return t.vm.initFloatObject(result)
				}
			},
		},
		{
			// Returns self divided by another Float.
			//
			// ```Ruby
			// 6 / 3 # => 2
			// ```
			// @return [Float]
			Name: "/",
			Fn: func(receiver Object) builtinMethodBody {
				return func(t *thread, args []Object, blockFrame *callFrame) Object {

					leftValue := receiver.(*FloatObject).value
					right, ok := args[0].(*FloatObject)

					if !ok {
						err := t.vm.initErrorObject(errors.TypeError, errors.WrongArgumentTypeFormat, classes.FloatClass, args[0].Class().Name)
						return err
					}

					rightValue := right.value
					return t.vm.initFloatObject(leftValue / rightValue)
				}
			},
		},
		{
			// Returns if self is larger than another Float.
			//
			// ```Ruby
			// 10 > -1 # => true
			// 3 > 3 # => false
			// ```
			// @return [Boolean]
			Name: ">",
			Fn: func(receiver Object) builtinMethodBody {
				return func(t *thread, args []Object, blockFrame *callFrame) Object {

					leftValue := receiver.(*FloatObject).value
					right, ok := args[0].(*FloatObject)

					if !ok {
						err := t.vm.initErrorObject(errors.TypeError, errors.WrongArgumentTypeFormat, classes.FloatClass, args[0].Class().Name)
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
			// Returns if self is larger than or equals to another Float.
			//
			// ```Ruby
			// 2 >= 1 # => true
			// 1 >= 1 # => true
			// ```
			// @return [Boolean]
			Name: ">=",
			Fn: func(receiver Object) builtinMethodBody {
				return func(t *thread, args []Object, blockFrame *callFrame) Object {

					leftValue := receiver.(*FloatObject).value
					right, ok := args[0].(*FloatObject)

					if !ok {
						err := t.vm.initErrorObject(errors.TypeError, errors.WrongArgumentTypeFormat, classes.FloatClass, args[0].Class().Name)
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
			// Returns if self is smaller than another Float.
			//
			// ```Ruby
			// 1 < 3 # => true
			// 1 < 1 # => false
			// ```
			// @return [Boolean]
			Name: "<",
			Fn: func(receiver Object) builtinMethodBody {
				return func(t *thread, args []Object, blockFrame *callFrame) Object {

					leftValue := receiver.(*FloatObject).value
					right, ok := args[0].(*FloatObject)

					if !ok {
						err := t.vm.initErrorObject(errors.TypeError, errors.WrongArgumentTypeFormat, classes.FloatClass, args[0].Class().Name)
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
			// Returns if self is smaller than or equals to another Float.
			//
			// ```Ruby
			// 1 <= 3 # => true
			// 1 <= 1 # => true
			// ```
			// @return [Boolean]
			Name: "<=",
			Fn: func(receiver Object) builtinMethodBody {
				return func(t *thread, args []Object, blockFrame *callFrame) Object {

					leftValue := receiver.(*FloatObject).value
					right, ok := args[0].(*FloatObject)

					if !ok {
						err := t.vm.initErrorObject(errors.TypeError, errors.WrongArgumentTypeFormat, classes.FloatClass, args[0].Class().Name)
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
			// Returns 1 if self is larger than the incoming Float, -1 if smaller. Otherwise 0.
			//
			// ```Ruby
			// 1 <=> 3 # => -1
			// 1 <=> 1 # => 0
			// 3 <=> 1 # => 1
			// ```
			// @return [Float]
			Name: "<=>",
			Fn: func(receiver Object) builtinMethodBody {
				return func(t *thread, args []Object, blockFrame *callFrame) Object {

					leftValue := receiver.(*FloatObject).value
					right, ok := args[0].(*FloatObject)

					if !ok {
						err := t.vm.initErrorObject(errors.TypeError, errors.WrongArgumentTypeFormat, classes.FloatClass, args[0].Class().Name)
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
			// Returns if self is equal to another Float.
			//
			// ```Ruby
			// 1 == 3 # => false
			// 1 == 1 # => true
			// ```
			// @return [Boolean]
			Name: "==",
			Fn: func(receiver Object) builtinMethodBody {
				return func(t *thread, args []Object, blockFrame *callFrame) Object {

					leftValue := receiver.(*FloatObject).value
					right, ok := args[0].(*FloatObject)

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
			// Returns if self is not equal to another Float.
			//
			// ```Ruby
			// 1 != 3 # => true
			// 1 != 1 # => false
			// ```
			// @return [Boolean]
			Name: "!=",
			Fn: func(receiver Object) builtinMethodBody {
				return func(t *thread, args []Object, blockFrame *callFrame) Object {

					leftValue := receiver.(*FloatObject).value
					right, ok := args[0].(*FloatObject)

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
			// Returns the `Integer` representation of self.
			//
			// ```Ruby
			// '100.1'.to_f.to_i # => 100
			// ```
			// @return [Integer]
			Name: "to_i",
			Fn: func(receiver Object) builtinMethodBody {
				return func(t *thread, args []Object, blockFrame *callFrame) Object {
					r := receiver.(*FloatObject)
					newInt := t.vm.initIntegerObject(int(r.value))
					newInt.flag = i
					return newInt
				}
			},
		},
		{
			Name: "ptr",
			Fn: func(receiver Object) builtinMethodBody {
				return func(t *thread, args []Object, blockFrame *callFrame) Object {
					r := receiver.(*FloatObject)
					return t.vm.initGoObject(&r.value)
				}
			},
		},
	}
}

// Internal functions ===================================================

// Functions for initialization -----------------------------------------

func (vm *VM) initFloatObject(value float64) *FloatObject {
	return &FloatObject{
		baseObj: &baseObj{class: vm.topLevelClass(classes.FloatClass)},
		value:   value,
	}
}

func (vm *VM) initFloatClass() *RClass {
	ic := vm.initializeClass(classes.FloatClass, false)
	ic.setBuiltinMethods(builtinFloatInstanceMethods(), false)
	ic.setBuiltinMethods(builtinFloatClassMethods(), true)
	return ic
}

// Polymorphic helper functions -----------------------------------------

// Value returns the object
func (f *FloatObject) Value() interface{} {
	return f.value
}

// toString returns the object's value as the string format, in non
// exponential format (straight number, without exponent `E<exp>`).
func (f *FloatObject) toString() string {
	return strconv.FormatFloat(f.value, 'f', -1, 64)
}

// toJSON just delegates to toString
func (f *FloatObject) toJSON() string {
	return f.toString()
}

// equal checks if the Float values between receiver and argument are equal
func (f *FloatObject) equal(e *FloatObject) bool {
	return f.value == e.value
}
