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
	PutSelf             = "putself"
	PutObject           = "putobject"
	PutNull             = "putnil"
	NewArray            = "newarray"
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
	Leave               = "leave"
)

type Instruction struct {
	Action string
	Params []string
	Line   int
	Anchor *Anchor
}

func (i *Instruction) compile() string {
	if i.Anchor != nil {
		return fmt.Sprintf("%d %s %d\n", i.Line, i.Action, i.Anchor.Line)
	}
	if len(i.Params) > 0 {
		return fmt.Sprintf("%d %s %s\n", i.Line, i.Action, strings.Join(i.Params, " "))
	}

	return fmt.Sprintf("%d %s\n", i.Line, i.Action)
}

type Label struct {
	Name string
}

type Anchor struct {
	Line int
}

func (l *Label) compile() string {
	return fmt.Sprintf("<%s>\n", l.Name)
}

type InstructionSet struct {
	Label        *Label
	Instructions []*Instruction
	Count        int
}

func (is *InstructionSet) setLabel(name string) {
	l := &Label{Name: name}
	is.Label = l
}

func (is *InstructionSet) define(action string, params ...interface{}) {
	ps := []string{}
	i := &Instruction{Action: action, Params: ps, Line: is.Count}
	for _, param := range params {
		switch p := param.(type) {
		case string:
			ps = append(ps, p)
		case *Anchor:
			i.Anchor = p
		case int:
			ps = append(ps, fmt.Sprint(p))
		case int64:
			ps = append(ps, fmt.Sprint(p))
		}
	}

	if len(ps) > 0 {
		i.Params = ps
	}

	is.Instructions = append(is.Instructions, i)
	is.Count++
}

func (is *InstructionSet) compile() string {
	var out bytes.Buffer
	out.WriteString(is.Label.compile())
	for _, i := range is.Instructions {
		out.WriteString(i.compile())
	}

	return out.String()
}
