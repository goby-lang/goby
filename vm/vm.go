package vm

import (
	"fmt"
	"github.com/goby-lang/goby/bytecode"
	"github.com/goby-lang/goby/parser"
	"io/ioutil"
	"path"
	"path/filepath"
	"strings"
)

var stackTrace int

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
}

// New initializes a vm to initialize state and returns it.
func New(fileDir string, args []string) *VM {
	vm := &VM{args: args}
	vm.mainThread = vm.newThread()

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

	p, _ := filepath.Abs("../")

	if !strings.HasSuffix(p, "goby") {
		vm.projectRoot = path.Join(p, "goby")
	} else {
		vm.projectRoot = p
	}

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

// ExecBytecodes accepts a sequence of bytecodes and use vm to evaluate them.
func (vm *VM) ExecBytecodes(bytecodes, fn string) {
	filename := filename(fn)
	p := newBytecodeParser(filename)
	p.vm = vm
	p.parseBytecode(bytecodes)

	// Keep update label table after parsed new files.
	// TODO: Find more efficient way to do this.
	for labelType, table := range p.labelTable {
		for labelName, is := range table {
			vm.isTables[labelType][labelName] = is
		}
	}

	vm.blockTables[p.filename] = p.blockTable
	vm.classISIndexTables[filename] = newISIndexTable()
	vm.methodISIndexTables[filename] = newISIndexTable()

	defer func() {
		if p := recover(); p != nil {
			switch p.(type) {
			case errorMessage:
				return
			default:
				panic(p)
			}
		}
	}()

	cf := newCallFrame(p.program)
	cf.self = mainObj
	vm.mainThread.callFrameStack.push(cf)
	vm.startFromTopFrame()
}

// GetExecResult returns stack's top most value. Normally it's used in tests.
func (vm *VM) GetExecResult() Object {
	return vm.mainThread.stack.top().Target
}

func (vm *VM) initConstants() {
	vm.constants = make(map[string]*Pointer)
	constants := make(map[string]*Pointer)

	builtInClasses := []Class{
		integerClass,
		stringClass,
		booleanClass,
		nullClass,
		arrayClass,
		hashClass,
		classClass,
		methodClass,
	}

	args := []Object{}

	for _, arg := range vm.args {
		args = append(args, initializeString(arg))
	}

	for _, c := range builtInClasses {
		p := &Pointer{Target: c}
		constants[c.ReturnName()] = p
	}

	constants["ARGV"] = &Pointer{Target: initializeArray(args)}
	objectClass.constants = constants
	vm.constants["Object"] = &Pointer{objectClass}
}

// Start evaluation from top most call frame
func (vm *VM) startFromTopFrame() {
	vm.mainThread.startFromTopFrame()
}

func (vm *VM) currentFilePath() string {
	return string(vm.mainThread.callFrameStack.top().instructionSet.filename)
}

func (vm *VM) printDebugInfo(i *instruction) {
	fmt.Println(i.inspect())
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

	vm.methodISIndexTables[filename].Data[name]++
	return is, ok
}

func (vm *VM) getClassIS(name string, filename filename) *instructionSet {
	iss, ok := vm.isTables[bytecode.LabelDefClass][name]

	if !ok {
		panic(fmt.Sprintf("Can't find class %s's instructions", name))
	}

	is := iss[vm.classISIndexTables[filename].Data[name]]

	vm.classISIndexTables[filename].Data[name]++
	return is
}

func (vm *VM) loadConstant(name string, isModule bool) *RClass {
	var c *RClass
	var ptr *Pointer

	ptr = objectClass.constants[name]

	if ptr == nil {
		c = initializeClass(name, isModule)
		objectClass.constants[name] = &Pointer{Target: c}
	} else {
		c = ptr.Target.(*RClass)
	}

	return c
}

func (vm *VM) lookupConstant(cf *callFrame, constName string) *Pointer {
	var constant *Pointer
	var namespace Class
	var hasNamespace bool

	top := vm.mainThread.stack.top()

	if top == nil {
		hasNamespace = false
	} else {
		namespace, hasNamespace = top.Target.(Class)
	}

	if hasNamespace {
		if namespace != cf.self {
			vm.mainThread.stack.pop()
		}

		constant = namespace.lookupConstant(constName, true)

		if constant != nil {
			return constant
		}
	}

	switch s := cf.self.(type) {
	case Class:
		constant = s.lookupConstant(constName, true)
		if constant != nil {
			return constant
		}
	default:
		c := s.returnClass()

		constant = c.lookupConstant(constName, true)
		if constant != nil {
			return constant
		}
	}

	constant = vm.constants[constName]
	return constant
}

func (vm *VM) execGobyLib(libName string) {
	libPath := path.Join(vm.projectRoot, "lib", libName)
	file, err := ioutil.ReadFile(libPath)

	if err != nil {
		vm.mainThread.returnError(err.Error())
	}

	vm.execRequiredFile(libPath, file)
}

func (vm *VM) execRequiredFile(filepath string, file []byte) {
	program := parser.BuildAST(file)
	g := bytecode.NewGenerator(program)
	bytecodes := g.GenerateByteCode(program)

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
	vm.ExecBytecodes(bytecodes, filepath)

	// Restore instruction sets.
	vm.isTables[bytecode.LabelDef] = oldMethodTable
	vm.isTables[bytecode.LabelDefClass] = oldClassTable
}

func newError(format string, args ...interface{}) *Error {
	return &Error{Message: fmt.Sprintf(format, args...)}
}
