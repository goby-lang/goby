package object

import (
	"fmt"
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

func (c *Class) LookupClassMethod(method_name string, args []Object) Object {
	method, ok := c.ClassMethods.Get(method_name)

	if !ok {
		if c.SuperClass == nil {
			method = c.Class.LookupClassMethod(method_name, args)
		} else {
			method = c.SuperClass.LookupClassMethod(method_name, args)
		}
	}

	if method == nil {
		return &Error{Message: fmt.Sprintf("%s's class method %s is nil.", c.Inspect(), method_name)}
	}

	return method
}

func (c *Class) LookUpInstanceMethod(method_name string, args []Object) Object {
	method, ok := c.InstanceMethods.Get(method_name)

	if !ok {

		for c != nil {
			method, ok = c.InstanceMethods.Get(method_name)

			if !ok {
				// search superclass's superclass
				c = c.SuperClass

				// but if no more superclasses, return an error.
				if c == nil {
					return &Error{Message: fmt.Sprintf("undefined instance method %s for class %s", method_name, c.Inspect())}
				}
			} else {
				// stop looping
				c = nil
			}
		}
	}

	if method == nil {
		return &Error{Message: fmt.Sprintf("%s's instance method %s is nil.", c.Inspect(), method_name)}
	}

	return method
}
