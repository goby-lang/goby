package object

type Class interface {
	LookupClassMethod(string) Object
	LookupInstanceMethod(string) Object
	ReturnClass() Class
	ReturnName() string
	Object
}

type RClass struct {
	Scope *Scope
	*BaseClass
}

type BaseClass struct {
	Name       string
	Methods    *Environment
	SuperClass *RClass
	Class      *RClass
}

func (c *BaseClass) Type() ObjectType {
	return CLASS_OBJ
}

func (c *BaseClass) Inspect() string {
	return "<Class:" + c.Name + ">"
}

func (c *BaseClass) LookupClassMethod(method_name string) Object {
	return c.Class.LookupInstanceMethod(method_name)
}

func (c *BaseClass) LookupInstanceMethod(method_name string) Object {
	method, ok := c.Methods.Get(method_name)

	if !ok {
		if c.SuperClass != nil {
			return c.SuperClass.LookupInstanceMethod(method_name)
		} else {
			return nil
		}
	}

	return method
}

func (c *BaseClass) ReturnClass() Class {
	return c.Class
}

func (c *BaseClass) ReturnName() string {
	return c.Name
}
