package parser

import (
	"fmt"
	"github.com/goby-lang/goby/compiler/ast"
	"github.com/goby-lang/goby/compiler/parser/arguments"
	"github.com/goby-lang/goby/compiler/parser/errors"
	"github.com/goby-lang/goby/compiler/parser/events"
	"github.com/goby-lang/goby/compiler/parser/precedence"
	"github.com/goby-lang/goby/compiler/parser/states"
	"github.com/goby-lang/goby/compiler/token"
)

type (
	prefixParseFn func() ast.Expression
	infixParseFn  func(ast.Expression) ast.Expression
)

func (p *Parser) parseAssignExpression(v ast.Expression) ast.Expression {
	var value ast.Expression
	var tok token.Token
	exp := &ast.AssignExpression{BaseNode: &ast.BaseNode{}}

	if !p.fsm.Is(states.ParsingFuncCall) {
		exp.MarkAsStmt()
	}

	oldState := p.fsm.Current()
	p.fsm.Event(events.ParseAssignment)

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

		errMsg := fmt.Sprintf("Can't assign value to %s. Line: %d", v.String(), p.curToken.Line)
		p.error = errors.InitError(errMsg, errors.InvalidAssignmentError)
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

	event, _ := events.EventTable[oldState]
	p.fsm.Event(event)

	return exp
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

func (p *Parser) parseDotExpression(receiver ast.Expression) ast.Expression {
	_, ok := receiver.(*ast.IntegerLiteral)

	// When both receiver & caller are integer => Float
	if ok && p.peekTokenIs(token.Int) {
		return p.parseFloatLiteral(receiver)
	}

	// Normal call method expression with receiver
	return p.parseCallExpressionWithReceiver(receiver)
}

func (p *Parser) parseExpression(precedence int) ast.Expression {
	parseFn := p.prefixParseFns[p.curToken.Type]
	if parseFn == nil {
		p.noPrefixParseFnError(p.curToken.Type)
		return nil
	}

	// Prohibit calling a capitalized method on toplevel:
	if p.curTokenIs(token.Constant) && (p.fsm.Is(states.Normal) || p.fsm.Is(states.ParsingAssignment)) {
		if p.peekTokenIs(token.LParen) {
			p.callConstantError(p.curToken.Type)
			return nil
		}
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
	if p.curTokenIs(token.Ident) && (p.fsm.Is(states.Normal) || p.fsm.Is(states.ParsingAssignment)) {
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
		if arguments.Tokens[p.peekToken.Type] && p.peekTokenAtSameLine() {
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

	for !p.peekTokenIs(token.Semicolon) && (precedence < p.peekPrecedence() || p.isParsingAssignment()) && p.peekTokenAtSameLine() {
		infixFn := p.infixParseFns[p.peekToken.Type]

		if infixFn == nil {
			return leftExp
		}

		prevTok := p.curToken
		p.nextToken()

		// normally when current token is Dot, we'll get function `parseCallExpressionWithoutReceiver`
		// so we need to make sure if we're parsing method call or float
		if p.isParsingFloat(prevTok) {
			infixFn = p.parseFloatLiteral
		}

		leftExp = infixFn(leftExp)

		for p.peekTokenIs(token.UnderScore) {
			p.nextToken()
			leftExp = infixFn(leftExp)
		}

	}

	if p.peekTokenIs(token.Semicolon) {
		p.nextToken()
	}

	return leftExp
}

func (p *Parser) parseGetBlockExpression() ast.Expression {
	return &ast.GetBlockExpression{BaseNode: &ast.BaseNode{Token: p.curToken}}
}

func (p *Parser) parseGroupedExpression() ast.Expression {
	p.nextToken()

	exp := p.parseExpression(precedence.Normal)

	if !p.expectPeek(token.RParen) {
		return nil
	}

	return exp
}

func (p *Parser) parseIdentifier() ast.Expression {
	return &ast.Identifier{BaseNode: &ast.BaseNode{Token: p.curToken}, Value: p.curToken.Literal}
}

func (p *Parser) parseIndexExpression(left ast.Expression) ast.Expression {
	callExpression := &ast.CallExpression{Receiver: left, Method: "[]", BaseNode: &ast.BaseNode{Token: p.curToken}}

	if p.peekTokenIs(token.RBracket) {
		callExpression.Arguments = []ast.Expression{}
		p.nextToken()
		return callExpression
	}

	p.nextToken()

	callExpression.Arguments = []ast.Expression{p.parseExpression(precedence.Normal)}

	// Accepting multiple indexing argument
	for p.peekTokenIs(token.Comma) {
		p.nextToken()
		p.nextToken()
		callExpression.Arguments = append(callExpression.Arguments, p.parseExpression(precedence.Normal))
	}

	if !p.expectPeek(token.RBracket) {
		return nil
	}

	// Assign value to index
	if p.peekTokenIs(token.Assign) {
		p.nextToken()
		p.nextToken()
		assignValue := p.parseExpression(precedence.Normal)
		callExpression.Method = "[]="
		callExpression.Arguments = append(callExpression.Arguments, assignValue)
	}

	return callExpression
}

func (p *Parser) parseInfixExpression(left ast.Expression) ast.Expression {
	operator := p.curToken
	preced := p.curPrecedence()

	// prevent "* *" from being parsed
	if p.curToken.Literal == token.Asterisk && p.peekToken.Literal == token.Asterisk {
		msg := fmt.Sprintf("unexpected %s Line: %d", p.curToken.Literal, p.peekToken.Line)
		p.error = errors.InitError(msg, errors.UnexpectedTokenError)
		return nil
	}

	if operator.Literal == "||" || operator.Literal == "&&" {
		preced = precedence.Normal
	}

	p.nextToken()

	return newInfixExpression(left, operator, p.parseExpression(preced))
}

func (p *Parser) parseInstanceVariable() ast.Expression {
	return &ast.InstanceVariable{BaseNode: &ast.BaseNode{Token: p.curToken}, Value: p.curToken.Literal}
}

func (p *Parser) parseMultiVariables(left ast.Expression) ast.Expression {
	var1, ok := left.(ast.Variable)

	if !ok {
		p.noPrefixParseFnError(p.curToken.Type)
	}

	vars := []ast.Expression{var1}

	p.nextToken()

	exp := p.parseExpression(precedence.Call)

	var2, ok := exp.(ast.Variable)

	if !ok {
		p.noPrefixParseFnError(p.curToken.Type)
	}

	vars = append(vars, var2)

	for p.peekTokenIs(token.Comma) {
		p.nextToken()
		p.nextToken()
		exp := p.parseExpression(precedence.Call) // Use highest precedence

		v, ok := exp.(ast.Variable)

		if !ok {
			p.noPrefixParseFnError(p.curToken.Type)
		}

		vars = append(vars, v)
	}

	result := &ast.MultiVariableExpression{Variables: vars}
	return result
}

// this function is only for parsing keyword arguments or keyword params
func (p *Parser) parseArgumentPairExpression(key ast.Expression) ast.Expression {
	exp := &ast.ArgumentPairExpression{BaseNode: &ast.BaseNode{Token: p.curToken}, Key: key}

	switch p.fsm.Current() {
	case states.ParsingMethodParam:
		// when parsing keyword params, there'll be some cases that don't have pair value like
		// ```ruby
		//   def foo(bar:, baz:); end
		// ```
		if p.peekTokenIs(token.Comma) || p.peekTokenIs(token.RParen) {
			return exp
		}

		p.nextToken()
		value := p.parseExpression(precedence.Normal)

		exp.Value = value

		return exp
	case states.ParsingFuncCall:
		p.nextToken()
		value := p.parseExpression(precedence.Normal)

		exp.Value = value

		return exp
	default:
		msg := fmt.Sprintf("unexpected %s Line: %d", p.curToken.Literal, p.peekToken.Line)
		p.error = errors.InitError(msg, errors.UnexpectedTokenError)
		return nil
	}
}

func (p *Parser) parsePrefixExpression() ast.Expression {
	pe := &ast.PrefixExpression{
		BaseNode: &ast.BaseNode{Token: p.curToken},
		Operator: p.curToken.Literal,
	}

	prevToken := p.curToken
	p.nextToken()

	if prevToken.Type == token.Bang {
		pe.Right = p.parseExpression(precedence.BangPrefix)
	} else {
		pe.Right = p.parseExpression(precedence.MinusPrefix)
	}

	return pe
}

func (p *Parser) parseSelfExpression() ast.Expression {
	return &ast.SelfExpression{BaseNode: &ast.BaseNode{Token: p.curToken}}
}

func (p *Parser) parseYieldExpression() ast.Expression {
	ye := &ast.YieldExpression{BaseNode: &ast.BaseNode{Token: p.curToken}}

	if p.peekTokenIs(token.LParen) {
		p.nextToken()
		ye.Arguments = p.parseCallArgumentsWithParens()
	}

	if arguments.Tokens[p.peekToken.Type] && p.peekTokenAtSameLine() { // yield 123
		p.nextToken()
		ye.Arguments = p.parseCallArguments()
	}

	return ye
}

// helpers

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

		return newInfixExpression(value, infixOperator, p.parseExpression(precedence.Lowest))
	default:
		p.peekError(p.curToken.Type)
		return nil
	}
}

func (p *Parser) isParsingFloat(leftTok token.Token) bool {
	sameLine := leftTok.Line == p.peekToken.Line
	bothAreInt := leftTok.Type == p.peekToken.Type && leftTok.Type == token.Int
	return p.curTokenIs(token.Dot) && bothAreInt && sameLine
}

func (p *Parser) isParsingAssignment() bool {
	return p.fsm.Is(states.ParsingAssignment) && p.peekTokenIs(token.Assign)
}

func newInfixExpression(left ast.Expression, operator token.Token, right ast.Expression) *ast.InfixExpression {
	return &ast.InfixExpression{
		BaseNode: &ast.BaseNode{Token: operator},
		Left:     left,
		Operator: operator.Literal,
		Right:    right,
	}
}
