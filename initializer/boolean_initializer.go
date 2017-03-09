package initializer

import "github.com/st0012/Rooby/object"

var builtinBooleanMethods = []*object.BuiltInMethod{
	{
		Fn: func(receiver object.Object) object.BuiltinMethodBody {
			return func(args ...object.Object) object.Object {
				err := checkArgumentLen(args, object.BooleanClass, "==")

				if err != nil {
					return err
				}

				leftValue := receiver.(*object.BooleanObject).Value
				right, ok := args[0].(*object.BooleanObject)

				if !ok {
					return wrongTypeError(object.BooleanClass)
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
				err := checkArgumentLen(args, object.BooleanClass, "!=")

				if err != nil {
					return err
				}

				leftValue := receiver.(*object.BooleanObject).Value
				right, ok := args[0].(*object.BooleanObject)

				if !ok {
					return wrongTypeError(object.BooleanClass)
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

func initializeBooleanClass() *object.RBool {
	methods := object.NewEnvironment()

	for _, m := range builtinBooleanMethods {
		methods.Set(m.Name, m)
	}

	bc := &object.BaseClass{Name: "Boolean", Methods: methods, Class: ClassClass, SuperClass: ObjectClass}
	b := &object.RBool{BaseClass: bc}
	object.BooleanClass = b

	object.TRUE = &object.BooleanObject{Value: true, Class: object.BooleanClass}
	object.FALSE = &object.BooleanObject{Value: false, Class: object.BooleanClass}
	return b
}

func initializeNullClass() *object.NullClass {
	baseClass := &object.BaseClass{Name: "Null", Methods: object.NewEnvironment(), Class: ClassClass, SuperClass: ObjectClass}
	nc := &object.NullClass{BaseClass: baseClass}
	object.NULL = &object.Null{Class: nc}
	return nc
}
