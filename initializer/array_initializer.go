package initializer

import (
	"github.com/st0012/Rooby/object"
)

var (
	ArrayClass *object.ArrayClass
)

var builtinArrayMethods = []*object.BuiltInMethod{
	{
		Fn: func(receiver object.Object) object.BuiltinMethodBody {
			return func(args ...object.Object) object.Object {
				if len(args) != 1 {
					return newError("Expect 1 arguments. got=%d", len(args))
				}

				i := args[0]
				index, ok := i.(*object.IntegerObject)

				if !ok {
					return newError("Expect index argument to be Integer. got=%T", i)
				}

				arr := receiver.(*object.ArrayObject)

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
		Fn: func(receiver object.Object) object.BuiltinMethodBody {
			return func(args ...object.Object) object.Object {
				// First arg is index
				// Second arg is assigned value
				if len(args) != 2 {
					return newError("Expect 2 arguments. got=%d", len(args))
				}

				i := args[0]
				index, ok := i.(*object.IntegerObject)
				indexValue := index.Value

				if !ok {
					return newError("Expect index argument to be Integer. got=%T", i)
				}

				arr := receiver.(*object.ArrayObject)

				// Expand the array
				if len(arr.Elements) < (indexValue + 1) {
					newArr := make([]object.Object, indexValue+1)
					copy(newArr, arr.Elements)
					for i, _ := range newArr[len(arr.Elements):] {
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
		Fn: func(receiver object.Object) object.BuiltinMethodBody {
			return func(args ...object.Object) object.Object {
				if len(args) != 0 {
					return newError("Expect 0 argument. got=%d", len(args))
				}

				arr := receiver.(*object.ArrayObject)
				return InitilaizeInteger(arr.Length())
			}
		},
		Name: "length",
	},
	{
		Fn: func(receiver object.Object) object.BuiltinMethodBody {
			return func(args ...object.Object) object.Object {
				if len(args) != 0 {
					return newError("Expect 0 argument. got=%d", len(args))
				}

				arr := receiver.(*object.ArrayObject)
				return arr.Pop()
			}
		},
		Name: "pop",
	},
	{
		Fn: func(receiver object.Object) object.BuiltinMethodBody {
			return func(args ...object.Object) object.Object {
				arr := receiver.(*object.ArrayObject)
				return arr.Push(args)
			}
		},
		Name: "push",
	},
}

func initializeArrayClass() *object.ArrayClass {
	methods := object.NewEnvironment()

	for _, m := range builtinArrayMethods {
		methods.Set(m.Name, m)
	}

	bc := &object.BaseClass{Name: "Array", Methods: methods, Class: ClassClass, SuperClass: ObjectClass}
	ac := &object.ArrayClass{BaseClass: bc}
	ArrayClass = ac
	return ac
}

func InitializeArray(elements []object.Object) *object.ArrayObject {
	return &object.ArrayObject{Elements: elements, Class: ArrayClass}
}
