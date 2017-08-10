package ast

import (
	"bytes"
	"strings"
)

// Variable interface represents assignable nodes in Goby, currently are Identifier, InstanceVariable and Constant
type Variable interface {
	variableNode()
	ReturnValue() string
	Expression
}

// MultiVariableExpression is not really an expression, it's just a container that holds multiple Variables
type MultiVariableExpression struct {
	*BaseNode
	Variables []Expression
}

func (m *MultiVariableExpression) expressionNode() {}
func (m *MultiVariableExpression) TokenLiteral() string {
	return ""
}
func (m *MultiVariableExpression) String() string {
	var out bytes.Buffer
	var variables []string

	for _, v := range m.Variables {
		variables = append(variables, v.String())
	}

	out.WriteString(strings.Join(variables, ", "))

	return out.String()
}

type Identifier struct {
	*BaseNode
	Value string
}

func (i *Identifier) variableNode() {}
func (i *Identifier) ReturnValue() string {
	return i.Value
}
func (i *Identifier) expressionNode() {}
func (i *Identifier) TokenLiteral() string {
	return i.Token.Literal
}
func (i *Identifier) String() string {
	return i.Value
}

type InstanceVariable struct {
	*BaseNode
	Value string
}

func (iv *InstanceVariable) variableNode() {}
func (iv *InstanceVariable) ReturnValue() string {
	return iv.Value
}
func (iv *InstanceVariable) expressionNode() {}
func (iv *InstanceVariable) TokenLiteral() string {
	return iv.Token.Literal
}
func (iv *InstanceVariable) String() string {
	return iv.Value
}

type Constant struct {
	*BaseNode
	Value       string
	IsNamespace bool
}

func (c *Constant) variableNode() {}
func (c *Constant) ReturnValue() string {
	return c.Value
}
func (c *Constant) expressionNode() {}
func (c *Constant) TokenLiteral() string {
	return c.Token.Literal
}
func (c *Constant) String() string {
	return c.Value
}
