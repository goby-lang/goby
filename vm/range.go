package vm

import (
	"fmt"

	"github.com/goby-lang/goby/vm/classes"
	"github.com/goby-lang/goby/vm/errors"
)

// RangeObject is the built in range class
// Range represents an interval: a set of values from the beginning to the end specified.
// Currently, only Integer objects or integer literal are supported.
//
// ```ruby
// r = 0
// (1..(1+4)).each do |i|
//   puts(r = r + i)
// end
// ```
//
// ```ruby
// r = 0
// a = 1
// b = 5
// (a..b).each do |i|
//   r = r + i
// end
// ```
//
type RangeObject struct {
	*baseObj
	Start int
	End   int
}

// Class methods --------------------------------------------------------
func builtinRangeClassMethods() []*BuiltinMethodObject {
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
func builtinRangeInstanceMethods() []*BuiltinMethodObject {
	return []*BuiltinMethodObject{
		{
			// Returns a Boolean of compared two ranges
			//
			// ```ruby
			// (1..5) == (1..5) # => true
			// (1..5) == (1..6) # => false
			// ```
			//
			// @return [Boolean]
			Name: "==",
			Fn: func(receiver Object, sourceLine int) builtinMethodBody {
				return func(t *thread, args []Object, blockFrame *normalCallFrame) Object {

					left := receiver.(*RangeObject)
					r := args[0]
					right, ok := r.(*RangeObject)

					if !ok {
						return FALSE
					}

					if left.Start == right.Start && left.End == right.End {
						return TRUE
					}

					return FALSE
				}
			},
		},
		{
			// Returns a Boolean of compared two ranges
			//
			// ```ruby
			// (1..5) != (1..5) # => false
			// (1..5) != (1..6) # => true
			// ```
			//
			// @return [Boolean]
			Name: "!=",
			Fn: func(receiver Object, sourceLine int) builtinMethodBody {
				return func(t *thread, args []Object, blockFrame *normalCallFrame) Object {

					left := receiver.(*RangeObject)
					r := args[0]
					right, ok := r.(*RangeObject)

					if !ok {
						return TRUE
					}

					if left.Start == right.Start && left.End == right.End {
						return FALSE
					}

					return TRUE
				}
			},
		},
		{
			// By using binary search, finds a value in range which meets the given condition in O(log n)
			// where n is the size of the range.
			//
			// You can use this method in two use cases: a find-minimum mode and a find-any mode. In either
			// case, the elements of the range must be monotone (or sorted) with respect to the block.
			//
			// In find-minimum mode (this is a good choice for typical use case), the block must return true
			// or false, and there must be a value x so that:
			//
			// - the block returns false for any value which is less than x
			// - the block returns true for any value which is greater than or equal to x.
			//
			// If x is within the range, this method returns the value x. Otherwise, it returns nil.
			//
			// ```ruby
			// ary = [0, 4, 7, 10, 12]
			// (0..4).bsearch {|i| ary[i] >= 4 } #=> 1
			// (0..4).bsearch {|i| ary[i] >= 6 } #=> 2
			// (0..4).bsearch {|i| ary[i] >= 8 } #=> 3
			// (0..4).bsearch {|i| ary[i] >= 100 } #=> nil
			// ```
			//
			// In find-any mode , the block must return a number, and there must be two values x and y
			// (x <= y) so that:
			//
			// - the block returns a positive number for v if v < x
			// - the block returns zero for v if x <= v < y
			// - the block returns a negative number for v if y <= v
			//
			// This method returns any value which is within the intersection of the given range and x…y
			// (if any). If there is no value that satisfies the condition, it returns nil.
			//
			// ```ruby
			// ary = [0, 100, 100, 100, 200]
			// (0..4).bsearch {|i| 100 - ary[i] } #=> 1, 2 or 3
			// (0..4).bsearch {|i| 300 - ary[i] } #=> nil
			// (0..4).bsearch {|i|  50 - ary[i] } #=> nil
			// ```
			//
			// @return [Integer]
			Name: "bsearch",
			Fn: func(receiver Object, sourceLine int) builtinMethodBody {
				return func(t *thread, args []Object, blockFrame *normalCallFrame) Object {
					ran := receiver.(*RangeObject)

					if ran.Start > ran.End || ran.Start < 0 {
						// if block is not used, it should be popped
						t.callFrameStack.pop()
						return NULL
					}

					start := ran.Start
					end := ran.End
					var mid int
					pivot := -1

					for {
						mid = (start + end) / 2
						if (start+end)%2 != 0 {
							mid++
						}

						result := t.builtinMethodYield(blockFrame, t.vm.initIntegerObject(mid))

						switch r := result.Target.(type) {
						case *BooleanObject:
							if r.value {
								pivot = mid
							}

							if start >= end {
								if pivot == -1 {
									return NULL
								}
								return t.vm.initIntegerObject(pivot)
							}

							if r.value {
								end = mid - 1
							} else if mid+1 > ran.End {
								return NULL
							} else {
								start = mid + 1
							}
						case *IntegerObject:
							if r.value == 0 {
								return t.vm.initIntegerObject(mid)
							}

							if start == end {
								return NULL
							}

							if r.value > 0 {
								start = mid + 1
							} else {
								end = mid - 1
							}
						default:
							return t.vm.initErrorObject(errors.TypeError, sourceLine, "Expect Integer or Boolean type. got=%s", r.Class().Name)
						}
					}
				}
			},
		},
		{
			// Iterates over the elements of range, passing each in turn to the block.
			// Returns `nil`.
			//
			// ```ruby
			// sum = 0
			// (1..5).each do |i|
			//   sum = sum + i
			// end
			// sum # => 15
			//
			// sum = 0
			// (-1..-5).each do |i|
			//   sum = sum + i
			// end
			// sum # => -15
			// ```
			//
			// **Note:**
			// - Only `do`-`end` block is supported for now: `{ }` block is unavailable.
			// - Three-dot range `...` is not supported yet.
			//
			// @return [Range]
			Name: "each",
			Fn: func(receiver Object, sourceLine int) builtinMethodBody {
				return func(t *thread, args []Object, blockFrame *normalCallFrame) Object {
					ran := receiver.(*RangeObject)

					if blockFrame == nil {
						return t.vm.initErrorObject(errors.InternalError, sourceLine, errors.CantYieldWithoutBlockFormat)
					}

					if ran.Start <= ran.End {
						for i := ran.Start; i <= ran.End; i++ {
							obj := t.vm.initIntegerObject(i)
							t.builtinMethodYield(blockFrame, obj)
						}
					} else {
						for i := ran.End; i <= ran.Start; i++ {
							obj := t.vm.initIntegerObject(i)
							t.builtinMethodYield(blockFrame, obj)
						}
					}
					return ran
				}
			},
		},
		{
			// Returns the first value of the range.
			//
			// ```ruby
			// (1..5).first   # => 1
			// (5..1).first   # => 5
			// (-2..3).first  # => -2
			// (-5..-7).first # => -5
			// ```
			//
			// @return [Integer]
			Name: "first",
			Fn: func(receiver Object, sourceLine int) builtinMethodBody {
				return func(t *thread, args []Object, blockFrame *normalCallFrame) Object {
					ran := receiver.(*RangeObject)
					return t.vm.initIntegerObject(ran.Start)
				}
			},
		},
		{
			// The include method will check whether the integer object is in the range
			//
			// ```ruby
			// (5..10).include?(10)  # => true
			// (5..10).include?(11)  # => false
			// (5..10).include?(7)   # => true
			// (5..10).include?(5)   # => true
			// (5..10).include?(4)   # => false
			// (-5..1).include?(-2)  # => true
			// (-5..-2).include?(-2) # => true
			// (-5..-3).include?(-2) # => false
			// (1..-5).include?(-2)  # => true
			// (-2..-5).include?(-2) # => true
			// (-3..-5).include?(-2) # => false
			// ```
			// @return [Boolean]
			Name: "include?",
			Fn: func(receiver Object, sourceLine int) builtinMethodBody {
				return func(t *thread, args []Object, blockFrame *normalCallFrame) Object {
					ran := receiver.(*RangeObject)

					value := args[0].(*IntegerObject).value
					ascendRangeBool := ran.Start <= ran.End && value >= ran.Start && value <= ran.End
					descendRangeBool := ran.End <= ran.Start && value <= ran.Start && value >= ran.End

					if ascendRangeBool || descendRangeBool {
						return TRUE
					}
					return FALSE
				}
			},
		},
		{
			// Returns the last value of the range.
			//
			// ```ruby
			// (1..5).last   # => 5
			// (5..1).last   # => 1
			// (-2..3).last  # => 3
			// (-5..-7).last # => -7
			// ```
			//
			// @return [Integer]
			Name: "last",
			Fn: func(receiver Object, sourceLine int) builtinMethodBody {
				return func(t *thread, args []Object, blockFrame *normalCallFrame) Object {
					ran := receiver.(*RangeObject)
					return t.vm.initIntegerObject(ran.End)
				}
			},
		},
		{
			// Returns the size of the range
			//
			// ```ruby
			// (1..5).size   # => 5
			// (3..9).size   # => 7
			// (-1..-5).size # => 5
			// (-1..7).size  # => 9
			// ```
			// @return [Integer]
			Name: "size",
			Fn: func(receiver Object, sourceLine int) builtinMethodBody {
				return func(t *thread, args []Object, blockFrame *normalCallFrame) Object {
					ran := receiver.(*RangeObject)

					if ran.Start <= ran.End {
						return t.vm.initIntegerObject(ran.End - ran.Start + 1)
					}
					return t.vm.initIntegerObject(ran.Start - ran.End + 1)
				}
			},
		},
		{
			// The step method can loop through the first to the last of the object with given steps.
			// An error will occur if not yielded to the block.
			//
			// ```ruby
			// sum = 0
			// (2..9).step(3) do |i|
			// 	 sum = sum + i
			// end
			// sum # => 15
			//
			// sum = 0
			// (2..-9).step(3) do |i|
			// 	 sum = sum + i
			// end
			// sum # => 0
			//
			// sum = 0
			// (-1..5).step(2) do |i|
			//   sum = sum + 1
			// end
			// sum # => 8
			//
			// sum = 0
			// (-1..-5).step(2) do |i|
			//   sum = sum + 1
			// end
			// sum # => 0
			// ```
			//
			// @return [Range]
			Name: "step",
			Fn: func(receiver Object, sourceLine int) builtinMethodBody {
				return func(t *thread, args []Object, blockFrame *normalCallFrame) Object {
					ran := receiver.(*RangeObject)

					if blockFrame == nil {
						return t.vm.initErrorObject(errors.InternalError, sourceLine, errors.CantYieldWithoutBlockFormat)
					}

					stepValue := args[0].(*IntegerObject).value
					if stepValue == 0 {
						return newError("Step can't be 0")
					} else if stepValue < 0 {
						return newError("Step can't be negative")
					}

					// range end must greater or equal than range start to execute the block
					if ran.End >= ran.Start {
						for i := ran.Start; i <= ran.End; i += stepValue {
							obj := t.vm.initIntegerObject(i)
							t.builtinMethodYield(blockFrame, obj)
						}

						return ran
					}

					// if block is not used, it should be popped
					t.callFrameStack.pop()

					return ran
				}
			},
		},
		{
			// Returns an Array object that contains the values of the range.
			//
			// ```ruby
			// (1..5).to_a     # => [1, 2, 3, 4, 5]
			// (1..5).to_a[2]  # => 3
			// (-1..-5).to_a   # => [-1, -2, -3, -4, -5]
			// (-1..3).to_a    # => [-1, 0, 1, 2, 3]
			// ```
			//
			// @return [Array]
			Name: "to_a",
			Fn: func(receiver Object, sourceLine int) builtinMethodBody {
				return func(t *thread, args []Object, blockFrame *normalCallFrame) Object {
					ro := receiver.(*RangeObject)

					elems := []Object{}

					if ro.Start <= ro.End {
						for i := ro.Start; i <= ro.End; i++ {
							elems = append(elems, t.vm.initIntegerObject(i))
						}
					} else {
						for i := ro.End; i <= ro.Start; i++ {
							elems = append(elems, t.vm.initIntegerObject(i))
						}
					}

					return t.vm.initArrayObject(elems)
				}
			},
		},
		{
			// The to_s method can convert range to string format
			//
			// ```ruby
			// (1..5).to_s   # "(1..5)"
			// (-1..-3).to_s # "(-1..-3)"
			// ```
			// @return [String]
			Name: "to_s",
			Fn: func(receiver Object, sourceLine int) builtinMethodBody {
				return func(t *thread, args []Object, blockFrame *normalCallFrame) Object {
					ran := receiver.(*RangeObject)

					return t.vm.initStringObject(ran.toString())
				}
			},
		},
	}
}

// Internal functions ===================================================

// Functions for initialization -----------------------------------------

func (vm *VM) initRangeObject(start, end int) *RangeObject {
	return &RangeObject{
		baseObj: &baseObj{class: vm.topLevelClass(classes.RangeClass)},
		Start:   start,
		End:     end,
	}
}

func (vm *VM) initRangeClass() *RClass {
	rc := vm.initializeClass(classes.RangeClass, false)
	rc.setBuiltinMethods(builtinRangeInstanceMethods(), false)
	rc.setBuiltinMethods(builtinRangeClassMethods(), true)
	return rc
}

// Polymorphic helper functions -----------------------------------------

// toString returns the object's name as the string format
func (ro *RangeObject) toString() string {
	return fmt.Sprintf("(%d..%d)", ro.Start, ro.End)
}

// toJSON just delegates to toString
func (ro *RangeObject) toJSON() string {
	return ro.toString()
}

// Value returns range object's string format
func (ro *RangeObject) Value() interface{} {
	return ro.toString()
}
