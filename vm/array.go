package vm

import (
	"bytes"
	"strings"

	"github.com/goby-lang/goby/vm/classes"
	"github.com/goby-lang/goby/vm/errors"
)

// ArrayObject represents instance from Array class.
// An array is a collection of different objects that are ordered and indexed.
// Elements in an array can belong to any class.
type ArrayObject struct {
	*baseObj
	Elements []Object
	splat    bool
}

// Class methods --------------------------------------------------------
func builtinArrayClassMethods() []*BuiltinMethodObject {
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
func builtinArrayInstanceMethods() []*BuiltinMethodObject {
	return []*BuiltinMethodObject{
		{
			// Retrieves an object in an array using Integer index.
			// The index starts from 0. It returns `null` if the given index is bigger than its size.
			//
			// ```ruby
			// a = [1, 2, 3, "a", "b", "c"]
			// a[0]  # => 1
			// a[3]  # => "a"
			// a[10] # => nil
			// a[-1] # => "c"
			// a[-3] # => "a"
			// a[-7] # => nil
			// ```
			Name: "[]",
			Fn: func(receiver Object, sourceLine int) builtinMethodBody {
				return func(t *thread, args []Object, blockFrame *normalCallFrame) Object {
					arr := receiver.(*ArrayObject)
					return arr.index(t, args, sourceLine)
				}
			},
		},
		{
			// Repetition — returns a new array built by concatenating the specified number of copies
			// of `self`.
			//
			// ```ruby
			// a = [1, 2, 3]
			// a * 2   # => [1, 2, 3, 1, 2, 3]
			// a * ',' # => "1,2,3"
			// ```
			Name: "*",
			Fn: func(receiver Object, sourceLine int) builtinMethodBody {
				return func(t *thread, args []Object, blockFrame *normalCallFrame) Object {
					if len(args) != 1 {
						return t.vm.initErrorObject(errors.ArgumentError, sourceLine, "Expect 1 arguments. got=%d", len(args))
					}

					arr := receiver.(*ArrayObject)

					copiesNumber, ok := args[0].(*IntegerObject)

					if !ok {
						return t.vm.initErrorObject(errors.TypeError, sourceLine, errors.WrongArgumentTypeFormat, classes.IntegerClass, args[0].Class().Name)
					}

					return arr.concatenateCopies(t, copiesNumber)
				}
			},
		},
		{
			// Returns a new array built by concatenating the two arrays together to produce a third array.
			//
			// ```ruby
			// a = [1, 2]
			// b + [3, 4]  # => [1, 2, 4, 4]
			// ```
			Name: "+",
			Fn: func(receiver Object, sourceLine int) builtinMethodBody {
				return func(t *thread, args []Object, blockFrame *normalCallFrame) Object {
					if len(args) != 1 {
						return t.vm.initErrorObject(errors.ArgumentError, sourceLine, "Expect 1 arguments. got=%d", len(args))
					}

					otherArrayArg := args[0]
					otherArray, ok := otherArrayArg.(*ArrayObject)

					if !ok {
						return t.vm.initErrorObject(errors.TypeError, sourceLine, errors.WrongArgumentTypeFormat, classes.ArrayClass, args[0].Class().Name)
					}

					selfArray := receiver.(*ArrayObject)

					newArrayelements := append(selfArray.Elements, otherArray.Elements...)

					newArray := t.vm.initArrayObject(newArrayelements)

					return newArray
				}
			},
		},
		{
			// Assigns value to an array. It requires an index and a value as argument.
			// The array will expand if the assigned index is bigger than its size.
			// Returns the assigned value.
			//
			// ```ruby
			// a = []
			// a[0] = 10  # => 10
			// a[3] = 20  # => 20
			// a          # => [10, nil, nil, 20]
			// a[-2] = 5  # => [10, nil, 5, 20]
			// ```
			Name: "[]=",
			Fn: func(receiver Object, sourceLine int) builtinMethodBody {
				return func(t *thread, args []Object, blockFrame *normalCallFrame) Object {

					// First argument is index, there exists two cases which will be described in the following code
					if len(args) != 2 && len(args) != 3 {
						return t.vm.initErrorObject(errors.ArgumentError, sourceLine, "Expect 2..3 arguments. got=%d", len(args))
					}

					i := args[0]
					index, ok := i.(*IntegerObject)

					if !ok {
						return t.vm.initErrorObject(errors.TypeError, sourceLine, errors.WrongArgumentTypeFormat, classes.IntegerClass, args[0].Class().Name)
					}

					indexValue := index.value
					arr := receiver.(*ArrayObject)

					// <Three Argument Case>
					// Second argument is the length of successive array values
					// Third argument is the assignment value
					if len(args) == 3 {
						// Negative index value too small
						if indexValue < 0 {
							if arr.normalizeIndex(index) == -1 {
								return t.vm.initErrorObject(errors.InternalError, sourceLine, "Index value %d too small for array. minimum: %d", indexValue, -arr.length())
							}
							indexValue = arr.normalizeIndex(index)
						}

						a := args[2]
						assignedValue, isArray := a.(*ArrayObject)

						// Expand the array with nil case, we don't need to care the second argument
						// because the count in this case is useless
						if indexValue >= arr.length() {

							for arr.length() < indexValue {
								arr.Elements = append(arr.Elements, NULL)
							}

							if isArray {
								arr.Elements = append(arr.Elements, assignedValue.Elements...)
							} else {
								arr.Elements = append(arr.Elements, a)
							}
							return a
						}

						c := args[1]
						count, ok := c.(*IntegerObject)

						// Second argument must be integer
						if !ok {
							return t.vm.initErrorObject(errors.TypeError, sourceLine, errors.WrongArgumentTypeFormat, classes.IntegerClass, args[1].Class().Name)
						}

						countValue := count.value
						// Second argument must be positive value
						if countValue < 0 {
							return t.vm.initErrorObject(errors.InternalError, sourceLine, "Expect second argument greater than or equal 0. got: %d", countValue)
						}

						endValue := indexValue + countValue
						// Addition of index and count is too large
						if endValue > arr.length() {
							endValue = arr.length()
						}

						arr.Elements = append(arr.Elements[:indexValue], arr.Elements[endValue:]...)

						// If assigned value is array, then splat the array and push each element in the receiver
						// according to the first and second argument
						if isArray {
							arr.Elements = append(arr.Elements[:indexValue], append(assignedValue.Elements, arr.Elements[indexValue:]...)...)
						} else {
							arr.Elements = append(arr.Elements[:indexValue], append([]Object{a}, arr.Elements[indexValue:]...)...)
						}

						return a
					}

					// <Two Argument Case>
					// Second argument is the assignment value

					// Negative index value condition
					if indexValue < 0 {
						if len(arr.Elements) < -indexValue {
							return t.vm.initErrorObject(errors.InternalError, sourceLine, "Index value %d too small for array. minimum: %d", indexValue, -arr.length())
						}
						arr.Elements[len(arr.Elements)+indexValue] = args[1]
						return arr.Elements[len(arr.Elements)+indexValue]
					}

					// Expand the array
					if len(arr.Elements) < (indexValue + 1) {
						newArr := make([]Object, indexValue+1)
						copy(newArr, arr.Elements)
						for i := len(arr.Elements); i <= indexValue; i++ {
							newArr[i] = NULL
						}
						arr.Elements = newArr
					}

					arr.Elements[indexValue] = args[1]

					return arr.Elements[indexValue]
				}
			},
		},
		{
			// Passes each element of the collection to the given block. The method returns true if
			// the block ever returns a value other than false or nil
			//
			// ```ruby
			// a = [1, 2, 3]
			//
			// a.any? do |e|
			//   e == 2
			// end            # => true
			// a.any? do |e|
			//   e
			// end            # => true
			// a.any? do |e|
			//   e == 5
			// end            # => false
			// a.any? do |e|
			//   nil
			// end            # => false
			//
			// a = []
			//
			// a.any? do |e|
			//   true
			// end            # => false
			// ```
			Name: "any?",
			Fn: func(receiver Object, sourceLine int) builtinMethodBody {
				return func(t *thread, args []Object, blockFrame *normalCallFrame) Object {
					arr := receiver.(*ArrayObject)

					if blockFrame == nil {
						return t.vm.initErrorObject(errors.InternalError, sourceLine, errors.CantYieldWithoutBlockFormat)
					}

					if len(arr.Elements) == 0 {
						t.callFrameStack.pop()
					}

					for _, obj := range arr.Elements {
						result := t.builtinMethodYield(blockFrame, obj)

						if result.Target.isTruthy() {
							return TRUE
						}
					}

					return FALSE
				}
			},
		},
		{
			// Retrieves an object in an array using the index argument.
			// The index is 0-based; nil is returned when trying to access the index out of bounds.
			//
			// ```ruby
			// a = [1, 2, 3]
			// a.at(0)  # => 1
			// a.at(10) # => nil
			// a.at(-2) # => 2
			// a.at(-4) # => nil
			// ```
			Name: "at",
			Fn: func(receiver Object, sourceLine int) builtinMethodBody {
				return func(t *thread, args []Object, blockFrame *normalCallFrame) Object {
					if len(args) != 1 {
						return t.vm.initErrorObject(errors.ArgumentError, sourceLine, "Expect 1 arguments. got=%d", len(args))
					}
					arr := receiver.(*ArrayObject)
					return arr.index(t, args, sourceLine)
				}
			},
		},
		{
			// Removes all elements in the array and returns an empty array.
			//
			// ```ruby
			// a = [1, 2, 3]
			// a.clear # => []
			// a       # => []
			// ```
			Name: "clear",
			Fn: func(receiver Object, sourceLine int) builtinMethodBody {
				return func(t *thread, args []Object, blockFrame *normalCallFrame) Object {
					if len(args) != 0 {
						return t.vm.initErrorObject(errors.ArgumentError, sourceLine, "Expect 0 argument. got=%d", len(args))
					}

					arr := receiver.(*ArrayObject)
					arr.Elements = []Object{}

					return arr
				}
			},
		},
		{
			// Appends any number of argument to the array.
			//
			// ```ruby
			// a = [1, 2, 3]
			// a.concat(4, 5, 6)
			// a # => [1, 2, 3, 4, 5, 6]
			// ```
			Name: "concat",
			Fn: func(receiver Object, sourceLine int) builtinMethodBody {
				return func(t *thread, args []Object, blockFrame *normalCallFrame) Object {
					arr := receiver.(*ArrayObject)

					for _, arg := range args {
						addAr, ok := arg.(*ArrayObject)

						if !ok {
							return t.vm.initErrorObject(errors.TypeError, sourceLine, errors.WrongArgumentTypeFormat, classes.ArrayClass, arg.Class().Name)
						}

						for _, el := range addAr.Elements {
							arr.Elements = append(arr.Elements, el)
						}
					}

					return arr
				}
			},
		},
		{
			// Loop through each element with the given block.
			// Return the sum of elements that return true from yield.
			//
			// ```ruby
			// a = [1, 2, 3, 4, 5]
			//
			// a.count do |e|
			//   e * 2 > 3
			// end
			// # => 4
			// ```
			Name: "count",
			Fn: func(receiver Object, sourceLine int) builtinMethodBody {
				return func(t *thread, args []Object, blockFrame *normalCallFrame) Object {
					arr := receiver.(*ArrayObject)
					var count int

					if len(args) > 1 {
						return t.vm.initErrorObject(errors.ArgumentError, sourceLine, "Expect 1 argument, got=%d", len(args))
					}

					if blockFrame != nil {
						if len(arr.Elements) == 0 {
							t.callFrameStack.pop()
						}

						for _, obj := range arr.Elements {
							result := t.builtinMethodYield(blockFrame, obj)
							if result.Target.isTruthy() {
								count++
							}
						}

						return t.vm.initIntegerObject(count)
					}

					if len(args) == 0 {
						return t.vm.initIntegerObject(len(arr.Elements))
					}

					arg := args[0]
					findInt, findIsInt := arg.(*IntegerObject)
					findString, findIsString := arg.(*StringObject)
					findBoolean, findIsBoolean := arg.(*BooleanObject)

					for i := 0; i < len(arr.Elements); i++ {
						el := arr.Elements[i]
						switch el.(type) {
						case *IntegerObject:
							elInt := el.(*IntegerObject)
							if findIsInt && findInt.equal(elInt) {
								count++
							}
						case *StringObject:
							elString := el.(*StringObject)
							if findIsString && findString.equal(elString) {
								count++
							}
						case *BooleanObject:
							elBoolean := el.(*BooleanObject)
							if findIsBoolean && findBoolean.equal(elBoolean) {
								count++
							}
						}
					}

					return t.vm.initIntegerObject(count)
				}
			},
		},
		{
			// Deletes the element at the given position.
			// Returns the removed element.
			// The index is 0-based; nil is returned when using an out of bounds index.
			//
			// ```ruby
			// a = ["a", "b", "c"]
			// a.delete_at(1) # => "b"
			// a.delete_at(-1) # => "c"
			// a       # => ["a"]
			// ```
			Name: "delete_at",
			Fn: func(receiver Object, sourceLine int) builtinMethodBody {
				return func(t *thread, args []Object, blockFrame *normalCallFrame) Object {
					if len(args) != 1 {
						return t.vm.initErrorObject(errors.ArgumentError, sourceLine, "Expect 1 argument. got=%d", len(args))
					}

					i := args[0]
					index, ok := i.(*IntegerObject)

					if !ok {
						return t.vm.initErrorObject(errors.TypeError, sourceLine, errors.WrongArgumentTypeFormat, classes.IntegerClass, args[0].Class().Name)
					}

					arr := receiver.(*ArrayObject)
					normalizedIndex := arr.normalizeIndex(index)

					if normalizedIndex == -1 {
						return NULL
					}

					// delete and slice

					deletedValue := arr.Elements[normalizedIndex]

					arr.Elements = append(arr.Elements[:normalizedIndex], arr.Elements[normalizedIndex+1:]...)

					return deletedValue
				}
			},
		},
		{
			// Extracts the nested value specified by the sequence of idx objects by calling `dig` at
			// each step, returning nil if any intermediate step is nil.
			//
			// ```Ruby
			// [1 , 2].dig(-2)      # => 1
			// [[], 2].dig(0, 1)    # => nil
			// [[], 2].dig(0, 1, 2) # => nil
			// [1, 2].dig(0, 1)     # => TypeError: Expect target to be Diggable
			// ```
			//
			// @return [Object]
			Name: "dig",
			Fn: func(receiver Object, sourceLine int) builtinMethodBody {
				return func(t *thread, args []Object, blockFrame *normalCallFrame) Object {
					if len(args) == 0 {
						return t.vm.initErrorObject(errors.ArgumentError, sourceLine, "Expected 1+ arguments, got 0")
					}

					array := receiver.(*ArrayObject)
					value := array.dig(t, args, sourceLine)

					return value
				}
			},
		},
		{
			// Loop through each element with the given block.
			//
			// ```ruby
			// a = ["a", "b", "c"]
			//
			// a.each do |e|
			//   puts(e + e)
			// end
			// # => "aa"
			// # => "bb"
			// # => "cc"
			// ```
			Name: "each",
			Fn: func(receiver Object, sourceLine int) builtinMethodBody {
				return func(t *thread, args []Object, blockFrame *normalCallFrame) Object {
					if len(args) != 0 {
						return t.vm.initErrorObject(errors.ArgumentError, sourceLine, "Expect 0 argument. got=%d", len(args))
					}

					if blockFrame == nil {
						return t.vm.initErrorObject(errors.InternalError, sourceLine, errors.CantYieldWithoutBlockFormat)
					}

					arr := receiver.(*ArrayObject)

					// If it's an empty array, pop the block's call frame
					if len(arr.Elements) == 0 {
						t.callFrameStack.pop()
					}

					for _, obj := range arr.Elements {
						t.builtinMethodYield(blockFrame, obj)
					}
					return arr
				}
			},
		},
		{
			Name: "each_index",
			Fn: func(receiver Object, sourceLine int) builtinMethodBody {
				return func(t *thread, args []Object, blockFrame *normalCallFrame) Object {
					if len(args) != 0 {
						return t.vm.initErrorObject(errors.ArgumentError, sourceLine, "Expect 0 argument. got=%d", len(args))
					}

					if blockFrame == nil {
						return t.vm.initErrorObject(errors.InternalError, sourceLine, errors.CantYieldWithoutBlockFormat)
					}

					arr := receiver.(*ArrayObject)

					// If it's an empty array, pop the block's call frame
					if len(arr.Elements) == 0 {
						t.callFrameStack.pop()
					}

					for i := range arr.Elements {
						t.builtinMethodYield(blockFrame, t.vm.initIntegerObject(i))
					}
					return arr
				}
			},
		},
		{
			// Returns if the array"s length is 0 or not.
			//
			// ```ruby
			// [1, 2, 3].empty? # => false
			// [].empty? # => true
			// ```
			// @return [Boolean]
			Name: "empty?",
			Fn: func(receiver Object, sourceLine int) builtinMethodBody {
				return func(t *thread, args []Object, blockFrame *normalCallFrame) Object {

					if len(args) != 0 {
						return t.vm.initErrorObject(errors.ArgumentError, sourceLine, "Expect 0 argument. got=%d", len(args))
					}

					arr := receiver.(*ArrayObject)

					if arr.length() == 0 {
						return TRUE
					}

					return FALSE
				}
			},
		},
		{
			// Returns the first element of the array.
			Name: "first",
			Fn: func(receiver Object, sourceLine int) builtinMethodBody {
				return func(t *thread, args []Object, blockFrame *normalCallFrame) Object {
					if len(args) > 1 {
						return t.vm.initErrorObject(errors.ArgumentError, sourceLine, "Expect 0..1 argument. got=%d", len(args))
					}

					arr := receiver.(*ArrayObject)
					arrLength := len(arr.Elements)

					if arrLength == 0 {
						return NULL
					}

					if len(args) == 0 {
						return arr.Elements[0]
					}

					arg, ok := args[0].(*IntegerObject)

					if !ok {
						return t.vm.initErrorObject(errors.TypeError, sourceLine, errors.WrongArgumentTypeFormat, classes.IntegerClass, args[0].Class().Name)
					}

					if arg.value < 1 {
						return t.vm.initErrorObject(errors.ArgumentError, sourceLine, "Expect argument to be positive value. got=%d", arg.value)
					}

					if arrLength > arg.value {
						return t.vm.initArrayObject(arr.Elements[:arg.value])
					}
					return arr
				}
			},
		},
		{
			// Returns a new array that is a one-dimensional flattening of self.
			//
			// ```ruby
			// a = [ 1, 2, 3 ]
			// b = [ 4, 5, 6, [7, 8] ]
			// c = [ a, b, 9, 10 ] # => [[1, 2, 3], [4, 5, 6, [7, 8]], 9, 10]
			// c.flatten # => [1, 2, 3, 4, 5, 6, 7, 8, 9, 10]
			// ```
			// @return [Array]
			Name: "flatten",
			Fn: func(receiver Object, sourceLine int) builtinMethodBody {
				return func(t *thread, args []Object, blockFrame *normalCallFrame) Object {
					arr := receiver.(*ArrayObject)

					if len(args) != 0 {
						return t.vm.initErrorObject(errors.ArgumentError, sourceLine, "Expect 0 argument. got=%d", len(args))
					}

					newElements := arr.flatten()

					return t.vm.initArrayObject(newElements)
				}
			},
		},
		{
			// Returns a string by concatenating each element to string, separated by given separator.
			// If separator is nil, it uses empty string.
			//
			// ```ruby
			// [ 1, 2, 3 ].join # => "123"
			// [ 1, 2, 3 ].join("-") # => "1-2-3"
			// ```
			// @param separator [String]
			// @return [String]
			Name: "join",
			Fn: func(receiver Object, sourceLine int) builtinMethodBody {
				return func(t *thread, args []Object, blockFrame *normalCallFrame) Object {
					arr := receiver.(*ArrayObject)

					var sep string
					if len(args) == 0 {
						sep = ""
					} else if len(args) == 1 {
						arg, ok := args[0].(*StringObject)

						if !ok {
							return t.vm.initErrorObject(errors.TypeError, sourceLine, errors.WrongArgumentTypeFormat, classes.StringClass, args[0].Class().Name)
						}

						sep = arg.value
					} else {
						return t.vm.initErrorObject(errors.ArgumentError, sourceLine, "Expect 0 or 1 argument. got=%d", len(args))
					}

					elements := []string{}
					for _, e := range arr.flatten() {
						elements = append(elements, e.toString())
					}

					return t.vm.initStringObject(strings.Join(elements, sep))
				}
			},
		},
		{
			// Returns the last element of the array.
			Name: "last",
			Fn: func(receiver Object, sourceLine int) builtinMethodBody {
				return func(t *thread, args []Object, blockFrame *normalCallFrame) Object {
					if len(args) > 1 {
						return t.vm.initErrorObject(errors.ArgumentError, sourceLine, "Expect 0..1 argument. got=%d", len(args))
					}

					arr := receiver.(*ArrayObject)
					arrLength := len(arr.Elements)

					if len(args) == 0 {
						return arr.Elements[arrLength-1]
					}

					arg, ok := args[0].(*IntegerObject)

					if !ok {
						return t.vm.initErrorObject(errors.TypeError, sourceLine, errors.WrongArgumentTypeFormat, classes.IntegerClass, args[0].Class().Name)
					}

					if arg.value < 1 {
						return t.vm.initErrorObject(errors.ArgumentError, sourceLine, "Expect argument to be positive value. got=%d", arg.value)
					}

					if arrLength > arg.value {
						return t.vm.initArrayObject(arr.Elements[arrLength-arg.value : arrLength])
					}
					return arr
				}
			},
		},
		{
			// Returns the length of the array.
			//
			// ```ruby
			// [1, 2, 3].length # => 3
			// ```
			// @return [Integer]
			Name: "length",
			Fn: func(receiver Object, sourceLine int) builtinMethodBody {
				return func(t *thread, args []Object, blockFrame *normalCallFrame) Object {

					if len(args) != 0 {
						return t.vm.initErrorObject(errors.ArgumentError, sourceLine, "Expect 0 argument. got=%d", len(args))
					}

					arr := receiver.(*ArrayObject)
					return t.vm.initIntegerObject(arr.length())
				}
			},
		},
		{
			// Loop through each element with the given block. Return a new array with each yield element. Block is required.
			//
			// ```ruby
			// a = ["a", "b", "c"]
			//
			// a.map do |e|
			//   e + e
			// end
			// # => ["aa", "bb", "cc"]
			// ```
			//
			// @return [Array]
			Name: "map",
			Fn: func(receiver Object, sourceLine int) builtinMethodBody {
				return func(t *thread, args []Object, blockFrame *normalCallFrame) Object {
					arr := receiver.(*ArrayObject)
					var elements = make([]Object, len(arr.Elements))

					if blockFrame == nil {
						return t.vm.initErrorObject(errors.InternalError, sourceLine, errors.CantYieldWithoutBlockFormat)
					}

					// If it's an empty array, pop the block's call frame
					if len(arr.Elements) == 0 {
						t.callFrameStack.pop()
					}

					for i, obj := range arr.Elements {
						result := t.builtinMethodYield(blockFrame, obj)
						elements[i] = result.Target
					}

					return t.vm.initArrayObject(elements)
				}
			},
		},
		{
			// Removes the last element in the array and returns it.
			//
			// ```ruby
			// a = [1, 2, 3]
			// a.pop # => 3
			// a     # => [1, 2]
			// ```
			Name: "pop",
			Fn: func(receiver Object, sourceLine int) builtinMethodBody {
				return func(t *thread, args []Object, blockFrame *normalCallFrame) Object {

					if len(args) != 0 {
						return t.vm.initErrorObject(errors.ArgumentError, sourceLine, "Expect 0 argument. got=%d", len(args))
					}

					arr := receiver.(*ArrayObject)
					return arr.pop()
				}
			},
		},
		{
			// Appends the given object to the array and returns the array.
			//
			// ```ruby
			// a = [1, 2, 3]
			// a.push(4) # => [1, 2, 3, 4]
			// ```
			Name: "push",
			Fn: func(receiver Object, sourceLine int) builtinMethodBody {
				return func(t *thread, args []Object, blockFrame *normalCallFrame) Object {

					arr := receiver.(*ArrayObject)
					return arr.push(args)
				}
			},
		},
		{
			// Loop through each elements and accumulate each results of given block in the first argument of the block
			// If you do not give an argument, the first element of collection is used as an initial value
			//
			// ```ruby
			// a = [1, 2, 7]
			//
			// a.reduce do |sum, n|
			//   sum + n
			// end
			// # => 10
			//
			// a.reduce(10) do |sum, n|
			//   sum + n
			// end
			// # => 20
			// ```
			Name: "reduce",
			Fn: func(receiver Object, sourceLine int) builtinMethodBody {
				return func(t *thread, args []Object, blockFrame *normalCallFrame) Object {
					arr := receiver.(*ArrayObject)
					if blockFrame == nil {
						return t.vm.initErrorObject(errors.InternalError, sourceLine, errors.CantYieldWithoutBlockFormat)
					}

					// If it's an empty array, pop the block's call frame
					if len(arr.Elements) == 0 {
						t.callFrameStack.pop()
					}

					var prev Object
					var start int
					if len(args) == 0 {
						prev = arr.Elements[0]
						start = 1
					} else if len(args) == 1 {
						prev = args[0]
						start = 0
					} else {
						return t.vm.initErrorObject(errors.ArgumentError, sourceLine, "Expect 0 or 1 argument. got=%d", len(args))
					}

					for i := start; i < len(arr.Elements); i++ {
						result := t.builtinMethodYield(blockFrame, prev, arr.Elements[i])
						prev = result.Target
					}

					return prev
				}
			},
		},
		{
			// Returns a new array containing self‘s elements in reverse order.
			//
			// ```ruby
			// a = [1, 2, 7]
			//
			// a.reverse # => [7, 2, 1]
			// ```
			Name: "reverse",
			Fn: func(receiver Object, sourceLine int) builtinMethodBody {
				return func(t *thread, args []Object, blockFrame *normalCallFrame) Object {
					if len(args) != 0 {
						return t.vm.initErrorObject(errors.ArgumentError, sourceLine, "Expect 0 argument. got=%d", len(args))
					}

					arr := receiver.(*ArrayObject)

					return arr.reverse()
				}
			},
		},
		{
			// Same as #each, but traverses self in reverse order.
			//
			// ```ruby
			// a = ["a", "b", "c"]
			//
			// a.each do |e|
			//   puts(e + e)
			// end
			// # => "cc"
			// # => "bb"
			// # => "aa"
			// ```
			Name: "reverse_each",
			Fn: func(receiver Object, sourceLine int) builtinMethodBody {
				return func(t *thread, args []Object, blockFrame *normalCallFrame) Object {
					if len(args) != 0 {
						return t.vm.initErrorObject(errors.ArgumentError, sourceLine, "Expect 0 argument. got=%d", len(args))
					}

					if blockFrame == nil {
						return t.vm.initErrorObject(errors.InternalError, sourceLine, errors.CantYieldWithoutBlockFormat)
					}

					arr := receiver.(*ArrayObject)

					// If it's an empty array, pop the block's call frame
					if len(arr.Elements) == 0 {
						t.callFrameStack.pop()
					}

					reversedArr := arr.reverse()

					for _, obj := range reversedArr.Elements {
						t.builtinMethodYield(blockFrame, obj)
					}

					return reversedArr
				}
			},
		},
		{
			// Returns a new array by putting the desired element as the first element.
			// Use integer index as an argument to retrieve the element.
			//
			// ```ruby
			// a = ["a", "b", "c", "d"]
			//
			// a.rotate    # => ["b", "c", "d", "a"]
			// a.rotate(2) # => ["c", "d", "a", "b"]
			// a.rotate(3) # => ["d", "a", "b", "c"]
			// ```
			Name: "rotate",
			Fn: func(receiver Object, sourceLine int) builtinMethodBody {
				return func(t *thread, args []Object, blockFrame *normalCallFrame) Object {
					if len(args) > 1 {
						return t.vm.initErrorObject(errors.ArgumentError, sourceLine, "Expect 0..1 argument. got=%d", len(args))
					}

					arr := receiver.(*ArrayObject)
					rotArr := t.vm.initArrayObject(arr.Elements)

					rotate := 1

					if len(args) != 0 {
						arg, ok := args[0].(*IntegerObject)
						if !ok {
							return t.vm.initErrorObject(errors.TypeError, sourceLine, errors.WrongArgumentTypeFormat, classes.IntegerClass, args[0].Class().Name)
						}
						rotate = arg.value
					}

					for i := 0; i < rotate; i++ {
						el := rotArr.shift()
						rotArr.push([]Object{el})
					}

					return rotArr
				}
			},
		},
		{
			// Loop through each element with the given block.
			// Return a new array with each element that returns true from yield.
			//
			// ```ruby
			// a = [1, 2, 3, 4, 5]
			//
			// a.select do |e|
			//   e + 1 > 3
			// end
			// # => [3, 4, 5]
			// ```
			Name: "select",
			Fn: func(receiver Object, sourceLine int) builtinMethodBody {
				return func(t *thread, args []Object, blockFrame *normalCallFrame) Object {
					arr := receiver.(*ArrayObject)
					var elements []Object

					if blockFrame == nil {
						return t.vm.initErrorObject(errors.InternalError, sourceLine, errors.CantYieldWithoutBlockFormat)
					}

					// If it's an empty array, pop the block's call frame
					if len(arr.Elements) == 0 {
						t.callFrameStack.pop()
					}

					for _, obj := range arr.Elements {
						result := t.builtinMethodYield(blockFrame, obj)
						if result.Target.isTruthy() {
							elements = append(elements, obj)
						}
					}

					return t.vm.initArrayObject(elements)
				}
			},
		},
		{
			// Removes the first element in the array and returns it.
			//
			// ```ruby
			// a = [1, 2, 3]
			// a.shift # => 1
			// a       # => [2, 3]
			// ```
			Name: "shift",
			Fn: func(receiver Object, sourceLine int) builtinMethodBody {
				return func(t *thread, args []Object, blockFrame *normalCallFrame) Object {
					if len(args) != 0 {
						return t.vm.initErrorObject(errors.ArgumentError, sourceLine, "Expect 0 argument. got=%d", len(args))
					}

					arr := receiver.(*ArrayObject)
					return arr.shift()
				}
			},
		},
		{
			// Inserts the specified element in the first position of the array.
			//
			// ```ruby
			// a = [1, 2]
			// a.unshift(0) # => [0, 1, 2]
			// a            # => [0, 1, 2]
			// ```
			Name: "unshift",
			Fn: func(receiver Object, sourceLine int) builtinMethodBody {
				return func(t *thread, args []Object, blockFrame *normalCallFrame) Object {
					arr := receiver.(*ArrayObject)
					return arr.unshift(args)
				}
			},
		},
		{
			// Returns an array containing the elements in self corresponding to the given indexes.
			//
			// ```ruby
			// a = ["a", "b", "c"]
			// a.values_at(1)     # => ["b"]
			// a.values_at(-1, 3) # => ["c", nil]
			// a.values_at()      # => []
			// ```
			Name: "values_at",
			Fn: func(receiver Object, sourceLine int) builtinMethodBody {
				return func(t *thread, args []Object, blockFrame *normalCallFrame) Object {
					arr := receiver.(*ArrayObject)
					var elements = make([]Object, len(args))

					for i, arg := range args {
						index, ok := arg.(*IntegerObject)

						if !ok {
							return t.vm.initErrorObject(errors.TypeError, sourceLine, errors.WrongArgumentTypeFormat, classes.IntegerClass, arg.Class().Name)
						}

						if index.value >= len(arr.Elements) {
							elements[i] = NULL
						} else if index.value < 0 && -index.value > len(arr.Elements) {
							elements[i] = NULL
						} else if index.value < 0 {
							elements[i] = arr.Elements[len(arr.Elements)+index.value]
						} else {
							elements[i] = arr.Elements[index.value]
						}
					}

					return t.vm.initArrayObject(elements)
				}
			},
		},
	}
}

// Internal functions ===================================================

// Functions for initialization -----------------------------------------

func (vm *VM) initArrayObject(elements []Object) *ArrayObject {
	return &ArrayObject{
		baseObj:  &baseObj{class: vm.topLevelClass(classes.ArrayClass)},
		Elements: elements,
	}
}

func (vm *VM) initArrayClass() *RClass {
	ac := vm.initializeClass(classes.ArrayClass, false)
	ac.setBuiltinMethods(builtinArrayInstanceMethods(), false)
	ac.setBuiltinMethods(builtinArrayClassMethods(), true)
	vm.libFiles = append(vm.libFiles, "array.gb")
	return ac
}

// Polymorphic helper functions -----------------------------------------

// Value returns the elements from the object
func (a *ArrayObject) Value() interface{} {
	return a.Elements
}

// toString returns the object's elements as the string format
func (a *ArrayObject) toString() string {
	var out bytes.Buffer

	elements := []string{}
	for _, e := range a.Elements {
		_, isString := e.(*StringObject)
		if isString {
			elements = append(elements, "\""+e.toString()+"\"")
		} else {
			elements = append(elements, e.toString())
		}
	}

	out.WriteString("[")
	out.WriteString(strings.Join(elements, ", "))
	out.WriteString("]")

	return out.String()
}

// toJSON returns the object's elements as the JSON string format
func (a *ArrayObject) toJSON() string {
	var out bytes.Buffer
	elements := []string{}
	for _, e := range a.Elements {
		elements = append(elements, e.toJSON())
	}

	out.WriteString("[")
	out.WriteString(strings.Join(elements, ", "))
	out.WriteString("]")

	return out.String()
}

// concatenateCopies returns a array composed of N copies of the array
func (a *ArrayObject) concatenateCopies(t *thread, n *IntegerObject) Object {
	aLen := len(a.Elements)
	result := make([]Object, 0, aLen*n.value)

	for i := 0; i < n.value; i++ {
		result = append(result, a.Elements...)
	}

	return t.vm.initArrayObject(result)
}

// recursive indexed access - see ArrayObject#dig documentation.
func (a *ArrayObject) dig(t *thread, keys []Object, sourceLine int) Object {
	currentKey := keys[0]
	intCurrentKey, ok := currentKey.(*IntegerObject)

	if !ok {
		return t.vm.initErrorObject(errors.TypeError, sourceLine, errors.WrongArgumentTypeFormat, classes.IntegerClass, currentKey.Class().Name)
	}

	normalizedIndex := a.normalizeIndex(intCurrentKey)

	if normalizedIndex == -1 {
		return NULL
	}

	nextKeys := keys[1:]
	currentValue := a.Elements[normalizedIndex]

	if len(nextKeys) == 0 {
		return currentValue
	}

	diggableCurrentValue, ok := currentValue.(Diggable)

	if !ok {
		return t.vm.initErrorObject(errors.TypeError, sourceLine, "Expect target to be Diggable, got %s", currentValue.Class().Name)
	}

	return diggableCurrentValue.dig(t, nextKeys, sourceLine)
}

// Retrieves an object in an array using Integer index; common to `[]` and `at()`.
func (a *ArrayObject) index(t *thread, args []Object, sourceLine int) Object {
	if len(args) > 2 || len(args) == 0 {
		return t.vm.initErrorObject(errors.ArgumentError, sourceLine, "Expect 1..2 arguments. got=%d", len(args))
	}

	i := args[0]
	index, ok := i.(*IntegerObject)

	if !ok {
		return t.vm.initErrorObject(errors.TypeError, sourceLine, errors.WrongArgumentTypeFormat, classes.IntegerClass, args[0].Class().Name)
	}

	normalizedIndex := a.normalizeIndex(index)
	if normalizedIndex == -1 {
		return NULL
	}

	if len(args) == 2 {
		j := args[1]
		count, ok := j.(*IntegerObject)

		if !ok {
			return t.vm.initErrorObject(errors.TypeError, sourceLine, errors.WrongArgumentTypeFormat, classes.IntegerClass, args[1].Class().Name)
		}

		if count.value < 0 {
			return NULL
		}

		if normalizedIndex+count.value > len(a.Elements) {
			return t.vm.initArrayObject(a.Elements[normalizedIndex:])
		}
		return t.vm.initArrayObject(a.Elements[normalizedIndex : normalizedIndex+count.value])
	}

	return a.Elements[normalizedIndex]
}

// flatten returns a array of Objects that is one-dimensional flattening of Elements
func (a *ArrayObject) flatten() []Object {
	var result []Object

	for _, e := range a.Elements {
		arr, isArray := e.(*ArrayObject)
		if isArray {
			result = append(result, arr.flatten()...)
		} else {
			result = append(result, e)
		}
	}

	return result
}

// length returns the length of array's elements
func (a *ArrayObject) length() int {
	return len(a.Elements)
}

// normalizes the index to the Ruby-style:
//
// 1. if the index is between o and the index length, returns the index
// 2. if it's a negative value (within bounds), returns the normalized positive version
// 3. if it's out of bounds (either positive or negative), returns -1
func (a *ArrayObject) normalizeIndex(objectIndex *IntegerObject) int {
	aLength := len(a.Elements)
	index := objectIndex.value

	// out of bounds

	if index >= aLength {
		return -1
	}

	if index < 0 && -index > aLength {
		return -1
	}

	// within bounds

	if index < 0 {
		return aLength + index
	}

	return index
}

// pop removes the last element in the array and returns it
func (a *ArrayObject) pop() Object {
	if len(a.Elements) < 1 {
		return NULL
	}

	value := a.Elements[len(a.Elements)-1]
	a.Elements = a.Elements[:len(a.Elements)-1]
	return value
}

// push appends given object into array and returns the array object
func (a *ArrayObject) push(objs []Object) *ArrayObject {
	a.Elements = append(a.Elements, objs...)
	return a
}

// returns a reversed copy of the passed array
func (a *ArrayObject) reverse() *ArrayObject {
	arrLen := len(a.Elements)
	reversedArrElems := make([]Object, arrLen)

	for i, element := range a.Elements {
		reversedArrElems[arrLen-i-1] = element
	}

	newArr := &ArrayObject{
		baseObj:  &baseObj{class: a.class},
		Elements: reversedArrElems,
	}

	return newArr
}

// shift removes the first element in the array and returns it
func (a *ArrayObject) shift() Object {
	if len(a.Elements) < 1 {
		return NULL
	}

	value := a.Elements[0]
	a.Elements = a.Elements[1:]
	return value
}

// copy returns the duplicate of the Array object
func (a *ArrayObject) copy() Object {
	elems := make([]Object, len(a.Elements))

	copy(elems, a.Elements)

	newArr := &ArrayObject{
		baseObj:  &baseObj{class: a.class},
		Elements: elems,
	}

	return newArr
}

// unshift inserts an element in the first position of the array
func (a *ArrayObject) unshift(objs []Object) *ArrayObject {
	a.Elements = append(objs, a.Elements...)
	return a
}
