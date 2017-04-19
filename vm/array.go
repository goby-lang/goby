package vm

import (
	"bytes"
	"strings"
)

var (
	ArrayClass *RArray
)

type RArray struct {
	*BaseClass
}

type ArrayObject struct {
	Class    *RArray
	Elements []Object
}

func (a *ArrayObject) Type() ObjectType {
	return ARRAY_OBJ
}

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

func (a *ArrayObject) ReturnClass() Class {
	return a.Class
}

func (a *ArrayObject) Length() int {
	return len(a.Elements)
}

func (a *ArrayObject) Pop() Object {
	value := a.Elements[len(a.Elements)-1]
	a.Elements = a.Elements[:len(a.Elements)-1]
	return value
}

func (a *ArrayObject) Push(objs []Object) *ArrayObject {
	a.Elements = append(a.Elements, objs...)
	return a
}

func InitializeArray(elements []Object) *ArrayObject {
	return &ArrayObject{Elements: elements, Class: ArrayClass}
}

func init() {
	methods := NewEnvironment()

	for _, m := range builtinArrayMethods {
		methods.Set(m.Name, m)
	}

	bc := &BaseClass{Name: "Array", Methods: methods, ClassMethods: NewEnvironment(), Class: ClassClass, SuperClass: ObjectClass}
	ac := &RArray{BaseClass: bc}
	ArrayClass = ac
}

var builtinArrayMethods = []*BuiltInMethod{
	{
		Fn: func(receiver Object) BuiltinMethodBody {
			return func(args []Object, block *Method) Object {
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
		Fn: func(receiver Object) BuiltinMethodBody {
			return func(args []Object, block *Method) Object {
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
		Fn: func(receiver Object) BuiltinMethodBody {
			return func(args []Object, block *Method) Object {
				if len(args) != 0 {
					return newError("Expect 0 argument. got=%d", len(args))
				}

				arr := receiver.(*ArrayObject)
				return InitilaizeInteger(arr.Length())
			}
		},
		Name: "length",
	},
	{
		Fn: func(receiver Object) BuiltinMethodBody {
			return func(args []Object, block *Method) Object {
				if len(args) != 0 {
					return newError("Expect 0 argument. got=%d", len(args))
				}

				arr := receiver.(*ArrayObject)
				return arr.Pop()
			}
		},
		Name: "pop",
	},
	{
		Fn: func(receiver Object) BuiltinMethodBody {
			return func(args []Object, block *Method) Object {
				arr := receiver.(*ArrayObject)
				return arr.Push(args)
			}
		},
		Name: "push",
	},
	//{
	//	Fn: func(receiver Object) BuiltinMethodBody {
	//		return func(args []Object, block *Method) Object {
	//			arr := receiver.(*ArrayObject)
	//			for _, obj := range arr.Elements {
	//				evalMethodObject(block.Scope.Self, block, []Object{obj}, nil)
	//			}
	//			return arr
	//		}
	//	},
	//	Name: "each",
	//},
}
