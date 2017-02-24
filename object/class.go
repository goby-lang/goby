package object

import (
	"github.com/st0012/rooby/ast"
)

type Class struct {
	Name            *ast.Constant
	Scope           *Scope
	InstanceMethods *Environment
	ClassMethods    *Environment
	SuperClass      *Class
	Class           *Class
}

func (c *Class) Type() ObjectType {
	return CLASS_OBJ
}

func (c *Class) Inspect() string {
	return "<Class:" + c.Name.Value + ">"
}

func (c *Class) LookupClassMethod(method_name string) Object {
	method, ok := c.ClassMethods.Get(method_name)

	if !ok {
		if c.SuperClass == nil {
			method = c.Class.LookupClassMethod(method_name)
		} else {
			method = c.SuperClass.LookupClassMethod(method_name)
		}
	}

	return method
}

func (c *Class) LookupInstanceMethod(method_name string) Object {
	method, ok := c.InstanceMethods.Get(method_name)

	if !ok {

		for c != nil {
			method, ok = c.InstanceMethods.Get(method_name)

			if !ok {
				// search superclass's superclass
				c = c.SuperClass

				if c == nil {
					return nil
				}
			} else {
				// stop looping
				c = nil
			}
		}
	}

	return method
}
