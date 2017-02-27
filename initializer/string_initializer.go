package initializer

import (
	"github.com/st0012/rooby/ast"
	"github.com/st0012/rooby/object"
)

var (
	StringClass *object.StringClass
)

var builtinStringMethods = map[string]*object.BuiltInMethod{
	"+": {
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
		Des:  "Add two strings",
		Name: "+",
	},
	">": {
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
		Des:  "Compare two strings",
		Name: ">",
	},
	"<": {
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
		Des:  "Compare two strings",
		Name: "<",
	},
	"==": {
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
		Des:  "Compare two strings",
		Name: "==",
	},
	"!=": {
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
		Des:  "Compare two strings",
		Name: "!=",
	},
}

func initializeStringClass() *object.StringClass {
	methods := object.NewEnvironment()

	for name, method := range builtinStringMethods {
		methods.Set(name, method)
	}

	n := &ast.Constant{Value: "String"}
	bc := &object.BaseClass{Name: n, Methods: methods, Class: ClassClass, SuperClass: ObjectClass}
	sc := &object.StringClass{BaseClass: bc}
	StringClass = sc
	return sc
}
