package parser

import (
	"fmt"
	"github.com/goby-lang/goby/compiler/ast"
	"github.com/goby-lang/goby/compiler/parser/arguments"
	"github.com/goby-lang/goby/compiler/parser/errors"
	"github.com/goby-lang/goby/compiler/parser/events"
	"github.com/goby-lang/goby/compiler/parser/precedence"
	"github.com/goby-lang/goby/compiler/token"
)

func (p *Parser) parseStatement() ast.Statement {
	switch p.curToken.Type {
	case token.Return:
		return p.parseReturnStatement()
	case token.Def:
		return p.parseDefMethodStatement()
	case token.Comment:
		return nil
	case token.While:
		return p.parseWhileStatement()
	case token.Class:
		return p.parseClassStatement()
	case token.Module:
		return p.parseModuleStatement()
	case token.Next:
		return &ast.NextStatement{BaseNode: &ast.BaseNode{Token: p.curToken}}
	case token.Break:
		return &ast.BreakStatement{BaseNode: &ast.BaseNode{Token: p.curToken}}
	default:
		exp := p.parseExpressionStatement()

		// If parseExpressionStatement got error exp.Expression would be nil
		if exp.Expression != nil {
			// In REPL mode everything should return a value.
			if p.Mode == REPLMode {
				exp.Expression.MarkAsExp()
			} else {
				exp.Expression.MarkAsStmt()
			}
		}

		return exp
	}
}

func (p *Parser) parseDefMethodStatement() *ast.DefStatement {
	var params []ast.Expression
	stmt := &ast.DefStatement{BaseNode: &ast.BaseNode{Token: p.curToken}}

	p.nextToken()

	// Method has specific receiver like `def self.foo` or `def bar.foo`
	if p.peekTokenIs(token.Dot) {
		switch p.curToken.Type {
		case token.Ident:
			stmt.Receiver = p.parseIdentifier()
		case token.InstanceVariable:
			stmt.Receiver = p.parseInstanceVariable()
		case token.Constant:
			stmt.Receiver = p.parseConstant()
		case token.Self:
			stmt.Receiver = &ast.SelfExpression{BaseNode: &ast.BaseNode{Token: p.curToken}}
		default:
			msg := fmt.Sprintf("Invalid method receiver: %s. Line: %d", p.curToken.Literal, p.curToken.Line)
			p.error = errors.InitError(msg, errors.MethodDefinitionError)
		}

		p.nextToken() // .
		if !p.expectPeek(token.Ident) {
			return nil
		}
	}

	stmt.Name = &ast.Identifier{BaseNode: &ast.BaseNode{Token: p.curToken}, Value: p.curToken.Literal}

	// Setter method def foo=()
	if p.peekTokenIs(token.Assign) {
		stmt.Name.Value = stmt.Name.Value + "="
		p.nextToken()
	}

	if p.peekTokenIs(token.Ident) && p.peekTokenAtSameLine() { // def foo x, next token is x and at same line
		msg := fmt.Sprintf("Please add parentheses around method \"%s\"'s parameters. Line: %d", stmt.Name.Value, p.curToken.Line)
		p.error = errors.InitError(msg, errors.MethodDefinitionError)
	}

	if p.peekTokenIs(token.LParen) {
		p.nextToken()

		switch p.peekToken.Type {
		case token.RParen:
			params = []ast.Expression{}
		default:
			params = p.parseParameters()
		}

		if !p.expectPeek(token.RParen) {
			return nil
		}
	} else {
		params = []ast.Expression{}
	}

	stmt.Parameters = params
	stmt.BlockStatement = p.parseBlockStatement(token.End)
	stmt.BlockStatement.KeepLastValue()

	return stmt
}

func (p *Parser) parseParameters() []ast.Expression {
	p.fsm.Event(events.ParseMethodParam)
	params := []ast.Expression{}

	p.nextToken()
	param := p.parseExpression(precedence.Normal)
	params = append(params, param)

	for p.peekTokenIs(token.Comma) {
		p.nextToken()
		p.nextToken()

		if p.curTokenIs(token.Asterisk) && !p.peekTokenIs(token.Ident) {
			p.expectPeek(token.Ident)
			break
		}

		param := p.parseExpression(precedence.Normal)
		params = append(params, param)
	}

	p.checkMethodParameters(params)

	p.fsm.Event(events.BackToNormal)
	return params
}

