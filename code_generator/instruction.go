package code_generator

import (
	"bytes"
	"fmt"
	"strings"
)

type Instruction struct {
	Action string
	Params []string
	Line   int
	Anchor *Anchor
}

func (i *Instruction) Compile() string {
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

func (l *Label) Compile() string {
	return fmt.Sprintf("<%s>\n", l.Name)
}

type InstructionSet struct {
	Label        *Label
	Instructions []*Instruction
	Count        int
}

func (is *InstructionSet) Define(action string, params ...interface{}) {
	ps := []string{}
	i := &Instruction{Action: action, Params: ps, Line: is.Count}
	for _, param := range params {
		switch p := param.(type) {
		case string:
			ps = append(ps, p)
		case *Anchor:
			i.Anchor = p
		}
	}

	if len(ps) > 0 {
		i.Params = ps
	}

	is.Instructions = append(is.Instructions, i)
	is.Count += 1
}

func (is *InstructionSet) Compile() string {
	var out bytes.Buffer
	out.WriteString(is.Label.Compile())
	for _, i := range is.Instructions {
		out.WriteString(i.Compile())
	}

	return out.String()
}
