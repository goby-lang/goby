package parser

import (
	"github.com/gooby-lang/gooby/compiler/ast"
	"github.com/gooby-lang/gooby/compiler/parser/precedence"
	"github.com/gooby-lang/gooby/compiler/token"
)

// Case expression forms if statement when parsing it
//
// ```ruby
// case 1
// when 0, 1
//  '0 or 1'
// else
//  'else'
// end
// ```
//
// is the same with if expression below
//
// ```ruby
// if 1 == 0 || 1 == 1
//  '0 or 1'
// else
//  'else'
// end
// ```
//
// TODO Implement '===' method and replace '==' to '===' in Case expression

func (p *Parser) parseCaseExpression() ast.Expression {
	ie := &ast.IfExpression{BaseNode: &ast.BaseNode{Token: p.curToken}}
	ie.Conditionals = p.parseCaseConditionals()

	if p.curTokenIs(token.Else) {
		ie.Alternative = p.parseBlockStatement(token.End)
		ie.Alternative.KeepLastValue()
	}

	return ie
}

// case expression parsing helpers
func (p *Parser) parseCaseConditionals() []*ast.ConditionalExpression {
	var ce []*ast.ConditionalExpression
	var base ast.Expression
	if p.peekTokenIs(token.When) {
		base = &ast.BooleanExpression{BaseNode: &ast.BaseNode{Token: token.Token{Type: token.True, Literal: "true", Line: p.curToken.Line}}, Value: true}
	} else {
		p.nextToken()
		base = p.parseExpression(precedence.Normal)
	}

	p.expectPeek(token.When)

	for p.curTokenIs(token.When) {
		ce = append(ce, p.parseCaseConditional(base))
	}

	return ce
}

func (p *Parser) parseCaseConditional(base ast.Expression) *ast.ConditionalExpression {
	ce := &ast.ConditionalExpression{BaseNode: &ast.BaseNode{Token: p.curToken}}
	p.nextToken()

	ce.Condition = p.parseCaseCondition(base)
	ce.Consequence = p.parseBlockStatement(token.When, token.Else, token.End)
	ce.Consequence.KeepLastValue()

	return ce
}

func (p *Parser) parseCaseCondition(base ast.Expression) *ast.InfixExpression {
	first := p.parseExpression(precedence.Normal)
	infix := newInfixExpression(base, token.Token{Type: token.Eq, Literal: token.Eq}, first)

	for p.peekTokenIs(token.Comma) {
		p.nextToken()
		p.nextToken()

		right := p.parseExpression(precedence.Normal)
		rightInfix := newInfixExpression(base, token.Token{Type: token.Eq, Literal: token.Eq}, right)
		infix = newInfixExpression(infix, token.Token{Type: token.Or, Literal: token.Or}, rightInfix)
	}

	return infix
}

func (p *Parser) parseIfExpression() ast.Expression {
	ie := &ast.IfExpression{BaseNode: &ast.BaseNode{Token: p.curToken}}
	// parse if and elsif expressions
	ie.Conditionals = p.parseConditionalExpressions()

	// curToken is now ELSE or RBRACE
	if p.curTokenIs(token.Else) {
		ie.Alternative = p.parseBlockStatement(token.End)
		ie.Alternative.KeepLastValue()
	}

	return ie
}

// infix expression parsing helpers
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
	ce.Condition = p.parseExpression(precedence.Normal)
	ce.Consequence = p.parseBlockStatement(token.ElsIf, token.Else, token.End)
	ce.Consequence.KeepLastValue()

	return ce
}
