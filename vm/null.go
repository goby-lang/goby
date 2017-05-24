package vm

var (
	nullClass *RNull
	NULL      *NullObject
)

type RNull struct {
	*BaseClass
}

type NullObject struct {
	Class *RNull
}

func (n *NullObject) objectType() objectType {
	return nullObj
}

func (n *NullObject) Inspect() string {
	return "null"
}

func (n *NullObject) returnClass() Class {
	return n.Class
}

func initNull() {
	methods := newEnvironment()

	for _, m := range builtInNullMethods {
		methods.set(m.Name, m)
	}

	baseClass := &BaseClass{Name: "Null", Methods: methods, ClassMethods: newEnvironment(), Class: classClass, pseudoSuperClass: objectClass}
	nc := &RNull{BaseClass: baseClass}
	nullClass = nc
	NULL = &NullObject{Class: nullClass}
}

var builtInNullMethods = []*BuiltInMethod{
	{
		Name: "!",
		Fn: func(receiver Object) builtinMethodBody {
			return func(vm *VM, args []Object, blockFrame *callFrame) Object {

				return TRUE
			}
		},
	},
}
