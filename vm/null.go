package vm

var (
	// NULL represents Goby's null objects.
	NULL *NullObject
)

// NullObject (`nil`) represents the null value in Goby.
// `nil` is convert into `null` when exported to JSON format.
// - `Null.new` is not supported.
type NullObject struct {
	Class *RClass
}

// toString returns the name of NullObject
func (n *NullObject) toString() string {
	return "nil"
}

func (n *NullObject) toJSON() string {
	return "null"
}

func (n *NullObject) returnClass() Class {
	return n.Class
}

func (vm *VM) initNullClass() *RClass {
	nc := vm.initializeClass(nullClass, false)
	nc.setBuiltInMethods(builtInNullInstanceMethods(), false)
	nc.setBuiltInMethods(builtInNullClassMethods(), true)
	NULL = &NullObject{Class: nc}
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
	}
}
