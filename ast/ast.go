package ast

import (
	"bytes"
	"fmt"
	"github.com/goby-lang/goby/token"
	"strings"
)

type node interface {
	TokenLiteral() string
	String() string
}

type Statement interface {
	node
	statementNode()
}

type Expression interface {
	node
	expressionNode()
}

// Program is the root node of entire AST
type Program struct {
	Statements []Statement
}

func (p *Program) TokenLiteral() string {
	if len(p.Statements) > 0 {
		return p.Statements[0].TokenLiteral()
	}

	return ""
}

func (p *Program) String() string {
	var out bytes.Buffer

	for _, s := range p.Statements {
		out.WriteString(s.String())
	}

	return out.String()
}

type DefStatement struct {
	Token          token.Token
	Name           *Identifier
	Receiver       Expression
	Parameters     []*Identifier
	BlockStatement *BlockStatement
}

func (ds *DefStatement) statementNode() {}
func (ds *DefStatement) TokenLiteral() string {
	return ds.Token.Literal
}
func (ds *DefStatement) String() string {
	var out bytes.Buffer

	out.WriteString("def ")
	out.WriteString(ds.Name.TokenLiteral())
	out.WriteString("(")

	for i, param := range ds.Parameters {
		out.WriteString(param.String())
		if i != len(ds.Parameters)-1 {
			out.WriteString(", ")
		}
	}

	out.WriteString(") ")
	out.WriteString("{\n")
	out.WriteString(ds.BlockStatement.String())
	out.WriteString("\n}")

	return out.String()
}

// NextStatement represents "next" keyword
type NextStatement struct {
	Token token.Token
}

func (ns *NextStatement) statementNode() {}

// TokenLiteral returns token's literal
func (ns *NextStatement) TokenLiteral() string {
	return ns.Token.Literal
}
func (ns *NextStatement) String() string {
	return "next"
}

type ClassStatement struct {
	Token          token.Token
	Name           *Constant
	Body           *BlockStatement
	SuperClass     Expression
	SuperClassName string
}

func (cs *ClassStatement) statementNode() {}
func (cs *ClassStatement) TokenLiteral() string {
	return cs.Token.Literal
}
func (cs *ClassStatement) String() string {
	var out bytes.Buffer

	out.WriteString("class ")
	out.WriteString(cs.Name.TokenLiteral())
	out.WriteString(" {\n")
	out.WriteString(cs.Body.String())
	out.WriteString("\n}")

	return out.String()
}

// ModuleStatement represents module node in AST
type ModuleStatement struct {
	Token      token.Token
	Name       *Constant
	Body       *BlockStatement
	SuperClass *Constant
}

func (ms *ModuleStatement) statementNode() {}

// TokenLiteral returns token's literal
func (ms *ModuleStatement) TokenLiteral() string {
	return ms.Token.Literal
}
func (ms *ModuleStatement) String() string {
	var out bytes.Buffer

	out.WriteString("module ")
	out.WriteString(ms.Name.TokenLiteral())
	out.WriteString(" {\n")
	out.WriteString(ms.Body.String())
	out.WriteString("\n}")

	return out.String()
}

type ReturnStatement struct {
	Token       token.Token
	ReturnValue Expression
}

func (rs *ReturnStatement) statementNode() {}
func (rs *ReturnStatement) TokenLiteral() string {
	return rs.Token.Literal
}
func (rs *ReturnStatement) String() string {
	var out bytes.Buffer

	out.WriteString(rs.TokenLiteral() + " ")

	if rs.ReturnValue != nil {
		out.WriteString(rs.ReturnValue.String())
	}

	out.WriteString(";")

	return out.String()
}

type ExpressionStatement struct {
	Token      token.Token
	Expression Expression
}

func (es *ExpressionStatement) statementNode() {}
func (es *ExpressionStatement) TokenLiteral() string {
	return es.Token.Literal
}
func (es *ExpressionStatement) String() string {
	if es.Expression != nil {
		return es.Expression.String()
	}

	return ""
}

type IntegerLiteral struct {
	Token token.Token
	Value int
}

func (il *IntegerLiteral) expressionNode() {}
func (il *IntegerLiteral) TokenLiteral() string {
	return il.Token.Literal
}
func (il *IntegerLiteral) String() string {
	return il.Token.Literal
}

type StringLiteral struct {
	Token token.Token
	Value string
}

func (sl *StringLiteral) expressionNode() {}
func (sl *StringLiteral) TokenLiteral() string {
	return sl.Token.Literal
}
func (sl *StringLiteral) String() string {
	var out bytes.Buffer

	out.WriteString("\"")
	out.WriteString(sl.Token.Literal)
	out.WriteString("\"")
	return out.String()
}

type ArrayExpression struct {
	Token    token.Token
	Elements []Expression
}

