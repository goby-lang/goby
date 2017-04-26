package vm

import (
	"bytes"
	"fmt"
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
	labelTables map[labelType]map[string][]*instructionSet
	// method instruction set table
	methodISTable *isIndexTable
	// class instruction set table
	classISTable *isIndexTable
	// block instruction set table
	blockList *isIndexTable

	fileDir string
}

type isIndexTable struct {
	Data map[string]int
}

type stack struct {
	Data []*Pointer
	VM   *VM
}

// New initializes a vm to initialize state and returns it.
func New(fileDir string) *VM {
	s := &stack{}
	cfs := &callFrameStack{callFrames: []*callFrame{}}
	vm := &VM{stack: s, callFrameStack: cfs, sp: 0, cfp: 0}
	s.VM = vm
	cfs.vm = vm

	vm.initConstants()
	vm.methodISTable = &isIndexTable{Data: make(map[string]int)}
	vm.classISTable = &isIndexTable{Data: make(map[string]int)}
	vm.blockList = &isIndexTable{Data: make(map[string]int)}
	vm.labelTables = map[labelType]map[string][]*instructionSet{
		LabelDef:      make(map[string][]*instructionSet),
		LabelDefClass: make(map[string][]*instructionSet),
		Block:         make(map[string][]*instructionSet),
		Program:       make(map[string][]*instructionSet),
	}
	vm.fileDir = fileDir
	return vm
}

// ExecBytecodes accepts a sequence of bytecodes and use vm to evaluate them.
func (vm *VM) ExecBytecodes(bytecodes string) {
	p := newBytecodeParser()
	p.vm = vm
	p.parseBytecode(bytecodes)

	// Keep update label table after parsed new files.
	// TODO: Find more efficient way to do this.
	for labelType, table := range p.labelTable {
		for labelName, is := range table {
			vm.labelTables[labelType][labelName] = is
		}
	}

	cf := newCallFrame(vm.labelTables[Program]["ProgramStart"][0])
	cf.self = mainObj
	vm.callFrameStack.push(cf)
	vm.start()
}

// GetExecResult returns stack's top most value. Normally it's used in tests.
func (vm *VM) GetExecResult() Object {
	return vm.stack.Top().Target
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

	for _, c := range builtInClasses {
		p := &Pointer{Target: c}
		constants[c.ReturnName()] = p
	}

	vm.constants = constants
}

func (vm *VM) evalCallFrame(cf *callFrame) {
	for cf.pc < len(cf.instructionSet.instructions) {
		i := cf.instructionSet.instructions[cf.pc]
		vm.execInstruction(cf, i)
	}
}

// Start evaluation from top most call frame
func (vm *VM) start() {
	cf := vm.callFrameStack.top()
	vm.evalCallFrame(cf)
}

func (vm *VM) execInstruction(cf *callFrame, i *instruction) {
	cf.pc++
	//fmt.Print(i.Inspect())
	i.action.operation(vm, cf, i.Params...)
	//fmt.Println(vm.callFrameStack.inspect())
	//fmt.Println(vm.stack.inspect())
}

func (vm *VM) getBlock(name string) (*instructionSet, bool) {
	// The "name" here is actually an index from label
	// for example <Block:1>'s name is "1"
	iss, ok := vm.labelTables[Block][name]

	if !ok {
		return nil, false
	}

	is := iss[0]

	return is, ok
}

func (vm *VM) getMethodIS(name string) (*instructionSet, bool) {
	iss, ok := vm.labelTables[LabelDef][name]

	if !ok {
		return nil, false
	}

	is := iss[vm.methodISTable.Data[name]]

	vm.methodISTable.Data[name]++
	return is, ok
}

func (vm *VM) getClassIS(name string) (*instructionSet, bool) {
	iss, ok := vm.labelTables[LabelDefClass][name]

	if !ok {
		return nil, false
	}

	is := iss[vm.classISTable.Data[name]]

	vm.classISTable.Data[name]++
	return is, ok
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

func (s *stack) Top() *Pointer {

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
