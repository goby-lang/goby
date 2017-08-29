package vm

import (
	"fmt"
	"sync"
)

func (vm *VM) initChannelClass() *RClass {
	class := vm.initializeClass(channelClass, false)
	class.setBuiltInMethods(builtinChannelClassMethods(), true)
	class.setBuiltInMethods(builtinChannelInstanceMethods(), false)
	return class
}

type objectMap struct {
	store *sync.Map
}

func (m *objectMap) storeObj(obj Object) int {
	m.store.Store(obj.id(), obj)

	return obj.id()
}

// storeObj store objects into the container map
// and update containerCount at the same time
func (m *objectMap) retrieveObj(num int) Object {
	obj, _ := m.store.Load(num)
	return obj.(Object)
}

// ChannelObject represents a goby channel, which carries a golang channel
type ChannelObject struct {
	*baseObj
	Chan chan int
}

func (co *ChannelObject) Value() interface{} {
	return co.Chan
}

// Polymorphic helper functions -----------------------------------------
func (co *ChannelObject) toString() string {
	return fmt.Sprintf("<Channel: %p>", co.Chan)
}

func (co *ChannelObject) toJSON() string {
	return co.toString()
}

func (co *ChannelObject) copy() Object {
	newC := &ChannelObject{baseObj: &baseObj{class: co.class}, Chan: make(chan int)}
	return newC
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
