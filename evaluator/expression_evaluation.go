package evaluator

import (
	"github.com/st0012/Rooby/ast"
)

func evalPrefixExpression(operator string, right Object) Object {
	switch operator {
	case "!":
		return evalBangPrefixExpression(right)
	case "-":
		return evalMinusPrefixExpression(right)
	}
	return newError("unknown operator: %s%s", operator, right.Type())
}

func evalBangPrefixExpression(right Object) *BooleanObject {
	switch right {
	case FALSE:
		return TRUE
	case NULL:
		return TRUE
	default:
		return FALSE
	}
}

func evalMinusPrefixExpression(right Object) Object {
	if right.Type() != INTEGER_OBJ {
		return newError("unknown operator: %s%s", "-", right.Type())
	}
	value := right.(*IntegerObject).Value
	return &IntegerObject{Value: -value, Class: IntegerClass}
}

func evalInfixExpression(left Object, operator string, right Object) Object {
	result := sendMethodCall(left, operator, []Object{right}, nil)

	if err, ok := result.(*Error); ok {
		return err
	}

	return result
}

func evalIfExpression(exp *ast.IfExpression, scope *Scope) Object {
	condition := Eval(exp.Condition, scope)
	if isError(condition) {
		return condition
	}

	if condition.Type() == INTEGER_OBJ || condition.(*BooleanObject).Value {
		return Eval(exp.Consequence, scope)
	} else {
		if exp.Alternative != nil {
			return Eval(exp.Alternative, scope)
		} else {
			return NULL
		}
	}
}

func evalIdentifier(node *ast.Identifier, scope *Scope) Object {
	// check if it's a variable
	if val, ok := scope.Env.Get(node.Value); ok {
		return val
	}

	// check if it's a method
	receiver := scope.Self
	method_name := node.Value
	args := []Object{}

	error := newError("undefined local variable or method `%s' for %s", method_name, receiver.Inspect())

	switch receiver := receiver.(type) {
	case *RClass:
		method := receiver.LookupClassMethod(method_name)

		if method == nil {
			return error
		} else {
			evaluated := evalClassMethod(receiver, method, args, nil)
			return unwrapReturnValue(evaluated)
		}
	case *RObject:
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

func evalConstant(node *ast.Constant, scope *Scope) Object {
	if val, ok := scope.Env.Get(node.Value); ok {
		return val
	}

	return newError("constant %s not found in: %s", node.Value, scope.Self.Inspect())
}

func evalInstanceVariable(node *ast.InstanceVariable, scope *Scope) Object {
	instance := scope.Self.(*RObject)
	if val, ok := instance.InstanceVariables.Get(node.Value); ok {
		return val
	}

	return newError("instance variable %s not found in: %s", node.Value, instance.Inspect())
}
