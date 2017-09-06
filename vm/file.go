package vm

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"syscall"
)

var fileModeTable = map[string]int{
	"r":  syscall.O_RDONLY,
	"r+": syscall.O_RDWR,
	"w":  syscall.O_WRONLY,
	"w+": syscall.O_RDWR,
}

func (vm *VM) initFileClass() *RClass {
	fc := vm.initializeClass(fileClass, false)
	fc.setBuiltInMethods(builtinFileClassMethods(), true)
	fc.setBuiltInMethods(builtinFileInstanceMethods(), false)

	vm.libFiles = append(vm.libFiles, "file.gb")

	return fc
}

func (vm *VM) initFileObject(f *os.File) *FileObject {
	return &FileObject{
		baseObj: &baseObj{class: vm.topLevelClass(fileClass)},
		File:    f,
	}
}

// FileObject is a special type that contains file pointer so we can keep track on target file.
type FileObject struct {
	*baseObj
	File *os.File
}

// Polymorphic helper functions -----------------------------------------
func (f *FileObject) toString() string {
	return "<File: " + f.File.Name() + ">"
}

func (f *FileObject) toJSON() string {
	return f.toString()
}

func builtinFileClassMethods() []*BuiltInMethodObject {
	return []*BuiltInMethodObject{
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
					filename := args[0].(*StringObject).value
					return t.vm.initStringObject(filepath.Base(filename))
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
					filemod := args[0].(*IntegerObject).value
					for i := 1; i < len(args); i++ {
						filename := args[i].(*StringObject).value
						if !filepath.IsAbs(filename) {
							filename = filepath.Join(t.vm.fileDir, filename)
						}

						err := os.Chmod(filename, os.FileMode(uint32(filemod)))
						if err != nil {
							return t.vm.initErrorObject(InternalError, err.Error())
						}
					}

					return t.vm.initIntegerObject(len(args) - 1)
				}
			},
		},
		{
			Name: "delete",
			Fn: func(receiver Object) builtinMethodBody {
				return func(t *thread, args []Object, blockFrame *callFrame) Object {
					for _, arg := range args {
						filename := arg.(*StringObject).value
						err := os.Remove(filename)

						if err != nil {
							return t.vm.initErrorObject(InternalError, err.Error())
						}
					}

					return t.vm.initIntegerObject(len(args))
				}
			},
		},
		{
			Name: "exist?",
			Fn: func(receiver Object) builtinMethodBody {
				return func(t *thread, args []Object, blockFrame *callFrame) Object {
					filename := args[0].(*StringObject).value
					_, err := os.Stat(filename)

					if err != nil {
						return FALSE
					}

					return TRUE
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
					filename := args[0].(*StringObject).value
					return t.vm.initStringObject(filepath.Ext(filename))
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
						next := args[i].(*StringObject).value
						elements = append(elements, next)
					}

					return t.vm.initStringObject(filepath.Join(elements...))
				}
			},
		},
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
						return t.vm.initErrorObject(InternalError, "Expect at least a filename to open file")
					}

					if len(args) >= 1 {
						fn = args[0].(*StringObject).value
						mode = syscall.O_RDONLY
						perm = os.FileMode(0755)

						if len(args) >= 2 {
							m := args[1].(*StringObject).value
							md, ok := fileModeTable[m]

							if !ok {
								return t.vm.initErrorObject(InternalError, "Unknown file mode: %s", m)
							}

							if md == syscall.O_RDWR || md == syscall.O_WRONLY {
								os.Create(fn)
							}

							mode = md
							perm = os.FileMode(0755)

							if len(args) == 3 {
								p := args[2].(*IntegerObject).value
								perm = os.FileMode(p)
							}
						}
					}

					f, err := os.OpenFile(fn, mode, perm)

					if err != nil {
						return t.vm.initErrorObject(InternalError, err.Error())
					}

					// TODO: Refactor this class retrieval mess
					fileObj := &FileObject{File: f, baseObj: &baseObj{class: t.vm.topLevelClass(fileClass)}}

					return fileObj
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
					filename := args[0].(*StringObject).value
					if !filepath.IsAbs(filename) {
						filename = filepath.Join(t.vm.fileDir, filename)
					}

					fileStats, err := os.Stat(filename)
					if err != nil {
						return t.vm.initErrorObject(InternalError, err.Error())
					}

					return t.vm.initIntegerObject(int(fileStats.Size()))
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
					filename := args[0].(*StringObject).value
					dir, file := filepath.Split(filename)

					dirObject := t.vm.initStringObject(dir)
					fileObject := t.vm.initStringObject(file)

					return t.vm.initArrayObject([]Object{dirObject, fileObject})
				}
			},
		},
	}
}

func builtinFileInstanceMethods() []*BuiltInMethodObject {
	return []*BuiltInMethodObject{
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
			Name: "read",
			Fn: func(receiver Object) builtinMethodBody {
				return func(t *thread, args []Object, blockFrame *callFrame) Object {
					file := receiver.(*FileObject).File
					data, err := ioutil.ReadFile(file.Name())

					if err != nil {
						return t.vm.initErrorObject(InternalError, err.Error())
					}

					return t.vm.initStringObject(string(data))
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
						return t.vm.initErrorObject(InternalError, err.Error())
					}

					return t.vm.initIntegerObject(int(fileStats.Size()))
				}
			},
		},
		{
			Name: "write",
			Fn: func(receiver Object) builtinMethodBody {
				return func(t *thread, args []Object, blockFrame *callFrame) Object {
					file := receiver.(*FileObject).File
					data := args[0].(*StringObject).value
					length, err := file.Write([]byte(data))

					if err != nil {
						return t.vm.initErrorObject(InternalError, err.Error())
					}

					return t.vm.initIntegerObject(length)
				}
			},
		},
	}
}
