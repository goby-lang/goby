package vm

import (
	"fmt"
)

var (
	// UndefinedMethodErrorClass ...
	UndefinedMethodErrorClass *RClass
	// ArgumentErrorClass ...
	ArgumentErrorClass *RClass
	// TypeErrorClass ...
	TypeErrorClass *RClass
)

const (
	// UndefinedMethodError describes the error type in string
	UndefinedMethodError = "UndefinedMethodError"
	// ArgumentError describes the error type in string
	ArgumentError = "ArgumentError"
	// TypeError describes the error type in string
	TypeError = "TypeError"
)

func init() {
	bc := createBaseClass(UndefinedMethodError)
	UndefinedMethodErrorClass = &RClass{BaseClass: bc}
	bc = createBaseClass(ArgumentError)
	ArgumentErrorClass = &RClass{BaseClass: bc}
	bc = createBaseClass(TypeError)
	TypeErrorClass = &RClass{BaseClass: bc}
}

// Error ...
type Error struct {
	Class   *RClass
	Message string
}

// Inspect ...
func (e *Error) Inspect() string {
	return "ERROR: " + e.Message
}

func (e *Error) toJSON() string {
	return e.Inspect()
}

func (e *Error) returnClass() Class {
	return e.Class
}

func initializeError(errorType *RClass, format string, args ...interface{}) *Error {
	return &Error{
		Class:   errorType,
		Message: fmt.Sprintf(errorType.Name+": "+format, args...),
	}
}
