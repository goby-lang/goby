package vm

import (
	"fmt"
)

// Nothing to describe here, just error classes
var (
	// UndefinedMethodErrorClass ...
	UndefinedMethodErrorClass *RClass
	// ArgumentErrorClass ...
	ArgumentErrorClass *RClass
	// TypeErrorClass ...
	TypeErrorClass *RClass
	// UnsupportedMethodClass ...
	UnsupportedMethodClass *RClass
	// NameErrorClass ...
	NameErrorClass *RClass
	// InternalErrorClass ...
	InternalErrorClass *RClass
)

const (
	// UndefinedMethodError describes the error type in string
	UndefinedMethodError = "UndefinedMethodError"
	// ArgumentError describes the error type in string
	ArgumentError = "ArgumentError"
	// TypeError describes the error type in string
	TypeError = "TypeError"
	// UnsupportedMethodError describes the error type in string
	UnsupportedMethodError = "UnsupportedMethodError"
	// NameError describes constant related errors
	NameError = "NameError"
	// InternalError is the default error type
	InternalError = "InternalError"
)

func (vm *VM) initErrorClasses() {
	UndefinedMethodErrorClass = vm.initializeClass(UndefinedMethodError, false)
	ArgumentErrorClass = vm.initializeClass(ArgumentError, false)
	TypeErrorClass = vm.initializeClass(TypeError, false)
	UnsupportedMethodClass = vm.initializeClass(UnsupportedMethodError, false)
	NameErrorClass = vm.initializeClass(NameError, false)
	InternalErrorClass = vm.initializeClass(InternalError, false)
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
