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

// ArrayObject represents array instance
type ArrayObject struct {
	Class    *RArray
	Elements []Object
}

// objectType returns array instance's type
func (a *ArrayObject) objectType() objectType {
	return arrayObj
}

// inspect returns detailed info of a array include elements it contains
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

// initializeArray returns an array that contains given objects
func initializeArray(elements []Object) *ArrayObject {
	return &ArrayObject{Elements: elements, Class: arrayClass}
}

func init() {
	methods := NewEnvironment()

	for _, m := range builtinArrayMethods {
		methods.Set(m.Name, m)
	}

	bc := &BaseClass{Name: "Array", Methods: methods, ClassMethods: NewEnvironment(), Class: classClass, SuperClass: objectClass}
	ac := &RArray{BaseClass: bc}
	arrayClass = ac
}

var builtinArrayMethods = []*BuiltInMethod{
	{
		Fn: func(receiver Object) builtinMethodBody {
			return func(vm *VM, args []Object, blockFrame *callFrame) Object {
				if len(args) != 1 {
					return newError("Expect 1 arguments. got=%d", len(args))
				}

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
		Name: "[]",
	},
	{
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
		Name: "[]=",
	},
	{
		Fn: func(receiver Object) builtinMethodBody {
			return func(vm *VM, args []Object, blockFrame *callFrame) Object {

				if len(args) != 0 {
					return newError("Expect 0 argument. got=%d", len(args))
				}

				arr := receiver.(*ArrayObject)
				return initilaizeInteger(arr.length())
			}
		},
		Name: "length",
	},
	{
		Fn: func(receiver Object) builtinMethodBody {
			return func(vm *VM, args []Object, blockFrame *callFrame) Object {

				if len(args) != 0 {
					return newError("Expect 0 argument. got=%d", len(args))
				}

				arr := receiver.(*ArrayObject)
				return arr.pop()
			}
		},
		Name: "pop",
	},
	{
		Fn: func(receiver Object) builtinMethodBody {
			return func(vm *VM, args []Object, blockFrame *callFrame) Object {

				arr := receiver.(*ArrayObject)
				return arr.push(args)
			}
		},
		Name: "push",
	},
	{
		Fn: func(receiver Object) builtinMethodBody {
			return func(vm *VM, args []Object, blockFrame *callFrame) Object {
				arr := receiver.(*ArrayObject)

				if blockFrame == nil {
					panic("Can't yield without a block")
				}

				for _, obj := range arr.Elements {
					builtInMethodYield(vm, blockFrame, obj)
				}
				return arr
			}
		},
		Name: "each",
	},
}
