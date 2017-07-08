package vm

import (
	"fmt"
	"io/ioutil"
	"path"
	"time"
)

const (
	objectClass  = "Object"
	classClass   = "Class"
	integerClass = "Integer"
	stringClass  = "String"
	arrayClass   = "Array"
	hashClass    = "Hash"
	booleanClass = "Boolean"
	nullClass    = "Null"
	channelClass = "Channel"
	rangeClass   = "Range"
	methodClass  = "method"
)

// RClass represents normal (not built in) class object
type RClass struct {
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
	class *RClass
	// Singleton is a flag marks if this class a singleton class
	Singleton bool
	isModule  bool
	constants map[string]*Pointer
	scope     *RClass
	*baseObj
}

func initClassClass() *RClass {
	classClass := &RClass{
		Name:         classClass,
		Methods:      newEnvironment(),
		ClassMethods: newEnvironment(),
		constants:    make(map[string]*Pointer),
	}

	classClass.setBuiltInMethods(builtinCommonInstanceMethods(), false)
	classClass.setBuiltInMethods(builtinCommonInstanceMethods(), true)
	classClass.setBuiltInMethods(builtinClassClassMethods(), true)

	return classClass
}

func initObjectClass(c *RClass) *RClass {
	objectClass := &RClass{
		Name:         objectClass,
		class:        c,
		ClassMethods: newEnvironment(),
		Methods:      newEnvironment(),
		constants:    make(map[string]*Pointer),
	}

	objectClass.setBuiltInMethods(builtinCommonInstanceMethods(), false)

	return objectClass
}

// initializeClass initializes and returns a class instance with given class name
func (vm *VM) initializeClass(name string, isModule bool) *RClass {
	class := vm.createRClass(name)
	class.isModule = isModule

	return class
}

func (vm *VM) createRClass(className string) *RClass {
	classClass := vm.builtInClasses[classClass]
	objectClass := vm.builtInClasses[objectClass]

	return &RClass{
		Name:             className,
		Methods:          newEnvironment(),
		ClassMethods:     newEnvironment(),
		pseudoSuperClass: objectClass,
		superClass:       objectClass,
		constants:        make(map[string]*Pointer),
		isModule:         false,
		baseObj:          &baseObj{class: classClass, InstanceVariables: newEnvironment()},
	}
}

// toString returns the basic inspected result (which is class name) of current class
// TODO: Singleton class's inspect() should also mark if it's a singleton class explicitly.
func (c *RClass) toString() string {
	if c.isModule {
		return "<Module:" + c.Name + ">"
	}
	return "<Class:" + c.Name + ">"
}

func (c *RClass) toJSON() string {
	return c.toString()
}

func (c *RClass) setBuiltInMethods(methodList []*BuiltInMethodObject, classMethods bool) {
	for _, m := range methodList {
		c.Methods.set(m.Name, m)
	}

	if classMethods {
		for _, m := range methodList {
			c.ClassMethods.set(m.Name, m)
		}
	}
}

func (c *RClass) lookupClassMethod(methodName string) Object {
	method, ok := c.ClassMethods.get(methodName)

	if !ok {
		if c.superClass != nil {
			return c.superClass.lookupClassMethod(methodName)
		}
		if c.class != nil {
			return c.class.lookupClassMethod(methodName)
		}
		return nil
	}

	return method
}

func (c *RClass) lookupInstanceMethod(methodName string) Object {
	method, ok := c.Methods.get(methodName)

	if !ok {
		if c.superClass != nil {
			return c.superClass.lookupInstanceMethod(methodName)
		}

		if c.class != nil {
			return c.class.lookupInstanceMethod(methodName)
		}

		return nil
	}

	return method
}

