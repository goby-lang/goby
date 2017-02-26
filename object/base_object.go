package object

type BaseObject struct {
	Class             *Class
	InstanceVariables *Environment
	Scope             *Scope
	InitializeMethod  *Method
}

func (bo *BaseObject) Type() ObjectType {
	return BASE_OBJECT_OBJ
}

func (bo *BaseObject) Inspect() string {
	return "<Instance of: " + bo.Class.Name.Value + ">"
}
