package vm

import (
	"bytes"
	"strings"
)

var (
	arrayClass *RArray
)

// RArray is the built in array class
type RArray struct {
	*BaseClass
}

// ArrayObject represents instance from Array class.
// An array is a collection of different objects that are ordered and indexed.
// Elements in an array can belong to any class.
type ArrayObject struct {
	Class    *RArray
	Elements []Object
}

// objectType returns array instance's type
func (a *ArrayObject) objectType() objectType {
	return arrayObj
}

// Inspect returns detailed info of a array include elements it contains
func (a *ArrayObject) Inspect() string {
	var out bytes.Buffer

	elements := []string{}
	for _, e := range a.Elements {
		elements = append(elements, e.Inspect())
	}

	out.WriteString("Array:")
	out.WriteString("[")
	out.WriteString(strings.Join(elements, ", "))
	out.WriteString("]")

	return out.String()
}

// returnClass returns current object's class, which is RArray
func (a *ArrayObject) returnClass() Class {
	return a.Class
}

// length returns the length of array's elements
func (a *ArrayObject) length() int {
	return len(a.Elements)
}

// pop removes the last element in the array and returns it
func (a *ArrayObject) pop() Object {
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
	value := a.Elements[0]
	a.Elements = a.Elements[1:]
	return value
}

// initializeArray returns an array that contains given objects
func initializeArray(elements []Object) *ArrayObject {
	return &ArrayObject{Elements: elements, Class: arrayClass}
}

func init() {
	methods := newEnvironment()

	for _, m := range builtinArrayInstanceMethods {
		methods.set(m.Name, m)
	}

	bc := &BaseClass{Name: "Array", Methods: methods, ClassMethods: newEnvironment(), Class: classClass, pseudoSuperClass: objectClass, superClass: objectClass}
	ac := &RArray{BaseClass: bc}
	arrayClass = ac
}

