package vm

import (
	"bytes"
	"github.com/goby-lang/goby/vm/classes"
	"github.com/goby-lang/goby/vm/errors"
	"strings"
)

func (vm *VM) initArrayObject(elements []Object) *ArrayObject {
	return &ArrayObject{
		baseObj:  &baseObj{class: vm.topLevelClass(classes.ArrayClass)},
		Elements: elements,
	}
}

func (vm *VM) initArrayClass() *RClass {
	ac := vm.initializeClass(classes.ArrayClass, false)
	ac.setBuiltInMethods(builtinArrayInstanceMethods(), false)
	ac.setBuiltInMethods(builtInArrayClassMethods(), true)
	return ac
}

// ArrayObject represents instance from Array class.
// An array is a collection of different objects that are ordered and indexed.
// Elements in an array can belong to any class.
type ArrayObject struct {
	*baseObj
	Elements []Object
	splat    bool
}

func (a *ArrayObject) Value() interface{} {
	return a.Elements
}

// Polymorphic helper functions -----------------------------------------
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

// shift removes the first element in the array and returns it
func (a *ArrayObject) shift() Object {
	if len(a.Elements) < 1 {
		return NULL
	}

	value := a.Elements[0]
	a.Elements = a.Elements[1:]
	return value
}

func (a *ArrayObject) copy() Object {
	elems := make([]Object, len(a.Elements))

	copy(elems, a.Elements)

	newArr := &ArrayObject{
		baseObj:  &baseObj{class: a.class},
		Elements: elems,
	}

	return newArr
}

func builtInArrayClassMethods() []*BuiltInMethodObject {
	return []*BuiltInMethodObject{
		{
			Name: "new",
			Fn: func(receiver Object) builtinMethodBody {
				return func(t *thread, args []Object, blockFrame *callFrame) Object {
					return t.unsupportedMethodError("#new", receiver)
				}
			},
		},
	}
}

