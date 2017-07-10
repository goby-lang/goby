package vm

import (
	"fmt"
)

var (
	// InternalErrorClass ...
	InternalErrorClass *RClass
	// ArgumentErrorClass ...
	ArgumentErrorClass *RClass
	// NameErrorClass ...
	NameErrorClass *RClass
	// TypeErrorClass ...
	TypeErrorClass *RClass
	// UndefinedMethodErrorClass ...
	UndefinedMethodErrorClass *RClass
	// UnsupportedMethodClass ...
	UnsupportedMethodClass *RClass
)

const (
	// InternalError is the default error type
	InternalError = "InternalError"
	// ArgumentError: an argument-related error
	ArgumentError = "ArgumentError"
	// NameError: a constant-related error
	NameError = "NameError"
	// TypeError: a type-related error
	TypeError = "TypeError"
	// UndefinedMethodError: undefined-method error
	UndefinedMethodError = "UndefinedMethodError"
	// UnsupportedMethodError: intentionally unsupported-error
	UnsupportedMethodError = "UnsupportedMethodError"
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

func initErrorObject(errorType *RClass, format string, args ...interface{}) *Error {
	return &Error{
		baseObj: &baseObj{class: errorType},
		Message: fmt.Sprintf(errorType.Name+": "+format, args...),
	}
}

func (vm *VM) initErrorClasses() {
	InternalErrorClass = vm.initializeClass(InternalError, false)
	ArgumentErrorClass = vm.initializeClass(ArgumentError, false)
	NameErrorClass = vm.initializeClass(NameError, false)
	TypeErrorClass = vm.initializeClass(TypeError, false)
	UndefinedMethodErrorClass = vm.initializeClass(UndefinedMethodError, false)
	UnsupportedMethodClass = vm.initializeClass(UnsupportedMethodError, false)
}

// Polymorphic helper functions -----------------------------------------

// toString converts error messages into string.
func (e *Error) toString() string {
	return "ERROR: " + e.Message
}

// toJSON converts the receiver into JSON string.
func (e *Error) toJSON() string {
	return e.toString()
}
