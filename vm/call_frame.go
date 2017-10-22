package vm

import (
	"sync"
)

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
	sourceLine int
	fileName   string
}

type callFrame interface {
	// Getters
	Self() Object
	BlockFrame() *normalCallFrame
	IsBlock() bool
	EP() *normalCallFrame
	Locals() []*Pointer
	LocalPtr() int
	SourceLine() int
	FileName() string

	getLCL(index, depth int) *Pointer
	insertLCL(index, depth int, value Object)
	storeConstant(constName string, constant interface{}) *Pointer
	lookupConstant(constName string) *Pointer
	inspect() string
	stopExecution()
}

type goMethodCallFrame struct {
	*baseFrame
	method builtinMethodBody
	name   string
}

func (cf *goMethodCallFrame) stopExecution() {}

type normalCallFrame struct {
	*baseFrame
	instructionSet *instructionSet
	// program counter
	pc int
}

func (n *normalCallFrame) instructionsCount() int {
	return len(n.instructionSet.instructions)
}

func (n *normalCallFrame) stopExecution() {
	n.pc = n.instructionsCount()
}

func (b *baseFrame) Self() Object {
	return b.self
}

func (b *baseFrame) BlockFrame() *normalCallFrame {
	return b.blockFrame
}

func (b *baseFrame) IsBlock() bool {
	return b.isBlock
}

func (b *baseFrame) EP() *normalCallFrame {
	return b.ep
}

func (b *baseFrame) Locals() []*Pointer {
	return b.locals
}

func (b *baseFrame) LocalPtr() int {
	return b.lPr
}

func (b *baseFrame) SourceLine() int {
	return b.sourceLine
}

func (b *baseFrame) FileName() string {
	return b.fileName
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

	return cf
}

func (cfs *callFrameStack) top() callFrame {
	if cfs.thread.cfp > 0 {
		return cfs.callFrames[cfs.thread.cfp-1]
	}

	return nil
}

func newNormalCallFrame(is *instructionSet, filename string, sourceLine int) *normalCallFrame {
	return &normalCallFrame{baseFrame: &baseFrame{locals: make([]*Pointer, 15), lPr: 0, fileName: filename, sourceLine: sourceLine}, instructionSet: is, pc: 0}
}

func newGoMethodCallFrame(m builtinMethodBody, n, filename string, sourceLine int) *goMethodCallFrame {
	return &goMethodCallFrame{baseFrame: &baseFrame{locals: make([]*Pointer, 15), lPr: 0, fileName: filename, sourceLine: sourceLine}, method: m, name: n}
}
