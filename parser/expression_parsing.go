package parser

import (
	"fmt"
	"github.com/st0012/Rooby/ast"
	"github.com/st0012/Rooby/token"
	"strconv"
)

var precedence = map[token.TokenType]int{
	token.EQ:       EQUALS,
	token.NOT_EQ:   EQUALS,
	token.LT:       LESSGREATER,
	token.GT:       LESSGREATER,
	token.PLUS:     SUM,
	token.MINUS:    SUM,
	token.INCR:     SUM,
	token.DECR:     SUM,
	token.SLASH:    PRODUCT,
	token.ASTERISK: PRODUCT,
	token.LBRACKET: INDEX,
	token.DOT:      CALL,
	token.LPAREN:   CALL,
}

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

func (p *Parser) parseExpression(precendence int) ast.Expression {
	prefix := p.prefixParseFns[p.curToken.Type]
	if prefix == nil {
		p.noPrefixParseFnError(p.curToken.Type)
		return nil
	}
	leftExp := prefix()

	for !p.peekTokenIs(token.SEMICOLON) && precendence < p.peekPrecedence() && p.peekTokenAtSameLine() {
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

	if p.peekTokenIs(token.RBRACE) {
		p.nextToken() // '}'
		return pairs
	}

	p.parseHashPair(pairs)

	for p.peekTokenIs(token.COMMA) {
		p.nextToken()

		p.parseHashPair(pairs)
	}

	if !p.expectPeek(token.RBRACE) {
		return nil
	}

	return pairs
}

func (p *Parser) parseHashPair(pairs map[string]ast.Expression) {
	var key string
	var value ast.Expression

	if !p.expectPeek(token.IDENT) {
		return
	}

	key = p.parseIdentifier().(*ast.Identifier).Value

	if !p.expectPeek(token.COLON) {
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

	if p.peekTokenIs(token.RBRACKET) {
		return nil
	}

	p.nextToken()

	callExpression.Arguments = []ast.Expression{p.parseExpression(LOWEST)}

	if !p.expectPeek(token.RBRACKET) {
		return nil
	}

	// Assign value to index
	if p.peekTokenIs(token.ASSIGN) {
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

	if p.peekTokenIs(token.RBRACKET) {
		p.nextToken() // ']'
		return elems
	}

	p.nextToken() // start of first expression
	elems = append(elems, p.parseExpression(LOWEST))

	for p.peekTokenIs(token.COMMA) {
		p.nextToken() // ","
		p.nextToken() // start of next expression
		elems = append(elems, p.parseExpression(LOWEST))
	}

	if !p.expectPeek(token.RBRACKET) {
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

	if !p.expectPeek(token.RPAREN) {
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
	if p.curTokenIs(token.ELSE) {
		ie.Alternative = p.parseBlockStatement()
	}

	return ie
}

func (p *Parser) parseCallExpression(receiver ast.Expression) ast.Expression {
	if p.curTokenIs(token.LPAREN) { // call expression doesn't have a receiver foo(x) || foo()
		// receiver is self
		selfTok := token.Token{Type: token.SELF, Literal: "self", Line: p.curToken.Line}
		self := &ast.SelfExpression{Token: selfTok}
		m := receiver.(*ast.Identifier).Value
		receiver = self

		// current token is identifier (method name)
		exp := &ast.CallExpression{Token: p.curToken, Receiver: receiver, Method: m}

		exp.Arguments = p.parseCallArguments()
		return exp

	} else { // call expression has a receiver like: p.foo
		exp := &ast.CallExpression{Token: p.curToken, Receiver: receiver}

		// check if method name is identifier
		if !p.expectPeek(token.IDENT) {
			return nil
		}

		exp.Method = p.curToken.Literal

		if p.peekTokenIs(token.LPAREN) { // p.foo.bar; || p.foo; || p.foo + 123
			p.nextToken()
			exp.Arguments = p.parseCallArguments()
		} else {
			exp.Arguments = []ast.Expression{}
		}

		return exp
	}
}

func (p *Parser) parseCallArguments() []ast.Expression {
	args := []ast.Expression{}

	if p.peekTokenIs(token.RPAREN) {
		p.nextToken() // ')'
		return args
	}

	p.nextToken() // start of first expression
	args = append(args, p.parseExpression(LOWEST))

	for p.peekTokenIs(token.COMMA) {
		p.nextToken() // ","
		p.nextToken() // start of next expression
		args = append(args, p.parseExpression(LOWEST))
	}

	if !p.expectPeek(token.RPAREN) {
		return nil
	}

	return args
}
