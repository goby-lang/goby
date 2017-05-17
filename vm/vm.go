package vm

import (
	"bytes"
	"fmt"
	"github.com/goby-lang/goby/bytecode"
	"strings"
)

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

var stackTrace int

type isIndexTable struct {
	Data map[string]int
}

type isTable map[string][]*instructionSet

type filename string

func newISIndexTable() *isIndexTable {
	return &isIndexTable{Data: make(map[string]int)}
}

type errorMessage string

type stack struct {
	Data []*Pointer
	VM   *VM
}

type standardLibraryInitMethod func(*VM)

var standardLibraris = map[string]standardLibraryInitMethod{
	"file": initializeFileClass,
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

	//fmt.Println(i.inspect())
	i.action.operation(vm, cf, i.Params...)
}

func (vm *VM) printDebugInfo(i *instruction) {
	fmt.Println(i.inspect())
}

func (vm *VM) getBlock(name string, filename filename) (*instructionSet, bool) {
	// The "name" here is actually an index from label
	// for example <Block:1>'s name is "1"
	is, ok := vm.blockTables[filename][name]

	if !ok {
		return nil, false
	}

	return is, ok
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

func (vm *VM) getClassIS(name string, filename filename) (*instructionSet, bool) {
	iss, ok := vm.isTables[bytecode.LabelDefClass][name]

	if !ok {
		return nil, false
	}

	is := iss[vm.classISIndexTables[filename].Data[name]]

	vm.classISIndexTables[filename].Data[name]++
	return is, ok
}

func (vm *VM) lookupConstant(cf *callFrame, constName string) *Pointer {
	var constant *Pointer
	var namespace Class
	var hasNamespace bool
	var ok bool

	top := vm.stack.top()

	if top == nil || top.Target.objectType() != classObj {
		hasNamespace = false
	} else {
		namespace, hasNamespace = top.Target.(Class)
	}

	if hasNamespace {
		if namespace != cf.self {
			vm.stack.pop()
		}
		constant = namespace.lookupConstant(constName)

		if constant == nil {
			constant, ok = vm.constants[constName]
		}
	} else if scope, inClass := cf.self.(Class); inClass {
		constant = scope.lookupConstant(constName)

		if constant == nil {
			constant, ok = vm.constants[constName]
		}
	} else {
		constant, ok = vm.constants[constName]
	}

	if !ok {
		msg := "Can't find constant: " + constName
		vm.returnError(msg)
	}

	return constant
}

func (s *stack) push(v *Pointer) {
	if len(s.Data) <= s.VM.sp {
		s.Data = append(s.Data, v)
	} else {
		s.Data[s.VM.sp] = v
	}

	s.VM.sp++
}

func (s *stack) pop() *Pointer {
	if len(s.Data) < 1 {
		panic("Nothing to pop!")
	}

	s.VM.sp--

	v := s.Data[s.VM.sp]
	s.Data[s.VM.sp] = nil
	return v
}

func (s *stack) top() *Pointer {

	if len(s.Data) == 0 {
		return nil
	}

	if s.VM.sp > 0 {
		return s.Data[s.VM.sp-1]
	}

	return s.Data[0]
}

func (s *stack) inspect() string {
	var out bytes.Buffer
	datas := []string{}

	for i, p := range s.Data {
		if p != nil {
			o := p.Target
			if i == s.VM.sp {
				datas = append(datas, fmt.Sprintf("%s (%T) %d <----", o.Inspect(), o, i))
			} else {
				datas = append(datas, fmt.Sprintf("%s (%T) %d", o.Inspect(), o, i))
			}

		} else {
			if i == s.VM.sp {
				datas = append(datas, "nil <----")
			} else {
				datas = append(datas, "nil")
			}

		}

	}

	out.WriteString("-----------\n")
	out.WriteString(strings.Join(datas, "\n"))
	out.WriteString("\n---------\n")

	return out.String()
}

func newError(format string, args ...interface{}) *Error {
	return &Error{Message: fmt.Sprintf(format, args...)}
}

// TODO: Use this method to replace unnecessary panics
func (vm *VM) returnError(msg string) {
	err := &Error{Message: msg}
	vm.stack.push(&Pointer{err})
	panic(errorMessage(msg))
}
