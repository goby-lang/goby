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
// 1.1 + 1.1 # => 2.2
// 2.1 * 2.1 # => 4.41
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
			Fn: func(receiver Object, instruction *instruction) builtinMethodBody {
				return func(t *thread, args []Object, blockFrame *normalCallFrame) Object {
					return t.initUnsupportedMethodError(instruction.sourceLine, "#new", receiver)
				}
			},
		},
	}
}

// Instance methods -----------------------------------------------------
func builtinFloatInstanceMethods() []*BuiltinMethodObject {
	return []*BuiltinMethodObject{
		{
			// Returns the sum of self and a Numeric.
			//
			// ```Ruby
			// 1.1 + 2 # => 3.1
			// ```
			// @return [Float]
			Name: "+",
			Fn: func(receiver Object, instruction *instruction) builtinMethodBody {
				return func(t *thread, args []Object, blockFrame *normalCallFrame) Object {
					operation := func(leftValue float64, rightValue float64) float64 {
						return leftValue + rightValue
					}

					return receiver.(*FloatObject).arithmeticOperation(t, args[0], operation, instruction)
				}
			},
		},
		{
			// Returns the modulo between self and a Numeric.
			//
			// ```Ruby
			// 5.5 % 2 # => 1.5
			// ```
			// @return [Float]
			Name: "%",
			Fn: func(receiver Object, instruction *instruction) builtinMethodBody {
				return func(t *thread, args []Object, blockFrame *normalCallFrame) Object {
					operation := func(leftValue float64, rightValue float64) float64 {
						return math.Mod(leftValue, rightValue)
					}

					return receiver.(*FloatObject).arithmeticOperation(t, args[0], operation, instruction)
				}
			},
		},
		{
			// Returns the subtraction of a Numeric from self.
			//
			// ```Ruby
			// 1.5 - 1 # => 0.5
			// ```
			// @return [Float]
			Name: "-",
			Fn: func(receiver Object, instruction *instruction) builtinMethodBody {
				return func(t *thread, args []Object, blockFrame *normalCallFrame) Object {
					operation := func(leftValue float64, rightValue float64) float64 {
						return leftValue - rightValue
					}

					return receiver.(*FloatObject).arithmeticOperation(t, args[0], operation, instruction)
				}
			},
		},
		{
			// Returns self multiplying a Numeric.
			//
			// ```Ruby
			// 2.5 * 10 # => 25.0
			// ```
			// @return [Float]
			Name: "*",
			Fn: func(receiver Object, instruction *instruction) builtinMethodBody {
				return func(t *thread, args []Object, blockFrame *normalCallFrame) Object {
					operation := func(leftValue float64, rightValue float64) float64 {
						return leftValue * rightValue
					}

					return receiver.(*FloatObject).arithmeticOperation(t, args[0], operation, instruction)
				}
			},
		},
		{
			// Returns self squaring a Numeric.
			//
			// ```Ruby
			// 4.0 ** 2.5 # => 32.0
			// ```
			// @return [Float]
			Name: "**",
			Fn: func(receiver Object, instruction *instruction) builtinMethodBody {
				return func(t *thread, args []Object, blockFrame *normalCallFrame) Object {
					operation := func(leftValue float64, rightValue float64) float64 {
						return math.Pow(leftValue, rightValue)
					}

					return receiver.(*FloatObject).arithmeticOperation(t, args[0], operation, instruction)
				}
			},
		},
		{
			// Returns self divided by a Numeric.
			//
			// ```Ruby
			// 7.5 / 3 # => 2.5
			// ```
			// @return [Float]
			Name: "/",
			Fn: func(receiver Object, instruction *instruction) builtinMethodBody {
				return func(t *thread, args []Object, blockFrame *normalCallFrame) Object {
					operation := func(leftValue float64, rightValue float64) float64 {
						return leftValue / rightValue
					}

					return receiver.(*FloatObject).arithmeticOperation(t, args[0], operation, instruction)
				}
			},
		},
		{
			// Returns if self is larger than a Numeric.
			//
			// ```Ruby
			// 10 > -1 # => true
			// 3 > 3 # => false
			// ```
			// @return [Boolean]
			Name: ">",
			Fn: func(receiver Object, instruction *instruction) builtinMethodBody {
				return func(t *thread, args []Object, blockFrame *normalCallFrame) Object {
					operation := func(leftValue float64, rightValue float64) bool {
						return leftValue > rightValue
					}

					return receiver.(*FloatObject).numericComparison(t, args[0], operation, instruction)
				}
			},
		},
		{
			// Returns if self is larger than or equals to a Numeric.
			//
			// ```Ruby
			// 2 >= 1 # => true
			// 1 >= 1 # => true
			// ```
			// @return [Boolean]
			Name: ">=",
			Fn: func(receiver Object, instruction *instruction) builtinMethodBody {
				return func(t *thread, args []Object, blockFrame *normalCallFrame) Object {
					operation := func(leftValue float64, rightValue float64) bool {
						return leftValue >= rightValue
					}

					return receiver.(*FloatObject).numericComparison(t, args[0], operation, instruction)
				}
			},
		},
		{
			// Returns if self is smaller than a Numeric.
			//
			// ```Ruby
			// 1 < 3 # => true
			// 1 < 1 # => false
			// ```
			// @return [Boolean]
			Name: "<",
			Fn: func(receiver Object, instruction *instruction) builtinMethodBody {
				return func(t *thread, args []Object, blockFrame *normalCallFrame) Object {
					operation := func(leftValue float64, rightValue float64) bool {
						return leftValue < rightValue
					}

					return receiver.(*FloatObject).numericComparison(t, args[0], operation, instruction)
				}
			},
		},
		{
			// Returns if self is smaller than or equals to a Numeric.
			//
			// ```Ruby
			// 1 <= 3 # => true
			// 1 <= 1 # => true
			// ```
			// @return [Boolean]
			Name: "<=",
			Fn: func(receiver Object, instruction *instruction) builtinMethodBody {
				return func(t *thread, args []Object, blockFrame *normalCallFrame) Object {
					operation := func(leftValue float64, rightValue float64) bool {
						return leftValue <= rightValue
					}

					return receiver.(*FloatObject).numericComparison(t, args[0], operation, instruction)
				}
			},
		},
		{
			// Returns 1 if self is larger than a Numeric, -1 if smaller. Otherwise 0.
			//
			// ```Ruby
			// 1.5 <=> 3 # => -1
			// 1.0 <=> 1 # => 0
			// 3.5 <=> 1 # => 1
			// ```
			// @return [Float]
			Name: "<=>",
			Fn: func(receiver Object, instruction *instruction) builtinMethodBody {
				return func(t *thread, args []Object, blockFrame *normalCallFrame) Object {
					rightNumeric, ok := args[0].(Numeric)

					if !ok {
						return t.vm.initErrorObject(errors.TypeError, instruction.sourceLine, errors.WrongArgumentTypeFormat, "Numeric", args[0].Class().Name)
					}

					leftValue := receiver.(*FloatObject).value
					rightValue := rightNumeric.floatValue()

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
			// Returns if self is equal to an Object.
			// If the Object is a Numeric, a comparison is performed, otherwise, the
			// result is always false.
			//
			// ```Ruby
			// 1.0 == 3     # => false
			// 1.0 == 1     # => true
			// 1.0 == '1.0' # => false
			// ```
			// @return [Boolean]
			Name: "==",
			Fn: func(receiver Object, instruction *instruction) builtinMethodBody {
				return func(t *thread, args []Object, blockFrame *normalCallFrame) Object {
					result := receiver.(*FloatObject).equalityTest(args[0])

					return toBooleanObject(result)
				}
			},
		},
		{
			// Returns if self is not equal to an Object.
			// If the Object is a Numeric, a comparison is performed, otherwise, the
			// result is always true.
			//
			// ```Ruby
			// 1.0 != 3     # => true
			// 1.0 != 1     # => false
			// 1.0 != '1.0' # => true
			// ```
			// @return [Boolean]
			Name: "!=",
			Fn: func(receiver Object, instruction *instruction) builtinMethodBody {
				return func(t *thread, args []Object, blockFrame *normalCallFrame) Object {
					result := !receiver.(*FloatObject).equalityTest(args[0])

					return toBooleanObject(result)
				}
			},
		},
		{
			// Returns the `Integer` representation of self.
			//
			// ```Ruby
			// 100.1.to_i # => 100
			// ```
			// @return [Integer]
			Name: "to_i",
			Fn: func(receiver Object, instruction *instruction) builtinMethodBody {
				return func(t *thread, args []Object, blockFrame *normalCallFrame) Object {
					r := receiver.(*FloatObject)
					newInt := t.vm.initIntegerObject(int(r.value))
					newInt.flag = i
					return newInt
				}
			},
		},
		{
			Name: "ptr",
			Fn: func(receiver Object, instruction *instruction) builtinMethodBody {
				return func(t *thread, args []Object, blockFrame *normalCallFrame) Object {
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

// Numeric interface
func (f *FloatObject) floatValue() float64 {
	return f.value
}

// TODO: Remove instruction argument
// Apply the passed arithmetic operation, while performing type conversion.
func (f *FloatObject) arithmeticOperation(t *thread, rightObject Object, operation func(leftValue float64, rightValue float64) float64, instruction *instruction) Object {
	rightNumeric, ok := rightObject.(Numeric)

	if !ok {
		return t.vm.initErrorObject(errors.TypeError, instruction.sourceLine, errors.WrongArgumentTypeFormat, "Numeric", rightObject.Class().Name)
	}

	leftValue := f.value
	rightValue := rightNumeric.floatValue()

	result := operation(leftValue, rightValue)

	return t.vm.initFloatObject(result)
}

// Apply an equality test, returning true if the objects are considered equal,
// and false otherwise.
func (f *FloatObject) equalityTest(rightObject Object) bool {
	rightNumeric, ok := rightObject.(Numeric)

	if !ok {
		return false
	}

	leftValue := f.value
	rightValue := rightNumeric.floatValue()

	return leftValue == rightValue
}

// TODO: Remove instruction argument
// Apply the passed numeric comparison, while performing type conversion.
func (f *FloatObject) numericComparison(t *thread, rightObject Object, operation func(leftValue float64, rightValue float64) bool, instruction *instruction) Object {
	rightNumeric, ok := rightObject.(Numeric)

	if !ok {
		return t.vm.initErrorObject(errors.TypeError, instruction.sourceLine, errors.WrongArgumentTypeFormat, "Numeric", rightObject.Class().Name)
	}

	leftValue := f.value
	rightValue := rightNumeric.floatValue()

	result := operation(leftValue, rightValue)

	return toBooleanObject(result)
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
