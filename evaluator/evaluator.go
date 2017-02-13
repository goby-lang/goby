package evaluator

import (
	"fmt"
	"github.com/st0012/rooby/ast"
	"github.com/st0012/rooby/object"
)

var (
	TRUE  = &object.Boolean{Value: true}
	FALSE = &object.Boolean{Value: false}
	NULL  = &object.Null{}
)

func Eval(node ast.Node, env *object.Environment) object.Object {
	switch node := node.(type) {

	// Statements
	case *ast.Program:
		return evalProgram(node.Statements, env)
	case *ast.ExpressionStatement:
		return Eval(node.Expression, env)
	case *ast.BlockStatement:
		return evalBlockStatements(node.Statements, env)
	case *ast.ReturnStatement:
		val := Eval(node.ReturnValue, env)
		if isError(val) {
			return val
		}
		return &object.ReturnValue{Value: val}
	case *ast.LetStatement:
		val := Eval(node.Value, env)
		if isError(val) {
			return val
		}
		return env.Set(node.Name.Value, val)
	case *ast.ClassStatement:
		return evalClass(node, env)
	case *ast.Identifier:
		return evalIdentifier(node, env)
	case *ast.Constant:
		return evalConstant(node, env)
	case *ast.DefStatement:
		return evalDefStatement(node, env)

	// Expressions
	case *ast.IfExpression:
		return evalIfExpression(node, env)
	//case *ast.DefStatement:
	//	return &object.Method{Parameters: node.Parameters, Body: node.BlockStatement, Env: env}
	//case *ast.CallExpression:
	//	function := Eval(node.Method, env)
	//	if isError(function) {
	//		return function
	//	}
	//
	//	args := evalArgs(node.Arguments, env)
	//
	//	return applyFunction(function, args)

	case *ast.PrefixExpression:
		val := Eval(node.Right, env)
		if isError(val) {
			return val
		}
		return evalPrefixExpression(node.Operator, val)
	case *ast.InfixExpression:
		valLeft := Eval(node.Left, env)
		if isError(valLeft) {
			return valLeft
		}

		valRight := Eval(node.Right, env)
		if isError(valRight) {
			return valRight
		}

		return evalInfixExpression(valLeft, node.Operator, valRight)
	case *ast.IntegerLiteral:
		return &object.Integer{Value: node.Value}
	case *ast.StringLiteral:
		return &object.String{Value: node.Value}
	case *ast.Boolean:
		if node.Value {
			return TRUE
		}
		return FALSE
	}

	return nil
}

func evalIdentifier(node *ast.Identifier, env *object.Environment) object.Object {
	if val, ok := env.Get(node.Value); ok {
		return val
	}

	return newError("identifier not found: %s", node.Value)
}

func evalConstant(node *ast.Constant, env *object.Environment) object.Object {
	if val, ok := env.Get(node.Value); ok {
		return val
	}

	return newError("constant not found: %s", node.Value)
}

func evalProgram(stmts []ast.Statement, env *object.Environment) object.Object {
	var result object.Object

	for _, statement := range stmts {
		result = Eval(statement, env)

		switch result := result.(type) {
		case *object.ReturnValue:
			return result.Value
		case *object.Error:
			return result
		}
	}

	return result
}

func evalBlockStatements(stmts []ast.Statement, env *object.Environment) object.Object {
	var result object.Object

	for _, statement := range stmts {

		result = Eval(statement, env)

		if result != nil {
			switch result := result.(type) {
			case *object.ReturnValue:
				return result
			case *object.Error:
				return result
			}
		}
	}

	return result
}

func evalPrefixExpression(operator string, right object.Object) object.Object {
	switch operator {
	case "!":
		return evalBangPrefixExpression(right)
	case "-":
		return evalMinusPrefixExpression(right)
	}
	return newError("unknown operator: %s%s", operator, right.Type())
}

func evalBangPrefixExpression(right object.Object) *object.Boolean {
	switch right {
	case FALSE:
		return TRUE
	case NULL:
		return TRUE
	default:
		return FALSE
	}
}

func evalMinusPrefixExpression(right object.Object) object.Object {
	if right.Type() != object.INTEGER_OBJ {
		return newError("unknown operator: %s%s", "-", right.Type())
	}
	value := right.(*object.Integer).Value
	return &object.Integer{Value: -value}
}

func evalInfixExpression(left object.Object, operator string, right object.Object) object.Object {
	switch {
	case left.Type() == object.INTEGER_OBJ && right.Type() == object.INTEGER_OBJ:
		return evalIntegerInfixExpression(left, operator, right)
	case left.Type() == object.BOOLEAN_OBJ && right.Type() == object.BOOLEAN_OBJ:
		return evalBooleanInfixExpression(left, operator, right)
	case left.Type() == object.STRING_OBJ && right.Type() == object.STRING_OBJ:
		return evalStringInfixExpression(left, operator, right)
	default:
		return newError("type mismatch: %s %s %s", left.Type(), operator, right.Type())
	}
}

