package parser

import (
	"fmt"
	"strconv"

	"github.com/rooby-lang/rooby/ast"
	"github.com/rooby-lang/rooby/token"
)

var precedence = map[token.Type]int{
	token.Eq:       EQUALS,
	token.NotEq:    EQUALS,
	token.LT:       LESSGREATER,
	token.LTE:      LESSGREATER,
	token.GT:       LESSGREATER,
	token.GTE:      LESSGREATER,
	token.COMP:     LESSGREATER,
	token.And:      LESSGREATER,
	token.Or:       LESSGREATER,
	token.Plus:     SUM,
	token.Minus:    SUM,
	token.Incr:     SUM,
	token.Decr:     SUM,
	token.Slash:    PRODUCT,
	token.Asterisk: PRODUCT,
	token.Pow:      PRODUCT,
	token.LBracket: INDEX,
	token.Dot:      CALL,
	token.LParen:   CALL,
}

// Constants for denoting precedence
const (
	_ int = iota
	LOWEST
	EQUALS
	LESSGREATER
	SUM
	PRODUCT
	PREFIX
	INDEX
	CALL
)

type (
	prefixParseFn func() ast.Expression
	infixParseFn  func(ast.Expression) ast.Expression
)

func (p *Parser) parseExpression(precedence int) ast.Expression {
	prefix := p.prefixParseFns[p.curToken.Type]
	if prefix == nil {
		p.noPrefixParseFnError(p.curToken.Type)
		return nil
	}
	leftExp := prefix()

	for !p.peekTokenIs(token.Semicolon) && precedence < p.peekPrecedence() && p.peekTokenAtSameLine() {
		infix := p.infixParseFns[p.peekToken.Type]
		if infix == nil {
			return leftExp
		}
		p.nextToken()
		leftExp = infix(leftExp)
	}

	return leftExp
}

func (p *Parser) parseSelfExpression() ast.Expression {
	return &ast.SelfExpression{Token: p.curToken}
}

func (p *Parser) parseIdentifier() ast.Expression {
	return &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
}

func (p *Parser) parseConstant() ast.Expression {
	return &ast.Constant{Token: p.curToken, Value: p.curToken.Literal}
}

func (p *Parser) parseInstanceVariable() ast.Expression {
	return &ast.InstanceVariable{Token: p.curToken, Value: p.curToken.Literal}
}

func (p *Parser) parseIntegerLiteral() ast.Expression {
	lit := &ast.IntegerLiteral{Token: p.curToken}

	value, err := strconv.ParseInt(lit.TokenLiteral(), 0, 64)
	if err != nil {
		msg := fmt.Sprintf("could not parse %q as integer", lit.TokenLiteral())
		p.errors = append(p.errors, msg)
		return nil
	}

	lit.Value = int(value)

	return lit
}

func (p *Parser) parseStringLiteral() ast.Expression {
	lit := &ast.StringLiteral{Token: p.curToken}
	lit.Value = p.curToken.Literal

	return lit
}

func (p *Parser) parseBooleanLiteral() ast.Expression {
	lit := &ast.Boolean{Token: p.curToken}

	value, err := strconv.ParseBool(lit.TokenLiteral())
	if err != nil {
		msg := fmt.Sprintf("could not parse %q as boolean", lit.TokenLiteral())
		p.errors = append(p.errors, msg)
		return nil
	}

	lit.Value = value

	return lit
}

func (p *Parser) parsePostfixExpression(receiver ast.Expression) ast.Expression {
	arguments := []ast.Expression{}
	return &ast.CallExpression{Token: p.curToken, Receiver: receiver, Method: p.curToken.Literal, Arguments: arguments}
}

func (p *Parser) parseHashExpression() ast.Expression {
	hash := &ast.HashExpression{Token: p.curToken}
	hash.Data = p.parseHashPairs()
	return hash
}

func (p *Parser) parseHashPairs() map[string]ast.Expression {
	pairs := map[string]ast.Expression{}

	if p.peekTokenIs(token.RBrace) {
		p.nextToken() // '}'
		return pairs
	}

	p.parseHashPair(pairs)

	for p.peekTokenIs(token.Comma) {
		p.nextToken()

		p.parseHashPair(pairs)
	}

	if !p.expectPeek(token.RBrace) {
		return nil
	}

	return pairs
}

func (p *Parser) parseHashPair(pairs map[string]ast.Expression) {
	var key string
	var value ast.Expression

	if !p.expectPeek(token.Ident) {
		return
	}

	key = p.parseIdentifier().(*ast.Identifier).Value

	if !p.expectPeek(token.Colon) {
		return
	}

	p.nextToken()
	value = p.parseExpression(LOWEST)
	pairs[key] = value
}

func (p *Parser) parseArrayExpression() ast.Expression {
	arr := &ast.ArrayExpression{Token: p.curToken}
	arr.Elements = p.parseArrayElements()
	return arr
}

func (p *Parser) parseArrayIndexExpression(left ast.Expression) ast.Expression {
	callExpression := &ast.CallExpression{Receiver: left, Method: "[]", Token: p.curToken}

	if p.peekTokenIs(token.RBracket) {
		return nil
	}

	p.nextToken()

	callExpression.Arguments = []ast.Expression{p.parseExpression(LOWEST)}

	if !p.expectPeek(token.RBracket) {
		return nil
	}

	// Assign value to index
	if p.peekTokenIs(token.Assign) {
		p.nextToken()
		p.nextToken()
		assignValue := p.parseExpression(LOWEST)
		callExpression.Method = "[]="
		callExpression.Arguments = append(callExpression.Arguments, assignValue)
	}

	return callExpression
}

