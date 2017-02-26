package object

import (
	"github.com/st0012/rooby/ast"
)

type Class struct {
	Name       *ast.Constant
	Scope      *Scope
	Methods    *Environment
	SuperClass *Class
	Class      *Class
}

func (c *Class) Type() ObjectType {
	return CLASS_OBJ
}

func (c *Class) Inspect() string {
	return "<Class:" + c.Name.Value + ">"
}

func (c *Class) LookupClassMethod(method_name string) Object {
	return c.Class.LookupInstanceMethod(method_name)
}

func (c *Class) LookupInstanceMethod(method_name string) Object {
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
