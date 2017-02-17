package evaluator

import (
	"github.com/st0012/rooby/ast"
	"github.com/st0012/rooby/object"
)

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

func evalClassStatement(exp *ast.ClassStatement, env *object.Environment) object.Object {
	class := &object.Class{Name: exp.Name}
	classEnv := object.NewClosedEnvironment(env)
	Eval(exp.Body, classEnv)

	for method_name, method := range builtins(class) {
		classEnv.Set(method_name, method)
	}

	class.Body = classEnv
	env.Set(class.Name.Value, class)
	return class
}

func evalDefStatement(exp *ast.DefStatement, env *object.Environment) object.Object {
	method := &object.Method{Name: exp.Name.Value, Parameters: exp.Parameters, Body: exp.BlockStatement, Env: env}
	env.Set("_method_"+method.Name, method)
	return method
}
