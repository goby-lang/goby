package parser

import (
	"fmt"
	"github.com/goby-lang/goby/compiler/ast"
	"github.com/goby-lang/goby/compiler/token"
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
	case token.Break:
		return &ast.BreakStatement{Token: p.curToken}
	default:
		return p.parseExpressionStatement()
	}
}

func (p *Parser) parseDefMethodStatement() *ast.DefStatement {
	var params []ast.Expression
	stmt := &ast.DefStatement{Token: p.curToken}

	p.nextToken()

	// Method has specific receiver like `def self.foo` or `def bar.foo`
	if p.peekTokenIs(token.Dot) {
		switch p.curToken.Type {
		case token.Ident:
			stmt.Receiver = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
		case token.Self:
			stmt.Receiver = &ast.SelfExpression{Token: p.curToken}
		default:
			p.error = &Error{Message: fmt.Sprintf("Invalid method receiver: %s. Line: %d", p.curToken.Literal, p.curToken.Line), errType: MethodDefinitionError}
		}

		p.nextToken() // .
		if !p.expectPeek(token.Ident) {
			return nil
		}
	}

	stmt.Name = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}

	// Setter method def foo=()
	if p.peekTokenIs(token.Assign) {
		stmt.Name.Value = stmt.Name.Value + "="
		p.nextToken()
	}

	if p.peekTokenIs(token.Ident) && p.peekTokenAtSameLine() { // def foo x, next token is x and at same line
		p.error = &Error{Message: fmt.Sprintf("Please add parentheses around method \"%s\"'s parameters. Line: %d", stmt.Name.Value, p.curToken.Line), errType: MethodDefinitionError}
	}

	if p.peekTokenIs(token.LParen) {
		p.nextToken()

		// empty params
		if p.peekTokenIs(token.RParen) {
			p.nextToken()
			params = []ast.Expression{}
		} else {
			params = p.parseParameters()

			if !p.expectPeek(token.RParen) {
				return nil
			}
		}
	} else {
		params = []ast.Expression{}
	}

	stmt.Parameters = params
	stmt.BlockStatement = p.parseBlockStatement()

	return stmt
}

func (p *Parser) parseParameters() []ast.Expression {
	p.fsm.Event(parseMethodParam)
	params := []ast.Expression{}

	p.nextToken()
	param := p.parseExpression(NORMAL)
	params = append(params, param)

	for p.peekTokenIs(token.Comma) {
		p.nextToken()
		p.nextToken()
		param := p.parseExpression(NORMAL)
		params = append(params, param)
	}

	p.fsm.Event(normal)
	return params
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
