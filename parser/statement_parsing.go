package parser

import (
	"github.com/rooby-lang/rooby/ast"
	"github.com/rooby-lang/rooby/token"
)

func (p *Parser) parseStatement() ast.Statement {
	switch p.curToken.Type {
	case token.InstanceVariable, token.Ident, token.Constant:
		if p.curToken.Literal == "class" {
			p.curToken.Type = token.Class
			return p.parseStatement()
		}

		if p.peekTokenIs(token.Assign) {
			return p.parseAssignStatement()
		}

		return p.parseExpressionStatement()

	case token.Return:
		return p.parseReturnStatement()
	case token.Def:
		return p.parseDefMethodStatement()
	case token.Class:
		return p.parseClassStatement()
	case token.Comment:
		return nil
	case token.While:
		return p.parseWhileStatement()
	case token.RequireRelative:
		return p.parseRequireRelativeStatement()
	default:
		return p.parseExpressionStatement()
	}
}


func (p *Parser) parseRequireRelativeStatement() *ast.RequireRelativeStatement {
	stmt := &ast.RequireRelativeStatement{Token: p.curToken}
	p.nextToken()

	filepath := p.curToken.Literal
	stmt.Filepath = filepath
	return stmt
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
		stmt.SuperClass = &ast.Constant{Token: p.curToken, Value: p.curToken.Literal}
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

func (p *Parser) parseAssignStatement() *ast.AssignStatement {
	stmt := &ast.AssignStatement{Token: p.curToken}

	switch p.curToken.Type {
	case token.Ident:
		stmt.Name = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
	case token.Constant:
		stmt.Name = &ast.Constant{Token: p.curToken, Value: p.curToken.Literal}
	case token.InstanceVariable:
		stmt.Name = &ast.InstanceVariable{Token: p.curToken, Value: p.curToken.Literal}
	}

	if !p.expectPeek(token.Assign) {
		return nil
	}

	p.nextToken()

	stmt.Value = p.parseExpression(LOWEST)

	if p.peekTokenIs(token.Semicolon) {
		p.nextToken()
	}

	return stmt
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
	ws.Condition = p.parseExpression(LOWEST)
	ws.Body = p.parseBlockStatement()

	return ws
}
