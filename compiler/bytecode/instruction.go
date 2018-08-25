package bytecode

import (
	"fmt"
	"strings"

	"github.com/goby-lang/goby/compiler/parser/arguments"
)

// instruction set types
const (
	MethodDef = "Def"
	ClassDef  = "DefClass"
	Block     = "Block"
	Program   = "ProgramStart"
)

// instruction actions
const (
	GetLocal uint8 = iota
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

// InstructionNameTable is the table the maps instruction's op code with its readable name
var InstructionNameTable = []string{
	GetLocal:            "getlocal",
	GetConstant:         "getconstant",
	GetInstanceVariable: "getinstancevariable",
	SetLocal:            "setlocal",
	SetConstant:         "setconstant",
	SetInstanceVariable: "setinstancevariable",
	PutBoolean:          "putboolean",
	PutString:           "putstring",
	PutFloat:            "putfloat",
	PutSelf:             "putself",
	PutObject:           "putobject",
	PutNull:             "putnil",
	NewArray:            "newarray",
	ExpandArray:         "expand_array",
	SplatArray:          "splat_array",
	NewHash:             "newhash",
	NewRange:            "newrange",
	BranchUnless:        "branchunless",
	BranchIf:            "branchif",
	Jump:                "jump",
	Break:               "break",
	DefMethod:           "def_method",
	DefSingletonMethod:  "def_singleton_method",
	DefClass:            "def_class",
	Send:                "send",
	InvokeBlock:         "invokeblock",
	GetBlock:            "getblock",
	Pop:                 "pop",
	Dup:                 "dup",
	Leave:               "leave",
}

// Instruction represents compiled bytecode instruction
type Instruction struct {
	Opcode     uint8
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
	return fmt.Sprintf("%s: %s. source line: %d", i.ActionName(), strings.Join(params, ", "), i.sourceLine)
}

// ActionName returns the human readable name of the instruction
func (i *Instruction) ActionName() string {
	return InstructionNameTable[i.Opcode]
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

type ArgSet arguments.ArgSet

// InstructionSet contains a set of Instructions and some metadata
type InstructionSet struct {
	name         string
	isType       string
	Instructions []*Instruction
	count        int
	argTypes     *ArgSet
}

func (as *ArgSet) FindIndex(name string) int {
	for i, n := range as.Names {
		if n == name {
			return i
		}
	}

	return -1
}

func (as *ArgSet) setArg(index int, name string, argType uint8) {
	as.Names[index] = name
	as.Types[index] = argType
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

func (is *InstructionSet) define(action uint8, sourceLine int, params ...interface{}) *Instruction {
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
