package vm

import (
	"fmt"
	"sync"
)

type objectMap struct {
	store   map[int]Object
	counter int
	sync.RWMutex
}

// ChannelObject represents a goby channel, which carries a golang channel
type ChannelObject struct {
	*baseObj
	Chan chan int
}

func (vm *VM) initChannelClass() *RClass {
	class := vm.initializeClass(channelClass, false)
	class.setBuiltInMethods(builtinChannelClassMethods(), true)
	class.setBuiltInMethods(builtinChannelInstanceMethods(), false)
	return class
}

func builtinChannelClassMethods() []*BuiltInMethodObject {
	return []*BuiltInMethodObject{
		{
			Name: "new",
			Fn: func(receiver Object) builtinMethodBody {
				return func(t *thread, args []Object, blockFrame *callFrame) Object {
					c := &ChannelObject{baseObj: &baseObj{class: t.vm.topLevelClass(channelClass)}, Chan: make(chan int)}
					return c
				}
			},
		},
	}
}

func builtinChannelInstanceMethods() []*BuiltInMethodObject {
	return []*BuiltInMethodObject{
		{
			Name: "close",
			Fn: func(receiver Object) builtinMethodBody {
				return func(t *thread, args []Object, blockFrame *callFrame) Object {
					c := receiver.(*ChannelObject)

					close(c.Chan)

					return NULL
				}
			},
		},
		{
			Name: "deliver",
			Fn: func(receiver Object) builtinMethodBody {
				return func(t *thread, args []Object, blockFrame *callFrame) Object {
					id := t.vm.channelObjectMap.storeObj(args[0])

					c := receiver.(*ChannelObject)

					c.Chan <- id

					return args[0]
				}
			},
		},
		{
			Name: "receive",
			Fn: func(receiver Object) builtinMethodBody {
				return func(t *thread, args []Object, blockFrame *callFrame) Object {
					c := receiver.(*ChannelObject)

					num := <-c.Chan

					return t.vm.channelObjectMap.retrieveObj(num)
				}
			},
		},
	}
}

// Polymorphic helper functions -----------------------------------------

// toString returns detailed info of a channel include elements it contains.
func (co *ChannelObject) toString() string {
	return fmt.Sprintf("<Channel: %p>", co.Chan)
}

// toJSON converts the receiver into JSON string.
func (co *ChannelObject) toJSON() string {
	return co.toString()
}

func (co *ChannelObject) Value() interface{} {
	return co.Chan
}

func (m *objectMap) storeObj(obj Object) int {
	m.Lock()
	defer m.Unlock()

	m.store[m.counter] = obj
	i := m.counter

	// containerCount here can be considered as deliveries' id
	// And this id will be unused once the delivery is completed (which will be really quick)
	// So we can assume that if we reach 1000th delivery,
	// previous deliveries are all completed and they don't need their id anymore.
	if m.counter > 1000 {
		m.counter = 0
	} else {
		m.counter++
	}

	return i
}

// storeObj store objects into the container map
// and update containerCount at the same time
func (m *objectMap) retrieveObj(num int) Object {
	m.RLock()

	defer m.RUnlock()
	return m.store[num]
}
