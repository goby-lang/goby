package vm

import (
	"bufio"
	"io/ioutil"
	"os"
	"path/filepath"
	"syscall"

	"github.com/goby-lang/goby/vm/classes"
	"github.com/goby-lang/goby/vm/errors"
)

// FileObject is a special type that contains file pointer so we can keep track on target file.
// Using `File.open` with block is recommended because the instance (block variable) automatically closes.
//
// ```ruby
// File.open("/tmp/goby/out.txt", "w", 0755) do |f|
//   a = f.read
//   f.write(a + "12345")
// end         # f automatically closes
// ```
//
type FileObject struct {
	*BaseObj
	File *os.File
}

var fileModeTable = map[string]int{
	"r":  syscall.O_RDONLY,
	"r+": syscall.O_RDWR,
	"w":  syscall.O_WRONLY,
	"w+": syscall.O_RDWR,
}

// Class methods --------------------------------------------------------
func builtinFileClassMethods() []*BuiltinMethodObject {
	return []*BuiltinMethodObject{
		{
			// Returns the last element from path.
			//
			// ```ruby
			// File.basename("/home/goby/plugin/loop.gb") # => loop.gb
			// ```
			// @param filePath [String]
			// @return [String]
			Name: "basename",
			Fn: func(receiver Object, sourceLine int, t *Thread, args []Object, blockFrame *normalCallFrame) Object {
				if len(args) != 1 {
					return t.vm.InitErrorObject(errors.ArgumentError, sourceLine, errors.WrongNumberOfArgument, 1, len(args))
				}

				fn, ok := args[0].(*StringObject)
				if !ok {
					return t.vm.InitErrorObject(errors.TypeError, sourceLine, errors.WrongArgumentTypeFormat, classes.StringClass, args[0].Class().Name)
				}

				return t.vm.InitStringObject(filepath.Base(fn.value))

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
			// @param fileName [String]
			// @return [Integer]
			Name: "chmod",
			Fn: func(receiver Object, sourceLine int, t *Thread, args []Object, blockFrame *normalCallFrame) Object {
				if len(args) < 2 {
					return t.vm.InitErrorObject(errors.ArgumentError, sourceLine, errors.WrongNumberOfArgumentMore, 2, len(args))
				}

				mod, ok := args[0].(*IntegerObject)
				if !ok {
					return t.vm.InitErrorObject(errors.TypeError, sourceLine, errors.WrongArgumentTypeFormatNum, 1, classes.IntegerClass, args[0].Class().Name)
				}

				if !os.FileMode(mod.value).IsRegular() {
					return t.vm.InitErrorObject(errors.ArgumentError, sourceLine, errors.InvalidChmodNumber, mod.value)
				}

				for i := 1; i < len(args); i++ {
					fn, ok := args[i].(*StringObject)
					if !ok {
						return t.vm.InitErrorObject(errors.TypeError, sourceLine, errors.WrongArgumentTypeFormatNum, i+1, classes.StringClass, args[0].Class().Name)
					}

					if !filepath.IsAbs(fn.value) {
						fn.value = filepath.Join(t.vm.fileDir, fn.value)
					}

					err := os.Chmod(fn.value, os.FileMode(uint32(mod.value)))
					if err != nil {
						return t.vm.InitErrorObject(errors.IOError, sourceLine, err.Error())
					}
				}

				return t.vm.InitIntegerObject(len(args) - 1)

			},
		},
		// Deletes the specified files.
		// Return the number of deleted files.
		// The number of the argument can be zero, but deleting non-existent files causes an error.
		//
		// ```ruby
		// File.delete("test.sh")             # => 1
		// File.delete("test.sh", "test2.sh") # => 2
		// File.delete()                      # => 0
		// File.delete("non-existent.txt")    # =>
		// ```
		// @param fileName [String]
		// @return [Integer]
		{
			Name: "delete",
			Fn: func(receiver Object, sourceLine int, t *Thread, args []Object, blockFrame *normalCallFrame) Object {
				for i, arg := range args {
					fn, ok := arg.(*StringObject)
					if !ok {
						return t.vm.InitErrorObject(errors.TypeError, sourceLine, errors.WrongArgumentTypeFormatNum, i+1, classes.StringClass, args[i].Class().Name)
					}
					err := os.Remove(fn.value)

					if err != nil {
						return t.vm.InitErrorObject(errors.IOError, sourceLine, err.Error())
					}
				}

				return t.vm.InitIntegerObject(len(args))

			},
		},
		// Determines if the specified file.
		//
		// ```ruby
		// File.exist?("test.sh")             # => false
		// File.open("test.sh, "w", 0755)
		// File.exist?("test.sh")             # => true
		// ```
		// @param fileName [String]
		// @return [Boolean]
		{
			Name: "exist?",
			Fn: func(receiver Object, sourceLine int, t *Thread, args []Object, blockFrame *normalCallFrame) Object {
				if len(args) != 1 {
					return t.vm.InitErrorObject(errors.ArgumentError, sourceLine, errors.WrongNumberOfArgument, 1, len(args))
				}

				fn, ok := args[0].(*StringObject)
				if !ok {
					return t.vm.InitErrorObject(errors.TypeError, sourceLine, errors.WrongArgumentTypeFormat, classes.StringClass, args[0].Class().Name)
				}
				_, err := os.Stat(fn.value)

				if err != nil {
					return FALSE
				}

				return TRUE

			},
		},
		{
			// Returns the extension part of file.
			//
			// ```ruby
			// File.extname("loop.gb") # => .gb
			// ```
			//
			// @param fileName [String]
			// @return [String]
			Name: "extname",
			Fn: func(receiver Object, sourceLine int, t *Thread, args []Object, blockFrame *normalCallFrame) Object {
				if len(args) != 1 {
					return t.vm.InitErrorObject(errors.ArgumentError, sourceLine, errors.WrongNumberOfArgument, 1, len(args))
				}

				fn, ok := args[0].(*StringObject)
				if !ok {
					return t.vm.InitErrorObject(errors.TypeError, sourceLine, errors.WrongArgumentTypeFormat, classes.StringClass, args[0].Class().Name)
				}

				return t.vm.InitStringObject(filepath.Ext(fn.value))

			},
		},
		{
			// Returns the string with joined elements.
			// Arguments can be zero.
			//
			// ```ruby
			// File.join("home", "goby", "plugin") # => home/goby/plugin
			// ```
			//
			// @param fileName [String]
			// @return [String]
			Name: "join",
			Fn: func(receiver Object, sourceLine int, t *Thread, args []Object, blockFrame *normalCallFrame) Object {
				var e []string
				for i := 0; i < len(args); i++ {
					next, ok := args[i].(*StringObject)
					if !ok {
						return t.vm.InitErrorObject(errors.TypeError, sourceLine, errors.WrongArgumentTypeFormat, classes.StringClass, args[0].Class().Name)
					}

					e = append(e, next.value)
				}

				return t.vm.InitStringObject(filepath.Join(e...))

			},
		},
		{
			// Finds the file with given fileName and initializes a file object with it.
			// File permissions can be specified at the second or third argument.
			//
			// ```ruby
			// File.new("./samples/server.gb")
			//
			// File.new("../test_fixtures/file_test/size.gb", "r")
			//
			// File.new("../test_fixtures/file_test/size.gb", "r", 0755)
			// ```
			// @param fileName [String]
			// @return [File]
			Name: "new",
			Fn: func(receiver Object, sourceLine int, t *Thread, args []Object, blockFrame *normalCallFrame) Object {
				aLen := len(args)
				if aLen < 1 || aLen > 3 {
					return t.vm.InitErrorObject(errors.ArgumentError, sourceLine, errors.WrongNumberOfArgumentRange, 1, 3, aLen)
				}

				fn, ok := args[0].(*StringObject)
				if !ok {
					return t.vm.InitErrorObject(errors.TypeError, sourceLine, errors.WrongArgumentTypeFormatNum, 1, classes.StringClass, args[0].Class().Name)
				}

				mod := syscall.O_RDONLY
				perm := os.FileMode(0755)
				if aLen >= 2 {
					m, ok := args[1].(*StringObject)
					if !ok {
						return t.vm.InitErrorObject(errors.TypeError, sourceLine, errors.WrongArgumentTypeFormatNum, 2, classes.StringClass, args[1].Class().Name)
					}

					md, ok := fileModeTable[m.value]
					if !ok {
						return t.vm.InitErrorObject(errors.ArgumentError, sourceLine, "Unknown file mode: %s", m.value)
					}

					if md == syscall.O_RDWR || md == syscall.O_WRONLY {
						os.Create(fn.value)
					}

					mod = md
					perm = os.FileMode(0755)

					if aLen == 3 {
						p, ok := args[2].(*IntegerObject)
						if !ok {
							return t.vm.InitErrorObject(errors.TypeError, sourceLine, errors.WrongArgumentTypeFormatNum, 3, classes.IntegerClass, args[2].Class().Name)
						}

						if !os.FileMode(p.value).IsRegular() {
							return t.vm.InitErrorObject(errors.ArgumentError, sourceLine, errors.InvalidChmodNumber, p.value)
						}

						perm = os.FileMode(p.value)
					}
				}

				f, err := os.OpenFile(fn.value, mod, perm)

				if err != nil {
					return t.vm.InitErrorObject(errors.IOError, sourceLine, err.Error())
				}

				// TODO: Refactor this class retrieval mess
				fo := &FileObject{File: f, BaseObj: &BaseObj{class: t.vm.TopLevelClass(classes.FileClass)}}

				return fo

			},
		},
		{
			// Returns size of file in bytes.
			//
			// ```ruby
			// File.size("loop.gb") # => 321123
			// ```
			//
			// @param fileName [String]
			// @return [Integer]
			Name: "size",
			Fn: func(receiver Object, sourceLine int, t *Thread, args []Object, blockFrame *normalCallFrame) Object {
				if len(args) != 1 {
					return t.vm.InitErrorObject(errors.ArgumentError, sourceLine, errors.WrongNumberOfArgument, 1, len(args))
				}

				fn, ok := args[0].(*StringObject)
				if !ok {
					return t.vm.InitErrorObject(errors.TypeError, sourceLine, errors.WrongArgumentTypeFormat, classes.StringClass, args[0].Class().Name)
				}

				if !filepath.IsAbs(fn.value) {
					fn.value = filepath.Join(t.vm.fileDir, fn.value)
				}

				fs, err := os.Stat(fn.value)
				if err != nil {
					return t.vm.InitErrorObject(errors.IOError, sourceLine, err.Error())
				}

				return t.vm.InitIntegerObject(int(fs.Size()))

			},
		},
		{
			// Returns array of path and file.
			//
			// ```ruby
			// File.split("/home/goby/.settings") # => ["/home/goby/", ".settings"]
			// ```
			//
			// @param filePath [String]
			// @return [Array]
			Name: "split",
			Fn: func(receiver Object, sourceLine int, t *Thread, args []Object, blockFrame *normalCallFrame) Object {
				if len(args) != 1 {
					return t.vm.InitErrorObject(errors.ArgumentError, sourceLine, errors.WrongNumberOfArgument, 1, len(args))
				}

				fn, ok := args[0].(*StringObject)
				if !ok {
					return t.vm.InitErrorObject(errors.TypeError, sourceLine, errors.WrongArgumentTypeFormat, classes.StringClass, args[0].Class().Name)
				}

				d, f := filepath.Split(fn.value)

				return t.vm.InitArrayObject([]Object{t.vm.InitStringObject(d), t.vm.InitStringObject(f)})

			},
		},
	}
}

// Instance methods -----------------------------------------------------
func builtinFileInstanceMethods() []*BuiltinMethodObject {
	return []*BuiltinMethodObject{
		{
			// Closes the instance of File class. Possible to close twice.
			//
			// ```ruby
			// File.open("/tmp/goby/out.txt", "w", 0755) do |f|
			//   f.close      # redundant: instance f will automatically close
			// end
			//
			// f = File.new("/tmp/goby/out.txt", "w", 0755)
			// f.close
			// f.close
			// ```
			//
			// @return [Null]
			Name: "close",
			Fn: func(receiver Object, sourceLine int, t *Thread, args []Object, blockFrame *normalCallFrame) Object {
				file := receiver.(*FileObject).File
				file.Close()

				return NULL

			},
		},
		// Returns the path and the file name.
		//
		// ```ruby
		// File.open("/tmp/goby/out.txt", "w", 0755) do |f|
		//   puts f.name      #=> "/tmp/goby/out.txt"
		// end
		// ```
		//
		// @return [String]
		{
			Name: "name",
			Fn: func(receiver Object, sourceLine int, t *Thread, args []Object, blockFrame *normalCallFrame) Object {
				name := receiver.(*FileObject).File.Name()
				return t.vm.InitStringObject(name)

			},
		},
		// Returns the contents of the specified file.
		//
		// ```ruby
		// File.open("/tmp/goby/out.txt", "w", 0755) do |f|
		//   f.write("Hello, Goby!")
		//   puts f.read      #=> "Hello, Goby!"
		// end
		// ```
		//
		// @return [String]
		{
			Name: "read",
			Fn: func(receiver Object, sourceLine int, t *Thread, args []Object, blockFrame *normalCallFrame) Object {
				var result string
				var f []byte
				var err error

				file := receiver.(*FileObject).File

				if file.Name() == "/dev/stdin" {
					reader := bufio.NewReader(os.Stdin)
					result, err = reader.ReadString('\n')
				} else {
					f, err = ioutil.ReadFile(file.Name())
					result = string(f)
				}

				if err != nil {
					return t.vm.InitErrorObject(errors.IOError, sourceLine, err.Error())
				}

				return t.vm.InitStringObject(result)

			},
		},
		{
			// Returns size of file in bytes.
			//
			// ```ruby
			// File.new("loop.gb").size # => 321123
			// ```
			//
			// @return [Integer]
			Name: "size",
			Fn: func(receiver Object, sourceLine int, t *Thread, args []Object, blockFrame *normalCallFrame) Object {
				file := receiver.(*FileObject).File

				fileStats, err := os.Stat(file.Name())
				if err != nil {
					return t.vm.InitErrorObject(errors.IOError, sourceLine, err.Error())
				}

				return t.vm.InitIntegerObject(int(fileStats.Size()))

			},
		},
		{
			Name: "write",
			Fn: func(receiver Object, sourceLine int, t *Thread, args []Object, blockFrame *normalCallFrame) Object {
				file := receiver.(*FileObject).File
				data := args[0].(*StringObject).value
				length, err := file.Write([]byte(data))

				if err != nil {
					return t.vm.InitErrorObject(errors.IOError, sourceLine, err.Error())
				}

				return t.vm.InitIntegerObject(length)

			},
		},
	}
}

// Internal functions ===================================================

// Functions for initialization -----------------------------------------

func (vm *VM) initFileObject(f *os.File) *FileObject {
	return &FileObject{
		BaseObj: &BaseObj{class: vm.TopLevelClass(classes.FileClass)},
		File:    f,
	}
}

func (vm *VM) initFileClass() *RClass {
	fc := vm.initializeClass(classes.FileClass)
	fc.setBuiltinMethods(builtinFileClassMethods(), true)
	fc.setBuiltinMethods(builtinFileInstanceMethods(), false)

	vm.libFiles = append(vm.libFiles, "file.gb")

	return fc
}

// Polymorphic helper functions -----------------------------------------

// ToString returns the object's name as the string format
func (f *FileObject) ToString() string {
	return "<File: " + f.File.Name() + ">"
}

// Inspect delegates to ToString
func (f *FileObject) Inspect() string {
	return f.ToString()
}

// ToJSON just delegates to `ToString`
func (f *FileObject) ToJSON(t *Thread) string {
	return f.ToString()
}

// Value returns file object's string format
func (f *FileObject) Value() interface{} {
	return f.File
}
