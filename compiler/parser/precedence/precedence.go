package precedence

import "github.com/goby-lang/goby/compiler/token"

// Constants for denoting precedence
const (
	_ = iota
	Lowest
	Normal
	Assign
	Logic
	Range
	Equals
	Compare
	Sum
	Product
	BangPrefix
	Index
	Call
	MinusPrefix
)

// LookupTable maps token to its corresponding precedence
var LookupTable = map[token.Type]int{
	token.Eq:                 Equals,
	token.NotEq:              Equals,
	token.LT:                 Compare,
	token.LTE:                Compare,
	token.GT:                 Compare,
	token.GTE:                Compare,
	token.COMP:               Compare,
	token.And:                Logic,
	token.Or:                 Logic,
	token.Range:              Range,
	token.Plus:               Sum,
	token.Minus:              Sum,
	token.Modulo:             Sum,
	token.Slash:              Product,
	token.Asterisk:           Product,
	token.Pow:                Product,
	token.LBracket:           Index,
	token.Dot:                Call,
	token.LParen:             Call,
	token.ResolutionOperator: Call,
	token.Assign:             Assign,
	token.PlusEq:             Assign,
	token.MinusEq:            Assign,
	token.OrEq:               Assign,
	token.Colon:              Assign,
}