func (p *Parser) checkMethodParameters(params []ast.Expression) {

	/*
		0 means previous arg is normal argument
		1 means previous arg is optioned argument
		2 means previous arg is keyword argument
		3 means previous arg is splat argument
	*/
	argState := arguments.NormalArg

	checkedParams := []ast.Expression{}

	for _, param := range params {
		switch exp := param.(type) {
		case *ast.Identifier:
			switch argState {
			case arguments.OptionedArg:
				p.error = errors.NewArgumentError(arguments.NormalArg, arguments.OptionedArg, exp.Value, p.curToken.Line)
			case arguments.RequiredKeywordArg:
				p.error = errors.NewArgumentError(arguments.NormalArg, arguments.RequiredKeywordArg, exp.Value, p.curToken.Line)
			case arguments.OptionalKeywordArg:
				p.error = errors.NewArgumentError(arguments.NormalArg, arguments.OptionalKeywordArg, exp.Value, p.curToken.Line)
			case arguments.SplatArg:
				p.error = errors.NewArgumentError(arguments.NormalArg, arguments.SplatArg, exp.Value, p.curToken.Line)
			}
		case *ast.AssignExpression:
			switch argState {
			case arguments.RequiredKeywordArg:
				p.error = errors.NewArgumentError(arguments.OptionedArg, arguments.RequiredKeywordArg, exp.String(), p.curToken.Line)
			case arguments.OptionalKeywordArg:
				p.error = errors.NewArgumentError(arguments.OptionedArg, arguments.OptionalKeywordArg, exp.String(), p.curToken.Line)
			case arguments.SplatArg:
				p.error = errors.NewArgumentError(arguments.OptionedArg, arguments.SplatArg, exp.String(), p.curToken.Line)
			}
			argState = arguments.OptionedArg
		case *ast.ArgumentPairExpression:
			if exp.Value == nil {
				switch argState {
				case arguments.OptionalKeywordArg:
					p.error = errors.NewArgumentError(arguments.RequiredKeywordArg, arguments.OptionalKeywordArg, exp.String(), p.curToken.Line)
				case arguments.SplatArg:
					p.error = errors.NewArgumentError(arguments.RequiredKeywordArg, arguments.SplatArg, exp.String(), p.curToken.Line)
				}

				argState = arguments.RequiredKeywordArg
			} else {
				switch argState {
				case arguments.SplatArg:
					p.error = errors.NewArgumentError(arguments.OptionalKeywordArg, arguments.SplatArg, exp.String(), p.curToken.Line)
				}

				argState = arguments.OptionalKeywordArg
			}
		case *ast.PrefixExpression:
			switch argState {
			case arguments.SplatArg:
				msg := fmt.Sprintf("Can't define splat argument more than once. Line: %d", p.curToken.Line)
				p.error = errors.InitError(msg, errors.ArgumentError)
			}
			argState = arguments.SplatArg
		}

		if p.error != nil {
			break
		}

		if paramDuplicated(checkedParams, param) {
			msg := fmt.Sprintf("Duplicate argument name: \"%s\". Line: %d", getArgName(param), p.curToken.Line)
			p.error = errors.InitError(msg, errors.ArgumentError)
		} else {
			checkedParams = append(checkedParams, param)
		}
	}
}

func (p *Parser) parseClassStatement() *ast.ClassStatement {
	stmt := &ast.ClassStatement{BaseNode: &ast.BaseNode{Token: p.curToken}}

	if !p.expectPeek(token.Constant) {
		return nil
	}

	stmt.Name = &ast.Constant{BaseNode: &ast.BaseNode{Token: p.curToken}, Value: p.curToken.Literal}

	// See if there is any inheritance
	if p.peekTokenIs(token.LT) {
		p.nextToken() // <
		p.nextToken() // Inherited class like 'Bar'
		stmt.SuperClass = p.parseExpression(precedence.Normal)

		switch exp := stmt.SuperClass.(type) {
		case *ast.InfixExpression:
			stmt.SuperClassName = exp.Right.(*ast.Constant).Value
		case *ast.Constant:
			stmt.SuperClassName = exp.Value
		}
	}

	stmt.Body = p.parseBlockStatement(token.End)

	return stmt
}