func builtinArrayInstanceMethods() []*BuiltInMethodObject {
	return []*BuiltInMethodObject{
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
			Fn: func(receiver Object) builtinMethodBody {
				return func(t *thread, args []Object, blockFrame *callFrame) Object {
					if len(args) != 1 {
						return t.vm.initErrorObject(errors.ArgumentError, "Expect 1 arguments. got=%d", len(args))
					}

					i := args[0]
					index, ok := i.(*IntegerObject)

					if !ok {
						return t.vm.initErrorObject(errors.TypeError, errors.WrongArgumentTypeFormat, classes.IntegerClass, args[0].Class().Name)
					}

					arr := receiver.(*ArrayObject)
					arrLength := len(arr.Elements)

					if int(index.value) < 0 {
						if -int(index.value) > arrLength {
							return NULL
						}
						calculatedIndex := arrLength + int(index.value)
						return arr.Elements[calculatedIndex]
					} else if int(index.value) >= arrLength {
						return NULL
					}

					return arr.Elements[index.value]
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
			Fn: func(receiver Object) builtinMethodBody {
				return func(t *thread, args []Object, blockFrame *callFrame) Object {

					// First arg is index
					// Second arg is assigned value
					if len(args) != 2 {
						return t.vm.initErrorObject(errors.ArgumentError, "Expect 2 arguments. got=%d", len(args))
					}

					i := args[0]
					index, ok := i.(*IntegerObject)
					indexValue := index.value

					if !ok {
						return t.vm.initErrorObject(errors.TypeError, errors.WrongArgumentTypeFormat, classes.IntegerClass, args[0].Class().Name)
					}

					arr := receiver.(*ArrayObject)

					// Negative index value condition
					if indexValue < 0 {
						if len(arr.Elements) < -indexValue {
							return t.vm.initErrorObject(errors.ArgumentError, "Index is too small for array. got=%s", i.Class().Name)
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
			// Retrieves an object in an array using the index argument.
			// It raises an error if index out of range.
			//
			// ```ruby
			// a = [1, 2, 3]
			// a.at(0)  # => 1
			// a.at(10) # => nil
			// a.at(-2) # => 2
			// a.at(-4) # => nil
			// ```
			Name: "at",
			Fn: func(receiver Object) builtinMethodBody {
				return func(t *thread, args []Object, blockFrame *callFrame) Object {
					if len(args) != 1 {
						return t.vm.initErrorObject(errors.ArgumentError, "Expect 1 argument. got=%d", len(args))
					}

					i := args[0]
					index, ok := i.(*IntegerObject)

					if !ok {
						return t.vm.initErrorObject(errors.TypeError, errors.WrongArgumentTypeFormat, classes.IntegerClass, args[0].Class().Name)
					}

					arr := receiver.(*ArrayObject)

					if index.value < 0 {
						if -index.value > len(arr.Elements) {
							return NULL
						}
						return arr.Elements[len(arr.Elements)+index.value]
					}

					if len(arr.Elements) == 0 || int(index.value) >= len(arr.Elements) {
						return NULL
					}

					return arr.Elements[index.value]
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
			Fn: func(receiver Object) builtinMethodBody {
				return func(t *thread, args []Object, blockFrame *callFrame) Object {
					if len(args) != 0 {
						return t.vm.initErrorObject(errors.ArgumentError, "Expect 0 argument. got=%d", len(args))
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
			Fn: func(receiver Object) builtinMethodBody {
				return func(t *thread, args []Object, blockFrame *callFrame) Object {
					arr := receiver.(*ArrayObject)

					for _, arg := range args {
						addAr, ok := arg.(*ArrayObject)

						if !ok {
							return t.vm.initErrorObject(errors.TypeError, errors.WrongArgumentTypeFormat, classes.ArrayClass, arg.Class().Name)
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
			Fn: func(receiver Object) builtinMethodBody {
				return func(t *thread, args []Object, blockFrame *callFrame) Object {
					arr := receiver.(*ArrayObject)
					var count int

					if len(args) > 1 {
						return t.vm.initErrorObject(errors.ArgumentError, "Expect 1 argument, got=%v", len(args))
					}

					if blockFrame != nil {
						for _, obj := range arr.Elements {
							result := t.builtInMethodYield(blockFrame, obj)
							if result.Target.(*BooleanObject).value {
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
			Fn: func(receiver Object) builtinMethodBody {
				return func(t *thread, args []Object, blockFrame *callFrame) Object {
					if len(args) != 0 {
						return t.vm.initErrorObject(errors.ArgumentError, "Expect 0 argument. got=%d", len(args))
					}

					if blockFrame == nil {
						return t.vm.initErrorObject(errors.InternalError, errors.CantYieldWithoutBlockFormat)
					}

					arr := receiver.(*ArrayObject)

					for _, obj := range arr.Elements {
						t.builtInMethodYield(blockFrame, obj)
					}
					return arr
				}
			},
		},
		{
			Name: "each_index",
			Fn: func(receiver Object) builtinMethodBody {
				return func(t *thread, args []Object, blockFrame *callFrame) Object {
					if len(args) != 0 {
						return t.vm.initErrorObject(errors.ArgumentError, "Expect 0 argument. got=%d", len(args))
					}

					if blockFrame == nil {
						return t.vm.initErrorObject(errors.InternalError, errors.CantYieldWithoutBlockFormat)
					}

					arr := receiver.(*ArrayObject)

					for i := range arr.Elements {
						t.builtInMethodYield(blockFrame, t.vm.initIntegerObject(i))
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
			Fn: func(receiver Object) builtinMethodBody {
				return func(t *thread, args []Object, blockFrame *callFrame) Object {

					if len(args) != 0 {
						return t.vm.initErrorObject(errors.ArgumentError, "Expect 0 argument. got=%d", len(args))
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
			Fn: func(receiver Object) builtinMethodBody {
				return func(t *thread, args []Object, blockFrame *callFrame) Object {
					if len(args) > 1 {
						return t.vm.initErrorObject(errors.ArgumentError, "Expect 0..1 argument. got=%d", len(args))
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
						return t.vm.initErrorObject(errors.TypeError, errors.WrongArgumentTypeFormat, classes.IntegerClass, args[0].Class().Name)
					}

					if arg.value < 1 {
						return t.vm.initErrorObject(errors.ArgumentError, "Expect argument to be positive value. got=%d", arg.value)
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
			Fn: func(receiver Object) builtinMethodBody {
				return func(t *thread, args []Object, blockFrame *callFrame) Object {
					arr := receiver.(*ArrayObject)

					if len(args) != 0 {
						return t.vm.initErrorObject(errors.ArgumentError, "Expect 0 argument. got=%d", len(args))
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
			Fn: func(receiver Object) builtinMethodBody {
				return func(t *thread, args []Object, blockFrame *callFrame) Object {
					arr := receiver.(*ArrayObject)

					var sep string
					if len(args) == 0 {
						sep = ""
					} else if len(args) == 1 {
						arg, ok := args[0].(*StringObject)

						if !ok {
							return t.vm.initErrorObject(errors.TypeError, errors.WrongArgumentTypeFormat, classes.StringClass, args[0].Class().Name)
						}

						sep = arg.value
					} else {
						return t.vm.initErrorObject(errors.ArgumentError, "Expect 0 or 1 argument. got=%d", len(args))
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
			Fn: func(receiver Object) builtinMethodBody {
				return func(t *thread, args []Object, blockFrame *callFrame) Object {
					if len(args) > 1 {
						return t.vm.initErrorObject(errors.ArgumentError, "Expect 0..1 argument. got=%d", len(args))
					}

					arr := receiver.(*ArrayObject)
					arrLength := len(arr.Elements)

					if len(args) == 0 {
						return arr.Elements[arrLength-1]
					}

					arg, ok := args[0].(*IntegerObject)

					if !ok {
						return t.vm.initErrorObject(errors.TypeError, errors.WrongArgumentTypeFormat, classes.IntegerClass, args[0].Class().Name)
					}

					if arg.value < 1 {
						return t.vm.initErrorObject(errors.ArgumentError, "Expect argument to be positive value. got=%d", arg.value)
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
			Fn: func(receiver Object) builtinMethodBody {
				return func(t *thread, args []Object, blockFrame *callFrame) Object {

					if len(args) != 0 {
						return t.vm.initErrorObject(errors.ArgumentError, "Expect 0 argument. got=%d", len(args))
					}

					arr := receiver.(*ArrayObject)
					return t.vm.initIntegerObject(arr.length())
				}
			},
		},
		{
			// Loop through each element with the given block. Return a new array with each yield element.
			//
			// ```ruby
			// a = ["a", "b", "c"]
			//
			// a.map do |e|
			//   e + e
			// end
			// # => ["aa", "bb", "cc"]
			// ```
			Name: "map",
			Fn: func(receiver Object) builtinMethodBody {
				return func(t *thread, args []Object, blockFrame *callFrame) Object {
					arr := receiver.(*ArrayObject)
					var elements = make([]Object, len(arr.Elements))

					if blockFrame == nil {
						return t.vm.initErrorObject(errors.InternalError, errors.CantYieldWithoutBlockFormat)
					}

					for i, obj := range arr.Elements {
						result := t.builtInMethodYield(blockFrame, obj)
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
			Fn: func(receiver Object) builtinMethodBody {
				return func(t *thread, args []Object, blockFrame *callFrame) Object {

					if len(args) != 0 {
						return t.vm.initErrorObject(errors.ArgumentError, "Expect 0 argument. got=%d", len(args))
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
			Fn: func(receiver Object) builtinMethodBody {
				return func(t *thread, args []Object, blockFrame *callFrame) Object {

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
			Fn: func(receiver Object) builtinMethodBody {
				return func(t *thread, args []Object, blockFrame *callFrame) Object {
					arr := receiver.(*ArrayObject)
					if blockFrame == nil {
						return t.vm.initErrorObject(errors.InternalError, errors.CantYieldWithoutBlockFormat)
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
						return t.vm.initErrorObject(errors.ArgumentError, "Expect 0 or 1 argument. got=%d", len(args))
					}

					for i := start; i < len(arr.Elements); i++ {
						result := t.builtInMethodYield(blockFrame, prev, arr.Elements[i])
						prev = result.Target
					}

					return prev
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
			Fn: func(receiver Object) builtinMethodBody {
				return func(t *thread, args []Object, blockFrame *callFrame) Object {
					if len(args) > 1 {
						return t.vm.initErrorObject(errors.ArgumentError, "Expect 0..1 argument. got=%d", len(args))
					}

					arr := receiver.(*ArrayObject)
					rotArr := t.vm.initArrayObject(arr.Elements)

					rotate := 1

					if len(args) != 0 {
						arg, ok := args[0].(*IntegerObject)
						if !ok {
							return t.vm.initErrorObject(errors.TypeError, errors.WrongArgumentTypeFormat, classes.IntegerClass, args[0].Class().Name)
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
			Fn: func(receiver Object) builtinMethodBody {
				return func(t *thread, args []Object, blockFrame *callFrame) Object {
					arr := receiver.(*ArrayObject)
					var elements []Object

					if blockFrame == nil {
						return t.vm.initErrorObject(errors.InternalError, errors.CantYieldWithoutBlockFormat)
					}

					for _, obj := range arr.Elements {
						result := t.builtInMethodYield(blockFrame, obj)
						if result.Target.(*BooleanObject).value {
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
			Fn: func(receiver Object) builtinMethodBody {
				return func(t *thread, args []Object, blockFrame *callFrame) Object {
					if len(args) != 0 {
						return t.vm.initErrorObject(errors.ArgumentError, "Expect 0 argument. got=%d", len(args))
					}

					arr := receiver.(*ArrayObject)
					return arr.shift()
				}
			},
		},
	}
}
