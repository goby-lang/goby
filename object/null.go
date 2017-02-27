package object

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
