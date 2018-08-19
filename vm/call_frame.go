package vm

import (
	"sync"
)

type callFrameStack struct {
	callFrames []callFrame
	pointer    int
}

type baseFrame struct {
	// environment pointer, points to the call frame we want to get locals from
	ep     *normalCallFrame
	self   Object
	locals []*Pointer
	// local pointer
	lPr           uint8
	isBlock       bool
	isSourceBlock bool
	// for helping stop the frame execution
	isRemoved  bool
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
	IsSourceBlock() bool
	IsRemoved() bool
	setAsRemoved()
	EP() *normalCallFrame
	Locals() []*Pointer
	LocalPtr() uint8
	SourceLine() int
	FileName() string

	getLCL(index, depth uint8) *Pointer
	insertLCL(index, depth uint8, value Object)
	storeConstant(constName string, constant interface{}) *Pointer
	lookupConstantUnderAllScope(constName string) *Pointer
	lookupConstantUnderCurrentScope(constName string) *Pointer
	lookupConstantInCurrentScope(constName string) *Pointer
	inspect() string
	stopExecution()
}

type goMethodCallFrame struct {
	*baseFrame
	method   builtinMethodBody
	argPtr   int
	argCount int
	receiver Object
	name     string
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

func (b *baseFrame) IsRemoved() bool {
	return b.isRemoved
}

func (b *baseFrame) setAsRemoved() {
	b.isRemoved = true
}

func (b *baseFrame) IsSourceBlock() bool {
	return b.isSourceBlock
}

func (b *baseFrame) EP() *normalCallFrame {
	return b.ep
}

func (b *baseFrame) Locals() []*Pointer {
	return b.locals
}

func (b *baseFrame) LocalPtr() uint8 {
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
func (b *baseFrame) getLCL(index, depth uint8) (p *Pointer) {

	if depth == 0 {
		b.RLock()

		if int(index) >= len(b.locals) {
			b.RUnlock()
			return
		}
		p = b.locals[index]

		b.RUnlock()
		return
	}

	return b.blockFrame.ep.getLCL(index, depth-1)
}

func (b *baseFrame) insertLCL(index, depth uint8, value Object) {
	existedLCL := b.getLCL(index, depth)

	if existedLCL != nil {
		existedLCL.Target = value
		return
	}

	b.Lock()

	if int(index) >= len(b.locals) {
		b.locals = append(b.locals, nil)
		copy(b.locals[index:], b.locals[index:])
	}

	b.locals[index] = &Pointer{Target: value}

	if index >= b.lPr {
		b.lPr = index + 1
	}
	b.Unlock()
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

func (b *baseFrame) lookupConstantUnderAllScope(constName string) *Pointer {
	var c *Pointer

	switch scope := b.self.(type) {
	case *RClass:
		c = scope.lookupConstantUnderAllScope(constName)
	default:
		scopeClass := scope.Class()
		c = scopeClass.lookupConstantUnderAllScope(constName)
	}

	return c
}

func (b *baseFrame) lookupConstantUnderCurrentScope(constName string) *Pointer {
	var c *Pointer

	switch scope := b.self.(type) {
	case *RClass:
		c = scope.lookupConstantUnderCurrentScope(constName)
	default:
		scopeClass := scope.Class()
		c = scopeClass.lookupConstantUnderCurrentScope(constName)
	}

	return c
}

func (b *baseFrame) lookupConstantInCurrentScope(constName string) *Pointer {
	var c *Pointer

	switch scope := b.self.(type) {
	case *RClass:
		c = scope.lookupConstantInCurrentScope(constName)
	default:
		scopeClass := scope.Class()
		c = scopeClass.lookupConstantInCurrentScope(constName)
	}

	return c
}

func (cfs *callFrameStack) push(cf callFrame) {
	if cf == nil {
		panic("Callframe can't be nil!")
	}

	if len(cfs.callFrames) <= cfs.pointer {
		cfs.callFrames = append(cfs.callFrames, cf)
	} else {
		cfs.callFrames[cfs.pointer] = cf
	}

	cfs.pointer++
}

func (cfs *callFrameStack) pop() callFrame {
	var cf callFrame

	if len(cfs.callFrames) < 1 {
		panic("Nothing to pop!")
	}

	if cfs.pointer > 0 {
		cfs.pointer--
	}

	cf = cfs.callFrames[cfs.pointer]
	cfs.callFrames[cfs.pointer] = nil

	return cf
}

func (cfs *callFrameStack) top() callFrame {
	if cfs.pointer > 0 {
		return cfs.callFrames[cfs.pointer-1]
	}

	return nil
}

func newNormalCallFrame(filename string, sourceLine int) *normalCallFrame {
	return &normalCallFrame{baseFrame: &baseFrame{locals: make([]*Pointer, 5), lPr: 0, fileName: filename, sourceLine: sourceLine}, pc: 0}
}

func newGoMethodCallFrame(m builtinMethodBody, receiver Object, argCount, argPtr int, n, filename string, sourceLine int, blockFrame *normalCallFrame) *goMethodCallFrame {
	return &goMethodCallFrame{
		baseFrame: &baseFrame{
			locals:     make([]*Pointer, 5),
			lPr:        0,
			fileName:   filename,
			sourceLine: sourceLine,
			blockFrame: blockFrame,
		},
		method:   m,
		name:     n,
		receiver: receiver,
		argCount: argCount,
		argPtr:   argPtr,
	}
}
