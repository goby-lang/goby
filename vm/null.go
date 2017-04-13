package vm

var (
	NULL *Null
)

type NullClass struct {
	*BaseClass
}

type Null struct {
	Class *NullClass
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
	nc := &NullClass{BaseClass: baseClass}
	NULL = &Null{Class: nc}
}
