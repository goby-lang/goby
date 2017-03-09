package evaluator

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

	ClassClass = &RClass{BaseClass: &BaseClass{Name: "Class", Methods: classMethods}}
	ObjectClass = &RClass{BaseClass: &BaseClass{Name: "Object", Class: ClassClass, Methods: globalMethods}}
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

func InitializeClass(name string, scope *Scope) *RClass {
	class := &RClass{BaseClass: &BaseClass{Name: name, Methods: NewEnvironment(), Class: ClassClass, SuperClass: ObjectClass}}
	classScope := &Scope{Self: class, Env: NewClosedEnvironment(scope.Env)}
	class.Scope = classScope

	return class
}

var BuiltinGlobalMethods = []*BuiltInMethod{
	{
		Fn: func(receiver Object) BuiltinMethodBody {
			return func(args ...Object) Object {
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
			return func(args ...Object) Object {
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
}

var BuiltinClassMethods = []*BuiltInMethod{
	{
		Fn: func(receiver Object) BuiltinMethodBody {
			return func(args ...Object) Object {
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
			return func(args ...Object) Object {
				name := receiver.(Class).ReturnName()
				nameString := &StringObject{Value: name}
				return nameString
			}
		},
		Name: "name",
	},
}
