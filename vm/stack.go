package vm

import (
	"os"
	"sync"
)

type stack struct {
	Data   []*Pointer
	thread *thread
	// Although every thread has its own stack, vm's main thread still can be accessed by other threads.
	// This is why we need a lock in stack
	// TODO: Find a way to fix this instead of put lock on every stack.
	*sync.RWMutex
}

func (s *stack) set(index int, pointer *Pointer) {
	t := s.thread

	if _, ok := pointer.Target.(*Error); ok {
		cf := t.callFrameStack.top()
		cf.pc = len(cf.instructionSet.instructions)

		if t.vm.mode == NormalMode {
			if t.isMainThread() {
				os.Exit(1)
			}
		}
	}

	s.Lock()

	defer s.Unlock()

	s.Data[index] = pointer
}

func (s *stack) push(v *Pointer) {
	s.Lock()
	defer s.Unlock()

	if len(s.Data) <= s.thread.sp {
		s.Data = append(s.Data, v)
	} else {
		s.Data[s.thread.sp] = v
	}

	if _, ok := v.Target.(*Error); ok {
		t := s.thread
		cf := t.callFrameStack.top()
		cf.pc = len(cf.instructionSet.instructions)

		if t.vm.mode == NormalMode {
			if t.isMainThread() {
				os.Exit(1)
			}
		}
	}

	s.thread.sp++
}

func (s *stack) pop() *Pointer {
	s.Lock()
	defer s.Unlock()

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
	s.RLock()
	defer s.RUnlock()

	if len(s.Data) == 0 {
		return nil
	}

	if s.thread.sp > 0 {
		return s.Data[s.thread.sp-1]
	}

	return s.Data[0]
}
