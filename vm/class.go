package vm

import (
	"fmt"
	"os"
	"path"
	"reflect"
	"sync"
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
	isSingleton           bool
	isModule              bool
	constants             map[string]*Pointer
	scope                 *RClass
	inheritsMethodMissing bool
	*baseObj
}

var externalClasses = map[string][]ClassLoader{}
var externalClassLock sync.Mutex

// RegisterExternalClass will add the given class to the global registery of available classes
func RegisterExternalClass(path string, c ...ClassLoader) {
	externalClassLock.Lock()
	externalClasses[path] = c
	externalClassLock.Unlock()
}

// ClassLoader can be registerd with a vm so that it can load this library at vm creation
type ClassLoader = func(*VM) error

func buildMethods(m map[string]MethodBuilder) []*BuiltinMethodObject {
	out := make([]*BuiltinMethodObject, len(m))
	var i int
	for k, v := range m {
		out[i] = ExternalBuiltinMethod(k, v)
		i++
	}
	return out
}

// ExternalClass helps define external go classes
func ExternalClass(name, path string, classMethods, instanceMethods map[string]MethodBuilder) ClassLoader {
	return func(v *VM) error {
		pg := v.initializeClass(name, false)
		pg.setBuiltinMethods(buildMethods(classMethods), true)
		pg.setBuiltinMethods(buildMethods(instanceMethods), false)
		v.objectClass.setClassConstant(pg)

		return v.mainThread.execGobyLib(path)
	}
}

