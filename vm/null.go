package vm

var (
	NullClass *RNull
	NULL      *Null
)

type RNull struct {
	*BaseClass
}

type Null struct {
	Class *RNull
}

func (n *Null) Type() ObjectType {
	return NULL_OBJ
}

func (n *Null) Inspect() string {
	return "null"
}

func (n *Null) ReturnClass() Class {
	return n.Class
}

func initNull() {
	methods := NewEnvironment()

	for _, m := range builtInNullMethods {
		methods.Set(m.Name, m)
	}

	baseClass := &BaseClass{Name: "Null", Methods: methods, ClassMethods: NewEnvironment(), Class: ClassClass, SuperClass: ObjectClass}
	nc := &RNull{BaseClass: baseClass}
	NullClass = nc
	NULL = &Null{Class: NullClass}
}

var builtInNullMethods = []*BuiltInMethod{
	{
		Fn: func(receiver Object) BuiltinMethodBody {
			return func(args []Object, block *Method) Object {
				return TRUE
			}
		},
		Name: "!",
	},
}
