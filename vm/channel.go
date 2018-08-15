package vm

import (
	"fmt"
	"sync"

	"github.com/goby-lang/goby/vm/classes"
	"github.com/goby-lang/goby/vm/errors"
)

// ChannelObject represents Goby's "channel", which equips the Golang' channel and works with `thread`.
// `thread` is actually a "goroutine".
// A channel object can relay any kind of objects and guarantees thread-safe communications.
// You should always use channel objects for safe communications between threads.
// `Channel#new` is available.
//
// Note that channels are not like files and you don't need to explicitly close them (e.g.: exiting a loop).
// See https://tour.golang.org/concurrency/4
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
//   c.deliver(i)     # sends `i` to channel `c`
// end
//
// # If we put a bare `i += 1` here, then it will execute along with other thread,
// # which will cause a race condition.
// # The following "receive" is needed to block the main process until thread is finished
// c.receive
// i += 1
//
// c.close           # Redundant: just for explanation and you don't need to call this here
// ```
type ChannelObject struct {
	*BaseObj
	Chan         chan int
	ChannelState int
}

// Channel's state.
// To Goby language contributors: Golang's channels should be carefully handled because:
// - You cannot write to closed channels, or got a panic.
// - You cannot close closed channels, or got a panic.
// - You cannot write to nil channels, or causes a deadlock.
// - You cannot read nil channels, or causes a deadlock.
// Ref: https://beatsync.net/main/log20150325.html
const (
	chOpen = iota
	chClosed
)

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
			Fn: func(receiver Object, sourceLine int, t *Thread, args []Object, blockFrame *normalCallFrame) Object {
				c := &ChannelObject{BaseObj: &BaseObj{class: t.vm.TopLevelClass(classes.ChannelClass)}, Chan: make(chan int, chOpen)}
				return c
			},
		},
	}
}

// Instance methods -----------------------------------------------------
func builtinChannelInstanceMethods() []*BuiltinMethodObject {
	return []*BuiltinMethodObject{
		{
			// Just to close and the channel to declare no more objects will be sent.
			// Channel is not like files, and you don't need to call `close` explicitly unless
			// you definitely need to notify that no more objects will be sent,
			// Well, you can call `#close` against the same channel twice or more, which is redundant.
			// (Go's channel cannot do that)
			// See https://tour.golang.org/concurrency/4
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
			// If you call `close` twice against the same channel, an error is returned.
			//
			// It takes no argument.
			//
			// @return [Null]
			Name: "close",
			Fn: func(receiver Object, sourceLine int, t *Thread, args []Object, blockFrame *normalCallFrame) Object {
				if len(args) != 0 {
					return t.vm.InitErrorObject(errors.ArgumentError, sourceLine, errors.WrongNumberOfArgumentFormat, 0, len(args))
				}
				c := receiver.(*ChannelObject)

				if c.ChannelState == chClosed {
					return t.vm.InitErrorObject(errors.ChannelCloseError, sourceLine, errors.ChannelIsClosed)
				}
				c.ChannelState = chClosed

				close(receiver.(*ChannelObject).Chan)
				receiver = nil
				return NULL
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
			// If you call `deliver` against the closed channel, an error is returned.
			//
			// It takes 1 argument.
			//
			// @param object [Object]
			// @return [Object]
			Name: "deliver",
			Fn: func(receiver Object, sourceLine int, t *Thread, args []Object, blockFrame *normalCallFrame) Object {
				if len(args) != 1 {
					return t.vm.InitErrorObject(errors.ArgumentError, sourceLine, errors.WrongNumberOfArgumentFormat, 1, len(args))
				}
				c := receiver.(*ChannelObject)

				if c.ChannelState == chClosed {
					return t.vm.InitErrorObject(errors.ChannelCloseError, sourceLine, errors.ChannelIsClosed)
				}

				id := t.vm.channelObjectMap.storeObj(args[0])
				c.Chan <- id

				return args[0]
			},
		},
		{
			// Receives objects from other threads' `deliver` method, then returns it.
			// The method works as if the channel would receive objects perpetually from outside.
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
			// If you call `receive` against the closed channel, an error is returned.
			//
			// It takes no arguments.
			//
			// @return [Object]
			Name: "receive",
			Fn: func(receiver Object, sourceLine int, t *Thread, args []Object, blockFrame *normalCallFrame) Object {
				if len(args) != 0 {
					if len(args) != 0 {
						return t.vm.InitErrorObject(errors.ArgumentError, sourceLine, errors.WrongNumberOfArgumentFormat, 0, len(args))
					}
				}
				c := receiver.(*ChannelObject)

				if c.ChannelState == chClosed {
					return t.vm.InitErrorObject(errors.ChannelCloseError, sourceLine, errors.ChannelIsClosed)
				}

				num := <-c.Chan

				return t.vm.channelObjectMap.retrieveObj(num)
			},
		},
	}
}

// Internal functions ===================================================

// Functions for initialization -----------------------------------------

func (vm *VM) initChannelClass() *RClass {
	class := vm.initializeClass(classes.ChannelClass)
	class.setBuiltinMethods(builtinChannelClassMethods(), true)
	class.setBuiltinMethods(builtinChannelInstanceMethods(), false)
	return class
}

// Polymorphic helper functions -----------------------------------------

// Value returns the object
func (co *ChannelObject) Value() interface{} {
	return co.Chan
}

// ToString returns the object's name as the string format
func (co *ChannelObject) ToString() string {
	return fmt.Sprintf("<Channel: %p>", co.Chan)
}

// Inspect delegates to ToString
func (co *ChannelObject) Inspect() string {
	return co.ToString()
}

// ToJSON just delegates to ToString
func (co *ChannelObject) ToJSON(t *Thread) string {
	return co.ToString()
}

// copy returns the duplicate of the Array object
func (co *ChannelObject) copy() Object {
	newC := &ChannelObject{BaseObj: &BaseObj{class: co.class}, Chan: make(chan int)}
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
