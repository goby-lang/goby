package initializer

import (
	"fmt"
	"github.com/st0012/rooby/object"
	"github.com/st0012/rooby/ast"
)

var(
	ObjectClass *object.RClass
	ClassClass  *object.RClass
)

var BuiltinGlobalMethods = []*object.BuiltInMethod{
	{
		Fn: func(receiver object.Object) object.BuiltinMethodBody {
			return func(args ...object.Object) object.Object {
				for _, arg := range args {
					fmt.Println(arg.Inspect())
				}

				return NULL
			}
		},
		Name: "puts",
	},
	{
		Fn: func(receiver object.Object) object.BuiltinMethodBody {
			return func(args ...object.Object) object.Object {
				switch r := receiver.(type) {
				case object.BaseObject:
					return r.ReturnClass()
				case object.Class:
					return r.ReturnClass()
				default:
					return &object.Error{Message: fmt.Sprint("Can't call class on %T", r)}
				}
			}
		},
		Name: "class",
	},
}

var BuiltinClassMethods = []*object.BuiltInMethod{
	{
		Fn: func(receiver object.Object) object.BuiltinMethodBody {
			return func(args ...object.Object) object.Object {
				class := receiver.(*object.RClass)
				instance := InitializeInstance(class)
				initMethod := class.LookupInstanceMethod("initialize")

				if initMethod != nil {
					instance.InitializeMethod = initMethod.(*object.Method)
				}

				return instance
			}
		},
		Name: "new",
	},
	{
		Fn: func(receiver object.Object) object.BuiltinMethodBody {
			return func(args ...object.Object) object.Object {
				name := receiver.(object.Class).ReturnName()
				nameString := &object.StringObject{Value: name.Value}
				return nameString
			}
		},
		Name: "name",
	},
}


func initializeObjectClass() *object.RClass {
	name := &ast.Constant{Value: "Object"}
	globalMethods := object.NewEnvironment()

	for _, m := range BuiltinGlobalMethods {
		globalMethods.Set(m.Name, m)
	}

	class := &object.RClass{BaseClass: &object.BaseClass{Name: name, Class: ClassClass, Methods: globalMethods}}
	ObjectClass = class
	return class
}

func initializeClassClass() *object.RClass {
	methods := object.NewEnvironment()

	for _, m := range BuiltinClassMethods {
		methods.Set(m.Name, m)
	}

	name := &ast.Constant{Value: "Class"}
	class := &object.RClass{BaseClass: &object.BaseClass{Name: name, Methods: methods}}
	ClassClass = class
	return class
}