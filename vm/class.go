package vm

import (
	"fmt"
	"io/ioutil"
	"path"
	"reflect"
	"time"

	"github.com/goby-lang/goby/vm/classes"
	"github.com/goby-lang/goby/vm/errors"
)

// RClass represents normal (not built in) class object
type RClass struct {
	// Name is the class's name
	Name string
	// Methods contains its instances' methods
	Methods *environment
	// pseudoSuperClass points to the class it inherits
	pseudoSuperClass *RClass
	// This is the class where we should looking for a method.
	// It can be normal class, singleton class or a module.
	superClass *RClass
	// Class points to this class's class, which should be ClassClass
	isSingleton bool
	isModule    bool
	constants   map[string]*Pointer
	scope       *RClass
	*baseObj
}

// Class methods --------------------------------------------------------
func builtinClassCommonClassMethods() []*BuiltinMethodObject {
	return []*BuiltinMethodObject{
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
			Fn: func(receiver Object, instruction *instruction) builtinMethodBody {
				return func(t *thread, args []Object, blockFrame *normalCallFrame) Object {
					r := receiver.(*RClass)
					r.setAttrAccessor(args)

					return r
				}
			},
		},
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
			Fn: func(receiver Object, instruction *instruction) builtinMethodBody {
				return func(t *thread, args []Object, blockFrame *normalCallFrame) Object {
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
			Fn: func(receiver Object, instruction *instruction) builtinMethodBody {
				return func(t *thread, args []Object, blockFrame *normalCallFrame) Object {
					r := receiver.(*RClass)
					r.setAttrWriter(args)

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
			//   def ten
			//     10
			//   end
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
			// puts(a.ten) # => ten (overridden)
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
			Fn: func(receiver Object, instruction *instruction) builtinMethodBody {
				return func(t *thread, args []Object, blockFrame *normalCallFrame) Object {
					var class *RClass
					module, ok := args[0].(*RClass)

					if !ok {
						return t.vm.initErrorObject(errors.TypeError, instruction, "Expect argument to be a module. got=%v", args[0].Class().Name)
					}

					switch r := receiver.(type) {
					case *RClass:
						class = r
					default:
						class = r.SingletonClass()
					}

					if class.alreadyInherit(module) {
						return class
					}

					module.superClass = class.superClass
					class.superClass = module

					return class
				}
			},
		},
		{
			Name: "extend",
			Fn: func(receiver Object, instruction *instruction) builtinMethodBody {
				return func(t *thread, args []Object, blockFrame *normalCallFrame) Object {
					var class *RClass
					module, ok := args[0].(*RClass)

					if !ok {
						return t.vm.initErrorObject(errors.TypeError, instruction, "Expect argument to be a module. got=%v", args[0].Class().Name)
					}

					class = receiver.SingletonClass()

					if class.alreadyInherit(module) {
						return class
					}

					module.superClass = class.superClass
					class.superClass = module

					return class
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
			Fn: func(receiver Object, instruction *instruction) builtinMethodBody {
				return func(t *thread, args []Object, blockFrame *normalCallFrame) Object {
					if len(args) != 0 {
						return t.vm.initErrorObject(errors.ArgumentError, instruction, "Expect 0 argument. got: %d", len(args))
					}

					n, ok := receiver.(*RClass)

					if !ok {
						return t.vm.initErrorObject(errors.UndefinedMethodError, instruction, "Undefined Method '%s' for %s", "#name", receiver.toString())
					}

					name := n.ReturnName()
					nameString := t.vm.initStringObject(name)
					return nameString
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
			Fn: func(receiver Object, instruction *instruction) builtinMethodBody {
				return func(t *thread, args []Object, blockFrame *normalCallFrame) Object {
					class, ok := receiver.(*RClass)

					if !ok {
						return t.initUnsupportedMethodError(instruction, "#new", receiver)
					}

					instance := class.initializeInstance()
					initMethod := class.lookupMethod("initialize")

					if initMethod != nil {
						instance.InitializeMethod = initMethod.(*MethodObject)
					}

					return instance
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
			Fn: func(receiver Object, instruction *instruction) builtinMethodBody {
				return func(t *thread, args []Object, blockFrame *normalCallFrame) Object {
					if len(args) != 0 {
						return t.vm.initErrorObject(errors.ArgumentError, instruction, "Expect 0 argument. got: %d", len(args))
					}

					c, ok := receiver.(*RClass)

					if !ok {
						return t.vm.initErrorObject(errors.UndefinedMethodError, instruction, "Undefined Method '%s' for %s", "#superclass", receiver.toString())
					}

					superClass := c.returnSuperClass()

					if superClass == nil {
						return NULL
					}

					return superClass
				}
			},
		},
		{
			Name: "ancestors",
			Fn: func(receiver Object, instruction *instruction) builtinMethodBody {
				return func(t *thread, args []Object, blockFrame *normalCallFrame) Object {
					c, ok := receiver.(*RClass)

					if !ok {
						return t.vm.initErrorObject(errors.UndefinedMethodError, instruction, "Undefined Method '%s' for %s", "#ancestors", receiver.toString())
					}

					a := c.ancestors()
					ancestors := make([]Object, len(a))
					for i := range a {
						ancestors[i] = a[i]
					}
					return t.vm.initArrayObject(ancestors)
				}
			},
		},
	}
}

// Instance methods -----------------------------------------------------
func builtinClassCommonInstanceMethods() []*BuiltinMethodObject {
	return []*BuiltinMethodObject{
		{
			Name: "singleton_class",
			Fn: func(receiver Object, instruction *instruction) builtinMethodBody {
				return func(t *thread, args []Object, blockFrame *normalCallFrame) Object {
					return receiver.SingletonClass()
				}
			},
		},
		{
			// General method for comparing equalty of the objects
			//
			// ```ruby
			// 123 == 123   # => true
			// 123 == "123" # => false
			//
			// # Hash will not concern about the key-value pair order
			// { a: 1, b: 2 } == { a: 1, b: 2 } # => true
			// { a: 1, b: 2 } == { b: 2, a: 1 } # => true
			//
			// # Hash key will be override if the key duplicated
			// { a: 1, b: 2 } == { a: 2, b: 2, a: 1 } # => true
			// { a: 1, b: 2 } == { a: 1, b: 2, a: 2 } # => false
			//
			// # Array will concern about the order of the elements
			// [1, 2, 3] == [1, 2, 3] # => true
			// [1, 2, 3] == [3, 2, 1] # => false
			// ```
			//
			// @return [@boolean]
			Name: "==",
			Fn: func(receiver Object, instruction *instruction) builtinMethodBody {
				return func(t *thread, args []Object, blockFrame *normalCallFrame) Object {
					className := receiver.Class().Name
					compareClassName := args[0].Class().Name

					if className == compareClassName && reflect.DeepEqual(receiver, args[0]) {
						return TRUE
					}
					return FALSE
				}
			},
		}, {
			// General method for comparing inequalty of the objects
			//
			// ```ruby
			// 123 != 123   # => false
			// 123 != "123" # => true
			//
			// # Hash will not concern about the key-value pair order
			// { a: 1, b: 2 } != { a: 1, b: 2 } # => false
			// { a: 1, b: 2 } != { b: 2, a: 1 } # => false
			//
			// # Hash key will be override if the key duplicated
			// { a: 1, b: 2 } != { a: 2, b: 2, a: 1 } # => false
			// { a: 1, b: 2 } != { a: 1, b: 2, a: 2 } # => true
			//
			// # Array will concern about the order of the elements
			// [1, 2, 3] != [1, 2, 3] # => false
			// [1, 2, 3] != [3, 2, 1] # => true
			// ```
			//
			// @return [@boolean]
			Name: "!=",
			Fn: func(receiver Object, instruction *instruction) builtinMethodBody {
				return func(t *thread, args []Object, blockFrame *normalCallFrame) Object {
					className := receiver.Class().Name
					compareClassName := args[0].Class().Name

					if className == compareClassName && reflect.DeepEqual(receiver, args[0]) {
						return FALSE
					}
					return TRUE
				}
			},
		},
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
			Fn: func(receiver Object, instruction *instruction) builtinMethodBody {
				return func(t *thread, args []Object, blockFrame *normalCallFrame) Object {
					libName := args[0].(*StringObject).value
					initFunc, ok := standardLibraries[libName]

					if !ok {
						return t.vm.initErrorObject(errors.InternalError, instruction, "Can't require \"%s\"", libName)
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
			Fn: func(receiver Object, instruction *instruction) builtinMethodBody {
				return func(t *thread, args []Object, blockFrame *normalCallFrame) Object {
					callerDir := path.Dir(t.vm.currentFilePath())
					filepath := args[0].(*StringObject).value

					filepath = path.Join(callerDir, filepath)

					file, err := ioutil.ReadFile(filepath + ".gb")

					if err != nil {
						return t.vm.initErrorObject(errors.InternalError, instruction, err.Error())
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
			Fn: func(receiver Object, instruction *instruction) builtinMethodBody {
				return func(t *thread, args []Object, blockFrame *normalCallFrame) Object {

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
			Fn: func(receiver Object, instruction *instruction) builtinMethodBody {
				return func(t *thread, args []Object, blockFrame *normalCallFrame) Object {

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
			Fn: func(receiver Object, instruction *instruction) builtinMethodBody {
				return func(t *thread, args []Object, blockFrame *normalCallFrame) Object {

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
			Fn: func(receiver Object, instruction *instruction) builtinMethodBody {
				return func(t *thread, args []Object, blockFrame *normalCallFrame) Object {
					if len(args) != 1 {
						return t.vm.initErrorObject(errors.ArgumentError, instruction, "Expect 1 argument. got: %d", len(args))
					}

					int := args[0].(*IntegerObject)
					seconds := int.value
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
			Fn: func(receiver Object, instruction *instruction) builtinMethodBody {
				return func(t *thread, args []Object, blockFrame *normalCallFrame) Object {
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
			//     if block_given?
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
			Name: "block_given?",
			Fn: func(receiver Object, instruction *instruction) builtinMethodBody {
				return func(t *thread, args []Object, blockFrame *normalCallFrame) Object {
					cf := t.callFrameStack.callFrames[t.cfp-2]

					if cf.BlockFrame() == nil {
						return FALSE
					}

					return TRUE
				}
			},
		},
		{
			Name: "send",
			Fn: func(receiver Object, instruction *instruction) builtinMethodBody {
				return func(t *thread, args []Object, blockFrame *normalCallFrame) Object {
					if len(args) < 1 {
						return t.vm.initErrorObject(errors.ArgumentError, instruction, "no method name given")
					}

					name, ok := args[0].(*StringObject)

					if !ok {
						return t.vm.initErrorObject(errors.TypeError, instruction, errors.WrongArgumentTypeFormat, classes.StringClass, args[0].Class().Name)
					}

					t.sendMethod(name.value, len(args), blockFrame, instruction)

					return t.stack.top().Target
				}
			},
		},
		{
			Name: "thread",
			Fn: func(receiver Object, instruction *instruction) builtinMethodBody {
				return func(t *thread, args []Object, blockFrame *normalCallFrame) Object {
					if blockFrame == nil {
						return t.vm.initErrorObject(errors.InternalError, instruction, errors.CantYieldWithoutBlockFormat)
					}

					newT := t.vm.newThread()

					go func() {
						newT.builtinMethodYield(blockFrame, args...)
					}()

					// We need to pop this frame from main thread manually,
					// because the block's 'leave' instruction is running on other process
					t.callFrameStack.pop()

					return NULL
				}
			},
		},
		{
			// Returns true if Object class is equal to the input argument class
			//
			// ```ruby
			// "Hello".is_a?(String)            # => true
			// 123.is_a?(Integer)               # => true
			// [1, true, "String"].is_a?(Array) # => true
			// { a: 1, b: 2 }.is_a?(Hash)       # => true
			// "Hello".is_a?(Integer)           # => false
			// 123.is_a?(Range)                 # => false
			// (2..4).is_a?(Hash)               # => false
			// nil.is_a?(Integer)               # => false
			// ```
			//
			// @param n/a []
			// @return [Boolean]
			Name: "is_a?",
			Fn: func(receiver Object, instruction *instruction) builtinMethodBody {
				return func(t *thread, args []Object, blockFrame *normalCallFrame) Object {
					if len(args) != 1 {
						return t.vm.initErrorObject(errors.ArgumentError, instruction, "Expect 1 argument. got: %d", len(args))
					}

					c := args[0]
					gobyClass, ok := c.(*RClass)

					if !ok {
						return t.vm.initErrorObject(errors.TypeError, instruction, errors.WrongArgumentTypeFormat, classes.ClassClass, c.Class().Name)
					}

					receiverClass := receiver.Class()

					for {
						if receiverClass.Name == gobyClass.Name {
							return TRUE
						}

						if receiverClass.Name == classes.ObjectClass {
							break
						}

						receiverClass = receiverClass.superClass
					}
					return FALSE
				}
			},
		},
		{
			// Returns true if Object is nil
			//
			// ```ruby
			// 123.nil?            # => false
			// "String".nil?       # => false
			// { a: 1, b: 2 }.nil? # => false
			// (3..5).nil?         # => false
			// nil.nil?            # => true  (See the implementation of Null#nil? in vm/null.go file)
			// ```
			//
			// @param n/a []
			// @return [Boolean]
			Name: "nil?",
			Fn: func(receiver Object, instruction *instruction) builtinMethodBody {
				return func(t *thread, args []Object, blockFrame *normalCallFrame) Object {
					if len(args) != 0 {
						return t.vm.initErrorObject(errors.ArgumentError, instruction, "Expect 0 argument. got: %d", len(args))
					}
					return FALSE
				}
			},
		},
		{
			Name: "instance_variable_get",
			Fn: func(receiver Object, instruction *instruction) builtinMethodBody {
				return func(t *thread, args []Object, blockFrame *normalCallFrame) Object {
					arg, isStr := args[0].(*StringObject)

					if !isStr {
						return t.vm.initErrorObject(errors.TypeError, instruction, errors.WrongArgumentTypeFormat, classes.StringClass, args[0].Class().Name)
					}

					obj, ok := receiver.instanceVariableGet(arg.value)

					if !ok {
						return NULL
					}

					return obj
				}
			},
		},
		{
			Name: "instance_variable_set",
			Fn: func(receiver Object, instruction *instruction) builtinMethodBody {
				return func(t *thread, args []Object, blockFrame *normalCallFrame) Object {
					if len(args) != 2 {
						return t.vm.initErrorObject(errors.ArgumentError, instruction, "Expect 2 arguments. got: %d", len(args))
					}

					argName, isStr := args[0].(*StringObject)
					obj := args[1]

					if !isStr {
						return t.vm.initErrorObject(errors.TypeError, instruction, errors.WrongArgumentTypeFormat, classes.StringClass, args[0].Class().Name)
					}

					receiver.instanceVariableSet(argName.value, obj)

					return obj
				}
			},
		},
		{
			Name: "methods",
			Fn: func(receiver Object, instruction *instruction) builtinMethodBody {
				return func(t *thread, args []Object, blockFrame *normalCallFrame) Object {
					methods := []Object{}
					set := map[string]interface{}{}
					klasses := receiver.Class().ancestors()
					if receiver.SingletonClass() != nil {
						klasses = append([]*RClass{receiver.SingletonClass()}, klasses...)
					}
					for _, klass := range klasses {
						for _, name := range klass.Methods.names() {
							if set[name] == nil {
								set[name] = true
								methods = append(methods, t.vm.initStringObject(name))
							}
						}
					}
					return t.vm.initArrayObject(methods)
				}
			},
		},
	}
}

// Internal functions ===================================================

// Functions for initialization -----------------------------------------

// initializeClass is a common function for vm, which initializes and returns
// a class instance with given class name.
func (vm *VM) initializeClass(name string, isModule bool) *RClass {
	class := vm.createRClass(name)
	class.isModule = isModule
	singletonClass := vm.createRClass(fmt.Sprintf("#<Class:%s>", name))
	singletonClass.isSingleton = true
	class.singletonClass = singletonClass
	class.inherits(vm.objectClass)

	return class
}

func (vm *VM) createRClass(className string) *RClass {
	objectClass := vm.objectClass
	classClass := vm.topLevelClass(classes.ClassClass)

	return &RClass{
		Name:             className,
		Methods:          newEnvironment(),
		pseudoSuperClass: objectClass,
		superClass:       objectClass,
		constants:        make(map[string]*Pointer),
		isModule:         false,
		baseObj:          &baseObj{class: classClass, InstanceVariables: newEnvironment()},
	}
}

func initClassClass() *RClass {
	classClass := &RClass{
		Name:      classes.ClassClass,
		Methods:   newEnvironment(),
		constants: make(map[string]*Pointer),
		baseObj:   &baseObj{},
	}

	singletonClass := &RClass{
		Name:        "#<Class:Class>",
		Methods:     newEnvironment(),
		constants:   make(map[string]*Pointer),
		isModule:    false,
		baseObj:     &baseObj{class: classClass, InstanceVariables: newEnvironment()},
		isSingleton: true,
	}

	classClass.class = classClass
	classClass.singletonClass = singletonClass

	classClass.setBuiltinMethods(builtinClassCommonClassMethods(), true)

	return classClass
}

func initObjectClass(c *RClass) *RClass {
	objectClass := &RClass{
		Name:      classes.ObjectClass,
		Methods:   newEnvironment(),
		constants: make(map[string]*Pointer),
		baseObj:   &baseObj{class: c},
	}

	singletonClass := &RClass{
		Name:        "#<Class:Object>",
		Methods:     newEnvironment(),
		constants:   make(map[string]*Pointer),
		isModule:    false,
		baseObj:     &baseObj{class: c, InstanceVariables: newEnvironment()},
		isSingleton: true,
		superClass:  c,
	}

	objectClass.singletonClass = singletonClass
	objectClass.superClass = objectClass
	objectClass.pseudoSuperClass = objectClass
	c.inherits(objectClass)

	objectClass.setBuiltinMethods(builtinClassCommonInstanceMethods(), true)
	objectClass.setBuiltinMethods(builtinClassCommonInstanceMethods(), false)

	return objectClass
}

// Polymorphic helper functions -----------------------------------------

// TODO: Remove the redundant functions

// ReturnName returns the object's name as the string format
func (c *RClass) ReturnName() string {
	return c.Name
}

// TODO: Singleton class's inspect() should also mark if it's a singleton class explicitly.

// toString returns the object's name as the string format
func (c *RClass) toString() string {
	return c.Name
}

// toJSON just delegates to `toString`
func (c *RClass) toJSON() string {
	return c.toString()
}

// Value returns class itself
func (c *RClass) Value() interface{} {
	return c
}

func (c *RClass) inherits(sc *RClass) {
	c.superClass = sc
	c.pseudoSuperClass = sc
	c.singletonClass.superClass = sc.singletonClass
	c.singletonClass.pseudoSuperClass = sc.singletonClass
}

func (c *RClass) setBuiltinMethods(methodList []*BuiltinMethodObject, classMethods bool) {
	for _, m := range methodList {
		c.Methods.set(m.Name, m)
	}

	if classMethods {
		for _, m := range methodList {
			c.singletonClass.Methods.set(m.Name, m)
		}
	}
}

func (c *RClass) findMethod(methodName string) (method Object) {
	if c.isSingleton {
		method = c.superClass.lookupMethod(methodName)
	} else {
		method = c.SingletonClass().lookupMethod(methodName)
	}

	return
}

func (c *RClass) lookupMethod(methodName string) Object {
	method, ok := c.Methods.get(methodName)

	if !ok {
		if c.superClass != nil && c.superClass != c {
			if c.Name == classes.ClassClass {
				return nil
			}

			return c.superClass.lookupMethod(methodName)
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

		if c.superClass != nil && c.Name != classes.ObjectClass {
			return c.superClass.lookupConstant(constName, false)
		}

		return nil
	}

	return constant
}

func (c *RClass) setClassConstant(constant *RClass) {
	c.constants[constant.Name] = &Pointer{Target: constant}
}

func (c *RClass) getClassConstant(constName string) (class *RClass) {
	t := c.constants[constName].Target
	class, ok := t.(*RClass)

	if ok {
		return
	}

	panic(constName + " is not a class.")
}

func (c *RClass) alreadyInherit(constant *RClass) bool {
	if c.superClass == constant {
		return true
	}

	if c.superClass.Name == classes.ObjectClass {
		return false
	}

	return c.superClass.alreadyInherit(constant)
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
			attrName := attr.(*StringObject).value
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
			attrName := attr.(*StringObject).value
			c.Methods.set(attrName, generateAttrReadMethod(attrName))
		}
	case []string:
		for _, attrName := range args {
			c.Methods.set(attrName, generateAttrReadMethod(attrName))
		}
	case string:
		c.Methods.set(args, generateAttrReadMethod(args))
	}

}

func (c *RClass) setAttrAccessor(args interface{}) {
	c.setAttrReader(args)
	c.setAttrWriter(args)
}

func (c *RClass) ancestors() []*RClass {
	klasses := []*RClass{c}
	for {
		if c.Name == classes.ObjectClass {
			break
		}
		c = c.superClass
		klasses = append(klasses, c)
	}

	return klasses
}

// Other helper functions -----------------------------------------------

func generateAttrWriteMethod(attrName string) *BuiltinMethodObject {
	return &BuiltinMethodObject{
		Name: attrName + "=",
		Fn: func(receiver Object, instruction *instruction) builtinMethodBody {
			return func(t *thread, args []Object, blockFrame *normalCallFrame) Object {
				v := receiver.instanceVariableSet("@"+attrName, args[0])
				return v
			}
		},
	}
}

func generateAttrReadMethod(attrName string) *BuiltinMethodObject {
	return &BuiltinMethodObject{
		Name: attrName,
		Fn: func(receiver Object, instruction *instruction) builtinMethodBody {
			return func(t *thread, args []Object, blockFrame *normalCallFrame) Object {
				v, ok := receiver.instanceVariableGet("@" + attrName)

				if ok {
					return v
				}

				return NULL
			}
		},
	}
}
