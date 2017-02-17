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
		return env.Set(node.Name.TokenLiteral(), val)
	case *ast.ClassStatement:
		return evalClassStatement(node, env)
	case *ast.Identifier:
		return evalIdentifier(node, env)
	case *ast.Constant:
		return evalConstant(node, env)
	case *ast.InstanceVariable:
		return evalInstanceVariable(node, env)
	case *ast.DefStatement:
		return evalDefStatement(node, env)

	// Expressions
	case *ast.IfExpression:
		return evalIfExpression(node, env)
	case *ast.CallExpression:
		receiver := Eval(node.Receiver, env)
		args := evalArgs(node.Arguments, env)
		return sendMethodCall(receiver, node.Method.Value, args)

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

func sendMethodCall(receiver object.Object, method_name string, args []object.Object) object.Object {
	switch receiver := receiver.(type) {
	case *object.Class:
		evaluated := evalClassMethod(receiver, method_name, args)

		return unwrapReturnValue(evaluated)
	case *object.BaseObject:

		evaluated := evalInstanceMethod(receiver, method_name, args)
		return unwrapReturnValue(evaluated)
	default:
		return newError("not a valid receiver: %s", receiver.Type())
	}
}

func evalClassMethod(receiver object.Object, method_name string, args []object.Object) object.Object {
	method, ok := receiver.(*object.Class).Body.Get(method_name)
	if !ok {
		return &object.Error{Message: fmt.Sprintf("undefined method %s for class %s", method_name, receiver.Inspect())}
	}

	switch m := method.(type) {
	case *object.Method:
		if len(m.Parameters) != len(args) {
			return newError("wrong arguments: expect=%d, got=%d", len(m.Parameters), len(args))
		}

		methodEnv := extendMethodEnv(m, args)
		return Eval(m.Body, methodEnv)
	case *object.BuiltInMethod:
		return m.Fn(args...)
	default:
		return newError("unknown method type")
	}

}

func evalInstanceMethod(receiver object.Object, method_name string, args []object.Object) object.Object {
	method, ok := receiver.(*object.BaseObject).Class.Body.Get("_method_" + method_name)
	if !ok {
		return &object.Error{Message: fmt.Sprintf("undefined method %s for class %s", method_name, receiver.Inspect())}
	}
	switch m := method.(type) {
	case *object.Method:
		if len(m.Parameters) != len(args) {
			return newError("wrong arguments: expect=%d, got=%d", len(m.Parameters), len(args))
		}

		methodEnv := extendMethodEnv(m, args)
		return Eval(m.Body, methodEnv)
	default:
		return newError("unknown method type")
	}

}

func evalArgs(exps []ast.Expression, env *object.Environment) []object.Object {
	args := []object.Object{}

	for _, exp := range exps {
		arg := Eval(exp, env)
		args = append(args, arg)
		if isError(arg) {
			return []object.Object{arg}
		}
	}

	return args
}

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

func unwrapReturnValue(obj object.Object) object.Object {
	if returnValue, ok := obj.(*object.ReturnValue); ok {
		return returnValue.Value
	}

	return obj
}
