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
}

func (i *Instruction) Compile() string {
	if len(i.Params) > 0 {
		return fmt.Sprintf("%d %s %s\n", i.Line, i.Action, strings.Join(i.Params, " "))
	}

	return fmt.Sprintf("%d %s\n", i.Line, i.Action)
}

type Label struct {
	Name string
}

func (l *Label) Compile() string {
	return fmt.Sprintf("<%s>\n", l.Name)
}

type InstructionSet struct {
	Label        *Label
	Instructions []*Instruction
	Count        int
}

func (is *InstructionSet) Define(action string, params ...string) {
	i := &Instruction{Action: action, Params: params, Line: is.Count}
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
