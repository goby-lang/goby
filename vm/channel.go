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

var objMap = &objectMap{store: map[int]Object{}}

// storeObj store objects into the container map
// and update containerCount at the same time
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

func (m *objectMap) retrieveObj(num int) Object {
	m.RLock()

	defer m.RUnlock()
	return m.store[num]
}

var channelID = 0

// ChannelObject represents a goby channel, which carries a golang channel
type ChannelObject struct {
	id    int
	Class *RClass
	Chan  chan int
}

func initializeChannelClass() *RClass {
	class := initializeClass("Channel", false)
	class.setBuiltInMethods(builtinChannelClassMethods(), true)
	class.setBuiltInMethods(builtinChannelInstanceMethods(), false)
	return class
}

func (co *ChannelObject) toString() string {
	return fmt.Sprintf("<Channel: %d>", co.id)
}

func (co *ChannelObject) toJSON() string {
	return co.toString()
}

func (co *ChannelObject) returnClass() Class {
	return co.Class
}

func builtinChannelClassMethods() []*BuiltInMethodObject {
	return []*BuiltInMethodObject{
		{
			Name: "new",
			Fn: func(receiver Object) builtinMethodBody {
				return func(t *thread, args []Object, blockFrame *callFrame) Object {
					c := &ChannelObject{Class: t.vm.builtInClasses["Channel"], id: channelID, Chan: make(chan int)}
					channelID++
					return c
				}
			},
		},
	}
}

func builtinChannelInstanceMethods() []*BuiltInMethodObject {
	return []*BuiltInMethodObject{
		{
			Name: "deliver",
			Fn: func(receiver Object) builtinMethodBody {
				return func(t *thread, args []Object, blockFrame *callFrame) Object {
					id := objMap.storeObj(args[0])

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

					return objMap.retrieveObj(num)
				}
			},
		},
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
	}
}
