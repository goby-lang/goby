package vm

import (
	"fmt"
	"io/ioutil"
	"path"
	"time"
)

var (
	objectClass *RClass
	classClass  *RClass
)

func initTopLevelClasses() {
	classClass = &RClass{
		BaseClass: &BaseClass{
			Name:         "Class",
			Methods:      newEnvironment(),
			ClassMethods: newEnvironment(),
			constants:    make(map[string]*Pointer),
		},
	}

	classClass.setBuiltInMethods(builtinGlobalMethods, false)
	classClass.setBuiltInMethods(builtinGlobalMethods, true)
	classClass.setBuiltInMethods(builtinClassClassMethods, true)

	objectClass = &RClass{
		BaseClass: &BaseClass{
			Name:         "Object",
			Class:        classClass,
			ClassMethods: newEnvironment(),
			Methods:      newEnvironment(),
			constants:    make(map[string]*Pointer),
		},
	}

	objectClass.setBuiltInMethods(builtinGlobalMethods, false)
}

// initializeClass initializes and returns a class instance with given class name
func initializeClass(name string, isModule bool) *RClass {
	class := &RClass{
		BaseClass: &BaseClass{
			Name:             name,
			Methods:          newEnvironment(),
			ClassMethods:     newEnvironment(),
			Class:            classClass,
			pseudoSuperClass: objectClass,
			superClass:       objectClass,
			constants:        make(map[string]*Pointer),
			isModule:         isModule,
		},
	}

	return class
}

// Class is an interface that implements a class's basic functions.
//
// - lookupClassMethod: search for current class's class method with given name.
// - lookupInstanceMethod: search for current class's instance method with given name.
// - ReturnName returns class's name
type Class interface {
	lookupClassMethod(string) Object
	lookupInstanceMethod(string) Object
	lookupConstant(string, bool) *Pointer
	ReturnName() string
	returnSuperClass() Class
	setSingletonMethod(string, *MethodObject)
	Object
}

// RClass represents normal (not built in) class object
type RClass struct {
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
	// pseudoSuperClass points to the class it inherits
	pseudoSuperClass *RClass
	// This is the class where we should looking for a method.
	// It can be normal class, singleton class or a module.
	superClass *RClass
	// Class points to this class's class, which should be ClassClass
	Class *RClass
	// Singleton is a flag marks if this class a singleton class
	Singleton bool
	isModule  bool
	constants map[string]*Pointer
	scope     Class
}

// Inspect returns the basic inspected result (which is class name) of current class
// TODO: Singleton class's inspect() should also mark if it's a singleton class explicitly.
func (c *BaseClass) Inspect() string {
	if c.isModule {
		return "<Module:" + c.Name + ">"
	}
	return "<Class:" + c.Name + ">"
}

func (c *BaseClass) setBuiltInMethods(methodList []*BuiltInMethodObject, classMethods bool) {
	for _, m := range methodList {
		c.Methods.set(m.Name, m)
		m.class = methodClass
	}

	if classMethods {
		for _, m := range methodList {
			c.ClassMethods.set(m.Name, m)
			m.class = methodClass
		}
	}
}

func (c *BaseClass) lookupClassMethod(methodName string) Object {
	method, ok := c.ClassMethods.get(methodName)

	if !ok {
		if c.superClass != nil {
			return c.superClass.lookupClassMethod(methodName)
		}
		if c.Class != nil {
			return c.Class.lookupClassMethod(methodName)
		}
		return nil
	}

	return method
}

func (c *BaseClass) lookupInstanceMethod(methodName string) Object {
	method, ok := c.Methods.get(methodName)

	if !ok {
		if c.superClass != nil {
			return c.superClass.lookupInstanceMethod(methodName)
		}

		if c.Class != nil {
			return c.Class.lookupInstanceMethod(methodName)
		}

		return nil
	}

	return method
}

