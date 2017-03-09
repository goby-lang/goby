package evaluator

import (
	"github.com/st0012/Rooby/ast"
	"github.com/st0012/Rooby/object"
)

func evalAssignStatement(stmt *ast.AssignStatement, scope *object.Scope) object.Object {
	value := Eval(stmt.Value, scope)

	if isError(value) {
		return value
	}

	switch variableName := stmt.Name.(type) {
	case *ast.Identifier:
		return scope.Env.Set(variableName.Value, value)
	case *ast.Constant:
		return scope.Env.Set(variableName.Value, value)
	case *ast.InstanceVariable:
		switch ivScope := scope.Self.(type) {
		case *object.RObject:
			return ivScope.InstanceVariables.Set(variableName.Value, value)
		case *object.RClass:
			return newError("Can not define instance variable %s in a class.", variableName.Value)
		default:
			return newError("Can not define instance variable %s in %T", ivScope)
		}
	default:
		return newError("Can not define variable to a %T", variableName)
	}

	return newError("Can not define variable %s in %s", stmt.Name.String(), scope.Self.Inspect())
}

func evalBlockStatements(stmts []ast.Statement, scope *object.Scope) object.Object {
	var result object.Object

	for _, statement := range stmts {

		result = Eval(statement, scope)

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

func evalClassStatement(exp *ast.ClassStatement, scope *object.Scope) object.Object {
	class := object.InitializeClass(exp.Name.Value, scope)

	// Evaluate superclass
	if exp.SuperClass != nil {

		constant := evalConstant(exp.SuperClass, scope)
		inheritedClass, ok := constant.(*object.RClass)
		if !ok {
			newError("Constant %s is not a class. got=%T", exp.SuperClass.Value, constant)
		}

		class.SuperClass = inheritedClass
	}

	Eval(exp.Body, class.Scope) // Eval class's content

	scope.Env.Set(class.Name, class)
	return class
}

func evalDefStatement(exp *ast.DefStatement, scope *object.Scope) object.Object {
	class, ok := scope.Self.(*object.RClass)
	// scope must be a class for now.
	if !ok {
		return newError("Method %s must be defined inside a Class. got=%T", exp.Name.Value, scope.Self)
	}

	method := &object.Method{Name: exp.Name.Value, Parameters: exp.Parameters, Body: exp.BlockStatement, Scope: scope}

	switch exp.Receiver.(type) {
	case nil:
		class.Methods.Set(method.Name, method)
	case *ast.SelfExpression:
		class.Class.Methods.Set(method.Name, method)
	}

	return method
}

func evalWhileStatement(exp *ast.WhileStatement, scope *object.Scope) object.Object {
	condition := exp.Condition

	con := Eval(condition, scope)

	for con != object.FALSE && con != object.NULL {
		Eval(exp.Body, scope)
		con = Eval(condition, scope)
	}

	return nil
}

func evalYieldStatement(node *ast.YieldStatement, scope *object.Scope) object.Object {
	block, ok := scope.Env.GetCurrent("block")
	if ok {
		b := block.(*object.Method)
		var args []object.Object

		for _, arg := range node.Arguments {
			args = append(args, Eval(arg, scope))
		}
		return evalMethodObject(scope.Self, b, args, nil)
	}

	return newError("Yield without a block")
}
