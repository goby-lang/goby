package vm

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"syscall"
)

func initializeFileClass(vm *VM) {
	class := vm.initializeClass("File", false)
	class.setBuiltInMethods(builtinFileClassMethods(), true)
	class.setBuiltInMethods(builtinFileInstanceMethods(), false)
	vm.builtInClasses[objectClass].constants["File"] = &Pointer{Target: class}

	vm.execGobyLib("file.gb")
}

// FileObject is a special type that contains file pointer so we can keep track on target file.
type FileObject struct {
	*baseObj
	File *os.File
}

// toString returns detailed infoof a array include elements it contains
func (f *FileObject) toString() string {
	return "<File: " + f.File.Name() + ">"
}

func (f *FileObject) toJSON() string {
	return f.toString()
}

var fileModeTable = map[string]int{
	"r":  syscall.O_RDONLY,
	"r+": syscall.O_RDWR,
	"w":  syscall.O_WRONLY,
	"w+": syscall.O_RDWR,
}

// Only initialize file related methods after it's being required.
func builtinFileClassMethods() []*BuiltInMethodObject {
	return []*BuiltInMethodObject{
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

							if md == syscall.O_RDWR || md == syscall.O_WRONLY {
								os.Create(fn)
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

					// TODO: Refactor this class retrieval mess
					fileObj := &FileObject{File: f, baseObj: &baseObj{class: t.vm.builtInClasses[objectClass].constants["File"].Target.(*RClass)}}

					return fileObj
				}
			},
		},
		{
			Name: "delete",
			Fn: func(receiver Object) builtinMethodBody {
				return func(t *thread, args []Object, blockFrame *callFrame) Object {
					for _, arg := range args {
						filename := arg.(*StringObject).Value
						err := os.Remove(filename)

						if err != nil {
							t.returnError(err.Error())
							return nil
						}
					}

					return t.vm.initIntegerObject(len(args))
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
					return t.vm.initStringObject(filepath.Ext(filename))
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
							filename = filepath.Join(t.vm.fileDir, filename)
						}

						err := os.Chmod(filename, os.FileMode(uint32(filemod)))
						if err != nil {
							t.returnError(err.Error())
						}
					}

					return t.vm.initIntegerObject(len(args) - 1)
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
						filename = filepath.Join(t.vm.fileDir, filename)
					}

					fileStats, err := os.Stat(filename)
					if err != nil {
						panic(err)
					}

					return t.vm.initIntegerObject(int(fileStats.Size()))
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
					return t.vm.initStringObject(filepath.Base(filename))
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

					return t.vm.initStringObject(filepath.Join(elements...))
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

					dirObject := t.vm.initStringObject(dir)
					fileObject := t.vm.initStringObject(file)

					return t.vm.initArrayObject([]Object{dirObject, fileObject})
				}
			},
		},
		{
			Name: "exist",
			Fn: func(receiver Object) builtinMethodBody {
				return func(t *thread, args []Object, blockFrame *callFrame) Object {
					filename := args[0].(*StringObject).Value
					_, err := os.Stat(filename)

					if err != nil {
						return FALSE
					}

					return TRUE
				}
			},
		},
	}

}

// Only initialize file related methods after it's being required.
func builtinFileInstanceMethods() []*BuiltInMethodObject {
	return []*BuiltInMethodObject{
		{
			Name: "name",
			Fn: func(receiver Object) builtinMethodBody {
				return func(t *thread, args []Object, blockFrame *callFrame) Object {
					name := receiver.(*FileObject).File.Name()
					return t.vm.initStringObject(name)
				}
			},
		},
		{
			// Returns size of file in bytes.
			//
			// ```ruby
			// File.new("loop.gb").size # => 321123
			// ```
			// @return [Integer]
			Name: "size",
			Fn: func(receiver Object) builtinMethodBody {
				return func(t *thread, args []Object, blockFrame *callFrame) Object {
					file := receiver.(*FileObject).File

					fileStats, err := os.Stat(file.Name())
					if err != nil {
						panic(err)
					}

					return t.vm.initIntegerObject(int(fileStats.Size()))
				}
			},
		},
		{
			Name: "read",
			Fn: func(receiver Object) builtinMethodBody {
				return func(t *thread, args []Object, blockFrame *callFrame) Object {
					file := receiver.(*FileObject).File
					data, err := ioutil.ReadFile(file.Name())

					if err != nil {
						t.returnError(err.Error())
					}

					return t.vm.initStringObject(string(data))
				}
			},
		},
		{
			Name: "write",
			Fn: func(receiver Object) builtinMethodBody {
				return func(t *thread, args []Object, blockFrame *callFrame) Object {
					file := receiver.(*FileObject).File
					data := args[0].(*StringObject).Value
					length, err := file.Write([]byte(data))

					if err != nil {
						t.returnError(err.Error())
					}

					return t.vm.initIntegerObject(length)
				}
			},
		},
		{
			Name: "close",
			Fn: func(receiver Object) builtinMethodBody {
				return func(t *thread, args []Object, blockFrame *callFrame) Object {
					file := receiver.(*FileObject).File
					file.Close()

					return NULL
				}
			},
		},
	}

}
