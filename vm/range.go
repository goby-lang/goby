package vm

import (
	"fmt"
)

var (
	rangeClass *RClass
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

func (ro *RangeObject) toArray() *ArrayObject {
	elems := []Object{}

	for i := ro.Start; i <= ro.End; i++ {
		elems = append(elems, initIntegerObject(i))
	}

	return initArrayObject(elems)
}

func initRangeObject(start, end int) *RangeObject {
	return &RangeObject{Class: rangeClass, Start: start, End: end}
}

func initRangeClass() {
	bc := &BaseClass{Name: "Range", ClassMethods: newEnvironment(), Methods: newEnvironment(), Class: classClass, pseudoSuperClass: objectClass, superClass: objectClass}
	rc := &RClass{BaseClass: bc}
	rc.setBuiltInMethods(builtInRangeInstanceMethods, false)
	rangeClass = rc
}

var builtInRangeInstanceMethods = []*BuiltInMethodObject{
	{
		// Returns the first value of the range.
		//
		// ```ruby
		// (1..5).first # => 1
		// ```
		//
		// @return [Integer]
		Name: "first",
		Fn: func(receiver Object) builtinMethodBody {
			return func(t *thread, args []Object, blockFrame *callFrame) Object {
				ran := receiver.(*RangeObject)
				return initIntegerObject(ran.Start)
			}
		},
	},
	{
		// Returns the last value of the range.
		//
		// ```ruby
		// (1..5).last # => 5
		// ```
		//
		// @return [Integer]
		Name: "last",
		Fn: func(receiver Object) builtinMethodBody {
			return func(t *thread, args []Object, blockFrame *callFrame) Object {
				ran := receiver.(*RangeObject)
				return initIntegerObject(ran.End)
			}
		},
	},
	{
		// Iterates over the elements of range, passing each in turn to the block.
		// Returns `nil`.
		//
		// ```ruby
		// (1..5).to_a     # => [1, 2, 3, 4, 5]
		// (1..5).to_a[2]  # => 3
		// ```
		//
		// **Note:**
		// - Only `do`-`end` block is supported for now: `{ }` block is unavailable.
		// - Three-dot range `...` is not supported yet.
		//
		// @return [Null]
		Name: "each",
		Fn: func(receiver Object) builtinMethodBody {
			return func(t *thread, args []Object, blockFrame *callFrame) Object {
				ran := receiver.(*RangeObject)

				if blockFrame == nil {
					t.returnError("Can't yield without a block")
				}

				for i := ran.Start; i <= ran.End; i++ {
					obj := initIntegerObject(i)
					t.builtInMethodYield(blockFrame, obj)
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
		// ```
		//
		// @return [Array]
		Name: "to_a",
		Fn: func(receiver Object) builtinMethodBody {
			return func(t *thread, args []Object, blockFrame *callFrame) Object {
				ran := receiver.(*RangeObject)

				return ran.toArray()
			}
		},
	},
	{
		// The to_s method can convert range to string format
		//
		// ```ruby
		// (1..5).to_s # "(1..5)"
		// ```
		// @return [String]
		Name: "to_s",
		Fn: func(receiver Object) builtinMethodBody {
			return func(t *thread, args []Object, blockFrame *callFrame) Object {
				ran := receiver.(*RangeObject)

				return initStringObject(ran.toString())
			}
		},
	},
	{
		// Returns the size of the range
		//
		// ```ruby
		// (1..5).size # => 5
		// (3..9).size # => 7
		// ```
		// @return [Integer]
		Name: "size",
		Fn: func(receiver Object) builtinMethodBody {
			return func(t *thread, args []Object, blockFrame *callFrame) Object {
				ran := receiver.(*RangeObject)

				return initIntegerObject(ran.End - ran.Start + 1)
			}
		},
	},
	{
		// The step method can loop through the first to the last of the object with given steps.
		// An error will occur if not yielded to the block.
		//
		// ```ruby
		// (2..9).step(3) do |i|
		// 	 puts i
		// end
		// # => 2
		// # => 5
		// # => 8
		// ```
		Name: "step",
		Fn: func(receiver Object) builtinMethodBody {
			return func(t *thread, args []Object, blockFrame *callFrame) Object {
				ran := receiver.(*RangeObject)

				if blockFrame == nil {
					t.returnError("Can't yield without a block")
				}

				stepValue := args[0].(*IntegerObject).Value
				for i := ran.Start; i <= ran.End; i += stepValue {
					obj := initIntegerObject(i)
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
		// (5..10).include(7) # => true
		// (5..10).include(5) # => true
		// (5..10).include(4) # => false
		// ```
		// @return [Boolean]
		Name: "include",
		Fn: func(receiver Object) builtinMethodBody {
			return func(t *thread, args []Object, blockFrame *callFrame) Object {
				ran := receiver.(*RangeObject)

				value := args[0].(*IntegerObject).Value
				start := ran.Start
				end := ran.End
				if value >= start || value <= end {
					return TRUE
				}
				return FALSE
			}
		},
	},
	{
	// By using binary search, finds a value in range which meets the given condition in O(log n)
	// where n is the size of the range.
	//
	// You can use this method to find minimum number. The elements of the range must be monotone
	// (or sorted) with respect to the block.
	//
	// In find minimum mode (this is a good choice for typical use case), the block must return
	// true or false, and there must be a value x so that:
	//
	// - the block returns false for any value which is less than x, and
	//
	// - the block returns true for any value which is greater than or equal to x.
	//
	// If x is within the range, this method returns the value x. Otherwise, it returns nil.
	//
	// ```ruby
	// ary = [0, 4, 7, 10, 12]
	// (0..(ary.length)).bsearch do
	//   |i| ary[i] >= 4
	// end
	// #=> 1
	//
	// (0..(ary.length)).bsearch do |i|
	//   ary[i] >= 6
	// end
	// #=> 2
	//
	// (0..(ary.length)).bsearch do |i|
	//   ary[i] >= 8
	// end
	// #=> 3
	//
	// (0..(ary.length)).bsearch do |i|
	//   ary[i] >= 100
	// end
	// #=> nil
	// ```
	// @return [Integer]
	//Name: "bsearch",
	//Fn: func(receiver Object) builtinMethodBody {
	//	return func(t *thread, args []Object, blockFrame *callFrame) Object {
	//		ran := receiver.(*RangeObject)
	//
	//		if blockFrame == nil {
	//			t.returnError("Can't yield without a block")
	//		}
	//
	//		start := ran.Start
	//		end := ran.End
	//		if start > end {
	//			return NULL
	//		} else {
	//			for start <= end {
	//				mid := (start + end) / 2
	//				if ((start + end) % 2 != 0) {
	//					mid += 1
	//				}
	//				fmt.Println(mid)
	//				obj := initIntegerObject(mid)
	//				result := t.builtInMethodYield(blockFrame, obj)
	//				if result.Target.(*BooleanObject).Value {
	//					start = mid + 1
	//				} else {
	//					end = mid - 1
	//				}
	//
	//				if start == end {
	//					return initIntegerObject(start)
	//				}
	//			}
	//			return NULL
	//		}
	//	}
	//},
	},
}
