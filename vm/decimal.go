package vm

import (
	"math/big"

	"github.com/goby-lang/goby/vm/classes"
	"github.com/goby-lang/goby/vm/errors"
	"math"
	"strings"
)

// A type alias for representing a decimal
type Decimal = big.Rat
type Int = big.Int
type Float = big.Float

// (Experiment)
// DecimalObject represents a comparable decimal number using Go's Rat from math/big
// representation, which consists of a numerator and a denominator with arbitrary size.
// The numerator can be 0, but the denominator cannot be 0.
//
// ```ruby
// "3.14".to_d            # => 3.14
// "-0.7238943".to_d      # => -0.7238943
// "355/113".to_d         # => 3.1415929203539823008849557522123893805309734513274336283185840
//
// a = "1.1".to_d
// b = "1.0".to_d
// c = "0.1".to_d
// a - b # => 0.1
// a - b == c # => true
// ```
//
// - `Decimal.new` is not supported.
type DecimalObject struct {
	*baseObj
	value Decimal
}

// Class methods --------------------------------------------------------
func builtinDecimalClassMethods() []*BuiltinMethodObject {
	return []*BuiltinMethodObject{
		{
			Name: "new",
			Fn: func(receiver Object, sourceLine int) builtinMethodBody {
				return func(t *thread, args []Object, blockFrame *normalCallFrame) Object {
					return t.initUnsupportedMethodError(sourceLine, "#new", receiver)
				}
			},
		},
	}
}