// Class methods --------------------------------------------------------
func builtinClassCommonClassMethods() []*BuiltinMethodObject {
	return []*BuiltinMethodObject{
		{
			// Returns an array that contains ancestor classes/modules of the receiver,
			// left to right.
			//
			// ```ruby
			// String.ancestors #=> [String, Object]
			//
			// module Foo
			//   def bar
			//     42
			//   end
			// end
			//
			// class Bar
			//   include Foo
			// end
			//
			// Bar.ancestors
			// #=> [Bar, Foo, Object]
			//
			// # you need `#singleton_class` to show the 'extended' modules
			// class Baz
			//   extend Foo
			// end
			//
			// Baz.singleton_class.ancestors
			// #=> [#<Class:Baz>, Foo, #<Class:Object>, Class, Object]
			// Baz.ancestors          # Foo is hidden
			// #=> [Baz, Object]
			// ```
			//
			// @param class [Class] Receiver
			// @return [Array]
			Name: "ancestors",
			Fn: func(receiver Object, sourceLine int) builtinMethodBody {
				return func(t *Thread, args []Object, blockFrame *normalCallFrame) Object {
					c, ok := receiver.(*RClass)

					if !ok {
						return t.vm.InitErrorObject(errors.UndefinedMethodError, sourceLine, "Undefined Method '%s' for %s", "#ancestors", receiver.toString())
					}

					a := c.ancestors()
					ancestors := make([]Object, len(a))
					for i := range a {
						ancestors[i] = a[i]
					}
					return t.vm.InitArrayObject(ancestors)
				}
			},
		},
		{
			// Returns true if self is an ancestor of another class/module.
			//
			// ```ruby
			// Object > Array #=> true
			// Array > Object #=> false
			// Object > Object #=> false
			// ```
			//
			// @param module [Class]
			// @return [Boolean, Null]
			Name: ">",
			Fn: func(receiver Object, sourceLine int) builtinMethodBody {
				return func(t *Thread, args []Object, blockFrame *normalCallFrame) Object {
					c, ok := receiver.(*RClass)

					if !ok {
						return t.vm.InitErrorObject(errors.UndefinedMethodError, sourceLine, "Undefined Method '%s' for %s", "#<", receiver.toString())
					}

					module, ok := args[0].(*RClass)

					if !ok {
						return t.vm.InitErrorObject(errors.TypeError, sourceLine, "Expect argument to be a module. got=%v", args[0].Class().Name)
					}

					if c == module {
						return FALSE
					}

					if module.alreadyInherit(c) {
						return TRUE
					}

					if c.alreadyInherit(module) {
						return FALSE
					}
					return NULL
				}
			},
		},
		{
			// Returns true if self is an ancestor or same class/module of another.
			//
			// ```ruby
			// Object >= Array #=> true
			// Array >= Object #=> false
			// Object >= Object #=> true
			// ```
			//
			// @param module [Class]
			// @return [Boolean, Null]
			Name: ">=",
			Fn: func(receiver Object, sourceLine int) builtinMethodBody {
				return func(t *Thread, args []Object, blockFrame *normalCallFrame) Object {
					c, ok := receiver.(*RClass)

					if !ok {
						return t.vm.InitErrorObject(errors.UndefinedMethodError, sourceLine, "Undefined Method '%s' for %s", "#<", receiver.toString())
					}

					module, ok := args[0].(*RClass)

					if !ok {
						return t.vm.InitErrorObject(errors.TypeError, sourceLine, "Expect argument to be a module. got=%v", args[0].Class().Name)
					}

					if c == module {
						return TRUE
					}

					if module.alreadyInherit(c) {
						return TRUE
					}

					if c.alreadyInherit(module) {
						return FALSE
					}
					return NULL
				}
			},
		},
		{
			// Returns true if another class/module is an ancestor of self.
			//
			// ```ruby
			// Object < Array #=> false
			// Array < Object #=> true
			// Object < Object #=> false
			// ```
			//
			// @param module [Class]
			// @return [Boolean, Null]
			Name: "<",
			Fn: func(receiver Object, sourceLine int) builtinMethodBody {
				return func(t *Thread, args []Object, blockFrame *normalCallFrame) Object {
					c, ok := receiver.(*RClass)

					if !ok {
						return t.vm.InitErrorObject(errors.UndefinedMethodError, sourceLine, "Undefined Method '%s' for %s", "#<", receiver.toString())
					}

					module, ok := args[0].(*RClass)

					if !ok {
						return t.vm.InitErrorObject(errors.TypeError, sourceLine, "Expect argument to be a module. got=%v", args[0].Class().Name)
					}

					if c == module {
						return FALSE
					}

					if module.alreadyInherit(c) {
						return FALSE
					}

					if c.alreadyInherit(module) {
						return TRUE
					}
					return NULL
				}
			},
		},
		{
			// Returns true if another is an ancestor or same class/module of self.
			//
			// ```ruby
			// Object <= Array #=> false
			// Array <= Object #=> true
			// Object <= Object #=> true
			// ```
			//
			// @param module [Class]
			// @return [Boolean, Null]
			Name: "<=",
			Fn: func(receiver Object, sourceLine int) builtinMethodBody {
				return func(t *Thread, args []Object, blockFrame *normalCallFrame) Object {
					c, ok := receiver.(*RClass)

					if !ok {
						return t.vm.InitErrorObject(errors.UndefinedMethodError, sourceLine, "Undefined Method '%s' for %s", "#<", receiver.toString())
					}

					module, ok := args[0].(*RClass)

					if !ok {
						return t.vm.InitErrorObject(errors.TypeError, sourceLine, "Expect argument to be a module. got=%v", args[0].Class().Name)
					}

					if c == module {
						return TRUE
					}

					if module.alreadyInherit(c) {
						return FALSE
					}

					if c.alreadyInherit(module) {
						return TRUE
					}
					return NULL
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
			Fn: func(receiver Object, sourceLine int) builtinMethodBody {
				return func(t *Thread, args []Object, blockFrame *normalCallFrame) Object {
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
			Fn: func(receiver Object, sourceLine int) builtinMethodBody {
				return func(t *Thread, args []Object, blockFrame *normalCallFrame) Object {
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
			Fn: func(receiver Object, sourceLine int) builtinMethodBody {
				return func(t *Thread, args []Object, blockFrame *normalCallFrame) Object {
					r := receiver.(*RClass)
					r.setAttrWriter(args)

					return r
				}
			},
		},
		// Inserts a module as a singleton class to make the module's methods class methods.
		// You can see the extended module by using `singleton_class.ancestors`
		//
		// ```ruby
		// String.ancestors #=> [String, Object]
		//
		// module Foo
		//   def bar
		//     42
		//   end
		// end
		//
		// class Bar
		//   extend Foo
		// end
		//
		// Bar.bar   #=> 42
		//
		// Bar.singleton_class.ancestors
		// #=> [#<Class:Bar>, Foo, #<Class:Object>, Class, Object]
		//
		// Bar.ancestors           # Foo is hidden
		// #=> [Bar, Object]
		// ```
		//
		// @param module [Class] Module name to extend
		// @return [Null]
		{
			Name: "extend",
			Fn: func(receiver Object, sourceLine int) builtinMethodBody {
				return func(t *Thread, args []Object, blockFrame *normalCallFrame) Object {
					var class *RClass
					module, ok := args[0].(*RClass)

					if !ok || !module.isModule {
						return t.vm.InitErrorObject(errors.TypeError, sourceLine, "Expect argument to be a module. got=%v", args[0].Class().Name)
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
			// Baz.ancestors
			// [Baz, Bar, Foo, Object]   # Bar is prioritized to Foo
			//
			// a = Baz.new
			// puts(a.ten) # => ten      # overridden
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
			// @param module [Class] Module name to include
			// @return [Null]
			Name: "include",
			Fn: func(receiver Object, sourceLine int) builtinMethodBody {
				return func(t *Thread, args []Object, blockFrame *normalCallFrame) Object {
					var class *RClass
					module, ok := args[0].(*RClass)

					if !ok || !module.isModule {
						return t.vm.InitErrorObject(errors.TypeError, sourceLine, "Expect argument to be a module. got=%v", args[0].Class().Name)
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
			Name: "inherits_method_missing",
			Fn: func(receiver Object, sourceLine int) builtinMethodBody {
				return func(t *Thread, args []Object, blockFrame *normalCallFrame) Object {
					var class *RClass

					if len(args) != 0 {
						return t.vm.InitErrorObject(errors.ArgumentError, sourceLine, "Expect 0 argument. got: %d", len(args))
					}

					switch r := receiver.(type) {
					case *RClass:
						class = r
					default:
						class = r.SingletonClass()
					}

					class.inheritsMethodMissing = true

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
			Fn: func(receiver Object, sourceLine int) builtinMethodBody {
				return func(t *Thread, args []Object, blockFrame *normalCallFrame) Object {
					if len(args) != 0 {
						return t.vm.InitErrorObject(errors.ArgumentError, sourceLine, "Expect 0 argument. got: %d", len(args))
					}

					n, ok := receiver.(*RClass)

					if !ok {
						return t.vm.InitErrorObject(errors.UndefinedMethodError, sourceLine, "Undefined Method '%s' for %s", "#name", receiver.toString())
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
			// Note that the built-in classes such as String are not open for creating instances
			// and you can't call `new` against them.
			//
			// ```ruby
			// a = String.new # => error
			// ```
			// @param class [Class] Receiver
			// @return [Object] Created object
			Name: "new",
			Fn: func(receiver Object, sourceLine int) builtinMethodBody {
				return func(t *Thread, args []Object, blockFrame *normalCallFrame) Object {
					class, ok := receiver.(*RClass)

					if !ok {
						return t.vm.initUnsupportedMethodError(sourceLine, "#new", receiver)
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
			// A predicate class method that returns `true` if the object has an ability to respond to the method, otherwise `false`.
			// Note that signs like `+` or `?` should be String literal.
			//
			// ```ruby
			// Class.respond_to? "respond_to?"            #=> true
			// Class.respond_to? :numerator        #=> false
			// ```
			//
			// @param [String]
			// @return [Boolean]
			Name: "respond_to?",
			Fn: func(receiver Object, sourceLine int) builtinMethodBody {
				return func(t *Thread, args []Object, blockFrame *normalCallFrame) Object {
					if len(args) != 1 {
						return t.vm.InitErrorObject(errors.ArgumentError, sourceLine, "Expect 1 argument. got=%d", len(args))
					}

					arg, ok := args[0].(*StringObject)
					if !ok {
						return t.vm.InitErrorObject(errors.TypeError, sourceLine, errors.WrongArgumentTypeFormat, classes.StringClass, arg.Class().Name)
					}

					r := receiver
					if r.findMethod(arg.value) == nil {
						return FALSE
					}
					return TRUE
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
			Fn: func(receiver Object, sourceLine int) builtinMethodBody {
				return func(t *Thread, args []Object, blockFrame *normalCallFrame) Object {
					if len(args) != 0 {
						return t.vm.InitErrorObject(errors.ArgumentError, sourceLine, "Expect 0 argument. got: %d", len(args))
					}

					c, ok := receiver.(*RClass)

					if !ok {
						return t.vm.InitErrorObject(errors.UndefinedMethodError, sourceLine, "Undefined Method '%s' for %s", "#superclass", receiver.toString())
					}

					superClass := c.returnSuperClass()

					if superClass == nil {
						return NULL
					}

					return superClass
				}
			},
		},
	}
}

// Instance methods -----------------------------------------------------
func builtinClassCommonInstanceMethods() []*BuiltinMethodObject {
	return []*BuiltinMethodObject{
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
			Fn: func(receiver Object, sourceLine int) builtinMethodBody {
				return func(t *Thread, args []Object, blockFrame *normalCallFrame) Object {
					className := receiver.Class().Name
					compareClassName := args[0].Class().Name

					if className == compareClassName && reflect.DeepEqual(receiver, args[0]) {
						return TRUE
					}
					return FALSE
				}
			},
		},
		{
			// Inverts the boolean value. Any objects other than `nil` and `false` are `true`, thus returns `false`.
			//
			// ```ruby
			// !true  # => false
			// !false # => true
			// !nil   # => true
			// !555   # => false
			// ```
			//
			// @param object [Object] object that return boolean value to invert
			// @return [Object] Inverted boolean value
			Name: "!",
			Fn: func(receiver Object, sourceLine int) builtinMethodBody {
				return func(t *Thread, args []Object, blockFrame *normalCallFrame) Object {

					rightValue, ok := receiver.(*BooleanObject)
					if !ok {
						return FALSE
					}

					if rightValue.value {
						return FALSE
					}
					return TRUE
				}
			},
		},
		{
			// General method for comparing inequality of the objects
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
			// @return [Boolean]
			Name: "!=",
			Fn: func(receiver Object, sourceLine int) builtinMethodBody {
				return func(t *Thread, args []Object, blockFrame *normalCallFrame) Object {
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
			// Returns true if a block is given in the current context and `yield` is ready to call.
			//
			// **Note:** The method name does not end with '?' because the sign is unavailable in Goby for now.
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
			Fn: func(receiver Object, sourceLine int) builtinMethodBody {
				return func(t *Thread, args []Object, blockFrame *normalCallFrame) Object {
					cf := t.callFrameStack.callFrames[t.callFrameStack.pointer-2]

					if cf.BlockFrame() == nil {
						return FALSE
					}

					return TRUE
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
			Fn: func(receiver Object, sourceLine int) builtinMethodBody {
				return func(t *Thread, args []Object, blockFrame *normalCallFrame) Object {

					switch r := receiver.(type) {
					case Object:
						return r.Class()
					default:
						return t.vm.InitErrorObject(errors.InternalError, sourceLine, errors.CantCallClass, string(r.Class().ReturnName()))
					}
				}
			},
		},
		// Exits from the interpreter, returning the specified exit code (if any).
		//
		// The method itself formally returns nil, although it's not usable.
		//
		// ```ruby
		// exit                    # exits with status code 0
		// exit(1)                 # exits with status code 1
		// ```
		//
		// @param [Integer] exit code (optional), defaults to 0
		// @return nil
		{
			Name: "exit",
			Fn: func(receiver Object, sourceLine int) builtinMethodBody {
				return func(t *Thread, args []Object, blockFrame *normalCallFrame) Object {
					switch len(args) {
					case 0:
						os.Exit(0)
					case 1:
						exitCode, ok := args[0].(*IntegerObject)

						if !ok {
							return t.vm.InitErrorObject(errors.TypeError, sourceLine, errors.WrongArgumentTypeFormat, classes.IntegerClass, args[0].Class().Name)
						}

						os.Exit(exitCode.value)
					default:
						return t.vm.InitErrorObject(errors.ArgumentError, sourceLine, "Expected at most 1 argument; got: %d", len(args))
					}

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
			Fn: func(receiver Object, sourceLine int) builtinMethodBody {
				return func(t *Thread, args []Object, blockFrame *normalCallFrame) Object {
					if len(args) != 1 {
						return t.vm.InitErrorObject(errors.ArgumentError, sourceLine, "Expect 1 argument. got: %d", len(args))
					}

					c := args[0]
					gobyClass, ok := c.(*RClass)

					if !ok {
						return t.vm.InitErrorObject(errors.TypeError, sourceLine, errors.WrongArgumentTypeFormat, classes.ClassClass, c.Class().Name)
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
			Name: "inherits_method_missing?",
			Fn: func(receiver Object, sourceLine int) builtinMethodBody {
				return func(t *Thread, args []Object, blockFrame *normalCallFrame) Object {
					if len(args) != 0 {
						return t.vm.InitErrorObject(errors.ArgumentError, sourceLine, "Expect 0 argument. got: %d", len(args))
					}

					if receiver.Class().inheritsMethodMissing {
						return TRUE
					}

					return FALSE
				}
			},
		},
		{
			Name: "instance_eval",
			Fn: func(receiver Object, sourceLine int) builtinMethodBody {
				return func(t *Thread, args []Object, blockFrame *normalCallFrame) Object {
					switch len(args) {
					case 0:
					case 1:
						if args[0].Class().Name == classes.BlockClass {
							blockObj := args[0].(*BlockObject)
							blockFrame = newNormalCallFrame(blockObj.instructionSet, blockObj.instructionSet.filename, sourceLine)
							blockFrame.ep = blockObj.ep
							blockFrame.self = receiver
							blockFrame.isBlock = true
						} else {
							return t.vm.InitErrorObject(errors.TypeError, sourceLine, errors.WrongArgumentTypeFormat, classes.BlockClass, args[0].Class().Name)
						}
					default:
						return t.vm.InitErrorObject(errors.ArgumentError, sourceLine, "Expect at most 1 arguments. got: %d", len(args))
					}

					if blockFrame == nil {
						return receiver
					}

					if blockIsEmpty(blockFrame) {
						return receiver
					}

					blockFrame.self = receiver
					result := t.builtinMethodYield(blockFrame)

					return result.Target
				}
			},
		},
		// Returns the value of the instance variable.
		// Only string literal with `@` is supported.
		//
		// ```ruby
		// class Foo
		//   def initialize
		//     @bar = 99
		//   end
		// end
		//
		// a = Foo.new
		// a.instance_variable_get("@bar")   #=> 99
		// ```
		//
		// @param string [String]
		// @return [Object], value
		{
			Name: "instance_variable_get",
			Fn: func(receiver Object, sourceLine int) builtinMethodBody {
				return func(t *Thread, args []Object, blockFrame *normalCallFrame) Object {
					if len(args) != 1 {
						return t.vm.InitErrorObject(errors.ArgumentError, sourceLine, "Expect 1 arguments. got: %d", len(args))
					}
					arg, isStr := args[0].(*StringObject)

					if !isStr {
						return t.vm.InitErrorObject(errors.TypeError, sourceLine, errors.WrongArgumentTypeFormat, classes.StringClass, args[0].Class().Name)
					}

					obj, ok := receiver.InstanceVariableGet(arg.value)

					if !ok {
						return NULL
					}

					return obj
				}
			},
		},
		{
			// Updates the specified instance variable with the value provided
			// Only string literal with `@` is supported for specifying an instance variable.
			//
			// ```ruby
			// class Foo
			//   def initialize
			//     @bar = 99
			//   end
			// end
			//
			// a = Foo.new
			// a.instance_variable_set("@bar", 42)
			// ```
			//
			// @param string [String], value [Object]
			// @return [Object] value
			Name: "instance_variable_set",
			Fn: func(receiver Object, sourceLine int) builtinMethodBody {
				return func(t *Thread, args []Object, blockFrame *normalCallFrame) Object {
					if len(args) != 2 {
						return t.vm.InitErrorObject(errors.ArgumentError, sourceLine, "Expect 2 arguments. got: %d", len(args))
					}

					argName, isStr := args[0].(*StringObject)
					obj := args[1]

					if !isStr {
						return t.vm.InitErrorObject(errors.TypeError, sourceLine, errors.WrongArgumentTypeFormat, classes.StringClass, args[0].Class().Name)
					}

					receiver.InstanceVariableSet(argName.value, obj)

					return obj
				}
			},
		},
		// Returns an array that contains the method names of the receiver.
		//
		// ```ruby
		// Class.methods
		// ["ancestors", "attr_accessor", "attr_reader", "attr_writer", "extend", "include", "name", "new", "superclass", "!", "!=", "==", "block_given?", "class", "instance_variable_get", "instance_variable_set", "is_a?", "methods", "nil?", "puts", "require", "require_relative", "send", "singleton_class", "sleep", "thread", "to_s"]
		// ```
		//
		// @param class [Class] Receiver
		// @return [Array]
		{
			Name: "methods",
			Fn: func(receiver Object, sourceLine int) builtinMethodBody {
				return func(t *Thread, args []Object, blockFrame *normalCallFrame) Object {
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
					return t.vm.InitArrayObject(methods)
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
			Fn: func(receiver Object, sourceLine int) builtinMethodBody {
				return func(t *Thread, args []Object, blockFrame *normalCallFrame) Object {
					if len(args) != 0 {
						return t.vm.InitErrorObject(errors.ArgumentError, sourceLine, "Expect 0 argument. got: %d", len(args))
					}
					return FALSE
				}
			},
		},
		{
			// Returns object's unique id from Go's `receiver.id()`
			// @param n/a []
			// @return [Integer] Object's address
			Name: "object_id",
			Fn: func(receiver Object, sourceLine int) builtinMethodBody {
				return func(t *Thread, args []Object, blockFrame *normalCallFrame) Object {
					return t.vm.InitIntegerObject(receiver.id())
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
			Fn: func(receiver Object, sourceLine int) builtinMethodBody {
				return func(t *Thread, args []Object, blockFrame *normalCallFrame) Object {
					for _, arg := range args {
						fmt.Println(arg.toString())
					}

					return NULL
				}
			},
		},
		{
			Name: "raise",
			Fn: func(receiver Object, sourceLine int) builtinMethodBody {
				return func(t *Thread, args []Object, blockFrame *normalCallFrame) Object {
					switch len(args) {
					case 0:
						return t.vm.InitErrorObject(errors.InternalError, sourceLine, "")
					case 1:
						return t.vm.InitErrorObject(errors.InternalError, sourceLine, "'%s'", args[0].toString())
					case 2:
						errorClass, ok := args[0].(*RClass)

						if !ok {
							return t.vm.InitErrorObject(errors.ArgumentError, sourceLine, "Expect error class, got: %s", args[0].Class().Name)
						}

						return t.vm.InitErrorObject(errorClass.Name, sourceLine, "'%s'", args[1].toString())
					}

					return t.vm.InitErrorObject(errors.ArgumentError, sourceLine, "Expect at most 2 arguments. got: %d", len(args))
				}
			},
		},
		{
			// A predicate class method that returns `true` if the object has an ability to respond to the method, otherwise `false`.
			// Note that signs like `+` or `?` should be String literal.
			//
			// ```ruby
			// 1.respond_to? :to_i               #=> true
			// "string".respond_to? "+"          #=> true
			// 1.respond_to? :numerator          #=> false
			// ```
			//
			// @param [String]
			// @return [Boolean]
			Name: "respond_to?",
			Fn: func(receiver Object, sourceLine int) builtinMethodBody {
				return func(t *Thread, args []Object, blockFrame *normalCallFrame) Object {
					if len(args) != 1 {
						return t.vm.InitErrorObject(errors.ArgumentError, sourceLine, "Expect 1 argument. got=%d", len(args))
					}

					arg, ok := args[0].(*StringObject)
					if !ok {
						return t.vm.InitErrorObject(errors.TypeError, sourceLine, errors.WrongArgumentTypeFormat, classes.StringClass, arg.Class().Name)
					}

					r := receiver
					if r.findMethod(arg.value) == nil {
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
			// ```ruby
			// require("db")
			// File.extname("foo.rb")
			// ```
			//
			// TBD: the load paths for `require`
			//
			// @param filename [String] Quoted file name of the library, without extension
			// @return [Boolean] Result of loading module
			Name: "require",
			Fn: func(receiver Object, sourceLine int) builtinMethodBody {
				return func(t *Thread, args []Object, blockFrame *normalCallFrame) Object {
					if len(args) != 1 {
						return t.vm.InitErrorObject(errors.ArgumentError, sourceLine, "Expect 1 argument. got: %d", len(args))
					}
					switch args[0].(type) {
					case *StringObject:
						libName := args[0].(*StringObject).value
						initFunc, ok := standardLibraries[libName]

						if !ok {
							externalClassLock.Lock()
							loaders, ok := externalClasses[libName]
							externalClassLock.Unlock()
							if !ok {
								err := t.execGobyLib(libName + ".gb")
								if err != nil {
									return t.vm.InitErrorObject(errors.InternalError, sourceLine, "Can't require \"%s\"", libName)
								}
							}
							initFunc = func(v *VM) {
								for _, l := range loaders {
									l(v)
								}
							}
						}

						initFunc(t.vm)

						return TRUE
					default:
						return t.vm.InitErrorObject(errors.InternalError, sourceLine, "Can't require \"%s\". Pass a string instead", args[0].(Object).Class().Name)
					}
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
			Fn: func(receiver Object, sourceLine int) builtinMethodBody {
				return func(t *Thread, args []Object, blockFrame *normalCallFrame) Object {
					if len(args) != 1 {
						return t.vm.InitErrorObject(errors.ArgumentError, sourceLine, "Expect 1 argument. got: %d", len(args))
					}
					switch args[0].(type) {
					case *StringObject:
						callerDir := path.Dir(t.vm.currentFilePath())
						filepath := args[0].(*StringObject).value
						filepath = path.Join(callerDir, filepath)
						filepath = filepath + ".gb"

						t.execFile(filepath)

						return TRUE
					default:
						return t.vm.InitErrorObject(errors.InternalError, sourceLine, "Can't require \"%s\". Pass a string instead", args[0].(Object).Class().Name)
					}
				}
			},
		},
		// Invoke the specified instance method or class method.
		// - Method name should be either a symbol or String (required).
		// - You can pass one or more arguments (option).
		// - A block can also be provided (option).
		//
		//
		// ```ruby
		// class Foo
		//   def self.bar
		//     10
		//   end
		// end
		//
		// Foo.send(:bar)  #=> 10
		//
		// class Math
		//   def self.add(x, y)
		//     x + y
		//   end
		// end
		//
		// Math.send(:add, 10, 15) #=> 25
		//
		// class Foo
		//   def bar(x, y)
		//     yield x, y
		//   end
		// end
		// a = Foo.new
		//
		// a.send(:bar, 7, 8) do |i, j| i * j; end   #=> 56
		// ```
		//
		// @param name [String/symbol], args [Object], block
		// @return [Object]
		{
			Name: "send",
			Fn: func(receiver Object, sourceLine int) builtinMethodBody {
				return func(t *Thread, args []Object, blockFrame *normalCallFrame) Object {
					if len(args) < 1 {
						return t.vm.InitErrorObject(errors.ArgumentError, sourceLine, "no method name given")
					}

					name, ok := args[0].(*StringObject)

					if !ok {
						return t.vm.InitErrorObject(errors.TypeError, sourceLine, errors.WrongArgumentTypeFormat, classes.StringClass, args[0].Class().Name)
					}

					t.sendMethod(name.value, len(args)-1, blockFrame, sourceLine)

					return t.Stack.top().Target
				}
			},
		},
		{
			// Returns the singleton class object of the receiver class.
			//
			// ```ruby
			// class Foo
			// end
			//
			// Foo.singleton_class
			// #=> #<Class:Foo>
			// Foo.singleton_class.ancestors
			// #=> [#<Class:Foo>, #<Class:Object>, Class, Object]
			// ```
			//
			// @param class [Class] receiver
			// @return [Object] singleton class
			Name: "singleton_class",
			Fn: func(receiver Object, sourceLine int) builtinMethodBody {
				return func(t *Thread, args []Object, blockFrame *normalCallFrame) Object {
					r := receiver
					if r.SingletonClass() == nil {
						id := t.vm.InitIntegerObject(r.id())
						singletonClass := t.vm.createRClass(fmt.Sprintf("#<Class:#<%s:%s>>", r.Class().Name, id.toString()))
						singletonClass.isSingleton = true
						return singletonClass
					}
					return receiver.SingletonClass()
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
			Fn: func(receiver Object, sourceLine int) builtinMethodBody {
				return func(t *Thread, args []Object, blockFrame *normalCallFrame) Object {
					if len(args) != 1 {
						return t.vm.InitErrorObject(errors.ArgumentError, sourceLine, "Expect 1 argument. got: %d", len(args))
					}

					int, ok := args[0].(*IntegerObject)

					if ok {
						seconds := int.value
						time.Sleep(time.Duration(seconds) * time.Second)
						return int
					}

					float, ok := args[0].(*FloatObject)

					if ok {
						nanoseconds := int64(float.value * float64(time.Second/time.Nanosecond))
						time.Sleep(time.Duration(nanoseconds) * time.Nanosecond)
						return float
					}

					return t.vm.InitErrorObject(errors.TypeError, sourceLine, errors.WrongArgumentTypeFormat, "Numeric", args[0].Class().Name)
				}
			},
		},
		{
			Name: "thread",
			Fn: func(receiver Object, sourceLine int) builtinMethodBody {
				return func(t *Thread, args []Object, blockFrame *normalCallFrame) Object {
					if blockFrame == nil {
						return t.vm.InitErrorObject(errors.BlockError, sourceLine, errors.CantGetBlockWithoutBlockArgument)
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
			// Returns object's string representation.
			// @param n/a []
			// @return [String] Object's string representation.
			Name: "to_s",
			Fn: func(receiver Object, sourceLine int) builtinMethodBody {
				return func(t *Thread, args []Object, blockFrame *normalCallFrame) Object {
					return t.vm.initStringObject(receiver.toString())
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
func (c *RClass) toJSON(t *Thread) string {
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

func (c *RClass) lookupConstantInCurrentScope(constName string) *Pointer {
	constant, ok := c.constants[constName]

	if !ok {
		return nil
	}

	return constant
}

func (c *RClass) lookupConstantUnderCurrentScope(constName string) *Pointer {
	constant, ok := c.constants[constName]

	if !ok {
		if c.scope != nil {
			return c.scope.lookupConstantUnderCurrentScope(constName)
		}

		return nil
	}

	return constant
}

func (c *RClass) lookupConstantUnderAllScope(constName string) *Pointer {
	constant, ok := c.constants[constName]

	if !ok {
		if c.scope != nil {
			return c.scope.lookupConstantUnderCurrentScope(constName)
		}

		// Finding constant in superclass means it's out of the scope
		if c.superClass != nil && c.Name != classes.ObjectClass {
			constant, _ = c.constants[constName]
			return constant
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
		Fn: func(receiver Object, sourceLine int) builtinMethodBody {
			return func(t *Thread, args []Object, blockFrame *normalCallFrame) Object {
				v := receiver.InstanceVariableSet("@"+attrName, args[0])
				return v
			}
		},
	}
}

func generateAttrReadMethod(attrName string) *BuiltinMethodObject {
	return &BuiltinMethodObject{
		Name: attrName,
		Fn: func(receiver Object, sourceLine int) builtinMethodBody {
			return func(t *Thread, args []Object, blockFrame *normalCallFrame) Object {
				v, ok := receiver.InstanceVariableGet("@" + attrName)

				if ok {
					return v
				}

				return NULL
			}
		},
	}
}
