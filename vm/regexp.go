package vm

import (
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
// ToDo: Regexp literals with '/.../'
type Regexp = regexp2.Regexp
type RegexpObject struct {
	*BaseObj
	regexp *Regexp
}

// Class methods --------------------------------------------------------
func builtInRegexpClassMethods() []*BuiltinMethodObject {
	return []*BuiltinMethodObject{
		{
			Name: "new",
			Fn: func(receiver Object, sourceLine int, t *Thread, args []Object, blockFrame *normalCallFrame) Object {
				if len(args) != 1 {
					return t.vm.InitErrorObject(errors.ArgumentError, sourceLine, "Expect 1 argument. got=%d", len(args))
				}

				arg, ok := args[0].(*StringObject)
				if !ok {
					return t.vm.InitErrorObject(errors.TypeError, sourceLine, errors.WrongArgumentTypeFormat, classes.StringClass, arg.Class().Name)
				}

				r := t.vm.initRegexpObject(args[0].ToString())
				if r == nil {
					return t.vm.InitErrorObject(errors.ArgumentError, sourceLine, "Invalid regexp: %v", args[0].ToString())
				}
				return r

			},
		},
	}
}

// Instance methods -----------------------------------------------------
func builtinRegexpInstanceMethods() []*BuiltinMethodObject {
	return []*BuiltinMethodObject{

		{
			// Returns true if the two regexp patterns are exactly the same, or returns false if not.
			// If comparing with non Regexp class, just returns false.
			//
			// ```ruby
			// r1 = Regexp.new("goby[0-9]+")
			// r2 = Regexp.new("goby[0-9]+")
			// r3 = Regexp.new("Goby[0-9]+")
			//
			// r1 == r2   # => true
			// r1 == r2   # => false
			// ```
			//
			// @param regexp [Regexp]
			// @return [Boolean]
			Name: "==",
			Fn: func(receiver Object, sourceLine int, t *Thread, args []Object, blockFrame *normalCallFrame) Object {

				if len(args) != 1 {
					return t.vm.InitErrorObject(errors.ArgumentError, sourceLine, "Expect 1 argument. got=%d", len(args))
				}

				right, ok := args[0].(*RegexpObject)
				if !ok {
					return FALSE
				}

				left := receiver.(*RegexpObject)

				if left.Value() == right.Value() {
					return TRUE
				}
				return FALSE

			},
		},
		{
			// Returns boolean value to indicate the result of regexp match with the string given. The methods evaluates a String object.
			//
			// ```ruby
			// r = Regexp.new("o")
			// r.match?("pow")  # => true
			// r.match?("gee")  # => false
			// ```
			//
			// @param string [String]
			// @return [Boolean]
			Name: "match?",
			Fn: func(receiver Object, sourceLine int, t *Thread, args []Object, blockFrame *normalCallFrame) Object {

				if len(args) != 1 {
					return t.vm.InitErrorObject(errors.ArgumentError, sourceLine, "Expect 1 argument. got=%d", len(args))
				}

				arg := args[0]
				input, ok := arg.(*StringObject)
				if !ok {
					return t.vm.InitErrorObject(errors.TypeError, sourceLine, errors.WrongArgumentTypeFormat, classes.StringClass, arg.Class().Name)
				}

				re := receiver.(*RegexpObject).regexp
				m, _ := re.MatchString(input.value)

				return toBooleanObject(m)

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
		BaseObj: &BaseObj{class: vm.TopLevelClass(classes.RegexpClass)},
		regexp:  r,
	}
}

func (vm *VM) initRegexpClass() *RClass {
	rc := vm.initializeClass(classes.RegexpClass)
	rc.setBuiltinMethods(builtinRegexpInstanceMethods(), false)
	rc.setBuiltinMethods(builtInRegexpClassMethods(), true)
	return rc
}

// Polymorphic helper functions -----------------------------------------

// Value returns the object
func (r *RegexpObject) Value() interface{} {
	return r.regexp.String()
}

// ToString returns the object's name as the string format
func (r *RegexpObject) ToString() string {
	return r.regexp.String()
}

// ToJSON just delegates to ToString
func (r *RegexpObject) ToJSON(t *Thread) string {
	return "\"" + r.ToString() + "\""
}

// equal checks if the string values between receiver and argument are equal
func (r *RegexpObject) equal(e *RegexpObject) bool {
	return r.ToString() == r.ToString()
}
