package parser

import (
	"github.com/st0012/Rooby/ast"
	"github.com/st0012/Rooby/token"
)

func (p *Parser) parseStatement() ast.Statement {
	switch p.curToken.Type {
	case token.INSTANCE_VARIABLE, token.IDENT, token.CONSTANT:
		if p.curToken.Literal == "class" {
			p.curToken.Type = token.CLASS
			return p.parseStatement()
		}

		if p.peekTokenIs(token.ASSIGN) {
			return p.parseAssignStatement()
		} else {
			return p.parseExpressionStatement()
		}
	case token.RETURN:
		return p.parseReturnStatement()
	case token.DEF:
		return p.parseDefMethodStatement()
	case token.CLASS:
		return p.parseClassStatement()
	case token.COMMENT:
		return nil
	default:
		return p.parseExpressionStatement()
	}
}

func (p *Parser) parseDefMethodStatement() *ast.DefStatement {
	stmt := &ast.DefStatement{Token: p.curToken}

	p.nextToken()

	switch p.curToken.Type {
	case token.IDENT:
		if p.peekTokenIs(token.DOT) {
			stmt.Receiver = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
			p.nextToken() // .
			if !p.expectPeek(token.IDENT) {
				return nil
			}
			stmt.Name = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
		} else {
			stmt.Name = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
		}
	case token.SELF:
		stmt.Receiver = &ast.SelfExpression{Token: p.curToken}
		p.nextToken() // .
		if !p.expectPeek(token.IDENT) {
			return nil
		}
		stmt.Name = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
	default:
		return nil
	}

	// def foo
	if p.peekTokenAtSameLine() { // def foo(), next token is ( and at same line
		if !p.expectPeek(token.LPAREN) {
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

	if !p.expectPeek(token.CONSTANT) {
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

	if p.peekTokenIs(token.RPAREN) {
		p.nextToken()
		return identifiers
	} // empty params

	p.nextToken()

	ident := &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
	identifiers = append(identifiers, ident)

	for p.peekTokenIs(token.COMMA) {
		p.nextToken()
		p.nextToken()
		identifier := &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
		identifiers = append(identifiers, identifier)
	}

	if !p.expectPeek(token.RPAREN) {
		return nil
	}

	return identifiers
}

func (p *Parser) parseAssignStatement() *ast.AssignStatement {
	stmt := &ast.AssignStatement{Token: p.curToken}

	switch p.curToken.Type {
	case token.IDENT:
		stmt.Name = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
	case token.CONSTANT:
		stmt.Name = &ast.Constant{Token: p.curToken, Value: p.curToken.Literal}
	case token.INSTANCE_VARIABLE:
		stmt.Name = &ast.InstanceVariable{Token: p.curToken, Value: p.curToken.Literal}
	}

	if !p.expectPeek(token.ASSIGN) {
		return nil
	}

	p.nextToken()

	stmt.Value = p.parseExpression(LOWEST)

	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) parseReturnStatement() *ast.ReturnStatement {
	stmt := &ast.ReturnStatement{Token: p.curToken}

	p.nextToken()

	stmt.ReturnValue = p.parseExpression(LOWEST)

	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) parseExpressionStatement() *ast.ExpressionStatement {
	stmt := &ast.ExpressionStatement{Token: p.curToken}

	stmt.Expression = p.parseExpression(LOWEST)
	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) parseBlockStatement() *ast.BlockStatement {
	// curToken is {
	bs := &ast.BlockStatement{Token: p.curToken}
	bs.Statements = []ast.Statement{}

	p.nextToken()

	for !p.curTokenIs(token.END) && !p.curTokenIs(token.ELSE) {
		stmt := p.parseStatement()
		if stmt != nil {
			bs.Statements = append(bs.Statements, stmt)
		}
		p.nextToken()
	}

	return bs
}
