package bytecode

import (
	"fmt"
	"github.com/goby-lang/goby/compiler/ast"
)

func (g *Generator) compileExpression(is *InstructionSet, exp ast.Expression, scope *scope, table *localTable) {
	// See fsm initialization's comment
	if g.fsm.Is(keepExp) {
		switch exp := exp.(type) {
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
		case *ast.SelfExpression:
			is.define(PutSelf)
		case *ast.PrefixExpression:
			g.compilePrefixExpression(is, exp, scope, table)
		case *ast.InfixExpression:
			if exp.Operator != "=" {
				g.compileInfixExpression(is, exp, scope, table)
			}
		}
	}

	switch exp := exp.(type) {
	case *ast.Identifier:
		g.compileIdentifier(is, exp, scope, table)
	case *ast.InfixExpression:
		if exp.Operator == "=" {
			// Because this is assignment so we do need the expression's value
			g.fsm.Event(keepExp)
			g.compileAssignExpression(is, exp, scope, table)
		}
	case *ast.IfExpression:
		g.compileIfExpression(is, exp, scope, table)
	case *ast.YieldExpression:
		g.compileYieldExpression(is, exp, scope, table)
	case *ast.CallExpression:
		g.fsm.Event(keepExp)
		g.compileCallExpression(is, exp, scope, table)
	}
}

func (g *Generator) compileIdentifier(is *InstructionSet, exp *ast.Identifier, scope *scope, table *localTable) {
	index, depth, ok := table.getLCL(exp.Value, table.depth)

	// it's local variable
	if ok && g.fsm.Is(keepExp) {
		is.define(GetLocal, depth, index)
		return
	}

	// otherwise it's a method call
	is.define(PutSelf)
	is.define(Send, exp.Value, 0)
}

func (g *Generator) compileYieldExpression(is *InstructionSet, exp *ast.YieldExpression, scope *scope, table *localTable) {
	is.define(PutSelf)

	for _, arg := range exp.Arguments {
		g.compileExpression(is, arg, scope, table)
	}

	is.define(InvokeBlock, len(exp.Arguments))
}

func (g *Generator) compileCallExpression(is *InstructionSet, exp *ast.CallExpression, scope *scope, table *localTable) {
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

	if exp.Method == "++" || exp.Method == "--" {
		// ++ and -- are methods with side effect but shouldn't return anything
		is.define(Pop)
	}
}

func (g *Generator) compileAssignExpression(is *InstructionSet, exp *ast.InfixExpression, scope *scope, table *localTable) {
	g.compileExpression(is, exp.Right, scope, table)

	switch name := exp.Left.(type) {
	case *ast.Identifier:
		index, depth := table.setLCL(name.Value, table.depth)

		if exp.Optioned != 0 {
			is.define(SetLocal, depth, index, exp.Optioned)
			return
		}

		is.define(SetLocal, depth, index)
	case *ast.InstanceVariable:
		is.define(SetInstanceVariable, name.Value)
	case *ast.Constant:
		is.define(SetConstant, name.Value)
	}
}

func (g *Generator) compileBlockArgExpression(index int, exp *ast.CallExpression, scope *scope, table *localTable) {
	is := &InstructionSet{}
	is.name = fmt.Sprint(index)
	is.isType = Block

	for i := 0; i < len(exp.BlockArguments); i++ {
		table.set(exp.BlockArguments[i].Value)
	}

	g.fsm.Event(removeExp)
	g.compileCodeBlock(is, exp.Block, scope, table)
	g.fsm.Event(keepExp)

	g.endInstructions(is)
	g.instructionSets = append(g.instructionSets, is)
}

func (g *Generator) compileIfExpression(is *InstructionSet, exp *ast.IfExpression, scope *scope, table *localTable) {
	g.compileExpression(is, exp.Condition, scope, table)

	anchor1 := &anchor{}
	anchor2 := &anchor{}

	is.define(BranchUnless, anchor1)

	g.fsm.Event(removeExp)
	g.compileCodeBlock(is, exp.Consequence, scope, table)
	g.fsm.Event(keepExp)

	anchor1.line = is.count + 1

	if g.fsm.Is(keepExp) {
		is.define(Jump, anchor2)
	}

	if exp.Alternative == nil {
		if g.fsm.Is(keepExp) {
			// jump over the `putnil` in false case
			anchor2.line = anchor1.line + 1
			is.define(PutNull)
		}

		return
	}

	g.fsm.Event(removeExp)
	g.compileCodeBlock(is, exp.Alternative, scope, table)
	g.fsm.Event(keepExp)

	anchor2.line = is.count
}

func (g *Generator) compilePrefixExpression(is *InstructionSet, exp *ast.PrefixExpression, scope *scope, table *localTable) {
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

func (g *Generator) compileInfixExpression(is *InstructionSet, node *ast.InfixExpression, scope *scope, table *localTable) {
	g.compileExpression(is, node.Left, scope, table)
	g.compileExpression(is, node.Right, scope, table)

	if node.Operator != "::" {
		is.define(Send, node.Operator, "1")
	}
}
