package vm

type CallFrameStack struct {
	CallFrames []*CallFrame
	VM         *VM
}

type CallFrame struct {
	InstructionSet *InstructionSet
	PC             int
	Self           BaseObject
	Local          []Object
	LPr            int
	ArgPr          int
}

func (cf *CallFrame) insertLCL(i int, value Object) {
	index := i
	cf.Local = append(cf.Local, nil)
	copy(cf.Local[index:], cf.Local[index:])
	cf.Local[index] = value

	if index >= cf.LPr {
		cf.LPr = index + 1
	}
}

func (cf *CallFrame) getLCL(index *IntegerObject) Object {

	return cf.Local[index.Value]
}

func (cfs *CallFrameStack) Push(cf *CallFrame) {
	if len(cfs.CallFrames) <= cfs.VM.CFP {
		cfs.CallFrames = append(cfs.CallFrames, cf)
	} else {
		cfs.CallFrames[cfs.VM.CFP] = cf
	}

	cfs.VM.CFP += 1
}

func (cfs *CallFrameStack) Pop() *CallFrame {
	if len(cfs.CallFrames) < 1 {
		panic("Nothing to pop!")
	}

	if cfs.VM.CFP > 0 {
		cfs.VM.CFP -= 1
	}

	cf := cfs.CallFrames[cfs.VM.CFP]
	cfs.CallFrames[cfs.VM.CFP] = nil
	return cf
}

func (cfs *CallFrameStack) Top() *CallFrame {
	if cfs.VM.CFP > 0 {
		return cfs.CallFrames[cfs.VM.CFP-1]
	}

	return nil
}

func NewCallFrame(is *InstructionSet) *CallFrame {
	return &CallFrame{Local: make([]Object, 100), InstructionSet: is, PC: 0, LPr: 0}
}
