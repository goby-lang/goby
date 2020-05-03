package vm

import (
	"fmt"

	"strings"

	"github.com/goby-lang/goby/vm/errors"
)

// Error class is actually a special struct to hold internal error types with messages.
// Goby developers need not to take care of the struct.
// Goby maintainers should consider using the appropriate error type.
// Cannot create instances of Error class, or inherit Error class.
//
// The type of internal errors:
//
// see vm/errors/error.go.
//
type Error struct {
	*BaseObj
	message      string
	stackTraces  []string
	storedTraces bool
}

// Internal functions ===================================================

// Functions for initialization -----------------------------------------

// InitNoMethodError is to print unsupported method errors. This is exported for using from sub-packages.
func (vm *VM) InitNoMethodError(sourceLine int, methodName string, receiver Object) *Error {
	return vm.InitErrorObject(errors.NoMethodError, sourceLine, errors.UndefinedMethod, methodName, receiver.Inspect())
}

// InitErrorObject initializes and returns Error object
func (vm *VM) InitErrorObjectFromClass(errClass *RClass, sourceLine int, format string, args ...interface{}) *Error {
	t := &vm.mainThread
	cf := t.callFrameStack.top()

	switch cf := cf.(type) {
	case *normalCallFrame:
		// If program counter is 0 means we need to trace back to previous call frame
		if cf.pc == 0 {
			t.callFrameStack.pop()
			cf = t.callFrameStack.top().(*normalCallFrame)
		}
	}

	return &Error{
		BaseObj: NewBaseObject(errClass),
		// Add 1 to source line because it's zero indexed
		message:     fmt.Sprintf(errClass.Name+": "+format, args...),
		stackTraces: []string{fmt.Sprintf("from %s:%d", cf.FileName(), sourceLine)},
	}
}

// InitErrorObject initializes and returns Error object
func (vm *VM) InitErrorObject(errorType errors.ErrorType, sourceLine int, format string, args ...interface{}) *Error {
	en := errors.GetErrorName(errorType)
	errClass := vm.objectClass.getClassConstant(en)
	return vm.InitErrorObjectFromClass(errClass, sourceLine, format, args...)
}

func (vm *VM) initErrorClasses() {
	for _, et := range errors.AllErrorTypes() {
		en := errors.GetErrorName(et)
		c := vm.initializeClass(en)
		vm.objectClass.setClassConstant(c)
	}
}

// Polymorphic helper functions -----------------------------------------

// ToString returns the object's name as the string format
func (e *Error) ToString() string {
	return e.message
}

// Inspect delegates to ToString
func (e *Error) Inspect() string {
	return e.ToString()
}

// ToJSON just delegates to `ToString`
func (e *Error) ToJSON(t *Thread) string {
	return e.ToString()
}

// Value is equivalent to ToString
func (e *Error) Value() interface{} {
	return e.message
}

// Message prints the error's message and its stack traces
func (e *Error) Message() string {
	return e.message + "\n" + strings.Join(e.stackTraces, "\n")
}
