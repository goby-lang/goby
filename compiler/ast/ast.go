package ast

import (
	"bytes"
	"github.com/goby-lang/goby/compiler/token"
)

type BaseNode struct {
	Token token.Token
}

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
