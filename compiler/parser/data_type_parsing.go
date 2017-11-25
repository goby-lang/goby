package parser

import (
	"fmt"
	"github.com/goby-lang/goby/compiler/ast"
	"github.com/goby-lang/goby/compiler/token"
	"strconv"
)

func (p *Parser) parseIntegerLiteral() ast.Expression {
	lit := &ast.IntegerLiteral{BaseNode: &ast.BaseNode{Token: p.curToken}}

	value, err := strconv.ParseInt(lit.TokenLiteral(), 0, 64)
	if err != nil {
		p.error = newTypeParsingError(lit.TokenLiteral(), "integer", p.curToken.Line)
		return nil
	}

	lit.Value = int(value)

	return lit
}

func (p *Parser) parseFloatLiteral(integerPart ast.Expression) ast.Expression {
	// Get the fractional part of the token
	p.nextToken()

	floatTok := token.Token{
		Type:    token.Float,
		Literal: fmt.Sprintf("%s.%s", integerPart.String(), p.curToken.Literal),
		Line:    p.curToken.Line,
	}
	lit := &ast.FloatLiteral{BaseNode: &ast.BaseNode{Token: floatTok}}
	value, err := strconv.ParseFloat(lit.TokenLiteral(), 64)
	if err != nil {
		p.error = newTypeParsingError(lit.TokenLiteral(), "float", p.curToken.Line)
		return nil
	}
	lit.Value = float64(value)
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
		p.error = newTypeParsingError(lit.TokenLiteral(), "boolean", p.curToken.Line)
		return nil
	}

	lit.Value = value

	return lit
}

func (p *Parser) parseNilExpression() ast.Expression {
	return &ast.NilExpression{BaseNode: &ast.BaseNode{Token: p.curToken}}
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
	case token.Constant, token.Ident:
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

func (p *Parser) parseArrayExpression() ast.Expression {
	arr := &ast.ArrayExpression{BaseNode: &ast.BaseNode{Token: p.curToken}}
	arr.Elements = p.parseArrayElements()
	return arr
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