// Instance methods -----------------------------------------------------
func builtinDecimalInstanceMethods() []*BuiltinMethodObject {
	return []*BuiltinMethodObject{
		{
			// Returns the sum of self and a decimal.
			//
			// ```Ruby
			// 1.1 + 2 # => 3.1
			// ```
			// @return [Decimal]
			Name: "+",
			Fn: func(receiver Object, sourceLine int) builtinMethodBody {
				return func(t *thread, args []Object, blockFrame *normalCallFrame) Object {
					operation := func(leftValue *Decimal, rightValue *Decimal) *Decimal {
						return new(Decimal).Add(leftValue, rightValue)
					}

					return receiver.(*DecimalObject).arithmeticOperation(t, args[0], operation, sourceLine)
				}
			},
		},
		{
			// Returns the subtraction of a decimal from self.
			//
			// ```Ruby
			// 1.5 - 1 # => 0.5
			// ```
			// @return [Decimal]
			Name: "-",
			Fn: func(receiver Object, sourceLine int) builtinMethodBody {
				return func(t *thread, args []Object, blockFrame *normalCallFrame) Object {
					operation := func(leftValue *Decimal, rightValue *Decimal) *Decimal {
						return new(Decimal).Sub(leftValue, rightValue)
					}

					return receiver.(*DecimalObject).arithmeticOperation(t, args[0], operation, sourceLine)
				}
			},
		},
		{
			// Returns self multiplying a decimal.
			//
			// ```Ruby
			// 2.5 * 10 # => 25.0
			// ```
			// @return [Decimal]
			Name: "*",
			Fn: func(receiver Object, sourceLine int) builtinMethodBody {
				return func(t *thread, args []Object, blockFrame *normalCallFrame) Object {
					operation := func(leftValue *Decimal, rightValue *Decimal) *Decimal {
						return new(Decimal).Mul(leftValue, rightValue)
					}

					return receiver.(*DecimalObject).arithmeticOperation(t, args[0], operation, sourceLine)
				}
			},
		},
		{
			// Returns self squaring a decimal.
			//
			// ```Ruby
			// 4.0 ** 2.5 # => 32.0
			// ```
			// @return [Decimal]
			Name: "**",
			Fn: func(receiver Object, sourceLine int) builtinMethodBody {
				return func(t *thread, args []Object, blockFrame *normalCallFrame) Object {
					operation := func(leftValue *Decimal, rightValue *Decimal) *Decimal {
						l, _ := leftValue.Float64()
						r, _ := rightValue.Float64()
						return new(Decimal).SetFloat64(math.Pow(l, r))
					}

					return receiver.(*DecimalObject).arithmeticOperation(t, args[0], operation, sourceLine)
				}
			},
		},
		{
			// Returns self divided by a decimal.
			//
			// ```Ruby
			// 7.5 / 3 # => 2.5
			// ```
			// @return [Decimal]
			Name: "/",
			Fn: func(receiver Object, sourceLine int) builtinMethodBody {
				return func(t *thread, args []Object, blockFrame *normalCallFrame) Object {
					decimalOperation := func(leftValue *Decimal, rightValue *Decimal) *Decimal {
						return new(Decimal).Quo(leftValue, rightValue)
					}

					return receiver.(*DecimalObject).arithmeticOperation(t, args[0], decimalOperation, sourceLine)
				}
			},
		},
		{
			// Returns if self is larger than a decimal.
			//
			// ```Ruby
			// a = "3.14".to_d
			// b = "3.16".to_d
			// a > b # => false
			// b > a # => true
			// ```
			// @return [Boolean]
			Name: ">",
			Fn: func(receiver Object, sourceLine int) builtinMethodBody {
				return func(t *thread, args []Object, blockFrame *normalCallFrame) Object {
					decimalOperation := func(leftValue *Decimal, rightValue *Decimal) bool {
						if leftValue.Cmp(rightValue) == 1 {
							return true
						} else {
							return false
						}
					}

					return receiver.(*DecimalObject).numericComparison(t, args[0], decimalOperation, sourceLine)
				}
			},
		},
		{
			// Returns if self is larger than or equals to a Numeric.
			//
			// ```Ruby
			// a = "3.14".to_d
			// b = "3.16".to_d
			// e = "3.14".to_d
			// a >= b # => false
			// b >= a # => true
			// a >= e # => true
			// ```
			// @return [Boolean]
			Name: ">=",
			Fn: func(receiver Object, sourceLine int) builtinMethodBody {
				return func(t *thread, args []Object, blockFrame *normalCallFrame) Object {
					decimalOperation := func(leftValue *Decimal, rightValue *Decimal) bool {
						switch leftValue.Cmp(rightValue) {
						case 1, 0:
							return true
						default:
							return false
						}
					}

					return receiver.(*DecimalObject).numericComparison(t, args[0], decimalOperation, sourceLine)
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
					decimalOperation := func(leftValue *Decimal, rightValue *Decimal) bool {
						if leftValue.Cmp(rightValue) == -1 {
							return true
						} else {
							return false
						}
					}

					return receiver.(*DecimalObject).numericComparison(t, args[0], decimalOperation, sourceLine)
				}
			},
		},
		{
			// Returns if self is smaller than or equals to a decimal.
			//
			// ```Ruby
			// 1 <= 3 # => true
			// 1 <= 1 # => true
			// ```
			// @return [Boolean]
			Name: "<=",
			Fn: func(receiver Object, sourceLine int) builtinMethodBody {
				return func(t *thread, args []Object, blockFrame *normalCallFrame) Object {
					decimalOperation := func(leftValue *Decimal, rightValue *Decimal) bool {
						switch leftValue.Cmp(rightValue) {
						case -1, 0:
							return true
						default:
							return false
						}
					}

					return receiver.(*DecimalObject).numericComparison(t, args[0], decimalOperation, sourceLine)
				}
			},
		},
		{
			// Returns 1 if self is larger than a Numeric, -1 if smaller. Otherwise 0.
			// returns -1 if x < y
			// returns 0 if x == 0 (including -0 == 0, -Infinity == +Infinity and vice versa
			// returns 1 if x > 0
			//
			// ```Ruby
			// 1.5 <=> 3 # => -1
			// 1.0 <=> 1 # => 0
			// 3.5 <=> 1 # => 1
			// ```
			// @return [Integer]
			Name: "<=>",
			Fn: func(receiver Object, sourceLine int) builtinMethodBody {
				return func(t *thread, args []Object, blockFrame *normalCallFrame) Object {
					decimalOperation := func(leftValue *Decimal, rightValue *Decimal) int {
						return leftValue.Cmp(rightValue)
					}

					return receiver.(*DecimalObject).rocketComparison(t, args[0], decimalOperation, sourceLine)
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
					decimalOperation := func(leftValue *Decimal, rightValue *Decimal) bool {
						if leftValue.Cmp(rightValue) == 0 {
							return true
						} else {
							return false
						}
					}

					return receiver.(*DecimalObject).equalityTest(t, args[0], decimalOperation, true, sourceLine)
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
					decimalOperation := func(leftValue *Decimal, rightValue *Decimal) bool {
						if leftValue.Cmp(rightValue) != 0 {
							return true
						} else {
							return false
						}
					}

					return receiver.(*DecimalObject).equalityTest(t, args[0], decimalOperation, false, sourceLine)
				}
			},
		},
		{
			// Returns a string with fraction format of the decimal.
			// If the denominator is 1, '/1` is omitted.
			// Minus sign will be preserved.
			// (Actually, the internal rational number is always deducted)
			//
			// ```Ruby
			// a = "-355/113".to_d
			// a.reduction #=> -355/113
			// b = "-331/1".to_d
			// b.reduction #=> -331
			// ```
			//
			// @return [String]
			Name: "reduction",
			Fn: func(receiver Object, sourceLine int) builtinMethodBody {
				return func(t *thread, args []Object, blockFrame *normalCallFrame) Object {
					return t.vm.initStringObject(receiver.(*DecimalObject).value.RatString())
				}
			},
		},
		{
			// Returns the denominator of the decimal value which contains Go's big.Rat type.
			// The value does not contain a minus sign.
			//
			// ```Ruby
			// a = "-355/113".to_d
			// a.denominator #=> 113
			// ```
			//
			// @return [int]
			Name: "denominator",
			Fn: func(receiver Object, sourceLine int) builtinMethodBody {
				return func(t *thread, args []Object, blockFrame *normalCallFrame) Object {
					return t.vm.initIntegerObject(int(receiver.(*DecimalObject).value.Denom().Int64()))
				}
			},
		},
		{
			// Returns a string with fraction format of the decimal.
			// Even though the denominator is 1, fraction format is used.
			// Minus sign will be preserved.
			//
			// ```Ruby
			// a = "-355/113".to_d
			// a.fraction #=> -355/113
			// b = "-331/1".to_d
			// b.fraction #=> -331/1
			// ```
			//
			// @return [String]
			Name: "fraction",
			Fn: func(receiver Object, sourceLine int) builtinMethodBody {
				return func(t *thread, args []Object, blockFrame *normalCallFrame) Object {
					return t.vm.initStringObject(receiver.(*DecimalObject).value.String())
				}
			},
		},
		{
			// Inverses the numerator and the denominator of the decimal and returns it.
			// Minus sign will move to the new numerator.
			//
			// ```Ruby
			// a = "-355/113".to_d
			// a.inverse.fraction #=> -113/355
			// ```
			//
			// @return [Decimal]
			Name: "inverse",
			Fn: func(receiver Object, sourceLine int) builtinMethodBody {
				return func(t *thread, args []Object, blockFrame *normalCallFrame) Object {
					d := receiver.(*DecimalObject).value
					return t.vm.initDecimalObject(d.Inv(&d))
				}
			},
		},
		{
			// Returns the numerator of the decimal value which contains Go's big.Rat type.
			// The value can contain a minus sign.
			//
			// ```Ruby
			// a = "-355/113".to_d
			// a.numerator #=> -355
			// ```
			//
			// @return [int]
			Name: "numerator",
			Fn: func(receiver Object, sourceLine int) builtinMethodBody {
				return func(t *thread, args []Object, blockFrame *normalCallFrame) Object {
					return t.vm.initIntegerObject(int(receiver.(*DecimalObject).value.Num().Int64()))
				}
			},
		},
		{
			// Returns the decimal value with a string style.
			// Maximum digit under the dots is 60.
			// This is just to print the final value should not be used for recalculation.
			//
			// ```Ruby
			// a = "355/113".to_d
			// a.to_s # => 3.1415929203539823008849557522123893805309734513274336283185840
			// ```
			//
			// @return [String]
			Name: "to_s",
			Fn: func(receiver Object, sourceLine int) builtinMethodBody {
				return func(t *thread, args []Object, blockFrame *normalCallFrame) Object {
					return t.vm.initStringObject(receiver.(*DecimalObject).toString())
				}
			},
		},
	}
}

