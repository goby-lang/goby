package evaluator

import (
	"github.com/st0012/Rooby/ast"
	"github.com/st0012/Rooby/object"
)

func evalPrefixExpression(operator string, right object.Object) object.Object {
	switch operator {
	case "!":
		return evalBangPrefixExpression(right)
	case "-":
		return evalMinusPrefixExpression(right)
	}
	return newError("unknown operator: %s%s", operator, right.Type())
}

func evalBangPrefixExpression(right object.Object) *object.BooleanObject {
	switch right {
	case object.FALSE:
		return object.TRUE
	case object.NULL:
		return object.TRUE
	default:
		return object.FALSE
	}
}

func evalMinusPrefixExpression(right object.Object) object.Object {
	if right.Type() != object.INTEGER_OBJ {
		return newError("unknown operator: %s%s", "-", right.Type())
	}
	value := right.(*object.IntegerObject).Value
	return &object.IntegerObject{Value: -value, Class: object.IntegerClass}
}

func evalInfixExpression(left object.Object, operator string, right object.Object) object.Object {
	result := sendMethodCall(left, operator, []object.Object{right}, nil)

	if err, ok := result.(*object.Error); ok {
		return err
	}

	return result
}

func evalIfExpression(exp *ast.IfExpression, scope *object.Scope) object.Object {
	condition := Eval(exp.Condition, scope)
	if isError(condition) {
		return condition
	}

	if condition.Type() == object.INTEGER_OBJ || condition.(*object.BooleanObject).Value {
		return Eval(exp.Consequence, scope)
	} else {
		if exp.Alternative != nil {
			return Eval(exp.Alternative, scope)
		} else {
			return object.NULL
		}
	}
}

func evalIdentifier(node *ast.Identifier, scope *object.Scope) object.Object {
	// check if it's a variable
	if val, ok := scope.Env.Get(node.Value); ok {
		return val
	}

	// check if it's a method
	receiver := scope.Self
	method_name := node.Value
	args := []object.Object{}

	error := newError("undefined local variable or method `%s' for %s", method_name, receiver.Inspect())

	switch receiver := receiver.(type) {
	case *object.RClass:
		method := receiver.LookupClassMethod(method_name)

		if method == nil {
			return error
		} else {
			evaluated := evalClassMethod(receiver, method, args, nil)
			return unwrapReturnValue(evaluated)
		}
	case *object.RObject:
		method := receiver.Class.LookupInstanceMethod(method_name)

		if method == nil {
			return error
		} else {
			evaluated := evalInstanceMethod(receiver, method, args, nil)
			return unwrapReturnValue(evaluated)

		}
	}

	return error
}

func evalConstant(node *ast.Constant, scope *object.Scope) object.Object {
	if val, ok := scope.Env.Get(node.Value); ok {
		return val
	}

	return newError("constant %s not found in: %s", node.Value, scope.Self.Inspect())
}

func evalInstanceVariable(node *ast.InstanceVariable, scope *object.Scope) object.Object {
	instance := scope.Self.(*object.RObject)
	if val, ok := instance.InstanceVariables.Get(node.Value); ok {
		return val
	}

	return newError("instance variable %s not found in: %s", node.Value, instance.Inspect())
}
