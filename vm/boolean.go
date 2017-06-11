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

// RBool is the built in class of goby's boolean objects.
type RBool struct {
	*BaseClass
}

// BooleanObject represents boolean object in goby.
// It includes `true` and `FALSE` which represents logically true and false value.
type BooleanObject struct {
	Class *RBool
	Value bool
}

// Inspect returns boolean object's value, which is either true or false.
func (b *BooleanObject) Inspect() string {
	return fmt.Sprintf("%t", b.Value)
}

func (b *BooleanObject) toJSON() string {
	return b.Inspect()
}

// returnClass returns boolean object's class, which is RBool
func (b *BooleanObject) returnClass() Class {
	return b.Class
}

func (b *BooleanObject) equal(e *BooleanObject) bool {
	return b.Value == e.Value
}

var builtinBooleanInstanceMethods = []*BuiltInMethodObject{
	{
		// Returns true if the receiver equals to the argument.
		//
		// ```Ruby
		// 1 == 1 # => true
		// 100 == 33 # => false
		// ```
		// @return [Boolean]
		Name: "==",
		Fn: func(receiver Object) builtinMethodBody {
			return func(t *thread, args []Object, blockFrame *callFrame) Object {

				err := checkArgumentLen(args, booleanClass, "==")
				if err != nil {
					return err
				}

				if receiver == args[0] {
					return TRUE
				}
				return FALSE
			}
		},
	},
	{
		// Returns true if the receiver is not equals to the argument.
		//
		// ```Ruby
		// 4 != 2 # => true
		// 45 != 45 # => false
		// ```
		// @return [Boolean]
		Name: "!=",
		Fn: func(receiver Object) builtinMethodBody {
			return func(t *thread, args []Object, blockFrame *callFrame) Object {

				err := checkArgumentLen(args, booleanClass, "!=")
				if err != nil {
					return err
				}

				if receiver != args[0] {
					return TRUE
				}
				return FALSE
			}
		},
	},
	{
		// Reverse the receiver.
		//
		// ```ruby
		// !true  # => false
		// !false # => true
		// ```
		// @return [Boolean]
		Name: "!",
		Fn: func(receiver Object) builtinMethodBody {
			return func(t *thread, args []Object, blockFrame *callFrame) Object {

				rightValue := receiver.(*BooleanObject).Value

				if rightValue {
					return FALSE
				}

				return TRUE
			}
		},
	},
	{
		// Returns true if both the receiver and the argument are true.
		//
		// ```ruby
		// 3 > 2 && 5 > 3  # => true
		// 3 > 2 && 5 > 10 # => false
		// ```
		// @return [Boolean]
		Name: "&&",
		Fn: func(receiver Object) builtinMethodBody {
			return func(t *thread, args []Object, blockFrame *callFrame) Object {

				leftValue := receiver.(*BooleanObject).Value
				right, ok := args[0].(*BooleanObject)

				if !ok {
					return wrongTypeError(booleanClass)
				}

				rightValue := right.Value

				if leftValue && rightValue {
					return TRUE
				}

				return FALSE
			}
		},
	},
	{
		// Returns true either if the receiver or argument is true.
		//
		// ```ruby
		// 3 > 2 || 5 > 3  # => true
		// 3 > 2 || 5 > 10 # => true
		// 2 > 3 || 5 > 10 # => false
		// ```
		// @return [Boolean]
		Name: "||",
		Fn: func(receiver Object) builtinMethodBody {
			return func(t *thread, args []Object, blockFrame *callFrame) Object {

				leftValue := receiver.(*BooleanObject).Value
				right, ok := args[0].(*BooleanObject)

				if !ok {
					return wrongTypeError(booleanClass)
				}

				rightValue := right.Value

				if leftValue || rightValue {
					return TRUE
				}

				return FALSE
			}
		},
	},
}

func initBool() {
	bc := &BaseClass{Name: "Boolean", Methods: newEnvironment(), ClassMethods: newEnvironment(), Class: classClass, pseudoSuperClass: objectClass, superClass: objectClass}
	b := &RBool{BaseClass: bc}
	b.setBuiltInMethods(builtinBooleanInstanceMethods, false)
	booleanClass = b

	TRUE = &BooleanObject{Value: true, Class: booleanClass}
	FALSE = &BooleanObject{Value: false, Class: booleanClass}
}
