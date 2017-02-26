package object

import (
	"github.com/st0012/rooby/ast"
)

type RClass struct {
	Name       *ast.Constant
	Scope      *Scope
	Methods    *Environment
	SuperClass *RClass
	Class      *RClass
}

func (c *RClass) Type() ObjectType {
	return CLASS_OBJ
}

func (c *RClass) Inspect() string {
	return "<Class:" + c.Name.Value + ">"
}

func (c *RClass) LookupClassMethod(method_name string) Object {
	return c.Class.LookupInstanceMethod(method_name)
}

func (c *RClass) LookupInstanceMethod(method_name string) Object {
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
