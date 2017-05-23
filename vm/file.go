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
		// Returns extension part of file.
		//
		// ```ruby
		// File.extname("loop.gb") # => .gb
		// ```
		Name: "extname",
		Fn: func(receiver Object) builtinMethodBody {
			return func(v *VM, args []Object, blockFrame *callFrame) Object {
				filename := args[0].(*StringObject).Value
				return initializeString(filepath.Ext(filename))
			}
		},
	},
	{
		// Changes the mode of the file.
		// Return number of files.
		//
		// ```ruby
		// File.chmod(0755, "test.sh") # => 1
		// File.chmod(0755, "goby", "../test.sh") # => 2
		// ```
		Name: "chmod",
		Fn: func(receiver Object) builtinMethodBody {
			return func(v *VM, args []Object, blockFrame *callFrame) Object {
				filemod := args[0].(*IntegerObject).Value
				for i := 1; i < len(args)-1; i++ {
					filename := args[i].(*StringObject).Value
					if !filepath.IsAbs(filename) {
						filename = v.fileDir + filename
					}

					err := os.Chmod(filename, os.FileMode(uint32(filemod)))
					if err != nil {
						panic(err)
					}
				}

				return initilaizeInteger(len(args) - 1)
			}
		},
	},
	{
		// Returns size of file in bytes.
		//
		// ```ruby
		// File.size("loop.gb") # => 321123
		// ```
		Name: "size",
		Fn: func(receiver Object) builtinMethodBody {
			return func(v *VM, args []Object, blockFrame *callFrame) Object {
				filename := args[0].(*StringObject).Value
				if !filepath.IsAbs(filename) {
					filename = v.fileDir + filename
				}

				fileStats, err := os.Stat(filename)
				if err != nil {
					panic(err)
				}

				return initilaizeInteger(int(fileStats.Size()))
			}
		},
	},
}
