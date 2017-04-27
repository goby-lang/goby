package vm

import (
	"fmt"
)

var (
	booleanClass *RBool

	// TRUE is shared boolean object that represents true
	TRUE *BooleanObject
	// FALSE is shared boolean object that represents false
	FALSE *BooleanObject
)

// RBool is the built in class of rooby's boolean objects.
type RBool struct {
	*BaseClass
}

// BooleanObject represents boolean object in rooby
type BooleanObject struct {
	Class *RBool
	Value bool
}

// Type returns boolean object's type
func (b *BooleanObject) Type() objectType {
	return booleanObj
}

// Inspect returns boolean object's value, which is either true or false.
func (b *BooleanObject) Inspect() string {
	return fmt.Sprintf("%t", b.Value)
}

// returnClass returns boolean object's class, which is RBool
func (b *BooleanObject) returnClass() Class {
	return b.Class
}

var builtinBooleanMethods = []*BuiltInMethod{
	{
		Fn: func(receiver Object) builtinMethodBody {
			return func(vm *VM, args []Object, blockFrame *callFrame) Object {

				err := checkArgumentLen(args, booleanClass, "==")

				if err != nil {
					return err
				}

				leftValue := receiver.(*BooleanObject).Value
				right, ok := args[0].(*BooleanObject)

				if !ok {
					return wrongTypeError(booleanClass)
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
		Fn: func(receiver Object) builtinMethodBody {
			return func(vm *VM, args []Object, blockFrame *callFrame) Object {

				err := checkArgumentLen(args, booleanClass, "!=")

				if err != nil {
					return err
				}

				leftValue := receiver.(*BooleanObject).Value
				right, ok := args[0].(*BooleanObject)

				if !ok {
					return wrongTypeError(booleanClass)
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
		Fn: func(receiver Object) builtinMethodBody {
			return func(vm *VM, args []Object, blockFrame *callFrame) Object {

				rightValue := receiver.(*BooleanObject).Value

				if rightValue {
					return FALSE
				}

				return TRUE
			}
		},
		Name: "!",
	},
}

func initBool() {
	methods := NewEnvironment()

	for _, m := range builtinBooleanMethods {
		methods.Set(m.Name, m)
	}

	bc := &BaseClass{Name: "Boolean", Methods: methods, ClassMethods: NewEnvironment(), Class: classClass, SuperClass: objectClass}
	b := &RBool{BaseClass: bc}
	booleanClass = b

	TRUE = &BooleanObject{Value: true, Class: booleanClass}
	FALSE = &BooleanObject{Value: false, Class: booleanClass}
}
