package arguments

import "github.com/gooby-lang/gooby/compiler/token"

// Enums for different kinds of arguments
const (
	NormalArg = iota
	OptionedArg
	SplatArg
	RequiredKeywordArg
	OptionalKeywordArg
)

// Types is a table maps argument types enum to the their real name
var Types = map[int]string{
	NormalArg:          "Normal argument",
	OptionedArg:        "Optioned argument",
	RequiredKeywordArg: "Keyword argument",
	OptionalKeywordArg: "Optioned keyword argument",
	SplatArg:           "Splat argument",
}

// Tokens marks token types that can be used as method call arguments
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
