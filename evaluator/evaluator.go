package evaluator

import (
	"fmt"
	"github.com/st0012/Rooby/ast"
	"github.com/st0012/Rooby/initializer"
	"github.com/st0012/Rooby/object"
)

func Eval(node ast.Node, scope *object.Scope) object.Object {
	switch node := node.(type) {

	// Statements
	case *ast.Program:
		return evalProgram(node.Statements, scope)
	case *ast.ExpressionStatement:
		return Eval(node.Expression, scope)
	case *ast.BlockStatement:
		return evalBlockStatements(node.Statements, scope)
	case *ast.ReturnStatement:
		val := Eval(node.ReturnValue, scope)
		if isError(val) {
			return val
		}
		return &object.ReturnValue{Value: val}
	case *ast.AssignStatement:
		return evalAssignStatement(node, scope)
	case *ast.ClassStatement:
		return evalClassStatement(node, scope)
	case *ast.Identifier:
		return evalIdentifier(node, scope)
	case *ast.Constant:
		return evalConstant(node, scope)
	case *ast.InstanceVariable:
		return evalInstanceVariable(node, scope)
	case *ast.DefStatement:
		return evalDefStatement(node, scope)
	case *ast.WhileStatement:
		return evalWhileStatement(node, scope)
	case *ast.YieldStatement:
		// Only retrieve current scope's block.
		return evalYieldStatement(node, scope)

	// Expressions
	case *ast.IfExpression:
		return evalIfExpression(node, scope)
	case *ast.CallExpression:
		var block *object.Method
		receiver := Eval(node.Receiver, scope)
		args := evalArgs(node.Arguments, scope)

		if node.Block != nil {
			block = &object.Method{Body: node.Block, Parameters: node.BlockArguments, Scope: scope}
		}

		return sendMethodCall(receiver, node.Method, args, block)

	case *ast.PrefixExpression:
		val := Eval(node.Right, scope)
		if isError(val) {
			return val
		}
		return evalPrefixExpression(node.Operator, val)
	case *ast.InfixExpression:
		valLeft := Eval(node.Left, scope)
		if isError(valLeft) {
			return valLeft
		}

		valRight := Eval(node.Right, scope)
		if isError(valRight) {
			return valRight
		}

		return evalInfixExpression(valLeft, node.Operator, valRight)
	case *ast.SelfExpression:
		return scope.Self
	case *ast.IntegerLiteral:
		return initializer.InitilaizeInteger(node.Value)
	case *ast.StringLiteral:
		return initializer.InitializeString(node.Value)
	case *ast.Boolean:
		if node.Value {
			return initializer.TRUE
		}
		return initializer.FALSE
	case *ast.ArrayExpression:
		elements := []object.Object{}

		for _, exp := range node.Elements {
			elements = append(elements, Eval(exp, scope))
		}

		arr := initializer.InitializeArray(elements)
		return arr
	case *ast.HashExpression:
		pairs := map[string]object.Object{}

		for key, value := range node.Data {
			pairs[key] = Eval(value, scope)
		}

		hash := initializer.InitializeHash(pairs)
		return hash
	}

	return nil
}

func evalProgram(stmts []ast.Statement, scope *object.Scope) object.Object {
	var result object.Object

	for _, statement := range stmts {
		result = Eval(statement, scope)

		switch result := result.(type) {
		case *object.ReturnValue:
			return result.Value
		case *object.Error:
			return result
		}
	}

	return result
}

func sendMethodCall(receiver object.Object, method_name string, args []object.Object, block *object.Method) object.Object {
	error := newError("undefined method `%s' for %s", method_name, receiver.Inspect())

	switch receiver := receiver.(type) {
	case object.Class:
		method := receiver.LookupClassMethod(method_name)

		if method == nil {
			return error
		}

		evaluated := evalClassMethod(receiver, method, args, block)

		return unwrapReturnValue(evaluated)
	case object.BaseObject:
		if _, ok := receiver.(*object.Error); ok {
			return receiver
		}

		method := receiver.ReturnClass().LookupInstanceMethod(method_name)

		if method == nil {
			return error
		}

		evaluated := evalInstanceMethod(receiver, method, args, block)

		return unwrapReturnValue(evaluated)
	case *object.Error:
		return receiver
	default:
		return newError("not a valid receiver: %s", receiver.Inspect())
	}
}

func evalClassMethod(receiver object.Class, method object.Object, args []object.Object, block *object.Method) object.Object {
	switch m := method.(type) {
	case *object.Method:
		return evalMethodObject(receiver, m, args, block)
	case *object.BuiltInMethod:
		methodBody := m.Fn(receiver)
		evaluated := methodBody(args...)

		if m.Name == "new" {
			instance := evaluated.(*object.RObject)
			if instance.InitializeMethod != nil {
				evalInstanceMethod(instance, instance.InitializeMethod, args, block)
			}

			return instance
		}

		return evaluated
	case *object.Error:
		return m
	default:
		return newError("unknown class method type: %T)", m)
	}
}

func evalInstanceMethod(receiver object.BaseObject, method object.Object, args []object.Object, block *object.Method) object.Object {
	switch m := method.(type) {
	case *object.Method:
		return evalMethodObject(receiver, m, args, block)
	case *object.BuiltInMethod:
		methodBody := m.Fn(receiver)
		return methodBody(args...)
	case *object.Error:
		return m
	default:
		return newError("unknown instance method type: %T)", m)
	}
}

func evalArgs(exps []ast.Expression, scope *object.Scope) []object.Object {
	args := []object.Object{}

	for _, exp := range exps {
		arg := Eval(exp, scope)
		args = append(args, arg)
		if isError(arg) {
			return []object.Object{arg}
		}
	}

	return args
}

func evalMethodObject(receiver object.Object, m *object.Method, args []object.Object, block *object.Method) object.Object {
	if len(m.Parameters) != len(args) {
		return newError("wrong arguments: expect=%d, got=%d", len(m.Parameters), len(args))
	}

	methodEnv := m.ExtendEnv(args)
	if block != nil {
		methodEnv.Set("block", block)
	}
	scope := &object.Scope{Self: receiver, Env: methodEnv}
	return Eval(m.Body, scope)
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
