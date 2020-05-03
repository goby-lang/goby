package errors

import "fmt"

// ErrorType is the enum representation for built-in error types
type ErrorType int8

const (
	// InternalError is the default error type
	InternalError ErrorType = iota
	// IOError is an IO error such as file error
	IOError
	// ArgumentError is for an argument-related error
	ArgumentError
	// NameError is for a constant-related error
	NameError
	// StopIteration is raised when there are no more elements in an iterator
	StopIteration
	// TypeError is for a type-related error
	TypeError
	// NoMethodError is for an intentionally unsupported-method error
	NoMethodError
	// ConstantAlreadyInitializedError means user re-declares twice
	ConstantAlreadyInitializedError
	// HTTPError is returned when when a request fails to return a proper response
	HTTPError
	// ZeroDivisionError is for zero-division by Integer/Float/Decimal value
	ZeroDivisionError
	// ChannelCloseError is for accessing to the closed channel
	ChannelCloseError
	// NotImplementedError means the method is missing
	NotImplementedError

	// EndOfErrorTypeConst is an anchor for getting all error types' enum values, see AllErrorTypes
	EndOfErrorTypeConst
)

var errorTypesMap = map[ErrorType]string{
	InternalError: "InternalError",
	IOError: "IOError",
	ArgumentError: "ArgumentError",
	NameError: "NameError",
	StopIteration: "StopIteration",
	TypeError: "TypeError",
	NoMethodError: "NoMethodError",
	ConstantAlreadyInitializedError: "ConstantAlreadyInitializedError",
	HTTPError: "HTTPError",
	ZeroDivisionError: "ZeroDivisionError",
	ChannelCloseError: "ChannelCloseError",
	NotImplementedError: "NotImplementedError",
}


// AllErrorTypes returns all error types defined in this package in their enum format.
func AllErrorTypes() []ErrorType {
	ts := make([]ErrorType, EndOfErrorTypeConst)
	for i := 0; i < int(EndOfErrorTypeConst); i++ {
		ts[i] = ErrorType(i)
	}
	return ts
}

// GetErrorName receives an ErrorType enum and returns the corresponding error name.
func GetErrorName(t ErrorType) string {
	v, ok := errorTypesMap[t]

	if ok {
		return v
	}

	panic(fmt.Errorf("expect to find ErrorType %d's name", t))
}

/*
	Here defines different error message formats for different types of errors
*/
const (
	WrongNumberOfArgument           = "Expect %d argument(s). got: %d"
	WrongNumberOfArgumentMore       = "Expect %d or more argument(s). got: %d"
	WrongNumberOfArgumentLess       = "Expect %d or less argument(s). got: %d"
	WrongNumberOfArgumentRange      = "Expect %d to %d argument(s). got: %d"
	WrongArgumentTypeFormat         = "Expect argument to be %s. got: %s"
	WrongArgumentTypeFormatNum      = "Expect argument #%d to be %s. got: %s"
	InvalidChmodNumber              = "Invalid chmod number. got: %d"
	InvalidNumericString            = "Invalid numeric string. got: %s"
	CantLoadFile                    = "Can't load \"%s\""
	CantRequireNonString            = "Can't require \"%s\": Pass a string instead"
	CantYieldWithoutBlockFormat     = "Can't yield without a block"
	NotDiggable                     = "Expect target to be Diggable, got %s"
	DividedByZero                   = "Divided by 0"
	ChannelIsClosed                 = "The channel is already closed."
	TooSmallIndexValue              = "Index value %d too small for array. minimum: %d"
	IndexOutOfRange                 = "Index value out of range. got: %v"
	InvalidCode                     = "invalid code: %s"
	RegexpFailure                   = "Replacement failure with the Regexp. got: %s"
	NegativeValue                   = "Expect argument to be positive value. got: %d"
	NegativeSecondValue             = "Expect second argument to be positive value. got: %d"
	NativeNotImplementedErrorFormat = "'%s' should be implemented on %s but haven't be done yet. Looking forward to see your PR for it ;-)"
	UndefinedMethod                 = "Undefined Method '%+v' for %+v"
)
