package vm

var (
	// NULL represents Goby's null objects.
	NULL *NullObject
)

// NullObject (`nil`) represents the null value in Goby.
// `nil` is convert into `null` when exported to JSON format.
// - `Null.new` is not supported.
type NullObject struct {
	*baseObj
}

func (vm *VM) initNullClass() *RClass {
	nc := vm.initializeClass(nullClass, false)
	nc.setBuiltInMethods(builtInNullInstanceMethods(), false)
	nc.setBuiltInMethods(builtInNullClassMethods(), true)
	NULL = &NullObject{baseObj: &baseObj{class: nc}}
	return nc
}

func builtInNullClassMethods() []*BuiltInMethodObject {
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

func builtInNullInstanceMethods() []*BuiltInMethodObject {
	return []*BuiltInMethodObject{
		{
			// Returns true: the flipped boolean value of nil object.
			//
			// ```ruby
			// a = nil
			// !a
			// # => true
			// ```
			Name: "!",
			Fn: func(receiver Object) builtinMethodBody {
				return func(t *thread, args []Object, blockFrame *callFrame) Object {

					return TRUE
				}
			},
		},
		{
			// Returns true because it is nil. (See the main implementation of is_nil method in vm/class.go)
			//
			// ```ruby
			// a = nil
			// a == nil
			// # => true
			// ```
			Name: "==",
			Fn: func(receiver Object) builtinMethodBody {
				return func(t *thread, args []Object, blockFrame *callFrame) Object {
					if len(args) != 1 {
						return t.vm.initErrorObject(ArgumentError, "Expect 1 argument. got: %d", len(args))
					}

					if _, ok := args[0].(*NullObject); ok {
						return TRUE
					}
					return FALSE
				}
			},
		},
		{
			// Returns true: the flipped boolean value of nil object.
			//
			// ```ruby
			// a = nil
			// a != nil
			// # => false
			// ```
			Name: "!=",
			Fn: func(receiver Object) builtinMethodBody {
				return func(t *thread, args []Object, blockFrame *callFrame) Object {
					if len(args) != 1 {
						return t.vm.initErrorObject(ArgumentError, "Expect 1 argument. got: %d", len(args))
					}

					if _, ok := args[0].(*NullObject); !ok {
						return TRUE
					}
					return FALSE
				}
			},
		},
		{
			// Returns true because it is nil.
			//
			// ```ruby
			// a = nil
			// a.is_nil
			// # => true
			// ```
			Name: "is_nil",
			Fn: func(receiver Object) builtinMethodBody {
				return func(t *thread, args []Object, blockFrame *callFrame) Object {
					if len(args) != 0 {
						return t.vm.initErrorObject(ArgumentError, "Expect 0 argument. got: %d", len(args))
					}
					return TRUE
				}
			},
		},
	}
}

// Polymorphic helper functions -----------------------------------------

// toString returns the name of NullObject
func (n *NullObject) toString() string {
	return "nil"
}

func (n *NullObject) toJSON() string {
	return "null"
}

func (n *NullObject) value() interface{} {
	return nil
}
