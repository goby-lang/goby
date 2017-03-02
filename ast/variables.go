package ast

import "github.com/st0012/Rooby/token"

type Variable interface {
	variableNode()
	Node
}

type Identifier struct {
	Token token.Token
	Value string
}

func (i *Identifier) variableNode() {}
func (i *Identifier) expressionNode() {}
func (i *Identifier) TokenLiteral() string {
	return i.Token.Literal
}
func (i *Identifier) String() string {
	return i.Value
}

type InstanceVariable struct {
	Token token.Token
	Value string
}

func (iv *InstanceVariable) variableNode() {}
func (iv *InstanceVariable) expressionNode() {}
func (iv *InstanceVariable) TokenLiteral() string {
	return iv.Token.Literal
}
func (iv *InstanceVariable) String() string {
	return iv.Value
}

type Constant struct {
	Token token.Token
	Value string
}

func (c *Constant) variableNode() {}
func (c *Constant) expressionNode() {}
func (c *Constant) TokenLiteral() string {
	return c.Token.Literal
}
func (c *Constant) String() string {
	return c.Value
}
