package object

import (
	"github.com/st0012/rooby/ast"
)

type Class interface {
	LookupClassMethod(string) Object
	LookupInstanceMethod(string) Object
	Object
}

type RClass struct {
	Scope *Scope
	BaseClass
}

type BaseClass struct {
	Name       *ast.Constant
	Methods    *Environment
	SuperClass *RClass
	Class      *RClass
}

func (c *BaseClass) Type() ObjectType {
	return CLASS_OBJ
}

func (c *BaseClass) Inspect() string {
	return "<Class:" + c.Name.Value + ">"
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
