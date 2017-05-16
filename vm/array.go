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

func initArray() {
	methods := newEnvironment()

	for _, m := range builtinArrayMethods {
		methods.set(m.Name, m)
	}

	bc := &BaseClass{Name: "Array", Methods: methods, ClassMethods: newEnvironment(), Class: classClass, pseudoSuperClass: objectClass, superClass: objectClass}
	ac := &RArray{BaseClass: bc}
	arrayClass = ac
}

var builtinArrayMethods = []*BuiltInMethod{
	{
		Name: "[]",
		Fn: func(receiver Object) builtinMethodBody {
			return func(ma methodArgs) Object {
				if len(ma.args) != 1 {
					return &Error{Message: "Expect 1 arguments. got=%d" + string(len(ma.args))}
				}

				i := ma.args[0]
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
		Name: "[]=",
		Fn: func(receiver Object) builtinMethodBody {
			return func(ma methodArgs) Object {

				// First arg is index
				// Second arg is assigned value
				if len(ma.args) != 2 {
					return newError("Expect 2 arguments. got=%d", len(ma.args))
				}

				i := ma.args[0]
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

				arr.Elements[indexValue] = ma.args[1]

				return arr.Elements[indexValue]
			}
		},
	},
	{
		Name: "length",
		Fn: func(receiver Object) builtinMethodBody {
			return func(ma methodArgs) Object {

				if len(ma.args) != 0 {
					return newError("Expect 0 argument. got=%d", len(ma.args))
				}

				arr := receiver.(*ArrayObject)
				return initilaizeInteger(arr.length())
			}
		},
	},
	{
		Name: "pop",
		Fn: func(receiver Object) builtinMethodBody {
			return func(ma methodArgs) Object {

				if len(ma.args) != 0 {
					return newError("Expect 0 argument. got=%d", len(ma.args))
				}

				arr := receiver.(*ArrayObject)
				return arr.pop()
			}
		},
	},
	{
		Name: "push",
		Fn: func(receiver Object) builtinMethodBody {
			return func(ma methodArgs) Object {

				arr := receiver.(*ArrayObject)
				return arr.push(ma.args)
			}
		},
	},
	{
		Name: "shift",
		Fn: func(receiver Object) builtinMethodBody {
			return func(ma methodArgs) Object {
				if len(ma.args) != 0 {
					return newError("Expect 0 argument. got=%d", len(ma.args))
				}

				arr := receiver.(*ArrayObject)
				return arr.shift()
			}
		},
	},
	{
		Name: "each",
		Fn: func(receiver Object) builtinMethodBody {
			return func(ma methodArgs) Object {
				arr := receiver.(*ArrayObject)

				if ma.blockFrame == nil {
					panic("Can't yield without a block")
				}

				for _, obj := range arr.Elements {
					builtInMethodYield(ma.vm, ma.blockFrame, obj)
				}
				return arr
			}
		},
	},
	{
		Name: "each_index",
		Fn: func(receiver Object) builtinMethodBody {
			return func(ma methodArgs) Object {
				arr := receiver.(*ArrayObject)

				if ma.blockFrame == nil {
					panic("Can't yield without a block")
				}

				for i := range arr.Elements {
					builtInMethodYield(ma.vm, ma.blockFrame, initilaizeInteger(i))
				}
				return arr
			}
		},
	},
	{
		Name: "map",
		Fn: func(receiver Object) builtinMethodBody {
			return func(ma methodArgs) Object {
				arr := receiver.(*ArrayObject)
				var elements = make([]Object, len(arr.Elements))

				if ma.blockFrame == nil {
					panic("Can't yield without a block")
				}

				for i, obj := range arr.Elements {
					result := builtInMethodYield(ma.vm, ma.blockFrame, obj)
					elements[i] = result.Target
				}

				return initializeArray(elements)
			}
		},
	},
	{
		Name: "select",
		Fn: func(receiver Object) builtinMethodBody {
			return func(ma methodArgs) Object {
				arr := receiver.(*ArrayObject)
				var elements []Object

				if ma.blockFrame == nil {
					panic("Can't yield without a block")
				}

				for _, obj := range arr.Elements {
					result := builtInMethodYield(ma.vm, ma.blockFrame, obj)
					if result.Target.(*BooleanObject).Value {
						elements = append(elements, obj)
					}
				}

				return initializeArray(elements)
			}
		},
	},
	{
		Name: "at",
		Fn: func(receiver Object) builtinMethodBody {
			return func(ma methodArgs) Object {
				i := ma.args[0]
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
		Name: "clear",
		Fn: func(receiver Object) builtinMethodBody {
			return func(ma methodArgs) Object {
				arr := receiver.(*ArrayObject)
				arr.Elements = []Object{}

				return arr
			}
		},
	},
	{
		Name: "concat",
		Fn: func(receiver Object) builtinMethodBody {
			return func(ma methodArgs) Object {
				arr := receiver.(*ArrayObject)

				for _, arg := range ma.args {
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
		Name: "count",
		Fn: func(receiver Object) builtinMethodBody {
			return func(ma methodArgs) Object {
				arr := receiver.(*ArrayObject)
				var count int

				if ma.blockFrame != nil {
					for _, obj := range arr.Elements {
						result := builtInMethodYield(ma.vm, ma.blockFrame, obj)
						if result.Target.(*BooleanObject).Value {
							count++
						}
					}

					return initilaizeInteger(count)
				}

				if len(ma.args) > 1 {
					return newError("Expect one argument. got=%d", len(ma.args))
				}

				if len(ma.args) == 0 {
					return initilaizeInteger(len(arr.Elements))
				}

				arg := ma.args[0]
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
		Name: "rotate",
		Fn: func(receiver Object) builtinMethodBody {
			return func(ma methodArgs) Object {
				arr := receiver.(*ArrayObject)
				rotArr := initializeArray(arr.Elements)

				rotate := 1

				if len(ma.args) != 0 {
					arg, ok := ma.args[0].(*IntegerObject)
					if !ok {
						return newError("Expect index argument to be Integer. got=%T", ma.args[0])
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
		Name: "first",
		Fn: func(receiver Object) builtinMethodBody {
			return func(ma methodArgs) Object {
				arr := receiver.(*ArrayObject)

				if len(ma.args) == 0 {
					return arr.Elements[0]
				}

				arg, ok := ma.args[0].(*IntegerObject)
				if !ok {
					return newError("Expect index argument to be Integer. got=%T", ma.args[0])
				}

				return initializeArray(arr.Elements[:arg.Value])
			}
		},
	},
	{
		Name: "last",
		Fn: func(receiver Object) builtinMethodBody {
			return func(ma methodArgs) Object {
				arr := receiver.(*ArrayObject)

				if len(ma.args) == 0 {
					return arr.Elements[len(arr.Elements)-1]
				}

				arg, ok := ma.args[0].(*IntegerObject)
				if !ok {
					return newError("Expect index argument to be Integer. got=%T", ma.args[0])
				}

				l := len(arr.Elements)
				return initializeArray(arr.Elements[l-arg.Value : l])
			}
		},
	},
}
