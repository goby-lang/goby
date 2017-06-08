package vm

import (
	"os"
	"path/filepath"
	"strings"
	"syscall"
)

var fileClass *RClass

func initializeFileClass(vm *VM) {
	class := initializeClass("File", false)
	class.setBuiltInMethods(builtinFileClassMethods, true)
	vm.constants["File"] = &Pointer{Target: class}
	fileClass = class
}

type FileObject struct {
	Class *RClass
	File  *os.File
}

// Inspect returns detailed infoof a array include elements it contains
func (f *FileObject) Inspect() string {
	return "<File: " + f.File.Name() + ">"
}

// returnClass returns current object's class, which is RArray
func (f *FileObject) returnClass() Class {
	return f.Class
}

var fileModeTable = map[string]int{
	"r":  syscall.O_RDONLY,
	"r+": syscall.O_RDWR,
	"w":  syscall.O_WRONLY,
	"w+": syscall.O_RDWR,
}

var builtinFileClassMethods = []*BuiltInMethodObject{
	{
		// Finds the file with given filename and initializes a file object with it.
		//
		// ```ruby
		// File.new("./samples/server.gb")
		// ```
		// @param filename [String]
		// @return [File]
		Name: "new",
		Fn: func(receiver Object) builtinMethodBody {
			return func(t *thread, args []Object, blockFrame *callFrame) Object {
				var fn string
				var mode int
				var perm os.FileMode

				if len(args) < 1 {
					return newError("Expect at least a filename to open file")
				}

				if len(args) >= 1 {
					fn = args[0].(*StringObject).Value
					mode = syscall.O_RDONLY
					perm = os.FileMode(0755)

					if len(args) >= 2 {
						m := args[1].(*StringObject).Value
						md, ok := fileModeTable[m]

						if !ok {
							t.returnError("Unknown file mode: " + m)
						}

						mode = md
						perm = os.FileMode(0755)

						if len(args) == 3 {
							p := args[2].(*IntegerObject).Value
							perm = os.FileMode(p)
						}
					}
				}

				f, err := os.OpenFile(fn, mode, perm)

				if err != nil {
					t.returnError(err.Error())
				}

				fileObj := &FileObject{File: f, Class: fileClass}

				return fileObj
			}
		},
	},
	{
		// Returns extension part of file.
		//
		// ```ruby
		// File.extname("loop.gb") # => .gb
		// ```
		// @param filename [String]
		// @return [String]
		Name: "extname",
		Fn: func(receiver Object) builtinMethodBody {
			return func(t *thread, args []Object, blockFrame *callFrame) Object {
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
		// @param filename [String]
		// @return [Integer]
		Name: "chmod",
		Fn: func(receiver Object) builtinMethodBody {
			return func(t *thread, args []Object, blockFrame *callFrame) Object {
				filemod := args[0].(*IntegerObject).Value
				for i := 1; i < len(args); i++ {
					filename := args[i].(*StringObject).Value
					if !filepath.IsAbs(filename) {
						filename = t.vm.fileDir + filename
					}

					err := os.Chmod(filename, os.FileMode(uint32(filemod)))
					if err != nil {
						t.returnError(err.Error())
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
		// @param filename [String]
		// @return [Integer]
		Name: "size",
		Fn: func(receiver Object) builtinMethodBody {
			return func(t *thread, args []Object, blockFrame *callFrame) Object {
				filename := args[0].(*StringObject).Value
				if !filepath.IsAbs(filename) {
					filename = t.vm.fileDir + filename
				}

				fileStats, err := os.Stat(filename)
				if err != nil {
					panic(err)
				}

				return initilaizeInteger(int(fileStats.Size()))
			}
		},
	},
	{
		// Returns the last element from path.
		//
		// ```ruby
		// File.basename("/home/goby/plugin/loop.gb") # => loop.gb
		// ```
		// @param filepath [String]
		// @return [String]
		Name: "basename",
		Fn: func(receiver Object) builtinMethodBody {
			return func(t *thread, args []Object, blockFrame *callFrame) Object {
				filename := args[0].(*StringObject).Value
				return initializeString(filepath.Base(filename))
			}
		},
	},
	{
		// Returns string with joined elements.
		//
		// ```ruby
		// File.join("home", "goby", "plugin") # => home/goby/plugin
		// ```
		// @return [String]
		Name: "join",
		Fn: func(receiver Object) builtinMethodBody {
			return func(t *thread, args []Object, blockFrame *callFrame) Object {
				var elements []string
				for i := 0; i < len(args); i++ {
					next := args[i].(*StringObject).Value
					elements = append(elements, next)
				}

				return initializeString(strings.Join(elements, string(filepath.Separator)))
			}
		},
	},
	{
		// Returns array of path and file.
		//
		// ```ruby
		// File.split("/home/goby/.settings") # => ["/home/goby/", ".settings"]
		// ```
		// @param filepath [String]
		// @return [Array]
		Name: "split",
		Fn: func(receiver Object) builtinMethodBody {
			return func(t *thread, args []Object, blockFrame *callFrame) Object {
				filename := args[0].(*StringObject).Value
				dir, file := filepath.Split(filename)

				dirObject := initializeString(dir)
				fileObject := initializeString(file)

				return initializeArray([]Object{dirObject, fileObject})
			}
		},
	},
}
