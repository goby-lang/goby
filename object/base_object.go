package object

type RObject struct {
	Class             *RClass
	InstanceVariables *Environment
	Scope             *Scope
	InitializeMethod  *Method
}

func (bo *RObject) Type() ObjectType {
	return BASE_OBJECT_OBJ
}

func (bo *RObject) Inspect() string {
	return "<Instance of: " + bo.Class.Name.Value + ">"
}
