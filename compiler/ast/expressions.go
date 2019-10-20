package ast

import (
	"bytes"
	"fmt"
	"strings"
)

// IntegerLiteral contains the node expression and its value
type IntegerLiteral struct {
	*BaseNode
	Value int
}

func (il *IntegerLiteral) expressionNode() {}

// IntegerLiteral.TokenLiteral gets the Integer type token
func (il *IntegerLiteral) TokenLiteral() string {
	return il.Token.Literal
}

// IntegerLiteral.String gets the string format of the Integer type token
func (il *IntegerLiteral) String() string {
	return il.Token.Literal
}

// FloatLiteral contains the node expression and its value
type FloatLiteral struct {
	*BaseNode
	Value float64
}

func (il *FloatLiteral) expressionNode() {}

// FloatLiteral.TokenLiteral gets the literal of the Float type token
func (il *FloatLiteral) TokenLiteral() string {
	return il.Token.Literal
}

// FloatLiteral.String gets the string format of the Float type token
func (il *FloatLiteral) String() string {
	return il.Token.Literal
}

// StringLiteral contains the node expression and its value
type StringLiteral struct {
	*BaseNode
	Value string
}

// Define the string literal which contains the node expression and its value
func (sl *StringLiteral) expressionNode() {}

// StringLiteral.TokenLiteral gets the literal of the String type token
func (sl *StringLiteral) TokenLiteral() string {
	return sl.Token.Literal
}

// StringLiteral.String gets the string format of the String type token
func (sl *StringLiteral) String() string {
	var out bytes.Buffer

	out.WriteString("\"")
	out.WriteString(sl.Token.Literal)
	out.WriteString("\"")
	return out.String()
}

// ArrayExpression defines the array expression literal which contains the node expression and its value
type ArrayExpression struct {
	*BaseNode
	Elements []Expression
}

func (ae *ArrayExpression) expressionNode() {}

// ArrayExpression gets the literal of the Array type token
func (ae *ArrayExpression) TokenLiteral() string {
	return ae.Token.Literal
}

// ArrayExpression.String gets the string format of the Array type token
func (ae *ArrayExpression) String() string {
	var out bytes.Buffer

	out.WriteString("[")

	if len(ae.Elements) == 0 {
		out.WriteString("]")
		return out.String()
	}

	out.WriteString(ae.Elements[0].String())

	for _, elem := range ae.Elements[1:] {
		out.WriteString(", ")
		out.WriteString(elem.String())
	}

	out.WriteString("]")
	return out.String()
}

// PairExpression represents a key/value pair in method parameters or arguments
type ArgumentPairExpression struct {
	*BaseNode
	Key   Expression
	Value Expression
}

func (ape *ArgumentPairExpression) expressionNode() {}

// TokenLiteral .....
func (ape *ArgumentPairExpression) TokenLiteral() string {
	return ape.Token.Literal
}

// String .....
func (ape *ArgumentPairExpression) String() string {
	if ape.Value == nil {
		return fmt.Sprintf("%s:", ape.Key.String())
	}

	return fmt.Sprintf("%s: %s", ape.Key.String(), ape.Value.String())
}

// HashExpression defines the hash expression literal which contains the node expression and its value
type HashExpression struct {
	*BaseNode
	Data map[string]Expression
}

func (he *HashExpression) expressionNode() {}

// HashExpression.TokenLiteral gets the literal of the Hash type token
func (he *HashExpression) TokenLiteral() string {
	return he.Token.Literal
}

// HashExpression.String gets the string format of the Hash type token
func (he *HashExpression) String() string {
	var out bytes.Buffer
	var pairs []string

	for key, value := range he.Data {
		pairs = append(pairs, fmt.Sprintf("%s: %s", key, value.String()))
	}

	out.WriteString("{")
	out.WriteString(strings.Join(pairs, ", "))
	out.WriteString("}")

	return out.String()
}

