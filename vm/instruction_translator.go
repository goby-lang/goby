package vm

import (
	"github.com/goby-lang/goby/compiler/bytecode"
)

// instructionTranslator is responsible for parsing bytecodes
type instructionTranslator struct {
	vm         *VM
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
	t := set.InstType
	n := set.Name

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

func (it *instructionTranslator) transferInstructionSets(sets []*bytecode.InstructionSet) {
	for _, set := range sets {
		is := &instructionSet{filename: it.filename}
		is.instructions = set.Instructions
		is.ArgSet = set.ArgSet
		it.setMetadata(is, set)
	}
}
