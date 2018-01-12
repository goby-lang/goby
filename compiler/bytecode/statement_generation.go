package bytecode

import (
	"github.com/goby-lang/goby/compiler/ast"
)

/*
	These constants are enums that represent argument's types
*/
const (
	NormalArg int = iota
	OptionedArg
	SplatArg
	RequiredKeywordArg
	OptionalKeywordArg
)

func (g *Generator) compileStatements(stmts []ast.Statement, scope *scope, table *localTable) {
	is := &InstructionSet{isType: Program, name: Program}

	for _, statement := range stmts {
		g.compileStatement(is, statement, scope, table)
	}

	// empty input so no statement is given
	if len(stmts) == 0 {
		return
	}

	g.endInstructions(is, stmts[len(stmts)-1].Line())
	g.instructionSets = append(g.instructionSets, is)
}

func (g *Generator) compileStatement(is *InstructionSet, statement ast.Statement, scope *scope, table *localTable) {
	scope.line++
	switch stmt := statement.(type) {
	case *ast.ExpressionStatement:
		if !g.REPL && stmt.Expression.IsStmt() {
			g.compileExpression(is, stmt.Expression, scope, table)
			is.define(Pop, statement.Line())

			return
		}

		g.compileExpression(is, stmt.Expression, scope, table)
	case *ast.DefStatement:
		g.compileDefStmt(is, stmt, scope)
	case *ast.ClassStatement:
		g.compileClassStmt(is, stmt, scope, table)
		/*
			```
			This is for pop 'Bar' in

			```
			class Foo < Bar
			end
			```
		*/
		if stmt.SuperClass != nil {
			is.define(Pop, statement.Line())
		}
	case *ast.ModuleStatement:
		g.compileModuleStmt(is, stmt, scope)
	case *ast.ReturnStatement:
		g.compileExpression(is, stmt.ReturnValue, scope, table)
		g.endInstructions(is, stmt.Line())
	case *ast.WhileStatement:
		g.compileWhileStmt(is, stmt, scope, table)
	case *ast.NextStatement:
		g.compileNextStatement(is, stmt, scope)
	case *ast.BreakStatement:
		g.compileBreakStatement(is, stmt, scope)
	}
}

func (g *Generator) compileWhileStmt(is *InstructionSet, stmt *ast.WhileStatement, scope *scope, table *localTable) {
	anchor1 := &anchor{}
	breakAnchor := &anchor{}

	is.define(Jump, stmt.Line(), anchor1)

	anchor2 := &anchor{is.count}

	scope.anchors["next"] = anchor1
	scope.anchors["break"] = breakAnchor

	g.compileCodeBlock(is, stmt.Body, scope, table)

	anchor1.line = is.count

	g.compileExpression(is, stmt.Condition, scope, table)

	is.define(BranchIf, stmt.Line(), anchor2)

	breakAnchor.line = is.count
}

func (g *Generator) compileNextStatement(is *InstructionSet, stmt ast.Statement, scope *scope) {
	is.define(Jump, stmt.Line(), scope.anchors["next"])
}

func (g *Generator) compileBreakStatement(is *InstructionSet, stmt ast.Statement, scope *scope) {
	is.define(Jump, stmt.Line(), scope.anchors["break"])
}

func (g *Generator) compileClassStmt(is *InstructionSet, stmt *ast.ClassStatement, scope *scope, table *localTable) {
	is.define(PutSelf, stmt.Line())

	if stmt.SuperClass != nil {
		g.compileExpression(is, stmt.SuperClass, scope, table)
		is.define(DefClass, stmt.Line(), "class:"+stmt.Name.Value, stmt.SuperClassName)
	} else {
		is.define(DefClass, stmt.Line(), "class:"+stmt.Name.Value)
	}

	is.define(Pop, stmt.Line())

	scope = newScope(stmt)

	// compile class's content
	newIS := &InstructionSet{}
	newIS.name = stmt.Name.Value
	newIS.isType = ClassDef

	g.compileCodeBlock(newIS, stmt.Body, scope, scope.localTable)
	newIS.define(Leave, stmt.Line())
	g.instructionSets = append(g.instructionSets, newIS)
}

func (g *Generator) compileModuleStmt(is *InstructionSet, stmt *ast.ModuleStatement, scope *scope) {
	is.define(PutSelf, stmt.Line())
	is.define(DefClass, stmt.Line(), "module:"+stmt.Name.Value)
	is.define(Pop, stmt.Line())

	scope = newScope(stmt)
	newIS := &InstructionSet{}
	newIS.name = stmt.Name.Value
	newIS.isType = ClassDef

	g.compileCodeBlock(newIS, stmt.Body, scope, scope.localTable)
	newIS.define(Leave, stmt.Line())
	g.instructionSets = append(g.instructionSets, newIS)
}

func (g *Generator) compileDefStmt(is *InstructionSet, stmt *ast.DefStatement, scope *scope) {
	switch stmt.Receiver.(type) {
	case nil:
		is.define(PutSelf, stmt.Line())
		is.define(PutString, stmt.Line(), stmt.Name.Value)
		is.define(DefMethod, stmt.Line(), len(stmt.Parameters))
	default:
		g.compileExpression(is, stmt.Receiver, scope, scope.localTable)
		is.define(PutString, stmt.Line(), stmt.Name.Value)
		is.define(DefSingletonMethod, stmt.Line(), len(stmt.Parameters))
	}

	scope = newScope(stmt)

	// compile method definition's content
	newIS := &InstructionSet{
		name:   stmt.Name.Value,
		isType: MethodDef,
		argTypes: &ArgSet{
			names: make([]string, len(stmt.Parameters)),
			types: make([]int, len(stmt.Parameters)),
		},
	}

	for i := 0; i < len(stmt.Parameters); i++ {
		switch exp := stmt.Parameters[i].(type) {
		case *ast.Identifier:
			scope.localTable.setLCL(exp.Value, scope.localTable.depth)

			newIS.argTypes.setArg(i, exp.Value, NormalArg)
		case *ast.AssignExpression:
			exp.Optioned = 1

			v := exp.Variables[0]
			varName := v.(*ast.Identifier)
			g.compileAssignExpression(newIS, exp, scope, scope.localTable)

			newIS.argTypes.setArg(i, varName.Value, OptionedArg)
		case *ast.PrefixExpression:
			if exp.Operator != "*" {
				continue
			}
			ident := exp.Right.(*ast.Identifier)
			scope.localTable.setLCL(ident.Value, scope.localTable.depth)

			newIS.argTypes.setArg(i, ident.Value, SplatArg)
		case *ast.PairExpression:
			key := exp.Key.(*ast.Identifier)
			index, depth := scope.localTable.setLCL(key.Value, scope.localTable.depth)

			if exp.Value != nil {
				g.compileExpression(newIS, exp.Value, scope, scope.localTable)
				newIS.define(SetLocal, exp.Line(), depth, index, 1)
				newIS.argTypes.setArg(i, key.Value, OptionalKeywordArg)
			} else {
				newIS.argTypes.setArg(i, key.Value, RequiredKeywordArg)
			}
		}
	}

	if len(stmt.BlockStatement.Statements) == 0 {
		newIS.define(PutNull, stmt.Line())
	} else {
		g.compileCodeBlock(newIS, stmt.BlockStatement, scope, scope.localTable)
	}

	g.endInstructions(newIS, stmt.Line())
	g.instructionSets = append(g.instructionSets, newIS)
}
