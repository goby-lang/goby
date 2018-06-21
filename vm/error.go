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
	Type         string
}

// Internal functions ===================================================

// Functions for initialization -----------------------------------------

func (vm *VM) initUnsupportedMethodError(sourceLine int, methodName string, receiver Object) *Error {
	return vm.InitErrorObject(errors.UnsupportedMethodError, sourceLine, "Unsupported Method %s for %+v", methodName, receiver.ToString())
}

func (vm *VM) InitErrorObject(errorType string, sourceLine int, format string, args ...interface{}) *Error {
	errClass := vm.objectClass.getClassConstant(errorType)

	t := vm.mainThread
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
		BaseObj: &BaseObj{class: errClass},
		// Add 1 to source line because it's zero indexed
		message:     fmt.Sprintf(errorType+": "+format, args...),
		stackTraces: []string{fmt.Sprintf("from %s:%d", cf.FileName(), sourceLine)},
		Type:        errorType,
	}
}

func (vm *VM) initErrorClasses() {
	errTypes := []string{errors.InternalError, errors.ArgumentError, errors.NameError, errors.StopIteration, errors.TypeError, errors.UndefinedMethodError, errors.UnsupportedMethodError, errors.ConstantAlreadyInitializedError, errors.HTTPError, errors.ZeroDivisionError, errors.ChannelCloseError}

	for _, errType := range errTypes {
		c := vm.initializeClass(errType)
		vm.objectClass.setClassConstant(c)
	}
}

// Polymorphic helper functions -----------------------------------------

// ToString returns the object's name as the string format
func (e *Error) ToString() string {
	return e.message
}

// ToJSON just delegates to `ToString`
func (e *Error) ToJSON(t *Thread) string {
	return e.ToString()
}

func (e *Error) Value() interface{} {
	return e.message
}

// Message prints the error's message and its stack traces
func (e *Error) Message() string {
	return e.message + "\n" + strings.Join(e.stackTraces, "\n")
}
