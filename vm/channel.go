package vm

import (
	"fmt"
)

var channelClass *RClass

type ChannelObject struct {
	Id int
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

// toString returns detailed infoof a array include elements it contains
func (co *ChannelObject) toString() string {
	return fmt.Sprintf("<Channel: %d>", co.Id)
}

func (co *ChannelObject) toJSON() string {
	return co.toString()
}

// returnClass returns current object's class, which is RArray
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
					num := args[0].(*IntegerObject).Value
					c := receiver.(*ChannelObject)

					c.Chan <- num

					return c
				}
			},
		},
		{
			Name: "receive",
			Fn: func(receiver Object) builtinMethodBody {
				return func(t *thread, args []Object, blockFrame *callFrame) Object {
					c := receiver.(*ChannelObject)

					num := <-c.Chan

					return initilaizeInteger(num)
				}
			},
		},
	}
}