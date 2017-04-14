package vm

import (
	"fmt"
)

var (
	ObjectClass *RClass
	ClassClass  *RClass
)

func initTopLevelClasses() {
	globalMethods := NewEnvironment()
	classMethods := NewEnvironment()

	for _, m := range BuiltinGlobalMethods {
		globalMethods.Set(m.Name, m)
	}

	for _, m := range BuiltinClassMethods {
		classMethods.Set(m.Name, m)
	}

	ClassClass = &RClass{BaseClass: &BaseClass{Name: "Class", Methods: globalMethods, ClassMethods: classMethods}}
	ObjectClass = &RClass{BaseClass: &BaseClass{Name: "Object", Class: ClassClass, Methods: globalMethods, ClassMethods: NewEnvironment()}}
}

func InitializeClass(name string) *RClass {
	class := &RClass{BaseClass: &BaseClass{Name: name, Methods: NewEnvironment(), ClassMethods: NewEnvironment(), Class: ClassClass, SuperClass: ObjectClass}}
	//classScope := &Scope{Self: class, Env: NewClosedEnvironment(scope.Env)}
	//class.Scope = classScope

	return class
}

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
	Name         string
	Methods      *Environment
	ClassMethods *Environment
	SuperClass   *RClass
	Class        *RClass
	Singleton    bool
}

func (c *BaseClass) Type() ObjectType {
	return CLASS_OBJ
}

func (c *BaseClass) Inspect() string {
	return "<Class:" + c.Name + ">"
}

func (c *BaseClass) LookupClassMethod(method_name string) Object {
	method, ok := c.ClassMethods.Get(method_name)

	if !ok {
		if c.SuperClass != nil {
			return c.SuperClass.LookupClassMethod(method_name)
		} else {
			if c.Class != nil {
				return c.Class.LookupClassMethod(method_name)
			}
			return nil
		}
	}

	return method
}

func (c *BaseClass) LookupInstanceMethod(method_name string) Object {
	method, ok := c.Methods.Get(method_name)

	if !ok {
		if c.SuperClass != nil {
			return c.SuperClass.LookupInstanceMethod(method_name)
		} else {
			if c.Class != nil {
				return c.Class.LookupInstanceMethod(method_name)
			}
			return nil
		}
	}

	return method
}

func (c *BaseClass) SetSingletonMethod(name string, method *Method) {
	if c.SuperClass.Singleton {
		c.SuperClass.ClassMethods.Set(name, method)
	}

	class := InitializeClass(fmt.Sprintf("%s:singleton", c.Name))
	class.Singleton = true
	class.ClassMethods.Set(name, method)
	class.SuperClass = c.SuperClass
	class.Class = ClassClass
	c.SuperClass = class
}

func (c *BaseClass) ReturnClass() Class {
	return c.Class
}

func (c *BaseClass) ReturnName() string {
	return c.Name
}

var BuiltinGlobalMethods = []*BuiltInMethod{
	{
		Fn: func(receiver Object) BuiltinMethodBody {
			return func(args []Object, block *Method) Object {
				for _, arg := range args {
					fmt.Println(arg.Inspect())
				}

				return NULL
			}
		},
		Name: "puts",
	},
	{
		Fn: func(receiver Object) BuiltinMethodBody {
			return func(args []Object, block *Method) Object {
				switch r := receiver.(type) {
				case BaseObject:
					return r.ReturnClass()
				case Class:
					return r.ReturnClass()
				default:
					return &Error{Message: fmt.Sprint("Can't call class on %T", r)}
				}
			}
		},
		Name: "class",
	},
	{
		Fn: func(receiver Object) BuiltinMethodBody {
			return func(args []Object, block *Method) Object {
				return FALSE
			}
		},
		Name: "!",
	},
}

var BuiltinClassMethods = []*BuiltInMethod{
	{
		Fn: func(receiver Object) BuiltinMethodBody {
			return func(args []Object, block *Method) Object {
				class := receiver.(*RClass)
				instance := InitializeInstance(class)
				initMethod := class.LookupInstanceMethod("initialize")

				if initMethod != nil {
					instance.InitializeMethod = initMethod.(*Method)
				}

				return instance
			}
		},
		Name: "new",
	},
	{
		Fn: func(receiver Object) BuiltinMethodBody {
			return func(args []Object, block *Method) Object {
				name := receiver.(Class).ReturnName()
				nameString := InitializeString(name)
				return nameString
			}
		},
		Name: "name",
	},
}
