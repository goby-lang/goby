package vm

import (
	"fmt"
)

const (
	// InternalError is the default error type
	InternalError = "InternalError"
	// ArgumentError is for an argument-related error
	ArgumentError = "ArgumentError"
	// NameError is for a constant-related error
	NameError = "NameError"
	// TypeError is for a type-related error
	TypeError = "TypeError"
	// UndefinedMethodError is for an undefined-method error
	UndefinedMethodError = "UndefinedMethodError"
	// UnsupportedMethodError is for an intentionally unsupported-method error
	UnsupportedMethodError = "UnsupportedMethodError"
	// ConstantAlreadyInitializedError means user re-declares twice
	ConstantAlreadyInitializedError = "ConstantAlreadyInitializedError"
	//HTTPError is for general errors returned from http functions
	HTTPError = "HTTP Error"
	//HTTPResponseError is for non 200 responses in general contexts
	//ex Net::HTTP.post()
	HTTPResponseError = "HTTP Response Error"
)

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
	}
}

func (vm *VM) initErrorClasses() {
	errTypes := []string{InternalError, ArgumentError, NameError, TypeError, UndefinedMethodError, UnsupportedMethodError, ConstantAlreadyInitializedError}

	for _, errType := range errTypes {
		c := vm.initializeClass(errType, false)
		vm.objectClass.setClassConstant(c)
	}
}

/*
	Here defines different error message formats for different types of errors
*/
const (
	WrongNumberOfArgumentFormat = "Expect %d arguments. got: %d"
	WrongArgumentTypeFormat     = "Expect argument to be %s. got: %s"
	CantYieldWithoutBlockFormat = "Can't yield without a block"
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
}

// Polymorphic helper functions -----------------------------------------
func (e *Error) toString() string {
	return "ERROR: " + e.Message
}

func (e *Error) toJSON() string {
	return e.toString()
}
