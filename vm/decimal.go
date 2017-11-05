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
// DecimalObject represents a comparable decimal number using Go's Rational representation `big.Rat` from math/big package,
// which consists of a numerator and a denominator with arbitrary size.
// By using Decimal you can avoid errors on float type during calculations.
// To keep accuracy, avoid conversions until all calculations have been finished.
// The numerator can be 0, but the denominator cannot be 0.
// Using Decimal for loop counters or like that is not recommended (TBD).
//
// ```ruby
// "3.14".to_d            # => 3.14
// "-0.7238943".to_d      # => -0.7238943
// "355/113".to_d         # => 3.1415929203539823008849557522123893805309734513274336283185840
//
// a = "16.1".to_d
// b = "1.1".to_d
// e = "17.2".to_d
// a + b # => 0.1
// a + b == e # => true
//
// ('16.1'.to_d  + "1.1".to_d).to_s #=> 17.2
// ('16.1'.to_f  + "1.1".to_f).to_s #=> 17.200000000000003
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
			// Returns the sum of self and a numeric.
			// If the second term is integer or float, they are converted into decimal and then perform calculation.
			//
			// ```Ruby
			// "1.1".to_d + "2.1".to_d # => 3.2
			// "1.1".to_d + 2          # => 3.2
			// "1.1".to_d + "2.1".to_f
			// # => 3.200000000000000088817841970012523233890533447265625
			// ```
			//
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
			// If the second term is integer or float, they are converted into decimal and then perform calculation.
			//
			// ```Ruby
			// ("1.5".to_d) - "1.1".to_d   # => 0.4
			// ("1.5".to_d) - 1            # => 0.5
			// ("1.5".to_d) - "1.1".to_f   # => 0.4
			// #=> 0.399999999999999911182158029987476766109466552734375
			// ```
			//
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
			// If the second term is integer or float, they are converted into decimal and then perform calculation.
			//
			// ```Ruby
			// "2.5".to_d * "10.1".to_d     # => 25.25
			// "2.5".to_d * 10              # => 25
			// "2.5".to_d * "10.1".to_f
			// #=> 25.24999999999999911182158029987476766109466552734375
			// ```
			//
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
			// If the second term is integer or float, they are converted into decimal and then perform calculation.
			// Note that the calculation is via float64 (math.Pow) for now.
			//
			// ```Ruby
			// "4.0".to_d ** "2.5".to_d     # => 32
			// "4.0".to_d ** 2              # => 16
			// "4.0".to_d ** "2.5".to_f     # => 32
			// "4.0".to_d ** "2.1".to_d
			// #=> 18.379173679952561570871694129891693592071533203125
			// "4.0".to_d ** "2.1".to_f
			// #=> 18.379173679952561570871694129891693592071533203125
			// ```
			//
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
			// If the second term is integer or float, they are converted into decimal and then perform calculation.
			//
			// ```Ruby
			// "7.5".to_d / "3.1".to_d.fraction      # => 75/31
			// "7.5".to_d / "3.1".to_d
			// # => 2.419354838709677419354838709677419354838709677419354838709677
			// "7.5".to_d / 3                        # => 2.5
			// "7.5".to_d / "3.1".to_f
			// #=> 2.419354838709677350038104601967335570360611893758448172620333
			// ```
			//
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
			// If the second term is integer or float, they are converted into decimal and then perform calculation.
			//
			// ```Ruby
			// a = "3.14".to_d
			// b = "3.16".to_d
			// a > b          # => false
			// b > a          # => true
			// a > 3          # => true
			// a > 4          # => false
			// a > "3.1".to_f # => true
			// a > "3.2".to_f # => false
			// ```
			//
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
			// If the second term is integer or float, they are converted into decimal and then perform calculation.
			//
			// ```Ruby
			// a = "3.14".to_d
			// b = "3.16".to_d
			// e = "3.14".to_d
			// a >= b          # => false
			// b >= a          # => true
			// a >= e          # => true
			// a >= 3          # => true
			// a >= 4          # => false
			// a >= "3.1".to_f # => true
			// a >= "3.2".to_f # => false
			// ```
			//
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
			// If the second term is integer or float, they are converted into decimal and then perform calculation.
			//
			// ```Ruby
			// a = "3.14".to_d
			// b = "3.16".to_d
			// a < b          # => true
			// b < a          # => false
			// a < 3          # => false
			// a < 4          # => true
			// a < "3.1".to_f # => false
			// a < "3.2".to_f # => true
			// ```
			//
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
			// If the second term is integer or float, they are converted into decimal and then perform calculation.
			//
			// ```Ruby
			// a = "3.14".to_d
			// b = "3.16".to_d
			// e = "3.14".to_d
			// a <= b          # => true
			// b <= a          # => false
			// a <= e          # => false
			// a <= 3          # => false
			// a <= 4          # => true
			// a <= "3.1".to_f # => false
			// a <= "3.2".to_f # => true
			// ```
			//
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
			// If the second term is integer or float, they are converted into decimal and then perform calculation.
			// x < y: -1
			// x == y: 0 (including -0 == 0, -Infinity == +Infinity and vice versa)
			// x > y: 1
			//
			// ```Ruby
			// "1.5".to_d <=> 3 # => -1
			// "1.0".to_d <=> 1 # => 0
			// "3.5".to_d <=> 1 # => 1
			// ```
			//
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
			// If the second term is integer or float, they are converted into decimal and then perform calculation.
			// If the Object is not a Numeric the result is always false.
			//
			// ```Ruby
			// "1.0".to_d == 3           # => false
			// "1.0".to_d == 1           # => true
			// "1.0".to_d == "1".to_d    # => true
			// "1.0".to_d == "1".to_f    # => false
			// "1.0".to_d == "1.0".to_f  # => false
			// "1.0".to_d == 'str'       # => false
			// "1.0".to_d == Array       # => false
			// ```
			//
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
			// If the second term is integer or float, they are converted into decimal and then perform calculation.
			// If the Object is not a Numeric the result is always false.
			//
			// ```Ruby
			// "1.0".to_d != 3           # => false
			// "1.0".to_d != 1           # => true
			// "1.0".to_d != "1".to_d    # => true
			// "1.0".to_d != "1".to_f    # => false
			// "1.0".to_d != "1.0".to_f  # => false
			// "1.0".to_d != 'str'       # => false
			// "1.0".to_d != Array       # => false
			// ```
			//
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
			// The value is Decimal.
			// The value does not contain a minus sign.
			//
			// ```Ruby
			// a = "-355/113".to_d
			// a.denominator #=> 113
			// ```
			//
			// @return [Decimal]
			Name: "denominator",
			Fn: func(receiver Object, sourceLine int) builtinMethodBody {
				return func(t *thread, args []Object, blockFrame *normalCallFrame) Object {
					return t.vm.initDecimalObject(new(Decimal).SetInt(receiver.(*DecimalObject).value.Denom()))
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
			// The value is Decimal.
			// The value can contain a minus sign.
			//
			// ```Ruby
			// a = "-355/113".to_d
			// a.numerator #=> -355
			// ```
			//
			// @return [Decimal]
			Name: "numerator",
			Fn: func(receiver Object, sourceLine int) builtinMethodBody {
				return func(t *thread, args []Object, blockFrame *normalCallFrame) Object {
					return t.vm.initDecimalObject(new(Decimal).SetInt(receiver.(*DecimalObject).value.Num()))
				}
			},
		},
		{
			// Returns an array with two Decimal elements: numerator and denominator.
			//
			// ```ruby
			// "129.30928304982039482039842".to_d.to_a
			// # => [6465464152491019741019921, 50000000000000000000000]
			// ```
			//
			// @return [Array]
			Name: "to_a",
			Fn: func(receiver Object, sourceLine int) builtinMethodBody {
				return func(t *thread, args []Object, blockFrame *normalCallFrame) Object {

					n := receiver.(*DecimalObject).value.Num()
					d := receiver.(*DecimalObject).value.Denom()
					elems := []Object{}

					elems = append(elems, t.vm.initDecimalObject(new(Decimal).SetInt(n)))
					elems = append(elems, t.vm.initDecimalObject(new(Decimal).SetInt(d)))

					return t.vm.initArrayObject(elems)
				}
			},
		},
		{
			// Returns an array with two integer elements: numerator and denominator.
			// Big number can be be less accurate than `to_a`.
			//
			// ```ruby
			// "-355.133".to_d.to_ai       # => [-133, 133]
			//
			// "129.30928304982039482039842".to_d.to_ai
			// #=> [-8964879735843078383, -9123183826594430976]
			// ```
			//
			// @return [Array]
			Name: "to_ai",
			Fn: func(receiver Object, sourceLine int) builtinMethodBody {
				return func(t *thread, args []Object, blockFrame *normalCallFrame) Object {

					n := int(receiver.(*DecimalObject).value.Num().Int64())
					d := int(receiver.(*DecimalObject).value.Denom().Int64())
					elems := []Object{}

					elems = append(elems, t.vm.initIntegerObject(n))
					elems = append(elems, t.vm.initIntegerObject(d))

					return t.vm.initArrayObject(elems)
				}
			},
		},
		{
			// Returns Float object from Decimal object.
			// In most case the number of digits in Float is shorter than the one in Decimal.
			//
			// ```Ruby
			// a = "355/113".to_d
			// a.to_s # => 3.1415929203539823008849557522123893805309734513274336283185840
			// a.to_f # => 3.1415929203539825
			// ```
			//
			// @return [Float]
			Name: "to_f",
			Fn: func(receiver Object, sourceLine int) builtinMethodBody {
				return func(t *thread, args []Object, blockFrame *normalCallFrame) Object {
					return t.vm.initFloatObject(receiver.(*DecimalObject).FloatValue())
				}
			},
		},
		{
			// Returns a hash with two keys (numerator and denominator) and their values.
			// The values are Decimal.
			//
			// ```ruby
			// "129.30928304982039482039842".to_d.to_h
			// # => { numerator: 6465464152491019741019921, denominator: 50000000000000000000000 }
			// ```
			//
			// @return [Hash]
			Name: "to_h",
			Fn: func(receiver Object, sourceLine int) builtinMethodBody {
				return func(t *thread, args []Object, blockFrame *normalCallFrame) Object {

					n := receiver.(*DecimalObject).value.Num()
					d := receiver.(*DecimalObject).value.Denom()
					h := make(map[string]Object)

					h["denominator"] = t.vm.initDecimalObject(new(Decimal).SetInt(d))
					h["numerator"] = t.vm.initDecimalObject(new(Decimal).SetInt(n))

					return t.vm.initHashObject(h)
				}
			},
		},
		{
			// Returns a hash with two keys (numerator and denominator) and their values.
			// Returned values are integer, and big number can be be less accurate than `to_h`.
			//
			// ```ruby
			// "129.3".to_d.to_hi   #=> { denominator: 10, numerator: 1293 }
			// "129.30928304982039482039842".to_d.to_hi
			// # => { numerator: -9123183826594430976, denominator: -8964879735843078383 }
			// ```
			//
			// @return [Hash]
			Name: "to_hi",
			Fn: func(receiver Object, sourceLine int) builtinMethodBody {
				return func(t *thread, args []Object, blockFrame *normalCallFrame) Object {

					n := int(receiver.(*DecimalObject).value.Num().Int64())
					d := int(receiver.(*DecimalObject).value.Denom().Int64())
					h := make(map[string]Object)

					h["denominator"] = t.vm.initIntegerObject(d)
					h["numerator"] = t.vm.initIntegerObject(n)

					return t.vm.initHashObject(h)
				}
			},
		},
		{
			// Returns the truncated Integer object from Decimal object.
			//
			// ```Ruby
			// a = "355/113".to_d
			// a.to_i # => 3
			// ```
			//
			// @return [Integer]
			Name: "to_i",
			Fn: func(receiver Object, sourceLine int) builtinMethodBody {
				return func(t *thread, args []Object, blockFrame *normalCallFrame) Object {
					return t.vm.initIntegerObject(receiver.(*DecimalObject).IntegerValue())
				}
			},
		},
		{
			// Returns the float-converted decimal value with a string style.
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

// Alias of Value()
func (f *DecimalObject) DecimalValue() interface{} {
	return f.Value()
}

// Returns integer part of decimal
func (f *DecimalObject) IntegerValue() int {
	return int(f.FloatValue())
}

// Float interface
func (f *DecimalObject) FloatValue() float64 {
	x, _ := f.value.Float64()
	return x
}

// Apply the passed arithmetic operation, while performing type conversion.
func (d *DecimalObject) arithmeticOperation(
	t *thread,
	rightObject Object,
	decimalOperation func(leftValue *Decimal, rightValue *Decimal) *Decimal,
	sourceLine int,
) Object {
	var result Decimal

	rightValue, ok := assertNumeric(rightObject)
	if ok == false {
		return t.vm.initErrorObject(errors.TypeError, sourceLine, errors.WrongArgumentTypeFormat, "Numeric", rightObject.Class().Name)
	}

	leftValue := &d.value
	result = *decimalOperation(leftValue, rightValue)
	return t.vm.initDecimalObject(&result)
}

// Apply an equality test, returning true if the objects are considered equal,
// and false otherwise.
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

// Apply the passed numeric comparison, while performing type conversion.
func (d *DecimalObject) numericComparison(
	t *thread,
	rightObject Object,
	decimalOperation func(leftValue *Decimal, rightValue *Decimal) bool,
	sourceLine int,
) Object {
	var result bool

	rightValue, ok := assertNumeric(rightObject)
	if ok == false {
		return t.vm.initErrorObject(errors.TypeError, sourceLine, errors.WrongArgumentTypeFormat, "Numeric", rightObject.Class().Name)
	}

	leftValue := &d.value
	result = decimalOperation(leftValue, rightValue)
	return toBooleanObject(result)
}

// Apply the passed numeric comparison for rocket operator '<=>', while performing type conversion.
func (d *DecimalObject) rocketComparison(
	t *thread,
	rightObject Object,
	decimalOperation func(leftValue *Decimal, rightValue *Decimal) int,
	sourceLine int,
) Object {
	var result int

	rightValue, ok := assertNumeric(rightObject)
	if ok == false {
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

// Type assertion for numeric
func assertNumeric(n Object) (v *Decimal, result bool) {
	result = true
	switch n.(type) {
	case *DecimalObject:
		v = &n.(*DecimalObject).value
	case *IntegerObject:
		v = intToDecimal(n)
	case *FloatObject:
		v = floatToDecimal(n)
	default:
		result = false
	}
	return v, result
}

// intToDecimal converts int to Decimal
func intToDecimal(i Object) *Decimal {
	return new(Decimal).SetInt64(int64(i.(*IntegerObject).value))
}

// floatToDecimal converts int to Decimal
func floatToDecimal(i Object) *Decimal {
	return new(Decimal).SetFloat64(float64(i.(*FloatObject).value))
}
