package vm

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"sync/atomic"
	"testing"

	"github.com/goby-lang/goby/compiler"
	"github.com/goby-lang/goby/compiler/bytecode"
	"github.com/goby-lang/goby/compiler/parser"
	"github.com/goby-lang/goby/vm/classes"
)

// Version stores current Goby version
const Version = "0.1.11"

// DefaultLibPath is used for overriding vm.libpath build-time.
var DefaultLibPath string

type isIndexTable struct {
	Data map[string]int
}

func newISIndexTable() *isIndexTable {
	return &isIndexTable{Data: make(map[string]int)}
}

type isTable map[string][]*instructionSet

type filename = string

var standardLibraries = map[string]func(*VM){
	"net/http":           initHTTPClass,
	"net/simple_server":  initSimpleServerClass,
	"uri":                initURIClass,
	"json":               initJSONClass,
	"concurrent/array":   initConcurrentArrayClass,
	"concurrent/hash":    initConcurrentHashClass,
	"concurrent/rw_lock": initConcurrentRWLockClass,
	"spec":               initSpecClass,
}

// VM represents a stack based virtual machine.
type VM struct {
	mainObj     *RObject
	mainThread  Thread
	objectClass *RClass
	// a map holds different types of instruction set tables
	isTables map[setType]isTable
	// method instruction set table
	methodISIndexTables map[filename]*isIndexTable
	// class instruction set table
	classISIndexTables map[filename]*isIndexTable
	// block instruction set table
	blockTables map[filename]map[string]*instructionSet
	// fileDir indicates executed file's directory
	fileDir string
	// args are command line arguments
	args []string
	// projectRoot is goby root's absolute path, which is $GOROOT/src/github.com/goby-lang/goby
	projectRoot string

	// libPath indicates the Goby (.gb) libraries path. Defaults to `<projectRoot>/lib`, unless
	// DefaultLibPath is specified.
	libPath string

	channelObjectMap *objectMap

	mode parser.ParserMode

	libFiles []string

	threadCount int64
}

// New initializes a vm to initialize state and returns it.
func New(fileDir string, args []string) (vm *VM, e error) {
	vm = &VM{args: args}
	vm.mainThread.vm = vm
	vm.threadCount++
	vm.mode = parser.NormalMode

	vm.methodISIndexTables = map[filename]*isIndexTable{
		fileDir: newISIndexTable(),
	}
	vm.classISIndexTables = map[filename]*isIndexTable{
		fileDir: newISIndexTable(),
	}
	vm.blockTables = make(map[filename]map[string]*instructionSet)
	vm.isTables = map[setType]isTable{
		bytecode.MethodDef: make(isTable),
		bytecode.ClassDef:  make(isTable),
	}
	vm.fileDir = fileDir

	gobyRoot := os.Getenv("GOBY_ROOT")

	if len(gobyRoot) == 0 {
		vm.projectRoot = fmt.Sprintf("/usr/local/Cellar/goby/%s", Version)

		_, err := os.Stat(vm.projectRoot)

		if err != nil {
			gp := os.Getenv("GOPATH")
			path, _ := filepath.Abs(gp + "/src/github.com/goby-lang/goby")
			_, err = os.Stat(path)

			if err != nil {
				e = fmt.Errorf("You haven't set $GOBY_ROOT properly")
				return nil, e
			}

			vm.projectRoot = path
		}
	} else {
		vm.projectRoot = gobyRoot
	}

	if DefaultLibPath != "" {
		vm.libPath = DefaultLibPath
	} else {
		vm.libPath = filepath.Join(vm.projectRoot, "lib")
	}

	vm.initConstants()
	vm.mainObj = vm.initMainObj()
	vm.channelObjectMap = &objectMap{store: &sync.Map{}}

	for _, fn := range vm.libFiles {
		err := vm.mainThread.execGobyLib(fn)
		if err != nil {
			fmt.Printf("An error occurs when loading lib file %s:\n", string(fn))
			fmt.Println(err.Error())
			break
		}
	}

	return
}

func (vm *VM) newThread() (t Thread) {
	t.vm = vm
	t.id = atomic.AddInt64(&vm.threadCount, 1)
	return
}

// ExecInstructions accepts a sequence of bytecodes and use vm to evaluate them.
func (vm *VM) ExecInstructions(sets []*bytecode.InstructionSet, fn string) {
	translator := newInstructionTranslator(fn)
	translator.vm = vm
	translator.transferInstructionSets(sets)

	// Keep instruction set table updated after parsed new files.
	// TODO: Find more efficient way to do this.
	for setType, table := range translator.setTable {
		for name, is := range table {
			vm.isTables[setType][name] = is
		}
	}

	vm.blockTables[translator.filename] = translator.blockTable
	vm.SetClassISIndexTable(translator.filename)
	vm.SetMethodISIndexTable(translator.filename)

	cf := newNormalCallFrame(translator.program, translator.filename, 1)
	cf.self = vm.mainObj
	vm.mainThread.callFrameStack.push(cf)

	defer func() {
		switch err := recover().(type) {
		case error:
			panic(err)
		case *Error:
			if vm.mode == parser.NormalMode {
				fmt.Fprintln(os.Stderr, err.Message())
				os.Exit(1)
			}
		}
	}()

	vm.mainThread.startFromTopFrame()
}

// SetClassISIndexTable adds new instruction set's index table to vm.classISIndexTables
func (vm *VM) SetClassISIndexTable(fn filename) {
	vm.classISIndexTables[fn] = newISIndexTable()
}

