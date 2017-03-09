package evaluator

import (
	"fmt"
)

var (
	IntegerClass *RInteger
)

type RInteger struct {
	*BaseClass
}

type IntegerObject struct {
	Class *RInteger
	Value int
}

func (i *IntegerObject) Type() ObjectType {
	return INTEGER_OBJ
}

func (i *IntegerObject) Inspect() string {
	return fmt.Sprintf("%d", i.Value)
}

func (i *IntegerObject) ReturnClass() Class {
	return i.Class
}

func InitilaizeInteger(value int) *IntegerObject {
	return &IntegerObject{Value: value, Class: IntegerClass}
}

var builtinIntegerMethods = []*BuiltInMethod{
	{
		Fn: func(receiver Object) BuiltinMethodBody {
			return func(args []Object, block *Method) Object {
				err := checkArgumentLen(args, IntegerClass, "+")

				if err != nil {
					return err
				}

				leftValue := receiver.(*IntegerObject).Value
				right, ok := args[0].(*IntegerObject)

				if !ok {
					return wrongTypeError(IntegerClass)
				}

				rightValue := right.Value
				return &IntegerObject{Value: leftValue + rightValue, Class: IntegerClass}
			}
		},
		Name: "+",
	},
	{
		Fn: func(receiver Object) BuiltinMethodBody {
			return func(args []Object, block *Method) Object {
				err := checkArgumentLen(args, IntegerClass, "-")

				if err != nil {
					return err
				}

				leftValue := receiver.(*IntegerObject).Value
				right, ok := args[0].(*IntegerObject)

				if !ok {
					return wrongTypeError(IntegerClass)
				}

				rightValue := right.Value
				return &IntegerObject{Value: leftValue - rightValue, Class: IntegerClass}
			}
		},
		Name: "-",
	},
	{
		Fn: func(receiver Object) BuiltinMethodBody {
			return func(args []Object, block *Method) Object {
				err := checkArgumentLen(args, IntegerClass, "+")

				if err != nil {
					return err
				}

				leftValue := receiver.(*IntegerObject).Value
				right, ok := args[0].(*IntegerObject)

				if !ok {
					return wrongTypeError(IntegerClass)
				}

				rightValue := right.Value
				return &IntegerObject{Value: leftValue * rightValue, Class: IntegerClass}
			}
		},
		Name: "*",
	},
	{
		Fn: func(receiver Object) BuiltinMethodBody {
			return func(args []Object, block *Method) Object {
				err := checkArgumentLen(args, IntegerClass, "+")

				if err != nil {
					return err
				}

				leftValue := receiver.(*IntegerObject).Value
				right, ok := args[0].(*IntegerObject)

				if !ok {
					return wrongTypeError(IntegerClass)
				}

				rightValue := right.Value
				return &IntegerObject{Value: leftValue / rightValue, Class: IntegerClass}
			}
		},
		Name: "/",
	},
	{
		Fn: func(receiver Object) BuiltinMethodBody {
			return func(args []Object, block *Method) Object {
				err := checkArgumentLen(args, IntegerClass, ">")
				if err != nil {
					return err
				}

				leftValue := receiver.(*IntegerObject).Value
				right, ok := args[0].(*IntegerObject)

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
		Fn: func(receiver Object) BuiltinMethodBody {
			return func(args []Object, block *Method) Object {
				err := checkArgumentLen(args, IntegerClass, "<")
				if err != nil {
					return err
				}

				leftValue := receiver.(*IntegerObject).Value
				right, ok := args[0].(*IntegerObject)

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
		Fn: func(receiver Object) BuiltinMethodBody {
			return func(args []Object, block *Method) Object {
				err := checkArgumentLen(args, IntegerClass, "==")

				if err != nil {
					return err
				}

				leftValue := receiver.(*IntegerObject).Value
				right, ok := args[0].(*IntegerObject)

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
		Fn: func(receiver Object) BuiltinMethodBody {
			return func(args []Object, block *Method) Object {
				err := checkArgumentLen(args, IntegerClass, "!=")

				if err != nil {
					return err
				}

				leftValue := receiver.(*IntegerObject).Value
				right, ok := args[0].(*IntegerObject)

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
	{
		Fn: func(receiver Object) BuiltinMethodBody {
			return func(args []Object, block *Method) Object {
				if len(args) > 0 {
					return &Error{Message: "Too many arguments for Integer#++"}
				}

				int := receiver.(*IntegerObject)
				int.Value += 1
				return int
			}
		},
		Name: "++",
	},
	{
		Fn: func(receiver Object) BuiltinMethodBody {
			return func(args []Object, block *Method) Object {
				if len(args) > 0 {
					return &Error{Message: "Too many arguments for Integer#--"}
				}

				int := receiver.(*IntegerObject)
				int.Value -= 1
				return int
			}
		},
		Name: "--",
	},
}

func initInteger() {
	methods := NewEnvironment()

	for _, m := range builtinIntegerMethods {
		methods.Set(m.Name, m)
	}

	bc := &BaseClass{Name: "Integer", Methods: methods, Class: ClassClass, SuperClass: ObjectClass}
	ic := &RInteger{BaseClass: bc}
	IntegerClass = ic
}
