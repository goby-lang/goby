package vm

import (
	"github.com/dlclark/regexp2"
	"github.com/goby-lang/goby/vm/classes"
	"github.com/goby-lang/goby/vm/errors"
)

// Regexp is an alias of regexp2.Regexp
type Regexp = regexp2.Regexp

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
type RegexpObject struct {
	*BaseObj
	regexp *Regexp
}

// Class methods --------------------------------------------------------
var builtInRegexpClassMethods = []*BuiltinMethodObject{
	{
		Name: "new",
		Fn: func(receiver Object, sourceLine int, t *Thread, args []Object, blockFrame *normalCallFrame) Object {
			if len(args) != 1 {
				return t.vm.InitErrorObject(errors.ArgumentError, sourceLine, errors.WrongNumberOfArgument, 1, len(args))
			}

			arg, ok := args[0].(*StringObject)
			if !ok {
				return t.vm.InitErrorObject(errors.TypeError, sourceLine, errors.WrongArgumentTypeFormat, classes.StringClass, arg.Class().Name)
			}

			r := t.vm.initRegexpObject(args[0].ToString())
			if r == nil {
				return t.vm.InitErrorObject(errors.ArgumentError, sourceLine, "Invalid regexp: %v", args[0].Inspect())
			}
			return r

		},
	},
}

// Instance methods -----------------------------------------------------
var builtinRegexpInstanceMethods = []*BuiltinMethodObject{
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
				return t.vm.InitErrorObject(errors.ArgumentError, sourceLine, errors.WrongNumberOfArgument, 1, len(args))
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

// Internal functions ===================================================

// Functions for initialization -----------------------------------------

func (vm *VM) initRegexpObject(regexp string) *RegexpObject {
	r, err := regexp2.Compile(regexp, 0)
	if err != nil {
		return nil
	}
	return &RegexpObject{
		BaseObj: NewBaseObject(vm.TopLevelClass(classes.RegexpClass)),
		regexp:  r,
	}
}

func (vm *VM) initRegexpClass() *RClass {
	rc := vm.initializeClass(classes.RegexpClass)
	rc.setBuiltinMethods(builtinRegexpInstanceMethods, false)
	rc.setBuiltinMethods(builtInRegexpClassMethods, true)
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

// Inspect delegates to ToString
func (r *RegexpObject) Inspect() string {
	return r.ToString()
}

// ToJSON just delegates to ToString
func (r *RegexpObject) ToJSON(t *Thread) string {
	return "\"" + r.ToString() + "\""
}

// equal checks if the string values between receiver and argument are equal
func (r *RegexpObject) equal(e *RegexpObject) bool {
	return r.ToString() == r.ToString()
}

func (r *RegexpObject) equalTo(with Object) bool {
	right, ok := with.(*RegexpObject)
	if !ok {
		return false
	}

	if r.Value() == right.Value() {
		return true
	}

	return false
}
