package utils

import (
	"fmt"
	"github.com/goby-lang/goby/compiler/ast"
)

func IsDigit(ch rune) bool {
	return '0' <= ch && ch <= '9'
}

func IsDigitString(str string) bool {
	for _, char := range str {
		if !IsDigit(char) {
			return false
		}
	}
	return true
}

func IsLetter(ch rune) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_'
}

func IsLetterString(str string) bool {
	for _, char := range str {
		if !IsLetter(char) {
			return false
		}
	}
	return true
}

func CheckNameOfVariables(variables ...ast.Expression) (isPassed bool,  errVariableIndex int) {
	isPassed = true;
	sliceIndex := 0

	for _, variable := range variables {
		variableName := variable.(ast.Variable).String()

		var firstCharAt = variableName[sliceIndex]

		if firstCharAt == '@' {
			if len(variableName) < 2 {
				isPassed =  false;
				return
			}

			firstCharAt = variableName[sliceIndex+1]
			sliceIndex++
		}
		

		if IsDigit(rune(firstCharAt)) {
			isPassed = false
			return
		}

		for index, char := range variableName[sliceIndex:] {
			if !IsDigit(rune(char)) && !IsLetter(rune(char)) {
				isPassed = false
				errVariableIndex = index
				return
			}
		}
	}
	return
}

func CheckInstanceVariable(variable string) bool {
	if (len(variable) < 2) {
		return false;
	}

	if (variable[0] != '@' || IsDigit(rune(variable[1]))) {
		return false;
	}

	for _, char := range variable[2:] {
		if !IsDigit(rune(char)) && !IsLetter(rune(char)) {
			return false
		}
	}

	return true
}

func IsInstanceVariableSymbol(ch rune) bool {
	return ch == '@'
}

func IsEscapedChar(ch rune) bool {
	return ch == '\\'
}