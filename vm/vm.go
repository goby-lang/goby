package vm

import (
	"fmt"
	"github.com/goby-lang/goby/compiler"
	"github.com/goby-lang/goby/compiler/bytecode"
	"github.com/goby-lang/goby/compiler/parser"
	"github.com/goby-lang/goby/vm/classes"
	"github.com/goby-lang/goby/vm/errors"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

// Version stores current Goby version
const Version = "0.1.3"

// These are the enums for marking parser's mode, which decides whether it should pop unused values.
const (
	NormalMode int = iota
	REPLMode
	TestMode
)

type isIndexTable struct {
	Data map[string]int
}

func newISIndexTable() *isIndexTable {
	return &isIndexTable{Data: make(map[string]int)}
}

type isTable map[string][]*instructionSet

type filename = string

type errorMessage = string

var standardLibraries = map[string]func(*VM){
	"net/http":          initHTTPClass,
	"net/simple_server": initSimpleServerClass,
	"uri":               initURIClass,
	"db":                initDBClass,
	"plugin":            initPluginClass,
	"json":              initJSONClass,
	"concurrent/array":  initConcurrentArrayClass,
	"concurrent/hash":   initConcurrentHashClass,
}

// VM represents a stack based virtual machine.
type VM struct {
	mainObj     *RObject
	mainThread  *thread
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

	stackTraceCount int

	channelObjectMap *objectMap

	sync.Mutex

	mode int

	libFiles []string
}

// New initializes a vm to initialize state and returns it.
func New(fileDir string, args []string) (vm *VM, e error) {
	vm = &VM{args: args}
	vm.mainThread = vm.newThread()

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

	vm.initConstants()
	vm.mainObj = vm.initMainObj()
	vm.channelObjectMap = &objectMap{store: &sync.Map{}}

	for _, fn := range vm.libFiles {
		vm.execGobyLib(fn)
	}

	return
}

func (vm *VM) newThread() *thread {
	s := &stack{RWMutex: new(sync.RWMutex)}
	cfs := &callFrameStack{callFrames: []callFrame{}}
	t := &thread{stack: s, callFrameStack: cfs, sp: 0, cfp: 0}
	s.thread = t
	cfs.thread = t
	t.vm = vm
	return t
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

	cf := newNormalCallFrame(translator.program, translator.filename)
	cf.self = vm.mainObj
	vm.mainThread.callFrameStack.push(cf)
	vm.startFromTopFrame()
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

// Start evaluation from top most call frame
func (vm *VM) startFromTopFrame() {
	vm.mainThread.startFromTopFrame()
}

func (vm *VM) currentFilePath() string {
	frame := vm.mainThread.callFrameStack.top()
	return frame.FileName()
}

func (vm *VM) getBlock(name string, filename filename) *instructionSet {
	// The "name" here is actually an index of block
	// for example <Block:1>'s name is "1"
	is, ok := vm.blockTables[filename][name]

	if !ok {
		panic(fmt.Sprintf("Can't find block %s", name))
	}

	return is
}

func (vm *VM) getMethodIS(name string, filename filename) (*instructionSet, bool) {
	iss, ok := vm.isTables[bytecode.MethodDef][name]

	if !ok {
		return nil, false
	}

	is := iss[vm.methodISIndexTables[filename].Data[name]]

	vm.methodISIndexTables[filename].Data[name]++

	return is, ok
}

func (vm *VM) getClassIS(name string, filename filename) *instructionSet {
	iss, ok := vm.isTables[bytecode.ClassDef][name]

	if !ok {
		panic(fmt.Sprintf("Can't find class %s's instructions", name))
	}

	is := iss[vm.classISIndexTables[filename].Data[name]]

	vm.classISIndexTables[filename].Data[name]++

	return is
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
		constant = namespace.lookupConstant(constName, true)

		if constant != nil {
			return
		}
	}

	constant = cf.lookupConstant(constName)

	if constant == nil {
		constant = vm.objectClass.constants[constName]
	}

	if constName == classes.ObjectClass {
		constant = &Pointer{Target: vm.objectClass}
	}

	return
}

func (vm *VM) execGobyLib(libName string) {
	libPath := filepath.Join(vm.projectRoot, "lib", libName)
	file, err := ioutil.ReadFile(libPath)

	if err != nil {
		vm.mainThread.pushErrorObject(errors.InternalError, -1, err.Error())
	}

	vm.execRequiredFile(libPath, file)
}

func (vm *VM) execRequiredFile(filepath string, file []byte) {
	instructionSets, err := compiler.CompileToInstructions(string(file), parser.NormalMode)

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	oldMethodTable := isTable{}
	oldClassTable := isTable{}

	// Copy current file's instruction sets.
	for name, is := range vm.isTables[bytecode.MethodDef] {
		oldMethodTable[name] = is
	}

	for name, is := range vm.isTables[bytecode.ClassDef] {
		oldClassTable[name] = is
	}

	// This creates new execution environments for required file, including new instruction set table.
	// So we need to copy old instruction sets and restore them later, otherwise current program's instruction set would be overwrite.
	vm.ExecInstructions(instructionSets, filepath)

	// Restore instruction sets.
	vm.isTables[bytecode.MethodDef] = oldMethodTable
	vm.isTables[bytecode.ClassDef] = oldClassTable
}

func newError(format string, args ...interface{}) *Error {
	return &Error{Message: fmt.Sprintf(format, args...)}
}