type PrefixExpression struct {
	*BaseNode
	Operator string
	Right    Expression
}

func (pe *PrefixExpression) expressionNode() {}
func (pe *PrefixExpression) TokenLiteral() string {
	return pe.Token.Literal
}
func (pe *PrefixExpression) String() string {
	var out bytes.Buffer

	out.WriteString(pe.Operator)
	out.WriteString(pe.Right.String())

	return out.String()
}

type InfixExpression struct {
	*BaseNode
	Left     Expression
	Operator string
	Right    Expression
}

func (ie *InfixExpression) expressionNode() {}
func (ie *InfixExpression) TokenLiteral() string {
	return ie.Token.Literal
}
func (ie *InfixExpression) String() string {
	var out bytes.Buffer

	out.WriteString("(")
	out.WriteString(ie.Left.String())
	out.WriteString(" ")
	out.WriteString(ie.Operator)
	out.WriteString(" ")
	out.WriteString(ie.Right.String())
	out.WriteString(")")

	return out.String()
}

// AssignExpression represents variable assignment in Goby.
type AssignExpression struct {
	*BaseNode
	Variables []Expression
	Value     Expression
	// Optioned attribute is only used when infix expression is local assignment in params.
	// For example: `foo(x = 10)`'s `x = 10` is an optioned assign expression
	// TODO: Remove this when we can put metadata inside bytecode.
	Optioned int
}

func (ae *AssignExpression) expressionNode() {}
func (ae *AssignExpression) TokenLiteral() string {
	return ae.Token.Literal
}
func (ae *AssignExpression) String() string {
	var out bytes.Buffer
	var variables []string

	for _, v := range ae.Variables {
		variables = append(variables, v.String())
	}

	out.WriteString(strings.Join(variables, ", "))
	out.WriteString(" = ")
	out.WriteString(ae.Value.String())

	return out.String()
}

// BooleanExpression defines the boolean expression literal which contains the node expression and its value
type BooleanExpression struct {
	*BaseNode
	Value bool
}

func (b *BooleanExpression) expressionNode() {}

// BooleanExpression.TokenLiteral gets the literal of the Boolean type token
func (b *BooleanExpression) TokenLiteral() string {
	return b.Token.Literal
}

// BooleanExpression.String gets the string format of the Boolean type token
func (b *BooleanExpression) String() string {
	return b.Token.Literal
}

// NilExpression represents nil node
type NilExpression struct {
	*BaseNode
}

func (n *NilExpression) expressionNode() {}

// TokenLiteral returns `nil`
func (n *NilExpression) TokenLiteral() string {
	return n.Token.Literal
}

// String returns `nil`
func (n *NilExpression) String() string {
	return "nil"
}

type TernaryExpression struct {
	*BaseNode
	Condition   Expression
	Consequence Expression
	Alternative Expression
}

func (te *TernaryExpression) expressionNode() {}
func (te *TernaryExpression) TokenLiteral() string {
	return te.Token.Literal
}
func (te *TernaryExpression) String() string {
	var out bytes.Buffer

	out.WriteString(te.Condition.String())
	out.WriteString(" ? ")
	out.WriteString(te.Consequence.String())
	out.WriteString(" : ")
	out.WriteString(te.Alternative.String())

	return out.String()
}

type IfExpression struct {
	*BaseNode
	Conditionals []*ConditionalExpression
	Alternative  *BlockStatement
}

func (ie *IfExpression) expressionNode() {}
func (ie *IfExpression) TokenLiteral() string {
	return ie.Token.Literal
}
func (ie *IfExpression) String() string {
	var out bytes.Buffer

	for i, c := range ie.Conditionals {
		if i == 0 {
			out.WriteString("if")
			out.WriteString(" ")
		} else {
			out.WriteString("elsif")
			out.WriteString(" ")
		}

		out.WriteString(c.String())
	}

	if ie.Alternative != nil {
		out.WriteString("\n")
		out.WriteString("else\n")
		out.WriteString(ie.Alternative.String())
	}

	out.WriteString("\nend")

	return out.String()
}

