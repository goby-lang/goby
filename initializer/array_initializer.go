package initializer

import "github.com/st0012/Rooby/object"

var (
	ArrayClass *object.ArrayClass
)

var builtinArrayMethods = []*object.BuiltInMethod{
	{
		Fn: func(receiver object.Object) object.BuiltinMethodBody {
			return func(args ...object.Object) object.Object {
				if len(args) < 1 {
					return newError("Too few arguments. expected 1, got=%d", len(args))
				}

				if len(args) > 1 {
					return newError("Too many arguments. expected 1, got=%d", len(args))
				}

				i := args[0]
				index, ok := i.(*object.IntegerObject)

				if !ok {
					return newError("Expect argument to be Integer. got=%T", i)
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
