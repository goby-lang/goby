package vm

import (
	"fmt"
	"github.com/goby-lang/goby/compiler/bytecode"
)

// instructionTranslator is responsible for parsing bytecodes
type instructionTranslator struct {
	vm         *VM
	line       int
	setTable   map[setType]map[string][]*instructionSet
	blockTable map[string]*instructionSet
	filename   filename
	program    *instructionSet
}

// newInstructionTranslator initializes instructionTranslator and its instruction set table then returns it
func newInstructionTranslator(file filename) *instructionTranslator {
	it := &instructionTranslator{filename: file}
	it.blockTable = make(map[string]*instructionSet)
	it.setTable = map[setType]map[string][]*instructionSet{
		bytecode.MethodDef: make(map[string][]*instructionSet),
		bytecode.ClassDef:  make(map[string][]*instructionSet),
	}

	return it
}

func (it *instructionTranslator) setMetadata(is *instructionSet, set *bytecode.InstructionSet) {
	t := set.Type()
	n := set.Name()

	is.name = n

	switch t {
	case bytecode.Program:
		it.program = is
	case bytecode.Block:
		it.blockTable[n] = is
	default:
		it.setTable[t][n] = append(it.setTable[t][n], is)
	}
}

func (it *instructionTranslator) transferInstructionSets(sets []*bytecode.InstructionSet) []*instructionSet {
	iss := []*instructionSet{}

	for _, set := range sets {
		it.transferInstructionSet(iss, set)
	}

	return iss
}

func (it *instructionTranslator) transferInstructionSet(iss []*instructionSet, set *bytecode.InstructionSet) {
	is := &instructionSet{filename: it.filename}
	it.setMetadata(is, set)

	for _, i := range set.Instructions {
		it.transferInstruction(is, i)
	}

	is.paramTypes = set.ArgTypes()

	iss = append(iss, is)
}

// transferInstruction transfer a bytecode.Instruction into an vm instruction and append it into given instruction set.
func (it *instructionTranslator) transferInstruction(is *instructionSet, i *bytecode.Instruction) {
	act := i.Action

	action := builtinActions[act]

	if action == nil {
		panic(fmt.Sprintf("Unknown command: %d. line: %d", act, i.Line()))
	}

	vmI := is.define(i.Line(), action, i.Params...)
	vmI.sourceLine = i.SourceLine() + 1
}
