package vm

import (
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/dlclark/regexp2"
	"github.com/goby-lang/goby/vm/classes"
	"github.com/goby-lang/goby/vm/errors"
)

// MatchDataObject represents the match data returned by a regular expression matching operation.
//
// ```ruby
// 'abcd'.match(Regexp.new('(b.)')) #=> #<MatchData "bc" 1:"bc">
// ```
//
// - `MatchData.new` is not supported.
type MatchDataObject struct {
	*baseObj
	captures  []string
	positions []int
	pattern   string // original regex
	text      string // original text
}

// Class methods --------------------------------------------------------
func builtInMatchDataClassMethods() []*BuiltinMethodObject {
	return []*BuiltinMethodObject{
		{
			Name: "new",
			Fn: func(receiver Object) builtinMethodBody {
				return func(t *thread, args []Object, blockFrame *normalCallFrame) Object {
					return t.initUnsupportedMethodError("#new", receiver)
				}
			},
		},
	}
}

// Instance methods -----------------------------------------------------
func builtinMatchDataInstanceMethods() []*BuiltinMethodObject {
	return []*BuiltinMethodObject{
		{
			// Returns the array of captures; equivalent to `match.to_a[1..-1]`.
			//
			// ```ruby
			// c1, c2 = 'abcd'.match(Regexp.new('a(b)(c)'))
			// c1    #=> "b"
			// c2    #=> "c"
			// ```
			//
			// @return [Array]
			Name: "captures",
			Fn: func(receiver Object) builtinMethodBody {
				return func(t *thread, args []Object, blockFrame *normalCallFrame) Object {
					if len(args) != 0 {
						return t.vm.initErrorObject(errors.ArgumentError, "Expect 0 argument. got=%d", len(args))
					}

					matchData, _ := receiver.(*MatchDataObject)

					sourceCaptures := matchData.captures[1:]
					destCaptures := make([]Object, len(sourceCaptures), len(sourceCaptures))

					for i, capture := range sourceCaptures {
						destCaptures[i] = t.vm.initStringObject(capture)
					}

					return t.vm.initArrayObject(destCaptures)
				}
			},
		},
		{
			// Returns the array of captures.
			//
			// ```ruby
			// c0, c1, c2 = 'abcd'.match(Regexp.new('a(b)(c)'))
			// c0    #=> "abc"
			// c1    #=> "b"
			// c2    #=> "c"
			// ```
			//
			// @return [Array]
			Name: "to_a",
			Fn: func(receiver Object) builtinMethodBody {
				return func(t *thread, args []Object, blockFrame *normalCallFrame) Object {
					if len(args) != 0 {
						return t.vm.initErrorObject(errors.ArgumentError, "Expect 0 argument. got=%d", len(args))
					}

					matchData, _ := receiver.(*MatchDataObject)

					destCaptures := make([]Object, len(matchData.captures), len(matchData.captures))

					for i, capture := range matchData.captures {
						destCaptures[i] = t.vm.initStringObject(capture)
					}

					return t.vm.initArrayObject(destCaptures)
				}
			},
		},
	}
}

// Internal functions ===================================================

// Functions for initialization -----------------------------------------

// Initializes a MatchDataObject from a Match object, and the original pattern/text.
// Nothing prevents the programmer to pass pattern/text unrelated to the match, but this will
// create an inconsistent MatchData object.
func (vm *VM) initMatchDataObject(match *regexp2.Match, pattern, text string) *MatchDataObject {
	captures := make([]string, len(match.Groups()), len(match.Groups()))
	positions := make([]int, len(match.Groups()), len(match.Groups()))

	for i, group := range match.Groups() {
		// Using as reference the Ruby MatchData implementation, we have a crucial difference with this
		// Go implementation; the former stores only the first capture when matching, while the latter
		// stores all, and uses as reference the last.
		captures[i] = group.Captures[0].String()
		positions[i] = group.Captures[0].Index
	}

	return &MatchDataObject{
		baseObj:   &baseObj{class: vm.topLevelClass(classes.MatchDataClass)},
		captures:  captures,
		positions: positions,
		pattern:   pattern,
		text:      text,
	}
}

func (vm *VM) initMatchDataClass() *RClass {
	klass := vm.initializeClass(classes.MatchDataClass, false)
	klass.setBuiltinMethods(builtinMatchDataInstanceMethods(), false)
	klass.setBuiltinMethods(builtInMatchDataClassMethods(), true)
	return klass
}

// Polymorphic helper functions -----------------------------------------

// redirects to toString()
func (m *MatchDataObject) Value() interface{} {
	return m.toString()
}

// returns a string representation of the object
func (m *MatchDataObject) toString() string {
	result := "#<MatchData"

	for i, capture := range m.captures {
		if i == 0 {
			result += fmt.Sprintf(" \"%s\"", capture)
		} else {
			result += fmt.Sprintf(" %d:\"%s\"", i, capture)
		}
	}

	result += ">"

	return result
}

// returns a `{ captureNumber: captureValue }` JSON-encoded string
func (m *MatchDataObject) toJSON() string {
	capturesMap := make(map[int]string)

	for i, capture := range m.captures {
		capturesMap[i] = capture
	}

	capturesJson, _ := json.Marshal(capturesMap)

	return string(capturesJson)
}

// equal checks if the string values between receiver and argument are equal
func (m *MatchDataObject) equal(other *MatchDataObject) bool {
	return reflect.DeepEqual(m.captures, other.captures) &&
		reflect.DeepEqual(m.positions, other.positions) &&
		m.pattern == other.pattern &&
		m.text == other.text
}
