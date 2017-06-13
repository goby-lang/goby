package vm

type stack struct {
	Data   []*Pointer
	thread *thread
}

func (s *stack) push(v *Pointer) {
	if len(s.Data) <= s.thread.sp {
		s.Data = append(s.Data, v)
	} else {
		s.Data[s.thread.sp] = v
	}

	s.thread.sp++
}

func (s *stack) pop() *Pointer {
	if len(s.Data) < 1 {
		panic("Nothing to pop!")
	}

	if s.thread.sp < 0 {
		panic("SP is not normal!")
	}

	if s.thread.sp > 0 {
		s.thread.sp--
	}

	v := s.Data[s.thread.sp]
	s.Data[s.thread.sp] = nil
	return v
}

func (s *stack) top() *Pointer {

	if len(s.Data) == 0 {
		return nil
	}

	if s.thread.sp > 0 {
		return s.Data[s.thread.sp-1]
	}

	return s.Data[0]
}
