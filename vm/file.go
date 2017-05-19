package vm

import (
	"os"
	"path/filepath"
)

func initializeFileClass(vm *VM) {
	class := initializeClass("File", false)

	for _, m := range builtinFileClassMethods {
		class.ClassMethods.set(m.Name, m)
	}

	vm.constants["File"] = &Pointer{Target: class}
}

var builtinFileClassMethods = []*BuiltInMethod{
	{
		Name: "extname",
		Fn: func(receiver Object) builtinMethodBody {
			return func(v *VM, args []Object, blockFrame *callFrame) Object {
				filename := args[0].(*StringObject).Value
				return initializeString(filepath.Ext(filename))
			}
		},
	},
	{
		Name: "chmod",
		Fn: func(receiver Object) builtinMethodBody {
			return func(v *VM, args []Object, blockFrame *callFrame) Object {
				filemod := args[0].(*IntegerObject).Value

				for i := 1; i < len(args)-1; i++ {
					filename := args[i].(*StringObject).Value
					err := os.Chmod(filename, os.FileMode(uint32(filemod)))
					if err != nil {
						panic(err)
					}
				}

				return initilaizeInteger(len(args) - 1)
			}
		},
	},
}
