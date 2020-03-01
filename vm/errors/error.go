package errors

const (
	// InternalError is the default error type
	InternalError = "InternalError"
	// IOError is an IO error such as file error
	IOError = "IOError"
	// ArgumentError is for an argument-related error
	ArgumentError = "ArgumentError"
	// NameError is for a constant-related error
	NameError = "NameError"
	// StopIteration is raised when there are no more elements in an iterator
	StopIteration = "StopIteration"
	// TypeError is for a type-related error
	TypeError = "TypeError"
	// NoMethodError is for an intentionally unsupported-method error
	NoMethodError = "NoMethodError"
	// ConstantAlreadyInitializedError means user re-declares twice
	ConstantAlreadyInitializedError = "ConstantAlreadyInitializedError"
	// HTTPError is returned when when a request fails to return a proper response
	HTTPError = "HTTPError"
	// ZeroDivisionError is for zero-division by Integer/Float/Decimal value
	ZeroDivisionError = "ZeroDivisionError"
	// ChannelCloseError is for accessing to the closed channel
	ChannelCloseError = "ChannelCloseError"
	// NotImplementedError means the method is missing
	NotImplementedError = "NotImplementedError"
)

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
