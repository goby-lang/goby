package bytecode

import (
	"bytes"
	"fmt"
	"strings"
)

const (
	// label types
	LabelDef      = "Def"
	LabelDefClass = "DefClass"
	Block         = "Block"
	Program       = "ProgramStart"

	// instruction actions
	GetLocal            = "getlocal"
	GetConstant         = "getconstant"
	GetInstanceVariable = "getinstancevariable"
	SetLocal            = "setlocal"
	SetConstant         = "setconstant"
	SetInstanceVariable = "setinstancevariable"
	PutString           = "putstring"
	PutRegexp           = "putregexp"
	PutSelf             = "putself"
	PutObject           = "putobject"
	PutNull             = "putnil"
	NewArray            = "newarray"
	NewHash             = "newhash"
	NewRange            = "newrange"
	NewRegexp           = "newregexp"
	BranchUnless        = "branchunless"
	BranchIf            = "branchif"
	Jump                = "jump"
	DefMethod           = "def_method"
	DefSingletonMethod  = "def_singleton_method"
	DefClass            = "def_class"
	Send                = "send"
	InvokeBlock         = "invokeblock"
	Pop                 = "pop"
	Leave               = "leave"
)

// Instruction represents compiled bytecode instruction
type Instruction struct {
	Action string
	Params []string
	line   int
	anchor *anchor
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

func (i *Instruction) compile() string {
	if i.anchor != nil {
		return fmt.Sprintf("%d %s %d\n", i.line, i.Action, i.anchor.line)
	}
	if len(i.Params) > 0 {
		return fmt.Sprintf("%d %s %s\n", i.line, i.Action, strings.Join(i.Params, " "))
	}

	return fmt.Sprintf("%d %s\n", i.line, i.Action)
}

type label struct {
	name string
}

type anchor struct {
	line int
}

func (l *label) compile() string {
	return fmt.Sprintf("<%s>\n", l.name)
}

// InstructionSet contains a set of Instructions and attaches a label
type InstructionSet struct {
	label        *label
	Instructions []*Instruction
	count        int
}

// LabelName returns the label name of instruction set
func (is *InstructionSet) LabelName() string {
	return is.label.name
}

func (is *InstructionSet) setLabel(name string) {
	l := &label{name: name}
	is.label = l
}

func (is *InstructionSet) define(action string, params ...interface{}) {
	ps := []string{}
	i := &Instruction{Action: action, Params: ps, line: is.count}
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
	out.WriteString(is.label.compile())
	for _, i := range is.Instructions {
		out.WriteString(i.compile())
	}

	return out.String()
}
