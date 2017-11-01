package vm

import (
	"sync"

	"github.com/goby-lang/goby/vm/errors"
)

// ConcurrentRWLockObject is a Readers-Writer Lock (readers can concurrently put a lock, while a
// writer requires exclusive access).
//
// The implementation internally uses Go's `sync.RWLock` type.
//
// ```ruby
// require 'concurrent/rw_lock'
// lock = Concurrent::RWLock.new
// lock.with_read_lock do
//   # critical section
// end
// lock.with_write_lock do
//   # critical section
// end
// ```
//
type ConcurrentRWLockObject struct {
	*baseObj
	mutex    sync.RWMutex
}

// Class methods --------------------------------------------------------
func builtinConcurrentRWLockClassMethods() []*BuiltinMethodObject {
	return []*BuiltinMethodObject{
		{
			Name: "new",
			Fn: func(receiver Object, sourceLine int) builtinMethodBody {
				return func(t *thread, args []Object, blockFrame *normalCallFrame) Object {
					if len(args) != 0 {
						return t.vm.initErrorObject(errors.ArgumentError, sourceLine, "Expected 0 arguments, got %d", len(args))
					}

					return t.vm.initConcurrentRWLockObject()
				}
			},
		},
	}
}

// Instance methods -----------------------------------------------------
func builtinConcurrentRWLockInstanceMethods() []*BuiltinMethodObject {
	return []*BuiltinMethodObject{
		{
			// Executes the block with a read lock.
			// The lock is freed upon exiting the block.
			//
			// ```Ruby
			// lock = Concurrent::RWLock.new
			// lock.with_read_lock do
			//   # critical section
			// end
			//
			// @return [Object] the yielded value of the block.
			// ```
			Name: "with_read_lock",
			Fn: func(receiver Object, sourceLine int) builtinMethodBody {
				return func(t *thread, args []Object, blockFrame *normalCallFrame) Object {
					if blockFrame == nil {
						return t.vm.initErrorObject(errors.InternalError, sourceLine, errors.CantYieldWithoutBlockFormat)
					}

					if len(args) != 0 {
						t.callFrameStack.pop()

						return t.vm.initErrorObject(errors.ArgumentError, sourceLine, "Expected 0 arguments, got %d", len(args))
					}

					lockObject := receiver.(*ConcurrentRWLockObject)

					lockObject.mutex.RLock();

					blockReturnValue := t.builtinMethodYield(blockFrame).Target

					lockObject.mutex.RUnlock();

					return blockReturnValue
				}
			},
		},
		{
			// Executes the block with a write lock.
			// The lock is freed upon exiting the block.
			//
			// ```Ruby
			// lock = Concurrent::RWLock.new
			// lock.with_write_lock do
			//   # critical section
			// end
			//
			// @return [Object] the yielded value of the block.
			// ```
			Name: "with_write_lock",
			Fn: func(receiver Object, sourceLine int) builtinMethodBody {
				return func(t *thread, args []Object, blockFrame *normalCallFrame) Object {
					if blockFrame == nil {
						return t.vm.initErrorObject(errors.InternalError, sourceLine, errors.CantYieldWithoutBlockFormat)
					}

					if len(args) != 0 {
						t.callFrameStack.pop()

						return t.vm.initErrorObject(errors.ArgumentError, sourceLine, "Expected 0 arguments, got %d", len(args))
					}

					lockObject := receiver.(*ConcurrentRWLockObject)

					lockObject.mutex.Lock();

					blockReturnValue := t.builtinMethodYield(blockFrame).Target

					lockObject.mutex.Unlock();

					return blockReturnValue
				}
			},
		},
	}
}

// Internal functions ===================================================

// Functions for initialization -----------------------------------------

func (vm *VM) initConcurrentRWLockObject() *ConcurrentRWLockObject {
	concurrentModule := vm.loadConstant("Concurrent", true)
	lockClass := concurrentModule.getClassConstant("RWLock")

	return &ConcurrentRWLockObject{
		baseObj: &baseObj{class: lockClass},
		mutex:   sync.RWMutex{},
	}
}

func initConcurrentRWLockClass(vm *VM) {
	concurrentModule := vm.loadConstant("Concurrent", true)
	lockClass := vm.initializeClass("RWLock", false)

	lockClass.setBuiltinMethods(builtinConcurrentRWLockInstanceMethods(), false)
	lockClass.setBuiltinMethods(builtinConcurrentRWLockClassMethods(), true)

	concurrentModule.setClassConstant(lockClass)
}

// Polymorphic helper functions -----------------------------------------

// Value returns the object
func (lock *ConcurrentRWLockObject) Value() interface{} {
	return lock.mutex
}

// toString returns the object's name as the string format
func (lock *ConcurrentRWLockObject) toString() string {
	return "<Instance of: " + lock.class.Name + ">"
}

// toJSON just delegates to toString
func (lock *ConcurrentRWLockObject) toJSON() string {
	return lock.toString()
}
