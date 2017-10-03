package bytecode

import (
	"bytes"
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

// instruction actions
const (
	GetLocal            = "getlocal"
	GetConstant         = "getconstant"
	GetInstanceVariable = "getinstancevariable"
	SetLocal            = "setlocal"
	SetConstant         = "setconstant"
	SetInstanceVariable = "setinstancevariable"
	PutString           = "putstring"
	PutSelf             = "putself"
	PutObject           = "putobject"
	PutNull             = "putnil"
	NewArray            = "newarray"
	ExpandArray         = "expand_array"
	SplatArray          = "splat_array"
	NewHash             = "newhash"
	NewRange            = "newrange"
	BranchUnless        = "branchunless"
	BranchIf            = "branchif"
	Jump                = "jump"
	DefMethod           = "def_method"
	DefSingletonMethod  = "def_singleton_method"
	DefClass            = "def_class"
	Send                = "send"
	InvokeBlock         = "invokeblock"
	Pop                 = "pop"
	Dup                 = "dup"
	Leave               = "leave"
)

// Instruction represents compiled bytecode instruction
type Instruction struct {
	Action     string
	Params     []string
	line       int
	anchor     *anchor
	sourceLine int
}

// AnchorLine returns instruction anchor's line number if it has an anchor
func (i *Instruction) AnchorLine() (int, error) {
	if i.anchor != nil {
		return i.anchor.line, nil
	}

	return 0, fmt.Errorf("Can't find anchor on action %s", i.Action)
}

// Line returns instruction's line number
func (i *Instruction) Line() int {
	return i.line
}

// SourceLine returns instruction's source line number
func (i *Instruction) SourceLine() int {
	return i.sourceLine
}

func (i *Instruction) compile() string {
	if i.anchor != nil {
		return fmt.Sprintf("%d %s %d\n", i.line, i.Action, i.anchor.line)
	}
	if len(i.Params) > 0 {
		return fmt.Sprintf("%d %s %s\n", i.line, i.Action, strings.Join(i.Params, " "))
	}

	return fmt.Sprintf("%d %s\n", i.line, i.Action)
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

type ArgSet struct {
	names []string
	types []int
	count int
}

func (as *ArgSet) Types() []int {
	return as.types
}

func (as *ArgSet) Names() []string {
	return as.names
}

func (as *ArgSet) Count() int {
	return as.count
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

func (is *InstructionSet) define(action string, sourceLine int, params ...interface{}) {
	ps := []string{}
	i := &Instruction{Action: action, Params: ps, line: is.count, sourceLine: sourceLine}
	for _, param := range params {
		switch p := param.(type) {
		case string:
			ps = append(ps, p)
		case *anchor:
			i.anchor = p
		case int:
			ps = append(ps, fmt.Sprint(p))
		}
	}

	if len(ps) > 0 {
		i.Params = ps
	}

	is.Instructions = append(is.Instructions, i)
	is.count++
}

func (is *InstructionSet) compile() string {
	var out bytes.Buffer
	if is.isType == Program {
		out.WriteString(fmt.Sprintf("<%s>\n", is.isType))
	} else {
		out.WriteString(fmt.Sprintf("<%s:%s>\n", is.isType, is.name))
	}

	for _, i := range is.Instructions {
		out.WriteString(i.compile())
	}

	return out.String()
}
