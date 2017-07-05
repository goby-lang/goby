package vm

import (
	"fmt"
	"github.com/goby-lang/goby/compiler/bytecode"
	"strconv"
	"strings"
)

// bytecodeParser is responsible for parsing bytecodes
type bytecodeParser struct {
	line       int
	labelTable map[labelType]map[string][]*instructionSet
	vm         *VM
	filename   filename
	blockTable map[string]*instructionSet
	program    *instructionSet
}

// newBytecodeParser initializes bytecodeParser and its label table then returns it
func newBytecodeParser(file filename) *bytecodeParser {
	p := &bytecodeParser{filename: file}
	p.blockTable = make(map[string]*instructionSet)
	p.labelTable = map[labelType]map[string][]*instructionSet{
		bytecode.LabelDef:      make(map[string][]*instructionSet),
		bytecode.LabelDefClass: make(map[string][]*instructionSet),
	}

	return p
}

func (p *bytecodeParser) setLabel(is *instructionSet, name string) {
	var l *label
	var ln string
	var lt labelType

	if name == bytecode.Program {
		p.program = is
		return
	}

	ln = strings.Split(name, ":")[1]
	lt = labelType(strings.Split(name, ":")[0])

	l = &label{name: name, Type: lt}
	is.label = l

	if lt == bytecode.Block {
		p.blockTable[ln] = is
		return
	}

	p.labelTable[lt][ln] = append(p.labelTable[lt][ln], is)
}

func (p *bytecodeParser) parseParam(param string) interface{} {
	integer, e := strconv.ParseInt(param, 0, 64)
	if e != nil {
		return param
	}

	i := int(integer)

	return i
}

func (p *bytecodeParser) transferInstructionSets(sets []*bytecode.InstructionSet) []*instructionSet {
	iss := []*instructionSet{}
	count := 0

	for _, set := range sets {
		count++
		p.transferInstructionSet(iss, set)
	}

	return iss
}

func (p *bytecodeParser) transferInstructionSet(iss []*instructionSet, set *bytecode.InstructionSet) {
	is := &instructionSet{filename: p.filename}
	count := 0
	p.setLabel(is, set.LabelName())

	for _, i := range set.Instructions {
		count++
		p.transferInstruction(is, i)
	}

	iss = append(iss, is)
}

// transferInstruction transfer a bytecode.Instruction into an vm instruction and append it into given instruction set.
func (p *bytecodeParser) transferInstruction(is *instructionSet, i *bytecode.Instruction) {
	var params []interface{}
	act := operationType(i.Action)

	action := builtInActions[act]

	if action == nil {
		panic(fmt.Sprintf("Unknown command: %s. line: %d", act, i.Line()))
	} else {
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
				params = append(params, p.parseParam(param))
			}
		}
	}

	is.define(i.Line(), action, params...)
}