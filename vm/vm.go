package vm

import (
	"fmt"
	"github.com/goby-lang/goby/bytecode"
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
	// a stack that holds call frames
	callFrameStack *callFrameStack
	// call frame pointer
	cfp int
	// data stack
	stack *stack
	// stack pointer
	sp int
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

	fileDir string

	args []string
}

// New initializes a vm to initialize state and returns it.
func New(fileDir string, args []string) *VM {
	s := &stack{}
	cfs := &callFrameStack{callFrames: []*callFrame{}}
	vm := &VM{stack: s, callFrameStack: cfs, sp: 0, cfp: 0, args: args}
	s.VM = vm
	cfs.vm = vm

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
	return vm
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
	vm.callFrameStack.push(cf)
	vm.startFromTopFrame()
}

// GetExecResult returns stack's top most value. Normally it's used in tests.
func (vm *VM) GetExecResult() Object {
	return vm.stack.top().Target
}

func (vm *VM) initConstants() {
	constants := make(map[string]*Pointer)

	builtInClasses := []Class{
		integerClass,
		stringClass,
		booleanClass,
		nullClass,
		arrayClass,
		hashClass,
		classClass,
		objectClass,
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

	vm.constants = constants
}

func (vm *VM) evalCallFrame(cf *callFrame) {
	for cf.pc < len(cf.instructionSet.instructions) {
		i := cf.instructionSet.instructions[cf.pc]
		vm.execInstruction(cf, i)
	}
}

// Start evaluation from top most call frame
func (vm *VM) startFromTopFrame() {
	cf := vm.callFrameStack.top()
	vm.evalCallFrame(cf)
}

func (vm *VM) execInstruction(cf *callFrame, i *instruction) {
	cf.pc++

	defer func() {
		if p := recover(); p != nil {
			if stackTrace == 0 {
				fmt.Printf("Internal Error: %s\n", p)
			}
			fmt.Printf("Instruction trace: %d. \"%s\"\n", stackTrace, i.inspect())
			stackTrace++
			panic(p)
		}
	}()

	i.action.operation(vm, cf, i.Params...)
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

	if vm != nil {
		ptr = vm.constants[name]
	}

	if ptr == nil {
		c = initializeClass(name, isModule)
		vm.constants[name] = &Pointer{Target: c}
	} else {
		c = ptr.Target.(*RClass)
	}

	return c
}

func (vm *VM) lookupConstant(cf *callFrame, constName string) *Pointer {
	var constant *Pointer
	var namespace Class
	var hasNamespace bool

	top := vm.stack.top()

	if top == nil {
		hasNamespace = false
	} else {
		namespace, hasNamespace = top.Target.(Class)
	}

	if hasNamespace {
		if namespace != cf.self {
			vm.stack.pop()
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

// builtInMethodYield is like invokeblock instruction for built in methods
func (vm *VM) builtInMethodYield(blockFrame *callFrame, args ...Object) *Pointer {
	c := newCallFrame(blockFrame.instructionSet)
	c.blockFrame = blockFrame
	c.ep = blockFrame.ep
	c.self = blockFrame.self

	for i := 0; i < len(args); i++ {
		c.insertLCL(i, 0, args[i])
	}

	vm.callFrameStack.push(c)
	vm.startFromTopFrame()

	return vm.stack.top()
}

// TODO: Use this method to replace unnecessary panics
func (vm *VM) returnError(msg string) {
	err := &Error{Message: msg}
	vm.stack.push(&Pointer{err})
	panic(errorMessage(msg))
}

func newError(format string, args ...interface{}) *Error {
	return &Error{Message: fmt.Sprintf(format, args...)}
}
