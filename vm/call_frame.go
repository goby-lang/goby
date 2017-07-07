package vm

import "sync"

type callFrameStack struct {
	callFrames []*callFrame
	thread     *thread
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
	sync.RWMutex
}

// We use lock on every local variable retrieval and insertion.
// The main scenario is when multiple threads want to access local variables outside it's block
// Since they share same block frame, they will all access to that frame's locals.
// TODO: Find a better way to fix this, or prevent thread from accessing outside locals.
func (cf *callFrame) getLCL(index, depth int) *Pointer {
	if depth == 0 {
		cf.RLock()

		defer cf.RUnlock()

		return cf.locals[index]
	}

	return cf.blockFrame.ep.getLCL(index, depth-1)
}

func (cf *callFrame) insertLCL(index, depth int, value Object) {
	existedLCL := cf.getLCL(index, depth)

	if existedLCL != nil {
		existedLCL.Target = value
		return
	}

	cf.Lock()

	defer cf.Unlock()

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
		c := cf.self.returnClass()
		c.constants[constName] = ptr
	}

	return ptr
}

func (cf *callFrame) lookupConstant(constName string) *Pointer {
	var c *Pointer

	switch scope := cf.self.(type) {
	case *RClass:
		c = scope.lookupConstant(constName, true)
	default:
		scopeClass := scope.returnClass()
		c = scopeClass.lookupConstant(constName, true)
	}

	return c
}

func (cfs *callFrameStack) push(cf *callFrame) {
	if cf == nil {
		panic("Callframe can't be nil!")
	}

	if len(cfs.callFrames) <= cfs.thread.cfp {
		cfs.callFrames = append(cfs.callFrames, cf)
	} else {
		cfs.callFrames[cfs.thread.cfp] = cf
	}

	cfs.thread.cfp++
}

func (cfs *callFrameStack) pop() *callFrame {
	if len(cfs.callFrames) < 1 {
		panic("Nothing to pop!")
	}

	if cfs.thread.cfp > 0 {
		cfs.thread.cfp--
	}

	cf := cfs.callFrames[cfs.thread.cfp]
	cfs.callFrames[cfs.thread.cfp] = nil
	return cf
}

func (cfs *callFrameStack) top() *callFrame {
	if cfs.thread.cfp > 0 {
		return cfs.callFrames[cfs.thread.cfp-1]
	}

	return nil
}

func newCallFrame(is *instructionSet) *callFrame {
	return &callFrame{locals: make([]*Pointer, 100), instructionSet: is, pc: 0, lPr: 0}
}