// SetMethodISIndexTable adds new instruction set's index table to vm.methodISIndexTables
func (vm *VM) SetMethodISIndexTable(fn filename) {
	vm.methodISIndexTables[fn] = newISIndexTable()
}

// main object singleton methods -----------------------------------------------------
func builtinMainObjSingletonMethods() []*BuiltinMethodObject {
	return []*BuiltinMethodObject{
		{
			Name: "to_s",
			Fn: func(receiver Object, sourceLine int, thread *Thread, objects []Object, frame *normalCallFrame) Object {
				return thread.vm.InitStringObject("main")

			},
		},
	}
}

func (vm *VM) initMainObj() *RObject {
	obj := vm.objectClass.initializeInstance()
	singletonClass := vm.initializeClass(fmt.Sprintf("#<Class:%s>", obj.ToString()))
	singletonClass.Methods.set("include", vm.TopLevelClass(classes.ClassClass).lookupMethod("include"))
	singletonClass.setBuiltinMethods(builtinMainObjSingletonMethods(), false)
	obj.singletonClass = singletonClass

	return obj
}

func (vm *VM) initConstants() {
	// Init Class and Object
	cClass := initClassClass()
	mClass := initModuleClass(cClass)
	vm.objectClass = initObjectClass(cClass)
	vm.TopLevelClass(classes.ObjectClass).setClassConstant(cClass)
	vm.TopLevelClass(classes.ObjectClass).setClassConstant(mClass)

	// Init builtin classes
	builtinClasses := []*RClass{
		vm.initIntegerClass(),
		vm.initFloatClass(),
		vm.initStringClass(),
		vm.initBoolClass(),
		vm.initNullClass(),
		vm.initArrayClass(),
		vm.initHashClass(),
		vm.initRangeClass(),
		vm.initMethodClass(),
		vm.initBlockClass(),
		vm.initChannelClass(),
		vm.initGoClass(),
		vm.initFileClass(),
		vm.initRegexpClass(),
		vm.initMatchDataClass(),
		vm.initGoMapClass(),
		vm.initDecimalClass(),
	}

	// Init error classes
	vm.initErrorClasses()

	for _, c := range builtinClasses {
		vm.objectClass.setClassConstant(c)
	}

	// Init ARGV
	args := []Object{}

	for _, arg := range vm.args {
		args = append(args, vm.InitStringObject(arg))
	}

	vm.objectClass.constants["ARGV"] = &Pointer{Target: vm.InitArrayObject(args)}

	// Init ENV
	envs := map[string]Object{}

	for _, e := range os.Environ() {
		pair := strings.Split(e, "=")
		envs[pair[0]] = vm.InitStringObject(pair[1])
	}

	vm.objectClass.constants["ENV"] = &Pointer{Target: vm.InitHashObject(envs)}
	vm.objectClass.constants["STDOUT"] = &Pointer{Target: vm.initFileObject(os.Stdout)}
	vm.objectClass.constants["STDERR"] = &Pointer{Target: vm.initFileObject(os.Stderr)}
	vm.objectClass.constants["STDIN"] = &Pointer{Target: vm.initFileObject(os.Stdin)}
}

// TopLevelClass returns a specified top-level class (stored under the Object constant)
func (vm *VM) TopLevelClass(cn string) *RClass {
	objClass := vm.objectClass

	if cn == classes.ObjectClass {
		return objClass
	}

	return objClass.constants[cn].Target.(*RClass)
}

func (vm *VM) currentFilePath() string {
	frame := vm.mainThread.callFrameStack.top()
	return frame.FileName()
}

// loadConstant makes sure we don't create a class twice.
func (vm *VM) loadConstant(name string, isModule bool) *RClass {
	var c *RClass
	var ptr *Pointer

	ptr = vm.objectClass.constants[name]

	if ptr == nil {
		if isModule {
			c = vm.initializeClass(name)
		} else {
			c = vm.initializeModule(name)
		}

		vm.objectClass.setClassConstant(c)
	} else {
		c = ptr.Target.(*RClass)
	}

	return c
}

func (vm *VM) lookupConstant(cf callFrame, constName string) (constant *Pointer) {
	var namespace *RClass
	var hasNamespace bool

	top := vm.mainThread.Stack.top()

	if top == nil {
		hasNamespace = false
	} else {
		namespace, hasNamespace = top.Target.(*RClass)
	}

	if hasNamespace {
		constant = namespace.lookupConstantUnderAllScope(constName)

		if constant != nil {
			return
		}
	}

	constant = cf.lookupConstantUnderAllScope(constName)

	if constant == nil {
		constant = vm.objectClass.constants[constName]
	}

	if constName == classes.ObjectClass {
		constant = &Pointer{Target: vm.objectClass}
	}

	return
}
func initTestVM() *VM {
	fn, err := os.Getwd()

	if err != nil {
		panic(err)
	}

	v, err := New(fn, []string{})

	if err != nil {
		panic(err)
	}

	v.mode = parser.TestMode
	return v
}

func getFilename() string {
	_, filename, _, _ := runtime.Caller(1)
	return filename
}

// ExecAndReturn is a test helper
func ExecAndReturn(t *testing.T, src string) Object {
	t.Helper()
	v := initTestVM()
	return v.testEval(t, src, getFilename())
}

func (vm *VM) testEval(t *testing.T, input, filepath string) Object {
	iss, err := compiler.CompileToInstructions(input, parser.TestMode)

	if err != nil {
		t.Helper()
		t.Errorf("Error when compiling input: %s", input)
		t.Fatal(err.Error())
	}

	vm.ExecInstructions(iss, filepath)

	return vm.mainThread.Stack.top().Target
}
