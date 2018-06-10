package vm

import (
	"fmt"

	"github.com/goby-lang/goby/vm/classes"
)

// BooleanObject represents boolean object in goby and no instance methods are contained within it.
// `Boolean` class is just a dummy to hold logical `true` and `false` representation and no other active usage.
// `Boolean.new` is not supported.
//
// Please note that class checking such as `#is_a?(Boolean)` **should be avoided in principle**.
// `#is_a?` often leads to redundant code. Consider using `respond_to?` first, but actually it is unnecessary
// in almost all cases.
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
			Fn: func(receiver Object, sourceLine int) builtinMethodBody {
				return func(t *Thread, args []Object, blockFrame *normalCallFrame) Object {
					return t.vm.initUnsupportedMethodError(sourceLine, "#new", receiver)
				}
			},
		},
	}
}

// Internal functions ===================================================

// Functions for initialization -----------------------------------------

func (vm *VM) initBoolClass() *RClass {
	b := vm.initializeClass(classes.BooleanClass)
	b.setBuiltinMethods(builtinBooleanClassMethods(), true)

	TRUE = &BooleanObject{value: true, baseObj: &baseObj{class: b}}
	FALSE = &BooleanObject{value: false, baseObj: &baseObj{class: b}}

	return b
}

// Polymorphic helper functions -----------------------------------------

// Value returns the object
func toBooleanObject(value bool) *BooleanObject {
	if value {
		return TRUE
	}

	return FALSE
}

// Value returns the object
func (b *BooleanObject) Value() interface{} {
	return b.value
}

// toString returns the object's name as the string format
func (b *BooleanObject) toString() string {
	return fmt.Sprintf("%t", b.value)
}

// toJSON just delegates to `toString`
func (b *BooleanObject) toJSON(t *Thread) string {
	return b.toString()
}

func (b *BooleanObject) isTruthy() bool {
	return b.value
}

// equal returns true if the Boolean values between receiver and parameter are equal
func (b *BooleanObject) equal(e *BooleanObject) bool {
	return b.value == e.value
}
