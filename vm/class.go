package vm

import (
	"fmt"
)

var (
	objectClass *RClass
	classClass  *RClass
)

func initTopLevelClasses() {
	globalMethods := newEnvironment()
	classMethods := newEnvironment()

	for _, m := range builtinGlobalMethods {
		globalMethods.set(m.Name, m)
	}

	for _, m := range BuiltinClassMethods {
		classMethods.set(m.Name, m)
	}

	classClass = &RClass{BaseClass: &BaseClass{Name: "Class", Methods: globalMethods, ClassMethods: classMethods}}
	objectClass = &RClass{BaseClass: &BaseClass{Name: "Object", Class: classClass, Methods: globalMethods, ClassMethods: newEnvironment()}}
}

// InitializeClass initializes and returns a class instance with given class name
func InitializeClass(name string) *RClass {
	class := &RClass{BaseClass: &BaseClass{Name: name, Methods: newEnvironment(), ClassMethods: newEnvironment(), Class: classClass, SuperClass: objectClass}}
	//classScope := &scope{self: class, Env: closedEnvironment(scope.Env)}
	//class.scope = classScope

	return class
}

// Class is an interface that implements a class's basic functions.
// - LookupClassMethod: search for current class's class method with given name.
// - LookupInstanceMethod: search for current class's instance method with given name.
// - ReturnName returns class's name
type Class interface {
	LookupClassMethod(string) Object
	LookupInstanceMethod(string) Object
	ReturnName() string
	BaseObject
}

// RClass represents normal (not built in) class object
type RClass struct {
	// Scope contains current class's scope information
	Scope *scope
	*BaseClass
}

// BaseClass is a embedded struct that contains all the essential fields for a class
type BaseClass struct {
	// Name is the class's name
	Name string
	// Methods contains its instances' methods
	Methods *environment
	// ClassMethods contains this class's methods
	ClassMethods *environment
	// SuperClass points to the class it inherits
	SuperClass *RClass
	// Class points to this class's class, which should be ClassClass
	Class *RClass
	// Singleton is a flag marks if this class a singleton class
	Singleton bool
}

// objectType returns class object's type
func (c *BaseClass) objectType() objectType {
	return classObj
}

// inspect returns the basic inspected result (which is class name) of current class
// TODO: Singleton class's inspect() should also mark if it's a singleton class explicitly.
func (c *BaseClass) Inspect() string {
	return "<Class:" + c.Name + ">"
}

func (c *BaseClass) LookupClassMethod(method_name string) Object {
	method, ok := c.ClassMethods.get(method_name)

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
	method, ok := c.Methods.get(method_name)

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

// SetSingletonMethod will sets method to class's singleton class
// However, if the class doesn't have a singleton class, it will create one for it first.
func (c *BaseClass) SetSingletonMethod(name string, method *Method) {
	if c.SuperClass.Singleton {
		c.SuperClass.ClassMethods.set(name, method)
	}

	class := InitializeClass(fmt.Sprintf("%s:singleton", c.Name))
	class.Singleton = true
	class.ClassMethods.set(name, method)
	class.SuperClass = c.SuperClass
	class.Class = classClass
	c.SuperClass = class
}

func (c *BaseClass) returnClass() Class {
	return c.Class
}

func (c *BaseClass) ReturnName() string {
	return c.Name
}

func (c *RClass) initializeInstance() *RObject {
	instance := &RObject{Class: c, InstanceVariables: newEnvironment()}

	return instance
}

var builtinGlobalMethods = []*BuiltInMethod{
	{
		Fn: func(receiver Object) builtinMethodBody {
			return func(vm *VM, args []Object, blockFrame *callFrame) Object {

				for _, arg := range args {
					fmt.Println(arg.Inspect())
				}

				return NULL
			}
		},
		Name: "puts",
	},
	{
		Fn: func(receiver Object) builtinMethodBody {
			return func(vm *VM, args []Object, blockFrame *callFrame) Object {

				switch r := receiver.(type) {
				case BaseObject:
					return r.returnClass()
				case Class:
					return r.returnClass()
				default:
					return &Error{Message: fmt.Sprintf("Can't call class on %T", r)}
				}
			}
		},
		Name: "class",
	},
	{
		Fn: func(receiver Object) builtinMethodBody {
			return func(vm *VM, args []Object, blockFrame *callFrame) Object {

				return FALSE
			}
		},
		Name: "!",
	},
}

var BuiltinClassMethods = []*BuiltInMethod{
	{
		Fn: func(receiver Object) builtinMethodBody {
			return func(vm *VM, args []Object, blockFrame *callFrame) Object {

				class := receiver.(*RClass)
				instance := class.initializeInstance()
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
		Fn: func(receiver Object) builtinMethodBody {
			return func(vm *VM, args []Object, blockFrame *callFrame) Object {

				name := receiver.(Class).ReturnName()
				nameString := initializeString(name)
				return nameString
			}
		},
		Name: "name",
	},
}
