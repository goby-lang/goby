package vm

import (
	"fmt"
	"strconv"
	"strings"
	"unicode"
)

var (
	stringClass *RString
)

// RString is the built in string class
type RString struct {
	*BaseClass
}

// StringObject represents string instances
// String object holds and manipulates a sequence of characters.
// String objects may be created using as string literals.
// Double or single quotations can be used for representation.
//
// ```ruby
// a = "Three"
// b = 'zero'
// c = '漢'
// d = 'Tiếng Việt'
// e = "😏️️"
// ```
//
// **Note:**
//
// - Currently, manipulations are based upon Golang's Unicode manipulations.
// - Currently, UTF-8 encoding is assumed based upon Golang's string manipulation, but the encoding is not actually specified(TBD).
// - `String.new` is not supported.
type StringObject struct {
	Class *RString
	Value string
}

func (s *StringObject) toString() string {
	return "\"" + s.Value + "\""
}

func (s *StringObject) toJSON() string {
	return "\"" + s.Value + "\""
}

func (s *StringObject) returnClass() Class {
	if s.Class == nil {
		panic(fmt.Sprintf("String %s doesn't have class.", s.toString()))
	}

	return s.Class
}

func (s *StringObject) toArray() *ArrayObject {
	elems := []Object{}

	for i := 0; i < len(s.Value); i++ {
		elems = append(elems, initIntegerObject(i))
	}

	return initArrayObject(elems)
}

func (s *StringObject) equal(e *StringObject) bool {
	return s.Value == e.Value
}

func initStringObject(value string) *StringObject {
	replacer := strings.NewReplacer("\\n", "\n", "\\r", "\r", "\\t", "\t", "\\v", "\v", "\\\\", "\\")
	return &StringObject{Value: replacer.Replace(value), Class: stringClass}
}

