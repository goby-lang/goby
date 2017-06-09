package parser

import (
	"fmt"
	"strings"

	"github.com/goby-lang/goby/ast"
	"github.com/goby-lang/goby/lexer"
	"github.com/goby-lang/goby/token"
)

// Parser represents lexical analyzer struct
type Parser struct {
	l      *lexer.Lexer
	errors []string

	curToken  token.Token
	peekToken token.Token

	prefixParseFns map[token.Type]prefixParseFn
	infixParseFns  map[token.Type]infixParseFn

	// Determine if call expression should accept block argument,
	// currently only used when parsing while statement.
	// However, this is not a very good practice should change it in the future.
	acceptBlock bool
}

// BuildAST tokenizes and parses given file to build AST
func BuildAST(file []byte) *ast.Program {
	input := string(file)
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	p.CheckErrors()

	return program
}

// New initializes a parser and returns it
func New(l *lexer.Lexer) *Parser {
	p := &Parser{
		l:           l,
		errors:      []string{},
		acceptBlock: true,
	}

	// Read two tokens, so curToken and peekToken are both set.
	p.nextToken()
	p.nextToken()

	p.prefixParseFns = make(map[token.Type]prefixParseFn)
	p.registerPrefix(token.Ident, p.parseIdentifier)
	p.registerPrefix(token.Constant, p.parseConstant)
	p.registerPrefix(token.InstanceVariable, p.parseInstanceVariable)
	p.registerPrefix(token.Int, p.parseIntegerLiteral)
	p.registerPrefix(token.String, p.parseStringLiteral)
	p.registerPrefix(token.True, p.parseBooleanLiteral)
	p.registerPrefix(token.False, p.parseBooleanLiteral)
	p.registerPrefix(token.Null, p.parseNilExpression)
	p.registerPrefix(token.Minus, p.parsePrefixExpression)
	p.registerPrefix(token.Bang, p.parsePrefixExpression)
	p.registerPrefix(token.LParen, p.parseGroupedExpression)
	p.registerPrefix(token.If, p.parseIfExpression)
	p.registerPrefix(token.Self, p.parseSelfExpression)
	p.registerPrefix(token.LBracket, p.parseArrayExpression)
	p.registerPrefix(token.LBrace, p.parseHashExpression)
	p.registerPrefix(token.Semicolon, p.parseSemicolon)
	p.registerPrefix(token.Yield, p.parseYieldExpression)

	p.infixParseFns = make(map[token.Type]infixParseFn)
	p.registerInfix(token.Plus, p.parseInfixExpression)
	p.registerInfix(token.Modulo, p.parseInfixExpression)
	p.registerInfix(token.Minus, p.parseInfixExpression)
	p.registerInfix(token.Slash, p.parseInfixExpression)
	p.registerInfix(token.Eq, p.parseInfixExpression)
	p.registerInfix(token.Asterisk, p.parseInfixExpression)
	p.registerInfix(token.Pow, p.parseInfixExpression)
	p.registerInfix(token.NotEq, p.parseInfixExpression)
	p.registerInfix(token.LT, p.parseInfixExpression)
	p.registerInfix(token.LTE, p.parseInfixExpression)
	p.registerInfix(token.GT, p.parseInfixExpression)
	p.registerInfix(token.GTE, p.parseInfixExpression)
	p.registerInfix(token.COMP, p.parseInfixExpression)
	p.registerInfix(token.Dot, p.parseCallExpression)
	p.registerInfix(token.LParen, p.parseCallExpression)
	p.registerInfix(token.Ident, p.parseCallExpression)
	//p.registerInfix(token.Int, p.parseCallExpression)

	p.registerInfix(token.LBracket, p.parseArrayIndexExpression)
	p.registerInfix(token.Incr, p.parsePostfixExpression)
	p.registerInfix(token.Decr, p.parsePostfixExpression)
	p.registerInfix(token.And, p.parseInfixExpression)
	p.registerInfix(token.Or, p.parseInfixExpression)
	p.registerInfix(token.ResolutionOperator, p.parseInfixExpression)
	p.registerInfix(token.Assign, p.parseInfixExpression)

	return p
}

// ParseProgram update program statements and return program
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

// Errors return parser errors
func (p *Parser) Errors() []string {
	return p.errors
}

// CheckErrors is checking for parser's errors existence
func (p *Parser) CheckErrors() {
	errors := p.Errors()
	if len(errors) == 0 {
		return
	}

	panic(fmt.Sprintf(strings.Join(errors, "\n")))
}
