package vm

var (
	nullClass *RNull
	// NULL represents Goby's null objects.
	NULL *NullObject
)

// RNull is the built in class of Goby's null objects.
type RNull struct {
	*BaseClass
}

// NullObject (`nil`) represents the null value in Goby.
// `nil` is convert into `null` when exported to JSON format.
// Cannot perform `Null.new`.
type NullObject struct {
	Class *RNull
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

func initNullClass() {
	baseClass := &BaseClass{Name: "Null", Methods: newEnvironment(), ClassMethods: newEnvironment(), Class: classClass, pseudoSuperClass: objectClass}
	nc := &RNull{BaseClass: baseClass}
	nc.setBuiltInMethods(builtInNullInstanceMethods, false)
	nullClass = nc
	NULL = &NullObject{Class: nullClass}
}

var builtInNullInstanceMethods = []*BuiltInMethodObject{
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
