package parser

import (
	"fmt"
	"github.com/st0012/rooby/ast"
	"github.com/st0012/rooby/token"
	"strconv"
)

var precedence = map[token.TokenType]int{
	token.EQ:       EQUALS,
	token.NOT_EQ:   EQUALS,
	token.LT:       LESSGREATER,
	token.GT:       LESSGREATER,
	token.PLUS:     SUM,
	token.MINUS:    SUM,
	token.SLASH:    PRODUCT,
	token.ASTERISK: PRODUCT,
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

	for !p.peekTokenIs(token.SEMICOLON) && precendence < p.peekPrecedence() {
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

	lit.Value = value

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

	p.parseCondition(ie)
	p.parseConsequence(ie)

	// curToken is now ELSE or RBRACE
	if p.peekTokenIs(token.ELSE) {
		p.parseAlternative(ie)
	}

	return ie
}

func (p *Parser) parseCallExpression(receiver ast.Expression) ast.Expression {
	if p.curTokenIs(token.LPAREN) { // foo(x)
		// receiver is self
		selfTok := token.Token{Type: token.SELF, Literal: "self", Line: p.curToken.Line}
		self := &ast.SelfExpression{Token: selfTok}

		// current token is identifier (method name)
		exp := &ast.CallExpression{Token: p.curToken, Receiver: self, Method: receiver.(*ast.Identifier)}

		exp.Arguments = p.parseCallArguments()
		return exp

	} else { // p.foo
		exp := &ast.CallExpression{Token: p.curToken, Receiver: receiver}

		// check if method name is identifier
		if !p.expectPeek(token.IDENT) {
			return nil
		}

		exp.Method = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}

		if p.peekTokenIs(token.DOT) || p.peekTokenIs(token.SEMICOLON) { // p.foo.bar; || p.foo;
			exp.Arguments = []ast.Expression{}
		} else if !p.expectPeek(token.LPAREN) {
			return nil
		} else {
			exp.Arguments = p.parseCallArguments()
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

func (p *Parser) parseBlockStatement() *ast.BlockStatement {
	// curToken is {
	bs := &ast.BlockStatement{Token: p.curToken}
	bs.Statements = []ast.Statement{}

	p.nextToken()

	for !p.curTokenIs(token.RBRACE) {
		stmt := p.parseStatement()
		if stmt != nil {
			bs.Statements = append(bs.Statements, stmt)
		}
		p.nextToken()
	}

	return bs
}

func (p *Parser) parseCondition(ie *ast.IfExpression) *ast.IfExpression {
	if !p.expectPeek(token.LPAREN) {
		return nil
	}

	p.nextToken()
	ie.Condition = p.parseExpression(LOWEST)

	if !p.expectPeek(token.RPAREN) {
		return nil
	}

	return ie
}

func (p *Parser) parseConsequence(ie *ast.IfExpression) *ast.IfExpression {
	if !p.expectPeek(token.LBRACE) {
		return nil
	}

	ie.Consequence = p.parseBlockStatement()

	return ie
}

func (p *Parser) parseAlternative(ie *ast.IfExpression) *ast.IfExpression {
	p.nextToken()

	if !p.expectPeek(token.LBRACE) {
		return nil
	}
	ie.Alternative = p.parseBlockStatement()

	return ie
}
