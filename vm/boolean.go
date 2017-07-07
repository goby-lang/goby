package vm

import (
	"fmt"
)

var (
	// TRUE is shared boolean object that represents true
	TRUE *BooleanObject
	// FALSE is shared boolean object that represents false
	FALSE *BooleanObject
)

// BooleanObject represents boolean object in goby.
// It includes `true` and `FALSE` which represents logically true and false value.
// - `Boolean.new` is not supported.
type BooleanObject struct {
	Class *RClass
	Value bool
}

// toString returns boolean object's value, which is either true or false.
func (b *BooleanObject) toString() string {
	return fmt.Sprintf("%t", b.Value)
}

func (b *BooleanObject) toJSON() string {
	return b.toString()
}

// returnClass returns boolean object's class, which is RBool
func (b *BooleanObject) returnClass() *RClass {
	return b.Class
}

func (b *BooleanObject) equal(e *BooleanObject) bool {
	return b.Value == e.Value
}

func builtInBooleanClassMethods() []*BuiltInMethodObject {
	return []*BuiltInMethodObject{
		{
			Name: "new",
			Fn: func(receiver Object) builtinMethodBody {
				return func(t *thread, args []Object, blockFrame *callFrame) Object {
					return t.UnsupportedMethodError("#new", receiver)
				}
			},
		},
	}
}

func builtinBooleanInstanceMethods() []*BuiltInMethodObject {
	return []*BuiltInMethodObject{
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

					err := checkArgumentLen(args, receiver.returnClass(), "==")
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

					err := checkArgumentLen(args, receiver.returnClass(), "!=")
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
						return wrongTypeError(receiver.returnClass())
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
						return wrongTypeError(receiver.returnClass())
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
}

func (vm *VM) initBoolClass() *RClass {
	b := vm.initializeClass(booleanClass, false)
	b.setBuiltInMethods(builtinBooleanInstanceMethods(), false)
	b.setBuiltInMethods(builtInBooleanClassMethods(), true)

	TRUE = &BooleanObject{Value: true, Class: b}
	FALSE = &BooleanObject{Value: false, Class: b}

	return b
}
