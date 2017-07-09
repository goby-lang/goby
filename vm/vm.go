package vm

import (
	"fmt"
	"github.com/goby-lang/goby/compiler"
	"github.com/goby-lang/goby/compiler/bytecode"
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"
)

type isIndexTable struct {
	Data map[string]int
}

func newISIndexTable() *isIndexTable {
	return &isIndexTable{Data: make(map[string]int)}
}

type isTable map[string][]*instructionSet

type filename string

type errorMessage string

var standardLibraries = map[string]func(*VM){
	"file":              initializeFileClass,
	"net/http":          initializeHTTPClass,
	"net/simple_server": initializeSimpleServerClass,
	"uri":               initializeURIClass,
}

// VM represents a stack based virtual machine.
type VM struct {
	mainObj    *RObject
	mainThread *thread
	// a map holds pointers of constants
	constants map[string]*Pointer
	// a map holds different types of label tables
	isTables map[labelType]isTable
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

	replMode bool

	stackTraceCount int

	channelObjectMap *objectMap

	sync.Mutex
}

// New initializes a vm to initialize state and returns it.
func New(fileDir string, args []string) *VM {
	vm := &VM{args: args}
	vm.mainThread = vm.newThread()
	vm.constants = make(map[string]*Pointer)

	vm.initConstants()
	vm.methodISIndexTables = map[filename]*isIndexTable{
		filename(fileDir): newISIndexTable(),
	}
	vm.classISIndexTables = map[filename]*isIndexTable{
		filename(fileDir): newISIndexTable(),
	}
	vm.blockTables = make(map[filename]map[string]*instructionSet)
	vm.isTables = map[labelType]isTable{
		bytecode.LabelDef:      make(isTable),
		bytecode.LabelDefClass: make(isTable),
	}
	vm.fileDir = fileDir
	vm.projectRoot = os.Getenv("GOBY_ROOT")
	vm.mainObj = vm.initMainObj()
	vm.channelObjectMap = &objectMap{store: map[int]Object{}}

	return vm
}

func (vm *VM) newThread() *thread {
	s := &stack{}
	cfs := &callFrameStack{callFrames: []*callFrame{}}
	t := &thread{stack: s, callFrameStack: cfs, sp: 0, cfp: 0}
	s.thread = t
	cfs.thread = t
	t.vm = vm
	return t
}

