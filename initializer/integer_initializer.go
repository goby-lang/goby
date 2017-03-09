package initializer

import (
	"github.com/st0012/Rooby/object"
)

var builtinIntegerMethods = []*object.BuiltInMethod{
	{
		Fn: func(receiver object.Object) object.BuiltinMethodBody {
			return func(args ...object.Object) object.Object {
				err := checkArgumentLen(args, object.IntegerClass, "+")

				if err != nil {
					return err
				}

				leftValue := receiver.(*object.IntegerObject).Value
				right, ok := args[0].(*object.IntegerObject)

				if !ok {
					return wrongTypeError(object.IntegerClass)
				}

				rightValue := right.Value
				return &object.IntegerObject{Value: leftValue + rightValue, Class: object.IntegerClass}
			}
		},
		Name: "+",
	},
	{
		Fn: func(receiver object.Object) object.BuiltinMethodBody {
			return func(args ...object.Object) object.Object {
				err := checkArgumentLen(args, object.IntegerClass, "-")

				if err != nil {
					return err
				}

				leftValue := receiver.(*object.IntegerObject).Value
				right, ok := args[0].(*object.IntegerObject)

				if !ok {
					return wrongTypeError(object.IntegerClass)
				}

				rightValue := right.Value
				return &object.IntegerObject{Value: leftValue - rightValue, Class: object.IntegerClass}
			}
		},
		Name: "-",
	},
	{
		Fn: func(receiver object.Object) object.BuiltinMethodBody {
			return func(args ...object.Object) object.Object {
				err := checkArgumentLen(args, object.IntegerClass, "+")

				if err != nil {
					return err
				}

				leftValue := receiver.(*object.IntegerObject).Value
				right, ok := args[0].(*object.IntegerObject)

				if !ok {
					return wrongTypeError(object.IntegerClass)
				}

				rightValue := right.Value
				return &object.IntegerObject{Value: leftValue * rightValue, Class: object.IntegerClass}
			}
		},
		Name: "*",
	},
	{
		Fn: func(receiver object.Object) object.BuiltinMethodBody {
			return func(args ...object.Object) object.Object {
				err := checkArgumentLen(args, object.IntegerClass, "+")

				if err != nil {
					return err
				}

				leftValue := receiver.(*object.IntegerObject).Value
				right, ok := args[0].(*object.IntegerObject)

				if !ok {
					return wrongTypeError(object.IntegerClass)
				}

				rightValue := right.Value
				return &object.IntegerObject{Value: leftValue / rightValue, Class: object.IntegerClass}
			}
		},
		Name: "/",
	},
	{
		Fn: func(receiver object.Object) object.BuiltinMethodBody {
			return func(args ...object.Object) object.Object {
				err := checkArgumentLen(args, object.IntegerClass, ">")
				if err != nil {
					return err
				}

				leftValue := receiver.(*object.IntegerObject).Value
				right, ok := args[0].(*object.IntegerObject)

				if !ok {
					return wrongTypeError(object.IntegerClass)
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
				err := checkArgumentLen(args, object.IntegerClass, "<")
				if err != nil {
					return err
				}

				leftValue := receiver.(*object.IntegerObject).Value
				right, ok := args[0].(*object.IntegerObject)

				if !ok {
					return wrongTypeError(object.IntegerClass)
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
				err := checkArgumentLen(args, object.IntegerClass, "==")

				if err != nil {
					return err
				}

				leftValue := receiver.(*object.IntegerObject).Value
				right, ok := args[0].(*object.IntegerObject)

				if !ok {
					return wrongTypeError(object.IntegerClass)
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
				err := checkArgumentLen(args, object.IntegerClass, "!=")

				if err != nil {
					return err
				}

				leftValue := receiver.(*object.IntegerObject).Value
				right, ok := args[0].(*object.IntegerObject)

				if !ok {
					return wrongTypeError(object.IntegerClass)
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
	{
		Fn: func(receiver object.Object) object.BuiltinMethodBody {
			return func(args ...object.Object) object.Object {
				if len(args) > 0 {
					return &object.Error{Message: "Too many arguments for Integer#++"}
				}

				int := receiver.(*object.IntegerObject)
				int.Value += 1
				return int
			}
		},
		Name: "++",
	},
	{
		Fn: func(receiver object.Object) object.BuiltinMethodBody {
			return func(args ...object.Object) object.Object {
				if len(args) > 0 {
					return &object.Error{Message: "Too many arguments for Integer#--"}
				}

				int := receiver.(*object.IntegerObject)
				int.Value -= 1
				return int
			}
		},
		Name: "--",
	},
}

func initializeIntegerClass() *object.RInteger {
	methods := object.NewEnvironment()

	for _, m := range builtinIntegerMethods {
		methods.Set(m.Name, m)
	}

	bc := &object.BaseClass{Name: "Integer", Methods: methods, Class: ClassClass, SuperClass: ObjectClass}
	ic := &object.RInteger{BaseClass: bc}
	object.IntegerClass = ic
	return ic
}
