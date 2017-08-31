package parser

import (
	"fmt"
	"github.com/goby-lang/goby/compiler/ast"
	"github.com/goby-lang/goby/compiler/token"
)

func (p *Parser) parseStatement() ast.Statement {
	switch p.curToken.Type {
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
		return &ast.NextStatement{BaseNode: &ast.BaseNode{Token: p.curToken}}
	case token.Break:
		return &ast.BreakStatement{BaseNode: &ast.BaseNode{Token: p.curToken}}
	default:
		exp := p.parseExpressionStatement()

		// If parseExpressionStatement got error exp.Expression would be nil
		if exp.Expression != nil {
			// In REPL mode everything should return a value.
			if p.Mode == REPLMode {
				exp.Expression.MarkAsExp()
			} else {
				exp.Expression.MarkAsStmt()
			}
		}

		return exp
	}
}

func (p *Parser) parseDefMethodStatement() *ast.DefStatement {
	var params []ast.Expression
	stmt := &ast.DefStatement{BaseNode: &ast.BaseNode{Token: p.curToken}}

	p.nextToken()

	// Method has specific receiver like `def self.foo` or `def bar.foo`
	if p.peekTokenIs(token.Dot) {
		switch p.curToken.Type {
		case token.Ident:
			stmt.Receiver = p.parseIdentifier()
		case token.InstanceVariable:
			stmt.Receiver = p.parseInstanceVariable()
		case token.Constant:
			stmt.Receiver = p.parseConstant()
		case token.Self:
			stmt.Receiver = &ast.SelfExpression{BaseNode: &ast.BaseNode{Token: p.curToken}}
		default:
			p.error = &Error{Message: fmt.Sprintf("Invalid method receiver: %s. Line: %d", p.curToken.Literal, p.curToken.Line), errType: MethodDefinitionError}
		}

		p.nextToken() // .
		if !p.expectPeek(token.Ident) {
			return nil
		}
	}

	stmt.Name = &ast.Identifier{BaseNode: &ast.BaseNode{Token: p.curToken}, Value: p.curToken.Literal}

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

			if p.peekTokenIs(token.Asterisk) {
				p.nextToken()
				if p.expectPeek(token.Ident) {
					stmt.SplatParameter = p.parseIdentifier().(*ast.Identifier)
				}
			}

			if !p.expectPeek(token.RParen) {
				return nil
			}
		}
	} else {
		params = []ast.Expression{}
	}

	stmt.Parameters = params
	stmt.BlockStatement = p.parseBlockStatement()
	stmt.BlockStatement.KeepLastValue()

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

		if p.peekTokenIs(token.Asterisk) {
			break
		}

		p.nextToken()

		param := p.parseExpression(NORMAL)

		if paramDuplicated(params, param) {
			p.error = &Error{Message: fmt.Sprintf("Duplicate argument name: \"%s\". Line: %d", getArgName(param), p.curToken.Line), errType: SyntaxError}
		}
		params = append(params, param)
	}

	p.fsm.Event(backToNormal)
	return params
}

func (p *Parser) parseClassStatement() *ast.ClassStatement {
	stmt := &ast.ClassStatement{BaseNode: &ast.BaseNode{Token: p.curToken}}

	if !p.expectPeek(token.Constant) {
		return nil
	}

	stmt.Name = &ast.Constant{BaseNode: &ast.BaseNode{Token: p.curToken}, Value: p.curToken.Literal}

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
	stmt := &ast.ModuleStatement{BaseNode: &ast.BaseNode{Token: p.curToken}}

	if !p.expectPeek(token.Constant) {
		return nil
	}

	stmt.Name = &ast.Constant{BaseNode: &ast.BaseNode{Token: p.curToken}, Value: p.curToken.Literal}
	stmt.Body = p.parseBlockStatement()

	return stmt
}

func (p *Parser) parseReturnStatement() *ast.ReturnStatement {
	stmt := &ast.ReturnStatement{BaseNode: &ast.BaseNode{Token: p.curToken}}

	p.nextToken()

	stmt.ReturnValue = p.parseExpression(NORMAL)

	if p.peekTokenIs(token.Semicolon) {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) parseExpressionStatement() *ast.ExpressionStatement {
	stmt := &ast.ExpressionStatement{BaseNode: &ast.BaseNode{Token: p.curToken}}
	if p.curTokenIs(token.Ident) || p.curTokenIs(token.InstanceVariable) {
		// This is used for identifying method call without parens
		// Or multiple variable assignment
		stmt.Expression = p.parseExpression(LOWEST)
	} else {
		stmt.Expression = p.parseExpression(NORMAL)
	}

	return stmt
}

func (p *Parser) parseBlockStatement() *ast.BlockStatement {

	// curToken is '{'
	bs := &ast.BlockStatement{BaseNode: &ast.BaseNode{Token: p.curToken}}
	bs.Statements = []ast.Statement{}

	p.nextToken()

	if p.curTokenIs(token.Semicolon) {
		p.nextToken()
	}

	for !p.curTokenIs(token.End) && !p.curTokenIs(token.Else) && !p.curTokenIs(token.ElsIf) {

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
	ws := &ast.WhileStatement{BaseNode: &ast.BaseNode{Token: p.curToken}}

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

func paramDuplicated(params []ast.Expression, param ast.Expression) bool {
	for _, p := range params {
		if getArgName(param) == getArgName(p) {
			return true
		}
	}
	return false
}

func getArgName(exp ast.Expression) string {
	assignExp, ok := exp.(*ast.AssignExpression)

	if ok {
		return assignExp.Variables[0].TokenLiteral()
	}

	return exp.TokenLiteral()
}