var builtinArrayInstanceMethods = []*BuiltInMethodObject{
	{
		// Retrieves an object in an array using Integer index.
		// The index starts from 0. It returns `null` if the given index is bigger than its size.
		//
		// ```ruby
		// a = [1, 2, 3, "a", "b", "c"]
		// a[0]  # => 1
		// a[3]  # => "a"
		// a[10] # => null
		// ```
		Name: "[]",
		Fn: func(receiver Object) builtinMethodBody {
			return func(vm *VM, args []Object, blockFrame *callFrame) Object {
				if len(args) != 1 {
					return &Error{Message: "Expect 1 arguments. got=%d" + string(len(args))}
				}

				i := args[0]
				index, ok := i.(*IntegerObject)

				if !ok {
					return newError("Expect index argument to be Integer. got=%T", i)
				}

				arr := receiver.(*ArrayObject)

				if int(index.Value) >= len(arr.Elements) {
					return NULL
				}

				return arr.Elements[index.Value]

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
		// a          # => [10, null, null, 20]
		// ```
		Name: "[]=",
		Fn: func(receiver Object) builtinMethodBody {
			return func(vm *VM, args []Object, blockFrame *callFrame) Object {

				// First arg is index
				// Second arg is assigned value
				if len(args) != 2 {
					return newError("Expect 2 arguments. got=%d", len(args))
				}

				i := args[0]
				index, ok := i.(*IntegerObject)
				indexValue := index.Value

				if !ok {
					return newError("Expect index argument to be Integer. got=%T", i)
				}

				arr := receiver.(*ArrayObject)

				// Expand the array
				if len(arr.Elements) < (indexValue + 1) {
					newArr := make([]Object, indexValue+1)
					copy(newArr, arr.Elements)
					for i := range newArr[len(arr.Elements):] {
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
		// Returns the length of the array.
		//
		// ```ruby
		// [1, 2, 3].length # => 3
		// ```
		// @return [Integer]
		Name: "length",
		Fn: func(receiver Object) builtinMethodBody {
			return func(vm *VM, args []Object, blockFrame *callFrame) Object {

				if len(args) != 0 {
					return newError("Expect 0 argument. got=%d", len(args))
				}

				arr := receiver.(*ArrayObject)
				return initilaizeInteger(arr.length())
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
			return func(vm *VM, args []Object, blockFrame *callFrame) Object {

				if len(args) != 0 {
					return newError("Expect 0 argument. got=%d", len(args))
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
			return func(vm *VM, args []Object, blockFrame *callFrame) Object {

				arr := receiver.(*ArrayObject)
				return arr.push(args)
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
			return func(vm *VM, args []Object, blockFrame *callFrame) Object {
				if len(args) != 0 {
					return newError("Expect 0 argument. got=%d", len(args))
				}

				arr := receiver.(*ArrayObject)
				return arr.shift()
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
			return func(vm *VM, args []Object, blockFrame *callFrame) Object {
				arr := receiver.(*ArrayObject)

				if blockFrame == nil {
					vm.returnError("Can't yield without a block")
				}

				for _, obj := range arr.Elements {
					builtInMethodYield(vm, blockFrame, obj)
				}
				return arr
			}
		},
	},
	{
		Name: "each_index",
		Fn: func(receiver Object) builtinMethodBody {
			return func(vm *VM, args []Object, blockFrame *callFrame) Object {
				arr := receiver.(*ArrayObject)

				if blockFrame == nil {
					vm.returnError("Can't yield without a block")
				}

				for i := range arr.Elements {
					builtInMethodYield(vm, blockFrame, initilaizeInteger(i))
				}
				return arr
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
			return func(vm *VM, args []Object, blockFrame *callFrame) Object {
				arr := receiver.(*ArrayObject)
				var elements = make([]Object, len(arr.Elements))

				if blockFrame == nil {
					vm.returnError("Can't yield without a block")
				}

				for i, obj := range arr.Elements {
					result := builtInMethodYield(vm, blockFrame, obj)
					elements[i] = result.Target
				}

				return initializeArray(elements)
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
			return func(vm *VM, args []Object, blockFrame *callFrame) Object {
				arr := receiver.(*ArrayObject)
				var elements []Object

				if blockFrame == nil {
					vm.returnError("Can't yield without a block")
				}

				for _, obj := range arr.Elements {
					result := builtInMethodYield(vm, blockFrame, obj)
					if result.Target.(*BooleanObject).Value {
						elements = append(elements, obj)
					}
				}

				return initializeArray(elements)
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
		// a.at(10) # => Error
		// ```
		Name: "at",
		Fn: func(receiver Object) builtinMethodBody {
			return func(vm *VM, args []Object, blockFrame *callFrame) Object {
				i := args[0]
				index, ok := i.(*IntegerObject)

				if !ok {
					return newError("Expect index argument to be Integer. got=%T", i)
				}

				arr := receiver.(*ArrayObject)

				if len(arr.Elements) == 0 {
					return NULL
				}

				if int(index.Value) >= len(arr.Elements) {
					return newError("Index out of range")
				}

				return arr.Elements[index.Value]
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
			return func(vm *VM, args []Object, blockFrame *callFrame) Object {
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
			return func(vm *VM, args []Object, blockFrame *callFrame) Object {
				arr := receiver.(*ArrayObject)

				for _, arg := range args {
					addAr, ok := arg.(*ArrayObject)

					if !ok {
						return newError("Expect argument to be Array. got=%T", arg)
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
			return func(vm *VM, args []Object, blockFrame *callFrame) Object {
				arr := receiver.(*ArrayObject)
				var count int

				if blockFrame != nil {
					for _, obj := range arr.Elements {
						result := builtInMethodYield(vm, blockFrame, obj)
						if result.Target.(*BooleanObject).Value {
							count++
						}
					}

					return initilaizeInteger(count)
				}

				if len(args) > 1 {
					return newError("Expect one argument. got=%d", len(args))
				}

				if len(args) == 0 {
					return initilaizeInteger(len(arr.Elements))
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

				return initilaizeInteger(count)
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
			return func(vm *VM, args []Object, blockFrame *callFrame) Object {
				arr := receiver.(*ArrayObject)
				rotArr := initializeArray(arr.Elements)

				rotate := 1

				if len(args) != 0 {
					arg, ok := args[0].(*IntegerObject)
					if !ok {
						return newError("Expect index argument to be Integer. got=%T", args[0])
					}
					rotate = arg.Value
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
		// Returns the first element of the array.
		Name: "first",
		Fn: func(receiver Object) builtinMethodBody {
			return func(vm *VM, args []Object, blockFrame *callFrame) Object {
				arr := receiver.(*ArrayObject)

				if len(args) == 0 {
					return arr.Elements[0]
				}

				arg, ok := args[0].(*IntegerObject)
				if !ok {
					return newError("Expect index argument to be Integer. got=%T", args[0])
				}

				return initializeArray(arr.Elements[:arg.Value])
			}
		},
	},
	{
		// Returns the last element of the array.
		Name: "last",
		Fn: func(receiver Object) builtinMethodBody {
			return func(vm *VM, args []Object, blockFrame *callFrame) Object {
				arr := receiver.(*ArrayObject)

				if len(args) == 0 {
					return arr.Elements[len(arr.Elements)-1]
				}

				arg, ok := args[0].(*IntegerObject)
				if !ok {
					return newError("Expect index argument to be Integer. got=%T", args[0])
				}

				l := len(arr.Elements)
				return initializeArray(arr.Elements[l-arg.Value : l])
			}
		},
	},
}
