package vm

import (
	"fmt"
	"io/ioutil"
	"path"

	"github.com/goby-lang/goby/bytecode"
	"github.com/goby-lang/goby/parser"
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

// initializeClass initializes and returns a class instance with given class name
func initializeClass(name string) *RClass {
	class := &RClass{BaseClass: &BaseClass{Name: name, Methods: newEnvironment(), ClassMethods: newEnvironment(), Class: classClass, pseudoSuperClass: objectClass, superClass: objectClass}}
	//classScope := &scope{self: class, Env: closedEnvironment(scope.Env)}
	//class.scope = classScope

	return class
}

// Class is an interface that implements a class's basic functions.
// - lookupClassMethod: search for current class's class method with given name.
// - lookupInstanceMethod: search for current class's instance method with given name.
// - ReturnName returns class's name
type Class interface {
	lookupClassMethod(string) Object
	lookupInstanceMethod(string) Object
	ReturnName() string
	returnSuperClass() Class
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

// setSingletonMethod will sets method to class's singleton class
// However, if the class doesn't have a singleton class, it will create one for it first.
func (c *BaseClass) setSingletonMethod(name string, method *Method) {
	if c.pseudoSuperClass.Singleton {
		c.pseudoSuperClass.ClassMethods.set(name, method)
	}

	class := initializeClass(c.Name + "singleton")
	class.Singleton = true
	class.ClassMethods.set(name, method)
	class.superClass = c.superClass
	class.Class = classClass
	c.superClass = class
}

func (c *BaseClass) returnClass() Class {
	return c.Class
}

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

var builtinGlobalMethods = []*BuiltInMethod{
	{
		Name: "require",
		Fn: func(receiver Object) builtinMethodBody {
			return func(vm *VM, args []Object, blockFrame *callFrame) Object {
				libName := args[0].(*StringObject).Value
				initFunc, ok := standardLibraris[libName]

				if !ok {
					msg := "Can't require \"" + libName + "\""
					vm.returnError(msg)
				}

				initFunc(vm)

				return TRUE
			}
		},
	},
	{
		Name: "require_relative",
		Fn: func(receiver Object) builtinMethodBody {
			return func(vm *VM, args []Object, blockFrame *callFrame) Object {
				filepath := args[0].(*StringObject).Value
				filepath = path.Join(vm.fileDir, filepath)

				file, err := ioutil.ReadFile(filepath + ".ro")

				if err != nil {
					vm.returnError(err.Error())
				}

				vm.execRequiredFile(filepath, file)

				return TRUE
			}
		},
	},
	{
		Name: "puts",
		Fn: func(receiver Object) builtinMethodBody {
			return func(vm *VM, args []Object, blockFrame *callFrame) Object {

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
			return func(vm *VM, args []Object, blockFrame *callFrame) Object {

				switch r := receiver.(type) {
				case BaseObject:
					return r.returnClass()
				case Class:
					return r.returnClass()
				default:
					return &Error{Message: "Can't call class on %T" + string(r.objectType())}
				}
			}
		},
	},
	{
		Name: "!",
		Fn: func(receiver Object) builtinMethodBody {
			return func(vm *VM, args []Object, blockFrame *callFrame) Object {

				return FALSE
			}
		},
	},
}

var BuiltinClassMethods = []*BuiltInMethod{
	{
		Name: "include",
		Fn: func(receiver Object) builtinMethodBody {
			return func(vm *VM, args []Object, blockFrame *callFrame) Object {
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
			return func(vm *VM, args []Object, blockFrame *callFrame) Object {

				class := receiver.(*RClass)
				instance := class.initializeInstance()
				initMethod := class.lookupInstanceMethod("initialize")

				if initMethod != nil {
					instance.InitializeMethod = initMethod.(*Method)
				}

				return instance
			}
		},
	},
	{
		Name: "name",
		Fn: func(receiver Object) builtinMethodBody {
			return func(vm *VM, args []Object, blockFrame *callFrame) Object {

				name := receiver.(Class).ReturnName()
				nameString := initializeString(name)
				return nameString
			}
		},
	},
	{
		Name: "superclass",
		Fn: func(receiver Object) builtinMethodBody {
			return func(vm *VM, args []Object, blockFrame *callFrame) Object {

				return receiver.(Class).returnSuperClass()
			}
		},
	},
}

func (vm *VM) execRequiredFile(filepath string, file []byte) {
	program := parser.BuildAST(file)
	g := bytecode.NewGenerator(program)
	bytecodes := g.GenerateByteCode(program)

	oldMethodTable := isTable{}
	oldClassTable := isTable{}

	// Copy current file's instruction sets.
	for name, is := range vm.isTables[bytecode.LabelDef] {
		oldMethodTable[name] = is
	}

	for name, is := range vm.isTables[bytecode.LabelDefClass] {
		oldClassTable[name] = is
	}

	// This creates new execution environments for required file, including new instruction set table.
	// So we need to copy old instruction sets and restore them later, otherwise current program's instruction set would be overwrite.
	vm.ExecBytecodes(bytecodes, filepath)

	// Restore instruction sets.
	vm.isTables[bytecode.LabelDef] = oldMethodTable
	vm.isTables[bytecode.LabelDefClass] = oldClassTable
}
