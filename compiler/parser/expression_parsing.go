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
	token.Range:              RANGE,
	token.Plus:               SUM,
	token.Minus:              SUM,
	token.Incr:               SUM,
	token.Decr:               SUM,
	token.Modulo:             SUM,
	token.Slash:              PRODUCT,
	token.Asterisk:           PRODUCT,
	token.Pow:                PRODUCT,
	token.LBracket:           INDEX,
	token.Dot:                CALL,
	token.LParen:             CALL,
	token.ResolutionOperator: CALL,
	token.Assign:             ASSIGN,
	token.PlusEq:             ASSIGN,
	token.MinusEq:            ASSIGN,
	token.OrEq:               ASSIGN,
}

// Constants for denoting precedence
const (
	_ int = iota
	LOWEST
	NORMAL
	ASSIGN
	LOGIC
	RANGE
	EQUALS
	COMPARE
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

	/*
		Parse method call without explicit receiver and doesn't have parens around arguments. Here's some examples:

		When state is normal:
		```
		foo 10
		```

		When state is parseAssignment:
		```
		a = foo 10
		```
	*/
	if p.curTokenIs(token.Ident) && (p.fsm.Is(normal) || p.fsm.Is(parsingAssignment)) {
		if p.peekTokenIs(token.Do) {
			method := p.parseIdentifier()
			return p.parseCallExpressionWithoutReceiver(method)
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
			method := p.parseIdentifier()
			p.nextToken()
			return p.parseCallExpressionWithoutReceiver(method)
		}
	}

	leftExp := parseFn()

	/*
		Precedence example:

		```
		1 + 1 * 5 == 1 + (1 * 5)

		```

		Because "*"'s precedence is PRODUCT which is higher than "+"'s precedence SUM, we'll parse "*" first.

	*/

	for !p.peekTokenIs(token.Semicolon) &&
		(precedence < p.peekPrecedence() || (p.fsm.Is(parsingAssignment) && p.peekTokenIs(token.Assign))) &&
		// This is for preventing parser treat next line's expression as function's argument.
		p.peekTokenAtSameLine() {

		infixFn := p.infixParseFns[p.peekToken.Type]
		if infixFn == nil {
			return leftExp
		}
		p.nextToken()
		leftExp = infixFn(leftExp)
	}

	if p.peekTokenIs(token.Semicolon) {
		p.nextToken()
	}

	return leftExp
}

func (p *Parser) parseSelfExpression() ast.Expression {
	return &ast.SelfExpression{BaseNode: &ast.BaseNode{Token: p.curToken}}
}

func (p *Parser) parseIdentifier() ast.Expression {
	return &ast.Identifier{BaseNode: &ast.BaseNode{Token: p.curToken}, Value: p.curToken.Literal}
}



func (p *Parser) parseConstant() ast.Expression {
	c := &ast.Constant{BaseNode: &ast.BaseNode{Token: p.curToken}, Value: p.curToken.Literal}

	if p.peekTokenIs(token.ResolutionOperator) {
		c.IsNamespace = true
		p.nextToken()
		return p.parseInfixExpression(c)
	}

	return c
}

func (p *Parser) parseInstanceVariable() ast.Expression {
	return &ast.InstanceVariable{BaseNode: &ast.BaseNode{Token: p.curToken}, Value: p.curToken.Literal}
}

func (p *Parser) parseIntegerLiteral() ast.Expression {
	lit := &ast.IntegerLiteral{BaseNode: &ast.BaseNode{Token: p.curToken}}

	value, err := strconv.ParseInt(lit.TokenLiteral(), 0, 64)
	if err != nil {
		msg := fmt.Sprintf("could not parse %q as integer", lit.TokenLiteral())
		panic(msg)
	}

	lit.Value = int(value)

	return lit
}

func (p *Parser) parseStringLiteral() ast.Expression {
	lit := &ast.StringLiteral{BaseNode: &ast.BaseNode{Token: p.curToken}}
	lit.Value = p.curToken.Literal

	return lit
}

func (p *Parser) parseBooleanLiteral() ast.Expression {
	lit := &ast.BooleanExpression{BaseNode: &ast.BaseNode{Token: p.curToken}}

	value, err := strconv.ParseBool(lit.TokenLiteral())
	if err != nil {
		msg := fmt.Sprintf("could not parse %q as boolean", lit.TokenLiteral())
		panic(msg)
	}

	lit.Value = value

	return lit
}

func (p *Parser) parseNilExpression() ast.Expression {
	return &ast.NilExpression{BaseNode: &ast.BaseNode{Token: p.curToken}}
}

