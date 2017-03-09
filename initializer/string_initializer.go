package initializer

import (
	"github.com/st0012/Rooby/object"
)

var builtinStringMethods = []*object.BuiltInMethod{
	{
		Fn: func(receiver object.Object) object.BuiltinMethodBody {
			return func(args ...object.Object) object.Object {
				err := checkArgumentLen(args, object.StringClass, "+")

				if err != nil {
					return err
				}

				leftValue := receiver.(*object.StringObject).Value
				right, ok := args[0].(*object.StringObject)

				if !ok {
					return wrongTypeError(object.StringClass)
				}

				rightValue := right.Value
				return &object.StringObject{Value: leftValue + rightValue, Class: object.StringClass}
			}
		},
		Name: "+",
	},
	{
		Fn: func(receiver object.Object) object.BuiltinMethodBody {
			return func(args ...object.Object) object.Object {
				err := checkArgumentLen(args, object.StringClass, ">")
				if err != nil {
					return err
				}

				leftValue := receiver.(*object.StringObject).Value
				right, ok := args[0].(*object.StringObject)

				if !ok {
					return wrongTypeError(object.StringClass)
				}

				rightValue := right.Value

				if leftValue > rightValue {
					return object.TRUE
				}

				return object.FALSE
			}
		},
		Name: ">",
	},
	{
		Fn: func(receiver object.Object) object.BuiltinMethodBody {
			return func(args ...object.Object) object.Object {
				err := checkArgumentLen(args, object.StringClass, "<")
				if err != nil {
					return err
				}

				leftValue := receiver.(*object.StringObject).Value
				right, ok := args[0].(*object.StringObject)

				if !ok {
					return wrongTypeError(object.StringClass)
				}

				rightValue := right.Value

				if leftValue < rightValue {
					return object.TRUE
				}

				return object.FALSE
			}
		},
		Name: "<",
	},
	{
		Fn: func(receiver object.Object) object.BuiltinMethodBody {
			return func(args ...object.Object) object.Object {
				err := checkArgumentLen(args, object.StringClass, "==")

				if err != nil {
					return err
				}

				leftValue := receiver.(*object.StringObject).Value
				right, ok := args[0].(*object.StringObject)

				if !ok {
					return wrongTypeError(object.StringClass)
				}

				rightValue := right.Value

				if leftValue == rightValue {
					return object.TRUE
				}

				return object.FALSE
			}
		},
		Name: "==",
	},
	{
		Fn: func(receiver object.Object) object.BuiltinMethodBody {
			return func(args ...object.Object) object.Object {
				err := checkArgumentLen(args, object.StringClass, "!=")

				if err != nil {
					return err
				}

				leftValue := receiver.(*object.StringObject).Value
				right, ok := args[0].(*object.StringObject)

				if !ok {
					return wrongTypeError(object.StringClass)
				}

				rightValue := right.Value

				if leftValue != rightValue {
					return object.TRUE
				}

				return object.FALSE
			}
		},
		Name: "!=",
	},
}

func initializeStringClass() *object.RString {
	methods := object.NewEnvironment()

	for _, m := range builtinStringMethods {
		methods.Set(m.Name, m)
	}

	bc := &object.BaseClass{Name: "String", Methods: methods, Class: ClassClass, SuperClass: ObjectClass}
	sc := &object.RString{BaseClass: bc}
	object.StringClass = sc
	return sc
}
