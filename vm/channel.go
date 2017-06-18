package vm

import (
	"fmt"
)

var (
	channelClass *RClass
)

var containerMap = map[int]Object{}
var containerCount = 0

// storeObj store objects into the container map
// and update containerCount at the same time
func storeObj(obj Object) int {
	mutex.Lock()
	containerMap[containerCount] = obj
	i := containerCount

	// containerCount here can be considered as deliveries' id
	// And this id will be unused once the delivery is completed (which will be really quick)
	// So we can assume that if we reach 1000th delivery,
	// previous deliveries are all completed and they don't need their id anymore.
	if containerCount > 1000 {
		containerCount = 0
	} else {
		containerCount += 1
	}

	mutex.Unlock()
	return i
}

type ChannelObject struct {
	Id    int
	Class *RClass
	Chan  chan int
}

func initializeChannelClass() {
	class := initializeClass("Channel", false)
	class.setBuiltInMethods(builtinChannelClassMethods(), true)
	class.setBuiltInMethods(builtinChannelInstanceMethods(), false)
	objectClass.constants["Channel"] = &Pointer{Target: class}
	channelClass = class
}

func (co *ChannelObject) toString() string {
	return fmt.Sprintf("<Channel: %d>", co.Id)
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
					return &ChannelObject{Class: channelClass, Id: 0, Chan: make(chan int)}
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
					id := storeObj(args[0])

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

					mutex.Lock()

					defer mutex.Unlock()
					return containerMap[num]
				}
			},
		},
	}
}
