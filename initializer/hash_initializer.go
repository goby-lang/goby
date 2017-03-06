package initializer

import (
	"github.com/st0012/Rooby/object"
)

var (
	HashClass *object.HashClass
)

var builtinHashMethods = []*object.BuiltInMethod{
	{
		Fn: func(receiver object.Object) object.BuiltinMethodBody {
			return func(args ...object.Object) object.Object {
				if len(args) != 1 {
					return newError("Expect 1 arguments. got=%d", len(args))
				}

				i := args[0]
				key, ok := i.(*object.StringObject)

				if !ok {
					return newError("Expect index argument to be String. got=%T", i)
				}

				hash := receiver.(*object.HashObject)

				if len(hash.Pairs) == 0 {
					return NULL
				}

				value, ok := hash.Pairs[key.Value]

				if !ok {
					return NULL
				}

				return value

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

				k := args[0]
				key, ok := k.(*object.StringObject)

				if !ok {
					return newError("Expect index argument to be String. got=%T", k)
				}

				hash := receiver.(*object.HashObject)
				hash.Pairs[key.Value] = args[1]

				return args[1]
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

				hash := receiver.(*object.HashObject)
				return InitilaizeInteger(hash.Length())
			}
		},
		Name: "length",
	},
}

func initializeHashClass() *object.HashClass {
	methods := object.NewEnvironment()

	for _, m := range builtinHashMethods {
		methods.Set(m.Name, m)
	}

	bc := &object.BaseClass{Name: "Hash", Methods: methods, Class: ClassClass, SuperClass: ObjectClass}
	hc := &object.HashClass{BaseClass: bc}
	HashClass = hc
	return hc
}

func InitializeHash(pairs map[string]object.Object) *object.HashObject {
	return &object.HashObject{Pairs: pairs, Class: HashClass}
}
