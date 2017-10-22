package vm

import "github.com/goby-lang/goby/compiler/bytecode"

// InitForREPL does following things:
// - Initialize instruction sets' index tables
// - Set vm to REPL mode
// - Create and push main object frame
func (vm *VM) InitForREPL() {
	vm.SetClassISIndexTable("REPL")
	vm.SetMethodISIndexTable("REPL")

	// REPL should maintain a base call frame so that the whole program won't exit
	cf := newNormalCallFrame(&instructionSet{name: "REPL base"}, "REPL", 1)
	cf.self = vm.mainObj
	vm.mode = REPLMode
	vm.mainThread.callFrameStack.push(cf)
}

// REPLExec executes instructions differently from normal program execution.
func (vm *VM) REPLExec(sets []*bytecode.InstructionSet) {
	p := newInstructionTranslator("REPL")
	p.vm = vm
	p.transferInstructionSets(sets)

	for setType, table := range p.setTable {
		for name, is := range table {
			vm.isTables[setType][name] = is
		}
	}

	vm.blockTables[p.filename] = p.blockTable

	oldFrame := vm.mainThread.callFrameStack.pop()
	cf := newNormalCallFrame(p.program, p.filename, oldFrame.SourceLine())
	cf.self = vm.mainObj
	cf.locals = oldFrame.Locals()
	cf.ep = oldFrame.EP()
	cf.isBlock = oldFrame.IsBlock()
	cf.self = oldFrame.Self()
	cf.lPr = oldFrame.LocalPtr()
	vm.mainThread.callFrameStack.push(cf)
	vm.startFromTopFrame()
}

// GetExecResult returns stack's top most value. Normally it's used in tests.
func (vm *VM) GetExecResult() Object {
	top := vm.mainThread.stack.top()
	if top != nil {
		return top.Target
	}
	return NULL
}

// GetREPLResult returns strings that should be showed after each evaluation.
func (vm *VM) GetREPLResult() string {
	top := vm.mainThread.stack.pop()

	if top != nil {
		return top.Target.toString()
	}

	return ""
}
