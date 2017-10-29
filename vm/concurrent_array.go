package vm

import (
	"sync"

	"github.com/goby-lang/goby/vm/classes"
	"github.com/goby-lang/goby/vm/errors"
)

// ConcurrentArrayObject is a thread-safe Array, implemented using an R/W mutex.
//
type ConcurrentArrayObject struct {
	ArrayObject

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

// Internal functions ===================================================

// Called from the superclass; executes a method using a R/W mutex.
func (cao *ConcurrentArrayObject) executeWithLock(arrayType int, method func() Object) Object {
	switch arrayType {
	case ReadArrayLock:
		cao.RLock()
	case WriteArrayLock:
		cao.Lock()
	}

	result := method()

	switch arrayType {
	case ReadArrayLock:
		cao.RUnlock()
	case WriteArrayLock:
		cao.Unlock()
	}

	return result
}

// Functions for initialization -----------------------------------------

func (vm *VM) initConcurrentArrayObject(elements []Object) *ConcurrentArrayObject {
	concurrent := vm.loadConstant("Concurrent", true)
	array := concurrent.getClassConstant("Array")

	return &ConcurrentArrayObject{
		ArrayObject{
			baseObj:  &baseObj{class: array},
			Elements: elements,
		},
		sync.RWMutex{},
	}
}

func initConcurrentArrayClass(vm *VM) {
	concurrent := vm.loadConstant("Concurrent", true)
	array := vm.initializeClass("Array", false)

	array.setBuiltinMethods(builtinArrayInstanceMethods(), false)
	array.setBuiltinMethods(builtinConcurrentArrayClassMethods(), true)

	concurrent.setClassConstant(array)
}
