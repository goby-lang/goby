package vm

import (
	"strconv"

	"github.com/dlclark/regexp2"
	"github.com/goby-lang/goby/vm/classes"
	"github.com/goby-lang/goby/vm/errors"
)

// RegexpObject represents regexp instances, which of the type is actually string.
// Regexp object holds regexp strings.
// Powered by: github.com/dlclark/regexp2
// The regexp2 package had been ported from .NET Framework's regexp library, which is PCRE-compatible
// and is almost all equivalent to Ruby's Onigmo regexp library.
//
// ```ruby
// a = Regexp.new("orl")
// a.match?("Hello World")   #=> true
// a.match?("Hello Regexp")  #=> false
//
// b = Regexp.new("😏")
// b.match?("🤡 😏 😐")   #=> true
// b.match?("😝 😍 😊")   #=> false
//
// c = Regexp.new("居(ら(?=れ)|さ(?=せ)|る|ろ|れ(?=[ばる])|よ|(?=な[いかくけそ]|ま[しすせ]|そう|た|て))")
// c.match?("居られればいいのに")  #=> true
// c.match?("居ずまいを正す")      #=> false
// ```
//
// **Note:**
//
// - Currently, manipulations are based upon Golang's Unicode manipulations.
// - Currently, UTF-8 encoding is assumed based upon Golang's string manipulation, but the encoding is not actually specified(TBD).
// - `Regexp.new` is exceptionally supported.
//
// **To Goby maintainers**: avoid using Go's standard regexp package (slow and not rich). Consider the faster `Trim` or `Split` etc in Go's "strings" package first, or just use the dlclark/regexp2 instead.
// ToDo: Regexp literals with '/.../' and match operator `=~`
type RegexpObject struct {
	*baseObj
	Regexp *regexp2.Regexp
}

// Class methods --------------------------------------------------------
func builtInRegexpClassMethods() []*BuiltinMethodObject {
	return []*BuiltinMethodObject{
		{
			Name: "new",
			Fn: func(receiver Object) builtinMethodBody {
				return func(t *thread, args []Object, blockFrame *callFrame) Object {
					r := t.vm.initRegexpObject(args[0].toString())
					if r == nil {
						return t.vm.initErrorObject(errors.ArgumentError, "Invalid regexp: %v", args[0].toString())
					}
					return r
				}
			},
		},
	}
}

// Instance methods -----------------------------------------------------
func builtinRegexpInstanceMethods() []*BuiltinMethodObject {
	return []*BuiltinMethodObject{

		{
			// Returns boolean value to indicate the result of regexp match with the string given.
			//
			// ```ruby
			// r = Regexp.new("o")
			// r.match("pow")  # => true
			// r.match("gee")  # => false
			// ```
			//
			// @param string [String]
			// @return [Boolean]
			Name: "match?",
			Fn: func(receiver Object) builtinMethodBody {
				return func(t *thread, args []Object, blockFrame *callFrame) Object {

					if len(args) != 1 {
						return t.vm.initErrorObject(errors.ArgumentError, "Expect 1 argument. got=%v", strconv.Itoa(len(args)))
					}

					re := receiver.(*RegexpObject).Regexp
					m, _ := re.MatchString(args[0].toString())
					if m {
						return TRUE
					}
					return FALSE
				}
			},
		},
	}
}

// Internal functions ===================================================

// Functions for initialization -----------------------------------------

func (vm *VM) initRegexpObject(regexp string) *RegexpObject {
	r, err := regexp2.Compile(regexp, 0)
	if err != nil {
		return nil
	}
	return &RegexpObject{
		baseObj: &baseObj{class: vm.topLevelClass(classes.RegexpClass)},
		Regexp:  r,
	}
}

func (vm *VM) initRegexpClass() *RClass {
	rc := vm.initializeClass(classes.RegexpClass, false)
	rc.setBuiltinMethods(builtinRegexpInstanceMethods(), false)
	rc.setBuiltinMethods(builtInRegexpClassMethods(), true)
	return rc
}

// Polymorphic helper functions -----------------------------------------

// Value returns the object
func (r *RegexpObject) Value() interface{} {
	return r.toString()
}

// toString returns the object's name as the string format
func (r *RegexpObject) toString() string {
	return r.Regexp.String()
}

// toJSON just delegates to toString
func (r *RegexpObject) toJSON() string {
	return "\"" + r.toString() + "\""
}

// equal checks if the string values between receiver and argument are equal
func (r *RegexpObject) equal(e *RegexpObject) bool {
	return r.toString() == r.toString()
}
