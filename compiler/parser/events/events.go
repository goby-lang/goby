package events

import (
	"github.com/gooby-lang/gooby/compiler/parser/states"
)

// These are state machine's events
const (
	BackToNormal     = "backToNormal"
	ParseFuncCall    = "parseFuncCall"
	ParseMethodParam = "parseMethodParam"
	ParseAssignment  = "parseAssignment"
)

// EventTable is the mapping of state and its corresponding event
var EventTable = map[string]string{
	states.Normal:             BackToNormal,
	states.ParsingFuncCall:    ParseFuncCall,
	states.ParsingMethodParam: ParseMethodParam,
	states.ParsingAssignment:  ParseAssignment,
}
