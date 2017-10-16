package vm

import "sync"

type callFrameStack struct {
	callFrames []callFrame
	thread     *thread
}

type baseFrame struct {
	// environment pointer, points to the call frame we want to get locals from
	ep     *normalCallFrame
	self   Object
	locals []*Pointer
	// local pointer
	lPr        int
	isBlock    bool
	blockFrame *normalCallFrame
	sync.RWMutex
}

type callFrame interface {
	getLCL(index, depth int) *Pointer
	insertLCL(index, depth int, value Object)
	storeConstant(constName string, constant interface{}) *Pointer
	lookupConstant(constName string) *Pointer
	inspect() string
}

type goMethodCallFrame struct {
	*baseFrame
	method *BuiltinMethodObject
}

type normalCallFrame struct {
	*baseFrame
	instructionSet *instructionSet
	// program counter
	pc int
}

// We use lock on every local variable retrieval and insertion.
// The main scenario is when multiple threads want to access local variables outside it's block
// Since they share same block frame, they will all access to that frame's locals.
// TODO: Find a better way to fix this, or prevent thread from accessing outside locals.
func (b *baseFrame) getLCL(index, depth int) *Pointer {
	if depth == 0 {
		b.RLock()

		defer b.RUnlock()

		return b.locals[index]
	}

	return b.blockFrame.ep.getLCL(index, depth-1)
}

func (b *baseFrame) insertLCL(index, depth int, value Object) {
	existedLCL := b.getLCL(index, depth)

	if existedLCL != nil {
		existedLCL.Target = value
		return
	}

	b.Lock()

	defer b.Unlock()

	b.locals = append(b.locals, nil)
	copy(b.locals[index:], b.locals[index:])
	b.locals[index] = &Pointer{Target: value}

	if index >= b.lPr {
		b.lPr = index + 1
	}
}

func (b *baseFrame) storeConstant(constName string, constant interface{}) *Pointer {
	var ptr *Pointer

	switch c := constant.(type) {
	case *Pointer:
		ptr = c
	case Object:
		ptr = &Pointer{Target: c}
	}

	switch scope := b.self.(type) {
	case *RClass:
		scope.constants[constName] = ptr

		if class, ok := ptr.Target.(*RClass); ok {
			class.scope = scope
		}
	default:
		c := b.self.Class()
		c.constants[constName] = ptr
	}

	return ptr
}

func (b *baseFrame) lookupConstant(constName string) *Pointer {
	var c *Pointer

	switch scope := b.self.(type) {
	case *RClass:
		c = scope.lookupConstant(constName, true)
	default:
		scopeClass := scope.Class()
		c = scopeClass.lookupConstant(constName, true)
	}

	return c
}

func (cfs *callFrameStack) push(cf callFrame) {
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

func (cfs *callFrameStack) pop() callFrame {
	var cf callFrame

	if len(cfs.callFrames) < 1 {
		panic("Nothing to pop!")
	}

	if cfs.thread.cfp > 0 {
		cfs.thread.cfp--
	}

	cf = cfs.callFrames[cfs.thread.cfp]
	cfs.callFrames[cfs.thread.cfp] = nil

	switch cf := cf.(type) {
	case *normalCallFrame:
		return cf
	default:
		return nil
	}
}

func (cfs *callFrameStack) top() callFrame {
	var topFrame callFrame

	if cfs.thread.cfp > 0 {
		topFrame = cfs.callFrames[cfs.thread.cfp-1]
	}

	switch f := topFrame.(type) {
	case *normalCallFrame:
		return f
	}

	return nil
}

func newNormalCallFrame(is *instructionSet) *normalCallFrame {
	return &normalCallFrame{baseFrame: &baseFrame{locals: make([]*Pointer, 100), lPr: 0}, instructionSet: is, pc: 0}
}

func newGoMethodCallFrame(m *BuiltinMethodObject) *goMethodCallFrame {
	return &goMethodCallFrame{baseFrame: &baseFrame{locals: make([]*Pointer, 100), lPr: 0}, method: m}
}
