package bytecode

import (
	"fmt"

	"github.com/goby-lang/goby/ast"
)

func (g *Generator) compileExpression(is *instructionSet, exp ast.Expression, scope *scope, table *localTable) {
	switch exp := exp.(type) {
	case *ast.Identifier:
		g.compileIdentifier(is, exp, scope, table)
	case *ast.Constant:
		is.define(GetConstant, exp.Value)
	case *ast.InstanceVariable:
		is.define(GetInstanceVariable, exp.Value)
	case *ast.IntegerLiteral:
		is.define(PutObject, fmt.Sprint(exp.Value))
	case *ast.StringLiteral:
		is.define(PutString, fmt.Sprintf("\"%s\"", exp.Value))
	case *ast.BooleanExpression:
		is.define(PutObject, fmt.Sprint(exp.Value))
	case *ast.NilExpression:
		is.define(PutNull)
	case *ast.RangeExpression:
		g.compileExpression(is, exp.Start, scope, table)
		g.compileExpression(is, exp.End, scope, table)
		is.define(NewRange, 0)
	case *ast.ArrayExpression:
		for _, elem := range exp.Elements {
			g.compileExpression(is, elem, scope, table)
		}
		is.define(NewArray, len(exp.Elements))
	case *ast.HashExpression:
		for key, value := range exp.Data {
			is.define(PutString, fmt.Sprintf("\"%s\"", key))
			g.compileExpression(is, value, scope, table)
		}
		is.define(NewHash, len(exp.Data)*2)
	case *ast.InfixExpression:
		if exp.Operator == "=" {
			g.compileAssignExpression(is, exp, scope, table)
			return
		}
		g.compileInfixExpression(is, exp, scope, table)
	case *ast.PrefixExpression:
		g.compilePrefixExpression(is, exp, scope, table)
	case *ast.IfExpression:
		g.compileIfExpression(is, exp, scope, table)
	case *ast.SelfExpression:
		is.define(PutSelf)
	case *ast.YieldExpression:
		g.compileYieldExpression(is, exp, scope, table)
	case *ast.CallExpression:
		g.compileCallExpression(is, exp, scope, table)
	case *ast.RegexLiteral:
		is.define(PutString, fmt.Sprintf("/%s/", exp.Value))
	}
}

func (g *Generator) compileIdentifier(is *instructionSet, exp *ast.Identifier, scope *scope, table *localTable) {
	index, depth, ok := table.getLCL(exp.Value, table.depth)

	// it's local variable
	if ok {
		is.define(GetLocal, depth, index)
		return
	}

	// otherwise it's a method call
	is.define(PutSelf)
	is.define(Send, exp.Value, 0)
}

func (g *Generator) compileYieldExpression(is *instructionSet, exp *ast.YieldExpression, scope *scope, table *localTable) {
	is.define(PutSelf)

	for _, arg := range exp.Arguments {
		g.compileExpression(is, arg, scope, table)
	}

	is.define(InvokeBlock, len(exp.Arguments))
}

func (g *Generator) compileCallExpression(is *instructionSet, exp *ast.CallExpression, scope *scope, table *localTable) {
	g.compileExpression(is, exp.Receiver, scope, table)

	for _, arg := range exp.Arguments {
		g.compileExpression(is, arg, scope, table)
	}

	if exp.Block != nil {
		newTable := newLocalTable(table.depth + 1)
		newTable.upper = table
		blockIndex := g.blockCounter
		g.blockCounter++
		g.compileBlockArgExpression(blockIndex, exp, scope, newTable)
		is.define(Send, exp.Method, len(exp.Arguments), fmt.Sprintf("block:%d", blockIndex))
		return
	}
	is.define(Send, exp.Method, len(exp.Arguments))
}

func (g *Generator) compileAssignExpression(is *instructionSet, exp *ast.InfixExpression, scope *scope, table *localTable) {
	g.compileExpression(is, exp.Right, scope, table)

	switch name := exp.Left.(type) {
	case *ast.Identifier:
		index, depth := table.setLCL(name.Value, table.depth)
		is.define(SetLocal, depth, index)
	case *ast.InstanceVariable:
		is.define(SetInstanceVariable, name.Value)
	case *ast.Constant:
		is.define(SetConstant, name.Value)
	}
}

func (g *Generator) compileBlockArgExpression(index int, exp *ast.CallExpression, scope *scope, table *localTable) {
	is := &instructionSet{}
	is.setLabel(fmt.Sprintf("%s:%d", Block, index))

	for i := 0; i < len(exp.BlockArguments); i++ {
		table.set(exp.BlockArguments[i].Value)
	}

	g.compileCodeBlock(is, exp.Block, scope, table)
	g.endInstructions(is)
	g.instructionSets = append(g.instructionSets, is)
}

func (g *Generator) compileIfExpression(is *instructionSet, exp *ast.IfExpression, scope *scope, table *localTable) {
	g.compileExpression(is, exp.Condition, scope, table)

	anchor1 := &anchor{}
	is.define(BranchUnless, anchor1)

	g.compileCodeBlock(is, exp.Consequence, scope, table)

	anchor1.line = is.Count + 1

	if exp.Alternative == nil {
		anchor1.line--
		is.define(PutNull)
		return
	}

	anchor2 := &anchor{}
	is.define(Jump, anchor2)

	g.compileCodeBlock(is, exp.Alternative, scope, table)

	anchor2.line = is.Count
}

func (g *Generator) compilePrefixExpression(is *instructionSet, exp *ast.PrefixExpression, scope *scope, table *localTable) {
	switch exp.Operator {
	case "!":
		g.compileExpression(is, exp.Right, scope, table)
		is.define(Send, exp.Operator, 0)
	case "-":
		is.define(PutObject, 0)
		g.compileExpression(is, exp.Right, scope, table)
		is.define(Send, exp.Operator, 1)
	}
}

func (g *Generator) compileInfixExpression(is *instructionSet, node *ast.InfixExpression, scope *scope, table *localTable) {
	g.compileExpression(is, node.Left, scope, table)
	g.compileExpression(is, node.Right, scope, table)

	if node.Operator != "::" {
		is.define(Send, node.Operator, "1")
	}
}
