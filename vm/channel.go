package vm

import (
	"fmt"
	"github.com/goby-lang/goby/vm/classes"
	"sync"
)

// ChannelObject represents Goby's "channel", which equips the Golang' channel and works with `thread`.
// `thread` is actually a "goroutine".
// A channel object can relay any kind of objects and guarantees thread-safe communications.
// You should always use channel objects for safe communications between threads.
// `Channel#new` is available.
//
// ```ruby
// def f(from)
//   i = 0
//   while i < 3 do
//     puts(from + ": " + i.to_s)
//     i += 1
//   end
// end
//
// f("direct")
//
// c = Channel.new    # spawning a channel object
//
// thread do
//   puts(c.receive)
//   f("thread")
// end
//
// thread do
//   puts("going")
//   c.deliver(10)
// end
//
// sleep(2) # This is to prevent main program finished before goroutine.
// ```
//
// Note that the possibility of race conditions still exists. Handle them with care.
//
// ```ruby
// c = Channel.new
//
// i = 0
// thread do
//   i += 1
//   c.deliver(i)   # sends `i` to channel `c`
// end
//
// # If we put a bare `i += 1` here, then it will execute along with other thread,
// # which will cause a race condition.
// # The following "receive" is needed to block the main process until thread is finished
// c.receive
// i += 1
//
// c.close           # you should call this to close the channel explicitly
// ```
type ChannelObject struct {
	*baseObj
	Chan chan int
}

// Class methods --------------------------------------------------------
func builtinChannelClassMethods() []*BuiltinMethodObject {
	return []*BuiltinMethodObject{
		{
			// Creates an instance of `Channel` class, taking no arguments.
			//
			// ```ruby
			// c = Channel.new
			// c.class         #=> Channel
			// ```
			//
			// @return [Channel]
			Name: "new",
			Fn: func(receiver Object, sourceLine int) builtinMethodBody {
				return func(t *Thread, args []Object, blockFrame *normalCallFrame) Object {
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
			// Closes and releases the channel.
			// You should explicitly call `close` when threading is finished.
			//
			// ```ruby
			// c = Channel.new
			//
			// 1001.times do |i|
			// 	 thread do
			//     c.deliver(i)
			//	 end
			// end
			//
			// r = 0
			// 1001.times do
			//   r = r + c.receive
			// end
			//
			// c.close           # close the channel
			//
			// puts(r)
			// ```
			//
			// It takes no argument. TODO: add argument checking.
			//
			// @return [Null]
			Name: "close",
			Fn: func(receiver Object, sourceLine int) builtinMethodBody {
				return func(t *Thread, args []Object, blockFrame *normalCallFrame) Object {
					c := receiver.(*ChannelObject)

					close(c.Chan)

					return NULL
				}
			},
		},
		{
			// Sends an object to the receiver (channel), then returns the object.
			// Note that the method suspends the process until the object is actually received.
			// Thus if you call `deliver` outside thread, the main process would suspend.
			// Note that you don't need to send dummy object just to resume; use `close` instead.
			//
			// ```ruby
			// c = Channel.new
			//
			// i = 0
			// thread do
			//   i += 1
			//   c.deliver(i)   # sends `i` to channel `c`
			// end
			//
			// c.receive        # receives `i`
			// ```
			//
			// It takes 1 argument. TODO: add argument checking.
			//
			// @param object [Object]
			// @return [Object]
			Name: "deliver",
			Fn: func(receiver Object, sourceLine int) builtinMethodBody {
				return func(t *Thread, args []Object, blockFrame *normalCallFrame) Object {
					id := t.vm.channelObjectMap.storeObj(args[0])

					c := receiver.(*ChannelObject)

					c.Chan <- id

					return args[0]
				}
			},
		},
		{
			// Receives objects from other threads' `deliver` method, then returns it.
			// The method works as if the channel would receive objects perpetually from outer space.
			// Note that the method suspends the process until it actually receives something via `deliver`.
			// Thus if you call `receive` outside thread, the main process would suspend.
			// This also means you can resume a code by using the `receive` method.
			//
			// ```ruby
			// c = Channel.new
			//
			// thread do
			//   puts(c.receive)    # prints the object received from other threads.
			//   f("thread")
			// end
			// ```
			//
			// It takes no arguments. TODO: add argument checking.
			//
			// @return [Object]
			Name: "receive",
			Fn: func(receiver Object, sourceLine int) builtinMethodBody {
				return func(t *Thread, args []Object, blockFrame *normalCallFrame) Object {
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

// Value returns the object
func (co *ChannelObject) Value() interface{} {
	return co.Chan
}

// toString returns the object's name as the string format
func (co *ChannelObject) toString() string {
	return fmt.Sprintf("<Channel: %p>", co.Chan)
}

// toJSON just delegates to toString
func (co *ChannelObject) toJSON(t *Thread) string {
	return co.toString()
}

// copy returns the duplicate of the Array object
func (co *ChannelObject) copy() Object {
	newC := &ChannelObject{baseObj: &baseObj{class: co.class}, Chan: make(chan int)}
	return newC
}

// objectMap ==========================================================

type objectMap struct {
	store *sync.Map
}

// Polymorphic helper functions -----------------------------------------

// storeObj stores objects into the container map
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
