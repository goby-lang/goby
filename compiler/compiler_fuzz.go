// +build gofuzz

package compiler

import (
	"github.com/goby-lang/goby/compiler/parser"
)

func Fuzz(input []byte) int {
	_, err := CompileToInstructions(string(input), parser.NormalMode)
	if err != nil {
		return 0
	}
	return 1
}
