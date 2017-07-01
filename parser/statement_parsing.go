package parser

import (
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
			stmt.Name =
				&ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
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

	if p.peekTokenAtSameLine() { // `def foo()` or `def foo x `, next token at same line
		if p.peekTokenIs(token.LParen) {
			p.nextToken()

			// empty params
			if p.peekTokenIs(token.RParen) {
				p.nextToken()
				stmt.Parameters = []*ast.Identifier{}
			} else {
				stmt.Parameters = p.parseParameters()

				if !p.expectPeek(token.RParen) {
					return nil
				}
			}

		} else if p.peekTokenIs(token.Ident) { // def foo x, next token is x and at same line
			stmt.Parameters = p.parseParameters()
		}

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
		stmt.SuperClass = p.parseExpression(NORMAL)

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
	stmt.Body = p.parseBlockStatement()

	return stmt
}

func (p *Parser) parseParameters() []*ast.Identifier {

	identifiers := []*ast.Identifier{}

	p.nextToken()
	ident := &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
	identifiers = append(identifiers, ident)

	for p.peekTokenIs(token.Comma) {
		p.nextToken()
		p.nextToken()
		identifier := &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
		identifiers = append(identifiers, identifier)
	}

	return identifiers
}

func (p *Parser) parseReturnStatement() *ast.ReturnStatement {
	stmt := &ast.ReturnStatement{Token: p.curToken}

	p.nextToken()

	stmt.ReturnValue = p.parseExpression(NORMAL)

	if p.peekTokenIs(token.Semicolon) {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) parseExpressionStatement() *ast.ExpressionStatement {
	stmt := &ast.ExpressionStatement{Token: p.curToken}

	if p.curTokenIs(token.Ident) {
		// I use precedence to identify call_without_parens case, this is not an appropriate way but it work in current situation
		stmt.Expression = p.parseExpression(LOWEST)
	} else {
		stmt.Expression = p.parseExpression(NORMAL)
	}

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
			p.error = &Error{Message: "Unexpected EOF", errType: EndOfFileError}
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
	ws.Condition = p.parseExpression(NORMAL)
	p.acceptBlock = true
	p.nextToken()

	if p.curTokenIs(token.Semicolon) {
		p.nextToken()
	}

	ws.Body = p.parseBlockStatement()

	return ws
}
