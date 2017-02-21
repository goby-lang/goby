package object

type BaseObject struct {
	Class             *Class
	InstanceVariables *Environment
	Scope             *Scope
}

func (bo *BaseObject) Type() ObjectType {
	return BASE_OBJECT_OBJ
}

func (bo *BaseObject) Inspect() string {
	return "<Instance of: " + bo.Class.Name.Value + ">"
}

func (bo *BaseObject) RespondTo(method_name string) bool {
	_, ok := bo.Class.InstanceMethods.Get(method_name)
	return ok
}
