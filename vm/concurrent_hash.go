package vm

import (
	"sync"

	"github.com/goby-lang/goby/vm/classes"
	"github.com/goby-lang/goby/vm/errors"
)

// Pseudo-constant definition of the forwarded methods, mapped to a boolean representing the
// requirement for a write lock (true) or read lock (false)
var ConcurrentHashMethodsForwardingTable = map[string]bool{
	"[]": false,
	"[]=": true,
	"any?": false,
	"clear": true,
	"default": false,
	"default=": true,
	"delete": true,
	"delete_if": true,
	"dig": false,
	"each": false,
	"each_key": false,
	"each_value": false,
	"empty?": false,
	"eql?": false,
	"fetch": false,
	"has_key?": false,
	"has_value?": false,
	"keys": false,
	"length": false,
	"map_values": false,
	"merge": false,
	"select": false,
	"sorted_keys": false,
	"to_a": false,
	"to_json": false,
	"to_s": false,
	"transform_values": false,
	"values": false,
	"values_at": false,
}

// ConcurrentHashObject is a thread-safe Hash.
//
// The current design is the simplest possible, via a R/W mutex; if/when performance will be a
// concern, it's trivial to write more sophisticated implementations, since the logic is entirely
// encapsulated.
//
// The class is current a subclass of HashObject internally, but an Object direct subclass at
// user level.
type ConcurrentHashObject struct {
	*baseObj
	InternalHash *HashObject

	sync.RWMutex
}

// Class methods --------------------------------------------------------
func builtinConcurrentHashClassMethods() []*BuiltinMethodObject {
	return []*BuiltinMethodObject{
		{
			Name: "new",
			Fn: func(receiver Object) builtinMethodBody {
				return func(t *thread, args []Object, blockFrame *callFrame) Object {
					if len(args) > 1 {
						return t.vm.initErrorObject(errors.ArgumentError, "Expect 0 or 1 arguments, got %d", len(args))
					}

					if len(args) == 0 {
						return t.vm.initConcurrentHashObject(t.vm.initHashObject(make(map[string]Object)))
					} else {
						arg := args[0]
						hashArg, ok := arg.(*HashObject)

						if !ok {
							return t.vm.initErrorObject(errors.TypeError, errors.WrongArgumentTypeFormat, classes.HashClass, arg.Class().Name)
						}

						return t.vm.initConcurrentHashObject(hashArg)
					}
				}
			},
		},
	}
}

// Instance methods -----------------------------------------------------
func builtinConcurrentHashInstanceMethods() []*BuiltinMethodObject {
	methodDefinitions := []*BuiltinMethodObject{}

	for methodName, requireWriteLock := range ConcurrentHashMethodsForwardingTable {
		methodFunction := DefineForwardedConcurrentHashMethod(methodName, requireWriteLock)
		methodDefinitions = append(methodDefinitions, methodFunction)
	}

	return methodDefinitions
}

// Internal functions ===================================================

// Functions for initialization -----------------------------------------

func (vm *VM) initConcurrentHashObject(internalHash *HashObject) *ConcurrentHashObject {
	return &ConcurrentHashObject{
		baseObj: &baseObj{class: vm.topLevelClass(classes.ConcurrentHashClass)},
		InternalHash: internalHash,
	}
}

func initConcurrentHashClass(vm *VM) {
	chc := vm.initializeClass(classes.ConcurrentHashClass, false)
	chc.setBuiltinMethods(builtinConcurrentHashInstanceMethods(), false)
	chc.setBuiltinMethods(builtinConcurrentHashClassMethods(), true)
	vm.objectClass.setClassConstant(chc)
}

// Object interface functions -------------------------------------------

// toJSON returns the object's name as the JSON string format
func (chc *ConcurrentHashObject) toJSON() string {
	return chc.InternalHash.toJSON()
}

// toString returns the object's name as the string format
func (chc *ConcurrentHashObject) toString() string {
	return chc.InternalHash.toString()
}

// Value returns the object
func (chc *ConcurrentHashObject) Value() interface{} {
	return chc.InternalHash.Pairs
}

// Helper functions -----------------------------------------------------

func DefineForwardedConcurrentHashMethod(methodName string, requireWriteLock bool) *BuiltinMethodObject {
	return &BuiltinMethodObject {
		Name: methodName,
		Fn: func(receiver Object) builtinMethodBody {
			return func(t *thread, args []Object, blockFrame *callFrame) Object {
				concurrentHash := receiver.(*ConcurrentHashObject)
				hashMethodObject := concurrentHash.InternalHash.findMethod(methodName).(*BuiltinMethodObject)

				if requireWriteLock {
					concurrentHash.Lock()
				} else {
					concurrentHash.RLock()
				}

				result := hashMethodObject.Fn(concurrentHash.InternalHash)(t, args, blockFrame)

				if requireWriteLock {
					concurrentHash.Unlock()
				} else {
					concurrentHash.RUnlock()
				}

				switch result.(type) {
				case *HashObject:
					return t.vm.initConcurrentHashObject(result.(*HashObject))
				default:
					return result
				}
			}
		},
	}
}
