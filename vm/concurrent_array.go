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

var concurrentArrayMethodsForwardingTable = map[string]bool{
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

// concurrentArrayObject is a thread-safe Array, implemented as a wrapper of an ArrayObject, coupled
// with an R/W mutex.
//
// Arrays returned by any of the methods are in turn thread-safe.
//
// For implementation simplicity, methods are simple redirection, and defined via a table.
//
type concurrentArrayObject struct {
	*BaseObj
	InternalArray *ArrayObject

	sync.RWMutex
}

// Class methods --------------------------------------------------------
var builtinConcurrentArrayClassMethods = []*BuiltinMethodObject{
	{
		Name: "new",
		Fn: func(receiver Object, sourceLine int, t *Thread, args []Object, blockFrame *normalCallFrame) Object {
			aLen := len(args)

			switch aLen {
			case 0:
				return t.vm.initConcurrentArrayObject([]Object{})
			case 1:
				arg := args[0]
				arrayArg, ok := arg.(*ArrayObject)

				if !ok {
					return t.vm.InitErrorObject(errors.TypeError, sourceLine, errors.WrongArgumentTypeFormat, classes.ArrayClass, arg.Class().Name)
				}

				return t.vm.initConcurrentArrayObject(arrayArg.Elements)
			default:
				return t.vm.InitErrorObject(errors.ArgumentError, sourceLine, errors.WrongNumberOfArgumentLess, 1, aLen)
			}

		},
	},
}

// Internal functions ===================================================

// Functions for initialization -----------------------------------------

func (vm *VM) initConcurrentArrayObject(elements []Object) *concurrentArrayObject {
	concurrent := vm.loadConstant("Concurrent", true)

	return &concurrentArrayObject{
		BaseObj:       NewBaseObject(concurrent.getClassConstant(classes.ArrayClass)),
		InternalArray: vm.InitArrayObject(elements[:]),
	}
}

func initConcurrentArrayClass(vm *VM) {
	concurrent := vm.loadConstant("Concurrent", true)
	array := vm.initializeClass(classes.ArrayClass)

	var arrayMethodDefinitions = []*BuiltinMethodObject{}

	for methodName, requireWriteLock := range concurrentArrayMethodsForwardingTable {
		methodFunction := defineForwardedConcurrentArrayMethod(methodName, requireWriteLock)
		arrayMethodDefinitions = append(arrayMethodDefinitions, methodFunction)
	}

	array.setBuiltinMethods(arrayMethodDefinitions, false)
	array.setBuiltinMethods(builtinConcurrentArrayClassMethods, true)

	concurrent.setClassConstant(array)
}

// Object interface functions -------------------------------------------

// ToJSON returns the object's name as the JSON string format
func (ca *concurrentArrayObject) ToJSON(t *Thread) string {
	return ca.InternalArray.ToJSON(t)
}

// ToString returns the object's name as the string format
func (ca *concurrentArrayObject) ToString() string {
	return ca.InternalArray.Inspect()
}

// Inspect delegates to ToString
func (ca *concurrentArrayObject) Inspect() string {
	return ca.ToString()
}

// Value returns the object
func (ca *concurrentArrayObject) Value() interface{} {
	return ca.InternalArray.Elements
}

func (ca *concurrentArrayObject) equalTo(compared Object) bool {
	c, ok := compared.(*concurrentArrayObject)

	if !ok {
		return false
	}

	return ca.InternalArray.equalTo(c.InternalArray)
}

// Helper functions -----------------------------------------------------

func defineForwardedConcurrentArrayMethod(methodName string, requireWriteLock bool) *BuiltinMethodObject {
	return &BuiltinMethodObject{
		Name: methodName,
		Fn: func(receiver Object, sourceLine int, t *Thread, args []Object, blockFrame *normalCallFrame) Object {
			concurrentArray := receiver.(*concurrentArrayObject)

			if requireWriteLock {
				concurrentArray.Lock()
			} else {
				concurrentArray.RLock()
			}

			arrayMethodObject := concurrentArray.InternalArray.findMethod(methodName).(*BuiltinMethodObject)
			result := arrayMethodObject.Fn(concurrentArray.InternalArray, sourceLine, t, args, blockFrame)

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
		},
	}
}
