package errors

import (
	"fmt"
	"github.com/goby-lang/goby/compiler/parser/arguments"
)

// Enums for different kinds of syntax errors
const (
	_ = iota
	// EndOfFileError represents normal EOF error
	EndOfFileError
	// UnexpectedTokenError means that token is not what we expected
	UnexpectedTokenError
	// UnexpectedEndError means we get unexpected "end" keyword (this is mainly created for REPL)
	UnexpectedEndError
	// MethodDefinitionError means there's an error on method definition's method name
	MethodDefinitionError
	// InvalidAssignmentError means user assigns value to wrong type of expressions
	InvalidAssignmentError
	// SyntaxError means there's a grammatical in the source code
	SyntaxError
	// ArgumentError means there's a method parameter's definition error
	ArgumentError
)

// Error represents parser's parsing error
type Error struct {
	// Message contains the readable message of error
	Message string
	ErrType int
}

// IsEOF checks if error is end of file error
func (e *Error) IsEOF() bool {
	return e.ErrType == EndOfFileError
}

// IsUnexpectedEnd checks if error is unexpected "end" keyword error
func (e *Error) IsUnexpectedEnd() bool {
	return e.ErrType == UnexpectedEndError
}

// InitError is a helper function for easily initializing error object
func InitError(msg string, errType int) *Error {
	return &Error{Message: msg, ErrType: errType}
}

// NewArgumentError is a helper function the helps initializing argument errors
func NewArgumentError(formerArgType, laterArgType int, argLiteral string, line int) *Error {
	formerArg := arguments.Types[formerArgType]
	laterArg := arguments.Types[laterArgType]
	msg := fmt.Sprintf("%s \"%s\" should be defined before %s. Line: %d", formerArg, argLiteral, laterArg, line)
	return InitError(msg, ArgumentError)
}

// NewTypeParsingError is a helper function the helps initializing type parsing errors
func NewTypeParsingError(tokenLiteral, targetType string, line int) *Error {
	msg := fmt.Sprintf("could not parse %q as %s. Line: %d", tokenLiteral, targetType, line)
	return InitError(msg, SyntaxError)
}