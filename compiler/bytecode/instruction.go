package bytecode

//go:generate stringer -type=InstructionAction
// CAUTION: when you change the file, be sure to perform `make generate` for regeneration.

import (
	"fmt"
	"strings"
)

// instruction set types
const (
	MethodDef = "Def"
	ClassDef  = "DefClass"
	Block     = "Block"
	Program   = "ProgramStart"
)

// InstructionAction represents the instruction actions
type InstructionAction uint8

// instruction actions
const (
	GetLocal InstructionAction = iota + 1
	GetConstant
	GetInstanceVariable
	SetLocal
	SetConstant
	SetInstanceVariable
	PutBoolean
	PutString
	PutFloat
	PutSelf
	PutObject
	PutNull
	NewArray
	ExpandArray
	SplatArray
	NewHash
	NewRange
	BranchUnless
	BranchIf
	Jump
	Break
	DefMethod
	DefSingletonMethod
	DefClass
	Send
	InvokeBlock
	GetBlock
	Pop
	Dup
	Leave
)

// Instruction represents compiled bytecode instruction
type Instruction struct {
	Opcode     InstructionAction
	Params     []interface{}
	line       int
	anchor     *anchor
	sourceLine int
}

// Inspect is for inspecting the instruction's content
func (i *Instruction) Inspect() string {
	var params []string

	for _, param := range i.Params {
		params = append(params, fmt.Sprint(param))
	}
	return fmt.Sprintf("%s: %s. source line: %d", i.Opcode.String(), strings.Join(params, ", "), i.sourceLine)
}

// AnchorLine returns instruction anchor's line number if it has an anchor
func (i *Instruction) AnchorLine() int {
	if i.anchor != nil {
		return i.anchor.line
	}

	panic("you are calling AnchorLine on an instruction without anchors")
}

// Line returns instruction's line number
func (i *Instruction) Line() int {
	return i.line
}

// SourceLine returns instruction's source line number
func (i *Instruction) SourceLine() int {
	return i.sourceLine
}

type anchor struct {
	line int
}

// InstructionSet contains a set of Instructions and some metadata
type InstructionSet struct {
	name         string
	isType       string
	Instructions []*Instruction
	count        int
	argTypes     *ArgSet
}

// ArgSet stores the metadata of a method definition's parameters.
type ArgSet struct {
	names []string
	types []uint8
}

// Types are the getter method of *ArgSet's types attribute
func (as *ArgSet) Types() []uint8 {
	return as.types
}

// Names are the getter method of *ArgSet's names attribute
func (as *ArgSet) Names() []string {
	return as.names
}

func (as *ArgSet) FindIndex(name string) int {
	for i, n := range as.names {
		if n == name {
			return i
		}
	}

	return -1
}

func (as *ArgSet) setArg(index int, name string, argType uint8) {
	as.names[index] = name
	as.types[index] = argType
}

// ArgTypes returns enums that represents each argument's type
func (is *InstructionSet) ArgTypes() *ArgSet {
	return is.argTypes
}

// Name returns instruction set's name
func (is *InstructionSet) Name() string {
	return is.name
}

// SetType returns instruction's type
func (is *InstructionSet) Type() string {
	return is.isType
}

func (is *InstructionSet) define(action InstructionAction, sourceLine int, params ...interface{}) *Instruction {
	i := &Instruction{Opcode: action, Params: params, line: is.count, sourceLine: sourceLine + 1}
	for _, param := range params {
		a, ok := param.(*anchor)

		if ok {
			i.anchor = a
		}
	}

	is.Instructions = append(is.Instructions, i)
	is.count++
	return i
}
