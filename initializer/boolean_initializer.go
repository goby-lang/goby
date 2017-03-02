package initializer

import "github.com/st0012/Rooby/object"

var (
	BooleanClass *object.BooleanClass

	TRUE  *object.BooleanObject
	FALSE *object.BooleanObject
	NULL  *object.Null
)

var builtinBooleanMethods = []*object.BuiltInMethod{
	{
		Fn: func(receiver object.Object) object.BuiltinMethodBody {
			return func(args ...object.Object) object.Object {
				err := checkArgumentLen(args, BooleanClass, "==")

				if err != nil {
					return err
				}

				leftValue := receiver.(*object.BooleanObject).Value
				right, ok := args[0].(*object.BooleanObject)

				if !ok {
					return wrongTypeError(BooleanClass)
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
				err := checkArgumentLen(args, BooleanClass, "!=")

				if err != nil {
					return err
				}

				leftValue := receiver.(*object.BooleanObject).Value
				right, ok := args[0].(*object.BooleanObject)

				if !ok {
					return wrongTypeError(BooleanClass)
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

func initializeBooleanClass() *object.BooleanClass {
	methods := object.NewEnvironment()

	for _, m := range builtinBooleanMethods {
		methods.Set(m.Name, m)
	}

	bc := &object.BaseClass{Name: "Boolean", Methods: methods, Class: ClassClass, SuperClass: ObjectClass}
	b := &object.BooleanClass{BaseClass: bc}
	BooleanClass = b

	TRUE = &object.BooleanObject{Value: true, Class: BooleanClass}
	FALSE = &object.BooleanObject{Value: false, Class: BooleanClass}
	return b
}

func initializeNullClass() *object.NullClass {
	baseClass := &object.BaseClass{Name: "Null", Methods: object.NewEnvironment(), Class: ClassClass, SuperClass: ObjectClass}
	nc := &object.NullClass{BaseClass: baseClass}
	NULL = &object.Null{Class: nc}
	return nc
}
