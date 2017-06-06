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

func initErrors() {
	bc := createBaseClass("UndefinedMethodError")
	UndefinedMethodErrorClass = &RClass{BaseClass: bc}
	bc = createBaseClass("ArgumentError")
	ArgumentErrorClass = &RClass{BaseClass: bc}
	bc = createBaseClass("TypeError")
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

func (e *Error) returnClass() Class {
	return e.Class
}

// UndefinedMethodErrorObject ...
type UndefinedMethodErrorObject struct {
	Class   *RClass
	Message string
}

// Inspect ...
func (e *UndefinedMethodErrorObject) Inspect() string {
	return "ArgumentError: " + e.Message
}

func (e *UndefinedMethodErrorObject) returnClass() Class {
	return e.Class
}

func initializeUndefinedMethodError(format string, args ...interface{}) *UndefinedMethodErrorObject {
	return &UndefinedMethodErrorObject{Class: UndefinedMethodErrorClass, Message: fmt.Sprintf(format, args...)}
}

// ArgumentErrorObject ...
type ArgumentErrorObject struct {
	Class   *RClass
	Message string
}

// Inspect ...
func (e *ArgumentErrorObject) Inspect() string {
	return "ArgumentError: " + e.Message
}

func (e *ArgumentErrorObject) returnClass() Class {
	return e.Class
}

func initializeArgumentError(format string, args ...interface{}) *ArgumentErrorObject {
	return &ArgumentErrorObject{Class: ArgumentErrorClass, Message: fmt.Sprintf(format, args...)}
}

// TypeErrorObject ...
type TypeErrorObject struct {
	Class   *RClass
	Message string
}

// Inspect ...
func (e *TypeErrorObject) Inspect() string {
	return "TypeError: " + e.Message
}

func (e *TypeErrorObject) returnClass() Class {
	return e.Class
}

func initializeTypeError(format string, args ...interface{}) *TypeErrorObject {
	return &TypeErrorObject{Class: TypeErrorClass, Message: fmt.Sprintf(format, args...)}
}
