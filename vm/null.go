package vm

var (
	// NULL represents Goby's null objects.
	NULL *NullObject
)

func (vm *VM) initNullClass() *RClass {
	nc := vm.initializeClass(nullClass, false)
	nc.setBuiltInMethods(builtInNullInstanceMethods(), false)
	nc.setBuiltInMethods(builtInNullClassMethods(), true)
	NULL = &NullObject{baseObj: &baseObj{class: nc}}
	return nc
}

// NullObject (`nil`) represents the null value in Goby.
// `nil` is convert into `null` when exported to JSON format.
// - `Null.new` is not supported.
type NullObject struct {
	*baseObj
}

func (n *NullObject) Value() interface{} {
	return nil
}

func (n *NullObject) toString() string {
	return "nil"
}

func (n *NullObject) toJSON() string {
	return "null"
}

func builtInNullClassMethods() []*BuiltInMethodObject {
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
			// Returns the right value because nil as the receiver in the OR operation is considered as false value
			//
			// ```ruby
			// a = nil; a ||= 123;       a; # => 123
			// a = nil; a ||= "string";  a; # => "string"
			// a = nil; a ||= false;     a; # => false
			// a = nil; a ||= (1..4);    a; # => 1..4
			// a = nil; a ||= { b: 1 };  a; # => { b: 1 }
			// a = nil; a ||= Object;    a; # => Object
			// a = nil; a ||= [1, 2, 3]; a; # => [1, 2, 3]
			// a = nil; a ||= nil;       a; # => nil
			// a = nil; a ||= nil || 1;  a; # => 1
			// a = nil; a ||= 1 || nil;  a; # => 1
			// ```
			//
			// @return [Object]
			Name: "||",
			Fn: func(receiver Object) builtinMethodBody {
				return func(t *thread, args []Object, blockFrame *callFrame) Object {
					if len(args) != 1 {
						return t.vm.initErrorObject(ArgumentError, "Expect 1 argument. got: %d", len(args))
					}
					return args[0]
				}
			},
		},
		{
			Name: "&&",
			Fn: func(receiver Object) builtinMethodBody {
				return func(t *thread, args []Object, blockFrame *callFrame) Object {
					if len(args) != 1 {
						return t.vm.initErrorObject(ArgumentError, "Expect 1 argument. got: %d", len(args))
					}
					return FALSE
				}
			},
		},
		{
			// Returns true because it is nil. (See the main implementation of nil? method in vm/class.go)
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
			// a.nil?
			// # => true
			// ```
			Name: "nil?",
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
