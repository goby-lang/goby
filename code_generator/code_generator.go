package code_generator

import (
	"github.com/st0012/Rooby/ast"
	"strings"
	"fmt"
	"regexp"
)

func GenerateByteCode(program *ast.Program) string {
	return strings.TrimSpace(removeEmptyLine(compile(program)))
}

func compileProgram(stmts []ast.Statement) string {
	var result []string

	for _, statement := range stmts {
		result = append(result, compile(statement))
	}

	return strings.Join(result, "\n")
}

func compile(node ast.Node) string {
	switch node := node.(type) {

	// Statements
	case *ast.Program:
		return compileProgram(node.Statements)
	case *ast.ExpressionStatement:
		return compile(node.Expression)
	case *ast.IntegerLiteral:
		return fmt.Sprintf(`putobject %d`, node.Value)
	case *ast.InfixExpression:
		var operation string
		left := compile(node.Left)
		right := compile(node.Right)
		switch node.Operator {
		case "+":
			operation = "opt_plus"
		case "-":
			operation = "opt_minus"
		case "*":
			operation = "opt_mult"
		case "/":
			operation = "opt_div"
		default:
			panic(fmt.Sprintf("Doesn't support %s operator", node.Operator))
		}
		return fmt.Sprintf(`
%s
%s
%s
`, left, right, operation)
	default:
		return ""
	}
	//case *ast.BlockStatement:
	//	return evalBlockStatements(node.Statements)
	//case *ast.ReturnStatement:
	//	val := Eval(node.ReturnValue, scope)
	//	if isError(val) {
	//		return val
	//	}
	//	return &ReturnValue{Value: val}
	//case *ast.AssignStatement:
	//	return evalAssignStatement(node, scope)
	//case *ast.ClassStatement:
	//	return evalClassStatement(node, scope)
	//case *ast.Identifier:
	//	return evalIdentifier(node, scope)
	//case *ast.Constant:
	//	return evalConstant(node, scope)
	//case *ast.InstanceVariable:
	//	return evalInstanceVariable(node, scope)
	//case *ast.DefStatement:
	//	return evalDefStatement(node, scope)
	//case *ast.WhileStatement:
	//	return evalWhileStatement(node, scope)
	//case *ast.YieldStatement:
	//	// Only retrieve current scope's block.
	//	return evalYieldStatement(node, scope)
	//
	//// Expressions
	//case *ast.IfExpression:
	//	return evalIfExpression(node, scope)
	//case *ast.CallExpression:
	//	var block *Method
	//	receiver := Eval(node.Receiver, scope)
	//	args := evalArgs(node.Arguments, scope)
	//
	//	if node.Block != nil {
	//		block = &Method{Body: node.Block, Parameters: node.BlockArguments, Scope: scope}
	//	}
	//
	//	return sendMethodCall(receiver, node.Method, args, block)
	//
	//case *ast.PrefixExpression:
	//	val := Eval(node.Right, scope)
	//	if isError(val) {
	//		return val
	//	}
	//	return evalPrefixExpression(node.Operator, val)
	//case *ast.InfixExpression:
	//	valLeft := Eval(node.Left, scope)
	//	if isError(valLeft) {
	//		return valLeft
	//	}
	//
	//	valRight := Eval(node.Right, scope)
	//	if isError(valRight) {
	//		return valRight
	//	}
	//
	//	return evalInfixExpression(valLeft, node.Operator, valRight)
	//case *ast.SelfExpression:
	//	return scope.Self
	//case *ast.IntegerLiteral:
	//	return InitilaizeInteger(node.Value)
	//case *ast.StringLiteral:
	//	return InitializeString(node.Value)
	//case *ast.Boolean:
	//	if node.Value {
	//		return TRUE
	//	}
	//	return FALSE
	//case *ast.ArrayExpression:
	//	elements := []Object{}
	//
	//	for _, exp := range node.Elements {
	//		elements = append(elements, Eval(exp, scope))
	//	}
	//
	//	arr := InitializeArray(elements)
	//	return arr
	//case *ast.HashExpression:
	//	pairs := map[string]Object{}
	//
	//	for key, value := range node.Data {
	//		pairs[key] = Eval(value, scope)
	//	}
	//
	//	hash := InitializeHash(pairs)
	//	return hash
	//}
}


func removeEmptyLine(s string) string {
	regex, err := regexp.Compile("\n\n")
	if err != nil {
		panic(err)
	}
	s = regex.ReplaceAllString(s, "\n")

	return s
}