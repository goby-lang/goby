package parser

import (
	"fmt"
	"github.com/st0012/Rooby/ast"
	"github.com/st0012/Rooby/lexer"
	"github.com/st0012/Rooby/token"
	"strings"
)

type Parser struct {
	l      *lexer.Lexer
	errors []string

	curToken  token.Token
	peekToken token.Token

	prefixParseFns map[token.TokenType]prefixParseFn
	infixParseFns  map[token.TokenType]infixParseFn
}

func New(l *lexer.Lexer) *Parser {
	p := &Parser{
		l:      l,
		errors: []string{},
	}

	// Read two tokens, so curToken and peekToken are both set.
	p.nextToken()
	p.nextToken()

	p.prefixParseFns = make(map[token.TokenType]prefixParseFn)
	p.registerPrefix(token.IDENT, p.parseIdentifier)
	p.registerPrefix(token.CONSTANT, p.parseConstant)
	p.registerPrefix(token.INSTANCE_VARIABLE, p.parseInstanceVariable)
	p.registerPrefix(token.INT, p.parseIntegerLiteral)
	p.registerPrefix(token.STRING, p.parseStringLiteral)
	p.registerPrefix(token.TRUE, p.parseBooleanLiteral)
	p.registerPrefix(token.FALSE, p.parseBooleanLiteral)
	p.registerPrefix(token.MINUS, p.parsePrefixExpression)
	p.registerPrefix(token.BANG, p.parsePrefixExpression)
	p.registerPrefix(token.LPAREN, p.parseGroupedExpression)
	p.registerPrefix(token.IF, p.parseIfExpression)
	p.registerPrefix(token.SELF, p.parseSelfExpression)
	p.registerPrefix(token.LBRACKET, p.parseArrayExpression)
	p.registerPrefix(token.LBRACE, p.parseHashExpression)
	p.registerPrefix(token.SEMICOLON, p.parseSemicolon)

	p.infixParseFns = make(map[token.TokenType]infixParseFn)
	p.registerInfix(token.PLUS, p.parseInfixExpression)
	p.registerInfix(token.MINUS, p.parseInfixExpression)
	p.registerInfix(token.SLASH, p.parseInfixExpression)
	p.registerInfix(token.EQ, p.parseInfixExpression)
	p.registerInfix(token.ASTERISK, p.parseInfixExpression)
	p.registerInfix(token.NOT_EQ, p.parseInfixExpression)
	p.registerInfix(token.LT, p.parseInfixExpression)
	p.registerInfix(token.GT, p.parseInfixExpression)
	p.registerInfix(token.DOT, p.parseCallExpression)
	p.registerInfix(token.LPAREN, p.parseCallExpression)
	p.registerInfix(token.LBRACKET, p.parseArrayIndexExpression)
	p.registerInfix(token.INCR, p.parsePostfixExpression)
	p.registerInfix(token.DECR, p.parsePostfixExpression)

	return p
}

func (p *Parser) ParseProgram() *ast.Program {
	program := &ast.Program{}
	program.Statements = []ast.Statement{}

	for !p.curTokenIs(token.EOF) {
		stmt := p.parseStatement()
		if stmt != nil {
			program.Statements = append(program.Statements, stmt)
		}
		p.nextToken()
	}

	return program
}

func (p *Parser) parseSemicolon() ast.Expression {
	return nil
}

func (p *Parser) Errors() []string {
	return p.errors
}

func (p *Parser) CheckErrors() {
	errors := p.Errors()
	if len(errors) == 0 {
		return
	}

	fmt.Errorf(strings.Join(errors, "\n"))
}
