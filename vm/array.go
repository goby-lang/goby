package vm

import (
	"bytes"
	"strings"

	"github.com/goby-lang/goby/vm/classes"
	"github.com/goby-lang/goby/vm/errors"
)

// ArrayObject represents an instance from Array class.
// An array is a collection of different objects that are ordered and indexed.
// Elements in an array can belong to any class and you can also build a "tuple" within an array.
// Array objects should always be enumerable.
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
				return func(t *Thread, args []Object, blockFrame *normalCallFrame) Object {
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
			// a[0]  #=> 1
			// a[3]  #=> "a"
			// a[10] #=> nil
			// a[-1] #=> "c"
			// a[-3] #=> "a"
			// a[-7] #=> nil
			//
			// # Double indexing, second argument specifies the count of the elements
			// a[1, 3]  #=> [2, 3, "a"]
			// a[1, 0]  #=> [] <-- Zero count is empty
			// a[1, 5]  #=> [2, 3, "a", "b", "c"]
			// a[1, 10] #=> [2, 3, "a", "b", "c"]
			// a[-3, 2] #=> ["a", "b"]
			// a[-3, 5] #=> ["a", "b", "c"]
			// a[5, 1]  #=> ["c"]
			// a[6, 1]  #=> []
			// a[7, 1]  #=> nil
			//
			// Special case 1:
			// a[6]    #=> nil
			// a[6, 1] #=> []  <-- Not nil!
			// a[7, 1] #=> nil <-- Because it is really out of the edge of the array
			//
			// Special case 2: Second argument is negative
			// This behaviour is different from Ruby itself, in Ruby, it returns "nil".
			// However, in Goby, it raises error because there cannot be negative count values.
			//
			// a[1, -1]  #=> ArgumentError: Expect second argument greater than or equal 0. got: -1
			// a[-4, -3] #=> ArgumentError: Expect second argument greater than or equal 0. got: -3
			//
			// Special case 3: First argument is negative and exceed the array length
			// a[-6, 1] #=> [1]
			// a[-6, 0] #=> []
			// a[-7, 1] #=> ArgumentError: Index value -7 too small for array. minimum: -6
			// a[-7, 0] #=> ArgumentError: Index value -7 too small for array. minimum: -6
			// ```
			//
			// Note:
			// * The notations such as `a.[](1)` or `a.[] 1` are unsupported.
			// * `Range` object is unsupported for now.
			//
			// @param index [Integer], (count [Integer])
			// @return [Array]
			Name: "[]",
			Fn: func(receiver Object, sourceLine int) builtinMethodBody {
				return func(t *Thread, args []Object, blockFrame *normalCallFrame) Object {
					arr := receiver.(*ArrayObject)
					return arr.index(t, args, sourceLine)
				}
			},
		},
		{
			// Repetition — returns a new array built by just concatenating the specified number of copies of `self`.
			//
			// ```ruby
			// a = [1, 2, 3]
			// a * 2   #=> [1, 2, 3, 1, 2, 3]
			// ```
			//
			// * The index should be a positive or zero Integer object.
			// * Ruby's syntax such as `[1, 2, 3] * ','` are unsupported. Use `#join` instead.
			//
			// @param zero or positive integer [Integer]
			// @return [Array]
			Name: "*",
			Fn: func(receiver Object, sourceLine int) builtinMethodBody {
				return func(t *Thread, args []Object, blockFrame *normalCallFrame) Object {
					if len(args) != 1 {
						return t.vm.InitErrorObject(errors.ArgumentError, sourceLine, "Expect 1 arguments. got=%d", len(args))
					}

					arr := receiver.(*ArrayObject)

					copiesNumber, ok := args[0].(*IntegerObject)

					if !ok {
						return t.vm.InitErrorObject(errors.TypeError, sourceLine, errors.WrongArgumentTypeFormat, classes.IntegerClass, args[0].Class().Name)
					}

					return arr.concatenateCopies(t, copiesNumber)
				}
			},
		},
		{
			// Concatenation: returns a new array by just concatenating the two arrays.
			//
			// ```ruby
			// a = [1, 2]
			// b + [3, 4]  #=> [1, 2, 3, 4]
			// ```
			//
			// @param array [Array]
			// @return [Array]
			Name: "+",
			Fn: func(receiver Object, sourceLine int) builtinMethodBody {
				return func(t *Thread, args []Object, blockFrame *normalCallFrame) Object {
					if len(args) != 1 {
						return t.vm.InitErrorObject(errors.ArgumentError, sourceLine, "Expect 1 arguments. got=%d", len(args))
					}

					otherArrayArg := args[0]
					otherArray, ok := otherArrayArg.(*ArrayObject)

					if !ok {
						return t.vm.InitErrorObject(errors.TypeError, sourceLine, errors.WrongArgumentTypeFormat, classes.ArrayClass, args[0].Class().Name)
					}

					selfArray := receiver.(*ArrayObject)

					newArrayelements := append(selfArray.Elements, otherArray.Elements...)

					newArray := t.vm.InitArrayObject(newArrayelements)

					return newArray
				}
			},
		},
		{
			// Assigns one or more values to an array. It requires one or two indices and a value as argument.
			// The first index should be Integer, and the second index should be zero or positive integer.
			// The array will expand if the assigned index is bigger than the current size of self.
			// Returns the assigned value.
			// The gaps will be filled with `nil`, but such operations should be avoided.
			//
			// ```ruby
			// a = []
			// a[0] = 10  #=> 10
			// a[3] = 20  #=> 20
			// a          #=> [10, nil, nil, 20]
			// a[-2] = 5  #=> [10, nil, 5, 20]
			//
			// # Double indexing, second argument specify the count of the arguments
			// a = [1, 2, 3, 4, 5]
			// a[2, 3] = [:a, :b, :c]   # <-- Common case: overridden
			// a #=> [1, 2, "a", "b", "c"]
			//
			// a = [1, 2, 3, 4, 5]
			// a[4, 4] = [:a, :b, :c]   # <- Exceeded case: the array will be expanded and `5` will be overridden
			// a #=> [1, 2, 3, 4, "a", "b", "c"]
			//
			// a = [1, 2, 3, 4, 5]
			// a[5, 1] = [:a, :b, :c]   # <-- Edge case: insertion
			// a #=> [1, 2, 3, 4, 5, "a", "b", "c"]
			//
			// a = [1, 2, 3, 4, 5]
			// a[8, 123] = [:a, :b, :c] # <-- Weak array case: the gaps will be filled with `nil` but the tailing ones not
			// a #=> [1, 2, 3, 4, 5, nil, nil, nil, "a", "b", "c"]
			//
			// a = [1, 2, 3, 4, 5]
			// a[3, 0] = [:a, :b, :c]   # <-- Insertion case: the second index `0` is to insert there
			// a #=> [1, 2, 3, "a", "b", "c", 4, 5]
			//
			// a = [1, 2, 3, 4, 5]
			// a[0, 3] = 12345          # <-- Assign non-array value case
			// a #=> [12345, 4, 5]
			//
			// a = [1, 2, 3, 4, 5]
			// a[-3, 2] = [:a, :b, :c]  # <-- Negative index assign case
			// a #=> [1, 2, "a", "b", "c", 5]
			//
			// a = [1, 2, 3, 4, 5]
			// a[-5, 3] = [:a, :b, :c]  # <-- Negative index edge case
			// a #=> ["a", "b", "c", 4, 5]
			//
			// a = [1, 2, 3, 4, 5]
			// a[-5, 4] = [:a, :b, :c]  # <-- Negative index exceeded case: `4` will be destroyed
			// a #=> ["a", "b", "c", 5]
			//
			// a = [1, 2, 3, 4, 5]
			// a[-5, 5] = [:a, :b, :c]  # <-- Negative index exceeded case: `4, 5` will be destroyed
			// a #=> ["a", "b", "c"]
			//
			// a = [1, 2, 3, 4, 5]
			// a[-6, 4] = [:a, :b, :c]     # <-- Invalid: Negative index too small case
			// # ArgumentError: Index value -6 too small for array. minimum: -5
			//
			// a = [1, 2, 3, 4, 5]
			// a[6, -4] = [9, 8, 7]     # <-- Weak array assignment with negative count case
			// # ArgumentError: Expect second argument greater than or equal 0. got: -4
			// ```
			//
			// Note that passing multiple values to the method is unavailable.
			//
			// @param index [Integer], object [Object]
			// @param index [Integer], count [Integer], object [Object]
			// @return [Array]
			Name: "[]=",
			Fn: func(receiver Object, sourceLine int) builtinMethodBody {
				return func(t *Thread, args []Object, blockFrame *normalCallFrame) Object {

					// First argument is an index: there exists two cases which will be described in the following code
					if len(args) != 2 && len(args) != 3 {
						return t.vm.InitErrorObject(errors.ArgumentError, sourceLine, "Expect 2..3 arguments. got=%d", len(args))
					}

					i := args[0]
					index, ok := i.(*IntegerObject)

					if !ok {
						return t.vm.InitErrorObject(errors.TypeError, sourceLine, errors.WrongArgumentTypeFormat, classes.IntegerClass, args[0].Class().Name)
					}

					indexValue := index.value
					arr := receiver.(*ArrayObject)

					// <Three Argument Case>
					// Second argument: the length of successive array values (zero or positive Integer)
					// Third argument: the assignment value (object)
					if len(args) == 3 {
						// Negative index value too small
						if indexValue < 0 {
							if arr.normalizeIndex(index) == -1 {
								return t.vm.InitErrorObject(errors.ArgumentError, sourceLine, "Index value %d too small for array. minimum: %d", indexValue, -arr.length())
							}
							indexValue = arr.normalizeIndex(index)
						}

						c := args[1]
						count, ok := c.(*IntegerObject)

						// Second argument must be an integer
						if !ok {
							return t.vm.InitErrorObject(errors.TypeError, sourceLine, errors.WrongArgumentTypeFormat, classes.IntegerClass, args[1].Class().Name)
						}

						countValue := count.value
						// Second argument must be a positive value
						if countValue < 0 {
							return t.vm.InitErrorObject(errors.ArgumentError, sourceLine, "Expect second argument greater than or equal 0. got: %d", countValue)
						}

						a := args[2]
						assignedValue, isArray := a.(*ArrayObject)

						// Expand the array with nil; the second index is unnecessary in the case
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

						endValue := indexValue + countValue
						// the case the addition of index and count is too large
						if endValue > arr.length() {
							endValue = arr.length()
						}

						arr.Elements = append(arr.Elements[:indexValue], arr.Elements[endValue:]...)

						// If assigned value is an array, then splat the array and push each element to the receiver
						// following the first and second indices
						if isArray {
							arr.Elements = append(arr.Elements[:indexValue], append(assignedValue.Elements, arr.Elements[indexValue:]...)...)
						} else {
							arr.Elements = append(arr.Elements[:indexValue], append([]Object{a}, arr.Elements[indexValue:]...)...)
						}

						return a
					}

					// <Two Argument Case>
					// Second argument is the assignment value (object)

					// Negative index value condition
					if indexValue < 0 {
						if len(arr.Elements) < -indexValue {
							return t.vm.InitErrorObject(errors.ArgumentError, sourceLine, "Index value %d too small for array. minimum: %d", indexValue, -arr.length())
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
			// A predicate method.
			// Evaluates the given block and returns `true` if the block ever returns a value.
			// Returns `false` if the evaluated block returns `false` or `nil`.
			//
			// ```ruby
			// a = [1, 2, 3]
			//
			// a.any? do |e|
			//   e == 2
			// end            #=> true
			// a.any? do |e|
			//   e
			// end            #=> true
			// a.any? do |e|
			//   e == 5
			// end            #=> false
			// a.any? do |e|
			//   nil
			// end            #=> false
			//
			// a = []
			//
			// a.any? do |e|
			//   true
			// end            #=> false
			// ```
			//
			// @param block [Block]
			// @return [Boolean]
			Name: "any?",
			Fn: func(receiver Object, sourceLine int) builtinMethodBody {
				return func(t *Thread, args []Object, blockFrame *normalCallFrame) Object {
					arr := receiver.(*ArrayObject)

					if blockFrame == nil {
						return t.vm.InitErrorObject(errors.InternalError, sourceLine, errors.CantYieldWithoutBlockFormat)
					}

					if blockIsEmpty(blockFrame) {
						return FALSE
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
			// Retrieves an object in an array using the given index.
			// The index is 0-based; `nil` is returned when trying to access the index out of bounds.
			//
			// ```ruby
			// a = [1, 2, 3]
			// a.at(0)  #=> 1
			// a.at(10) #=> nil
			// a.at(-2) #=> 2
			// a.at(-4) #=> nil
			// ```
			//
			// @param index [Integer]
			// @return [Object]
			Name: "at",
			Fn: func(receiver Object, sourceLine int) builtinMethodBody {
				return func(t *Thread, args []Object, blockFrame *normalCallFrame) Object {
					if len(args) != 1 {
						return t.vm.InitErrorObject(errors.ArgumentError, sourceLine, "Expect 1 arguments. got=%d", len(args))
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
			// a.clear #=> []
			// a       #=> []
			// ```
			//
			// @return [Array]
			Name: "clear",
			Fn: func(receiver Object, sourceLine int) builtinMethodBody {
				return func(t *Thread, args []Object, blockFrame *normalCallFrame) Object {
					if len(args) != 0 {
						return t.vm.InitErrorObject(errors.ArgumentError, sourceLine, "Expect 0 argument. got=%d", len(args))
					}

					arr := receiver.(*ArrayObject)
					arr.Elements = []Object{}

					return arr
				}
			},
		},
		{
			// Concatenation: returns a new array by just concatenating the two arrays.
			//
			// ```ruby
			// a = [1, 2, 3]
			// a.concat([4, 5, 6])
			// a #=> [1, 2, 3, 4, 5, 6]
			// ```
			//
			// @param array [Array]
			// @return [Array]
			Name: "concat",
			Fn: func(receiver Object, sourceLine int) builtinMethodBody {
				return func(t *Thread, args []Object, blockFrame *normalCallFrame) Object {
					arr := receiver.(*ArrayObject)

					for _, arg := range args {
						addAr, ok := arg.(*ArrayObject)

						if !ok {
							return t.vm.InitErrorObject(errors.TypeError, sourceLine, errors.WrongArgumentTypeFormat, classes.ArrayClass, arg.Class().Name)
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
			// If no block is given, just returns the count of the elements within the array.
			// If a block is given, evaluate each element of the array by the given block,
			// and then return the count of elements that return `true` by the block.
			//
			// ```ruby
			// a = [1, 2, 3, 4, 5]
			//
			// a.count do |e|
			//   e * 2 > 3
			// end
			// #=> 4
			// ```
			//
			// @param
			// @param block [Block]
			// @return [Integer]
			Name: "count",
			Fn: func(receiver Object, sourceLine int) builtinMethodBody {
				return func(t *Thread, args []Object, blockFrame *normalCallFrame) Object {
					arr := receiver.(*ArrayObject)
					var count int

					if len(args) > 1 {
						return t.vm.InitErrorObject(errors.ArgumentError, sourceLine, "Expect 1 argument, got=%d", len(args))
					}

					if blockFrame != nil {
						if blockIsEmpty(blockFrame) {
							return t.vm.InitIntegerObject(0)
						}
						if len(arr.Elements) == 0 {
							t.callFrameStack.pop()
						}

						for _, obj := range arr.Elements {
							result := t.builtinMethodYield(blockFrame, obj)
							if result.Target.isTruthy() {
								count++
							}
						}

						return t.vm.InitIntegerObject(count)
					}

					if len(args) == 0 {
						return t.vm.InitIntegerObject(len(arr.Elements))
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

					return t.vm.InitIntegerObject(count)
				}
			},
		},
		{
			// Deletes the element pointed by the given index.
			// Returns the removed element.
			// The method is destructive and the self is mutated.
			// The index is 0-based; `nil` is returned when using an out-of-bounds index.
			//
			// ```ruby
			// a = ["a", "b", "c"]
			// a.delete_at(1) #=> "b"
			// a.delete_at(-1) #=> "c"
			// a       #=> ["a"]
			// ```
			//
			// @param index [Integer]
			// @return [Object]
			Name: "delete_at",
			Fn: func(receiver Object, sourceLine int) builtinMethodBody {
				return func(t *Thread, args []Object, blockFrame *normalCallFrame) Object {
					if len(args) != 1 {
						return t.vm.InitErrorObject(errors.ArgumentError, sourceLine, "Expect 1 argument. got=%d", len(args))
					}

					i := args[0]
					index, ok := i.(*IntegerObject)

					if !ok {
						return t.vm.InitErrorObject(errors.TypeError, sourceLine, errors.WrongArgumentTypeFormat, classes.IntegerClass, args[0].Class().Name)
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
			// Returns the value from the nested array, specified by one or more indices,
			// Returns `nil` if one of the intermediate values are `nil`.
			//
			// ```Ruby
			// [1 , 2].dig(-2)      #=> 1
			// [[], 2].dig(0, 1)    #=> nil
			// [[], 2].dig(0, 1, 2) #=> nil
			// [[1, 2, [3, [8, [9]]]], 4, 5].dig(0, 2, 1, 1, 0) #=> 9
			// [1, 2].dig(0, 1)     #=> TypeError: Expect target to be Diggable
			// ```
			//
			// @param index [Integer]...
			// @return [Object]
			Name: "dig",
			Fn: func(receiver Object, sourceLine int) builtinMethodBody {
				return func(t *Thread, args []Object, blockFrame *normalCallFrame) Object {
					if len(args) == 0 {
						return t.vm.InitErrorObject(errors.ArgumentError, sourceLine, "Expected 1+ arguments, got 0")
					}

					array := receiver.(*ArrayObject)
					value := array.dig(t, args, sourceLine)

					return value
				}
			},
		},
		{
			// Loops through each element in the array, with the given block.
			// Returns self.
			// A block literal is required.
			//
			// ```ruby
			// a = ["a", "b", "c"]
			//
			// b = a.each do |e|
			//   puts(e + e)
			// end
			// #=> "aa"
			// #=> "bb"
			// #=> "cc"
			// puts b
			// #=> ["a", "b", "c"]
			// ```
			//
			// @param block literal
			// @return [Array]
			Name: "each",
			Fn: func(receiver Object, sourceLine int) builtinMethodBody {
				return func(t *Thread, args []Object, blockFrame *normalCallFrame) Object {
					if len(args) != 0 {
						return t.vm.InitErrorObject(errors.ArgumentError, sourceLine, "Expect 0 argument. got=%d", len(args))
					}

					if blockFrame == nil {
						return t.vm.InitErrorObject(errors.InternalError, sourceLine, errors.CantYieldWithoutBlockFormat)
					}

					arr := receiver.(*ArrayObject)
					if blockIsEmpty(blockFrame) {
						return arr
					}

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
		// Works like #each, but passes the index of the element instead of the element itself.
		// Returns self.
		// A block literal is required.
		//
		// ```ruby
		// a = [:apple, :orange, :grape, :melon]
		//
		// b = a.each_index do |i|
		//   puts(i*i)
		// end
		// #=> 0
		// #=> 1
		// #=> 4
		// #=> 9
		// puts b
		// #=> ["a", "b", "c"]
		// ```
		//
		// @param block literal
		// @return [Array]
		{
			Name: "each_index",
			Fn: func(receiver Object, sourceLine int) builtinMethodBody {
				return func(t *Thread, args []Object, blockFrame *normalCallFrame) Object {
					if len(args) != 0 {
						return t.vm.InitErrorObject(errors.ArgumentError, sourceLine, "Expect 0 argument. got=%d", len(args))
					}

					if blockFrame == nil {
						return t.vm.InitErrorObject(errors.InternalError, sourceLine, errors.CantYieldWithoutBlockFormat)
					}

					arr := receiver.(*ArrayObject)
					if blockIsEmpty(blockFrame) {
						return arr
					}

					// If it's an empty array, pop the block's call frame
					if len(arr.Elements) == 0 {
						t.callFrameStack.pop()
					}

					for i := range arr.Elements {
						t.builtinMethodYield(blockFrame, t.vm.InitIntegerObject(i))
					}
					return arr
				}
			},
		},
		{
			// A predicate method.
			// Returns if the array"s length is 0 or not.
			//
			// ```ruby
			// [1, 2, 3].empty? #=> false
			// [].empty?        #=> true
			// [[]].empty?      #=> false
			// ```
			//
			// @return [Boolean]
			Name: "empty?",
			Fn: func(receiver Object, sourceLine int) builtinMethodBody {
				return func(t *Thread, args []Object, blockFrame *normalCallFrame) Object {

					if len(args) != 0 {
						return t.vm.InitErrorObject(errors.ArgumentError, sourceLine, "Expect 0 argument. got=%d", len(args))
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
			// If a count 'n' is provided as an argument, it returns the array of the first n elements.
			//
			// ```ruby
			// [1, 2, 3].first                            #=> 1
			// [:apple, :orange, :grape, :melon].first    #=> "apple"
			// [:apple, :orange, :grape, :melon].first(2) #=> ["apple", "orange"]
			// ```
			//
			// @param count [Integer]
			// @return [Object]
			Name: "first",
			Fn: func(receiver Object, sourceLine int) builtinMethodBody {
				return func(t *Thread, args []Object, blockFrame *normalCallFrame) Object {
					if len(args) > 1 {
						return t.vm.InitErrorObject(errors.ArgumentError, sourceLine, "Expect 0..1 argument. got=%d", len(args))
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
						return t.vm.InitErrorObject(errors.TypeError, sourceLine, errors.WrongArgumentTypeFormat, classes.IntegerClass, args[0].Class().Name)
					}

					if arg.value < 1 {
						return t.vm.InitErrorObject(errors.ArgumentError, sourceLine, "Expect argument to be positive value. got=%d", arg.value)
					}

					if arrLength > arg.value {
						return t.vm.InitArrayObject(arr.Elements[:arg.value])
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
			// c = [ a, b, 9, 10 ] #=> [[1, 2, 3], [4, 5, 6, [7, 8]], 9, 10]
			// c.flatten #=> [1, 2, 3, 4, 5, 6, 7, 8, 9, 10]
			//
			// [[[1, 2], [[[3, 4]], [5, 6]]]].flatten
			// #=> [1, 2, 3, 4, 5, 6]
			// ```
			//
			// @return [Array]
			Name: "flatten",
			Fn: func(receiver Object, sourceLine int) builtinMethodBody {
				return func(t *Thread, args []Object, blockFrame *normalCallFrame) Object {
					arr := receiver.(*ArrayObject)

					if len(args) != 0 {
						return t.vm.InitErrorObject(errors.ArgumentError, sourceLine, "Expect 0 argument. got=%d", len(args))
					}

					newElements := arr.flatten()

					return t.vm.InitArrayObject(newElements)
				}
			},
		},
		{
			// Returns a string by concatenating each element to string, separated by given separator.
			// If the array is nested, they will be flattened and then concatenated.
			// If separator is nil, it uses empty string.
			//
			// ```ruby
			// [ 1, 2, 3 ].join                #=> "123"
			// [[:h, :e, :l], [[:l], :o]].join #=> "hello"
			// [[:hello],{k: :v}].join         #=> 'hello{ k: "v" }'
			// [ 1, 2, 3 ].join("-")           #=> "1-2-3"
			// ```
			//
			// @param separator [String]
			// @return [String]
			Name: "join",
			Fn: func(receiver Object, sourceLine int) builtinMethodBody {
				return func(t *Thread, args []Object, blockFrame *normalCallFrame) Object {
					arr := receiver.(*ArrayObject)

					var sep string
					if len(args) == 0 {
						sep = ""
					} else if len(args) == 1 {
						arg, ok := args[0].(*StringObject)

						if !ok {
							return t.vm.InitErrorObject(errors.TypeError, sourceLine, errors.WrongArgumentTypeFormat, classes.StringClass, args[0].Class().Name)
						}

						sep = arg.value
					} else {
						return t.vm.InitErrorObject(errors.ArgumentError, sourceLine, "Expect 0 or 1 argument. got=%d", len(args))
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
			// If a count 'n' is provided as an argument, it returns the array of the last n elements.
			//
			// ```ruby
			// [1, 2, 3].last                            #=> 3
			// [:apple, :orange, :grape, :melon].last    #=> "melon"
			// [:apple, :orange, :grape, :melon].last(2) #=> ["grape", "melon"]
			// ```
			//
			// @param count [Integer]
			// @return [Object]
			Name: "last",
			Fn: func(receiver Object, sourceLine int) builtinMethodBody {
				return func(t *Thread, args []Object, blockFrame *normalCallFrame) Object {
					if len(args) > 1 {
						return t.vm.InitErrorObject(errors.ArgumentError, sourceLine, "Expect 0..1 argument. got=%d", len(args))
					}

					arr := receiver.(*ArrayObject)
					arrLength := len(arr.Elements)

					if len(args) == 0 {
						return arr.Elements[arrLength-1]
					}

					arg, ok := args[0].(*IntegerObject)

					if !ok {
						return t.vm.InitErrorObject(errors.TypeError, sourceLine, errors.WrongArgumentTypeFormat, classes.IntegerClass, args[0].Class().Name)
					}

					if arg.value < 1 {
						return t.vm.InitErrorObject(errors.ArgumentError, sourceLine, "Expect argument to be positive value. got=%d", arg.value)
					}

					if arrLength > arg.value {
						return t.vm.InitArrayObject(arr.Elements[arrLength-arg.value : arrLength])
					}
					return arr
				}
			},
		},
		{
			// Returns the length of the array.
			// The method does not take a block literal and is just to check the length of the array.
			//
			// ```ruby
			// [1, 2, 3].length #=> 3
			// ```
			//
			// @return [Integer]
			Name: "length",
			Fn: func(receiver Object, sourceLine int) builtinMethodBody {
				return func(t *Thread, args []Object, blockFrame *normalCallFrame) Object {

					if len(args) != 0 {
						return t.vm.InitErrorObject(errors.ArgumentError, sourceLine, "Expect 0 argument. got=%d", len(args))
					}

					arr := receiver.(*ArrayObject)
					return t.vm.InitIntegerObject(arr.length())
				}
			},
		},
		{
			// Loops through each element with the given block literal, and then returns the yielded elements as an array.
			// A block literal is required.
			//
			// ```ruby
			// a = ["a", "b", "c"]
			//
			// a.map do |e|
			//   e + e
			// end
			// #=> ["aa", "bb", "cc"]
			//
			// -------------------------
			//
			// a = [:apple, :orange, :lemon, :grape].map do |i|
			//   i + "s"
			// end
			// puts a
			// #=> ["apples", "oranges", "lemons", "grapes"]
			// ```
			//
			// @param block literal
			// @return [Array]
			Name: "map",
			Fn: func(receiver Object, sourceLine int) builtinMethodBody {
				return func(t *Thread, args []Object, blockFrame *normalCallFrame) Object {
					arr := receiver.(*ArrayObject)
					var elements = make([]Object, len(arr.Elements))

					if blockFrame == nil {
						return t.vm.InitErrorObject(errors.InternalError, sourceLine, errors.CantYieldWithoutBlockFormat)
					}

					// If it's an empty array, pop the block's call frame
					if len(arr.Elements) == 0 {
						t.callFrameStack.pop()
					}

					if blockIsEmpty(blockFrame) {
						for i := 0; i < len(arr.Elements); i++ {
							elements[i] = NULL
						}
					} else {
						for i, obj := range arr.Elements {
							result := t.builtinMethodYield(blockFrame, obj)
							elements[i] = result.Target
						}
					}

					return t.vm.InitArrayObject(elements)
				}
			},
		},
		{
			// A destructive method.
			// Removes the last element in the array and returns it.
			//
			// ```ruby
			// a = [1, 2, 3]
			// a.pop #=> 3
			// a     #=> [1, 2]
			// ```
			//
			// @return [Object]
			Name: "pop",
			Fn: func(receiver Object, sourceLine int) builtinMethodBody {
				return func(t *Thread, args []Object, blockFrame *normalCallFrame) Object {

					if len(args) != 0 {
						return t.vm.InitErrorObject(errors.ArgumentError, sourceLine, "Expect 0 argument. got=%d", len(args))
					}

					arr := receiver.(*ArrayObject)
					return arr.pop()
				}
			},
		},
		{
			// A destructive method.
			// Appends the given object to the array and returns the array.
			// One or more arguments can be passed to the method.
			// If no argument have been given, nothing will be added to the array,
			// and returns the unchanged array.
			// Even `nil` or empty strings `""` will be added to the array.
			//
			// ```ruby
			// a = [1, 2, 3]
			// a.push(4)       #=> [1, 2, 3, 4]
			// a.push(5, 6, 7) #=> [1, 2, 3, 4, 5, 6, 7]
			// a.push          #=> [1, 2, 3, 4, 5, 6, 7]
			// a               #=> [1, 2, 3, 4, 5, 6, 7]
			// a.push(nil, "") #=> [1, 2, 3, 4, 5, 6, 7, nil, ""]
			// ```
			//
			// @param object [Object]...
			// @return [Array]
			Name: "push",
			Fn: func(receiver Object, sourceLine int) builtinMethodBody {
				return func(t *Thread, args []Object, blockFrame *normalCallFrame) Object {

					arr := receiver.(*ArrayObject)
					return arr.push(args)
				}
			},
		},
		{
			// Accumulates the given argument and the results from evaluating each elements
			// with the first block parameter of the given block.
			// Takes one block with two block arguments (less than two block arguments are meaningless).
			// The first block argument is to succeed the initial value or previous result,
			// and the second block arguments is to enumerate the elements of the array.
			// You can also pass an argument as an initial value.
			// If you do not pass an argument, the first element of collection is used as an initial value.
			//
			// ```ruby
			// a = [1, 2, 7]
			//
			// a.reduce do |sum, n|
			//   sum + n
			// end
			// #=> 10
			//
			// a.reduce(10) do |sum, n|
			//   sum + n
			// end
			// #=> 20
			//
			// a = ["this", "is", "a", "test!"]
			// a.reduce("Yes, ") do |prev, s|
			//   prev + s + " "
			// end
			// #=> "Yes, this is a test! "
			// ```
			//
			// @param initial value [Object], block literal with two block parameters
			// @return [Object]
			Name: "reduce",
			Fn: func(receiver Object, sourceLine int) builtinMethodBody {
				return func(t *Thread, args []Object, blockFrame *normalCallFrame) Object {
					arr := receiver.(*ArrayObject)
					if blockFrame == nil {
						return t.vm.InitErrorObject(errors.InternalError, sourceLine, errors.CantYieldWithoutBlockFormat)
					}

					// If it's an empty array, pop the block's call frame
					if len(arr.Elements) == 0 {
						t.callFrameStack.pop()
					}

					if blockIsEmpty(blockFrame) {
						return NULL
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
						return t.vm.InitErrorObject(errors.ArgumentError, sourceLine, "Expect 0 or 1 argument. got=%d", len(args))
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
			// Returns a new array containing self‘s elements in reverse order. Not destructive.
			//
			// ```ruby
			// a = [1, 2, 7]
			//
			// a.reverse #=> [7, 2, 1]
			// ```
			//
			// @return [Array]
			Name: "reverse",
			Fn: func(receiver Object, sourceLine int) builtinMethodBody {
				return func(t *Thread, args []Object, blockFrame *normalCallFrame) Object {
					if len(args) != 0 {
						return t.vm.InitErrorObject(errors.ArgumentError, sourceLine, "Expect 0 argument. got=%d", len(args))
					}

					arr := receiver.(*ArrayObject)

					return arr.reverse()
				}
			},
		},
		{
			// Behaves as the same as #each, but traverses self in reverse order.
			// Returns self.
			// A block literal is required.
			//
			// ```ruby
			// a = [:a, :b, :c]
			//
			// a.reverse_each do |e|
			//   puts(e + e)
			// end
			// #=> "cc"
			// #=> "bb"
			// #=> "aa"
			// ```
			//
			// @param block literal
			// @return [Array]
			Name: "reverse_each",
			Fn: func(receiver Object, sourceLine int) builtinMethodBody {
				return func(t *Thread, args []Object, blockFrame *normalCallFrame) Object {
					if len(args) != 0 {
						return t.vm.InitErrorObject(errors.ArgumentError, sourceLine, "Expect 0 argument. got=%d", len(args))
					}

					if blockFrame == nil {
						return t.vm.InitErrorObject(errors.InternalError, sourceLine, errors.CantYieldWithoutBlockFormat)
					}

					arr := receiver.(*ArrayObject)
					if blockIsEmpty(blockFrame) {
						return arr
					}

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
			// Returns a new rotated array from the self.
			// The method is not destructive.
			// If zero `0` is passed, it returns a new array that has been rotated 1 time to left (default).
			// If an optional positive integer `n` is passed, it returns a new array that has been rotated `n` times to left.
			//
			// ```ruby
			// a = [:a, :b, :c, :d]
			//
			// a.rotate    #=> ["b", "c", "d", "a"]
			// a.rotate(2) #=> ["c", "d", "a", "b"]
			// a.rotate(3) #=> ["d", "a", "b", "c"]
			// ```
			//
			// If an optional negative integer `-n` is passed, it returns a new array that has been rotated `n` times to right.
			//
			// ```ruby
			// a = [:a, :b, :c, :d]
			//
			// a.rotate(-1) #=> ["d", "a", "b", "c"]
			// ```
			//
			// @param index [Integer]
			// @return [Array]
			Name: "rotate",
			Fn: func(receiver Object, sourceLine int) builtinMethodBody {
				return func(t *Thread, args []Object, blockFrame *normalCallFrame) Object {
					if len(args) > 1 {
						return t.vm.InitErrorObject(errors.ArgumentError, sourceLine, "Expect 0..1 argument. got=%d", len(args))
					}
					var rotate int
					arr := receiver.(*ArrayObject)
					rotArr := t.vm.InitArrayObject(arr.Elements)

					if len(args) == 0 {
						rotate = 1
					} else {
						arg, ok := args[0].(*IntegerObject)
						if !ok {
							return t.vm.InitErrorObject(errors.TypeError, sourceLine, errors.WrongArgumentTypeFormat, classes.IntegerClass, args[0].Class().Name)
						}
						rotate = arg.value
					}

					if rotate < 0 {
						for i := 0; i > rotate; i-- {
							el := rotArr.pop()
							rotArr.unshift([]Object{el})
						}
					} else {
						for i := 0; i < rotate; i++ {
							el := rotArr.shift()
							rotArr.push([]Object{el})
						}
					}

					return rotArr
				}
			},
		},
		{
			// Loops through each element with the given block literal that contains conditional expressions.
			// Returns a new array that contains elements that have been evaluated as `true` by the block.
			// A block literal is required.
			//
			// ```ruby
			// a = [1, 2, 3, 4, 5]
			//
			// a.select do |e|
			//   e + 1 > 3
			// end
			// #=> [3, 4, 5]
			// ```
			//
			// @param conditional block literal
			// @return [Array]
			Name: "select",
			Fn: func(receiver Object, sourceLine int) builtinMethodBody {
				return func(t *Thread, args []Object, blockFrame *normalCallFrame) Object {
					if len(args) > 0 {
						return t.vm.InitErrorObject(errors.ArgumentError, sourceLine, "Expect 0 argument. got=%d", len(args))
					}

					arr := receiver.(*ArrayObject)
					var elements []Object

					if blockFrame == nil {
						return t.vm.InitErrorObject(errors.InternalError, sourceLine, errors.CantYieldWithoutBlockFormat)
					}

					if blockIsEmpty(blockFrame) {
						return t.vm.InitArrayObject(elements)
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

					return t.vm.InitArrayObject(elements)
				}
			},
		},
		{
			// A destructive method.
			// Removes the first element from the array and returns the removed element.
			//
			// ```ruby
			// a = [1, 2, 3]
			// a.shift #=> 1
			// a       #=> [2, 3]
			// ```
			//
			// @return [Object]
			Name: "shift",
			Fn: func(receiver Object, sourceLine int) builtinMethodBody {
				return func(t *Thread, args []Object, blockFrame *normalCallFrame) Object {
					if len(args) != 0 {
						return t.vm.InitErrorObject(errors.ArgumentError, sourceLine, "Expect 0 argument. got=%d", len(args))
					}

					arr := receiver.(*ArrayObject)
					return arr.shift()
				}
			},
		},
		{
			// A destructive method.
			// Inserts one or more arguments at the first position of the array, and then returns the self.
			//
			// ```ruby
			// a = [1, 2]
			// a.unshift(0)             #=> [0, 1, 2]
			// a                        #=> [0, 1, 2]
			// a.unshift(:hello, :goby) #=> ["hello", "goby", 0, 1, 2]
			// a                        #=> ["hello", "goby", 0, 1, 2]
			// ```
			//
			// @param element [Object]
			// @return [Array]
			Name: "unshift",
			Fn: func(receiver Object, sourceLine int) builtinMethodBody {
				return func(t *Thread, args []Object, blockFrame *normalCallFrame) Object {
					arr := receiver.(*ArrayObject)
					return arr.unshift(args)
				}
			},
		},
		{
			// Returns a new array that contains the elements pointed by zero or more indices given.
			// If no arguments have been passed, an empty array `[]` will be returned.
			// If the index is out of range, `nil` is used as the element.
			//
			// ```ruby
			// a = ["a", "b", "c"]
			// a.values_at(1)     #=> ["b"]
			// a.values_at(-1, 3) #=> ["c", nil]
			// a.values_at()      #=> []
			// ```
			//
			// @param index [Integer]...
			// @return [Array]
			Name: "values_at",
			Fn: func(receiver Object, sourceLine int) builtinMethodBody {
				return func(t *Thread, args []Object, blockFrame *normalCallFrame) Object {
					arr := receiver.(*ArrayObject)
					var elements = make([]Object, len(args))

					for i, arg := range args {
						index, ok := arg.(*IntegerObject)

						if !ok {
							return t.vm.InitErrorObject(errors.TypeError, sourceLine, errors.WrongArgumentTypeFormat, classes.IntegerClass, arg.Class().Name)
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

					return t.vm.InitArrayObject(elements)
				}
			},
		},
	}
}

// Internal functions ===================================================

// Functions for initialization -----------------------------------------

// InitArrayObject returns a new object with the given elemnts
func (vm *VM) InitArrayObject(elements []Object) *ArrayObject {
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
	vm.libFiles = append(vm.libFiles, "array_enumerator.gb")
	vm.libFiles = append(vm.libFiles, "lazy_enumerator.gb")
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
func (a *ArrayObject) toJSON(t *Thread) string {
	var out bytes.Buffer
	elements := []string{}
	for _, e := range a.Elements {
		elements = append(elements, e.toJSON(t))
	}

	out.WriteString("[")
	out.WriteString(strings.Join(elements, ", "))
	out.WriteString("]")

	return out.String()
}

// concatenateCopies returns a array composed of N copies of the array
func (a *ArrayObject) concatenateCopies(t *Thread, n *IntegerObject) Object {
	aLen := len(a.Elements)
	result := make([]Object, 0, aLen*n.value)

	for i := 0; i < n.value; i++ {
		result = append(result, a.Elements...)
	}

	return t.vm.InitArrayObject(result)
}

// recursive indexed access - see ArrayObject#dig documentation.
func (a *ArrayObject) dig(t *Thread, keys []Object, sourceLine int) Object {
	currentKey := keys[0]
	intCurrentKey, ok := currentKey.(*IntegerObject)

	if !ok {
		return t.vm.InitErrorObject(errors.TypeError, sourceLine, errors.WrongArgumentTypeFormat, classes.IntegerClass, currentKey.Class().Name)
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
		return t.vm.InitErrorObject(errors.TypeError, sourceLine, "Expect target to be Diggable, got %s", currentValue.Class().Name)
	}

	return diggableCurrentValue.dig(t, nextKeys, sourceLine)
}

// Retrieves an object in an array using Integer index; common to `[]` and `at()`.
func (a *ArrayObject) index(t *Thread, args []Object, sourceLine int) Object {
	if len(args) > 2 || len(args) == 0 {
		return t.vm.InitErrorObject(errors.ArgumentError, sourceLine, "Expect 1..2 arguments. got=%d", len(args))
	}

	i := args[0]
	index, ok := i.(*IntegerObject)
	arrLength := a.length()

	if !ok {
		return t.vm.InitErrorObject(errors.TypeError, sourceLine, errors.WrongArgumentTypeFormat, classes.IntegerClass, args[0].Class().Name)
	}

	if index.value < 0 && index.value < -arrLength {
		return t.vm.InitErrorObject(errors.ArgumentError, sourceLine, "Index value %d too small for array. minimum: %d", index.value, -arrLength)
	}

	/* Validation for the second argument if exists */
	if len(args) == 2 {
		j := args[1]
		count, ok := j.(*IntegerObject)

		if !ok {
			return t.vm.InitErrorObject(errors.TypeError, sourceLine, errors.WrongArgumentTypeFormat, classes.IntegerClass, args[1].Class().Name)
		}
		if count.value < 0 {
			return t.vm.InitErrorObject(errors.ArgumentError, sourceLine, "Expect second argument greater than or equal 0. got: %d", count.value)
		}

		/*
		 *  This condition meets the special case (Don't know why ~ ? Ask Ruby or try it on irb!):
		 *
		 *  a = [1, 2, 3, 4, 5]
		 *  a[5, 5] #=> []
		 */
		if index.value > 0 && index.value == arrLength {
			return t.vm.InitArrayObject([]Object{})
		}
	}

	/* Start Indexing */
	normalizedIndex := a.normalizeIndex(index)
	if normalizedIndex == -1 {
		return NULL
	}

	if len(args) == 2 {
		j := args[1]
		count, _ := j.(*IntegerObject)

		if normalizedIndex+count.value > arrLength {
			return t.vm.InitArrayObject(a.Elements[normalizedIndex:])
		}
		return t.vm.InitArrayObject(a.Elements[normalizedIndex : normalizedIndex+count.value])
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
	e := make([]Object, len(a.Elements))

	copy(e, a.Elements)

	newArr := &ArrayObject{
		baseObj:  &baseObj{class: a.class},
		Elements: e,
	}

	return newArr
}

// unshift inserts an element in the first position of the array
func (a *ArrayObject) unshift(objs []Object) *ArrayObject {
	a.Elements = append(objs, a.Elements...)
	return a
}
