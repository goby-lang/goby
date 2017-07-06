package vm

import (
	"fmt"
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
	Class *RClass
	Start int
	End   int
}

func (ro *RangeObject) toString() string {
	return fmt.Sprintf("(%d..%d)", ro.Start, ro.End)
}

func (ro *RangeObject) toJSON() string {
	return ro.toString()
}

func (ro *RangeObject) returnClass() Class {
	return ro.Class
}

func (vm *VM) initRangeObject(start, end int) *RangeObject {
	return &RangeObject{Class: vm.builtInClasses[rangeClass], Start: start, End: end}
}

func (vm *VM) initRangeClass() *RClass {
	rc := vm.initializeClass(rangeClass, false)
	rc.setBuiltInMethods(builtInRangeInstanceMethods(), false)
	rc.setBuiltInMethods(builtInRangeClassMethods(), true)
	return rc
}

func builtInRangeClassMethods() []*BuiltInMethodObject {
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

func builtInRangeInstanceMethods() []*BuiltInMethodObject {
	return []*BuiltInMethodObject{
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
			Fn: func(receiver Object) builtinMethodBody {
				return func(t *thread, args []Object, blockFrame *callFrame) Object {
					ran := receiver.(*RangeObject)
					return t.vm.initIntegerObject(ran.Start)
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
			Fn: func(receiver Object) builtinMethodBody {
				return func(t *thread, args []Object, blockFrame *callFrame) Object {
					ran := receiver.(*RangeObject)
					return t.vm.initIntegerObject(ran.End)
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
			Fn: func(receiver Object) builtinMethodBody {
				return func(t *thread, args []Object, blockFrame *callFrame) Object {
					ran := receiver.(*RangeObject)

					if blockFrame == nil {
						t.returnError("Can't yield without a block")
					}

					if ran.Start <= ran.End {
						for i := ran.Start; i <= ran.End; i++ {
							obj := t.vm.initIntegerObject(i)
							t.builtInMethodYield(blockFrame, obj)
						}
					} else {
						for i := ran.End; i <= ran.Start; i++ {
							obj := t.vm.initIntegerObject(i)
							t.builtInMethodYield(blockFrame, obj)
						}
					}
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
			Fn: func(receiver Object) builtinMethodBody {
				return func(t *thread, args []Object, blockFrame *callFrame) Object {
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
			Fn: func(receiver Object) builtinMethodBody {
				return func(t *thread, args []Object, blockFrame *callFrame) Object {
					ran := receiver.(*RangeObject)

					return t.vm.initStringObject(ran.toString())
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
			Fn: func(receiver Object) builtinMethodBody {
				return func(t *thread, args []Object, blockFrame *callFrame) Object {
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
			Fn: func(receiver Object) builtinMethodBody {
				return func(t *thread, args []Object, blockFrame *callFrame) Object {
					ran := receiver.(*RangeObject)

					if blockFrame == nil {
						t.returnError("Can't yield without a block")
					}

					stepValue := args[0].(*IntegerObject).Value
					if stepValue == 0 {
						return newError("Step can't be 0")
					} else if stepValue < 0 {
						return newError("Step can't be negative")
					}

					for i := ran.Start; i <= ran.End; i += stepValue {
						obj := t.vm.initIntegerObject(i)
						t.builtInMethodYield(blockFrame, obj)
					}
					return ran
				}
			},
		},
		{
			// The include method will check whether the integer object is in the range
			//
			// ```ruby
			// (5..10).include(10)  # => true
			// (5..10).include(11)  # => false
			// (5..10).include(7)   # => true
			// (5..10).include(5)   # => true
			// (5..10).include(4)   # => false
			// (-5..1).include(-2)  # => true
			// (-5..-2).include(-2) # => true
			// (-5..-3).include(-2) # => false
			// (1..-5).include(-2)  # => true
			// (-2..-5).include(-2) # => true
			// (-3..-5).include(-2) # => false
			// ```
			// @return [Boolean]
			Name: "include",
			Fn: func(receiver Object) builtinMethodBody {
				return func(t *thread, args []Object, blockFrame *callFrame) Object {
					ran := receiver.(*RangeObject)

					value := args[0].(*IntegerObject).Value
					ascendRangeBool := ran.Start <= ran.End && value >= ran.Start && value <= ran.End
					descendRangeBool := ran.End <= ran.Start && value <= ran.Start && value >= ran.End

					if ascendRangeBool || descendRangeBool {
						return TRUE
					}
					return FALSE
				}
			},
		},
	}
}