// ConditionalExpression represents if or elsif expression
type ConditionalExpression struct {
	*BaseNode
	Condition   Expression
	Consequence *BlockStatement
}

func (ce *ConditionalExpression) expressionNode() {}

// TokenLiteral returns `if` or `elsif`
func (ce *ConditionalExpression) TokenLiteral() string {
	return ce.Token.Literal
}

func (ce *ConditionalExpression) String() string {
	var out bytes.Buffer

	out.WriteString(ce.Condition.String())
	out.WriteString("\n")
	out.WriteString(ce.Consequence.String())

	return out.String()
}

type CallExpression struct {
	*BaseNode
	Receiver       Expression
	Method         string
	Arguments      []Expression
	Block          *BlockStatement
	BlockArguments []*Identifier
}

func (tce *CallExpression) expressionNode() {}
func (tce *CallExpression) TokenLiteral() string {
	return tce.Token.Literal
}
func (tce *CallExpression) String() string {
	var out bytes.Buffer

	out.WriteString(tce.Receiver.String())
	out.WriteString(".")
	out.WriteString(tce.Method)

	var args = []string{}
	for _, arg := range tce.Arguments {
		args = append(args, arg.String())
	}

	out.WriteString("(")
	out.WriteString(strings.Join(args, ", "))
	out.WriteString(")")

	if tce.Block != nil {
		var blockArgs []string
		out.WriteString(" do")

		if len(tce.BlockArguments) > 0 {
			for _, arg := range tce.BlockArguments {
				blockArgs = append(blockArgs, arg.String())
			}
			out.WriteString(" |")
			out.WriteString(strings.Join(blockArgs, ", "))
			out.WriteString("|")
		}

		out.WriteString("\n")
		out.WriteString(tce.Block.String())
		out.WriteString("\nend")
	}

	return out.String()
}

type SelfExpression struct {
	*BaseNode
}

func (se *SelfExpression) expressionNode() {}
func (se *SelfExpression) TokenLiteral() string {
	return se.Token.Literal
}
func (se *SelfExpression) String() string {
	return "self"
}

type YieldExpression struct {
	*BaseNode
	Arguments []Expression
}

func (ye *YieldExpression) expressionNode() {}
func (ye *YieldExpression) TokenLiteral() string {
	return ye.Token.Literal
}
func (ye *YieldExpression) String() string {
	var out bytes.Buffer
	var args []string

	for _, arg := range ye.Arguments {
		args = append(args, arg.String())
	}

	out.WriteString(ye.TokenLiteral())
	out.WriteString("(")
	out.WriteString(strings.Join(args, ", "))
	out.WriteString(")")

	return out.String()
}

// GetBlockExpression represents `get_block` call in the AST
type GetBlockExpression struct {
	*BaseNode
}

func (gbe *GetBlockExpression) expressionNode() {}

// TokenLiteral ...
func (gbe *GetBlockExpression) TokenLiteral() string {
	return gbe.Token.Literal
}

// String ...
func (gbe *GetBlockExpression) String() string {
	return gbe.TokenLiteral()
}

// RangeExpression defines the range expression literal which contains the node expression and its start/end value
type RangeExpression struct {
	*BaseNode
	Start Expression
	End   Expression
}

func (re *RangeExpression) expressionNode() {}

// RangeExpression.TokenLiteral gets the literal of the Range type token
func (re *RangeExpression) TokenLiteral() string {
	return re.Token.Literal
}

// RangeExpression.String gets the string format of the Range type token
func (re *RangeExpression) String() string {
	var out bytes.Buffer

	out.WriteString("(")
	out.WriteString(re.Start.String())
	out.WriteString("..")
	out.WriteString(re.End.String())
	out.WriteString(")")

	return out.String()
}