func evalIntegerInfixExpression(left object.Object, operator string, right object.Object) object.Object {
	leftValue := left.(*object.Integer).Value
	rightValue := right.(*object.Integer).Value

	switch operator {
	case "+":
		return &object.Integer{Value: leftValue + rightValue}
	case "-":
		return &object.Integer{Value: leftValue - rightValue}
	case "*":
		return &object.Integer{Value: leftValue * rightValue}
	case "/":
		return &object.Integer{Value: leftValue / rightValue}
	case ">":
		return &object.Boolean{Value: leftValue > rightValue}
	case "<":
		return &object.Boolean{Value: leftValue < rightValue}
	case "==":
		return &object.Boolean{Value: leftValue == rightValue}
	case "!=":
		return &object.Boolean{Value: leftValue != rightValue}
	default:
		return newError("unknown operator: %s %s %s", left.Type(), operator, right.Type())
	}
}

func evalBooleanInfixExpression(left object.Object, operator string, right object.Object) object.Object {
	leftValue := left.(*object.Boolean).Value
	rightValue := right.(*object.Boolean).Value
	switch operator {
	case "==":
		return &object.Boolean{Value: leftValue == rightValue}
	case "!=":
		return &object.Boolean{Value: leftValue != rightValue}
	default:
		return newError("unknown operator: %s %s %s", left.Type(), operator, right.Type())
	}

}

func evalStringInfixExpression(left object.Object, operator string, right object.Object) object.Object {
	leftValue := left.(*object.String).Value
	rightValue := right.(*object.String).Value

	switch operator {
	case "+":
		return &object.String{Value: leftValue + rightValue}
	case ">":
		return &object.Boolean{Value: leftValue > rightValue}
	case "<":
		return &object.Boolean{Value: leftValue < rightValue}
	case "==":
		return &object.Boolean{Value: leftValue == rightValue}
	case "!=":
		return &object.Boolean{Value: leftValue != rightValue}
	default:
		return newError("unknown operator: %s %s %s", left.Type(), operator, right.Type())
	}
}

func evalIfExpression(exp *ast.IfExpression, env *object.Environment) object.Object {
	condition := Eval(exp.Condition, env)
	if isError(condition) {
		return condition
	}

	if condition.Type() == object.INTEGER_OBJ || condition.(*object.Boolean).Value {
		return Eval(exp.Consequence, env)
	} else {
		if exp.Alternative != nil {
			return Eval(exp.Alternative, env)
		} else {
			return NULL
		}
	}
}

func evalClass(exp *ast.ClassStatement, env *object.Environment) object.Object {
	class := &object.Class{Name: exp.Name}
	classEnv := object.NewClosedEnvironment(env)
	Eval(exp.Body, classEnv)
	class.Body = classEnv
	env.Set("Class", class)
	return class
}

func evalDefStatement(exp *ast.DefStatement, env *object.Environment) object.Object {
	method := &object.Method{Name: exp.Name.Value, Parameters: exp.Parameters, Body: exp.BlockStatement, Env: env}
	env.Set("_method_"+method.Name, method)
	return method
}

//func applyFunction(fn object.Object, args []object.Object) object.Object {
//	switch fn := fn.(type) {
//	case *object.Function:
//
//		if len(fn.Parameters) != len(args) {
//			return newError("wrong arguments: expect=%d, got=%d", len(fn.Parameters), len(args))
//		}
//
//		extendedEnv := extendFunctionEnv(fn, args)
//		evaluated := Eval(fn.Body, extendedEnv)
//		return unwrapReturnValue(evaluated)
//
//	case *object.BuiltInFunction:
//		return fn.Fn(args...)
//
//	default:
//		return newError("not a function: %s", fn.Type())
//	}
//}
//
//func evalArgs(exps []ast.Expression, env *object.Environment) []object.Object {
//	args := []object.Object{}
//
//	for _, exp := range exps {
//		arg := Eval(exp, env)
//		args = append(args, arg)
//		if isError(arg) {
//			return []object.Object{arg}
//		}
//	}
//
//	return args
//}

func extendMethodEnv(method *object.Method, args []object.Object) *object.Environment {
	e := object.NewClosedEnvironment(method.Env)

	for i, arg := range args {
		argName := method.Parameters[i].Value
		e.Set(argName, arg)
	}

	return e
}

func newError(format string, args ...interface{}) *object.Error {
	return &object.Error{Message: fmt.Sprintf(format, args...)}
}

func isError(obj object.Object) bool {
	if obj != nil {
		return obj.Type() == object.ERROR_OBJ
	}
	return false
}

//func unwrapReturnValue(obj object.Object) object.Object {
//	if returnValue, ok := obj.(*object.ReturnValue); ok {
//		return returnValue.Value
//	}
//
//	return obj
//}