func (c *RClass) lookupConstant(constName string, findInScope bool) *Pointer {
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

// ReturnName returns the name of the class
func (c *RClass) ReturnName() string {
	return c.Name
}

func (c *RClass) returnSuperClass() *RClass {
	return c.pseudoSuperClass
}

func (c *RClass) initializeInstance() *RObject {
	instance := &RObject{baseObj: &baseObj{class: c, InstanceVariables: newEnvironment()}}

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

func builtinCommonInstanceMethods() []*BuiltInMethodObject {
	return []*BuiltInMethodObject{
		{
			// Loads the given Goby library name without extension (mainly for modules), returning `true`
			// if successful and `false` if the feature is already loaded.
			//
			// Currently, only the following embedded Goby libraries are targeted:
			//
			// - "file"
			// - "net/http"
			// - "net/simple_server"
			// - "uri"
			//
			// ```ruby
			// require("file")
			// File.extname("foo.rb")
			// ```
			//
			// TBD: the load paths for `require`
			//
			// @param filename [String] Quoted file name of the library, without extension
			// @return [Boolean] Result of loading module
			Name: "require",
			Fn: func(receiver Object) builtinMethodBody {
				return func(t *thread, args []Object, blockFrame *callFrame) Object {
					libName := args[0].(*StringObject).Value
					initFunc, ok := standardLibraries[libName]

					if !ok {
						return initErrorObject(InternalErrorClass, "Can't require \"%s\"", libName)
					}

					initFunc(t.vm)

					return TRUE
				}
			},
		},
		{
			// Loads the Goby library (mainly for modules) from the given local path plus name
			// without extension from the current directory, returning `true` if successful,
			// and `false` if the feature is already loaded.
			//
			// ```ruby
			// require_relative("../test_fixtures/require_test/foo")
			// fifty = Foo.bar(5)
			// ```
			//
			// @param path/name [String] Quoted file path to library plus name, without extension
			// @return [Boolean] Result of loading module
			Name: "require_relative",
			Fn: func(receiver Object) builtinMethodBody {
				return func(t *thread, args []Object, blockFrame *callFrame) Object {
					callerDir := path.Dir(t.vm.currentFilePath())
					filepath := args[0].(*StringObject).Value

					filepath = path.Join(callerDir, filepath)

					file, err := ioutil.ReadFile(filepath + ".gb")

					if err != nil {
						return initErrorObject(InternalErrorClass, err.Error())
					}

					t.vm.execRequiredFile(filepath, file)

					return TRUE
				}
			},
		},
		{
			// Puts string literals or objects into stdout with a tailing line feed, converting into String
			// if needed.
			//
			// ```ruby
			// puts("foo", "bar")
			// # => foo
			// # => bar
			// puts("baz", String.name)
			// # => baz
			// # => String
			// puts("foo" + "bar")
			// # => foobar
			// ```
			// TODO: interpolation is needed to be implemented.
			//
			// @param *args [Class] String literals, or other objects that can be converted into String.
			// @return [Null]
			Name: "puts",
			Fn: func(receiver Object) builtinMethodBody {
				return func(t *thread, args []Object, blockFrame *callFrame) Object {

					for _, arg := range args {
						fmt.Println(arg.toString())
					}

					return NULL
				}
			},
		},
		{
			// Returns the class of the object. Receiver cannot be omitted.
			//
			// FYI: You can convert the class into String with `#name`.
			//
			// ```ruby
			// puts(100.class)         # => <Class:Integer>
			// puts(100.class.name)    # => Integer
			// puts("123".class)       # => <Class:String>
			// puts("123".class.name)  # => String
			// ```
			//
			// @param object [Object] Receiver (required)
			// @return [Class] The class of the receiver
			Name: "class",
			Fn: func(receiver Object) builtinMethodBody {
				return func(t *thread, args []Object, blockFrame *callFrame) Object {

					switch r := receiver.(type) {
					case Object:
						return r.Class()
					default:
						return &Error{Message: "Can't call class on %T" + string(r.Class().ReturnName())}
					}
				}
			},
		},
		{
			// Inverts the boolean value.
			//
			// ```ruby
			// !true  # => false
			// !false # => true
			// ```
			//
			// @param object [Object] object that return boolean value to invert
			// @return [Object] Inverted boolean value
			Name: "!",
			Fn: func(receiver Object) builtinMethodBody {
				return func(t *thread, args []Object, blockFrame *callFrame) Object {

					return FALSE
				}
			},
		},
		{
			// Suspends the current thread for duration (sec).
			//
			// **Note:** currently, parameter cannot be omitted, and only Integer can be specified.
			//
			// ```ruby
			// a = sleep(2)
			// puts(a)     # => 2
			// ```
			//
			// @param sec [Integer] time to wait in sec
			// @return [Integer] actual time slept in sec
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
		{
			// Returns object's string representation.
			// @param n/a []
			// @return [String] Object's string representation.
			Name: "to_s",
			Fn: func(receiver Object) builtinMethodBody {
				return func(t *thread, args []Object, blockFrame *callFrame) Object {
					return t.vm.initStringObject(receiver.toString())
				}
			},
		},
		{
			// Returns true if a block is given in the current context and `yield` is ready to call.
			//
			// **Note:** The method name does not end with '?' because the sign is unavalable in Goby for now.
			//
			// ```ruby
			// class File
			//   def self.open(filename, mode, perm)
			//     file = new(filename, mode, perm)
			//
			//     if block_given
			//       yield(file)
			//     end
			//
			//     file.close
			//   end
			// end
			// ```
			//
			// @param n/a []
			// @return [Boolean] true/false
			Name: "block_given",
			Fn: func(receiver Object) builtinMethodBody {
				return func(t *thread, args []Object, blockFrame *callFrame) Object {
					cf := t.callFrameStack.top()

					if cf.blockFrame == nil {
						return FALSE
					}

					return TRUE
				}
			},
		},
		{
			Name: "thread",
			Fn: func(receiver Object) builtinMethodBody {
				return func(t *thread, args []Object, blockFrame *callFrame) Object {
					newT := t.vm.newThread()

					go func() {
						newT.builtInMethodYield(blockFrame, args...)
					}()

					// We need to pop this frame from main thread manually,
					// because the block's 'leave' instruction is running on other process
					t.callFrameStack.pop()

					return NULL
				}
			},
		},
	}
}

func builtinClassClassMethods() []*BuiltInMethodObject {
	return []*BuiltInMethodObject{
		{
			// Creates instance variables and corresponding methods that return the value of each
			// instance variable.
			//
			// Only string literal can be used for now.
			//
			// ```ruby
			// class Foo
			//   attr_reader("bar", "buz")
			// end
			// ```
			// is equivalent to:
			//
			// ```ruby
			// class Foo
			//   def bar
			//     @bar
			//   end
			//   def buz
			//     @buz
			//   end
			// end
			// ```
			//
			// @param *args [String] One or more quoted method names for 'getter'
			// @return [Null]
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
			// Creates instance variables and corresponding methods that assign an argument to each
			// instance variable. No return value.
			//
			// Only string literal can be used for now.
			//
			// ```ruby
			// class Foo
			//   attr_writer("bar", "buz")
			// end
			// ```
			// is equivalent to:
			//
			// ```ruby
			// class Foo
			//   def bar=(val)
			//     @bar = val
			//   end
			//   def buz=(val)
			//     @buz = val
			//   end
			// end
			// ```
			//
			// @param *args [String] One or more quoted method names for 'setter'
			// @return [Null]
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
			// Creates instance variables and corresponding methods that return the value of
			// each instance variable and assign an argument to each instance variable.
			// Only string literal can be used for now.
			//
			// ```ruby
			// class Foo
			//   attr_accessor("bar", "buz")
			// end
			// ```
			// is equivalent to:
			//
			// ```ruby
			// class Foo
			//   def bar
			//     @bar
			//   end
			//   def buz
			//     @buz
			//   end
			//   def bar=(val)
			//     @bar = val
			//   end
			//   def buz=(val)
			//     @buz = val
			//   end
			// end
			// ```
			//
			// @param *args [String] One or more quoted method names for 'getter/setter'
			// @return [Null]
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
			// Includes a module for mixin, which inherits only methods and constants from the module.
			// The included module is inserted into the path of the inheritance tree, between the class
			// and the superclass so that the methods of the module is prioritized to superclasses.
			//
			// The order of `include` affects: the modules that included later are prioritized.
			// If multiple modules include the same methods, the method will only come from
			// the last included module.
			//
			// ```ruby
			// module Foo
			// def ten
			//    10
			// end
			// end
			//
			// module Bar
			//   def ten
			//     "ten"
			//   end
			// end
			//
			// class Baz
			//   include(Foo)
			//   include(Bar) # method `ten` is only included from this module
			// end
			//
			// a = Baz.new
			// puts(a.ten) # => ten (overriden)
			// ```
			//
			// **Note**:
			//
			// You cannot use string literal, or pass two or more arguments to `include`.
			//
			// ```ruby
			//   include("Foo")    # => error
			//   include(Foo, Bar) # => error
			// ```
			//
			// Including modules into built-in classes such as String are not supported:
			//
			// ```ruby
			// module Foo
			//   def ten
			//     10
			//   end
			// end
			// class String
			//   include(Foo) # => error
			// end
			// ```
			//
			// @param module [Class] Module name to include
			// @return [Null]
			Name: "include",
			Fn: func(receiver Object) builtinMethodBody {
				return func(t *thread, args []Object, blockFrame *callFrame) Object {
					var class *RClass
					module := args[0].(*RClass)

					switch r := receiver.(type) {
					case *RClass:
						class = r
					case *RObject:
						objectClass := t.vm.builtInClasses["Object"]

						if r.class == objectClass {
							class = objectClass
						}
					}

					module.superClass = class.superClass
					class.superClass = module

					return class
				}
			},
		},
		{
			// Creates and returns a new anonymous class from a receiver.
			// You can use any classes you defined as the receiver:
			//
			// ```ruby
			// class Foo
			// end
			// a = Foo.new
			// ```
			//
			// Note that the built-in classes such as Class or String are not open for creating instances
			// and you can't call `new` against them.
			//
			// ```ruby
			// a = Class.new  # => error
			// a = String.new # => error
			// ```
			// @param class [Class] Receiver
			// @return [Object] Created object
			Name: "new",
			Fn: func(receiver Object) builtinMethodBody {
				return func(t *thread, args []Object, blockFrame *callFrame) Object {
					class, ok := receiver.(*RClass)

					if !ok {
						return t.UnsupportedMethodError("#new", receiver)
					}

					if class.pseudoSuperClass.isModule {
						return initErrorObject(InternalErrorClass, "Module inheritance is not supported: %s", class.pseudoSuperClass.Name)
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
			// Returns the name of the class (receiver).
			//
			// ```ruby
			// puts(Array.name)  # => Array
			// puts(Class.name)  # => Class
			// puts(Object.name) # => Object
			// ```
			// @param class [Class] Receiver
			// @return [String] Converted receiver name
			Name: "name",
			Fn: func(receiver Object) builtinMethodBody {
				return func(t *thread, args []Object, blockFrame *callFrame) Object {

					name := receiver.(*RClass).ReturnName()
					nameString := t.vm.initStringObject(name)
					return nameString
				}
			},
		},
		{
			// Returns the superclass object of the receiver.
			//
			// ```ruby
			// puts(Array.superclass)  # => <Class:Object>
			// puts(String.superclass) # => <Class:Object>
			//
			// class Foo;end
			// class Bar < Foo
			// end
			// puts(Foo.superclass)    # => <Class:Object>
			// puts(Bar.superclass)    # => <Class:Foo>
			// ```
			//
			// **Note**: the following is not supported:
			//
			// - Class class
			//
			// - Object class
			//
			// - instance objects or object literals
			//
			// ```ruby
			// puts("string".superclass) # => error
			// puts(Class.superclass)    # => error
			// puts(Object.superclass)   # => error
			// ```
			// @param class [Class] Receiver
			// @return [Object] Superclass object of the receiver
			Name: "superclass",
			Fn: func(receiver Object) builtinMethodBody {
				return func(t *thread, args []Object, blockFrame *callFrame) Object {
					c := receiver.(*RClass).returnSuperClass()

					if c == nil {
						return NULL
					}

					return c
				}
			},
		},
	}
}
