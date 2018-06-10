package vm

import (
	"bytes"
	"fmt"
	"strings"
)

// Polymorphic helper functions for inspecting internal info.

func (i *instruction) inspect() string {
	var params []string

	for _, param := range i.Params {
		params = append(params, fmt.Sprint(param))
	}
	return fmt.Sprintf("%s: %s. source line: %d", i.action.name, strings.Join(params, ", "), i.sourceLine)
}

func (is *instructionSet) inspect() string {
	var out bytes.Buffer

	for _, i := range is.instructions {
		out.WriteString(i.inspect())
		out.WriteString("\n")
	}

	return out.String()
}

func (cf *goMethodCallFrame) inspect() string {
	return fmt.Sprintf("Go method frame. File name: %s. Method name: %s.", cf.FileName(), cf.name)
}

func (cf *normalCallFrame) inspect() string {
	if cf.ep != nil {
		return fmt.Sprintf("Normal frame. File name: %s. IS name: %s. is block: %t. ep: %d. source line: %d", cf.FileName(), cf.instructionSet.name, cf.isBlock, len(cf.ep.locals), cf.SourceLine())
	}
	return fmt.Sprintf("Normal frame. File name: %s. IS name: %s. is block: %t. source line: %d", cf.FileName(), cf.instructionSet.name, cf.isBlock, cf.SourceLine())
}

func (cfs *callFrameStack) inspect() string {
	var out bytes.Buffer

	for _, cf := range cfs.callFrames {
		if cf != nil {
			out.WriteString(fmt.Sprintln(cf.inspect()))
		}
	}

	return out.String()
}

func (s *Stack) inspect() string {
	var out bytes.Buffer
	datas := []string{}

	for i, p := range s.data {
		if p != nil {
			o := p.Target
			if i == s.pointer {
				datas = append(datas, fmt.Sprintf("%s (%T) %d <----", o.toString(), o, i))
			} else {
				datas = append(datas, fmt.Sprintf("%s (%T) %d", o.toString(), o, i))
			}

		} else {
			if i == s.pointer {
				datas = append(datas, "nil <----")
			} else {
				datas = append(datas, "nil")
			}

		}

	}

	out.WriteString("-----------\n")
	out.WriteString(strings.Join(datas, "\n"))
	out.WriteString("\n---------\n")

	return out.String()
}
