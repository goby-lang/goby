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
	// UndefinedMethodError is for an undefined-method error
	UndefinedMethodError = "UndefinedMethodError"
	// UnsupportedMethodError is for an intentionally unsupported-method error
	UnsupportedMethodError = "UnsupportedMethodError"
	// ConstantAlreadyInitializedError means user re-declares twice
	ConstantAlreadyInitializedError = "ConstantAlreadyInitializedError"
	// HTTPError is returned when when a request fails to return a proper response
	HTTPError = "HTTPError"
	// ZeroDivisionError is for zero-division by Integer/Float/Decimal value
	ZeroDivisionError = "ZeroDivisionError"
	// BlockError is for indicating errors regarding block
	BlockError = "BlockError"
	// FileError is for indicating errors regarding files
	FileError = "FileError"
)

// Error messages for native methods
const (
	WrongNumberOfArgumentFormat      = "Expects %d argument(s). got: %d"
	WrongNumberOfArgumentFormatMore  = "Expects %d or more argument(s). got: %d"
	WrongNumberOfArgumentFormatRange = "Expects %d to %d argument(s). got: %d"
	WrongArgumentTypeFormat          = "Expects argument to be %s. got: %s"
	WrongArgumentSignNegative        = "Expects argument to be positive. got: %d"

	IndexValueTooSmall = "Index value %d is too small for array. minimum: %d"
	UndefinedMethodFor = "Undefined method '%s' for %s"

	CantYieldWithoutBlockFormat = "Can't yield without a block"
	CantCallClass               = "Can't call class on %T"
	CantCompleteRequest         = "Can't complete request: %s"

	CantInitializeBlockObjectWithoutBlockArgument = "Can't initialize block object without a block argument"

	InvalidFileMode = "Invalid file mode: %"

	DividedByZero = "Divided by 0"

	CantParseStringAsJson = "Can't parse string %s as json: %s"

	CantCreatePlubing = "Can't create plugin: %s"

	ShouldBePositiveInteger = "Step should be a positive integer."
)

// Error messages for VM internal errors
const (
	UninitializedConstant            = "Uninitialized constant: %s"
	CantGetInstructionSet            = "Can't get instruction set from method: %s"
	ConstantIsNotClass               = "Constant %s is not a class. got: %s"
	ModuleInheritanceUnsupported     = "Module inheritance is unsupported: %s"
	CantGetBlockWithoutBlockArgument = "Can't get block without a block argument"
)
