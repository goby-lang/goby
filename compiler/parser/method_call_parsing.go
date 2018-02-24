package parser

import (
	"github.com/goby-lang/goby/compiler/ast"
	"github.com/goby-lang/goby/compiler/parser/arguments"
	"github.com/goby-lang/goby/compiler/parser/events"
	"github.com/goby-lang/goby/compiler/parser/precedence"
	"github.com/goby-lang/goby/compiler/token"
)

func (p *Parser) parseCallExpressionWithoutReceiver(receiver ast.Expression) ast.Expression {
	methodToken := receiver.(*ast.Identifier).Token

	exp := &ast.CallExpression{BaseNode: &ast.BaseNode{}}

	oldState := p.fsm.Current()
	p.fsm.Event(events.ParseFuncCall)
	// real receiver is self
	selfTok := token.Token{Type: token.Self, Literal: "self", Line: p.curToken.Line}
	self := &ast.SelfExpression{BaseNode: &ast.BaseNode{Token: selfTok}}

	// current token might be the first argument
	//     method name      |       argument
	// foo <- method token     x <- current token
	exp.Token = methodToken
	exp.Receiver = self
	exp.Method = methodToken.Literal

	if p.curTokenIs(token.LParen) {
		exp.Arguments = p.parseCallArgumentsWithParens()
	} else if p.curToken.Line == methodToken.Line && p.curToken != methodToken { // 'foo x' but not 'thread'
		exp.Arguments = p.parseCallArguments()
	}

	p.fsm.Event(events.EventTable[oldState])

	if p.peekTokenIs(token.Do) && p.acceptBlock { // foo do
		p.parseBlockArgument(exp)
	}

	return exp
}

func (p *Parser) parseCallExpressionWithReceiver(receiver ast.Expression) ast.Expression {
	exp := &ast.CallExpression{BaseNode: &ast.BaseNode{}}

	oldState := p.fsm.Current()
	p.fsm.Event(events.ParseFuncCall)

	// check if method name is identifier
	if !p.expectPeek(token.Ident) {
		return nil
	}

	exp.Token = p.curToken
	exp.Receiver = receiver
	exp.Method = p.curToken.Literal

	// The next token needs to be the same line
	// Otherwise something like Range value might be treated as arguments
	// See issue #532
	if p.peekTokenAtSameLine() {
		switch p.peekToken.Type {
		case token.LParen: // p.foo(x)
			p.nextToken()
			exp.Arguments = p.parseCallArgumentsWithParens()
		case token.Assign: // Setter method call like: p.foo = x
			exp.Method = exp.Method + "="
			p.nextToken()
			p.nextToken()
			exp.Arguments = append(exp.Arguments, p.parseExpression(precedence.Normal))
		default:
			if arguments.Tokens[p.peekToken.Type] && p.peekTokenAtSameLine() { // p.foo x, y, z || p.foo x
				p.nextToken()
				exp.Arguments = p.parseCallArguments()
			} else {
				exp.Arguments = []ast.Expression{}
			}
		}
	} else {
		exp.Arguments = []ast.Expression{}
	}

	p.fsm.Event(events.EventTable[oldState])

	// Parse block
	if p.peekTokenIs(token.Do) && p.acceptBlock {
		p.parseBlockArgument(exp)
	}

	return exp
}

func (p *Parser) parseCallArgumentsWithParens() []ast.Expression {
	args := []ast.Expression{}

	if p.peekTokenIs(token.RParen) {
		p.nextToken() // ')'
		return args
	}

	p.nextToken() // move to first argument token

	args = p.parseCallArguments()

	if !p.expectPeek(token.RParen) {
		return nil
	}

	return args
}

func (p *Parser) parseCallArguments() []ast.Expression {
	args := []ast.Expression{}

	args = append(args, p.parseExpression(precedence.Normal))

	for p.peekTokenIs(token.Comma) {
		p.nextToken() // ","
		p.nextToken() // start of next expression
		args = append(args, p.parseExpression(precedence.Normal))
	}

	return args
}

func (p *Parser) parseBlockArgument(exp *ast.CallExpression) {
	p.nextToken()

	// Parse block arguments
	if p.peekTokenIs(token.Bar) {
		var params []*ast.Identifier

		p.nextToken()
		p.nextToken()

		param := &ast.Identifier{BaseNode: &ast.BaseNode{Token: p.curToken}, Value: p.curToken.Literal}
		params = append(params, param)

		for p.peekTokenIs(token.Comma) {
			p.nextToken()
			p.nextToken()
			param := &ast.Identifier{BaseNode: &ast.BaseNode{Token: p.curToken}, Value: p.curToken.Literal}
			params = append(params, param)
		}

		if !p.expectPeek(token.Bar) {
			return
		}

		exp.BlockArguments = params
	}

	exp.Block = p.parseBlockStatement(token.End)
	exp.Block.KeepLastValue()
}
