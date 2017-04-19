package bytecode

import (
	"bytes"
	"fmt"
	"strings"
)

type instruction struct {
	action string
	params []string
	line   int
	anchor *anchor
}

func (i *instruction) compile() string {
	if i.anchor != nil {
		return fmt.Sprintf("%d %s %d\n", i.line, i.action, i.anchor.line)
	}
	if len(i.params) > 0 {
		return fmt.Sprintf("%d %s %s\n", i.line, i.action, strings.Join(i.params, " "))
	}

	return fmt.Sprintf("%d %s\n", i.line, i.action)
}

type label struct {
	Name string
}

type anchor struct {
	line int
}

func (l *label) compile() string {
	return fmt.Sprintf("<%s>\n", l.Name)
}

type instructionSet struct {
	label        *label
	Instructions []*instruction
	Count        int
}

func (is *instructionSet) setLabel(name string) {
	l := &label{Name: name}
	is.label = l
}

func (is *instructionSet) define(action string, params ...interface{}) {
	ps := []string{}
	i := &instruction{action: action, params: ps, line: is.Count}
	for _, param := range params {
		switch p := param.(type) {
		case string:
			ps = append(ps, p)
		case *anchor:
			i.anchor = p
		case int:
			ps = append(ps, fmt.Sprint(p))
		case int64:
			ps = append(ps, fmt.Sprint(p))
		}
	}

	if len(ps) > 0 {
		i.params = ps
	}

	is.Instructions = append(is.Instructions, i)
	is.Count += 1
}

func (is *instructionSet) compile() string {
	var out bytes.Buffer
	out.WriteString(is.label.compile())
	for _, i := range is.Instructions {
		out.WriteString(i.compile())
	}

	return out.String()
}
