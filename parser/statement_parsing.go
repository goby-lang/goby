package parser

import (
	"fmt"
	"github.com/goby-lang/goby/ast"
	"github.com/goby-lang/goby/token"
)

func (p *Parser) parseStatement() ast.Statement {
	switch p.curToken.Type {
	case token.InstanceVariable, token.Ident, token.Constant:

		return p.parseExpressionStatement()
	case token.Return:
		return p.parseReturnStatement()
	case token.Def:
		return p.parseDefMethodStatement()
	case token.Comment:
		return nil
	case token.While:
		return p.parseWhileStatement()
	case token.Class:
		return p.parseClassStatement()
	case token.Module:
		return p.parseModuleStatement()
	case token.Next:
		return &ast.NextStatement{Token: p.curToken}
	default:
		return p.parseExpressionStatement()
	}
}

func (p *Parser) parseDefMethodStatement() *ast.DefStatement {
	stmt := &ast.DefStatement{Token: p.curToken}

	p.nextToken()

	switch p.curToken.Type {
	case token.Ident:
		if p.peekTokenIs(token.Dot) {
			stmt.Receiver = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
			p.nextToken() // .
			if !p.expectPeek(token.Ident) {
				return nil
			}
			stmt.Name = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
		} else {
			stmt.Name = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
		}
	case token.Self:
		stmt.Receiver = &ast.SelfExpression{Token: p.curToken}
		p.nextToken() // .
		if !p.expectPeek(token.Ident) {
			return nil
		}
		stmt.Name = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
	default:
		return nil
	}

	// Setter method def foo=()
	if p.peekTokenIs(token.Assign) {
		stmt.Name.Value = stmt.Name.Value + "="
		p.nextToken()
	}
	// def foo
	if p.peekTokenAtSameLine() { // def foo(), next token is ( and at same line
		if !p.expectPeek(token.LParen) {
			return nil
		}

		stmt.Parameters = p.parseParameters()
	} else {
		stmt.Parameters = []*ast.Identifier{}
	}

	stmt.BlockStatement = p.parseBlockStatement()

	return stmt
}

func (p *Parser) parseClassStatement() *ast.ClassStatement {
	stmt := &ast.ClassStatement{Token: p.curToken}

	if !p.expectPeek(token.Constant) {
		return nil
	}

	stmt.Name = &ast.Constant{Token: p.curToken, Value: p.curToken.Literal}

	// See if there is any inheritance
	if p.peekTokenIs(token.LT) {
		p.nextToken() // <
		p.nextToken() // Inherited class like 'Bar'
		stmt.SuperClass = p.parseExpression(LOWEST)

		switch exp := stmt.SuperClass.(type) {
		case *ast.InfixExpression:
			stmt.SuperClassName = exp.Right.(*ast.Constant).Value
		case *ast.Constant:
			stmt.SuperClassName = exp.Value
		}
	}

	stmt.Body = p.parseBlockStatement()

	return stmt
}

func (p *Parser) parseModuleStatement() *ast.ModuleStatement {
	stmt := &ast.ModuleStatement{Token: p.curToken}

	if !p.expectPeek(token.Constant) {
		return nil
	}

	stmt.Name = &ast.Constant{Token: p.curToken, Value: p.curToken.Literal}

	// See if there is any inheritance
	if p.peekTokenIs(token.LT) {
		msg := fmt.Sprintf("Module doesn't support inheritance. Line: %d", p.curToken.Line)
		p.errors = append(p.errors, msg)
		return nil
	}

	stmt.Body = p.parseBlockStatement()

	return stmt
}

func (p *Parser) parseParameters() []*ast.Identifier {
	identifiers := []*ast.Identifier{}

	if p.peekTokenIs(token.RParen) {
		p.nextToken()
		return identifiers
	} // empty params

	p.nextToken()

	ident := &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
	identifiers = append(identifiers, ident)

	for p.peekTokenIs(token.Comma) {
		p.nextToken()
		p.nextToken()
		identifier := &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
		identifiers = append(identifiers, identifier)
	}

	if !p.expectPeek(token.RParen) {
		return nil
	}

	return identifiers
}

func (p *Parser) parseReturnStatement() *ast.ReturnStatement {
	stmt := &ast.ReturnStatement{Token: p.curToken}

	p.nextToken()

	stmt.ReturnValue = p.parseExpression(LOWEST)

	if p.peekTokenIs(token.Semicolon) {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) parseExpressionStatement() *ast.ExpressionStatement {
	stmt := &ast.ExpressionStatement{Token: p.curToken}

	stmt.Expression = p.parseExpression(LOWEST)
	if p.peekTokenIs(token.Semicolon) {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) parseBlockStatement() *ast.BlockStatement {
	// curToken is {
	bs := &ast.BlockStatement{Token: p.curToken}
	bs.Statements = []ast.Statement{}

	p.nextToken()

	for !p.curTokenIs(token.End) && !p.curTokenIs(token.Else) {

		if p.curTokenIs(token.EOF) {
			p.errors = append(p.errors, syntaxError("end", "EOF"))
			return bs
		}

		stmt := p.parseStatement()
		if stmt != nil {
			bs.Statements = append(bs.Statements, stmt)
		}
		p.nextToken()
	}

	return bs
}

func (p *Parser) parseWhileStatement() *ast.WhileStatement {
	ws := &ast.WhileStatement{Token: p.curToken}

	p.nextToken()
	// Prevent expression's method call to consume while's block as argument.
	p.acceptBlock = false
	ws.Condition = p.parseExpression(LOWEST)
	p.acceptBlock = true
	p.nextToken()

	if p.curTokenIs(token.Semicolon) {
		p.nextToken()
	}

	ws.Body = p.parseBlockStatement()

	return ws
}

func syntaxError(expecting string, unexpected string) string {
	return "Syntax error: Expecting '" + expecting + "', unexpected '" + unexpected + "'"
}
