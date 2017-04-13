package vm

var (
	NullClass *RNull
	NULL *Null
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
	baseClass := &BaseClass{Name: "Null", Methods: NewEnvironment(), ClassMethods: NewEnvironment(), Class: ClassClass, SuperClass: ObjectClass}
	nc := &RNull{BaseClass: baseClass}
	NullClass = nc
	NULL = &Null{Class: NullClass}

}
