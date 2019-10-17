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
	*BaseObj
	value bool
}

var (
	// TRUE is shared boolean object that represents true
	TRUE *BooleanObject
	// FALSE is shared boolean object that represents false
	FALSE *BooleanObject
)

// Class methods --------------------------------------------------------
var builtinBooleanClassMethods = []*BuiltinMethodObject{
	{
		Name: "new",
		Fn: func(receiver Object, sourceLine int, t *Thread, args []Object, blockFrame *normalCallFrame) Object {
			return t.vm.InitNoMethodError(sourceLine, "new", receiver)
		},
	},
}

// Internal functions ===================================================

// Functions for initialization -----------------------------------------

func (vm *VM) initBoolClass() *RClass {
	b := vm.initializeClass(classes.BooleanClass)
	b.setBuiltinMethods(builtinBooleanClassMethods, true)

	TRUE = &BooleanObject{value: true, BaseObj: &BaseObj{class: b}}
	FALSE = &BooleanObject{value: false, BaseObj: &BaseObj{class: b}}

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

// ToString returns the object's name as the string format
func (b *BooleanObject) ToString() string {
	return fmt.Sprintf("%t", b.value)
}

// Inspect delegates to ToString
func (b *BooleanObject) Inspect() string {
	return b.ToString()
}

// ToJSON just delegates to `ToString`
func (b *BooleanObject) ToJSON(t *Thread) string {
	return b.ToString()
}

func (b *BooleanObject) isTruthy() bool {
	return b.value
}

func (b *BooleanObject) equalTo(with Object) bool {
	bool, ok := with.(*BooleanObject)

	if !ok {
		return false
	}

	return b.value == bool.value
}

// equal returns true if the Boolean values between receiver and parameter are equal
func (b *BooleanObject) equal(e *BooleanObject) bool {
	return b.value == e.value
}
