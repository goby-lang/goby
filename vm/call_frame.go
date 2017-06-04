package vm

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

func (cf *callFrame) storeConstant(constName string, constant interface{}) *Pointer {
	var ptr *Pointer

	switch c := constant.(type) {
	case *RClass:
		ptr = &Pointer{Target: c}
	case *Pointer:
		ptr = c
	}

	switch scope := cf.self.(type) {
	case *RClass:
		scope.constants[constName] = ptr

		if class, ok := ptr.Target.(*RClass); ok {
			class.scope = scope
		}
	default:
		c := cf.self.returnClass().(*RClass)
		c.constants[constName] = ptr
	}

	return ptr
}

func (cf *callFrame) lookupConstant(constName string) (*Pointer, bool) {
	switch scope := cf.self.(type) {
	case *RClass:
		p, ok := scope.constants[constName]

		if ok {
			return p, true
		}
	default:
		c := cf.self.returnClass().(*RClass)
		p, ok := c.constants[constName]

		if ok {
			return p, true
		}
	}

	return nil, false
}

func (cf *callFrame) getLCL(index, depth int) *Pointer {
	if depth == 0 {
		return cf.locals[index]
	}

	return cf.blockFrame.ep.getLCL(index, depth-1)
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

func newCallFrame(is *instructionSet) *callFrame {
	return &callFrame{locals: make([]*Pointer, 100), instructionSet: is, pc: 0, lPr: 0}
}
