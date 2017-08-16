package bytecode

import (
	"fmt"
	"github.com/goby-lang/goby/compiler/ast"
)

func (g *Generator) compileExpression(is *InstructionSet, exp ast.Expression, scope *scope, table *localTable) {
	sourceLine := exp.Line()
	switch exp := exp.(type) {
	case *ast.Constant:
		is.define(GetConstant, sourceLine, exp.Value, fmt.Sprint(exp.IsNamespace))
	case *ast.InstanceVariable:
		is.define(GetInstanceVariable, sourceLine, exp.Value)
	case *ast.IntegerLiteral:
		is.define(PutObject, sourceLine, fmt.Sprint(exp.Value))
	case *ast.StringLiteral:
		is.define(PutString, sourceLine, exp.Value)
	case *ast.BooleanExpression:
		is.define(PutObject, sourceLine, fmt.Sprint(exp.Value))
	case *ast.NilExpression:
		is.define(PutNull, sourceLine)
	case *ast.RangeExpression:
		g.compileExpression(is, exp.Start, scope, table)
		g.compileExpression(is, exp.End, scope, table)
		is.define(NewRange, sourceLine, 0)
	case *ast.ArrayExpression:
		for _, elem := range exp.Elements {
			g.compileExpression(is, elem, scope, table)
		}
		is.define(NewArray, sourceLine, len(exp.Elements))
	case *ast.HashExpression:
		for key, value := range exp.Data {
			is.define(PutString, sourceLine, key)
			g.compileExpression(is, value, scope, table)
		}
		is.define(NewHash, sourceLine, len(exp.Data)*2)
	case *ast.SelfExpression:
		is.define(PutSelf, sourceLine)
	case *ast.PrefixExpression:
		g.compilePrefixExpression(is, exp, scope, table)
	case *ast.InfixExpression:
		g.compileInfixExpression(is, exp, scope, table)
	case *ast.Identifier:
		g.compileIdentifier(is, exp, scope, table)
	case *ast.AssignExpression:
		g.compileAssignExpression(is, exp, scope, table)
	case *ast.IfExpression:
		g.compileIfExpression(is, exp, scope, table)
	case *ast.YieldExpression:
		g.compileYieldExpression(is, exp, scope, table)
	case *ast.CallExpression:
		g.compileCallExpression(is, exp, scope, table)
	}
}

func (g *Generator) compileIdentifier(is *InstructionSet, exp *ast.Identifier, scope *scope, table *localTable) {
	index, depth, ok := table.getLCL(exp.Value, table.depth)

	if ok {
		is.define(GetLocal, exp.Line(), depth, index)
		return
	}

	// otherwise it's a method call
	is.define(PutSelf, exp.Line())
	is.define(Send, exp.Line(), exp.Value, 0)
}

func (g *Generator) compileYieldExpression(is *InstructionSet, exp *ast.YieldExpression, scope *scope, table *localTable) {
	is.define(PutSelf, exp.Line())

	for _, arg := range exp.Arguments {
		g.compileExpression(is, arg, scope, table)
	}

	is.define(InvokeBlock, exp.Line(), len(exp.Arguments))
}

func (g *Generator) compileCallExpression(is *InstructionSet, exp *ast.CallExpression, scope *scope, table *localTable) {
	g.compileExpression(is, exp.Receiver, scope, table)

	for _, arg := range exp.Arguments {
		g.compileExpression(is, arg, scope, table)
	}

	if exp.Block != nil {
		// Inside block should be one level deeper than outside
		newTable := newLocalTable(table.depth + 1)
		newTable.upper = table
		blockIndex := g.blockCounter
		g.blockCounter++
		g.compileBlockArgExpression(blockIndex, exp, scope, newTable)
		is.define(Send, exp.Line(), exp.Method, len(exp.Arguments), fmt.Sprintf("block:%d", blockIndex))
		return
	}
	is.define(Send, exp.Line(), exp.Method, len(exp.Arguments))
}

func (g *Generator) compileAssignExpression(is *InstructionSet, exp *ast.AssignExpression, scope *scope, table *localTable) {
	g.compileExpression(is, exp.Value, scope, table)

	if len(exp.Variables) > 1 {
		is.define(ExpandArray, exp.Line(), len(exp.Variables))
	}

	for i, v := range exp.Variables {
		switch name := v.(type) {
		case *ast.Identifier:
			index, depth := table.setLCL(name.Value, table.depth)

			if exp.Optioned != 0 {
				is.define(SetLocal, exp.Line(), depth, index, exp.Optioned)
				return
			}

			is.define(SetLocal, exp.Line(), depth, index)
		case *ast.InstanceVariable:
			is.define(SetInstanceVariable, exp.Line(), name.Value)
		case *ast.Constant:
			is.define(SetConstant, exp.Line(), name.Value)
		}

		/*
			Keep last value so we can have value to pop

			```ruby
			a, b = [1, 2]

			Here we only pop '2', and the statement compilation will add another pop to pop '1'
		*/
		if exp.IsStmt() && !g.REPL && i != len(exp.Variables)-1 {
			is.define(Pop, exp.Line())
		}
	}
}

func (g *Generator) compileBlockArgExpression(index int, exp *ast.CallExpression, scope *scope, table *localTable) {
	is := &InstructionSet{}
	is.name = fmt.Sprint(index)
	is.isType = Block

	for i := 0; i < len(exp.BlockArguments); i++ {
		table.set(exp.BlockArguments[i].Value)
	}

	g.compileCodeBlock(is, exp.Block, scope, table)
	g.endInstructions(is, exp.Line())
	g.instructionSets = append(g.instructionSets, is)
}

func (g *Generator) compileIfExpression(is *InstructionSet, exp *ast.IfExpression, scope *scope, table *localTable) {
	anchorLast := &anchor{}

	for _, c := range exp.Conditionals {
		anchorConditional := &anchor{}

		g.compileExpression(is, c.Condition, scope, table)
		is.define(BranchUnless, exp.Line(), anchorConditional)

		g.compileCodeBlock(is, c.Consequence, scope, table)
		anchorConditional.line = is.count + 1
		is.define(Jump, exp.Line(), anchorLast)
	}

	if exp.Alternative == nil {
		// jump over the `putnil` in false case
		anchorLast.line = is.count + 1
		is.define(PutNull, exp.Line())

		return
	}

	g.compileCodeBlock(is, exp.Alternative, scope, table)

	anchorLast.line = is.count
}

func (g *Generator) compilePrefixExpression(is *InstructionSet, exp *ast.PrefixExpression, scope *scope, table *localTable) {
	switch exp.Operator {
	case "!":
		g.compileExpression(is, exp.Right, scope, table)
		is.define(Send, exp.Line(), exp.Operator, 0)
	case "-":
		is.define(PutObject, exp.Line(), 0)
		g.compileExpression(is, exp.Right, scope, table)
		is.define(Send, exp.Line(), exp.Operator, 1)
	}
}

func (g *Generator) compileInfixExpression(is *InstructionSet, node *ast.InfixExpression, scope *scope, table *localTable) {
	g.compileExpression(is, node.Left, scope, table)
	g.compileExpression(is, node.Right, scope, table)

	if node.Operator != "::" {
		is.define(Send, node.Line(), node.Operator, "1")
	}
}
