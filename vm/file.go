package vm

import "path/filepath"

func initializeFileClass(vm *VM) {
	class := initializeClass("File")

	for _, m := range builtinFileClassMethods {
		class.ClassMethods.set(m.Name, m)
	}

	vm.constants["File"] = &Pointer{Target: class}
}

var builtinFileClassMethods = []*BuiltInMethod{
	{
		Name: "extname",
		Fn: func(receiver Object) builtinMethodBody {
			return func(ma methodArgs) Object {
				filename := ma.args[0].(*StringObject).Value
				return initializeString(filepath.Ext(filename))
			}
		},
	},
}
