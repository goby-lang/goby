package evaluator

import (
	"github.com/st0012/rooby/ast"
	"github.com/st0012/rooby/object"
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

func evalBangPrefixExpression(right object.Object) *object.Boolean {
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

func evalIfExpression(exp *ast.IfExpression, scope *object.Scope) object.Object {
	condition := Eval(exp.Condition, scope)
	if isError(condition) {
		return condition
	}

	if condition.Type() == object.INTEGER_OBJ || condition.(*object.Boolean).Value {
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
	case *object.Class:
		method := receiver.LookupClassMethod(method_name)

		if method == nil {
			return error
		} else {
			evaluated := evalClassMethod(receiver, method, args)
			return unwrapReturnValue(evaluated)
		}
	case *object.BaseObject:
		method := receiver.Class.LookUpInstanceMethod(method_name)

		if method == nil {
			return error
		} else {
			evaluated := evalInstanceMethod(receiver, method, args)
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
	instance := scope.Self.(*object.BaseObject)
	if val, ok := instance.InstanceVariables.Get(node.Value); ok {
		return val
	}

	return newError("instance variable %s not found in: %s", node.Value, instance.Inspect())
}
