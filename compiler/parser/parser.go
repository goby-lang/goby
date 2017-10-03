package parser

import (
	"fmt"

	"github.com/goby-lang/goby/compiler/ast"
	"github.com/goby-lang/goby/compiler/lexer"
	"github.com/goby-lang/goby/compiler/token"
	"github.com/looplab/fsm"
)

const (
	_ int = iota
	// EndOfFileError represents normal EOF error
	EndOfFileError
	// WrongTokenError means that token is not what we expected
	WrongTokenError
	// UnexpectedTokenError means that token is not expected to appear in current condition
	UnexpectedTokenError
	// UnexpectedEndError means we get unexpected "end" keyword (this is mainly created for REPL)
	UnexpectedEndError
	// MethodDefinitionError means there's an error on method definition's method name
	MethodDefinitionError
	// InvalidAssignmentError means user assigns value to wrong type of expressions
	InvalidAssignmentError
	// SyntaxError means there's a grammatical in the source code
	SyntaxError
)

// Error represents parser's parsing error
type Error struct {
	// Message contains the readable message of error
	Message string
	errType int
}

// IsEOF checks if error is end of file error
func (e *Error) IsEOF() bool {
	return e.errType == EndOfFileError
}

// IsUnexpectedEnd checks if error is unexpected "end" keyword error
func (e *Error) IsUnexpectedEnd() bool {
	return e.errType == UnexpectedEndError
}

// Parser represents lexical analyzer struct
type Parser struct {
	Lexer *lexer.Lexer
	error *Error

	curToken  token.Token
	peekToken token.Token

	prefixParseFns map[token.Type]prefixParseFn
	infixParseFns  map[token.Type]infixParseFn

	// Determine if call expression should accept block argument,
	// currently only used when parsing while statement.
	// However, this is not a very good practice should change it in the future.
	acceptBlock bool
	fsm         *fsm.FSM
	Mode        int
}

// These are the enums for marking parser's mode, which decides whether it should pop unused values.
const (
	NormalMode int = iota
	REPLMode
	TestMode
)

// These are state machine's events
const (
	backToNormal     = "backToNormal"
	parseFuncCall    = "parseFuncCall"
	parseMethodParam = "parseMethodParam"
	parseAssignment  = "parseAssignment"
)

// These are state machine's states
const (
	normal             = "normal"
	parsingFuncCall    = "parsingFuncCall"
	parsingMethodParam = "parsingMethodParam"
	parsingAssignment  = "parsingAssignment"
)

var eventTable = map[string]string{
	normal:             backToNormal,
	parsingFuncCall:    parseFuncCall,
	parsingMethodParam: parseMethodParam,
	parsingAssignment:  parseAssignment,
}

// New initializes a parser and returns it
func New(l *lexer.Lexer) *Parser {
	p := &Parser{
		Lexer:       l,
		acceptBlock: true,
	}

	p.fsm = fsm.NewFSM(
		normal,
		fsm.Events{
			{Name: parseFuncCall, Src: []string{normal}, Dst: parsingFuncCall},
			{Name: parseMethodParam, Src: []string{normal, parsingAssignment}, Dst: parsingMethodParam},
			{Name: parseAssignment, Src: []string{normal, parsingFuncCall}, Dst: parsingAssignment},
			{Name: backToNormal, Src: []string{parsingFuncCall, parsingMethodParam, parsingAssignment}, Dst: normal},
		},
		fsm.Callbacks{},
	)

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
	p.registerInfix(token.PlusEq, p.parseAssignExpression)
	p.registerInfix(token.Modulo, p.parseInfixExpression)
	p.registerInfix(token.Minus, p.parseInfixExpression)
	p.registerInfix(token.MinusEq, p.parseAssignExpression)
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
	p.registerInfix(token.Incr, p.parsePostfixExpression)
	p.registerInfix(token.Decr, p.parsePostfixExpression)
	p.registerInfix(token.And, p.parseInfixExpression)
	p.registerInfix(token.Or, p.parseInfixExpression)
	p.registerInfix(token.OrEq, p.parseAssignExpression)
	p.registerInfix(token.Comma, p.parseMultiVariables)
	p.registerInfix(token.ResolutionOperator, p.parseInfixExpression)
	p.registerInfix(token.Assign, p.parseAssignExpression)
	p.registerInfix(token.Range, p.parseRangeExpression)
	p.registerInfix(token.Dot, p.parseCallExpressionWithReceiver)
	p.registerInfix(token.LParen, p.parseCallExpressionWithoutReceiver)
	p.registerInfix(token.LBracket, p.parseIndexExpression)

	return p
}

// ParseProgram update program statements and return program
func (p *Parser) ParseProgram() (*ast.Program, *Error) {
	p.error = nil
	// Read two tokens, so curToken and peekToken are both set.
	p.nextToken()
	p.nextToken()
	program := &ast.Program{}
	program.Statements = []ast.Statement{}

	for !p.curTokenIs(token.EOF) {
		stmt := p.parseStatement()

		if stmt != nil {
			program.Statements = append(program.Statements, stmt)
		} else {
			if p.error != nil {
				return nil, p.error
			}
		}

		p.nextToken()

		if p.error != nil {
			return nil, p.error
		}
	}

	if p.Mode == TestMode {
		stmt := program.Statements[len(program.Statements)-1]
		expStmt, ok := stmt.(*ast.ExpressionStatement)

		if ok {
			expStmt.Expression.MarkAsExp()
		}
	}

	return program, nil
}

func (p *Parser) parseSemicolon() ast.Expression {
	return nil
}

func (p *Parser) peekPrecedence() int {
	if p, ok := precedence[p.peekToken.Type]; ok {
		return p
	}

	return NORMAL
}

func (p *Parser) curPrecedence() int {
	if p, ok := precedence[p.curToken.Type]; ok {
		return p
	}

	return NORMAL
}

func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.Lexer.NextToken()
}

func (p *Parser) curTokenIs(t token.Type) bool {
	return p.curToken.Type == t
}

func (p *Parser) peekTokenIs(t token.Type) bool {
	return p.peekToken.Type == t
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

func (p *Parser) peekError(t token.Type) {
	msg := fmt.Sprintf("expected next token to be %s, got %s instead. Line: %d", t, p.peekToken.Type, p.peekToken.Line)
	p.error = &Error{Message: msg, errType: WrongTokenError}
}

func (p *Parser) noPrefixParseFnError(t token.Type) {
	msg := fmt.Sprintf("unexpected %s Line: %d", p.curToken.Literal, p.curToken.Line)

	if t == token.End {
		p.error = &Error{Message: msg, errType: UnexpectedEndError}
	} else {
		p.error = &Error{Message: msg, errType: UnexpectedTokenError}
	}
}

func (p *Parser) peekTokenAtSameLine() bool {
	return p.curToken.Line == p.peekToken.Line && p.peekToken.Type != token.EOF
}