func (p *Parser) parseArrayElements() []ast.Expression {
	elems := []ast.Expression{}

	if p.peekTokenIs(token.RBracket) {
		p.nextToken() // ']'
		return elems
	}

	p.nextToken() // start of first expression
	elems = append(elems, p.parseExpression(LOWEST))

	for p.peekTokenIs(token.Comma) {
		p.nextToken() // ","
		p.nextToken() // start of next expression
		elems = append(elems, p.parseExpression(LOWEST))
	}

	if !p.expectPeek(token.RBracket) {
		return nil
	}

	return elems
}

func (p *Parser) parsePrefixExpression() ast.Expression {
	pe := &ast.PrefixExpression{
		Token:    p.curToken,
		Operator: p.curToken.Literal,
	}

	p.nextToken()

	pe.Right = p.parseExpression(PREFIX)

	return pe
}

func (p *Parser) parseInfixExpression(left ast.Expression) ast.Expression {
	exp := &ast.InfixExpression{
		Token:    p.curToken,
		Left:     left,
		Operator: p.curToken.Literal,
	}

	precedence := p.curPrecedence()
	p.nextToken()
	exp.Right = p.parseExpression(precedence)

	return exp
}

func (p *Parser) parseGroupedExpression() ast.Expression {
	p.nextToken()

	exp := p.parseExpression(LOWEST)

	if !p.expectPeek(token.RParen) {
		return nil
	}

	return exp
}

func (p *Parser) parseIfExpression() ast.Expression {
	ie := &ast.IfExpression{Token: p.curToken}
	p.nextToken()
	ie.Condition = p.parseExpression(LOWEST)
	ie.Consequence = p.parseBlockStatement()

	// curToken is now ELSE or RBRACE
	if p.curTokenIs(token.Else) {
		ie.Alternative = p.parseBlockStatement()
	}

	return ie
}

func (p *Parser) parseCallExpression(receiver ast.Expression) ast.Expression {
	var exp *ast.CallExpression

	if p.curTokenIs(token.LParen) { // call expression doesn't have a receiver foo(x) || foo()
		// method name is receiver, for example 'foo' of foo(x)
		m := receiver.(*ast.Identifier).Value
		// receiver is self
		selfTok := token.Token{Type: token.Self, Literal: "self", Line: p.curToken.Line}
		self := &ast.SelfExpression{Token: selfTok}
		receiver = self

		// current token is identifier (method name)
		exp = &ast.CallExpression{Token: p.curToken, Receiver: receiver, Method: m}
		exp.Arguments = p.parseCallArguments()
	} else { // call expression has a receiver like: p.foo
		exp = &ast.CallExpression{Token: p.curToken, Receiver: receiver}

		// check if method name is identifier
		if !p.expectPeek(token.Ident) {
			return nil
		}

		exp.Method = p.curToken.Literal

		if p.peekTokenIs(token.LParen) {
			p.nextToken()
			exp.Arguments = p.parseCallArguments()
		} else { // p.foo.bar; || p.foo; || p.foo + 123
			exp.Arguments = []ast.Expression{}
		}
	}

	// Setter method call like: p.foo = x
	if p.peekTokenIs(token.Assign) {
		exp.Method = exp.Method + "="
		p.nextToken()
		p.nextToken()
		exp.Arguments = append(exp.Arguments, p.parseExpression(LOWEST))
	}

	// Parse block
	if p.peekTokenIs(token.Do) {
		p.parseBlockParameters(exp)
	}

	return exp
}

func (p *Parser) parseBlockParameters(exp *ast.CallExpression) {
	p.nextToken()

	// Parse block arguments
	if p.peekTokenIs(token.Bar) {
		var params []*ast.Identifier

		p.nextToken()
		p.nextToken()

		param := &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
		params = append(params, param)

		for p.peekTokenIs(token.Comma) {
			p.nextToken()
			p.nextToken()
			param := &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
			params = append(params, param)
		}

		if !p.expectPeek(token.Bar) {
			return
		}

		exp.BlockArguments = params
	}

	exp.Block = p.parseBlockStatement()
}

func (p *Parser) parseCallArguments() []ast.Expression {
	args := []ast.Expression{}

	if p.peekTokenIs(token.RParen) {
		p.nextToken() // ')'
		return args
	}

	p.nextToken() // start of first expression
	args = append(args, p.parseExpression(LOWEST))

	for p.peekTokenIs(token.Comma) {
		p.nextToken() // ","
		p.nextToken() // start of next expression
		args = append(args, p.parseExpression(LOWEST))
	}

	if !p.expectPeek(token.RParen) {
		return nil
	}

	return args
}

func (p *Parser) parseYieldExpression() ast.Expression {
	ye := &ast.YieldExpression{Token: p.curToken}

	if p.peekTokenIs(token.LParen) {
		p.nextToken()
		ye.Arguments = p.parseCallArguments()
	}

	return ye
}
