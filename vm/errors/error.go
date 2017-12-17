package errors

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
	// ConstantAlreadyInitializedError means user re-declares twice
	ConstantAlreadyInitializedError = "ConstantAlreadyInitializedError"
	// HTTPError is returned when when a request fails to return a proper response
	HTTPError = "HTTPError"
	// For zero-division by Integer or Decimal value
	ZeroDivisionError = "ZeroDivisionError"
)

/*
	Here defines different error message formats for different types of errors
*/
const (
	WrongNumberOfArgumentFormat = "Expect %d arguments. got: %d"
	WrongArgumentTypeFormat     = "Expect argument to be %s. got: %s"
	CantYieldWithoutBlockFormat = "Can't yield without a block"
	DividedByZero               = "Divided by 0"
)
