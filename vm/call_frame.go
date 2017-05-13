package vm

import (
	"bytes"
	"fmt"
)

type callFrameStack struct {
	callFrames []*callFrame
	vm         *VM
}

type callFrame struct {
	instructionSet *instructionSet
	// program counter
	pc int
	// environment pointer, points to the call frame we want to get locals from
	ep     *callFrame
	self   Object
	locals []*Pointer
	// local pointer
	lPr        int
	isBlock    bool
	blockFrame *callFrame
}

func (cf *callFrame) insertLCL(index, depth int, value Object) {
	existedLCL := cf.getLCL(index, depth)

	if existedLCL != nil {
		existedLCL.Target = value
		return
	}

	cf.locals = append(cf.locals, nil)
	copy(cf.locals[index:], cf.locals[index:])
	cf.locals[index] = &Pointer{Target: value}

	if index >= cf.lPr {
		cf.lPr = index + 1
	}
}

func (cf *callFrame) getLCL(index, depth int) *Pointer {
	if depth == 0 {
		return cf.locals[index]
	}

	return cf.blockFrame.ep.getLCL(index, depth-1)
}

func (cf *callFrame) inspect() string {
	if cf.ep != nil {
		return fmt.Sprintf("Name: %s. is block: %t. ep: %d", cf.instructionSet.label.name, cf.isBlock, len(cf.ep.locals))
	}
	return fmt.Sprintf("Name: %s. is block: %t", cf.instructionSet.label.name, cf.isBlock)
}

func (cfs *callFrameStack) push(cf *callFrame) {
	if cf == nil {
		panic("Callframe can't be nil!")
	}

	if len(cfs.callFrames) <= cfs.vm.cfp {
		cfs.callFrames = append(cfs.callFrames, cf)
	} else {
		cfs.callFrames[cfs.vm.cfp] = cf
	}

	cfs.vm.cfp++
}

func (cfs *callFrameStack) pop() *callFrame {
	if len(cfs.callFrames) < 1 {
		panic("Nothing to pop!")
	}

	if cfs.vm.cfp > 0 {
		cfs.vm.cfp--
	}

	cf := cfs.callFrames[cfs.vm.cfp]
	cfs.callFrames[cfs.vm.cfp] = nil
	return cf
}

func (cfs *callFrameStack) top() *callFrame {
	if cfs.vm.cfp > 0 {
		return cfs.callFrames[cfs.vm.cfp-1]
	}

	return nil
}

func (cfs *callFrameStack) inspect() string {
	var out bytes.Buffer

	for _, cf := range cfs.callFrames {
		if cf != nil {
			out.WriteString(fmt.Sprintln(cf.inspect()))
		}
	}

	return out.String()
}

func newCallFrame(is *instructionSet) *callFrame {
	return &callFrame{locals: make([]*Pointer, 100), instructionSet: is, pc: 0, lPr: 0}
}
