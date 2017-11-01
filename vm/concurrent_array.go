package vm

import (
	"sync"

	"github.com/goby-lang/goby/vm/classes"
	"github.com/goby-lang/goby/vm/errors"
)

// Pseudo-constant definition of the forwarded methods, mapped to a boolean representing the
// requirement for a write lock (true) or read lock (false)
//
// We don't implement dig, as it has no concurrency guarantees.
var ConcurrentArrayMethodsForwardingTable = map[string]bool{
	"[]":           false,
	"*":            false,
	"+":            false,
	"[]=":          true,
	"any?":         false,
	"at":           false,
	"clear":        true,
	"concat":       true,
	"count":        false,
	"delete_at":    true,
	"each":         false,
	"each_index":   false,
	"empty?":       false,
	"first":        false,
	"flatten":      false,
	"join":         false,
	"last":         false,
	"length":       false,
	"map":          false,
	"pop":          true,
	"push":         true,
	"reduce":       false,
	"reverse":      false,
	"reverse_each": false,
	"rotate":       false,
	"select":       false,
	"shift":        true,
	"unshift":      true,
	"values_at":    false,
}

// ConcurrentArrayObject is a thread-safe Array, implemented as a wrapper of an ArrayObject, coupled
// with an R/W mutex.
//
// Arrays returned by any of the methods are in turn thread-safe.
//
// For implementation simplicity, methods are simple redirection, and defined via a table.
//
type ConcurrentArrayObject struct {
	*baseObj
	InternalArray *ArrayObject

	sync.RWMutex
}

// Class methods --------------------------------------------------------
func builtinConcurrentArrayClassMethods() []*BuiltinMethodObject {
	return []*BuiltinMethodObject{
		{
			Name: "new",
			Fn: func(receiver Object, sourceLine int) builtinMethodBody {
				return func(t *thread, args []Object, blockFrame *normalCallFrame) Object {
					if len(args) > 1 {
						return t.vm.initErrorObject(errors.ArgumentError, sourceLine, "Expect 0 or 1 arguments, got %d", len(args))
					}

					if len(args) == 0 {
						return t.vm.initConcurrentArrayObject([]Object{})
					} else {
						arg := args[0]
						arrayArg, ok := arg.(*ArrayObject)

						if !ok {
							return t.vm.initErrorObject(errors.TypeError, sourceLine, errors.WrongArgumentTypeFormat, classes.ArrayClass, arg.Class().Name)
						}

						return t.vm.initConcurrentArrayObject(arrayArg.Elements)
					}
				}
			},
		},
	}
}

// Instance methods -----------------------------------------------------
func builtinConcurrentArrayInstanceMethods() []*BuiltinMethodObject {
	methodDefinitions := []*BuiltinMethodObject{}

	for methodName, requireWriteLock := range ConcurrentArrayMethodsForwardingTable {
		methodFunction := DefineForwardedConcurrentArrayMethod(methodName, requireWriteLock)
		methodDefinitions = append(methodDefinitions, methodFunction)
	}

	return methodDefinitions
}

// Internal functions ===================================================

// Functions for initialization -----------------------------------------

func (vm *VM) initConcurrentArrayObject(elements []Object) *ConcurrentArrayObject {
	concurrent := vm.loadConstant("Concurrent", true)
	array := concurrent.getClassConstant("Array")

	return &ConcurrentArrayObject{
		baseObj:     &baseObj{class: array},
		InternalArray: vm.initArrayObject(elements[:]),
	}
}

func initConcurrentArrayClass(vm *VM) {
	concurrent := vm.loadConstant("Concurrent", true)
	array := vm.initializeClass("Array", false)

	array.setBuiltinMethods(builtinConcurrentArrayInstanceMethods(), false)
	array.setBuiltinMethods(builtinConcurrentArrayClassMethods(), true)

	concurrent.setClassConstant(array)
}


// Object interface functions -------------------------------------------

// toJSON returns the object's name as the JSON string format
func (cac *ConcurrentArrayObject) toJSON() string {
	return cac.InternalArray.toJSON()
}

// toString returns the object's name as the string format
func (cac *ConcurrentArrayObject) toString() string {
	return cac.InternalArray.toString()
}

// Value returns the object
func (cac *ConcurrentArrayObject) Value() interface{} {
	return cac.InternalArray.Elements
}

// Helper functions -----------------------------------------------------

func DefineForwardedConcurrentArrayMethod(methodName string, requireWriteLock bool) *BuiltinMethodObject {
	return &BuiltinMethodObject{
		Name: methodName,
		Fn: func(receiver Object, sourceLine int) builtinMethodBody {
			return func(t *thread, args []Object, blockFrame *normalCallFrame) Object {
				concurrentArray := receiver.(*ConcurrentArrayObject)

				if requireWriteLock {
					concurrentArray.Lock()
				} else {
					concurrentArray.RLock()
				}

        arrayMethodObject := concurrentArray.InternalArray.findMethod(methodName).(*BuiltinMethodObject)
        result := arrayMethodObject.Fn(concurrentArray.InternalArray, sourceLine)(t, args, blockFrame)

				if requireWriteLock {
					concurrentArray.Unlock()
				} else {
					concurrentArray.RUnlock()
				}

				switch result.(type) {
				case *ArrayObject:
          return t.vm.initConcurrentArrayObject(result.(*ArrayObject).Elements)
				default:
					return result
				}
			}
		},
	}
}
