package vm

import (
	"fmt"

	"github.com/goby-lang/goby/vm/classes"
)

// BooleanObject represents boolean object in goby.
// It includes `true` and `FALSE` which represents logically true and false value.
// - `Boolean.new` is not supported.
type BooleanObject struct {
	*baseObj
	value bool
}

var (
	// TRUE is shared boolean object that represents true
	TRUE *BooleanObject
	// FALSE is shared boolean object that represents false
	FALSE *BooleanObject
)

// Class methods --------------------------------------------------------
func builtinBooleanClassMethods() []*BuiltinMethodObject {
	return []*BuiltinMethodObject{
		{
			Name: "new",
			Fn: func(receiver Object) builtinMethodBody {
				return func(t *thread, args []Object, blockFrame *callFrame) Object {
					return t.unsupportedMethodError("#new", receiver)
				}
			},
		},
	}
}

// Instance methods -----------------------------------------------------
func builtinBooleanInstanceMethods() []*BuiltinMethodObject {
	return []*BuiltinMethodObject{
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

					rightValue := receiver.(*BooleanObject).value

					if rightValue {
						return FALSE
					}

					return TRUE
				}
			},
		},
	}
}

// Internal functions ===================================================

// Functions for initialization -----------------------------------------

func (vm *VM) initBoolClass() *RClass {
	b := vm.initializeClass(classes.BooleanClass, false)
	b.setBuiltinMethods(builtinBooleanInstanceMethods(), false)
	b.setBuiltinMethods(builtinBooleanClassMethods(), true)

	TRUE = &BooleanObject{value: true, baseObj: &baseObj{class: b}}
	FALSE = &BooleanObject{value: false, baseObj: &baseObj{class: b}}

	return b
}

// Polymorphic helper functions -----------------------------------------

// Value returns the object
func (b *BooleanObject) Value() interface{} {
	return b.value
}

// toString returns the object's name as the string format
func (b *BooleanObject) toString() string {
	return fmt.Sprintf("%t", b.value)
}

// toJSON just delegates to `toString`
func (b *BooleanObject) toJSON() string {
	return b.toString()
}

// equal returns true if the Boolean values between receiver and parameter are equal
func (b *BooleanObject) equal(e *BooleanObject) bool {
	return b.value == e.value
}
