package parser

import (
	"fmt"
	"github.com/goby-lang/goby/token"
)

func (p *Parser) peekPrecedence() int {
	if p, ok := precedence[p.peekToken.Type]; ok {
		return p
	}

	return LOWEST
}

func (p *Parser) curPrecedence() int {
	if p, ok := precedence[p.curToken.Type]; ok {
		return p
	}

	return LOWEST
}

func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.l.NextToken()
}

func (p *Parser) curTokenIs(t token.Type) bool {
	return p.curToken.Type == t
}

func (p *Parser) peekTokenIs(t token.Type) bool {
	return p.peekToken.Type == t
}

func (p *Parser) peekError(t token.Type) {
	msg := fmt.Sprintf("expected next token to be %s, got %s instead. Line: %d", t, p.peekToken.Type, p.peekToken.Line)
	p.errors = append(p.errors, msg)
}

func (p *Parser) expectPeek(t token.Type) bool {
	if p.peekTokenIs(t) {
		p.nextToken()
		return true
	}
	p.peekError(t)
	return false
}

func (p *Parser) registerPrefix(tokenType token.Type, fn prefixParseFn) {
	p.prefixParseFns[tokenType] = fn
}

func (p *Parser) registerInfix(tokenType token.Type, fn infixParseFn) {
	p.infixParseFns[tokenType] = fn
}

func (p *Parser) noPrefixParseFnError(t token.Type) {
	msg := fmt.Sprintf("no prefix function for %s. Line: %d", t, p.curToken.Line)
	p.errors = append(p.errors, msg)
}

func (p *Parser) peekTokenAtSameLine() bool {
	return p.curToken.Line == p.peekToken.Line
}
