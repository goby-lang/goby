package vm

import (
	"fmt"
	"bytes"
)

type CallFrameStack struct {
	CallFrames []*CallFrame
	VM         *VM
}

type CallFrame struct {
	InstructionSet *InstructionSet
	PC             int
	EP             []*Pointer
	Self           BaseObject
	Local          []*Pointer
	LPr            int
	IsBlock        bool
	BlockFrame     *CallFrame
}

func (cf *CallFrame) insertLCL(i int, value Object) {
	index := i
	existedLCL := cf.getLCL(i)

	if existedLCL != nil {
		existedLCL.Target = value
		return
	}

	cf.Local = append(cf.Local, nil)
	copy(cf.Local[index:], cf.Local[index:])
	cf.Local[index] = &Pointer{Target: value}

	if index >= cf.LPr {
		cf.LPr = index + 1
	}
}

func (cf *CallFrame) getLCL(index int) *Pointer {
	var p *Pointer

	p = cf.Local[index]

	for p == nil {
		if cf.BlockFrame != nil {
			p = getLCLFromEP(cf.BlockFrame, index)

			if p == nil {
				return nil
			}
		} else {
			return nil
		}
	}

	return p
}

func (cf *CallFrame) inspect() string {
	return fmt.Sprintf("Name: %s. is block: %t, lens of EP: %d", cf.InstructionSet.Label.Name, cf.IsBlock, len(cf.EP))
}

func getLCLFromEP(cf *CallFrame, index int) *Pointer {
	var v *Pointer

	v = cf.EP[index]

	if v != nil {
		return v
	}

	if cf.BlockFrame != nil {

		return getLCLFromEP(cf.BlockFrame, index)
	}

	return nil
}

func (cfs *CallFrameStack) Push(cf *CallFrame) {
	if cf == nil {
		panic("Callfame can't be nil!")
	}

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

func (cfs *CallFrameStack) inspect() string {
	var out bytes.Buffer

	for _, cf := range cfs.CallFrames {
		if cf != nil {
			out.WriteString(fmt.Sprintln(cf.inspect()))
		}
	}

	return out.String()
}
func NewCallFrame(is *InstructionSet) *CallFrame {
	return &CallFrame{Local: make([]*Pointer, 100), InstructionSet: is, PC: 0, LPr: 0}
}
