package ast

import (
	"bytes"
	"github.com/goby-lang/goby/compiler/token"
)

// BaseNode holds the attribute every expression or statement should have
type BaseNode struct {
	Token  token.Token
	isStmt bool
}

// IsExp returns if current node should be considered as an expression
func (b *BaseNode) IsExp() bool {
	return !b.isStmt
}

// IsStmt returns if current node should be considered as a statement
func (b *BaseNode) IsStmt() bool {
	return b.isStmt
}

// MarkAsStmt marks current node to be statement
func (b *BaseNode) MarkAsStmt() {
	b.isStmt = true
}

// MarkAsExp marks current node to be expression
func (b *BaseNode) MarkAsExp() {
	b.isStmt = false
}

type node interface {
	TokenLiteral() string
	String() string
	IsExp() bool
	IsStmt() bool
	MarkAsStmt()
	MarkAsExp()
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
