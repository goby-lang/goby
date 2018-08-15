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
	*BaseObj
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
			Fn: func(receiver Object, sourceLine int, t *Thread, args []Object, blockFrame *normalCallFrame) Object {
				return t.vm.InitNoMethodError(sourceLine, "new", receiver)

			},
		},
	}
}

// Instance methods -----------------------------------------------------
func builtinIntegerInstanceMethods() []*BuiltinMethodObject {
	return []*BuiltinMethodObject{
		{
			// Returns the sum of self and another Numeric.
			//
			// ```Ruby
			// 1 + 2 # => 3
			// ```
			// @return [Numeric]
			Name: "+",
			Fn: func(receiver Object, sourceLine int, t *Thread, args []Object, blockFrame *normalCallFrame) Object {
				intOperation := func(leftValue int, rightValue int) int {
					return leftValue + rightValue
				}
				floatOperation := func(leftValue float64, rightValue float64) float64 {
					return leftValue + rightValue
				}

				return receiver.(*IntegerObject).arithmeticOperation(t, args[0], intOperation, floatOperation, sourceLine, false)

			},
		},
		{
			// Divides left hand operand by right hand operand and returns remainder.
			//
			// ```Ruby
			// 5 % 2 # => 1
			// ```
			// @return [Numeric]
			Name: "%",
			Fn: func(receiver Object, sourceLine int, t *Thread, args []Object, blockFrame *normalCallFrame) Object {
				intOperation := func(leftValue int, rightValue int) int {
					return leftValue % rightValue
				}
				floatOperation := func(leftValue float64, rightValue float64) float64 {
					return math.Mod(leftValue, rightValue)
				}

				return receiver.(*IntegerObject).arithmeticOperation(t, args[0], intOperation, floatOperation, sourceLine, true)

			},
		},
		{
			// Returns the subtraction of another Numeric from self.
			//
			// ```Ruby
			// 1 - 1 # => 0
			// ```
			// @return [Numeric]
			Name: "-",
			Fn: func(receiver Object, sourceLine int, t *Thread, args []Object, blockFrame *normalCallFrame) Object {
				intOperation := func(leftValue int, rightValue int) int {
					return leftValue - rightValue
				}
				floatOperation := func(leftValue float64, rightValue float64) float64 {
					return leftValue - rightValue
				}

				return receiver.(*IntegerObject).arithmeticOperation(t, args[0], intOperation, floatOperation, sourceLine, false)

			},
		},
		{
			// Returns self multiplying another Numeric.
			//
			// ```Ruby
			// 2 * 10 # => 20
			// ```
			// @return [Numeric]
			Name: "*",
			Fn: func(receiver Object, sourceLine int, t *Thread, args []Object, blockFrame *normalCallFrame) Object {
				intOperation := func(leftValue int, rightValue int) int {
					return leftValue * rightValue
				}
				floatOperation := func(leftValue float64, rightValue float64) float64 {
					return leftValue * rightValue
				}

				return receiver.(*IntegerObject).arithmeticOperation(t, args[0], intOperation, floatOperation, sourceLine, false)

			},
		},
		{
			// Returns self squaring another Numeric.
			//
			// ```Ruby
			// 2 ** 8 # => 256
			// ```
			// @return [Numeric]
			Name: "**",
			Fn: func(receiver Object, sourceLine int, t *Thread, args []Object, blockFrame *normalCallFrame) Object {
				intOperation := func(leftValue int, rightValue int) int {
					return int(math.Pow(float64(leftValue), float64(rightValue)))
				}
				floatOperation := func(leftValue float64, rightValue float64) float64 {
					return math.Pow(leftValue, rightValue)
				}

				return receiver.(*IntegerObject).arithmeticOperation(t, args[0], intOperation, floatOperation, sourceLine, false)

			},
		},
		{
			// Returns self divided by another Numeric.
			//
			// ```Ruby
			// 6 / 3 # => 2
			// ```
			// @return [Numeric]
			Name: "/",
			Fn: func(receiver Object, sourceLine int, t *Thread, args []Object, blockFrame *normalCallFrame) Object {

				intOperation := func(leftValue int, rightValue int) int {
					return leftValue / rightValue
				}
				floatOperation := func(leftValue float64, rightValue float64) float64 {
					return leftValue / rightValue
				}

				return receiver.(*IntegerObject).arithmeticOperation(t, args[0], intOperation, floatOperation, sourceLine, true)

			},
		},
		{
			// Returns if self is larger than another Numeric.
			//
			// ```Ruby
			// 10 > -1 # => true
			// 3 > 3 # => false
			// ```
			// @return [Boolean]
			Name: ">",
			Fn: func(receiver Object, sourceLine int, t *Thread, args []Object, blockFrame *normalCallFrame) Object {
				intComparison := func(leftValue int, rightValue int) bool {
					return leftValue > rightValue
				}
				floatComparison := func(leftValue float64, rightValue float64) bool {
					return leftValue > rightValue
				}

				return receiver.(*IntegerObject).numericComparison(t, args[0], intComparison, floatComparison, sourceLine)

			},
		},
		{
			// Returns if self is larger than or equals to another Numeric.
			//
			// ```Ruby
			// 2 >= 1 # => true
			// 1 >= 1 # => true
			// ```
			// @return [Boolean]
			Name: ">=",
			Fn: func(receiver Object, sourceLine int, t *Thread, args []Object, blockFrame *normalCallFrame) Object {
				intComparison := func(leftValue int, rightValue int) bool {
					return leftValue >= rightValue
				}
				floatComparison := func(leftValue float64, rightValue float64) bool {
					return leftValue >= rightValue
				}

				return receiver.(*IntegerObject).numericComparison(t, args[0], intComparison, floatComparison, sourceLine)

			},
		},
		{
			// Returns if self is smaller than another Numeric.
			//
			// ```Ruby
			// 1 < 3 # => true
			// 1 < 1 # => false
			// ```
			// @return [Boolean]
			Name: "<",
			Fn: func(receiver Object, sourceLine int, t *Thread, args []Object, blockFrame *normalCallFrame) Object {
				intComparison := func(leftValue int, rightValue int) bool {
					return leftValue < rightValue
				}
				floatComparison := func(leftValue float64, rightValue float64) bool {
					return leftValue < rightValue
				}

				return receiver.(*IntegerObject).numericComparison(t, args[0], intComparison, floatComparison, sourceLine)

			},
		},
		{
			// Returns if self is smaller than or equals to another Numeric.
			//
			// ```Ruby
			// 1 <= 3 # => true
			// 1 <= 1 # => true
			// ```
			// @return [Boolean]
			Name: "<=",
			Fn: func(receiver Object, sourceLine int, t *Thread, args []Object, blockFrame *normalCallFrame) Object {
				intComparison := func(leftValue int, rightValue int) bool {
					return leftValue <= rightValue
				}
				floatComparison := func(leftValue float64, rightValue float64) bool {
					return leftValue <= rightValue
				}

				return receiver.(*IntegerObject).numericComparison(t, args[0], intComparison, floatComparison, sourceLine)

			},
		},
		{
			// Returns 1 if self is larger than the incoming Numeric, -1 if smaller. Otherwise 0.
			//
			// ```Ruby
			// 1 <=> 3 # => -1
			// 1 <=> 1 # => 0
			// 3 <=> 1 # => 1
			// ```
			// @return [Integer]
			Name: "<=>",
			Fn: func(receiver Object, sourceLine int, t *Thread, args []Object, blockFrame *normalCallFrame) Object {
				rightObject := args[0]

				switch rightObject.(type) {
				case *IntegerObject:
					leftValue := receiver.(*IntegerObject).value
					rightValue := rightObject.(*IntegerObject).value

					if leftValue < rightValue {
						return t.vm.InitIntegerObject(-1)
					}
					if leftValue > rightValue {
						return t.vm.InitIntegerObject(1)
					}

					return t.vm.InitIntegerObject(0)
				case *FloatObject:
					leftValue := float64(receiver.(*IntegerObject).value)
					rightValue := rightObject.(*FloatObject).value

					if leftValue < rightValue {
						return t.vm.InitIntegerObject(-1)
					}
					if leftValue > rightValue {
						return t.vm.InitIntegerObject(1)
					}

					return t.vm.InitIntegerObject(0)
				default:
					return t.vm.InitErrorObject(errors.TypeError, sourceLine, errors.WrongArgumentTypeFormat, "Numeric", rightObject.Class().Name)
				}

			},
		},
		{
			// Returns if self is equal to an Object.
			// If the Object is a Numeric, a comparison is performed, otherwise, the
			// result is always false.
			//
			// ```Ruby
			// 1 == 3   # => false
			// 1 == 1   # => true
			// 1 == '1' # => false
			// ```
			// @return [Boolean]
			Name: "==",
			Fn: func(receiver Object, sourceLine int, t *Thread, args []Object, blockFrame *normalCallFrame) Object {
				result := receiver.(*IntegerObject).equalityTest(args[0])

				return toBooleanObject(result)

			},
		},
		{
			// Returns if self is not equal to an Object.
			// If the Object is a Numeric, a comparison is performed, otherwise, the
			// result is always true.
			//
			// ```Ruby
			// 1 != 3   # => true
			// 1 != 1   # => false
			// 1 != '1' # => true
			// ```
			// @return [Boolean]
			Name: "!=",
			Fn: func(receiver Object, sourceLine int, t *Thread, args []Object, blockFrame *normalCallFrame) Object {
				result := !receiver.(*IntegerObject).equalityTest(args[0])

				return toBooleanObject(result)

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
			Fn: func(receiver Object, sourceLine int, t *Thread, args []Object, blockFrame *normalCallFrame) Object {

				i := receiver.(*IntegerObject)
				even := i.value%2 == 0

				if even {
					return TRUE
				}

				return FALSE

			},
		},
		// Returns the `Decimal` conversion of self.
		//
		// ```Ruby
		// 100.to_d # => '100'.to_d
		// ```
		// @return [Decimal]
		{
			Name: "to_d",
			Fn: func(receiver Object, sourceLine int, t *Thread, args []Object, blockFrame *normalCallFrame) Object {
				if len(args) > 0 {
					return t.vm.InitErrorObject(errors.ArgumentError, sourceLine, "Expect 0 arguments. got: %d", len(args))
				}
				r := receiver.(*IntegerObject)
				return t.vm.initDecimalObject(intToDecimal(r))

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
			Fn: func(receiver Object, sourceLine int, t *Thread, args []Object, blockFrame *normalCallFrame) Object {
				r := receiver.(*IntegerObject)
				newFloat := t.vm.initFloatObject(float64(r.value))
				return newFloat

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
			Fn: func(receiver Object, sourceLine int, t *Thread, args []Object, blockFrame *normalCallFrame) Object {
				return receiver

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
			Fn: func(receiver Object, sourceLine int, t *Thread, args []Object, blockFrame *normalCallFrame) Object {

				int := receiver.(*IntegerObject)

				return t.vm.InitStringObject(strconv.Itoa(int.value))

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
			Fn: func(receiver Object, sourceLine int, t *Thread, args []Object, blockFrame *normalCallFrame) Object {
				i := receiver.(*IntegerObject)
				return t.vm.InitIntegerObject(i.value + 1)

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
			Fn: func(receiver Object, sourceLine int, t *Thread, args []Object, blockFrame *normalCallFrame) Object {

				i := receiver.(*IntegerObject)
				odd := i.value%2 != 0
				if odd {
					return TRUE
				}

				return FALSE

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
			Fn: func(receiver Object, sourceLine int, t *Thread, args []Object, blockFrame *normalCallFrame) Object {
				i := receiver.(*IntegerObject)
				return t.vm.InitIntegerObject(i.value - 1)

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
			Fn: func(receiver Object, sourceLine int, t *Thread, args []Object, blockFrame *normalCallFrame) Object {
				n := receiver.(*IntegerObject)

				if n.value < 0 {
					return t.vm.InitErrorObject(errors.InternalError, sourceLine, "Expect integer greater than or equal 0. got: %d", n.value)
				}

				if blockFrame == nil {
					return t.vm.InitErrorObject(errors.InternalError, sourceLine, errors.CantYieldWithoutBlockFormat)
				}

				for i := 0; i < n.value; i++ {
					t.builtinMethodYield(blockFrame, t.vm.InitIntegerObject(i))
				}

				return n

			},
		},
		{
			Name: "to_int",
			Fn: func(receiver Object, sourceLine int, t *Thread, args []Object, blockFrame *normalCallFrame) Object {
				r := receiver.(*IntegerObject)
				newInt := t.vm.InitIntegerObject(r.value)
				newInt.flag = i
				return newInt

			},
		},
		{
			Name: "to_int8",
			Fn: func(receiver Object, sourceLine int, t *Thread, args []Object, blockFrame *normalCallFrame) Object {
				r := receiver.(*IntegerObject)
				newInt := t.vm.InitIntegerObject(r.value)
				newInt.flag = i8
				return newInt

			},
		},
		{
			Name: "to_int16",
			Fn: func(receiver Object, sourceLine int, t *Thread, args []Object, blockFrame *normalCallFrame) Object {
				r := receiver.(*IntegerObject)
				newInt := t.vm.InitIntegerObject(r.value)
				newInt.flag = i16
				return newInt

			},
		},
		{
			Name: "to_int32",
			Fn: func(receiver Object, sourceLine int, t *Thread, args []Object, blockFrame *normalCallFrame) Object {
				r := receiver.(*IntegerObject)
				newInt := t.vm.InitIntegerObject(r.value)
				newInt.flag = i32
				return newInt

			},
		},
		{
			Name: "to_int64",
			Fn: func(receiver Object, sourceLine int, t *Thread, args []Object, blockFrame *normalCallFrame) Object {
				r := receiver.(*IntegerObject)
				newInt := t.vm.InitIntegerObject(r.value)
				newInt.flag = i64
				return newInt

			},
		},
		{
			Name: "to_uint",
			Fn: func(receiver Object, sourceLine int, t *Thread, args []Object, blockFrame *normalCallFrame) Object {
				r := receiver.(*IntegerObject)
				newInt := t.vm.InitIntegerObject(r.value)
				newInt.flag = ui
				return newInt

			},
		},
		{
			Name: "to_uint8",
			Fn: func(receiver Object, sourceLine int, t *Thread, args []Object, blockFrame *normalCallFrame) Object {
				r := receiver.(*IntegerObject)
				newInt := t.vm.InitIntegerObject(r.value)
				newInt.flag = ui8
				return newInt

			},
		},
		{
			Name: "to_uint16",
			Fn: func(receiver Object, sourceLine int, t *Thread, args []Object, blockFrame *normalCallFrame) Object {
				r := receiver.(*IntegerObject)
				newInt := t.vm.InitIntegerObject(r.value)
				newInt.flag = ui16
				return newInt

			},
		},
		{
			Name: "to_uint32",
			Fn: func(receiver Object, sourceLine int, t *Thread, args []Object, blockFrame *normalCallFrame) Object {
				r := receiver.(*IntegerObject)
				newInt := t.vm.InitIntegerObject(r.value)
				newInt.flag = ui32
				return newInt

			},
		},
		{
			Name: "to_uint64",
			Fn: func(receiver Object, sourceLine int, t *Thread, args []Object, blockFrame *normalCallFrame) Object {
				r := receiver.(*IntegerObject)
				newInt := t.vm.InitIntegerObject(r.value)
				newInt.flag = ui64
				return newInt

			},
		},
		{
			Name: "to_float32",
			Fn: func(receiver Object, sourceLine int, t *Thread, args []Object, blockFrame *normalCallFrame) Object {
				r := receiver.(*IntegerObject)
				newInt := t.vm.InitIntegerObject(r.value)
				newInt.flag = f32
				return newInt

			},
		},
		{
			Name: "to_float64",
			Fn: func(receiver Object, sourceLine int, t *Thread, args []Object, blockFrame *normalCallFrame) Object {
				r := receiver.(*IntegerObject)
				newInt := t.vm.InitIntegerObject(r.value)
				newInt.flag = f64
				return newInt

			},
		},
		{
			Name: "ptr",
			Fn: func(receiver Object, sourceLine int, t *Thread, args []Object, blockFrame *normalCallFrame) Object {
				r := receiver.(*IntegerObject)
				return t.vm.initGoObject(&r.value)

			},
		},
	}
}

// Internal functions ===================================================

// Functions for initialization -----------------------------------------

func (vm *VM) InitIntegerObject(value int) *IntegerObject {
	return &IntegerObject{
		BaseObj: &BaseObj{class: vm.TopLevelClass(classes.IntegerClass)},
		value:   value,
		flag:    i,
	}
}

func (vm *VM) initIntegerClass() *RClass {
	ic := vm.initializeClass(classes.IntegerClass)
	ic.setBuiltinMethods(builtinIntegerInstanceMethods(), false)
	ic.setBuiltinMethods(builtinIntegerClassMethods(), true)
	return ic
}

// Polymorphic helper functions -----------------------------------------

// Value returns the object
func (i *IntegerObject) Value() interface{} {
	return i.value
}

// Numeric interface
func (i *IntegerObject) floatValue() float64 {
	return float64(i.value)
}

// TODO: Remove instruction argument
// Apply the passed arithmetic operation, while performing type conversion.
func (i *IntegerObject) arithmeticOperation(
	t *Thread,
	rightObject Object,
	intOperation func(leftValue int, rightValue int) int,
	floatOperation func(leftValue float64, rightValue float64) float64,
	sourceLine int,
	division bool,
) Object {
	switch rightObject.(type) {
	case *IntegerObject:
		leftValue := i.value
		rightValue := rightObject.(*IntegerObject).value
		if division && rightValue == 0 {
			return t.vm.InitErrorObject(errors.ZeroDivisionError, sourceLine, errors.DividedByZero)
		}

		result := intOperation(leftValue, rightValue)

		return t.vm.InitIntegerObject(result)
	case *FloatObject:
		leftValue := float64(i.value)
		rightValue := rightObject.(*FloatObject).value

		if division && rightValue == 0 {
			return t.vm.InitErrorObject(errors.ZeroDivisionError, sourceLine, errors.DividedByZero)
		}

		result := floatOperation(leftValue, rightValue)

		return t.vm.initFloatObject(result)
	default:
		return t.vm.InitErrorObject(errors.TypeError, sourceLine, errors.WrongArgumentTypeFormat, "Numeric", rightObject.Class().Name)
	}
}

// Apply an equality test, returning true if the objects are considered equal,
// and false otherwise.
// See comment on numericComparison().
func (i *IntegerObject) equalityTest(rightObject Object) bool {
	switch rightObject.(type) {
	case *IntegerObject:
		leftValue := i.value
		rightValue := rightObject.(*IntegerObject).value

		return leftValue == rightValue
	case *FloatObject:
		leftValue := i.floatValue()
		rightValue := rightObject.(*FloatObject).value

		return leftValue == rightValue
	default:
		return false
	}
}

// TODO: Remove instruction argument
// Apply the passed numeric comparison, while performing type conversion.
// 64-bit floats cover all the 32-bit integers, but since int is defined
// as *at least* 32 bit, we use two separate functions for safety.
func (i *IntegerObject) numericComparison(
	t *Thread,
	rightObject Object,
	intComparison func(leftValue int, rightValue int) bool,
	floatComparison func(leftValue float64, rightValue float64) bool,
	sourceLine int,
) Object {
	switch rightObject.(type) {
	case *IntegerObject:
		leftValue := i.value
		rightValue := rightObject.(*IntegerObject).value

		result := intComparison(leftValue, rightValue)

		return toBooleanObject(result)
	case *FloatObject:
		leftValue := i.floatValue()
		rightValue := rightObject.(*FloatObject).value

		result := floatComparison(leftValue, rightValue)

		return toBooleanObject(result)
	default:
		return t.vm.InitErrorObject(errors.TypeError, sourceLine, errors.WrongArgumentTypeFormat, "Numeric", rightObject.Class().Name)
	}
}

// ToString returns the object's name as the string format
func (i *IntegerObject) ToString() string {
	return strconv.Itoa(i.value)
}

func (i *IntegerObject) Inspect() string {
	return i.ToString()
}

// ToJSON just delegates to ToString
func (i *IntegerObject) ToJSON(t *Thread) string {
	return i.ToString()
}

// equal checks if the integer values between receiver and argument are equal
func (i *IntegerObject) equal(e *IntegerObject) bool {
	return i.value == e.value
}
