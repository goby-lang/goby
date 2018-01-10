package arguments

import "github.com/goby-lang/goby/compiler/token"

const (
	NormalArg = iota
	OptionedArg
	SplatArg
	RequiredKeywordArg
	OptionalKeywordArg
)

var Types = map[int]string{
	NormalArg:          "Normal argument",
	OptionedArg:        "Optioned argument",
	RequiredKeywordArg: "Keyword argument",
	OptionalKeywordArg: "Optioned keyword argument",
	SplatArg:           "Splat argument",
}

var Tokens = map[token.Type]bool{
	token.Int:              true,
	token.String:           true,
	token.True:             true,
	token.False:            true,
	token.Null:             true,
	token.InstanceVariable: true,
	token.Ident:            true,
	token.Constant:         true,
}
