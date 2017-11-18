package vm

import (
	"github.com/goby-lang/goby/vm/classes"
	"github.com/goby-lang/goby/vm/errors"
	"strings"
	"time"
)

// TimeObject represents an absolute point on a time flow, in year/month/day/hour/minute/second or RFC/ISO representations or like that, considering timezone, geo-location, DST, and calendar system.
// TBD: Note that TimeObject itself does not contain a concept of "duration".
// Duration is represented by using DurationObject.
//
// The followings are implemented in lib/time.gb (standard lib)
// - `Time.now`
//
// ```ruby
// Time.now #=> 2017-11-09 23:10:49 +0900
// Time.new #=> 2017-11-09 23:10:49 +0900
// Time.new('2017-05-30')               #=> 2017-05-30 00:00:00 +0000 UTC
// Time.new('2017-05-30 18:00')         #=> 2017-05-30 18:00:00 +0000 UTC
// Time.new('2017-05-30 18:00:34')      #=> 2017-05-30 18:00:34 +0000 UTC
// Time.new('2017-05-30 9:00')          #=> 2017-05-30 09:00:00 +0000 UTC
// Time.new('2017-05-30 9:00:56')       #=> 2017-05-30 09:00:56 +0000 UTC
// Time.new('2017-05-30 09:00 JST')     #=> 2017-05-30 09:00:00 +0900 JST
// Time.new('2017-05-30 09:00:59 JST')  #=> 2017-05-30 09:00:59 +0900 JST
// ```
//
// - `Time.new` is supported.
type Time = time.Time
type TimeObject struct {
	*baseObj
	value Time
}

// Class methods --------------------------------------------------------
func builtinTimeClassMethods() []*BuiltinMethodObject {
	return []*BuiltinMethodObject{
		{
			Name: "new",
			Fn: func(receiver Object, sourceLine int) builtinMethodBody {
				return func(t *thread, args []Object, blockFrame *normalCallFrame) Object {
					switch len(args) {
					case 0:
						return t.vm.initTimeObject(currentTime())
					case 1:
						arg, ok := args[0].(*StringObject)
						if !ok {
							return t.vm.initErrorObject(errors.TypeError, sourceLine, errors.WrongArgumentTypeFormat, classes.StringClass, arg.Class().Name)
						}
						if current, err := parseTime(arg.toString()); err != nil {
							return t.vm.initErrorObject(errors.ArgumentError, sourceLine, "Invalid time format. got=%s", arg.value)
						} else {
							return t.vm.initTimeObject(current)
						}
					default:
						return t.vm.initErrorObject(errors.ArgumentError, sourceLine, "Expect 0 or 1 argument. got=%d", len(args))
					}
				}
			},
		},
	}
}

// Instance methods -----------------------------------------------------
func builtinTimeInstanceMethods() []*BuiltinMethodObject {
	return []*BuiltinMethodObject{
		{
			// Parses a string-format time/date/timezone and updates the Time object.
			Name: "parse",
			Fn: func(receiver Object, sourceLine int) builtinMethodBody {
				return func(t *thread, args []Object, blockFrame *normalCallFrame) Object {
					if len(args) != 1 {
						return t.vm.initErrorObject(errors.ArgumentError, sourceLine, "Expect 1 argument. got=%d", len(args))
					}

					arg, ok := args[0].(*StringObject)
					if !ok {
						return t.vm.initErrorObject(errors.TypeError, sourceLine, errors.WrongArgumentTypeFormat, classes.StringClass, arg.Class().Name)
					}

					ts, err := parseTime(arg.toString())
					if err != nil {
						return t.vm.initErrorObject(errors.ArgumentError, sourceLine, "Invalid time format. got=%s", arg.value)
					}

					m := receiver.(*TimeObject)
					m.value = ts
					return m
				}
			},
		},
		{
			// Converts a Time object into a fixed-format string.
			Name: "to_s",
			Fn: func(receiver Object, sourceLine int) builtinMethodBody {
				return func(t *thread, args []Object, blockFrame *normalCallFrame) Object {
					m := receiver.(*TimeObject)
					return t.vm.initStringObject(m.toString())
				}
			},
		},
	}
}

// Internal functions ===================================================

// Functions for initialization -----------------------------------------

func (vm *VM) initTimeObject(t Time) *TimeObject {
	return &TimeObject{
		baseObj: &baseObj{class: vm.topLevelClass(classes.TimeClass)},
		value:   t,
	}
}

func (vm *VM) initTimeClass() *RClass {
	tc := vm.initializeClass("Time", true)
	tc.setBuiltinMethods(builtinTimeInstanceMethods(), false)
	tc.setBuiltinMethods(builtinTimeClassMethods(), true)
	vm.objectClass.setClassConstant(tc)
	vm.libFiles = append(vm.libFiles, "time.gb")
	return tc
}

// Polymorphic helper functions -----------------------------------------

// Value returns the object
func (t *TimeObject) Value() interface{} {
	return t.timeValue()
}

// Numeric interface
func (t *TimeObject) timeValue() Time {
	return t.value
}

// toString returns the object's value as the string format, in non
// exponential format (straight number, without exponent `E<exp>`).
func (t *TimeObject) toString() string {
	return t.value.String()
}

// toJSON just delegates to toString
func (t *TimeObject) toJSON() string {
	return t.toString()
}

// equal checks if the Float values between receiver and argument are equal
func (t *TimeObject) equal(e *TimeObject) bool {
	return t.value == e.value
}

// Other helper functions -----------------------------------------------

// Obtains current local time
func currentTime() Time {
	return time.Now()
}

// Parses the date/time strings into Time object
// Follow https://golang.org/src/time/format.go for formatting
func parseTime(t string) (Time, error) {
	const (
		dateForm            = "2006-01-02"
		dateTimeForm        = "2006-01-02 15:04"
		dateTimeSecForm     = "2006-01-02 15:04:05"
		dateTimezoneForm    = "2006-01-02 15:04 MST"
		dateTimezoneSecForm = "2006-01-02 15:04:05 MST"
	)
	switch strings.Count(t, ":") {
	case 0:
		return time.Parse(dateForm, t)
	case 1:
		if strings.Count(t, " ") > 1 {
			return time.Parse(dateTimezoneForm, t)
		} else {
			return time.Parse(dateTimeForm, t)
		}
	case 2:
		if strings.Count(t, " ") > 1 {
			return time.Parse(dateTimezoneSecForm, t)
		} else {
			return time.Parse(dateTimeSecForm, t)
		}
	default:
		return time.Parse("", t)
	}
}
