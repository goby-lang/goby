package vm

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
