package initializer

import (
	"github.com/st0012/rooby/ast"
	"github.com/st0012/rooby/object"
)

var (
	StringClass *object.StringClass
)

var builtinStringMethods = []*object.BuiltInMethod{
	{
		Fn: func(receiver object.Object) object.BuiltinMethodBody {
			return func(args ...object.Object) object.Object {
				err := checkArgumentLen(args, StringClass, "+")

				if err != nil {
					return err
				}

				leftValue := receiver.(*object.StringObject).Value
				right, ok := args[0].(*object.StringObject)

				if !ok {
					return wrongTypeError(StringClass)
				}

				rightValue := right.Value
				return &object.StringObject{Value: leftValue + rightValue, Class: StringClass}
			}
		},
		Name: "+",
	},
	{
		Fn: func(receiver object.Object) object.BuiltinMethodBody {
			return func(args ...object.Object) object.Object {
				err := checkArgumentLen(args, StringClass, ">")
				if err != nil {
					return err
				}

				leftValue := receiver.(*object.StringObject).Value
				right, ok := args[0].(*object.StringObject)

				if !ok {
					return wrongTypeError(StringClass)
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
				err := checkArgumentLen(args, StringClass, "<")
				if err != nil {
					return err
				}

				leftValue := receiver.(*object.StringObject).Value
				right, ok := args[0].(*object.StringObject)

				if !ok {
					return wrongTypeError(StringClass)
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
				err := checkArgumentLen(args, StringClass, "==")

				if err != nil {
					return err
				}

				leftValue := receiver.(*object.StringObject).Value
				right, ok := args[0].(*object.StringObject)

				if !ok {
					return wrongTypeError(StringClass)
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
				err := checkArgumentLen(args, StringClass, "!=")

				if err != nil {
					return err
				}

				leftValue := receiver.(*object.StringObject).Value
				right, ok := args[0].(*object.StringObject)

				if !ok {
					return wrongTypeError(StringClass)
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

func initializeStringClass() *object.StringClass {
	methods := object.NewEnvironment()

	for _, m := range builtinStringMethods {
		methods.Set(m.Name, m)
	}

	n := &ast.Constant{Value: "String"}
	bc := &object.BaseClass{Name: n, Methods: methods, Class: ClassClass, SuperClass: ObjectClass}
	sc := &object.StringClass{BaseClass: bc}
	StringClass = sc
	return sc
}
