package object

type BaseObject interface {
	ReturnClass() Class
	Object
}

type RObject struct {
	Class             *RClass
	InstanceVariables *Environment
	Scope             *Scope
	InitializeMethod  *Method
}

func (ro *RObject) Type() ObjectType {
	return BASE_OBJECT_OBJ
}

func (ro *RObject) Inspect() string {
	return "<Instance of: " + ro.Class.Name + ">"
}

func (ro *RObject) ReturnClass() Class {
	return ro.Class
}