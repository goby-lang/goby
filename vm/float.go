package vm

import (
	"math"

	"github.com/goby-lang/goby/vm/classes"
	"github.com/goby-lang/goby/vm/errors"
	"strconv"
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
			Fn: func(receiver Object, sourceLine int) builtinMethodBody {
				return func(t *thread, args []Object, blockFrame *normalCallFrame) Object {
					return t.vm.initUnsupportedMethodError(sourceLine, "#new", receiver)
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
			Fn: func(receiver Object, sourceLine int) builtinMethodBody {
				return func(t *thread, args []Object, blockFrame *normalCallFrame) Object {
					operation := func(leftValue float64, rightValue float64) float64 {
						return leftValue + rightValue
					}

					return receiver.(*FloatObject).arithmeticOperation(t, args[0], operation, sourceLine, false)
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
			Fn: func(receiver Object, sourceLine int) builtinMethodBody {
				return func(t *thread, args []Object, blockFrame *normalCallFrame) Object {
					operation := func(leftValue float64, rightValue float64) float64 {
						return math.Mod(leftValue, rightValue)
					}

					return receiver.(*FloatObject).arithmeticOperation(t, args[0], operation, sourceLine, true)
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
			Fn: func(receiver Object, sourceLine int) builtinMethodBody {
				return func(t *thread, args []Object, blockFrame *normalCallFrame) Object {
					operation := func(leftValue float64, rightValue float64) float64 {
						return leftValue - rightValue
					}

					return receiver.(*FloatObject).arithmeticOperation(t, args[0], operation, sourceLine, false)
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
			Fn: func(receiver Object, sourceLine int) builtinMethodBody {
				return func(t *thread, args []Object, blockFrame *normalCallFrame) Object {
					operation := func(leftValue float64, rightValue float64) float64 {
						return leftValue * rightValue
					}

					return receiver.(*FloatObject).arithmeticOperation(t, args[0], operation, sourceLine, false)
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
			Fn: func(receiver Object, sourceLine int) builtinMethodBody {
				return func(t *thread, args []Object, blockFrame *normalCallFrame) Object {
					operation := func(leftValue float64, rightValue float64) float64 {
						return math.Pow(leftValue, rightValue)
					}

					return receiver.(*FloatObject).arithmeticOperation(t, args[0], operation, sourceLine, false)
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
			Fn: func(receiver Object, sourceLine int) builtinMethodBody {
				return func(t *thread, args []Object, blockFrame *normalCallFrame) Object {
					operation := func(leftValue float64, rightValue float64) float64 {
						return leftValue / rightValue
					}

					return receiver.(*FloatObject).arithmeticOperation(t, args[0], operation, sourceLine, true)
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
			Fn: func(receiver Object, sourceLine int) builtinMethodBody {
				return func(t *thread, args []Object, blockFrame *normalCallFrame) Object {
					operation := func(leftValue float64, rightValue float64) bool {
						return leftValue > rightValue
					}

					return receiver.(*FloatObject).numericComparison(t, args[0], operation, sourceLine)
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
			Fn: func(receiver Object, sourceLine int) builtinMethodBody {
				return func(t *thread, args []Object, blockFrame *normalCallFrame) Object {
					operation := func(leftValue float64, rightValue float64) bool {
						return leftValue >= rightValue
					}

					return receiver.(*FloatObject).numericComparison(t, args[0], operation, sourceLine)
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
			Fn: func(receiver Object, sourceLine int) builtinMethodBody {
				return func(t *thread, args []Object, blockFrame *normalCallFrame) Object {
					operation := func(leftValue float64, rightValue float64) bool {
						return leftValue < rightValue
					}

					return receiver.(*FloatObject).numericComparison(t, args[0], operation, sourceLine)
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
			Fn: func(receiver Object, sourceLine int) builtinMethodBody {
				return func(t *thread, args []Object, blockFrame *normalCallFrame) Object {
					operation := func(leftValue float64, rightValue float64) bool {
						return leftValue <= rightValue
					}

					return receiver.(*FloatObject).numericComparison(t, args[0], operation, sourceLine)
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
			Fn: func(receiver Object, sourceLine int) builtinMethodBody {
				return func(t *thread, args []Object, blockFrame *normalCallFrame) Object {
					rightNumeric, ok := args[0].(Numeric)

					if !ok {
						return t.vm.initErrorObject(errors.TypeError, sourceLine, errors.WrongArgumentTypeFormat, "Numeric", args[0].Class().Name)
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
			Fn: func(receiver Object, sourceLine int) builtinMethodBody {
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
			Fn: func(receiver Object, sourceLine int) builtinMethodBody {
				return func(t *thread, args []Object, blockFrame *normalCallFrame) Object {
					result := !receiver.(*FloatObject).equalityTest(args[0])

					return toBooleanObject(result)
				}
			},
		},
		{
			// Converts the Integer object into Decimal object and returns it.
			// Each digit of the float is literally transferred to the corresponding digit
			// of the Decimal, via a string representation of the float.
			//
			// ```Ruby
			// "100.1".to_f.to_d # => 100.1
			//
			// a = "3.14159265358979".to_f
			// b = a.to_d #=> 3.14159265358979
			// ```
			//
			// @return [Decimal]
			Name: "to_d",
			Fn: func(receiver Object, sourceLine int) builtinMethodBody {
				return func(t *thread, args []Object, blockFrame *normalCallFrame) Object {

					if len(args) != 0 {
						return t.vm.initErrorObject(errors.ArgumentError, sourceLine, "Expect 0 argument. got=%v", strconv.Itoa(len(args)))
					}

					fl := receiver.(*FloatObject).value
					fs := strconv.FormatFloat(fl, 'f', -1, 64)
					de, err := new(Decimal).SetString(fs)
					if err == false {
						return t.vm.initErrorObject(errors.ArgumentError, sourceLine, "Invalid numeric string. got=%v", fs)
					}

					return t.vm.initDecimalObject(de)
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
			Fn: func(receiver Object, sourceLine int) builtinMethodBody {
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
			Fn: func(receiver Object, sourceLine int) builtinMethodBody {
				return func(t *thread, args []Object, blockFrame *normalCallFrame) Object {
					r := receiver.(*FloatObject)
					return t.vm.initGoObject(&r.value)
				}
			},
		},
		{
			// Returns the Float as a positive value.
			//
			// ```Ruby
			// -34.56.abs # => 34.56
			// 34.56.abs # => 34.56
			// ```
			// @return [Float]
			Name: "abs",
			Fn: func(receiver Object, sourceLine int) builtinMethodBody {
				return func(t *thread, args []Object, blockFrame *normalCallFrame) Object {
					if len(args) != 0 {
						return t.vm.initErrorObject(errors.ArgumentError, sourceLine, "Expect 0 argument. got=%v", strconv.Itoa(len(args)))
					}
					r := receiver.(*FloatObject)
					result := math.Abs(r.value)
					return t.vm.initFloatObject(result)
				}
			},
		},
		{
			// Returns the smallest Integer greater than or equal to self.
			//
			// ```Ruby
			// 1.2.ceil  # => 2
			// 2.ceil    # => 2
			// -1.2.ceil # => -1
			// -2.ceil   # => -2
			// ```
			// @return [Integer]
			Name: "ceil",
			Fn: func(receiver Object, sourceLine int) builtinMethodBody {
				return func(t *thread, args []Object, blockFrame *normalCallFrame) Object {
					// TODO: Make ceil accept arguments
					if len(args) != 0 {
						return t.vm.initErrorObject(errors.ArgumentError, sourceLine, "Expect 0 argument. got=%v", strconv.Itoa(len(args)))
					}
					r := receiver.(*FloatObject)
					result := math.Ceil(r.value)
					newInt := t.vm.initIntegerObject(int(result))
					newInt.flag = i
					return newInt
				}
			},
		},
		{
			// Returns the largest Integer less than or equal to self.
			//
			// ```Ruby
			// 1.2.floor  # => 1
			// 2.0.floor  # => 2
			// -1.2.floor # => -2
			// -2.0.floor # => -2
			// ```
			// @return [Integer]
			Name: "floor",
			Fn: func(receiver Object, sourceLine int) builtinMethodBody {
				return func(t *thread, args []Object, blockFrame *normalCallFrame) Object {
					// TODO: Make floor accept arguments
					if len(args) != 0 {
						return t.vm.initErrorObject(errors.ArgumentError, sourceLine, "Expect 0 argument. got=%v", strconv.Itoa(len(args)))
					}
					r := receiver.(*FloatObject)
					result := math.Floor(r.value)
					newInt := t.vm.initIntegerObject(int(result))
					newInt.flag = i
					return newInt
				}
			},
		},
		{
			// Returns true if Float is equal to 0.0
			//
			// ```Ruby
			// 0.0.zero? # => true
			// 1.0.zero? # => false
			// ```
			// @return [Boolean]
			Name: "zero?",
			Fn: func(receiver Object, sourceLine int) builtinMethodBody {
				return func(t *thread, args []Object, blockFrame *normalCallFrame) Object {
					if len(args) != 0 {
						return t.vm.initErrorObject(errors.ArgumentError, sourceLine, "Expect 0 argument. got=%v", strconv.Itoa(len(args)))
					}
					r := receiver.(*FloatObject)
					return toBooleanObject(r.value == 0.0)
				}
			},
		},
		{
			// Returns true if Float is larger than 0.0
			//
			// ```Ruby
			// -1.0.positive? # => false
			// 0.0.positive?  # => false
			// 1.0.positive?  # => true
			// ```
			// @return [Boolean]
			Name: "positive?",
			Fn: func(receiver Object, sourceLine int) builtinMethodBody {
				return func(t *thread, args []Object, blockFrame *normalCallFrame) Object {
					if len(args) != 0 {
						return t.vm.initErrorObject(errors.ArgumentError, sourceLine, "Expect 0 argument. got=%v", strconv.Itoa(len(args)))
					}
					r := receiver.(*FloatObject)
					return toBooleanObject(r.value > 0.0)
				}
			},
		},
		{
			// Returns true if Float is less than 0.0
			//
			// ```Ruby
			// -1.0.negative? # => true
			// 0.0.negative?  # => false
			// 1.0.negative?  # => false
			// ```
			// @return [Boolean]
			Name: "negative?",
			Fn: func(receiver Object, sourceLine int) builtinMethodBody {
				return func(t *thread, args []Object, blockFrame *normalCallFrame) Object {
					if len(args) != 0 {
						return t.vm.initErrorObject(errors.ArgumentError, sourceLine, "Expect 0 argument. got=%v", strconv.Itoa(len(args)))
					}
					r := receiver.(*FloatObject)
					return toBooleanObject(r.value < 0.0)
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
func (f *FloatObject) arithmeticOperation(t *thread, rightObject Object, operation func(leftValue float64, rightValue float64) float64, sourceLine int, division bool) Object {
	rightNumeric, ok := rightObject.(Numeric)

	if !ok {
		return t.vm.initErrorObject(errors.TypeError, sourceLine, errors.WrongArgumentTypeFormat, "Numeric", rightObject.Class().Name)
	}

	leftValue := f.value
	rightValue := rightNumeric.floatValue()

	if division && rightValue == 0 {
		return t.vm.initErrorObject(errors.ZeroDivisionError, sourceLine, errors.DividedByZero)
	}

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
func (f *FloatObject) numericComparison(t *thread, rightObject Object, operation func(leftValue float64, rightValue float64) bool, sourceLine int) Object {
	rightNumeric, ok := rightObject.(Numeric)

	if !ok {
		return t.vm.initErrorObject(errors.TypeError, sourceLine, errors.WrongArgumentTypeFormat, "Numeric", rightObject.Class().Name)
	}

	leftValue := f.value
	rightValue := rightNumeric.floatValue()

	result := operation(leftValue, rightValue)

	return toBooleanObject(result)
}

// toString returns the object's value as the string format, in non
// exponential format (straight number, without exponent `E<exp>`).
func (f *FloatObject) toString() string {
	return strconv.FormatFloat(f.value, 'f', -1, 64) // fmt.Sprintf("%f", f.value)
}

// toJSON just delegates to toString
func (f *FloatObject) toJSON(t *thread) string {
	return f.toString()
}

// equal checks if the Float values between receiver and argument are equal
func (f *FloatObject) equal(e *FloatObject) bool {
	return f.value == e.value
}