func (c *BaseClass) lookupConstant(constName string, findInScope bool) *Pointer {
	constant, ok := c.constants[constName]

	if !ok {
		if findInScope && c.scope != nil {
			return c.scope.lookupConstant(constName, true)
		}

		if c.superClass != nil {
			return c.superClass.lookupConstant(constName, false)
		}

		return nil
	}

	return constant
}

// setSingletonMethod will sets method to class's singleton class
// However, if the class doesn't have a singleton class, it will create one for it first.
func (c *BaseClass) setSingletonMethod(name string, method *MethodObject) {
	if c.pseudoSuperClass.Singleton {
		c.pseudoSuperClass.ClassMethods.set(name, method)
	}

	class := initializeClass(c.Name+"singleton", false)
	class.Singleton = true
	class.ClassMethods.set(name, method)
	class.superClass = c.superClass
	class.Class = classClass
	c.superClass = class
}

func (c *BaseClass) returnClass() Class {
	return c.Class
}

// ReturnName returns the name of the class
func (c *BaseClass) ReturnName() string {
	return c.Name
}

func (c *BaseClass) returnSuperClass() Class {
	return c.pseudoSuperClass
}

func (c *RClass) initializeInstance() *RObject {
	instance := &RObject{Class: c, InstanceVariables: newEnvironment()}

	return instance
}

func (c *RClass) setAttrWriter(args interface{}) {

	switch args := args.(type) {
	case []Object:
		for _, attr := range args {
			attrName := attr.(*StringObject).Value
			c.Methods.set(attrName+"=", generateAttrWriteMethod(attrName))
		}
	case []string:
		for _, attrName := range args {
			c.Methods.set(attrName+"=", generateAttrWriteMethod(attrName))
		}
	}

}

func (c *RClass) setAttrReader(args interface{}) {
	switch args := args.(type) {
	case []Object:
		for _, attr := range args {
			attrName := attr.(*StringObject).Value
			c.Methods.set(attrName, generateAttrReadMethod(attrName))
		}
	case []string:
		for _, attrName := range args {
			c.Methods.set(attrName, generateAttrReadMethod(attrName))
		}
	}

}

func generateAttrWriteMethod(attrName string) *BuiltInMethodObject {
	return &BuiltInMethodObject{
		Name: attrName + "=",
		Fn: func(receiver Object) builtinMethodBody {
			return func(t *thread, args []Object, blockFrame *callFrame) Object {
				v := receiver.(*RObject).InstanceVariables.set("@"+attrName, args[0])
				return v
			}
		},
	}
}

func generateAttrReadMethod(attrName string) *BuiltInMethodObject {
	return &BuiltInMethodObject{
		Name: attrName,
		Fn: func(receiver Object) builtinMethodBody {
			return func(t *thread, args []Object, blockFrame *callFrame) Object {
				v, ok := receiver.(*RObject).InstanceVariables.get("@" + attrName)

				if ok {
					return v
				}

				return NULL
			}
		},
	}
}

func (c *RClass) setAttrAccessor(args interface{}) {
	c.setAttrReader(args)
	c.setAttrWriter(args)
}

func createBaseClass(className string) *BaseClass {
	return &BaseClass{
		Name:             className,
		Methods:          newEnvironment(),
		ClassMethods:     newEnvironment(),
		Class:            classClass,
		pseudoSuperClass: objectClass,
		superClass:       objectClass,
	}
}