func (p *Parser) parsePostfixExpression(receiver ast.Expression) ast.Expression {
	arguments := []ast.Expression{}
	return &ast.CallExpression{BaseNode: &ast.BaseNode{Token: p.curToken}, Receiver: receiver, Method: p.curToken.Literal, Arguments: arguments}
}

func (p *Parser) parseHashExpression() ast.Expression {
	hash := &ast.HashExpression{BaseNode: &ast.BaseNode{Token: p.curToken}}
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

	p.nextToken()

	switch p.curToken.Type {
	case token.Ident:
		key = p.parseIdentifier().(ast.Variable).ReturnValue()
	case token.Constant:
		key = p.parseIdentifier().(ast.Variable).ReturnValue()
	default:
		return
	}

	if !p.expectPeek(token.Colon) {
		return
	}

	p.nextToken()
	value = p.parseExpression(NORMAL)
	pairs[key] = value
}

func (p *Parser) parseKeywordArgumentsExpression() ast.Expression {
	hash := &ast.HashExpression{BaseNode: &ast.BaseNode{Token: p.curToken}}
	hash.Data = p.parseKeywordArguments()
	return hash
}

func (p *Parser) parseKeywordArguments() map[string]ast.Expression {
	pairs := map[string]ast.Expression{}

	if p.peekTokenIs(token.RParen) {
		p.nextToken() // ')'
		return pairs
	}

	p.parseKeywordArgument(pairs)

	for p.peekTokenIs(token.Comma) {
		p.nextToken()

		p.parseKeywordArgument(pairs)
	}

	if !p.peekTokenIs(token.RParen) {
		return nil
	}

	return pairs
}

func (p *Parser) parseKeywordArgument(pairs map[string]ast.Expression) {
	var key string

	p.nextToken()

	switch p.curToken.Type {
	case token.Ident:
		key = p.parseIdentifier().(ast.Variable).ReturnValue()
	case token.Constant:
		key = p.parseIdentifier().(ast.Variable).ReturnValue()
	default:
		return
	}

	if !p.expectPeek(token.Colon) {
		return
	}

	// Keyword argument without default value
	if p.peekTokenIs(token.Comma) || p.peekTokenIs(token.RParen){
		pairs[key] = nil
	} else {
		p.nextToken()
		pairs[key] = p.parseExpression(NORMAL)
	}
}

func (p *Parser) parseArrayExpression() ast.Expression {
	arr := &ast.ArrayExpression{BaseNode: &ast.BaseNode{Token: p.curToken}}
	arr.Elements = p.parseArrayElements()
	return arr
}

func (p *Parser) parseIndexExpression(left ast.Expression) ast.Expression {
	callExpression := &ast.CallExpression{Receiver: left, Method: "[]", BaseNode: &ast.BaseNode{Token: p.curToken}}

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
		BaseNode: &ast.BaseNode{Token: p.curToken},
		Operator: p.curToken.Literal,
	}

	p.nextToken()

	pe.Right = p.parseExpression(PREFIX)

	return pe
}

func (p *Parser) parseInfixExpression(left ast.Expression) ast.Expression {
	exp := &ast.InfixExpression{
		BaseNode: &ast.BaseNode{Token: p.curToken},
		Left:     left,
		Operator: p.curToken.Literal,
	}

	precedence := p.curPrecedence()
	p.nextToken()

	if exp.Operator == "||" || exp.Operator == "&&" {
		precedence = NORMAL
	}

	exp.Right = p.parseExpression(precedence)

	return exp
}

func (p *Parser) parseAssignExpression(v ast.Expression) ast.Expression {
	var value ast.Expression
	var tok token.Token
	exp := &ast.AssignExpression{BaseNode: &ast.BaseNode{}}

	if !p.fsm.Is(parsingFuncCall) {
		exp.MarkAsStmt()
	}

	oldState := p.fsm.Current()
	p.fsm.Event(parseAssignment)

	switch v := v.(type) {
	case ast.Variable:
		exp.Variables = []ast.Expression{v}
	case *ast.MultiVariableExpression:
		exp.Variables = v.Variables
	case *ast.CallExpression:
		/*
			for cases like: `a[i] += b`
			which needs to be expand to

			a[i] = a[i] + b

			CallExp = CallExp + Expression
		*/

		if v.Method == "[]" {
			value = p.expandAssignmentValue(v)

			callExp := &ast.CallExpression{
				BaseNode:  &ast.BaseNode{},
				Method:    "[]=",
				Arguments: []ast.Expression{v.Arguments[0], value},
				Receiver:  v.Receiver,
			}
			return callExp
		}

		p.error = &Error{Message: fmt.Sprintf("Can't assign value to %s. Line: %d", v.String(), p.curToken.Line), errType: InvalidAssignmentError}
	default:
		p.error = &Error{Message: fmt.Sprintf("Can't assign value to %s. Line: %d", v.String(), p.curToken.Line), errType: InvalidAssignmentError}
	}

	if len(exp.Variables) == 1 {
		tok = token.Token{Type: token.Assign, Literal: "=", Line: p.curToken.Line}
		value = p.expandAssignmentValue(v)
	} else {
		tok = p.curToken
		precedence := p.curPrecedence()
		p.nextToken()
		value = p.parseExpression(precedence)
	}

	exp.Token = tok
	exp.Value = value

	event, _ := eventTable[oldState]
	p.fsm.Event(event)

	return exp
}

