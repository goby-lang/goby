package vm

import (
	"fmt"

	"github.com/goby-lang/goby/vm/errors"
	"strings"
)

// Error class is actually a special struct to hold internal error types with messages.
// Goby developers need not to take care of the struct.
// Goby maintainers should consider using the appropriate error type.
// Cannot create instances of Error class, or inherit Error class.
//
// The type of internal errors:
//
// * `InternalError`: default error type
// * `ArgumentError`: an argument-related error
// * `NameError`: a constant-related error
// * `TypeError`: a type-related error
// * `UndefinedMethodError`: undefined-method error
// * `UnsupportedMethodError`: intentionally unsupported-method error
//
type Error struct {
	*baseObj
	message      string
	stackTraces  []string
	storedTraces bool
	Type         string
}

// Internal functions ===================================================

// Functions for initialization -----------------------------------------

func (vm *VM) initUnsupportedMethodError(sourceLine int, methodName string, receiver Object) *Error {
	return vm.initErrorObject(errors.UnsupportedMethodError, sourceLine, "Unsupported Method %s for %+v", methodName, receiver.toString())
}

func (vm *VM) initErrorObject(errorType string, sourceLine int, format string, args ...interface{}) *Error {
	errClass := vm.objectClass.getClassConstant(errorType)

	t := vm.mainThread
	cf := t.callFrameStack.top()

	switch cf := cf.(type) {
	case *normalCallFrame:
		// If program counter is 0 means we need to trace back to previous call frame
		if cf.pc == 0 {
			t.callFrameStack.pop()
			cf = t.callFrameStack.top().(*normalCallFrame)
		}
	case *goMethodCallFrame:
		t.callFrameStack.pop()
	}

	return &Error{
		baseObj: &baseObj{class: errClass},
		// Add 1 to source line because it's zero indexed
		message:     fmt.Sprintf(errorType+": "+format, args...),
		stackTraces: []string{fmt.Sprintf("from %s:%d", cf.FileName(), sourceLine)},
		Type:        errorType,
	}
}

func (vm *VM) initErrorClasses() {
	errTypes := []string{errors.InternalError, errors.ArgumentError, errors.NameError, errors.TypeError, errors.UndefinedMethodError, errors.UnsupportedMethodError, errors.ConstantAlreadyInitializedError, errors.HTTPError}

	for _, errType := range errTypes {
		c := vm.initializeClass(errType, false)
		vm.objectClass.setClassConstant(c)
	}
}

// Polymorphic helper functions -----------------------------------------

// toString returns the object's name as the string format
func (e *Error) toString() string {
	return e.Message()
}

// toJSON just delegates to `toString`
func (e *Error) toJSON() string {
	return e.toString()
}

func (e *Error) Value() interface{} {
	return e.message
}

// Message prints the error's message and its stack traces
func (e *Error) Message() string {
	return e.message + "\n" + strings.Join(e.stackTraces, "\n")
}
