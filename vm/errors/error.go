package errors

// Please be noted that only global errors should be placed here.
const (
	// InternalError is the default error type
	InternalError = "InternalError"

	// CantGetInstructionSet is for internal error message literal
	CantGetInstructionSet = "Can't get instruction set from method: %s"
	// CantCallClass is for internal error message literal
	CantCallClass = "Can't call #class on %T"
)

const (
	// BlockError is for a generic block error
	BlockError = "BlockError"

	// CantGetBlockWithoutBlockArgument is for block error message literal
	CantGetBlockWithoutBlockArgument = "Can't get block without a block argument"
)

const (
	// TypeError is for a type(meaning class)-related error
	TypeError = "TypeError"

	// WrongArgumentTypeFormat is for type error message literal
	WrongArgumentTypeFormat = "Expects argument to be %s. got: %s"
)

const (
	// ConstantError is for a constant-related error
	ConstantError = "ConstantError"

	// UninitializedConstant is for onstant error message literal
	UninitializedConstant = "Uninitialized constant: %s"
	// ConstantIsNotClass is for onstant error message literal
	ConstantIsNotClass = "Constant %s is not a class. got: %s"
	// ConstantAlreadyInitializedError is for onstant error message literal
	ConstantAlreadyInitializedError = "Constant %s already been initialized. Can't assign value to a constant twice."
)

const (
	// UndefinedMethodError is for an unexpected undefined-method error
	UndefinedMethodError = "UndefinedMethodError"

	// UndefinedMethodFor is for undefined method error message literal
	UndefinedMethodFor = "Undefined method '%s' for %s"
)

const (
	// UnsupportedFeatureError is for an intentionally unsupported-feature error
	UnsupportedFeatureError = "UnsupportedFeatureError"

	// UnsupportedMethodFor is for unsupported feature error message literal
	UnsupportedMethodFor = "Method %s is unsupported for %+v"
	// ModuleInheritanceUnsupported is for unsupported feature error message literal
	ModuleInheritanceUnsupported = "Module inheritance is unsupported: %s"
)

const (
	// ArgumentError is for an argument-related error
	ArgumentError = "ArgumentError"

	// WrongNumberOfArgumentFormat is for argument error message literal
	WrongNumberOfArgumentFormat = "Expects %d argument(s). got: %d"
	// WrongNumberOfArgumentFormatMore is for argument error message literal
	WrongNumberOfArgumentFormatMore = "Expects %d or more arguments. got: %d"
	// WrongNumberOfArgumentFormatRange is for argument error message literal
	WrongNumberOfArgumentFormatRange = "Expects %d to %d argument(s). got: %d"
)

const (
	// ZeroDivisionError is for zero-division by Integer/Float/Decimal value
	ZeroDivisionError = "ZeroDivisionError"

	// DividedByZero is for Zero division error message literal
	DividedByZero = "Divided by 0"
)

const (
	// StopIterationError is raised when there are no more elements in an iterator
	StopIterationError = "StopIterationError"
)
