package vm

import (
	"bytes"
	"fmt"
	"strings"
)

type stack struct {
	Data []*Pointer
	VM   *VM
}

func (s *stack) push(v *Pointer) {
	if len(s.Data) <= s.VM.sp {
		s.Data = append(s.Data, v)
	} else {
		s.Data[s.VM.sp] = v
	}

	s.VM.sp++
}

func (s *stack) pop() *Pointer {
	if len(s.Data) < 1 {
		panic("Nothing to pop!")
	}

	s.VM.sp--

	v := s.Data[s.VM.sp]
	s.Data[s.VM.sp] = nil
	return v
}

func (s *stack) top() *Pointer {

	if len(s.Data) == 0 {
		return nil
	}

	if s.VM.sp > 0 {
		return s.Data[s.VM.sp-1]
	}

	return s.Data[0]
}

func (s *stack) inspect() string {
	var out bytes.Buffer
	datas := []string{}

	for i, p := range s.Data {
		if p != nil {
			o := p.Target
			if i == s.VM.sp {
				datas = append(datas, fmt.Sprintf("%s (%T) %d <----", o.Inspect(), o, i))
			} else {
				datas = append(datas, fmt.Sprintf("%s (%T) %d", o.Inspect(), o, i))
			}

		} else {
			if i == s.VM.sp {
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
