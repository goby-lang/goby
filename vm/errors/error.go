package errors

// Please be noted that only global errors should be placed here.
const (
	// InternalError is the default error type
	InternalError = "InternalError"

	CantGetInstructionSet = "Can't get instruction set from method: %s"
	CantCallClass         = "Can't call #class on %T"
)

const (
	// BlockError is for a generic block error
	BlockError = "BlockError"

	CantGetBlockWithoutBlockArgument = "Can't get block without a block argument"
)

const (
	// TypeError is for a type(meaning class)-related error
	TypeError = "TypeError"

	WrongArgumentTypeFormat = "Expects argument to be %s. got: %s"
)

const (
	// ConstantError is for a constant-related error
	ConstantError = "ConstantError"

	UninitializedConstant           = "Uninitialized constant: %s"
	ConstantIsNotClass              = "Constant %s is not a class. got: %s"
	ConstantAlreadyInitializedError = "Constant %s already been initialized. Can't assign value to a constant twice."
)

const (
	// UndefinedMethodError is for an unexpected undefined-method error
	UndefinedMethodError = "UndefinedMethodError"

	UndefinedMethodFor = "Undefined method '%s' for %s"
)

const (
	// UnsupportedFeatureError is for an intentionally unsupported-feature error
	UnsupportedFeatureError = "UnsupportedFeatureError"

	UnsupportedMethodFor         = "Method %s for %+v is unsupported"
	ModuleInheritanceUnsupported = "Module inheritance is unsupported: %s"
)

const (
	// ArgumentError is for an argument-related error
	ArgumentError = "ArgumentError"

	WrongNumberOfArgumentFormat      = "Expects %d argument(s). got: %d"
	WrongNumberOfArgumentFormatMore  = "Expects %d or more argument(s). got: %d"
	WrongNumberOfArgumentFormatRange = "Expects %d to %d argument(s). got: %d"
)

const (
	// ZeroDivisionError is for zero-division by Integer/Float/Decimal value
	ZeroDivisionError = "ZeroDivisionError"

	DividedByZero = "Divided by 0"
)

const (
	// StopIterationError is raised when there are no more elements in an iterator
	StopIterationError = "StopIterationError"
)
