package vm

import (
	"sync"
)

// Stack is a basic stack implementation
type Stack struct {
	data    []*Pointer
	pointer int
	// Although every thread has its own stack, vm's main thread still can be accessed by other threads.
	// This is why we need a lock in stack
	// TODO: Find a way to fix this instead of put lock on every stack.
	sync.RWMutex
}

// Set a value at a given index in the stack. TODO: Maybe we should be checking for size before we do this.
func (s *Stack) Set(index int, pointer *Pointer) {
	s.Lock()

	s.data[index] = pointer

	s.Unlock()
}

// Push an element to the top of the stack
func (s *Stack) Push(v *Pointer) {
	s.Lock()

	if len(s.data) <= s.pointer {
		s.data = append(s.data, v)
	} else {
		s.data[s.pointer] = v
	}

	s.pointer++
	s.Unlock()
}

// Pop an element off the top of the stack
func (s *Stack) Pop() *Pointer {
	s.Lock()

	if len(s.data) < 1 {
		panic("Nothing to pop!")
	}

	if s.pointer < 0 {
		panic("SP is not normal!")
	}

	if s.pointer > 0 {
		s.pointer--
	}

	v := s.data[s.pointer]
	s.data[s.pointer] = nil
	s.Unlock()
	return v
}

func (s *Stack) top() *Pointer {
	var r *Pointer
	s.RLock()

	if len(s.data) == 0 {
		r = nil
	} else if s.pointer > 0 {
		r = s.data[s.pointer-1]
	} else {
		r = s.data[0]
	}

	s.RUnlock()

	return r
}
