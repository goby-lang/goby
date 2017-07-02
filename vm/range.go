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
		// @return [nil]
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
}
