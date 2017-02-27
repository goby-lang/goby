package initializer

import (
	"github.com/st0012/rooby/object"
)

var (
	IntegerClass *object.IntegerClass
)

var builtinIntegerMethods = []*object.BuiltInMethod{
	{
		Fn: func(receiver object.Object) object.BuiltinMethodBody {
			return func(args ...object.Object) object.Object {
				err := checkArgumentLen(args, IntegerClass, "+")

				if err != nil {
					return err
				}

				leftValue := receiver.(*object.IntegerObject).Value
				right, ok := args[0].(*object.IntegerObject)

				if !ok {
					return wrongTypeError(IntegerClass)
				}

				rightValue := right.Value
				return &object.IntegerObject{Value: leftValue + rightValue, Class: IntegerClass}
			}
		},
		Name: "+",
	},
	{
		Fn: func(receiver object.Object) object.BuiltinMethodBody {
			return func(args ...object.Object) object.Object {
				err := checkArgumentLen(args, IntegerClass, "-")

				if err != nil {
					return err
				}

				leftValue := receiver.(*object.IntegerObject).Value
				right, ok := args[0].(*object.IntegerObject)

				if !ok {
					return wrongTypeError(IntegerClass)
				}

				rightValue := right.Value
				return &object.IntegerObject{Value: leftValue - rightValue, Class: IntegerClass}
			}
		},
		Name: "-",
	},
	{
		Fn: func(receiver object.Object) object.BuiltinMethodBody {
			return func(args ...object.Object) object.Object {
				err := checkArgumentLen(args, IntegerClass, "+")

				if err != nil {
					return err
				}

				leftValue := receiver.(*object.IntegerObject).Value
				right, ok := args[0].(*object.IntegerObject)

				if !ok {
					return wrongTypeError(IntegerClass)
				}

				rightValue := right.Value
				return &object.IntegerObject{Value: leftValue * rightValue, Class: IntegerClass}
			}
		},
		Name: "*",
	},
	{
		Fn: func(receiver object.Object) object.BuiltinMethodBody {
			return func(args ...object.Object) object.Object {
				err := checkArgumentLen(args, IntegerClass, "+")

				if err != nil {
					return err
				}

				leftValue := receiver.(*object.IntegerObject).Value
				right, ok := args[0].(*object.IntegerObject)

				if !ok {
					return wrongTypeError(IntegerClass)
				}

				rightValue := right.Value
				return &object.IntegerObject{Value: leftValue / rightValue, Class: IntegerClass}
			}
		},
		Name: "/",
	},
	{
		Fn: func(receiver object.Object) object.BuiltinMethodBody {
			return func(args ...object.Object) object.Object {
				err := checkArgumentLen(args, IntegerClass, ">")
				if err != nil {
					return err
				}

				leftValue := receiver.(*object.IntegerObject).Value
				right, ok := args[0].(*object.IntegerObject)

				if !ok {
					return wrongTypeError(IntegerClass)
				}

				rightValue := right.Value

				if leftValue > rightValue {
					return TRUE
				}

				return FALSE
			}
		},
		Name: ">",
	},
	{
		Fn: func(receiver object.Object) object.BuiltinMethodBody {
			return func(args ...object.Object) object.Object {
				err := checkArgumentLen(args, IntegerClass, "<")
				if err != nil {
					return err
				}

				leftValue := receiver.(*object.IntegerObject).Value
				right, ok := args[0].(*object.IntegerObject)

				if !ok {
					return wrongTypeError(IntegerClass)
				}

				rightValue := right.Value

				if leftValue < rightValue {
					return TRUE
				}

				return FALSE
			}
		},
		Name: "<",
	},
	{
		Fn: func(receiver object.Object) object.BuiltinMethodBody {
			return func(args ...object.Object) object.Object {
				err := checkArgumentLen(args, IntegerClass, "==")

				if err != nil {
					return err
				}

				leftValue := receiver.(*object.IntegerObject).Value
				right, ok := args[0].(*object.IntegerObject)

				if !ok {
					return wrongTypeError(IntegerClass)
				}

				rightValue := right.Value

				if leftValue == rightValue {
					return TRUE
				}

				return FALSE
			}
		},
		Name: "==",
	},
	{
		Fn: func(receiver object.Object) object.BuiltinMethodBody {
			return func(args ...object.Object) object.Object {
				err := checkArgumentLen(args, IntegerClass, "!=")

				if err != nil {
					return err
				}

				leftValue := receiver.(*object.IntegerObject).Value
				right, ok := args[0].(*object.IntegerObject)

				if !ok {
					return wrongTypeError(IntegerClass)
				}

				rightValue := right.Value

				if leftValue != rightValue {
					return TRUE
				}

				return FALSE
			}
		},
		Name: "!=",
	},
}

func initializeIntegerClass() *object.IntegerClass {
	methods := object.NewEnvironment()

	for _, m := range builtinIntegerMethods {
		methods.Set(m.Name, m)
	}

	bc := &object.BaseClass{Name: "Integer", Methods: methods, Class: ClassClass, SuperClass: ObjectClass}
	ic := &object.IntegerClass{BaseClass: bc}
	IntegerClass = ic
	return ic
}
