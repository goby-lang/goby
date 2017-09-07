package vm

import (
	"fmt"
	"github.com/goby-lang/goby/vm/classes"
)

var (
	// TRUE is shared boolean object that represents true
	TRUE *BooleanObject
	// FALSE is shared boolean object that represents false
	FALSE *BooleanObject
)

func (vm *VM) initBoolClass() *RClass {
	b := vm.initializeClass(classes.BooleanClass, false)
	b.setBuiltInMethods(builtinBooleanInstanceMethods(), false)
	b.setBuiltInMethods(builtInBooleanClassMethods(), true)

	TRUE = &BooleanObject{value: true, baseObj: &baseObj{class: b}}
	FALSE = &BooleanObject{value: false, baseObj: &baseObj{class: b}}

	return b
}

// BooleanObject represents boolean object in goby.
// It includes `true` and `FALSE` which represents logically true and false value.
// - `Boolean.new` is not supported.
type BooleanObject struct {
	*baseObj
	value bool
}

func (b *BooleanObject) Value() interface{} {
	return b.value
}

// Polymorphic helper functions -----------------------------------------
func (b *BooleanObject) toString() string {
	return fmt.Sprintf("%t", b.value)
}

func (b *BooleanObject) toJSON() string {
	return b.toString()
}

func (b *BooleanObject) equal(e *BooleanObject) bool {
	return b.value == e.value
}

func builtInBooleanClassMethods() []*BuiltInMethodObject {
	return []*BuiltInMethodObject{
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
