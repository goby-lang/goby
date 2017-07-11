package vm

import (
	"fmt"
	"github.com/goby-lang/goby/compiler/bytecode"
	"strconv"
	"strings"
)

// instructionTranslator is responsible for parsing bytecodes
type instructionTranslator struct {
	vm         *VM
	line       int
	labelTable map[labelType]map[string][]*instructionSet
	blockTable map[string]*instructionSet
	filename   filename
	program    *instructionSet
}

// newInstructionTranslator initializes instructionTranslator and its label table then returns it
func newInstructionTranslator(file filename) *instructionTranslator {
	it := &instructionTranslator{filename: file}
	it.blockTable = make(map[string]*instructionSet)
	it.labelTable = map[labelType]map[string][]*instructionSet{
		bytecode.LabelDef:      make(map[string][]*instructionSet),
		bytecode.LabelDefClass: make(map[string][]*instructionSet),
	}

	return it
}

func (it *instructionTranslator) setLabel(is *instructionSet, set *bytecode.InstructionSet) {
	t := labelType(set.SetType())
	n := set.Name()

	is.name = n

	if t == bytecode.Program {
		it.program = is
		return
	}

	if t == bytecode.Block {
		it.blockTable[n] = is
		return
	}

	it.labelTable[t][n] = append(it.labelTable[t][n], is)
}

func (it *instructionTranslator) parseParam(param string) interface{} {
	integer, e := strconv.ParseInt(param, 0, 64)
	if e != nil {
		return param
	}

	i := int(integer)

	return i
}

func (it *instructionTranslator) transferInstructionSets(sets []*bytecode.InstructionSet) []*instructionSet {
	iss := []*instructionSet{}
	count := 0

	for _, set := range sets {
		count++
		it.transferInstructionSet(iss, set)
	}

	return iss
}

func (it *instructionTranslator) transferInstructionSet(iss []*instructionSet, set *bytecode.InstructionSet) {
	is := &instructionSet{filename: it.filename}
	count := 0
	it.setLabel(is, set)

	for _, i := range set.Instructions {
		count++
		it.transferInstruction(is, i)
	}

	is.argTypes = set.ArgTypes()

	iss = append(iss, is)
}

// transferInstruction transfer a bytecode.Instruction into an vm instruction and append it into given instruction set.
func (it *instructionTranslator) transferInstruction(is *instructionSet, i *bytecode.Instruction) {
	var params []interface{}
	act := operationType(i.Action)

	action := builtInActions[act]

	if action == nil {
		panic(fmt.Sprintf("Unknown command: %s. line: %d", act, i.Line()))
	}

	switch act {
	case bytecode.PutString:
		text := strings.Split(i.Params[0], "\"")[1]
		params = append(params, text)
	case bytecode.BranchUnless, bytecode.BranchIf, bytecode.Jump:
		line, err := i.AnchorLine()

		if err != nil {
			panic(err.Error())
		}

		params = append(params, line)
	default:
		for _, param := range i.Params {
			params = append(params, it.parseParam(param))
		}
	}

	is.define(i.Line(), action, params...)
}