func (p *Parser) expandAssignmentValue(value ast.Expression) ast.Expression {
	switch p.curToken.Type {
	case token.Assign:
		precedence := p.curPrecedence()
		p.nextToken()
		return p.parseExpression(precedence)
	case token.MinusEq, token.PlusEq, token.OrEq:
		// Syntax Surgar: Assignment with operator case
		infixOperator := token.Token{Line: p.curToken.Line}
		switch p.curToken.Type {
		case token.PlusEq:
			infixOperator.Type = token.Plus
			infixOperator.Literal = "+"
		case token.MinusEq:
			infixOperator.Type = token.Minus
			infixOperator.Literal = "-"
		case token.OrEq:
			infixOperator.Type = token.Or
			infixOperator.Literal = "||"
		}

		p.nextToken()

		return &ast.InfixExpression{
			BaseNode: &ast.BaseNode{Token: infixOperator},
			Left:     value,
			Operator: infixOperator.Literal,
			Right:    p.parseExpression(LOWEST),
		}
	default:
		p.error = &Error{errType: UnexpectedTokenError, Message: fmt.Sprintf("Unexpect token '%s' for assgin expression", p.curToken.Literal)}
		return nil
	}
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
	ie := &ast.IfExpression{BaseNode: &ast.BaseNode{Token: p.curToken}}
	// parse if and elsif expressions
	ie.Conditionals = p.parseConditionalExpressions()

	// curToken is now ELSE or RBRACE
	if p.curTokenIs(token.Else) {
		ie.Alternative = p.parseBlockStatement()
		ie.Alternative.KeepLastValue()
	}

	return ie
}

func (p *Parser) parseConditionalExpressions() []*ast.ConditionalExpression {
	// first conditional expression should start with if
	cs := []*ast.ConditionalExpression{p.parseConditionalExpression()}

	for p.curTokenIs(token.ElsIf) {
		cs = append(cs, p.parseConditionalExpression())
	}

	return cs
}

func (p *Parser) parseConditionalExpression() *ast.ConditionalExpression {
	ce := &ast.ConditionalExpression{BaseNode: &ast.BaseNode{Token: p.curToken}}
	p.nextToken()
	ce.Condition = p.parseExpression(NORMAL)

	ce.Consequence = p.parseBlockStatement()
	ce.Consequence.KeepLastValue()

	return ce
}

func (p *Parser) parseYieldExpression() ast.Expression {
	ye := &ast.YieldExpression{BaseNode: &ast.BaseNode{Token: p.curToken}}

	if p.peekTokenIs(token.LParen) {
		p.nextToken()
		ye.Arguments = p.parseCallArgumentsWithParens()
	}

	if arguments[p.peekToken.Type] && p.peekTokenAtSameLine() { // yield 123
		p.nextToken()
		ye.Arguments = p.parseCallArguments()
	}

	return ye
}

func (p *Parser) parseRangeExpression(left ast.Expression) ast.Expression {
	exp := &ast.RangeExpression{
		BaseNode: &ast.BaseNode{Token: p.curToken},
		Start:    left,
	}

	precedence := p.curPrecedence()
	p.nextToken()
	exp.End = p.parseExpression(precedence)

	return exp
}

func (p *Parser) parseMultiVariables(left ast.Expression) ast.Expression {
	var1, ok := left.(ast.Variable)

	if !ok {
		p.noPrefixParseFnError(p.curToken.Type)
	}

	vars := []ast.Expression{var1}

	p.nextToken()

	exp := p.parseExpression(CALL)

	var2, ok := exp.(ast.Variable)

	if !ok {
		p.noPrefixParseFnError(p.curToken.Type)
	}

	vars = append(vars, var2)

	for p.peekTokenIs(token.Comma) {
		p.nextToken()
		p.nextToken()
		exp := p.parseExpression(CALL) // Use highest precedence

		v, ok := exp.(ast.Variable)

		if !ok {
			p.noPrefixParseFnError(p.curToken.Type)
		}

		vars = append(vars, v)
	}

	result := &ast.MultiVariableExpression{Variables: vars}
	return result
}
