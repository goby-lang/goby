package vm

import (
	"fmt"
	"sync"

	"github.com/goby-lang/goby/vm/classes"
)

// ChannelObject represents a goby channel, which carries a golang channel
type ChannelObject struct {
	*baseObj
	Chan chan int
}

// Class methods --------------------------------------------------------
func builtinChannelClassMethods() []*BuiltinMethodObject {
	return []*BuiltinMethodObject{
		{
			Name: "new",
			Fn: func(receiver Object) builtinMethodBody {
				return func(t *thread, args []Object, blockFrame *callFrame) Object {
					c := &ChannelObject{baseObj: &baseObj{class: t.vm.topLevelClass(classes.ChannelClass)}, Chan: make(chan int)}
					return c
				}
			},
		},
	}
}

// Instance methods -----------------------------------------------------
func builtinChannelInstanceMethods() []*BuiltinMethodObject {
	return []*BuiltinMethodObject{
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

// Internal functions ===================================================

// Functions for initialization -----------------------------------------

func (vm *VM) initChannelClass() *RClass {
	class := vm.initializeClass(classes.ChannelClass, false)
	class.setBuiltinMethods(builtinChannelClassMethods(), true)
	class.setBuiltinMethods(builtinChannelInstanceMethods(), false)
	return class
}

// Polymorphic helper functions -----------------------------------------

// Returns the object
func (co *ChannelObject) Value() interface{} {
	return co.Chan
}

// Returns the object's name as the string format
func (co *ChannelObject) toString() string {
	return fmt.Sprintf("<Channel: %p>", co.Chan)
}

// Alias of toString
func (co *ChannelObject) toJSON() string {
	return co.toString()
}

// Returns the duplicate of the Array object
func (co *ChannelObject) copy() Object {
	newC := &ChannelObject{baseObj: &baseObj{class: co.class}, Chan: make(chan int)}
	return newC
}

// objectMap ==========================================================

type objectMap struct {
	store *sync.Map
}

// Polymorphic helper functions -----------------------------------------

// storeObj store objects into the container map
// and update containerCount at the same time
func (m *objectMap) storeObj(obj Object) int {
	m.store.Store(obj.id(), obj)

	return obj.id()
}

// retrieveObj returns the objects with the number specified
func (m *objectMap) retrieveObj(num int) Object {
	obj, _ := m.store.Load(num)
	return obj.(Object)
}
