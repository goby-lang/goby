package parser

import (
	"fmt"
	"strconv"

	"github.com/goby-lang/goby/compiler/ast"
	"github.com/goby-lang/goby/compiler/token"
)

var arguments = map[token.Type]bool{
	token.Int:              true,
	token.String:           true,
	token.True:             true,
	token.False:            true,
	token.Null:             true,
	token.InstanceVariable: true,
	token.Ident:            true,
	token.Constant:         true,
}

var precedence = map[token.Type]int{
	token.Eq:                 EQUALS,
	token.NotEq:              EQUALS,
	token.LT:                 COMPARE,
	token.LTE:                COMPARE,
	token.GT:                 COMPARE,
	token.GTE:                COMPARE,
	token.COMP:               COMPARE,
	token.And:                LOGIC,
	token.Or:                 LOGIC,
	token.OrEq:               LOGIC,
	token.Range:              RANGE,
	token.Plus:               SUM,
	token.PlusEq:             SUM,
	token.Minus:              SUM,
	token.MinusEq:            SUM,
	token.Incr:               SUM,
	token.Decr:               SUM,
	token.Modulo:             SUM,
	token.Assign:             ASSIGN,
	token.Slash:              PRODUCT,
	token.Asterisk:           PRODUCT,
	token.Pow:                PRODUCT,
	token.LBracket:           INDEX,
	token.Dot:                CALL,
	token.LParen:             CALL,
	token.ResolutionOperator: CALL,
}

