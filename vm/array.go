package vm

import (
	"bytes"
	"strings"

	"sort"

	"github.com/goby-lang/goby/vm/classes"
	"github.com/goby-lang/goby/vm/errors"
)

// ArrayObject represents an instance from Array class.
// An array is a collection of different objects that are ordered and indexed.
// Elements in an array can belong to any class and you can also build a "tuple" within an array.
// Array objects should always be enumerable.
type ArrayObject struct {
	*BaseObj
	Elements []Object
	splat    bool
}

// Class methods --------------------------------------------------------
var builtinArrayClassMethods = []*BuiltinMethodObject{
	{
		Name: "new",
		Fn: func(receiver Object, sourceLine int, t *Thread, args []Object, blockFrame *normalCallFrame) Object {
			if len(args) > 2 {
				return t.vm.InitErrorObject(errors.ArgumentError, sourceLine, errors.WrongNumberOfArgumentRange, 1, 2, len(args))
			}

			if len(args) >= 1 {
				n, ok := args[0].(*IntegerObject)

				if !ok {
					return t.vm.InitErrorObject(errors.ArgumentError, sourceLine, errors.WrongArgumentTypeFormat, "Integer", args[0].Class().Name)
				}

				if n.value < 0 {
					return t.vm.InitErrorObject(errors.ArgumentError, sourceLine, "Negative Array Size")
				}

				elems := make([]Object, n.value)

				if blockFrame != nil && !blockIsEmpty(blockFrame) {
					for i := range elems {
						elems[i] = t.builtinMethodYield(blockFrame, t.vm.InitIntegerObject(i))
					}
				} else {
					var elem Object

					if len(args) == 2 {
						elem = args[1]
					} else {
						elem = NULL
					}

					for i := 0; i < n.value; i++ {
						elems[i] = elem
					}
				}

				return t.vm.InitArrayObject(elems)
			}

			return t.vm.InitArrayObject([]Object{})
		},
	},
}

