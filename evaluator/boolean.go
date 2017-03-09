package evaluator

import (
	"fmt"
)

var (
	BooleanClass *RBool
	TRUE         *BooleanObject
	FALSE        *BooleanObject
)

type RBool struct {
	*BaseClass
}

type BooleanObject struct {
	Class *RBool
	Value bool
}

func (b *BooleanObject) Type() ObjectType {
	return BOOLEAN_OBJ
}

func (b *BooleanObject) Inspect() string {
	return fmt.Sprintf("%t", b.Value)
}

func (b *BooleanObject) ReturnClass() Class {
	return b.Class
}

var builtinBooleanMethods = []*BuiltInMethod{
	{
		Fn: func(receiver Object) BuiltinMethodBody {
			return func(args ...Object) Object {
				err := checkArgumentLen(args, BooleanClass, "==")

				if err != nil {
					return err
				}

				leftValue := receiver.(*BooleanObject).Value
				right, ok := args[0].(*BooleanObject)

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
		Fn: func(receiver Object) BuiltinMethodBody {
			return func(args ...Object) Object {
				err := checkArgumentLen(args, BooleanClass, "!=")

				if err != nil {
					return err
				}

				leftValue := receiver.(*BooleanObject).Value
				right, ok := args[0].(*BooleanObject)

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

func initBool() {
	methods := NewEnvironment()

	for _, m := range builtinBooleanMethods {
		methods.Set(m.Name, m)
	}

	bc := &BaseClass{Name: "Boolean", Methods: methods, Class: ClassClass, SuperClass: ObjectClass}
	b := &RBool{BaseClass: bc}
	BooleanClass = b

	TRUE = &BooleanObject{Value: true, Class: BooleanClass}
	FALSE = &BooleanObject{Value: false, Class: BooleanClass}
}