// Internal functions ===================================================

// Functions for initialization -----------------------------------------

func (vm *VM) initDecimalObject(value *Decimal) *DecimalObject {
	return &DecimalObject{
		baseObj: &baseObj{class: vm.topLevelClass(classes.DecimalClass)},
		value:   *value,
	}
}

func (vm *VM) initDecimalClass() *RClass {
	dc := vm.initializeClass(classes.DecimalClass, false)
	dc.setBuiltinMethods(builtinDecimalInstanceMethods(), false)
	dc.setBuiltinMethods(builtinDecimalClassMethods(), true)
	return dc
}

// Polymorphic helper functions -----------------------------------------

// Value returns the object
func (f *DecimalObject) Value() interface{} {
	return f.value
}

//// Returns integer part of decimal
//func (f *DecimalObject) IntegerValue() interface{} {
//	return int(f.value)
//}
//
//// Float interface
//func (f *DecimalObject) FloatValue() float64 {
//	x, _ := f.value.Float64()
//	return x
//}

// TODO: Remove instruction argument
// Apply the passed arithmetic operation, while performing type conversion.
func (d *DecimalObject) arithmeticOperation(
	t *thread,
	rightObject Object,
	decimalOperation func(leftValue *Decimal, rightValue *Decimal) *Decimal,
	sourceLine int,
) Object {
	var rightValue *Decimal
	var result Decimal

	switch rightObject.(type) {
	case *DecimalObject:
		rightValue = &rightObject.(*DecimalObject).value
	case *IntegerObject:
		rightValue = intToDecimal(rightObject)
	case *FloatObject:
		rightValue = floatToDecimal(rightObject)
	default:
		return t.vm.initErrorObject(errors.TypeError, sourceLine, errors.WrongArgumentTypeFormat, "Numeric", rightObject.Class().Name)
	}

	leftValue := &d.value
	result = *decimalOperation(leftValue, rightValue)
	return t.vm.initDecimalObject(&result)
}

