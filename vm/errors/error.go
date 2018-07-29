package errors

const (
	// InternalError is the default error type
	InternalError = "InternalError"
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
)

/*
	Here defines different error message formats for different types of errors
*/
const (
	WrongNumberOfArgumentFormat = "Expect %d arguments. got: %d"
	WrongArgumentTypeFormat     = "Expect argument to be %s. got: %s"
	CantYieldWithoutBlockFormat = "Can't yield without a block"
	DividedByZero               = "Divided by 0"
	ChannelIsClosed             = "The channel is already closed."
)
