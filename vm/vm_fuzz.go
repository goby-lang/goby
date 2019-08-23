// +build gofuzz

package vm

import (
	"strings"

	"github.com/goby-lang/goby/compiler"
	"github.com/goby-lang/goby/compiler/parser"
)

// Fuzz tests evaluation
func Fuzz(fuzz []byte) int {
	source, ok := decodeSourceFuzz(fuzz)
	if !ok {
		return -1
	}
	iss, err := compiler.CompileToInstructions(source, parser.TestMode)
	if err != nil {
		return 0
	}
	vm := initTestVM()
	vm.ExecInstructions(iss, "fuzzer")
	return 1
}

// decodeSourceFuzz maps fuzz to Goby source code
func decodeSourceFuzz(fuzz []byte) (source string, ok bool) {
	var builder strings.Builder
	var value string
	for len(fuzz) > 0 {
		switch fuzz[0] {
		case 0:
			// Constant
			value, fuzz, ok = extractConstant(fuzz)
			if !ok {
				return "", false
			}
			builder.WriteString(value)
		case 1:
			// Ident
			value, fuzz, ok = extractIdent(fuzz)
			if !ok {
				return "", false
			}
			builder.WriteString(value)
		case 2:
			// InstanceVariable
			value, fuzz, ok = extractInstanceVariable(fuzz)
			if !ok {
				return "", false
			}
			builder.WriteString("@")
			builder.WriteString(value)
		case 3:
			// Int
			value, fuzz, ok = extractInt(fuzz)
			if !ok {
				return "", false
			}
			builder.WriteString(value)
		case 4:
			// Float
			value, fuzz, ok = extractFloat(fuzz)
			if !ok {
				return "", false
			}
			builder.WriteString(value)
		case 5:
			// Comment
			value, fuzz, ok = extractComment(fuzz)
			if !ok {
				return "", false
			}
			builder.WriteString("#")
			builder.WriteString(value)
		case 6:
			// String, double quotes
			value, fuzz, ok = extractString(fuzz)
			if !ok {
				return "", false
			}
			builder.WriteString("\"")
			builder.WriteString(value)
			builder.WriteString("\"")
		case 7:
			// String, single quotes
			value, fuzz, ok = extractString(fuzz)
			if !ok {
				return "", false
			}
			builder.WriteString("'")
			builder.WriteString(value)
			builder.WriteString("'")
		case 8:
			// String, colon
			value, fuzz, ok = extractString(fuzz)
			if !ok {
				return "", false
			}
			builder.WriteString(":")
			builder.WriteString(value)
		case 9:
			// Space
			builder.WriteString(" ")
			fuzz = fuzz[1:]
		case 10:
			// Assign
			builder.WriteString("=")
			fuzz = fuzz[1:]
		case 11:
			// Plus
			builder.WriteString("+")
			fuzz = fuzz[1:]
		case 12:
			// PlusEq
			builder.WriteString("+=")
			fuzz = fuzz[1:]
		case 13:
			// Minus
			builder.WriteString("-")
			fuzz = fuzz[1:]
		case 14:
			// MinusEq
			builder.WriteString("-=")
			fuzz = fuzz[1:]
		case 15:
			// Bang
			builder.WriteString("!")
			fuzz = fuzz[1:]
		case 16:
			// Asterisk
			builder.WriteString("*")
			fuzz = fuzz[1:]
		case 17:
			// Pow
			builder.WriteString("**")
			fuzz = fuzz[1:]
		case 18:
			// Slash
			builder.WriteString("/")
			fuzz = fuzz[1:]
		case 19:
			// Dot
			builder.WriteString(".")
			fuzz = fuzz[1:]
		case 20:
			// And
			builder.WriteString("&&")
			fuzz = fuzz[1:]
		case 21:
			// Or
			builder.WriteString("||")
			fuzz = fuzz[1:]
		case 22:
			// OrEq
			builder.WriteString("||=")
			fuzz = fuzz[1:]
		case 23:
			// Modulo
			builder.WriteString("%")
			fuzz = fuzz[1:]
		case 24:
			// Match
			builder.WriteString("=~")
			fuzz = fuzz[1:]
		case 25:
			// LT
			builder.WriteString("<")
			fuzz = fuzz[1:]
		case 26:
			// LTE
			builder.WriteString("<=")
			fuzz = fuzz[1:]
		case 27:
			// GT
			builder.WriteString(">")
			fuzz = fuzz[1:]
		case 28:
			// GTE
			builder.WriteString(">=")
			fuzz = fuzz[1:]
		case 29:
			// COMP
			builder.WriteString("<=>")
			fuzz = fuzz[1:]
		case 30:
			// Comma
			builder.WriteString(",")
			fuzz = fuzz[1:]
		case 31:
			// Semicolon
			builder.WriteString(";")
			fuzz = fuzz[1:]
		case 32:
			// Colon
			builder.WriteString(":")
			fuzz = fuzz[1:]
		case 33:
			// Bar
			builder.WriteString("|")
			fuzz = fuzz[1:]
		case 34:
			// LParen
			builder.WriteString("(")
			fuzz = fuzz[1:]
		case 35:
			// RParen
			builder.WriteString(")")
			fuzz = fuzz[1:]
		case 36:
			// LBrace
			builder.WriteString("{")
			fuzz = fuzz[1:]
		case 37:
			// RBrace
			builder.WriteString("}")
			fuzz = fuzz[1:]
		case 38:
			// LBracket
			builder.WriteString("[")
			fuzz = fuzz[1:]
		case 39:
			// RBracket
			builder.WriteString("]")
			fuzz = fuzz[1:]
		case 40:
			// Eq
			builder.WriteString("==")
			fuzz = fuzz[1:]
		case 41:
			// NotEq
			builder.WriteString("!=")
			fuzz = fuzz[1:]
		case 42:
			// Range
			builder.WriteString("..")
			fuzz = fuzz[1:]
		case 43:
			// True
			builder.WriteString("true")
			fuzz = fuzz[1:]
		case 44:
			// False
			builder.WriteString("false")
			fuzz = fuzz[1:]
		case 45:
			// Null
			builder.WriteString("nil")
			fuzz = fuzz[1:]
		case 46:
			// If
			builder.WriteString("if")
			fuzz = fuzz[1:]
		case 47:
			// ElsIf
			builder.WriteString("elsif")
			fuzz = fuzz[1:]
		case 48:
			// Else
			builder.WriteString("else")
			fuzz = fuzz[1:]
		case 49:
			// Case
			builder.WriteString("case")
			fuzz = fuzz[1:]
		case 50:
			// When
			builder.WriteString("when")
			fuzz = fuzz[1:]
		case 51:
			// Return
			builder.WriteString("return")
			fuzz = fuzz[1:]
		case 52:
			// Next
			builder.WriteString("next")
			fuzz = fuzz[1:]
		case 53:
			// Break
			builder.WriteString("break")
			fuzz = fuzz[1:]
		case 54:
			// Def
			builder.WriteString("def")
			fuzz = fuzz[1:]
		case 55:
			// Self
			builder.WriteString("self")
			fuzz = fuzz[1:]
		case 56:
			// End
			builder.WriteString("end")
			fuzz = fuzz[1:]
		case 57:
			// While
			builder.WriteString("while")
			fuzz = fuzz[1:]
		case 58:
			// Do
			builder.WriteString("do")
			fuzz = fuzz[1:]
		case 59:
			// Yield
			builder.WriteString("yield")
			fuzz = fuzz[1:]
		case 60:
			// GetBlock
			builder.WriteString("get_block")
			fuzz = fuzz[1:]
		case 61:
			// Class
			builder.WriteString("class")
			fuzz = fuzz[1:]
		case 62:
			// Module
			builder.WriteString("module")
			fuzz = fuzz[1:]
		case 63:
			// ResolutionOperator
			builder.WriteString("::")
			fuzz = fuzz[1:]
		default:
			return "", false
		}
	}
	return builder.String(), true
}