// Apply an equality test, returning true if the objects are considered equal,
// and false otherwise.
// TODO: Remove instruction argument
func (d *DecimalObject) equalityTest(
	t *thread,
	rightObject Object,
	decimalOperation func(leftValue *Decimal, rightValue *Decimal) bool,
	nonInverse bool,
	sourceLine int,
) Object {
	var rightValue *Decimal
	var result bool

	switch rightObject.(type) {
	case *DecimalObject:
		rightValue = &rightObject.(*DecimalObject).value
	default:
		return toBooleanObject(nonInverse == false)
	}

	leftValue := &d.value
	result = decimalOperation(leftValue, rightValue)
	return toBooleanObject(result)
}

// TODO: Remove instruction argument
// Apply the passed numeric comparison, while performing type conversion.
func (d *DecimalObject) numericComparison(
	t *thread,
	rightObject Object,
	decimalOperation func(leftValue *Decimal, rightValue *Decimal) bool,
	sourceLine int,
) Object {
	var rightValue *Decimal
	var result bool

	switch rightObject.(type) {
	case *DecimalObject:
		rightValue = &rightObject.(*DecimalObject).value
	case *IntegerObject:
		rightValue = intToDecimal(rightObject)
	case *FloatObject:
		rightValue = floatToDecimal(rightObject)
	default:
		return t.vm.initErrorObject(errors.TypeError, sourceLine, errors.WrongArgumentTypeFormat, "Numeric", rightObject.Class().Name)
	}

	leftValue := &d.value
	result = decimalOperation(leftValue, rightValue)
	return toBooleanObject(result)
}

// TODO: Remove instruction argument
// Apply the passed numeric comparison for rocket operator '<=>', while performing type conversion.
func (d *DecimalObject) rocketComparison(
	t *thread,
	rightObject Object,
	decimalOperation func(leftValue *Decimal, rightValue *Decimal) int,
	sourceLine int,
) Object {
	var rightValue *Decimal
	var result int

	switch rightObject.(type) {
	case *DecimalObject:
		rightValue = &rightObject.(*DecimalObject).value
	case *IntegerObject:
		rightValue = intToDecimal(rightObject)
	case *FloatObject:
		rightValue = floatToDecimal(rightObject)
	default:
		return t.vm.initErrorObject(errors.TypeError, sourceLine, errors.WrongArgumentTypeFormat, "Numeric", rightObject.Class().Name)
	}

	leftValue := &d.value
	result = decimalOperation(leftValue, rightValue)
	newInt := t.vm.initIntegerObject(result)
	newInt.flag = i
	return newInt
}

// toString returns the object's approximate float value as the string format.
func (d *DecimalObject) toString() string {
	fs := d.value.FloatString(60)
	fs = strings.TrimRight(fs, "0")
	fs = strings.TrimRight(fs, ".")
	return fs
}

// toJSON just delegates to toString
func (d *DecimalObject) toJSON() string {
	return d.toString()
}

// Other helper functions  ----------------------------------------------

// intToDecimal converts int to Decimal
func intToDecimal(i Object) *Decimal {
	return new(Decimal).SetInt64(int64(i.(*IntegerObject).value))
}

// floatToDecimal converts int to Decimal
func floatToDecimal(i Object) *Decimal {
	return new(Decimal).SetFloat64(float64(i.(*FloatObject).value))
}