// ExecInstructions accepts a sequence of bytecodes and use vm to evaluate them.
func (vm *VM) ExecInstructions(sets []*bytecode.InstructionSet, fn string) {
	filename := filename(fn)
	p := newInstructionTranslator(filename)
	p.vm = vm
	p.transferInstructionSets(sets)

	// Keep update label table after parsed new files.
	// TODO: Find more efficient way to do this.
	for labelType, table := range p.labelTable {
		for labelName, is := range table {
			vm.isTables[labelType][labelName] = is
		}
	}

	vm.blockTables[p.filename] = p.blockTable
	vm.SetClassISIndexTable(p.filename)
	vm.SetMethodISIndexTable(p.filename)

	cf := newCallFrame(p.program)
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

func (vm *VM) initMainObj() *RObject {
	return vm.constants[objectClass].Target.(*RClass).initializeInstance()
}

func (vm *VM) initConstants() {
	cClass := initClassClass()
	objClass := initObjectClass(cClass)

	vm.constants[objectClass] = &Pointer{objClass}
	vm.topLevelClass(objectClass).setClassConstant(cClass)

	builtInClasses := []*RClass{
		vm.initIntegerClass(),
		vm.initStringClass(),
		vm.initBoolClass(),
		vm.initNullClass(),
		vm.initArrayClass(),
		vm.initHashClass(),
		vm.initRangeClass(),
		vm.initMethodClass(),
		vm.initChannelClass(),
	}

	vm.initErrorClasses()

	for _, c := range builtInClasses {
		objClass.setClassConstant(c)
	}

	args := []Object{}

	for _, arg := range vm.args {
		args = append(args, vm.initStringObject(arg))
	}

	objClass.constants["ARGV"] = &Pointer{Target: vm.initArrayObject(args)}
}

func (vm *VM) topLevelClass(cn string) *RClass {
	objClass := vm.constants[objectClass].Target.(*RClass)

	if cn == objectClass {
		return objClass
	}

	return objClass.constants[cn].Target.(*RClass)
}

// Start evaluation from top most call frame
func (vm *VM) startFromTopFrame() {
	vm.mainThread.startFromTopFrame()
}

func (vm *VM) currentFilePath() string {
	return string(vm.mainThread.callFrameStack.top().instructionSet.filename)
}

func (vm *VM) getBlock(name string, filename filename) *instructionSet {
	// The "name" here is actually an index from label
	// for example <Block:1>'s name is "1"
	is, ok := vm.blockTables[filename][name]

	if !ok {
		panic(fmt.Sprintf("Can't find block %s", name))
	}

	return is
}

func (vm *VM) getMethodIS(name string, filename filename) (*instructionSet, bool) {
	iss, ok := vm.isTables[bytecode.LabelDef][name]

	if !ok {
		return nil, false
	}

	is := iss[vm.methodISIndexTables[filename].Data[name]]

	if !vm.replMode {
		vm.methodISIndexTables[filename].Data[name]++
	}

	return is, ok
}

func (vm *VM) getClassIS(name string, filename filename) *instructionSet {
	iss, ok := vm.isTables[bytecode.LabelDefClass][name]

	if !ok {
		panic(fmt.Sprintf("Can't find class %s's instructions", name))
	}

	is := iss[vm.classISIndexTables[filename].Data[name]]

	if !vm.replMode {
		vm.classISIndexTables[filename].Data[name]++
	}

	return is
}

// loadConstant makes sure we don't create a class twice.
func (vm *VM) loadConstant(name string, isModule bool) *RClass {
	var c *RClass
	var ptr *Pointer

	ptr = vm.topLevelClass(objectClass).constants[name]

	if ptr == nil {
		c = vm.initializeClass(name, isModule)
		vm.topLevelClass(objectClass).setClassConstant(c)
	} else {
		c = ptr.Target.(*RClass)
	}

	return c
}

func (vm *VM) lookupConstant(cf *callFrame, constName string) (constant *Pointer) {
	var namespace *RClass
	var hasNamespace bool

	top := vm.mainThread.stack.top()

	if top == nil {
		hasNamespace = false
	} else {
		namespace, hasNamespace = top.Target.(*RClass)
	}

	if hasNamespace {
		// pop namespace since we don't need it anymore
		if namespace != cf.self {
			vm.mainThread.stack.pop()
		}

		constant = namespace.lookupConstant(constName, true)

		if constant != nil {
			return
		}
	}

	constant = cf.lookupConstant(constName)

	if constant == nil {
		constant = vm.constants[constName]
	}

	return
}

func (vm *VM) execGobyLib(libName string) {
	libPath := filepath.Join(vm.projectRoot, "lib", libName)
	file, err := ioutil.ReadFile(libPath)

	if err != nil {
		vm.mainThread.returnError(err.Error())
	}

	vm.execRequiredFile(libPath, file)
}

func (vm *VM) execRequiredFile(filepath string, file []byte) {
	instructionSets, err := compiler.CompileToInstructions(string(file))

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	oldMethodTable := isTable{}
	oldClassTable := isTable{}

	// Copy current file's instruction sets.
	for name, is := range vm.isTables[bytecode.LabelDef] {
		oldMethodTable[name] = is
	}

	for name, is := range vm.isTables[bytecode.LabelDefClass] {
		oldClassTable[name] = is
	}

	// This creates new execution environments for required file, including new instruction set table.
	// So we need to copy old instruction sets and restore them later, otherwise current program's instruction set would be overwrite.
	vm.ExecInstructions(instructionSets, filepath)

	// Restore instruction sets.
	vm.isTables[bytecode.LabelDef] = oldMethodTable
	vm.isTables[bytecode.LabelDefClass] = oldClassTable
}

func newError(format string, args ...interface{}) *Error {
	return &Error{Message: fmt.Sprintf(format, args...)}
}
