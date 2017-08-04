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
)

func (vm *VM) initErrorObject(errorType, format string, args ...interface{}) *Error {
	errClass := vm.objectClass.getClassConstant(errorType)

	return &Error{
		baseObj: &baseObj{class: errClass},
		Message: fmt.Sprintf(errorType+": "+format, args...),
	}
}

func (vm *VM) initErrorClasses() {
	errTypes := []string{InternalError, ArgumentError, NameError, TypeError, UndefinedMethodError, UnsupportedMethodError}

	for _, errType := range errTypes {
		c := vm.initializeClass(errType, false)
		vm.topLevelClass(objectClass).setClassConstant(c)
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
