package vm

import (
	"fmt"

	"github.com/dlclark/regexp2"
	"github.com/goby-lang/goby/vm/classes"
	"github.com/goby-lang/goby/vm/errors"
)

// MatchDataObject represents the match data returned by a regular expression matching operation.
// You can use named-captures via `(?<name>)`.
//
// ```ruby
// 'abcd'.match(Regexp.new('(b.)'))
// #=> #<MatchData 0:"bc" 1:"bc">
//
// 'abcd'.match(Regexp.new('a(?<first>b)(?<second>c)'))
// #=> #<MatchData 0:"abc" first:"b" second:"c">
// ```
//
// - `MatchData.new` is not supported.
type Match = regexp2.Match
type MatchDataObject struct {
	*BaseObj
	match *Match
}

// Class methods --------------------------------------------------------
func builtInMatchDataClassMethods() []*BuiltinMethodObject {
	return []*BuiltinMethodObject{
		{
			Name: "new",
			Fn: func(receiver Object, sourceLine int, t *Thread, args []Object, blockFrame *normalCallFrame) Object {
				return t.vm.initUnsupportedMethodError(sourceLine, "#new", receiver)

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
			// c1, c2 = 'abcd'.match(Regexp.new('a(b)(c)')).captures
			// c1    #=> "b"
			// c2    #=> "c"
			// ```
			//
			// @return [Array]
			Name: "captures",
			Fn: func(receiver Object, sourceLine int, t *Thread, args []Object, blockFrame *normalCallFrame) Object {
				if len(args) != 0 {
					return t.vm.InitErrorObject(errors.ArgumentError, sourceLine, "Expect 0 argument. got=%d", len(args))
				}
				offset := 1

				g := receiver.(*MatchDataObject).match
				n := len(g.Groups()) - offset
				destCaptures := make([]Object, n, n)

				for i := 0; i < n; i++ {
					destCaptures[i] = t.vm.InitStringObject(g.GroupByNumber(i + offset).String())
				}

				return t.vm.InitArrayObject(destCaptures)

			},
		},
		{
			// Returns the array of captures.
			//
			// ```ruby
			// c0, c1, c2 = 'abcd'.match(Regexp.new('a(b)(c)')).to_a
			// c0    #=> "abc"
			// c1    #=> "b"
			// c2    #=> "c"
			// ```
			//
			// @return [Array]
			Name: "to_a",
			Fn: func(receiver Object, sourceLine int, t *Thread, args []Object, blockFrame *normalCallFrame) Object {
				if len(args) != 0 {
					return t.vm.InitErrorObject(errors.ArgumentError, sourceLine, "Expect 0 argument. got=%d", len(args))
				}

				g := receiver.(*MatchDataObject).match
				n := len(g.Groups())
				destCaptures := make([]Object, n, n)

				for i := 0; i < n; i++ {
					destCaptures[i] = t.vm.InitStringObject(g.GroupByNumber(i).String())
				}

				return t.vm.InitArrayObject(destCaptures)

			},
		},
		{
			// Returns the hash of captures, including the whole matched text(`0:`).
			// You can use named-capture as well.
			//
			// ```ruby
			// h = 'abcd'.match(Regexp.new('a(b)(c)')).to_h
			// puts h #=> { "0": "abc", "1": "b", "2": "c" }
			//
			// h = 'abcd'.match(Regexp.new('a(?<first>b)(?<second>c)')).to_h
			// puts h #=> { "0": "abc", "first": "b", "second": "c" }
			// ```
			//
			// @return [Hash]
			Name: "to_h",
			Fn: func(receiver Object, sourceLine int, t *Thread, args []Object, blockFrame *normalCallFrame) Object {
				if len(args) != 0 {
					return t.vm.InitErrorObject(errors.ArgumentError, sourceLine, "Expect 0 argument. got=%d", len(args))
				}

				groups := receiver.(*MatchDataObject).match
				result := make(map[string]Object)

				for _, g := range groups.Groups() {
					result[g.Name] = t.vm.InitStringObject(g.String())
				}

				return t.vm.InitHashObject(result)

			},
		},
		{
			// Returns the length of the array; equivalent to `match.to_a.length`.
			//
			// ```ruby
			// 'abcd'.match(Regexp.new('a(b)(c)')).length # => 3
			// ```
			// @return [Integer]
			Name: "length",
			Fn: func(receiver Object, sourceLine int, t *Thread, args []Object, blockFrame *normalCallFrame) Object {
				if len(args) != 0 {
					return t.vm.InitErrorObject(errors.ArgumentError, sourceLine, "Expect 0 argument. got=%d", len(args))
				}

				m := receiver.(*MatchDataObject).match

				return t.vm.InitIntegerObject(m.GroupCount())

			},
		},
	}
}

// Internal functions ===================================================

// Functions for initialization -----------------------------------------

// Initializes a MatchDataObject from a Match object, and the original pattern/text.
// Nothing prevents the programmer to pass pattern/text unrelated to the match, but this will
// create an inconsistent MatchData object.
func (vm *VM) initMatchDataObject(match *Match, pattern, text string) *MatchDataObject {
	return &MatchDataObject{
		BaseObj: &BaseObj{class: vm.TopLevelClass(classes.MatchDataClass)},
		match:   match,
	}
}

func (vm *VM) initMatchDataClass() *RClass {
	klass := vm.initializeClass(classes.MatchDataClass)
	klass.setBuiltinMethods(builtinMatchDataInstanceMethods(), false)
	klass.setBuiltinMethods(builtInMatchDataClassMethods(), true)
	return klass
}

// Polymorphic helper functions -----------------------------------------

// redirects to ToString()
func (m *MatchDataObject) Value() interface{} {
	return m.ToString()
}

// returns a string representation of the object
func (m *MatchDataObject) ToString() string {
	result := "#<MatchData"

	for _, c := range m.match.Groups() {
		result += fmt.Sprintf(" %s:\"%s\"", c.Name, c.String())
	}

	result += ">"

	return result
}

// returns a `{ captureNumber: captureValue }` JSON-encoded string
func (m *MatchDataObject) ToJSON(t *Thread) string {
	result := "{"

	for _, c := range m.match.Groups() {
		result += fmt.Sprintf(" %s:\"%s\"", c.Name, c.String())
	}

	result += "}"

	return result
}

// equal checks if the string values between receiver and argument are equal
func (m *MatchDataObject) equal(other *MatchDataObject) bool {
	return m.match == other.match
}
