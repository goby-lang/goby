package vm

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"sync/atomic"

	"github.com/goby-lang/goby/compiler/bytecode"
	"github.com/goby-lang/goby/vm/classes"
)

// Version stores current Goby version
const Version = "0.1.9"

// These are the enums for marking parser's mode, which decides whether it should pop unused values.
const (
	NormalMode int = iota
	REPLMode
	TestMode
)

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
	"db":                 initDBClass,
	"plugin":             initPluginClass,
	"json":               initJSONClass,
	"concurrent/array":   initConcurrentArrayClass,
	"concurrent/hash":    initConcurrentHashClass,
	"concurrent/rw_lock": initConcurrentRWLockClass,
	"spec":               initSpecClass,
}

// VM represents a stack based virtual machine.
type VM struct {
	mainObj     *RObject
	mainThread  thread
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

	mode int

	libFiles []string

	threadCount int64
}

// New initializes a vm to initialize state and returns it.
func New(fileDir string, args []string) (vm *VM, e error) {
	vm = &VM{args: args}
	vm.mainThread.vm = vm
	vm.threadCount++

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
			path, _ := filepath.Abs("$GOPATH/src/github.com/goby-lang/goby")
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

func (vm *VM) newThread() (t thread) {
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
		err, ok := recover().(*Error)

		if ok && vm.mode == NormalMode {
			fmt.Println(err.Message())
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
			Fn: func(receiver Object, sourceLine int) builtinMethodBody {
				return func(thread *thread, objects []Object, frame *normalCallFrame) Object {
					return thread.vm.initStringObject("main")
				}
			},
		},
	}
}

func (vm *VM) initMainObj() *RObject {
	obj := vm.objectClass.initializeInstance()
	singletonClass := vm.initializeClass(fmt.Sprintf("#<Class:%s>", obj.toString()), false)
	singletonClass.Methods.set("include", vm.topLevelClass(classes.ClassClass).lookupMethod("include"))
	singletonClass.setBuiltinMethods(builtinMainObjSingletonMethods(), false)
	obj.singletonClass = singletonClass

	return obj
}

func (vm *VM) initConstants() {
	// Init Class and Object
	cClass := initClassClass()
	vm.objectClass = initObjectClass(cClass)
	vm.topLevelClass(classes.ObjectClass).setClassConstant(cClass)

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
		args = append(args, vm.initStringObject(arg))
	}

	vm.objectClass.constants["ARGV"] = &Pointer{Target: vm.initArrayObject(args)}

	// Init ENV
	envs := map[string]Object{}

	for _, e := range os.Environ() {
		pair := strings.Split(e, "=")
		envs[pair[0]] = vm.initStringObject(pair[1])
	}

	vm.objectClass.constants["ENV"] = &Pointer{Target: vm.initHashObject(envs)}
	vm.objectClass.constants["STDOUT"] = &Pointer{Target: vm.initFileObject(os.Stdout)}
	vm.objectClass.constants["STDERR"] = &Pointer{Target: vm.initFileObject(os.Stderr)}
	vm.objectClass.constants["STDIN"] = &Pointer{Target: vm.initFileObject(os.Stdin)}
}

func (vm *VM) topLevelClass(cn string) *RClass {
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
		c = vm.initializeClass(name, isModule)
		vm.objectClass.setClassConstant(c)
	} else {
		c = ptr.Target.(*RClass)
	}

	return c
}

func (vm *VM) lookupConstant(cf callFrame, constName string) (constant *Pointer) {
	var namespace *RClass
	var hasNamespace bool

	top := vm.mainThread.stack.top()

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