var builtinGlobalMethods = []*BuiltInMethodObject{
	{
		Name: "require",
		Fn: func(receiver Object) builtinMethodBody {
			return func(t *thread, args []Object, blockFrame *callFrame) Object {
				libName := args[0].(*StringObject).Value
				initFunc, ok := standardLibraries[libName]

				if !ok {
					msg := "Can't require \"" + libName + "\""
					t.returnError(msg)
				}

				initFunc(t.vm)

				return TRUE
			}
		},
	},
	{
		Name: "require_relative",
		Fn: func(receiver Object) builtinMethodBody {
			return func(t *thread, args []Object, blockFrame *callFrame) Object {
				callerDir := path.Dir(t.vm.currentFilePath())
				filepath := args[0].(*StringObject).Value

				filepath = path.Join(callerDir, filepath)

				file, err := ioutil.ReadFile(filepath + ".gb")

				if err != nil {
					t.returnError(err.Error())
				}

				t.vm.execRequiredFile(filepath, file)

				return TRUE
			}
		},
	},
	{
		Name: "puts",
		Fn: func(receiver Object) builtinMethodBody {
			return func(t *thread, args []Object, blockFrame *callFrame) Object {

				for _, arg := range args {
					fmt.Println(arg.Inspect())
				}

				return NULL
			}
		},
	},
	{
		Name: "class",
		Fn: func(receiver Object) builtinMethodBody {
			return func(t *thread, args []Object, blockFrame *callFrame) Object {

				switch r := receiver.(type) {
				case Object:
					return r.returnClass()
				default:
					return &Error{Message: "Can't call class on %T" + string(r.returnClass().ReturnName())}
				}
			}
		},
	},
	{
		Name: "!",
		Fn: func(receiver Object) builtinMethodBody {
			return func(t *thread, args []Object, blockFrame *callFrame) Object {

				return FALSE
			}
		},
	},
	{
		Name: "sleep",
		Fn: func(receiver Object) builtinMethodBody {
			return func(t *thread, args []Object, blockFrame *callFrame) Object {
				int := args[0].(*IntegerObject)
				seconds := int.Value
				time.Sleep(time.Duration(seconds) * time.Second)
				return int
			}
		},
	},
}

// BuiltinClassMethods is a collection of class methods used by Class
var builtinClassClassMethods = []*BuiltInMethodObject{
	{
		Name: "attr_reader",
		Fn: func(receiver Object) builtinMethodBody {
			return func(t *thread, args []Object, blockFrame *callFrame) Object {
				r := receiver.(*RClass)
				r.setAttrReader(args)

				return r
			}
		},
	},
	{
		Name: "attr_writer",
		Fn: func(receiver Object) builtinMethodBody {
			return func(t *thread, args []Object, blockFrame *callFrame) Object {
				r := receiver.(*RClass)
				r.setAttrWriter(args)

				return r
			}
		},
	},
	{
		Name: "attr_accessor",
		Fn: func(receiver Object) builtinMethodBody {
			return func(t *thread, args []Object, blockFrame *callFrame) Object {
				r := receiver.(*RClass)
				r.setAttrAccessor(args)

				return r
			}
		},
	},
	{
		Name: "include",
		Fn: func(receiver Object) builtinMethodBody {
			return func(t *thread, args []Object, blockFrame *callFrame) Object {
				module := args[0].(*RClass)
				class := receiver.(*RClass)
				module.superClass = class.superClass
				class.superClass = module

				return class
			}
		},
	},
	{
		Name: "new",
		Fn: func(receiver Object) builtinMethodBody {
			return func(t *thread, args []Object, blockFrame *callFrame) Object {
				class := receiver.(*RClass)

				if class.pseudoSuperClass.isModule {
					t.returnError("Module inheritance is not supported: " + class.pseudoSuperClass.Name)
				}

				instance := class.initializeInstance()
				initMethod := class.lookupInstanceMethod("initialize")

				if initMethod != nil {
					instance.InitializeMethod = initMethod.(*MethodObject)
				}

				return instance
			}
		},
	},
	{
		Name: "name",
		Fn: func(receiver Object) builtinMethodBody {
			return func(t *thread, args []Object, blockFrame *callFrame) Object {

				name := receiver.(Class).ReturnName()
				nameString := initializeString(name)
				return nameString
			}
		},
	},
	{
		Name: "superclass",
		Fn: func(receiver Object) builtinMethodBody {
			return func(t *thread, args []Object, blockFrame *callFrame) Object {

				return receiver.(Class).returnSuperClass()
			}
		},
	},
}