func (p *Parser) parseModuleStatement() *ast.ModuleStatement {
	stmt := &ast.ModuleStatement{BaseNode: &ast.BaseNode{Token: p.curToken}}

	if !p.expectPeek(token.Constant) {
		return nil
	}

	stmt.Name = &ast.Constant{BaseNode: &ast.BaseNode{Token: p.curToken}, Value: p.curToken.Literal}
	stmt.Body = p.parseBlockStatement(token.End)

	return stmt
}

func (p *Parser) parseReturnStatement() *ast.ReturnStatement {
	stmt := &ast.ReturnStatement{BaseNode: &ast.BaseNode{Token: p.curToken}}

	if !p.peekTokenAtSameLine() {
		null := &ast.NilExpression{BaseNode: &ast.BaseNode{Token: p.curToken}}
		stmt.ReturnValue = null
		return stmt
	}

	p.nextToken()

	stmt.ReturnValue = p.parseExpression(precedence.Normal)

	return stmt
}

func (p *Parser) parseExpressionStatement() *ast.ExpressionStatement {
	stmt := &ast.ExpressionStatement{BaseNode: &ast.BaseNode{Token: p.curToken}}
	if p.curTokenIs(token.Ident) || p.curTokenIs(token.InstanceVariable) {
		// This is used for identifying method call without parens
		// Or multiple variable assignment
		stmt.Expression = p.parseExpression(precedence.Lowest)
	} else {
		stmt.Expression = p.parseExpression(precedence.Normal)
	}

	return stmt
}

func (p *Parser) parseBlockStatement(endTokens ...token.Type) *ast.BlockStatement {

	// curToken is '{'
	bs := &ast.BlockStatement{BaseNode: &ast.BaseNode{Token: p.curToken}}
	bs.Statements = []ast.Statement{}

	p.nextToken()

	if p.curTokenIs(token.Semicolon) {
		p.nextToken()
	}

ParseBlockLoop:
	for {
		for _, t := range endTokens {
			if p.curTokenIs(t) {
				break ParseBlockLoop
			}
		}

		if p.curTokenIs(token.EOF) {
			p.error = errors.InitError("Unexpected EOF", errors.EndOfFileError)
			return bs
		}
		stmt := p.parseStatement()

		if stmt != nil {
			bs.Statements = append(bs.Statements, stmt)
		}
		p.nextToken()
	}

	return bs
}

func (p *Parser) parseWhileStatement() *ast.WhileStatement {
	ws := &ast.WhileStatement{BaseNode: &ast.BaseNode{Token: p.curToken}}

	p.nextToken()
	// Prevent expression's method call to consume while's block as argument.
	p.acceptBlock = false

	oldState := p.fsm.Current()
	p.fsm.Event(events.ParseFuncCall)

	ws.Condition = p.parseExpression(precedence.Normal)

	event, _ := events.EventTable[oldState]
	p.fsm.Event(event)
	p.acceptBlock = true
	p.expectPeek(token.Do)

	if p.curTokenIs(token.Semicolon) {
		p.nextToken()
	}

	ws.Body = p.parseBlockStatement(token.End)

	return ws
}

func paramDuplicated(params []ast.Expression, param ast.Expression) bool {
	for _, p := range params {
		if getArgName(param) == getArgName(p) {
			return true
		}
	}
	return false
}

func getArgName(exp ast.Expression) string {
	assignExp, ok := exp.(*ast.AssignExpression)

	if ok {
		return assignExp.Variables[0].TokenLiteral()
	}

	switch exp := exp.(type) {
	case *ast.ArgumentPairExpression:
		return exp.Key.(*ast.Identifier).Value
	}

	return exp.TokenLiteral()
}