var builtinStringInstanceMethods = []*BuiltInMethodObject{
	{
		// Returns the concatenation of self and another String
		//
		// ```Ruby
		// "first" + "-second" # => "first-second"
		// ```
		//
		// @return [String]
		Name: "+",
		Fn: func(receiver Object) builtinMethodBody {
			return func(t *thread, args []Object, blockFrame *callFrame) Object {

				leftValue := receiver.(*StringObject).Value
				right, ok := args[0].(*StringObject)

				if !ok {
					return wrongTypeError(stringClass)
				}

				rightValue := right.Value
				return &StringObject{Value: leftValue + rightValue, Class: stringClass}
			}
		},
	},
	{
		// Returns self multiplying another Integer
		//
		// ```Ruby
		// "string " * 2 # => "string string string "
		// ```
		//
		// @return [String]
		Name: "*",
		Fn: func(receiver Object) builtinMethodBody {
			return func(t *thread, args []Object, blockFrame *callFrame) Object {

				leftValue := receiver.(*StringObject).Value
				right, ok := args[0].(*IntegerObject)

				if !ok {
					return wrongTypeError(stringClass)
				}

				if right.Value < 0 {
					return newError("Second argument must be greater than or equal to 0 String#*")
				}

				var result string

				for i := 0; i < right.Value; i++ {
					result += leftValue
				}

				return &StringObject{Value: result, Class: stringClass}
			}
		},
	},
	{
		// Returns a Boolean if first string greater than second string
		//
		// ```Ruby
		// "a" < "b" # => true
		// ```
		//
		// @return [Boolean]
		Name: ">",
		Fn: func(receiver Object) builtinMethodBody {
			return func(t *thread, args []Object, blockFrame *callFrame) Object {

				leftValue := receiver.(*StringObject).Value
				right, ok := args[0].(*StringObject)

				if !ok {
					return wrongTypeError(stringClass)
				}

				rightValue := right.Value

				if leftValue > rightValue {
					return TRUE
				}

				return FALSE
			}
		},
	},
	{
		// Returns a Boolean if first string less than second string
		//
		// ```Ruby
		// "a" < "b" # => true
		// ```
		//
		// @return [Boolean]
		Name: "<",
		Fn: func(receiver Object) builtinMethodBody {
			return func(t *thread, args []Object, blockFrame *callFrame) Object {

				leftValue := receiver.(*StringObject).Value
				right, ok := args[0].(*StringObject)

				if !ok {
					return wrongTypeError(stringClass)
				}

				rightValue := right.Value

				if leftValue < rightValue {
					return TRUE
				}

				return FALSE
			}
		},
	},
	{
		// Returns a Boolean of compared two strings
		//
		// ```Ruby
		// "first" == "second" # => false
		// "two" == "two" # => true
		// ```
		//
		// @return [Boolean]
		Name: "==",
		Fn: func(receiver Object) builtinMethodBody {
			return func(t *thread, args []Object, blockFrame *callFrame) Object {

				leftValue := receiver.(*StringObject).Value
				right, ok := args[0].(*StringObject)

				if !ok {
					return wrongTypeError(stringClass)
				}

				rightValue := right.Value

				if leftValue == rightValue {
					return TRUE
				}

				return FALSE
			}
		},
	},
	{
		// Returns a Integer. If first string is less than second string returns -1, if equal to returns 0, if greater returns 1
		//
		//
		// ```Ruby
		// "abc" <=> "abcd" # => -1
		// "abc" <=> "abc" # => 0
		// "abcd" <=> "abc" # => 1
		// ```
		//
		// @return [Integer]
		Name: "<=>",
		Fn: func(receiver Object) builtinMethodBody {
			return func(t *thread, args []Object, blockFrame *callFrame) Object {

				leftValue := receiver.(*StringObject).Value
				right, ok := args[0].(*StringObject)

				if !ok {
					return wrongTypeError(stringClass)
				}

				rightValue := right.Value

				if leftValue < rightValue {
					return initIntegerObject(-1)
				}
				if leftValue > rightValue {
					return initIntegerObject(1)
				}

				return initIntegerObject(0)
			}
		},
	},
	{
		// Returns a Boolean of compared two strings
		//
		// ```Ruby
		// "first" != "second" # => true
		// "two" != "two" # => false
		// ```
		//
		// @return [Boolean]
		Name: "!=",
		Fn: func(receiver Object) builtinMethodBody {
			return func(t *thread, args []Object, blockFrame *callFrame) Object {

				leftValue := receiver.(*StringObject).Value
				right, ok := args[0].(*StringObject)

				if !ok {
					return wrongTypeError(stringClass)
				}

				rightValue := right.Value

				if leftValue != rightValue {
					return TRUE
				}

				return FALSE
			}
		},
	},
	{
		// Return a new String with the first character converted to uppercase but the rest of string converted to lowercase.
		//
		// ```Ruby
		// "test".capitalize         # => "Test"
		// "tEST".capitalize         # => "Test"
		// "heLlo\nWoRLd".capitalize # => "Hello\nworld"
		// ```
		//
		// @return [String]
		Name: "capitalize",
		Fn: func(receiver Object) builtinMethodBody {
			return func(t *thread, args []Object, blockFrame *callFrame) Object {

				str := []byte(receiver.(*StringObject).Value)
				start := string(str[0])
				rest := string(str[1:])
				result := strings.ToUpper(start) + strings.ToLower(rest)

				return initStringObject(result)
			}
		},
	},
	{
		// Returns a new String with all characters is upcase
		//
		// ```Ruby
		// "very big".upcase # => "VERY BIG"
		// ```
		//
		// @return [String]
		Name: "upcase",
		Fn: func(receiver Object) builtinMethodBody {
			return func(t *thread, args []Object, blockFrame *callFrame) Object {

				str := receiver.(*StringObject).Value

				return initStringObject(strings.ToUpper(str))
			}
		},
	},
	{
		// Returns a new String with all characters is lowercase
		//
		// ```Ruby
		// "erROR".downcase        # => "error"
		// "HeLlO\tWorLD".downcase # => "hello\tworld"
		// ```
		//
		// @return [String]
		Name: "downcase",
		Fn: func(receiver Object) builtinMethodBody {
			return func(t *thread, args []Object, blockFrame *callFrame) Object {

				str := receiver.(*StringObject).Value

				return initStringObject(strings.ToLower(str))
			}
		},
	},
	{
		// Returns the character length of self
		// **Note:** the length is currently byte-based, instead of charcode-based.
		//
		// ```Ruby
		// "zero".size # => 4
		// "".size # => 0
		// ```
		//
		// @return [Integer]
		Name: "size",
		Fn: func(receiver Object) builtinMethodBody {
			return func(t *thread, args []Object, blockFrame *callFrame) Object {

				str := receiver.(*StringObject).Value

				return initIntegerObject(len(str))
			}
		},
	},
	{
		// Returns the character length of self
		// **Note:** the length is currently byte-based, instead of charcode-based.
		//
		// ```Ruby
		// "zero".size # => 4
		// "".size # => 0
		// ```
		//
		// @return [Integer]
		Name: "length",
		Fn: func(receiver Object) builtinMethodBody {
			return func(t *thread, args []Object, blockFrame *callFrame) Object {

				str := receiver.(*StringObject).Value

				return initIntegerObject(len(str))
			}
		},
	},
	{
		// Returns a new String with reverse order of self
		// **Note:** the length is currently byte-based, instead of charcode-based.
		//
		// ```Ruby
		// "reverse".reverse      # => "esrever"
		// "Hello\nWorld".reverse # => "dlroW\nolleH"
		// ```
		//
		// @return [String]
		Name: "reverse",
		Fn: func(receiver Object) builtinMethodBody {
			return func(t *thread, args []Object, blockFrame *callFrame) Object {

				str := receiver.(*StringObject).Value

				var revert string
				for i := len(str) - 1; i >= 0; i-- {
					revert += string(str[i])
				}

				return initStringObject(revert)
			}
		},
	},
	{
		// Returns a new String with self value
		//
		// ```Ruby
		// "string".to_s # => "string"
		// ```
		//
		// @return [String]
		Name: "to_s",
		Fn: func(receiver Object) builtinMethodBody {
			return func(t *thread, args []Object, blockFrame *callFrame) Object {

				str := receiver.(*StringObject).Value

				return initStringObject(str)
			}
		},
	},
	{
		// Returns the result of converting self to Integer
		//
		// ```Ruby
		// "123".to_i # => 123
		// "3d print".to_i # => 3
		// "some text".to_i # => 0
		// ```
		//
		// @return [Integer]
		Name: "to_i",
		Fn: func(receiver Object) builtinMethodBody {
			return func(t *thread, args []Object, blockFrame *callFrame) Object {

				str := receiver.(*StringObject).Value
				parsedStr, err := strconv.ParseInt(str, 10, 0)

				if err == nil {
					return initIntegerObject(int(parsedStr))
				}

				var digits string
				for _, char := range str {
					if unicode.IsDigit(char) {
						digits += string(char)
					} else {
						break
					}
				}

				if len(digits) > 0 {
					parsedStr, _ = strconv.ParseInt(digits, 10, 0)
					return initIntegerObject(int(parsedStr))
				}

				return initIntegerObject(0)
			}
		},
	},
	{
		// Checks if the specified string is included in the receiver
		//
		// ```Ruby
		// "Hello\nWorld".include("\n") # => true
		// ```
		//
		// @return [Bool]
		Name: "include",
		Fn: func(receiver Object) builtinMethodBody {
			return func(t *thread, args []Object, blockFrame *callFrame) Object {

				rcv := receiver.(*StringObject).Value
				arg := args[0].(*StringObject).Value

				if strings.Contains(rcv, arg) {
					return TRUE
				}

				return FALSE
			}
		},
	},
	{
		// Returns a string which is concatenate with the input string or character
		//
		// ```ruby
		// "Hello ".concat("World") # => "Hello World"
		// ```
		// @return [String]
		Name: "concat",
		Fn: func(receiver Object) builtinMethodBody {
			return func(t *thread, args []Object, blockFrame *callFrame) Object {

				str := receiver.(*StringObject).Value
				concatStr, ok := args[0].(*StringObject)

				if !ok {
					return wrongTypeError(stringClass)
				}

				return initStringObject(str + concatStr.Value)
			}
		},
	},
	{
		// Returns the character of the string with specified index
		// It will raise error if the input is not an Integer type
		//
		// ```ruby
		// "Hello"[1]        # => "e"
		// "Hello"[5]        # => nil
		//
		// # TODO: Carriage Return Case
		// "Hello\nWorld"[5] # => "\n"
		// # TODO: Negative Index Case
		// "Hello"[-1]       # => "o"
		// ```
		// @return [String]
		Name: "[]",
		Fn: func(receiver Object) builtinMethodBody {
			return func(t *thread, args []Object, blockFrame *callFrame) Object {
				if len(args) != 1 {
					return &Error{Message: "Expect 1 arguments. got=%d" + string(len(args))}
				}

				str := receiver.(*StringObject).Value
				i := args[0]
				index, ok := i.(*IntegerObject)
				indexValue := index.Value

				if !ok {
					return newError("Expect index argument to be Integer. got=%T", i)
				}

				if len(str) > indexValue {
					return initStringObject(string([]rune(str)[indexValue]))
				}
				return NULL
			}
		},
	},
	{
		// Replace character of the string with input string
		// It will raise error if the index is not Integer type or the index value is out of
		// range of the string length
		//
		// ```ruby
		// "Ruby"[1] = "oo" # => "Rooby"
		// "Go"[2] = "by"   # => "Goby"
		// # TODO: Carriage Return Case
		// "Hello\nWorld"[5] = " " # => "Hello World"
		// # TODO: Negative Index Case
		// "Ruby"[-3] = "oo" # => "Rooby"
		// ```
		// @return [String]
		Name: "[]=",
		Fn: func(receiver Object) builtinMethodBody {
			return func(t *thread, args []Object, blockFrame *callFrame) Object {

				str := receiver.(*StringObject).Value
				i := args[0]
				index, ok := i.(*IntegerObject)
				indexValue := index.Value

				if !ok {
					return newError("Expect index argument to be Integer. got=%T", i)
				}

				if len(str) < indexValue {
					return newError("Index value out of range. got=%T", i)
				}

				replaceStr := args[1].(*StringObject).Value
				if len(str) == indexValue {
					return initStringObject(str + replaceStr)
				}
				result := str[:indexValue] + replaceStr + str[indexValue+1:]
				return initStringObject(result)
			}
		},
	},
	{
		// Returns an array of characters converted from a string
		// ```ruby
		// "Goby".to_a # => ["G", "o", "b", "y"]
		// ```
		// @return [String]
		Name: "to_a",
		Fn: func(receiver Object) builtinMethodBody {
			return func(t *thread, args []Object, blockFrame *callFrame) Object {

				str := receiver.(*StringObject)

				return str.toArray()
			}
		},
	},
	//{
	//	// Doc
	//	Name: "count",
	//	Fn: func(receiver Object) builtinMethodBody {
	//		return func(t *thread, args []Object, blockFrame *callFrame) Object {
	//
	//			str := receiver.(*StringObject).Value
	//
	//			return initStringObject(str)
	//		}
	//	},
	//},
	{
		// Returns true if string is empty value
		//
		// ```ruby
		// "".empty      # => true
		// "Hello".empty # => false
		// ```
		// @return [Boolean]
		Name: "empty",
		Fn: func(receiver Object) builtinMethodBody {
			return func(t *thread, args []Object, blockFrame *callFrame) Object {

				str := receiver.(*StringObject).Value

				if str == "" {
					return TRUE
				}
				return FALSE
			}
		},
	},
	{
		// Returns true if receiver string is equal to argument string
		//
		// ```ruby
		// "Hello".eql("Hello") # => true
		// "Hello".eql("World") # => false
		// ```
		// @return [Boolean]
		Name: "eql",
		Fn: func(receiver Object) builtinMethodBody {
			return func(t *thread, args []Object, blockFrame *callFrame) Object {

				str := receiver.(*StringObject).Value
				compareStr, ok := args[0].(*StringObject)

				if !ok {
					return wrongTypeError(stringClass)
				}

				if compareStr.Value == str {
					return TRUE
				}
				return FALSE
			}
		},
	},
	{
		// Returns true if receiver string start with the argument string
		//
		// ```ruby
		// "Hello".start_with("Hel") # => true
		// "Hello".start_with("hel") # => false
		// ```
		// @return [Boolean]
		Name: "start_with",
		Fn: func(receiver Object) builtinMethodBody {
			return func(t *thread, args []Object, blockFrame *callFrame) Object {

				str := receiver.(*StringObject).Value
				compareStr, ok := args[0].(*StringObject)

				if !ok {
					return wrongTypeError(stringClass)
				}

				index := len(compareStr.Value) - 1
				if compareStr.Value == str[:index] {
					return TRUE
				}
				return FALSE
			}
		},
	},
	{
		// Returns true if receiver string end with the argument string
		//
		// ```ruby
		// "Hello".end_with("llo") # => true
		// "Hello".end_with("ell") # => false
		// ```
		// @return [Boolean]
		Name: "end_with",
		Fn: func(receiver Object) builtinMethodBody {
			return func(t *thread, args []Object, blockFrame *callFrame) Object {

				str := receiver.(*StringObject).Value
				compareStr, ok := args[0].(*StringObject)

				if !ok {
					return wrongTypeError(stringClass)
				}

				index := len(compareStr.Value)
				if compareStr.Value == str[index:] {
					return TRUE
				}
				return FALSE
			}
		},
	},
	{
		// Insert a string input in specified index value of the receiver string
		//
		// It will raise error if index value is not an integer or index value is out
		// of receiver string's range
		//
		// It will also raise error if the input string value is not type string
		//
		// ```ruby
		// "Hello".insert(0, "X") # => "XHello"
		// "Hello".insert(2, "X") # => "HeXllo"
		// "Hello".insert(5, "X") # => "HelloX"
		// # TODO: Negative Index Case
		// "Hello".insert(-1, "X") # => "HelloX"
		// "Hello".insert(-3, "X") # => "HelXlo"
		// ```
		// @return [String]
		Name: "insert",
		Fn: func(receiver Object) builtinMethodBody {
			return func(t *thread, args []Object, blockFrame *callFrame) Object {

				str := receiver.(*StringObject).Value
				i := args[0]
				index, ok := i.(*IntegerObject)
				indexValue := index.Value

				if !ok {
					return newError("Expect index argument to be Integer. got=%T", i)
				}

				if len(str) < indexValue {
					return newError("Index value out of range. got=%T", i)
				}

				insertStr, ok := args[1].(*StringObject)

				if !ok {
					return wrongTypeError(stringClass)
				}

				return initStringObject(str[:indexValue] + insertStr.Value + str[indexValue:])
			}
		},
	},
	//{
	//	// Doc
	//	Name: "delete",
	//	Fn: func(receiver Object) builtinMethodBody {
	//		return func(t *thread, args []Object, blockFrame *callFrame) Object {
	//
	//			str := receiver.(*StringObject).Value
	//			deleteStr, ok := args[0].(*StringObject)
	//
	//			if !ok {
	//				return wrongTypeError(stringClass)
	//			}
	//
	//			return initStringObject(str)
	//		}
	//	},
	//},
	{
		// Returns a string with the last character chopped
		//
		// ```ruby
		// "Hello".chop # => "Hell"
		// # TODO: Carriage Return Case
		// "Hello World\n".chop => "Hello World"
		// ```
		Name: "chop",
		Fn: func(receiver Object) builtinMethodBody {
			return func(t *thread, args []Object, blockFrame *callFrame) Object {

				str := receiver.(*StringObject).Value

				return initStringObject(str[:len(str)-1])
			}
		},
	},
	{
		// If input integer is greater than the length of receiver string, returns a new String of
		// length integer with receiver string left justified and padded with default " "; otherwise,
		// returns receiver string.
		//
		// It will raise error if the input string length is not integer type
		//
		// ```ruby
		// "Hello".ljust(2) # => "Hello"
		// "Hello".ljust(7) # => "Hello  "
		// # TODO: Default PadString
		// "Hello".ljust(10, "xo") => "Helloxoxox"
		// ```
		Name: "ljust",
		Fn: func(receiver Object) builtinMethodBody {
			return func(t *thread, args []Object, blockFrame *callFrame) Object {

				str := receiver.(*StringObject).Value
				l := args[0]
				strLength, ok := l.(*IntegerObject)
				strLengthValue := strLength.Value

				if !ok {
					return newError("Expect index argument to be Integer. got=%T", l)
				}

				padString := " "
				if strLengthValue > len(str) {
					for i := len(str); i < strLengthValue; i += len(padString) {
						str += padString
					}
				}

				return initStringObject(str)
			}
		},
	},
	{
		// If input integer is greater than the length of receiver string, returns a new String of
		// length integer with receiver string right justified and padded with default " "; otherwise,
		// returns receiver string.
		//
		// It will raise error if the input string length is not integer type
		//
		// ```ruby
		// "Hello".rjust(2) # => "Hello"
		// "Hello".rjust(7) # => "  Hello"
		// # TODO: Default PadString
		// "Hello".ljust(10, "xo") => "xoxoxHello"
		// ```
		Name: "rjust",
		Fn: func(receiver Object) builtinMethodBody {
			return func(t *thread, args []Object, blockFrame *callFrame) Object {

				str := receiver.(*StringObject).Value
				l := args[0]
				strLength, ok := l.(*IntegerObject)
				strLengthValue := strLength.Value

				if !ok {
					return newError("Expect index argument to be Integer. got=%T", l)
				}

				padString := " "
				if strLengthValue > len(str) {
					for i := len(str); i < strLengthValue; i += len(padString) {
						str = padString + str
					}
				}

				return initStringObject(str)
			}
		},
	},
	//{
	//	// Doc
	//	Name: "strip",
	//	Fn: func(receiver Object) builtinMethodBody {
	//		return func(t *thread, args []Object, blockFrame *callFrame) Object {
	//
	//			str := receiver.(*StringObject).Value
	//
	//			return initStringObject(str)
	//		}
	//	},
	//},
	//{
	//	// Doc
	//	Name: "split",
	//	Fn: func(receiver Object) builtinMethodBody {
	//		return func(t *thread, args []Object, blockFrame *callFrame) Object {
	//
	//			str := receiver.(*StringObject).Value
	//
	//			return initStringObject(str)
	//		}
	//	},
	//},
	//{
	//	// Doc
	//	Name: "slice",
	//	Fn: func(receiver Object) builtinMethodBody {
	//		return func(t *thread, args []Object, blockFrame *callFrame) Object {
	//
	//			str := receiver.(*StringObject).Value
	//
	//			return initStringObject(str)
	//		}
	//	},
	//},
	//{
	//	// Doc
	//	Name: "replace",
	//	Fn: func(receiver Object) builtinMethodBody {
	//		return func(t *thread, args []Object, blockFrame *callFrame) Object {
	//
	//			str := receiver.(*StringObject).Value
	//
	//			return initStringObject(str)
	//		}
	//	},
	//},
	//{
	//	// Doc
	//	Name: "gsub",
	//	Fn: func(receiver Object) builtinMethodBody {
	//		return func(t *thread, args []Object, blockFrame *callFrame) Object {
	//
	//			str := receiver.(*StringObject).Value
	//
	//			return initStringObject(str)
	//		}
	//	},
	//}
}

func initStringClass() {
	bc := &BaseClass{Name: "String", Methods: newEnvironment(), ClassMethods: newEnvironment(), Class: classClass, pseudoSuperClass: objectClass, superClass: objectClass}
	sc := &RString{BaseClass: bc}
	sc.setBuiltInMethods(builtinStringInstanceMethods, false)
	stringClass = sc
}
