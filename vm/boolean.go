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

func (vm *VM) initBoolClass() *RClass {
	b := vm.initializeClass(booleanClass, false)
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

					leftValue := receiver.(*BooleanObject).value
					right, ok := args[0].(*BooleanObject)

					if !ok {
						err := t.vm.initErrorObject(TypeError, WrongArgumentTypeFormat, booleanClass, right.Class().Name)
						return err
					}

					rightValue := right.value

					if leftValue && rightValue {
						return TRUE
					}

					return FALSE
				}
			},
		},
		{
			// If the receiver is truthy value, it returns true. In contrary, if the receiver is the falsey value,
			// it returns the right value
			//
			// ```ruby
			// a = true;  a ||= 123;       a; # => true
			// a = true;  a ||= "string";  a; # => true
			// a = true;  a ||= false;     a; # => true
			// a = true;  a ||= (1..4);    a; # => true
			// a = true;  a ||= { b: 1 };  a; # => true
			// a = true;  a ||= Object;    a; # => true
			// a = true;  a ||= [1, 2, 3]; a; # => true
			// a = true;  a ||= nil;       a; # => true
			// a = true;  a ||= nil || 1;  a; # => true
			// a = true;  a ||= 1 || nil;  a; # => true
			// a = false; a ||= 123;       a; # => 123
			// a = false; a ||= "string";  a; # => "string"
			// a = false; a ||= false;     a; # => false
			// a = false; a ||= (1..4);    a; # => 1..4
			// a = false; a ||= { b: 1 };  a; # => { b: 1 }
			// a = false; a ||= Object;    a; # => Object
			// a = false; a ||= [1, 2, 3]; a; # => [1, 2, 3]
			// a = false; a ||= nil;       a; # => nil
			// a = false; a ||= nil || 1;  a; # => 1
			// a = false; a ||= 1 || nil;  a; # => 1
			// ```
			//
			// @return [Object]
			Name: "||",
			Fn: func(receiver Object) builtinMethodBody {
				return func(t *thread, args []Object, blockFrame *callFrame) Object {
					if len(args) != 1 {
						return t.vm.initErrorObject(ArgumentError, "Expect 1 argument. got: %d", len(args))
					}

					leftValue := receiver.(*BooleanObject).value

					if leftValue {
						return receiver
					}
					return args[0]
				}
			},
		},
	}
}
