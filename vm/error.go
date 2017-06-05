package vm

import (
  "fmt"
)

var (
  // UndefinedMethodErrorClass ...
  UndefinedMethodErrorClass *RError
  // ArgumentErrorClass ...
  ArgumentErrorClass *RError
)

func initErrors() {
  bc := &BaseClass{Name: "UndefinedMethodError", Methods: newEnvironment(), ClassMethods: newEnvironment(), Class: classClass, pseudoSuperClass: objectClass, superClass: objectClass}
  UndefinedMethodErrorClass = &RError{BaseClass: bc}
  bc = &BaseClass{Name: "ArgumentError", Methods: newEnvironment(), ClassMethods: newEnvironment(), Class: classClass, pseudoSuperClass: objectClass, superClass: objectClass}
  ArgumentErrorClass = &RError{BaseClass: bc}
}

// RError ...
type RError struct {
  *BaseClass
}

// Error ...
type Error struct {
  Class   *RError
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
  Class *RError
  // *BaseClass
  // *Error
  Message string
}

// ArgumentErrorObject ...
type ArgumentErrorObject struct {
  Class *RError
  // *BaseClass
  // *Error
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