// charIdentifierStart detects a valid identifier start character
func charIdentifierStart(value byte) bool {
	return charUnderscore(value) || charLetter(value)
}

// charIdentifier detects a valid identifier character
func charIdentifier(value byte) bool {
	return charUnderscore(value) || charDigit(value) || charLetter(value)
}

// charSimple detects a simple character, valid in strings or comments
func charSimple(value byte) bool {
	return charUnderscore(value) || charDigit(value) || charLetter(value)
}

// charUnderscore detects an underscore
func charUnderscore(value byte) bool {
	return value == 95
}

// charDash detects a dash
func charDash(value byte) bool {
	return value == 45
}

// charDot detects a dot
func charDot(value byte) bool {
	return value == 46
}

// charDigit detects a digit
func charDigit(value byte) bool {
	return value >= 48 && value <= 57
}

// charLetter detects a letter
func charLetter(value byte) bool {
	return charUppercaseLetter(value) || charLowercaseLetter(value)
}

// charUppercaseLetter detects an uppercase letter
func charUppercaseLetter(value byte) bool {
	return value >= 65 && value <= 90
}

// charLowercaseLetter detects a lowercase letter
func charLowercaseLetter(value byte) bool {
	return value >= 97 && value <= 122
}

// extractValue extracts a length prefixed value
func extractValue(fuzz []byte) ([]byte, []byte, bool) {
	if len(fuzz) <= 2 {
		// Prohibit incomplete construct
		return nil, nil, false
	}
	length := uint8(fuzz[1])
	if length == 0 {
		// Prohibit empty name
		return nil, nil, false
	}
	if len(fuzz) < int(length+2) {
		// Prohibit length past end
		return nil, nil, false
	}
	value := fuzz[2 : length+2]
	fuzz = fuzz[2+length:]
	return value, fuzz, true
}