// Instance methods -----------------------------------------------------
var builtinArrayInstanceMethods = []*BuiltinMethodObject{
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
		// a[1, -1]  #=> ArgumentError: Expect second argument to be positive value. got: -1
		// a[-4, -3] #=> ArgumentError: Expect second argument to be positive value. got: -3
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
		Fn: func(receiver Object, sourceLine int, t *Thread, args []Object, blockFrame *normalCallFrame) Object {
			arr := receiver.(*ArrayObject)
			return arr.index(t, args, sourceLine)
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
		Fn: func(receiver Object, sourceLine int, t *Thread, args []Object, blockFrame *normalCallFrame) Object {
			if len(args) != 1 {
				return t.vm.InitErrorObject(errors.ArgumentError, sourceLine, errors.WrongNumberOfArgument, 1, len(args))
			}

			typeErr := t.vm.checkArgTypes(args, sourceLine, classes.IntegerClass)

			if typeErr != nil {
				return typeErr
			}

			return receiver.(*ArrayObject).concatenateCopies(t, args[0].Value().(int))
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
		Fn: func(receiver Object, sourceLine int, t *Thread, args []Object, blockFrame *normalCallFrame) Object {
			if len(args) != 1 {
				return t.vm.InitErrorObject(errors.ArgumentError, sourceLine, errors.WrongNumberOfArgument, 1, len(args))
			}

			otherArrayArg := args[0]
			otherArray, ok := otherArrayArg.(*ArrayObject)

			if !ok {
				return t.vm.InitErrorObject(errors.TypeError, sourceLine, errors.WrongArgumentTypeFormat, classes.ArrayClass, args[0].Class().Name)
			}

			selfArray := receiver.(*ArrayObject)

			newArrayElements := append(selfArray.Elements, otherArray.Elements...)

			newArray := t.vm.InitArrayObject(newArrayElements)

			return newArray
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
		// # ArgumentError: Expect second argument to be positive. got: -4
		// ```
		//
		// Note that passing multiple values to the method is unavailable.
		//
		// @param index [Integer], object [Object]
		// @param index [Integer], count [Integer], object [Object]
		// @return [Array]
		Name: "[]=",
		Fn: func(receiver Object, sourceLine int, t *Thread, args []Object, blockFrame *normalCallFrame) Object {

			// First argument is an index: there exists two cases which will be described in the following code
			aLen := len(args)
			if aLen < 2 || aLen > 3 {
				return t.vm.InitErrorObject(errors.ArgumentError, sourceLine, errors.WrongNumberOfArgumentRange, 2, 3, aLen)
			}

			typeErr := t.vm.checkArgTypes(args, sourceLine, classes.IntegerClass)

			if typeErr != nil {
				return typeErr
			}

			indexValue := args[0].Value().(int)
			arr := receiver.(*ArrayObject)

			// <Three Argument Case>
			// Second argument: the length of successive array values (zero or positive Integer)
			// Third argument: the assignment value (object)
			if aLen == 3 {
				// Negative index value too small
				if indexValue < 0 {
					if arr.normalizeIndex(indexValue) == -1 {
						return t.vm.InitErrorObject(errors.ArgumentError, sourceLine, errors.TooSmallIndexValue, indexValue, -arr.Len())
					}
					indexValue = arr.normalizeIndex(indexValue)
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
					return t.vm.InitErrorObject(errors.ArgumentError, sourceLine, errors.NegativeSecondValue, count.value)
				}

				a := args[2]
				assignedValue, isArray := a.(*ArrayObject)

				// Expand the array with nil; the second index is unnecessary in the case
				if indexValue >= arr.Len() {

					for arr.Len() < indexValue {
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
				if endValue > arr.Len() {
					endValue = arr.Len()
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
					return t.vm.InitErrorObject(errors.ArgumentError, sourceLine, errors.TooSmallIndexValue, indexValue, -arr.Len())
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
		Fn: func(receiver Object, sourceLine int, t *Thread, args []Object, blockFrame *normalCallFrame) Object {
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

				if result.isTruthy() {
					return TRUE
				}
			}

			return FALSE

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
		Fn: func(receiver Object, sourceLine int, t *Thread, args []Object, blockFrame *normalCallFrame) Object {
			if len(args) != 1 {
				return t.vm.InitErrorObject(errors.ArgumentError, sourceLine, errors.WrongNumberOfArgument, 1, len(args))
			}

			arr := receiver.(*ArrayObject)
			return arr.index(t, args, sourceLine)

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
		Fn: func(receiver Object, sourceLine int, t *Thread, args []Object, blockFrame *normalCallFrame) Object {
			if len(args) != 0 {
				return t.vm.InitErrorObject(errors.ArgumentError, sourceLine, errors.WrongNumberOfArgument, 0, len(args))
			}

			arr := receiver.(*ArrayObject)
			arr.Elements = []Object{}

			return arr
		},
	},
	{
		// Concatenation: returns a new array by just concatenating the arrays.
		// Empty or multiple arrays can be taken.
		//
		// ```ruby
		// a = [1, 2, 3]
		// a.concat([4, 5, 6])
		// a #=> [1, 2, 3, 4, 5, 6]
		//
		// [1, 2, 3].concat([])                 #=> [1, 2, 3]
		//
		// [1, 2, 3].concat([4, 5], [6, 7], []) #=> [1, 2, 3, 4, 5, 6, 7]
		// ```
		//
		// @param array [Array]
		// @return [Array]
		Name: "concat",
		Fn: func(receiver Object, sourceLine int, t *Thread, args []Object, blockFrame *normalCallFrame) Object {
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
		Fn: func(receiver Object, sourceLine int, t *Thread, args []Object, blockFrame *normalCallFrame) Object {
			aLen := len(args)
			if aLen > 1 {
				return t.vm.InitErrorObject(errors.ArgumentError, sourceLine, errors.WrongNumberOfArgumentLess, 1, aLen)
			}

			arr := receiver.(*ArrayObject)
			var count int
			if blockFrame != nil {
				if blockIsEmpty(blockFrame) {
					return t.vm.InitIntegerObject(0)
				}
				if len(arr.Elements) == 0 {
					t.callFrameStack.pop()
				}

				for _, obj := range arr.Elements {
					result := t.builtinMethodYield(blockFrame, obj)
					if result.isTruthy() {
						count++
					}
				}

				return t.vm.InitIntegerObject(count)
			}

			if aLen == 0 {
				return t.vm.InitIntegerObject(len(arr.Elements))
			}

			arg := args[0]
			findInt, findIsInt := arg.(*IntegerObject)
			findString, findIsString := arg.(*StringObject)
			findBoolean, findIsBoolean := arg.(*BooleanObject)

			for i := 0; i < len(arr.Elements); i++ {
				el := arr.Elements[i]
				switch el := el.(type) {
				case *IntegerObject:
					if findIsInt && findInt.equal(el) {
						count++
					}
				case *StringObject:
					if findIsString && findString.equal(el) {
						count++
					}
				case *BooleanObject:
					if findIsBoolean && findBoolean.equal(el) {
						count++
					}
				}
			}

			return t.vm.InitIntegerObject(count)

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
		Fn: func(receiver Object, sourceLine int, t *Thread, args []Object, blockFrame *normalCallFrame) Object {
			if len(args) != 1 {
				return t.vm.InitErrorObject(errors.ArgumentError, sourceLine, errors.WrongNumberOfArgument, 1, len(args))
			}

			typeErr := t.vm.checkArgTypes(args, sourceLine, classes.IntegerClass)

			if typeErr != nil {
				return typeErr
			}

			arr := receiver.(*ArrayObject)
			normalizedIndex := arr.normalizeIndex(args[0].Value().(int))

			if normalizedIndex == -1 {
				return NULL
			}

			// delete and slice

			deletedValue := arr.Elements[normalizedIndex]

			arr.Elements = append(arr.Elements[:normalizedIndex], arr.Elements[normalizedIndex+1:]...)

			return deletedValue

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
		Fn: func(receiver Object, sourceLine int, t *Thread, args []Object, blockFrame *normalCallFrame) Object {
			if len(args) < 1 {
				return t.vm.InitErrorObject(errors.ArgumentError, sourceLine, errors.WrongNumberOfArgumentMore, 1, len(args))
			}

			array := receiver.(*ArrayObject)
			value := array.dig(t, args, sourceLine)

			return value

		},
	},
	{
		// Performs a 'shallow' copy of the array and returns it.
		// Any arguments are ignored.
		// The object_id of the returned object is different from the one of the receiver.

		// Note that any elements of the array are NOT copied.
		//
		// See also `Object#dup`, `String#dup`, `Hash#dup`.
		//
		// ```ruby
		// a = ["s", "t", "r"]
		// a.object_id  #» 824635637568
		// a.each do |i|
		//   puts i.object_id
		// end
		// #» 824635637248
		// #» 824635637344
		// #» 824635637440
		//
		// b = a.dup
		// b.each do |i|
		//   puts i.object_id
		// end
		// #» 824635637248
		// #» 824635637344
		// #» 824635637440
		// b.object_id  #» 824637392704
		// ```
		//
		// @return [Array]
		Name: "dup",
		Fn: func(receiver Object, sourceLine int, t *Thread, args []Object, blockFrame *normalCallFrame) Object {
			arr, _ := receiver.(*ArrayObject)
			newArr := make([]Object, len(arr.Elements))
			copy(newArr, arr.Elements)
			newObj := t.vm.InitArrayObject(newArr)
			newObj.setInstanceVariables(arr.instanceVariables().copy())

			return newObj
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
		Fn: func(receiver Object, sourceLine int, t *Thread, args []Object, blockFrame *normalCallFrame) Object {
			if len(args) != 0 {
				return t.vm.InitErrorObject(errors.ArgumentError, sourceLine, errors.WrongNumberOfArgument, 0, len(args))
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
		Fn: func(receiver Object, sourceLine int, t *Thread, args []Object, blockFrame *normalCallFrame) Object {
			if len(args) != 0 {
				return t.vm.InitErrorObject(errors.ArgumentError, sourceLine, errors.WrongNumberOfArgument, 0, len(args))
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
		Fn: func(receiver Object, sourceLine int, t *Thread, args []Object, blockFrame *normalCallFrame) Object {
			if len(args) != 0 {
				return t.vm.InitErrorObject(errors.ArgumentError, sourceLine, errors.WrongNumberOfArgument, 0, len(args))
			}

			arr := receiver.(*ArrayObject)

			if arr.Len() == 0 {
				return TRUE
			}

			return FALSE

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
		Fn: func(receiver Object, sourceLine int, t *Thread, args []Object, blockFrame *normalCallFrame) Object {
			aLen := len(args)
			if aLen > 1 {
				return t.vm.InitErrorObject(errors.ArgumentError, sourceLine, errors.WrongNumberOfArgumentLess, 1, aLen)
			}

			arr := receiver.(*ArrayObject)
			arrLength := len(arr.Elements)
			if arrLength == 0 {
				return NULL
			}

			if aLen == 0 {
				return arr.Elements[0]
			}

			typeErr := t.vm.checkArgTypes(args, sourceLine, classes.IntegerClass)

			if typeErr != nil {
				return typeErr
			}

			value := args[0].Value().(int)

			if value < 1 {
				return t.vm.InitErrorObject(errors.ArgumentError, sourceLine, errors.NegativeValue, value)
			}

			if arrLength > value {
				return t.vm.InitArrayObject(arr.Elements[:value])
			}
			return arr

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
		Fn: func(receiver Object, sourceLine int, t *Thread, args []Object, blockFrame *normalCallFrame) Object {
			if len(args) != 0 {
				return t.vm.InitErrorObject(errors.ArgumentError, sourceLine, errors.WrongNumberOfArgument, 0, len(args))
			}

			arr := receiver.(*ArrayObject)
			newElements := arr.flatten()

			return t.vm.InitArrayObject(newElements)

		},
	},
	{
		// Returns a new hash from the element of the receiver (array) as keys, and generates respective values of hash from the keys by using the block provided.
		// The method can take a default value, and a block is required.
		// `index_with` is equivalent to `receiver.map do |e| e, e._do_something end.to_h`
		// Ref: https://github.com/rails/rails/pull/32523
		//
		// ```ruby
		// ary = [:Mon, :Tue, :Wed, :Thu, :Fri, :Sat, :Sun]
		// ary.index_with("weekday") do |d|
		//   if d == :Sat || d == :Sun
		//     "off day"
		//   end
		// end
		// #=> {Mon: "weekday",
		//      Tue: "weekday"
		//      Wed: "weekday"
		//      Thu: "weekday"
		//      Fri: "weekday"
		//      Sat: "off day"
		//      Sun: "off day"
		// }
		// ```
		//
		// @param optional default value [Object], block
		// @return [Hash]
		Name: "index_with",
		Fn: func(receiver Object, sourceLine int, t *Thread, args []Object, blockFrame *normalCallFrame) Object {
			if blockFrame == nil {
				return t.vm.InitErrorObject(errors.InternalError, sourceLine, errors.CantYieldWithoutBlockFormat)
			}

			a := receiver.(*ArrayObject)
			// If it's an empty array, pop the block's call frame
			if len(a.Elements) == 0 {
				t.callFrameStack.pop()
			}

			hash := make(map[string]Object)
			switch len(args) {
			case 0:
				for _, obj := range a.Elements {
					hash[obj.ToString()] = t.builtinMethodYield(blockFrame, obj)
				}
			case 1:
				arg := args[0]
				for _, obj := range a.Elements {
					switch b := t.builtinMethodYield(blockFrame, obj); b.(type) {
					case *NullObject:
						hash[obj.ToString()] = arg
					default:
						hash[obj.ToString()] = b
					}
				}
			default:
				return t.vm.InitErrorObject(errors.ArgumentError, sourceLine, errors.WrongNumberOfArgumentLess, 1, len(args))
			}

			return t.vm.InitHashObject(hash)

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
		Fn: func(receiver Object, sourceLine int, t *Thread, args []Object, blockFrame *normalCallFrame) Object {
			aLen := len(args)
			if aLen < 0 || aLen > 1 {
				return t.vm.InitErrorObject(errors.ArgumentError, sourceLine, errors.WrongNumberOfArgumentRange, 0, 1, aLen)
			}

			var sep string
			if aLen == 0 {
				sep = ""
			} else {
				typeErr := t.vm.checkArgTypes(args, sourceLine, classes.StringClass)

				if typeErr != nil {
					return typeErr
				}

				sep = args[0].Value().(string)
			}

			arr := receiver.(*ArrayObject)
			elements := []string{}
			for _, e := range arr.flatten() {
				elements = append(elements, e.ToString())
			}

			return t.vm.InitStringObject(strings.Join(elements, sep))

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
		Fn: func(receiver Object, sourceLine int, t *Thread, args []Object, blockFrame *normalCallFrame) Object {
			aLen := len(args)
			if aLen > 1 {
				return t.vm.InitErrorObject(errors.ArgumentError, sourceLine, errors.WrongNumberOfArgumentLess, 1, aLen)
			}

			arr := receiver.(*ArrayObject)
			arrLength := len(arr.Elements)

			if aLen == 0 {
				if arrLength == 0 {
					return NULL
				}

				return arr.Elements[arrLength-1]
			}

			typeErr := t.vm.checkArgTypes(args, sourceLine, classes.IntegerClass)

			if typeErr != nil {
				return typeErr
			}

			value := args[0].Value().(int)

			if value < 1 {
				return t.vm.InitErrorObject(errors.ArgumentError, sourceLine, errors.NegativeValue, value)
			}

			if arrLength > value {
				return t.vm.InitArrayObject(arr.Elements[arrLength-value : arrLength])
			}
			return arr

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
		Fn: func(receiver Object, sourceLine int, t *Thread, args []Object, blockFrame *normalCallFrame) Object {
			if len(args) != 0 {
				return t.vm.InitErrorObject(errors.ArgumentError, sourceLine, errors.WrongNumberOfArgument, 0, len(args))
			}

			arr := receiver.(*ArrayObject)
			return t.vm.InitIntegerObject(arr.Len())

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
		Fn: func(receiver Object, sourceLine int, t *Thread, args []Object, blockFrame *normalCallFrame) Object {
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
					elements[i] = t.builtinMethodYield(blockFrame, obj)
				}
			}

			return t.vm.InitArrayObject(elements)

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
		Fn: func(receiver Object, sourceLine int, t *Thread, args []Object, blockFrame *normalCallFrame) Object {
			if len(args) != 0 {
				return t.vm.InitErrorObject(errors.ArgumentError, sourceLine, errors.WrongNumberOfArgument, 0, len(args))
			}

			arr := receiver.(*ArrayObject)
			return arr.pop()

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
		Fn: func(receiver Object, sourceLine int, t *Thread, args []Object, blockFrame *normalCallFrame) Object {

			arr := receiver.(*ArrayObject)
			return arr.push(args)

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
		Fn: func(receiver Object, sourceLine int, t *Thread, args []Object, blockFrame *normalCallFrame) Object {
			aLen := len(args)
			if aLen > 1 {
				return t.vm.InitErrorObject(errors.ArgumentError, sourceLine, errors.WrongNumberOfArgumentLess, 1, aLen)
			}
			if blockFrame == nil {
				return t.vm.InitErrorObject(errors.InternalError, sourceLine, errors.CantYieldWithoutBlockFormat)
			}

			// If it's an empty array, pop the block's call frame
			arr := receiver.(*ArrayObject)
			if len(arr.Elements) == 0 {
				t.callFrameStack.pop()
			}

			if blockIsEmpty(blockFrame) {
				return NULL
			}

			var prev Object
			var start int
			switch aLen {
			case 0:
				prev = arr.Elements[0]
				start = 1
			case 1:
				prev = args[0]
				start = 0
			}

			for i := start; i < len(arr.Elements); i++ {
				prev = t.builtinMethodYield(blockFrame, prev, arr.Elements[i])
			}

			return prev

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
		Fn: func(receiver Object, sourceLine int, t *Thread, args []Object, blockFrame *normalCallFrame) Object {
			if len(args) != 0 {
				return t.vm.InitErrorObject(errors.ArgumentError, sourceLine, errors.WrongNumberOfArgument, 0, len(args))
			}

			arr := receiver.(*ArrayObject)
			return arr.reverse()

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
		Fn: func(receiver Object, sourceLine int, t *Thread, args []Object, blockFrame *normalCallFrame) Object {
			if len(args) != 0 {
				return t.vm.InitErrorObject(errors.ArgumentError, sourceLine, errors.WrongNumberOfArgument, 0, len(args))
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
		Fn: func(receiver Object, sourceLine int, t *Thread, args []Object, blockFrame *normalCallFrame) Object {
			aLen := len(args)
			if aLen > 1 {
				return t.vm.InitErrorObject(errors.ArgumentError, sourceLine, errors.WrongNumberOfArgumentLess, 1, aLen)
			}

			var rotate int
			arr := receiver.(*ArrayObject)
			rotArr := t.vm.InitArrayObject(arr.Elements)

			if aLen == 0 {
				rotate = 1
			} else {
				typeErr := t.vm.checkArgTypes(args, sourceLine, classes.IntegerClass)

				if typeErr != nil {
					return typeErr
				}

				rotate = args[0].Value().(int)
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
		Fn: func(receiver Object, sourceLine int, t *Thread, args []Object, blockFrame *normalCallFrame) Object {
			if len(args) != 0 {
				return t.vm.InitErrorObject(errors.ArgumentError, sourceLine, errors.WrongNumberOfArgument, 0, len(args))
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
				if result.isTruthy() {
					elements = append(elements, obj)
				}
			}

			return t.vm.InitArrayObject(elements)

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
		Fn: func(receiver Object, sourceLine int, t *Thread, args []Object, blockFrame *normalCallFrame) Object {
			if len(args) != 0 {
				return t.vm.InitErrorObject(errors.ArgumentError, sourceLine, errors.WrongNumberOfArgument, 0, len(args))
			}

			arr := receiver.(*ArrayObject)
			return arr.shift()

		},
	},
	{
		// Return a sorted array
		//
		// ```ruby
		// a = [3, 2, 1]
		// a.sort #=> [1, 2, 3]
		// ```
		//
		// @return [Object]
		Name: "sort",
		Fn: func(receiver Object, sourceLine int, t *Thread, args []Object, blockFrame *normalCallFrame) Object {
			if len(args) != 0 {
				return t.vm.InitErrorObject(errors.ArgumentError, sourceLine, "Expect 0 argument. got=%d", len(args))
			}

			arr := receiver.(*ArrayObject)
			newArr := arr.copy().(*ArrayObject)
			sort.Sort(newArr)
			return newArr

		},
	},
	{
		// Returns the result of interpreting ary as an array of [key value] array pairs.
		// Note that the keys should always be String or symbol literals (using symbol literal is preferable).
		// Each value can be any objects.
		//
		// ```ruby
		// ary = [[:john, [:guitar, :harmonica]], [:paul, :base], [:george, :guitar], [:ringo, :drum]]
		// ary.to_h
		// #=> { john: ["guitar", "harmonica"], paul: "base", george: "guitar", ringo: "drum" }
		// ```
		//
		// @return [Hash]
		Name: "to_h",
		Fn: func(receiver Object, sourceLine int, t *Thread, args []Object, blockFrame *normalCallFrame) Object {
			ary := receiver.(*ArrayObject)

			hash := make(map[string]Object)
			if len(ary.Elements) == 0 {
				return t.vm.InitHashObject(hash)
			}

			for i, el := range ary.Elements {
				kv, ok := el.(*ArrayObject)
				if !ok {
					return t.vm.InitErrorObject(errors.TypeError, sourceLine, "Expect the Array's element #%d to be Array. got: %s", i, el.Class().Name)
				}

				if len(kv.Elements) != 2 {
					return t.vm.InitErrorObject(errors.ArgumentError, sourceLine, "Expect element #%d to have 2 elements as a key-value pair. got: %s", i, kv.Inspect())
				}

				k := kv.Elements[0]
				if _, ok := k.(*StringObject); !ok {
					return t.vm.InitErrorObject(errors.TypeError, sourceLine, "Expect the key in the Array's element #%d to be String. got: %s", i, k.Class().Name)
				}

				hash[k.ToString()] = kv.Elements[1]

			}

			return t.vm.InitHashObject(hash)

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
		Fn: func(receiver Object, sourceLine int, t *Thread, args []Object, blockFrame *normalCallFrame) Object {
			arr := receiver.(*ArrayObject)
			return arr.unshift(args)

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
		Fn: func(receiver Object, sourceLine int, t *Thread, args []Object, blockFrame *normalCallFrame) Object {
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

		},
	},
}

// Internal functions ===================================================

// Functions for initialization -----------------------------------------

// InitArrayObject returns a new object with the given elemnts
func (vm *VM) InitArrayObject(elements []Object) *ArrayObject {
	return &ArrayObject{
		BaseObj:  NewBaseObject(vm.TopLevelClass(classes.ArrayClass)),
		Elements: elements,
	}
}

func (vm *VM) initArrayClass() *RClass {
	ac := vm.initializeClass(classes.ArrayClass)
	ac.setBuiltinMethods(builtinArrayInstanceMethods, false)
	ac.setBuiltinMethods(builtinArrayClassMethods, true)
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

// ToString returns the object's elements as the string format
func (a *ArrayObject) ToString() string {
	var out bytes.Buffer

	elements := []string{}
	for _, e := range a.Elements {
		elements = append(elements, e.Inspect())
	}

	out.WriteString("[")
	out.WriteString(strings.Join(elements, ", "))
	out.WriteString("]")

	return out.String()
}

// Inspect delegates to ToString
func (a *ArrayObject) Inspect() string {
	return a.ToString()
}

// ToJSON returns the object's elements as the JSON string format
func (a *ArrayObject) ToJSON(t *Thread) string {
	var out bytes.Buffer
	elements := []string{}
	for _, e := range a.Elements {
		elements = append(elements, e.ToJSON(t))
	}

	out.WriteString("[")
	out.WriteString(strings.Join(elements, ", "))
	out.WriteString("]")

	return out.String()
}

// concatenateCopies returns a array composed of N copies of the array
func (a *ArrayObject) concatenateCopies(t *Thread, n int) Object {
	aLen := len(a.Elements)
	result := make([]Object, 0, aLen*n)

	for i := 0; i < n; i++ {
		result = append(result, a.Elements...)
	}

	return t.vm.InitArrayObject(result)
}

// recursive indexed access - see ArrayObject#dig documentation.
func (a *ArrayObject) dig(t *Thread, keys []Object, sourceLine int) Object {
	typeErr := t.vm.checkArgTypes(keys, sourceLine, classes.IntegerClass)

	if typeErr != nil {
		return typeErr
	}

	normalizedIndex := a.normalizeIndex(keys[0].Value().(int))

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
		return t.vm.InitErrorObject(errors.TypeError, sourceLine, errors.NotDiggable, currentValue.Class().Name)
	}

	return diggableCurrentValue.dig(t, nextKeys, sourceLine)
}

// Retrieves an object in an array using Integer index; common to `[]` and `at()`.
func (a *ArrayObject) index(t *Thread, args []Object, sourceLine int) Object {
	aLen := len(args)
	if aLen < 1 || aLen > 2 {
		return t.vm.InitErrorObject(errors.ArgumentError, sourceLine, errors.WrongNumberOfArgumentRange, 1, 2, aLen)
	}

	typeErr := t.vm.checkArgTypes(args, sourceLine, classes.IntegerClass)

	if typeErr != nil {
		return typeErr
	}

	index := args[0].Value().(int)
	arrLength := a.Len()

	if index < 0 && index < -arrLength {
		return t.vm.InitErrorObject(errors.ArgumentError, sourceLine, errors.TooSmallIndexValue, index, -arrLength)
	}

	/* Validation for the second argument if exists */
	if aLen == 2 {
		j := args[1]
		count, ok := j.(*IntegerObject)

		if !ok {
			return t.vm.InitErrorObject(errors.TypeError, sourceLine, errors.WrongArgumentTypeFormat, classes.IntegerClass, args[1].Class().Name)
		}
		if count.value < 0 {
			return t.vm.InitErrorObject(errors.ArgumentError, sourceLine, errors.NegativeSecondValue, count.value)
		}

		/*
		 *  This condition meets the special case (Don't know why ~ ? Ask Ruby or try it on irb!):
		 *
		 *  a = [1, 2, 3, 4, 5]
		 *  a[5, 5] #=> []
		 */
		if index > 0 && index == arrLength {
			return t.vm.InitArrayObject([]Object{})
		}
	}

	/* Start Indexing */
	normalizedIndex := a.normalizeIndex(index)
	if normalizedIndex == -1 {
		return NULL
	}

	if aLen == 2 {
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

// Len returns the length of array's elements
func (a *ArrayObject) Len() int {
	return len(a.Elements)
}

// Swap is one of the required method to fulfill sortable interface
func (a *ArrayObject) Swap(i, j int) {
	a.Elements[i], a.Elements[j] = a.Elements[j], a.Elements[i]
}

// Less is one of the required method to fulfill sortable interface
func (a *ArrayObject) Less(i, j int) bool {
	leftObj, rightObj := a.Elements[i], a.Elements[j]
	switch leftObj := leftObj.(type) {
	case Numeric:
		return leftObj.lessThan(rightObj)
	case *StringObject:
		right, ok := rightObj.(*StringObject)

		if ok {
			return leftObj.value < right.value
		}

		return false
	default:
		return false
	}
}

// normalizes the index to the Ruby-style:
//
// 1. if the index is between o and the index length, returns the index
// 2. if it's a negative value (within bounds), returns the normalized positive version
// 3. if it's out of bounds (either positive or negative), returns -1
func (a *ArrayObject) normalizeIndex(index int) int {
	aLength := len(a.Elements)

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

	return &ArrayObject{
		BaseObj:  NewBaseObject(a.class),
		Elements: reversedArrElems,
	}
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

	return &ArrayObject{
		BaseObj:  NewBaseObject(a.class),
		Elements: e,
	}
}

func (a *ArrayObject) equalTo(compared Object) bool {
	c, ok := compared.(*ArrayObject)

	if !ok {
		return false
	}

	if len(a.Elements) != len(c.Elements) {
		return false
	}

	for i, e := range a.Elements {
		if !e.equalTo(c.Elements[i]) {
			return false
		}
	}

	return true
}

// unshift inserts an element in the first position of the array
func (a *ArrayObject) unshift(objs []Object) *ArrayObject {
	a.Elements = append(objs, a.Elements...)
	return a
}
