package ast

import (
	"bytes"
)

// ClassStatement represents a class node in AST
type ClassStatement struct {
	*BaseNode
	Name           *Constant
	Body           *BlockStatement
	SuperClass     Expression
	SuperClassName string
}

func (cs *ClassStatement) statementNode() {}

// TokenLiteral is a polymorphic function to return a token literal
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

// ModuleStatement represents a module node in AST
type ModuleStatement struct {
	*BaseNode
	Name       *Constant
	Body       *BlockStatement
	SuperClass *Constant
}

func (ms *ModuleStatement) statementNode() {}

// TokenLiteral is a polymorphic function to return a token literal
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

// ReturnStatement represents an expression as a return value
type ReturnStatement struct {
	*BaseNode
	ReturnValue Expression
}

func (rs *ReturnStatement) statementNode() {}

// TokenLiteral is a polymorphic function to return a token literal
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

// ExpressionStatement represents an expression statement
type ExpressionStatement struct {
	*BaseNode
	Expression Expression
}

func (es *ExpressionStatement) statementNode() {}

// TokenLiteral is a polymorphic function to return a token literal
func (es *ExpressionStatement) TokenLiteral() string {
	return es.Token.Literal
}
func (es *ExpressionStatement) String() string {
	if es.Expression != nil {
		return es.Expression.String()
	}

	return ""
}

// DefStatement represents a "def" keyword with  block
type DefStatement struct {
	*BaseNode
	Name           *Identifier
	Receiver       Expression
	Parameters     []Expression
	BlockStatement *BlockStatement
}

func (tds *DefStatement) statementNode() {}

// TokenLiteral is a polymorphic function to return a token literal
func (tds *DefStatement) TokenLiteral() string {
	return tds.Token.Literal
}
func (tds *DefStatement) String() string {
	var out bytes.Buffer

	out.WriteString("def ")
	out.WriteString(tds.Name.TokenLiteral())
	out.WriteString("(")

	for i, param := range tds.Parameters {
		out.WriteString(param.String())
		if i != len(tds.Parameters)-1 {
			out.WriteString(", ")
		}
	}

	out.WriteString(") ")
	out.WriteString("{\n")
	out.WriteString(tds.BlockStatement.String())
	out.WriteString("\n}")

	return out.String()
}

// NextStatement represents "next" keyword
type NextStatement struct {
	*BaseNode
}

func (ns *NextStatement) statementNode() {}

// TokenLiteral is a polymorphic function to return a token literal
func (ns *NextStatement) TokenLiteral() string {
	return ns.Token.Literal
}
func (ns *NextStatement) String() string {
	return "next"
}

// BreakStatement represents "break" keyword
type BreakStatement struct {
	*BaseNode
}

func (bs *BreakStatement) statementNode() {}

// TokenLiteral is a polymorphic function to return a token literal
func (bs *BreakStatement) TokenLiteral() string {
	return bs.Token.Literal
}
func (bs *BreakStatement) String() string {
	return bs.TokenLiteral()
}

// WhileStatement represents a "while" keyword with a block
type WhileStatement struct {
	*BaseNode
	Condition Expression
	Body      *BlockStatement
}

func (ws *WhileStatement) statementNode() {}

// TokenLiteral is a polymorphic function to return a token literal
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

// BlockStatement represents a block statement
type BlockStatement struct {
	*BaseNode
	Statements []Statement
}

func (bs *BlockStatement) statementNode() {}

// TokenLiteral is a polymorphic function to return a token literal
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

// IsEmpty checks if the statements in the block is empty
func (bs *BlockStatement) IsEmpty() bool {
	return len(bs.Statements) == 0
}

// KeepLastValue prevents block's last expression statement to be popped.
func (bs *BlockStatement) KeepLastValue() {
	if len(bs.Statements) == 0 {
		return
	}

	stmt := bs.Statements[len(bs.Statements)-1]
	expStmt, ok := stmt.(*ExpressionStatement)

	if ok {
		expStmt.Expression.MarkAsExp()
	}
}
