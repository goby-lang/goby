package events

import (
	"github.com/goby-lang/goby/compiler/parser/states"
)

// These are state machine's events
const (
	BackToNormal     = "backToNormal"
	ParseFuncCall    = "parseFuncCall"
	ParseMethodParam = "parseMethodParam"
	ParseAssignment  = "parseAssignment"
)

var EventTable = map[string]string{
	states.Normal:             BackToNormal,
	states.ParsingFuncCall:    ParseFuncCall,
	states.ParsingMethodParam: ParseMethodParam,
	states.ParsingAssignment:  ParseAssignment,
}