func (ae *ArrayExpression) expressionNode() {}
func (ae *ArrayExpression) TokenLiteral() string {
	return ae.Token.Literal
}
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

type HashExpression struct {
	Token token.Token
	Data  map[string]Expression
}

func (he *HashExpression) expressionNode() {}
func (he *HashExpression) TokenLiteral() string {
	return he.Token.Literal
}
func (he *HashExpression) String() string {
	var out bytes.Buffer
	var pairs []string

	for key, value := range he.Data {
		pairs = append(pairs, fmt.Sprintf("%s: %s", key, value.String()))
	}

	out.WriteString("{ ")
	out.WriteString(strings.Join(pairs, ", "))
	out.WriteString(" }")

	return out.String()
}

type PrefixExpression struct {
	Token    token.Token
	Operator string
	Right    Expression
}

func (pe *PrefixExpression) expressionNode() {}
func (pe *PrefixExpression) TokenLiteral() string {
	return pe.Token.Literal
}
func (pe *PrefixExpression) String() string {
	var out bytes.Buffer

	out.WriteString("(")
	out.WriteString(pe.Operator)
	out.WriteString(pe.Right.String())
	out.WriteString(")")

	return out.String()
}

type InfixExpression struct {
	Token    token.Token
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

type BooleanExpression struct {
	Token token.Token
	Value bool
}

func (b *BooleanExpression) expressionNode() {}
func (b *BooleanExpression) TokenLiteral() string {
	return b.Token.Literal
}
func (b *BooleanExpression) String() string {
	return b.Token.Literal
}

// NilExpression represents nil node
type NilExpression struct {
	Token token.Token
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

type IfExpression struct {
	Token       token.Token
	Condition   Expression
	Consequence *BlockStatement
	Alternative *BlockStatement
}

func (ie *IfExpression) expressionNode() {}
func (ie *IfExpression) TokenLiteral() string {
	return ie.Token.Literal
}
func (ie *IfExpression) String() string {
	var out bytes.Buffer

	out.WriteString("if")
	out.WriteString(" ")
	out.WriteString(ie.Condition.String())
	out.WriteString("\n")
	out.WriteString(ie.Consequence.String())

	if ie.Alternative != nil {
		out.WriteString("\n")
		out.WriteString("else\n")
		out.WriteString(ie.Alternative.String())
	}

	out.WriteString("\nend")

	return out.String()
}

type BlockStatement struct {
	Token      token.Token // {
	Statements []Statement
}

func (bs *BlockStatement) statementNode() {}
func (bs *BlockStatement) TokenLiteral() string {
	return bs.Token.Literal
}
func (bs *BlockStatement) String() string {
	var out bytes.Buffer

	for _, stmt := range bs.Statements {
		out.WriteString(stmt.String())
	}

	return out.String()
}

type CallExpression struct {
	Receiver       Expression
	Token          token.Token
	Method         string
	Arguments      []Expression
	Block          *BlockStatement
	BlockArguments []*Identifier
}

func (ce *CallExpression) expressionNode() {}
func (ce *CallExpression) TokenLiteral() string {
	return ce.Token.Literal
}
func (ce *CallExpression) String() string {
	var out bytes.Buffer

	out.WriteString(ce.Receiver.String())
	out.WriteString(".")
	out.WriteString(ce.Method)

	var args = []string{}
	for _, arg := range ce.Arguments {
		args = append(args, arg.String())
	}

	out.WriteString("(")
	out.WriteString(strings.Join(args, ", "))
	out.WriteString(")")

	if ce.Block != nil {
		var blockArgs []string
		out.WriteString(" do")

		if len(ce.BlockArguments) > 0 {
			for _, arg := range ce.BlockArguments {
				blockArgs = append(blockArgs, arg.String())
			}
			out.WriteString(" |")
			out.WriteString(strings.Join(blockArgs, ", "))
			out.WriteString("|")
		}

		out.WriteString("\n")
		out.WriteString(ce.Block.String())
		out.WriteString("\nend")
	}

	return out.String()
}

type SelfExpression struct {
	Token token.Token
}

func (se *SelfExpression) expressionNode() {}
func (se *SelfExpression) TokenLiteral() string {
	return se.Token.Literal
}
func (se *SelfExpression) String() string {
	return "self"
}

type WhileStatement struct {
	Token     token.Token
	Condition Expression
	Body      *BlockStatement
}

func (ws *WhileStatement) statementNode() {}
func (ws *WhileStatement) TokenLiteral() string {
	return ws.Token.Literal
}
func (ws *WhileStatement) String() string {
	var out bytes.Buffer

	out.WriteString("while ")
	out.WriteString(ws.Condition.String())
	out.WriteString(" do\n")
	out.WriteString(ws.Body.String())
	out.WriteString("\nend")

	return out.String()
}

type YieldExpression struct {
	Token     token.Token
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