// Constants for denoting precedence
const (
	_ int = iota
	LOWEST
	NORMAL
	LOGIC
	RANGE
	EQUALS
	COMPARE
	ASSIGN
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
	parseFn := p.prefixParseFns[p.curToken.Type]
	if parseFn == nil {
		p.noPrefixParseFnError(p.curToken.Type)
		return nil
	}

	if p.curTokenIs(token.Ident) && p.fsm.Is(normal) {
		if p.peekTokenIs(token.Do) {
			return p.parseCallExpressionWithoutParenAndReceiver(p.curToken)
		}

		/*
			This means method call with arguments but without parens like:

			```
			foo x
			foo @bar
			foo 10
			```

			Program like

			```
			foo
			x
			```

			will also enter this condition first, but we'll check if those two token is at same line in the parsing function
		*/
		if arguments[p.peekToken.Type] && p.peekTokenAtSameLine() {
			// Method token
			tok := p.curToken
			p.nextToken()
			return p.parseCallExpressionWithoutParenAndReceiver(tok)
		}
	}

	leftExp := parseFn()

	for !p.peekTokenIs(token.Semicolon) && precedence < p.peekPrecedence() && p.peekTokenAtSameLine() {

		infixFn := p.infixParseFns[p.peekToken.Type]
		if infixFn == nil {
			return leftExp
		}
		p.nextToken()
		leftExp = infixFn(leftExp)
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
	c := &ast.Constant{Token: p.curToken, Value: p.curToken.Literal}

	if p.peekTokenIs(token.ResolutionOperator) {
		c.IsNamespace = true
		p.nextToken()
		return p.parseInfixExpression(c)
	}

	return c
}

func (p *Parser) parseInstanceVariable() ast.Expression {
	return &ast.InstanceVariable{Token: p.curToken, Value: p.curToken.Literal}
}

func (p *Parser) parseIntegerLiteral() ast.Expression {
	lit := &ast.IntegerLiteral{Token: p.curToken}

	value, err := strconv.ParseInt(lit.TokenLiteral(), 0, 64)
	if err != nil {
		msg := fmt.Sprintf("could not parse %q as integer", lit.TokenLiteral())
		panic(msg)
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
	lit := &ast.BooleanExpression{Token: p.curToken}

	value, err := strconv.ParseBool(lit.TokenLiteral())
	if err != nil {
		msg := fmt.Sprintf("could not parse %q as boolean", lit.TokenLiteral())
		panic(msg)
	}

	lit.Value = value

	return lit
}

func (p *Parser) parseNilExpression() ast.Expression {
	return &ast.NilExpression{Token: p.curToken}
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
	value = p.parseExpression(NORMAL)
	pairs[key] = value
}

func (p *Parser) parseArrayExpression() ast.Expression {
	arr := &ast.ArrayExpression{Token: p.curToken}
	arr.Elements = p.parseArrayElements()
	return arr
}

func (p *Parser) parseIndexExpression(left ast.Expression) ast.Expression {
	callExpression := &ast.CallExpression{Receiver: left, Method: "[]", Token: p.curToken}

	if p.peekTokenIs(token.RBracket) {
		callExpression.Arguments = []ast.Expression{}
		p.nextToken()
		return callExpression
	}

	p.nextToken()

	callExpression.Arguments = []ast.Expression{p.parseExpression(NORMAL)}

	if !p.expectPeek(token.RBracket) {
		return nil
	}

	// Assign value to index
	if p.peekTokenIs(token.Assign) {
		p.nextToken()
		p.nextToken()
		assignValue := p.parseExpression(NORMAL)
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
	elems = append(elems, p.parseExpression(NORMAL))

	for p.peekTokenIs(token.Comma) {
		p.nextToken() // ","
		p.nextToken() // start of next expression
		elems = append(elems, p.parseExpression(NORMAL))
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

func (p *Parser) parseAssignExpression(v ast.Expression) ast.Expression {
	variable, ok := v.(ast.Variable)

	if !ok {
		p.error = &Error{Message: fmt.Sprintf("Can't assign value to %s. Line: %d", v.String(), p.curToken.Line), errType: InvalidAssignmentError}
	}

	// Pure assignment case
	if p.curTokenIs(token.Assign) {
		exp := &ast.AssignExpression{
			Token:    p.curToken,
			Variable: variable,
			Operator: p.curToken.Literal,
		}

		precedence := p.curPrecedence()
		p.nextToken()
		exp.Value = p.parseExpression(precedence)

		return exp
	}

	// Syntax Surgar: Assignment with operator case
	operator := token.Token{Line: p.curToken.Line}
	assignment := token.Token{Type: token.Assign, Literal: "=", Line: p.curToken.Line}
	switch p.curToken.Type {
	case token.PlusEq:
		operator.Type = token.Plus
		operator.Literal = "+"
	case token.MinusEq:
		operator.Type = token.Minus
		operator.Literal = "-"
	case token.OrEq:
		operator.Type = token.Or
		operator.Literal = "||"
	}
	p.nextToken()
	infixExp := &ast.InfixExpression{
		Token:    operator,
		Left:     variable,
		Operator: operator.Literal,
		Right:    p.parseExpression(LOWEST),
	}

	exp := &ast.AssignExpression{
		Token:    assignment,
		Variable: variable,
		Operator: assignment.Literal,
		Value:    infixExp,
	}

	return exp
}

func (p *Parser) parseGroupedExpression() ast.Expression {
	p.nextToken()

	exp := p.parseExpression(NORMAL)

	if !p.expectPeek(token.RParen) {
		return nil
	}

	return exp
}

func (p *Parser) parseIfExpression() ast.Expression {
	ie := &ast.IfExpression{Token: p.curToken}
	p.nextToken()
	ie.Condition = p.parseExpression(NORMAL)
	ie.Consequence = p.parseBlockStatement()

	// curToken is now ELSE or RBRACE
	if p.curTokenIs(token.Else) {
		ie.Alternative = p.parseBlockStatement()
	}

	return ie
}

func (p *Parser) parseCallExpressionWithoutParenAndReceiver(methodToken token.Token) ast.Expression {
	p.fsm.Event(parseFuncCall)
	// real receiver is self
	selfTok := token.Token{Type: token.Self, Literal: "self", Line: p.curToken.Line}
	self := &ast.SelfExpression{Token: selfTok}

	// current token might be the first argument
	//     method name      |       argument
	// foo <- method token     x <- current token
	exp := &ast.CallExpression{Token: methodToken, Receiver: self, Method: methodToken.Literal}

	if p.curToken.Line == methodToken.Line { // foo x
		exp.Arguments = p.parseCallArgumentsWithoutParens()
	}

	p.fsm.Event(normal)

	// Parse block
	if p.peekTokenIs(token.Do) && p.acceptBlock {
		p.parseBlockParameters(exp)
	}

	return exp
}

func (p *Parser) parseCallExpressionWithParen(receiver ast.Expression) ast.Expression {
	p.fsm.Event(parseFuncCall)
	m := receiver.(*ast.Identifier)
	mn := m.Value

	// real receiver is self
	selfTok := token.Token{Type: token.Self, Literal: "self", Line: p.curToken.Line}
	self := &ast.SelfExpression{Token: selfTok}
	receiver = self

	exp := &ast.CallExpression{Token: m.Token, Receiver: receiver, Method: mn}
	exp.Arguments = p.parseCallArguments()

	p.fsm.Event(normal)

	// Parse block
	if p.peekTokenIs(token.Do) && p.acceptBlock {
		p.parseBlockParameters(exp)
	}

	return exp
}

func (p *Parser) parseCallExpressionWithDot(receiver ast.Expression) ast.Expression {
	p.fsm.Event(parseFuncCall)
	exp := &ast.CallExpression{Token: p.curToken, Receiver: receiver}

	// check if method name is identifier
	if !p.expectPeek(token.Ident) {
		return nil
	}

	exp.Method = p.curToken.Literal

	if p.peekTokenIs(token.LParen) { // p.foo(x)
		p.nextToken()
		exp.Arguments = p.parseCallArguments()
	} else if p.peekTokenIs(token.Dot) { // p.foo.bar
		exp.Arguments = []ast.Expression{}
	} else if arguments[p.peekToken.Type] && p.peekTokenAtSameLine() { // p.foo x, y, z || p.foo x
		p.nextToken()
		exp.Arguments = p.parseCallArgumentsWithoutParens()
	}

	p.fsm.Event(normal)

	// Setter method call like: p.foo = x
	if p.peekTokenIs(token.Assign) {
		exp.Method = exp.Method + "="
		p.nextToken()
		p.nextToken()
		exp.Arguments = append(exp.Arguments, p.parseExpression(NORMAL))
	}

	// Parse block
	if p.peekTokenIs(token.Do) && p.acceptBlock {
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
	args = append(args, p.parseExpression(NORMAL))

	for p.peekTokenIs(token.Comma) {
		p.nextToken() // ","
		p.nextToken() // start of next expression
		args = append(args, p.parseExpression(NORMAL))
	}

	if !p.expectPeek(token.RParen) {
		return nil
	}

	return args
}

func (p *Parser) parseCallArgumentsWithoutParens() []ast.Expression {
	args := []ast.Expression{}

	args = append(args, p.parseExpression(NORMAL))

	for p.peekTokenIs(token.Comma) {
		p.nextToken() // ","
		p.nextToken() // start of next expression
		args = append(args, p.parseExpression(NORMAL))
	}

	if p.peekTokenAtSameLine() {
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

	if arguments[p.peekToken.Type] && p.peekTokenAtSameLine() { // yield 123
		p.nextToken()
		ye.Arguments = p.parseCallArgumentsWithoutParens()
	}

	return ye
}

func (p *Parser) parseRangeExpression(left ast.Expression) ast.Expression {
	exp := &ast.RangeExpression{
		Token: p.curToken,
		Start: left,
	}

	precedence := p.curPrecedence()
	p.nextToken()
	exp.End = p.parseExpression(precedence)

	return exp
}
