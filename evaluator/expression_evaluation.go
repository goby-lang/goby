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
			return NULL
		}
	}
}

func evalIdentifier(node *ast.Identifier, scope *object.Scope) object.Object {
	if val, ok := scope.Env.Get(node.Value); ok {
		return val
	}

	return newError("identifier not found: %s", node.Value)
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
