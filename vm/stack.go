package vm

import (
	"sync"
)

type stack struct {
	Data    []*Pointer
	pointer int
	// Although every thread has its own stack, vm's main thread still can be accessed by other threads.
	// This is why we need a lock in stack
	// TODO: Find a way to fix this instead of put lock on every stack.
	sync.RWMutex
}

func (s *stack) set(index int, pointer *Pointer) {
	s.Lock()

	s.Data[index] = pointer

	s.Unlock()
}

func (s *stack) push(v *Pointer) {
	s.Lock()

	if len(s.Data) <= s.pointer {
		s.Data = append(s.Data, v)
	} else {
		s.Data[s.pointer] = v
	}

	s.pointer++
	s.Unlock()
}

func (s *stack) pop() *Pointer {
	s.Lock()

	if len(s.Data) < 1 {
		panic("Nothing to pop!")
	}

	if s.pointer < 0 {
		panic("SP is not normal!")
	}

	if s.pointer > 0 {
		s.pointer--
	}

	v := s.Data[s.pointer]
	s.Data[s.pointer] = nil
	s.Unlock()
	return v
}

func (s *stack) top() *Pointer {
	var r *Pointer
	s.RLock()

	if len(s.Data) == 0 {
		r = nil
	} else if s.pointer > 0 {
		r = s.Data[s.pointer-1]
	} else {
		r = s.Data[0]
	}

	s.RUnlock()

	return r
}
