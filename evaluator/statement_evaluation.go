package evaluator

import (
	"github.com/st0012/rooby/ast"
	"github.com/st0012/rooby/object"
)

func evalLetStatement(stmt *ast.LetStatement, scope *object.Scope) object.Object {
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
		case *object.BaseObject:
			return ivScope.InstanceVariables.Set(variableName.Value, value)
		case *object.Class:
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
	class := &object.Class{Name: exp.Name, Scope: scope, ClassMethods: object.NewEnvironment(), InstanceMethods: object.NewEnvironment(), SuperClass: object.ObjectClass, Class: object.ObjectClass}

	// Evaluate superclass
	if exp.SuperClass != nil {
		constant := evalConstant(exp.SuperClass, scope)
		inheritedClass, ok := constant.(*object.Class)

		if !ok {
			newError("Constant %s is not a class. got=%T", exp.SuperClass.Value, constant)
		}

		class.SuperClass = inheritedClass
	}

	// Create scope for this class
	classEnv := object.NewClosedEnvironment(scope.Env)
	classScope := &object.Scope{Env: classEnv, Self: class}

	Eval(exp.Body, classScope) // Eval class's content

	// Class's built in methods like `new`
	for method_name, method := range builtinClassMethods(class) {
		class.ClassMethods.Set(method_name, method)
	}

	scope.Env.Set(class.Name.Value, class)
	return class
}

func evalDefStatement(exp *ast.DefStatement, scope *object.Scope) object.Object {
	class, ok := scope.Self.(*object.Class)
	// scope must be a class for now.
	if !ok {
		return newError("Method %s must be defined inside a Class. got=%T", exp.Name.Value, scope.Self)
	}

	method := &object.Method{Name: exp.Name.Value, Parameters: exp.Parameters, Body: exp.BlockStatement, Scope: scope}
	class.InstanceMethods.Set(method.Name, method)
	return method
}
