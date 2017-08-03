package bytecode

import (
	"fmt"
	"github.com/goby-lang/goby/compiler/ast"
)

/*
	These constants are enums that represent argument's types
*/
const (
	NormalArg int = iota
	OptionedArg
)

func (g *Generator) compileStatements(stmts []ast.Statement, scope *scope, table *localTable) {
	is := &InstructionSet{isType: Program, name: Program}

	for _, statement := range stmts {
		g.compileStatement(is, statement, scope, table)

		if !g.REPL {
			expStmt, ok := statement.(*ast.ExpressionStatement)

			if ok && expStmt.Expression.IsExp() {
				continue
			}

			is.define(Pop)
		}
	}

	g.endInstructions(is)
	g.instructionSets = append(g.instructionSets, is)
}

func (g *Generator) compileStatement(is *InstructionSet, statement ast.Statement, scope *scope, table *localTable) {
	scope.line++
	switch stmt := statement.(type) {
	case *ast.ExpressionStatement:
		g.compileExpression(is, stmt.Expression, scope, table)

		if stmt.Expression.IsStmt() {
			is.define(Pop)
		}
	case *ast.DefStatement:
		g.compileDefStmt(is, stmt, scope)
	case *ast.ClassStatement:
		g.compileClassStmt(is, stmt, scope, table)
	case *ast.ModuleStatement:
		g.compileModuleStmt(is, stmt, scope)
	case *ast.ReturnStatement:
		g.compileExpression(is, stmt.ReturnValue, scope, table)
		g.endInstructions(is)
	case *ast.WhileStatement:
		g.compileWhileStmt(is, stmt, scope, table)
	case *ast.NextStatement:
		g.compileNextStatement(is, scope)
	case *ast.BreakStatement:
		g.compileBreakStatement(is, scope)
	}
}

func (g *Generator) compileWhileStmt(is *InstructionSet, stmt *ast.WhileStatement, scope *scope, table *localTable) {
	anchor1 := &anchor{}
	breakAnchor := &anchor{}

	is.define(Jump, anchor1)

	is.define(PutNull)
	is.define(Pop)
	is.define(Jump, anchor1)

	anchor2 := &anchor{is.count}

	scope.anchors["next"] = anchor1
	scope.anchors["break"] = breakAnchor

	g.compileCodeBlock(is, stmt.Body, scope, table)

	anchor1.line = is.count

	g.compileExpression(is, stmt.Condition, scope, table)

	is.define(BranchIf, anchor2)
	is.define(PutNull)
	is.define(Pop)

	breakAnchor.line = is.count
}

func (g *Generator) compileNextStatement(is *InstructionSet, scope *scope) {
	is.define(Jump, scope.anchors["next"])
}

func (g *Generator) compileBreakStatement(is *InstructionSet, scope *scope) {
	is.define(Jump, scope.anchors["break"])
}

func (g *Generator) compileClassStmt(is *InstructionSet, stmt *ast.ClassStatement, scope *scope, table *localTable) {
	is.define(PutSelf)

	if stmt.SuperClass != nil {
		g.compileExpression(is, stmt.SuperClass, scope, table)
		is.define(DefClass, "class:"+stmt.Name.Value, stmt.SuperClassName)
	} else {
		is.define(DefClass, "class:"+stmt.Name.Value)
	}

	is.define(Pop)

	scope = newScope(stmt)

	// compile class's content
	newIS := &InstructionSet{}
	newIS.name = stmt.Name.Value
	newIS.isType = ClassDef

	g.compileCodeBlock(newIS, stmt.Body, scope, scope.localTable)
	newIS.define(Leave)
	g.instructionSets = append(g.instructionSets, newIS)
}

func (g *Generator) compileModuleStmt(is *InstructionSet, stmt *ast.ModuleStatement, scope *scope) {
	is.define(PutSelf)
	is.define(DefClass, "module:"+stmt.Name.Value)
	is.define(Pop)

	scope = newScope(stmt)
	newIS := &InstructionSet{}
	newIS.name = stmt.Name.Value
	newIS.isType = ClassDef

	g.compileCodeBlock(newIS, stmt.Body, scope, scope.localTable)
	newIS.define(Leave)
	g.instructionSets = append(g.instructionSets, newIS)
}

func (g *Generator) compileDefStmt(is *InstructionSet, stmt *ast.DefStatement, scope *scope) {
	is.define(PutSelf)
	is.define(PutString, fmt.Sprintf("\"%s\"", stmt.Name.Value))

	switch stmt.Receiver.(type) {
	case *ast.SelfExpression:
		is.define(DefSingletonMethod, len(stmt.Parameters))
	case nil:
		is.define(DefMethod, len(stmt.Parameters))
	}

	scope = newScope(stmt)

	// compile method definition's content
	newIS := &InstructionSet{}
	newIS.name = stmt.Name.Value
	newIS.isType = MethodDef

	for i := 0; i < len(stmt.Parameters); i++ {
		var argType int
		switch exp := stmt.Parameters[i].(type) {
		case *ast.Identifier:
			argType = NormalArg
			scope.localTable.setLCL(exp.Value, scope.localTable.depth)
		case *ast.AssignExpression:
			argType = OptionedArg
			exp.Optioned = 1
			g.compileAssignExpression(newIS, exp, scope, scope.localTable)
		}

		newIS.argTypes = append(newIS.argTypes, argType)
	}

	if len(stmt.BlockStatement.Statements) == 0 {
		newIS.define(PutNull)
	} else {
		g.compileCodeBlock(newIS, stmt.BlockStatement, scope, scope.localTable)
	}

	g.endInstructions(newIS)
	g.instructionSets = append(g.instructionSets, newIS)
}
