package initializer

import (
	"github.com/st0012/rooby/ast"
	"github.com/st0012/rooby/object"
)

var (
	IntegerClass *object.IntegerClass
)

var builtinIntegerMethods = map[string]*object.BuiltInMethod{
	"+": {
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
		Des:  "Addition",
		Name: "+",
	},
	"-": {
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
		Des:  "Subtraction",
		Name: "-",
	},
	"*": {
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
		Des:  "Multiplication",
		Name: "*",
	},
	"/": {
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
		Des:  "Division",
		Name: "/",
	},
	">": {
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
		Des:  "Compare two integers",
		Name: ">",
	},
	"<": {
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
		Des:  "Compare two integers",
		Name: "<",
	},
	"==": {
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
		Des:  "Compare two integers",
		Name: "==",
	},
	"!=": {
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
		Des:  "Compare two integers",
		Name: "!=",
	},
}

func initializeIntegerClass() *object.IntegerClass {
	methods := object.NewEnvironment()

	for name, method := range builtinIntegerMethods {
		methods.Set(name, method)
	}

	n := &ast.Constant{Value: "Integer"}
	bc := &object.BaseClass{Name: n, Methods: methods, Class: ClassClass, SuperClass: ObjectClass}
	ic := &object.IntegerClass{BaseClass: bc}
	IntegerClass = ic
	return ic
}
