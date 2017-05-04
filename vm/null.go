package vm

var (
	nullClass *RNull
	NULL      *Null
)

type RNull struct {
	*BaseClass
}

type Null struct {
	Class *RNull
}

func (n *Null) objectType() objectType {
	return nullObj
}

func (n *Null) Inspect() string {
	return "null"
}

func (n *Null) returnClass() Class {
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
	NULL = &Null{Class: nullClass}
}

var builtInNullMethods = []*BuiltInMethod{
	{
		Fn: func(receiver Object) builtinMethodBody {
			return func(vm *VM, args []Object, blockFrame *callFrame) Object {

				return TRUE
			}
		},
		Name: "!",
	},
}