// extractConstant extracts a constant name
func extractConstant(fuzz []byte) (string, []byte, bool) {
	value, fuzz, ok := extractValue(fuzz)
	if !ok {
		return "", nil, false
	}
	if !charUppercaseLetter(value[0]) {
		// Require uppercase letter at start
		return "", nil, false
	}
	if len(value) > 1 {
		for _, char := range value[1:] {
			if !charIdentifier(char) {
				// Require valid identifier character
				return "", nil, false
			}
		}
	}
	return string(value), fuzz, ok
}

// extractIdent extracts an identifier
func extractIdent(fuzz []byte) (string, []byte, bool) {
	value, fuzz, ok := extractValue(fuzz)
	if !ok {
		return "", nil, false
	}
	for _, char := range value {
		if !charIdentifier(char) {
			// Require valid identifier character
			return "", nil, false
		}
	}
	return string(value), fuzz, ok
}

// extractInstanceVariable extracts an instance variable name
func extractInstanceVariable(fuzz []byte) (string, []byte, bool) {
	value, fuzz, ok := extractValue(fuzz)
	if !ok {
		return "", nil, false
	}
	for _, char := range value {
		if !charIdentifier(char) {
			// Require valid identifier character
			return "", nil, false
		}
	}
	return string(value), fuzz, ok
}

// extractInt extracts an integer literal
func extractInt(fuzz []byte) (string, []byte, bool) {
	value, fuzz, ok := extractValue(fuzz)
	if !ok {
		return "", nil, false
	}
	if !(charDigit(value[0]) || charDash(value[0])) {
		// Require digit or dash at start
		return "", nil, false
	}
	if len(value) == 0 && charDash(value[0]) {
		// Prohibit dash without digits
		return "", nil, false
	}
	if len(value) > 1 {
		for _, char := range value[1:] {
			if !charDigit(char) {
				// Require digit
				return "", nil, false
			}
		}
	}
	return string(value), fuzz, ok
}

// extractFloat extracts a float literal
func extractFloat(fuzz []byte) (string, []byte, bool) {
	value, fuzz, ok := extractValue(fuzz)
	if !ok {
		return "", nil, false
	}
	if !(charDigit(value[0]) || charDash(value[0])) {
		// Require digit or dash at start
		return "", nil, false
	}
	if len(value) == 0 && charDash(value[0]) {
		// Porhibit dash without digits
		return "", nil, false
	}
	if len(value) > 1 {
		for _, char := range value[1:] {
			if !(charDigit(char) || charDot(char)) {
				// Require digit or dot
				return "", nil, false
			}
		}
	}
	var dots uint8
	for _, char := range value {
		if charDot(char) {
			dots++
		}
	}
	if dots > 1 {
		// Prohibit multiple dots
		return "", nil, false
	}
	return string(value), fuzz, ok
}

// extractComment extracts comment content
func extractComment(fuzz []byte) (string, []byte, bool) {
	value, fuzz, ok := extractValue(fuzz)
	if !ok {
		return "", nil, false
	}
	for _, char := range value {
		// Require simple character
		if !charSimple(char) {
			return "", nil, false
		}
	}
	return string(value), fuzz, ok
}

// extractString extracts string content
func extractString(fuzz []byte) (string, []byte, bool) {
	value, fuzz, ok := extractValue(fuzz)
	if !ok {
		return "", nil, false
	}
	for _, char := range value {
		// Require simple character
		if !charSimple(char) {
			return "", nil, false
		}
	}
	return string(value), fuzz, ok
}
