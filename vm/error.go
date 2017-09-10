package vm

import (
	"fmt"

	"github.com/goby-lang/goby/vm/errors"
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
	Message string
	Type    string
}

// Internal functions ===================================================

// Functions for initialization -----------------------------------------

func (vm *VM) initErrorObject(errorType, format string, args ...interface{}) *Error {
	errClass := vm.objectClass.getClassConstant(errorType)

	t := vm.mainThread
	cf := t.callFrameStack.top()

	// If program counter is 0 means we need to trace back to previous call frame
	if cf.pc == 0 {
		t.callFrameStack.pop()
		cf = t.callFrameStack.top()
	}

	i := cf.instructionSet.instructions[cf.pc-1]

	return &Error{
		baseObj: &baseObj{class: errClass},
		// Add 1 to source line because it's zero indexed
		Message: fmt.Sprintf("%s. At %s:%d", fmt.Sprintf(errorType+": "+format, args...), cf.instructionSet.filename, i.sourceLine+1),
		Type:    errorType,
	}
}

func (vm *VM) initErrorClasses() {
	errTypes := []string{errors.InternalError, errors.ArgumentError, errors.NameError, errors.TypeError, errors.UndefinedMethodError, errors.UnsupportedMethodError, errors.ConstantAlreadyInitializedError}

	for _, errType := range errTypes {
		c := vm.initializeClass(errType, false)
		vm.objectClass.setClassConstant(c)
	}
}

// Polymorphic helper functions -----------------------------------------

// Returns the object's name as the string format
func (e *Error) toString() string {
	return "ERROR: " + e.Message
}

// Alias of toString
func (e *Error) toJSON() string {
	return e.toString()
}
