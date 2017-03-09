package evaluator

import (
	"github.com/st0012/Rooby/ast"
)

func evalAssignStatement(stmt *ast.AssignStatement, scope *Scope) Object {
	value := Eval(stmt.Value, scope)

	if isError(value) {
		return value
	}

	varName := stmt.Name.ReturnValue()
	env, ok := scope.Env.GetValueLocation(varName)

	if ok {
		env.Set(varName, value)
	}

	switch variableName := stmt.Name.(type) {
	case *ast.Identifier, *ast.Constant:
		return scope.Env.Set(varName, value)
	case *ast.InstanceVariable:
		switch ivScope := scope.Self.(type) {
		case *RObject:
			return ivScope.InstanceVariables.Set(varName, value)
		case *RClass:
			return newError("Can not define instance variable %s in a class.", variableName.Value)
		default:
			return newError("Can not define instance variable %s in %T", ivScope)
		}
	default:
		return newError("Can not define variable to a %T", variableName)
	}

	return newError("Can not define variable %s in %s", stmt.Name.String(), scope.Self.Inspect())
}

func evalBlockStatements(stmts []ast.Statement, scope *Scope) Object {
	var result Object

	for _, statement := range stmts {

		result = Eval(statement, scope)

		if result != nil {
			switch result := result.(type) {
			case *ReturnValue:
				return result
			case *Error:
				return result
			}
		}
	}

	return result
}

func evalClassStatement(exp *ast.ClassStatement, scope *Scope) Object {
	class := InitializeClass(exp.Name.Value, scope)

	// Evaluate superclass
	if exp.SuperClass != nil {

		constant := evalConstant(exp.SuperClass, scope)
		inheritedClass, ok := constant.(*RClass)
		if !ok {
			newError("Constant %s is not a class. got=%T", exp.SuperClass.Value, constant)
		}

		class.SuperClass = inheritedClass
	}

	Eval(exp.Body, class.Scope) // Eval class's content

	scope.Env.Set(class.Name, class)
	return class
}

func evalDefStatement(exp *ast.DefStatement, scope *Scope) Object {
	class, ok := scope.Self.(*RClass)
	// scope must be a class for now.
	if !ok {
		return newError("Method %s must be defined inside a Class. got=%T", exp.Name.Value, scope.Self)
	}

	method := &Method{Name: exp.Name.Value, Parameters: exp.Parameters, Body: exp.BlockStatement, Scope: scope}

	switch exp.Receiver.(type) {
	case nil:
		class.Methods.Set(method.Name, method)
	case *ast.SelfExpression:
		class.Class.Methods.Set(method.Name, method)
	}

	return method
}

func evalWhileStatement(exp *ast.WhileStatement, scope *Scope) Object {
	condition := exp.Condition

	con := Eval(condition, scope)

	for con != FALSE && con != NULL {
		Eval(exp.Body, scope)
		con = Eval(condition, scope)
	}

	return nil
}

func evalYieldStatement(node *ast.YieldStatement, scope *Scope) Object {
	block, ok := scope.Env.GetCurrent("block")
	if ok {
		b := block.(*Method)
		var args []Object

		for _, arg := range node.Arguments {
			args = append(args, Eval(arg, scope))
		}
		return evalMethodObject(scope.Self, b, args, nil)
	}

	return newError("Yield without a block")
}
