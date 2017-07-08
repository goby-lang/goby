package vm

import (
	"fmt"
)

// Nothing to describe here, just error classes
var (
	// ArgumentErrorClass ...
	ArgumentErrorClass *RClass
	// InternalErrorClass ...
	InternalErrorClass *RClass
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
	// ArgumentError describes the error type in string
	ArgumentError = "ArgumentError"
	// InternalError is the default error type
	InternalError = "InternalError"
	// NameError describes constant related errors
	NameError = "NameError"
	// TypeError describes the error type in string
	TypeError = "TypeError"
	// UndefinedMethodError describes the error type in string
	UndefinedMethodError = "UndefinedMethodError"
	// UnsupportedMethodError describes the error type in string
	UnsupportedMethodError = "UnsupportedMethodError"
)

func (vm *VM) initErrorClasses() {
	ArgumentErrorClass = vm.initializeClass(ArgumentError, false)
	InternalErrorClass = vm.initializeClass(InternalError, false)
	NameErrorClass = vm.initializeClass(NameError, false)
	TypeErrorClass = vm.initializeClass(TypeError, false)
	UndefinedMethodErrorClass = vm.initializeClass(UndefinedMethodError, false)
	UnsupportedMethodClass = vm.initializeClass(UnsupportedMethodError, false)
}

// Error ...
type Error struct {
	*baseObj
	Message string
}

// toString ...
func (e *Error) toString() string {
	return "ERROR: " + e.Message
}

func (e *Error) toJSON() string {
	return e.toString()
}

func initErrorObject(errorType *RClass, format string, args ...interface{}) *Error {
	return &Error{
		baseObj: &baseObj{class: errorType},
		Message: fmt.Sprintf(errorType.Name+": "+format, args...),
	}
}
