package ripper

import (
	"fmt"
	"strings"

	"github.com/goby-lang/goby/compiler"
	"github.com/goby-lang/goby/compiler/bytecode"
	"github.com/goby-lang/goby/compiler/lexer"
	"github.com/goby-lang/goby/compiler/parser"
	"github.com/goby-lang/goby/compiler/token"
	"github.com/goby-lang/goby/vm"
	"github.com/goby-lang/goby/vm/classes"
	"github.com/goby-lang/goby/vm/errors"
)

// Ripper is a loadable library and has abilities to parse/lex/tokenize/get instructions of Goby codes from String.
// The library would be convenient for validating Goby codes when building lint tools,
// as well as the tests for Goby's compiler.
// For now, Ripper is a class and has only class methods, but I think this should finally be a 'newable' module with more sophisticated instance methods.

// Object is an imported object from vm
type Object = vm.Object

// VM is an imported object from vm
type VM = vm.VM

// Thread is an imported object from vm
type Thread = vm.Thread

// Method is an imported object from vm
type Method = vm.Method

// StringObject is an imported object from vm
type StringObject = vm.StringObject

// HashObject is an imported object from vm
type HashObject = vm.HashObject

// ArrayObject is an imported object from vm
type ArrayObject = vm.ArrayObject

// Class methods --------------------------------------------------------

// Returns the list of instruction code generated from Goby code.
// Returns `[]` when the Goby code is invalid.
// The return value is a "tuple" style nested array:
// - `Array`: contains an instruction set
//   - `arg_types:` (none if `nil`)
// 		 - `names:` array of names (string)
//     - `types:` array of types (integer)
//   - `instructions:` array of instructions
//     - `action:` string
//     - `line:` integer
//     - `params:` array of parameters (string)
//     - `source_line:` integer
//     - `arg_set:` (none if `nil`)
//   		 - `names:` array of names (string)
//       - `types:` array of types (integer)
//
// ```ruby
// require 'ripper'; Ripper.instruction "10.times do |i| puts i end"
// #=> [{ arg_set: { names: ["i"], types: [0] }, instructions: [{ action: "putself", line: 0, params: [], source_line: 1 }, { action: "getlocal", line: 1, params: ["0", "0"], source_line: 1 }, { action: "send", line: 2, params: ["puts", "1", "", "&{[i] [0]}"], source_line: 1 }, { action: "leave", line: 3, params: [], source_line: 1 }], name: "0", type: "Block" }, { arg_set: { names: ["i"], types: [0] }, instructions: [{ action: "putself", line: 0, params: [], source_line: 1 }, { action: "getlocal", line: 1, params: ["0", "0"], source_line: 1 }, { action: "send", line: 2, params: ["puts", "1", "", "&{[i] [0]}"], source_line: 1 }, { action: "leave", line: 3, params: [], source_line: 1 }], name: "0", type: "Block" }, { arg_set: { names: [], types: [] }, instructions: [{ action: "putobject", line: 0, params: ["10"], source_line: 1 }, { action: "send", line: 1, params: ["times", "0", "block:0", "&{[] []}"], source_line: 1 }, { action: "pop", line: 2, params: [], source_line: 1 }, { action: "leave", line: 3, params: [], source_line: 1 }], name: "ProgramStart", type: "ProgramStart" }, { arg_set: { names: [], types: [] }, instructions: [{ action: "putobject", line: 0, params: ["10"], source_line: 1 }, { action: "send", line: 1, params: ["times", "0", "block:0", "&{[] []}"], source_line: 1 }, { action: "pop", line: 2, params: [], source_line: 1 }, { action: "leave", line: 3, params: [], source_line: 1 }], name: "ProgramStart", type: "ProgramStart" }]
//
// require 'ripper'; Ripper.instruction "10.times do |i| puts i" # the code is invalid
// #=> InternalError: invalid code: 10.times do |i| puts i
// ```
//
// @param Goby code [String]
// @return [Array]
func instruction(receiver Object, sourceLine int, t *Thread, args []Object) Object {
	if len(args) != 1 {
		return t.VM().InitErrorObject(errors.ArgumentError, sourceLine, "Expect 1 argument. got=%d", len(args))
	}

	arg, ok := args[0].(*StringObject)
	if !ok {
		return t.VM().InitErrorObject(errors.TypeError, sourceLine, errors.WrongArgumentTypeFormat, classes.StringClass, args[0].Class().Name)
	}

	i, err := compiler.CompileToInstructions(arg.Value().(string), parser.NormalMode)
	if err != nil {
		return t.VM().InitErrorObject(errors.InternalError, sourceLine, errors.InvalidCode, arg.ToString())
	}

	return convertToTuple(i, t.VM())
}

// Returns a nested array that contains the line #, type of the tokenize, and the literal of the tokenize.
// Note that the class method does not return any errors even though the provided Goby code is invalid.
//
// ```ruby
// require 'ripper'; Ripper.lex "10.times do |i| puts i end"
// #=> [[0, "on_int", "10"], [0, "on_dot", "."], [0, "on_ident", "times"], [0, "on_do", "do"], [0, "on_bar", "|"], [0, "on_ident", "i"], [0, "on_bar", "|"], [0, "on_ident", "puts"], [0, "on_ident", "i"], [0, "on_end", "end"], [0, "on_eof", ""]]
//
// require 'ripper'; Ripper.lex "10.times do |i| puts i" # the code is invalid
// #=> [[0, "on_int", "10"], [0, "on_dot", "."], [0, "on_ident", "times"], [0, "on_do", "do"], [0, "on_bar", "|"], [0, "on_ident", "i"], [0, "on_bar", "|"], [0, "on_ident", "puts"], [0, "on_ident", "i"], [0, "on_eof", ""]]
// ```
//
// @param Goby code [String]
// @return [Array]
func lex(receiver Object, sourceLine int, t *Thread, args []Object) Object {
	if len(args) != 1 {
		return t.VM().InitErrorObject(errors.ArgumentError, sourceLine, "Expect 1 argument. got=%d", len(args))
	}

	arg, ok := args[0].(*StringObject)
	if !ok {
		return t.VM().InitErrorObject(errors.TypeError, sourceLine, errors.WrongArgumentTypeFormat, classes.StringClass, args[0].Class().Name)
	}

	l := lexer.New(arg.Value().(string))
	array := t.VM().InitArrayObject([]Object{})
	var elements []Object
	var nextToken token.Token
	for i := 0; ; i++ {
		nextToken = l.NextToken()
		elements = append(elements, t.VM().InitIntegerObject(nextToken.Line))
		elements = append(elements, t.VM().InitStringObject(convertLex(nextToken.Type)))
		elements = append(elements, t.VM().InitStringObject(nextToken.Literal))
		array.Elements = append(array.Elements, t.VM().InitArrayObject(elements))
		if nextToken.Type == token.EOF {
			break
		}
		elements = nil
	}
	return array
}

// Just to disable creating instances.
func new(receiver Object, sourceLine int, t *Thread, args []Object) Object {
	return t.VM().InitNoMethodError(sourceLine, "new", receiver)
}

// Returns the parsed Goby codes as a String object.
// Returns an error when the code is invalid.
//
// ```ruby
// require 'ripper'; Ripper.parse "10.times do |i| puts i end"
// #=> "10.times() do |i|
// #=> self.puts(i)
// #=> end"
//
// require 'ripper'; Ripper.parse "10.times do |i| puts i" # the code is invalid
// #=> TypeError: InternalError%!(EXTRA string=String, string=Invalid Goby code)
// ```
//
// @param Goby code [String]
// @return [String]
func parse(receiver Object, sourceLine int, t *Thread, args []Object) Object {
	if len(args) != 1 {
		return t.VM().InitErrorObject(errors.ArgumentError, sourceLine, "Expect 1 argument. got=%d", len(args))
	}

	arg, ok := args[0].(*StringObject)
	if !ok {
		return t.VM().InitErrorObject(errors.TypeError, sourceLine, errors.WrongArgumentTypeFormat, classes.StringClass, args[0].Class().Name)
	}

	l := lexer.New(arg.Value().(string))
	p := parser.New(l)
	program, _ := p.ParseProgram()

	return t.VM().InitStringObject(program.String())
}

// Returns a tokenized Goby codes as an Array object.
// Note that this does not return any errors even though the provided code is invalid.
//
// ```ruby
// require 'ripper'; Ripper.tokenize "10.times do |i| puts i end"
// #=> ["10", ".", "times", "do", "|", "i", "|", "puts", "i", "end", "EOF"]
//
// require 'ripper'; Ripper.parse "10.times do |i| puts i" # the code is invalid
// #=> ["10", ".", "times", "do", "|", "i", "|", "puts", "i", "EOF"]
// ```
//
// @param Goby code [String]
// @return [String]
func tokenize(receiver Object, sourceLine int, t *Thread, args []Object) Object {
	if len(args) != 1 {
		return t.VM().InitErrorObject(errors.ArgumentError, sourceLine, "Expect 1 argument. got=%d", len(args))
	}

	arg, ok := args[0].(*StringObject)
	if !ok {
		return t.VM().InitErrorObject(errors.TypeError, sourceLine, errors.WrongArgumentTypeFormat, classes.StringClass, args[0].Class().Name)
	}

	l := lexer.New(arg.Value().(string))
	el := []Object{}
	var nt token.Token
	for i := 0; ; i++ {
		nt = l.NextToken()
		if nt.Type == token.EOF {
			el = append(el, t.VM().InitStringObject("EOF"))
			break
		}
		el = append(el, t.VM().InitStringObject(nt.Literal))
	}
	return t.VM().InitArrayObject(el)
}

// Internal functions ===================================================
func init() {
	vm.RegisterExternalClass("ripper", vm.ExternalClass("Ripper", "ripper.gb",
		// class methods
		map[string]vm.Method{
			"instruction": instruction,
			"lex":         lex,
			"new":         new,
			"parse":       parse,
			"tokenize":    tokenize,
		},
		// instance methods
		map[string]vm.Method{},
	))
}

// Other helper functions ----------------------------------------------

func convertToTuple(instSet []*bytecode.InstructionSet, v *VM) *ArrayObject {
	ary := []Object{}
	for _, instruction := range instSet {
		hashInstLevel1 := make(map[string]Object)
		hashInstLevel1["name"] = v.InitStringObject(instruction.Name())
		hashInstLevel1["type"] = v.InitStringObject(instruction.Type())
		if instruction.ArgTypes() != nil {
			hashInstLevel1["arg_types"] = getArgNameType(instruction.ArgTypes(), v)
		}
		ary = append(ary, v.InitHashObject(hashInstLevel1))

		arrayInst := []Object{}
		for _, ins := range instruction.Instructions {
			hashInstLevel2 := make(map[string]Object)
			hashInstLevel2["action"] = v.InitStringObject(ins.ActionName())
			hashInstLevel2["line"] = v.InitIntegerObject(ins.Line())
			hashInstLevel2["source_line"] = v.InitIntegerObject(ins.SourceLine())

			arrayParams := []Object{}
			for _, param := range ins.Params {
				arrayParams = append(arrayParams, v.InitStringObject(covertTypesToString(param)))
			}
			hashInstLevel2["params"] = v.InitArrayObject(arrayParams)

			if ins.Opcode == bytecode.Send {
				hashInstLevel1["arg_set"] = getArgNameType(ins.Params[3].(*bytecode.ArgSet), v)
			}

			arrayInst = append(arrayInst, v.InitHashObject(hashInstLevel2))
		}

		hashInstLevel1["instructions"] = v.InitArrayObject(arrayInst)
		ary = append(ary, v.InitHashObject(hashInstLevel1))
	}
	return v.InitArrayObject(ary)
}

func getArgNameType(argSet *bytecode.ArgSet, v *VM) *HashObject {
	h := make(map[string]Object)

	aName := []Object{}
	for _, argname := range argSet.Names() {
		aName = append(aName, v.InitStringObject(argname))
	}
	h["names"] = v.InitArrayObject(aName)

	aType := []Object{}
	for _, argtype := range argSet.Types() {
		aType = append(aType, v.InitIntegerObject(int(argtype)))
	}

	h["types"] = v.InitArrayObject(aType)
	return v.InitHashObject(h)
}

// TODO: This should finally be auto-generated from tokenize.go
func convertLex(t token.Type) string {
	var s string
	switch t {
	case token.Asterisk:
		s = "asterisk"
	case token.And:
		s = "and"
	case token.Assign:
		s = "assign"
	case token.Bang:
		s = "bang"
	case token.Bar:
		s = "bar"
	case token.Colon:
		s = "colon"
	case token.Comma:
		s = "comma"
	case token.COMP:
		s = "comp"
	case token.Dot:
		s = "dot"
	case token.Eq:
		s = "eq"
	case token.GT:
		s = "gt"
	case token.GTE:
		s = "gte"
	case token.LBrace:
		s = "lbrace"
	case token.LBracket:
		s = "lbracket"
	case token.LParen:
		s = "lparen"
	case token.LT:
		s = "lt"
	case token.LTE:
		s = "lte"
	case token.Minus:
		s = "minus"
	case token.MinusEq:
		s = "minuseq"
	case token.Modulo:
		s = "modulo"
	case token.NotEq:
		s = "noteq"
	case token.Or:
		s = "or"
	case token.OrEq:
		s = "oreq"
	case token.Plus:
		s = "plus"
	case token.PlusEq:
		s = "pluseq"
	case token.Pow:
		s = "pow"
	case token.Range:
		s = "range"
	case token.RBrace:
		s = "rbrace"
	case token.RBracket:
		s = "rbracket"
	case token.ResolutionOperator:
		s = "resolutionoperator"
	case token.RParen:
		s = "rparen"
	case token.Semicolon:
		s = "semicolon"
	case token.Slash:
		s = "slash"
	default:
		s = strings.ToLower(string(t))
	}

	return "on_" + s
}

func covertTypesToString(v interface{}) string {
	return fmt.Sprint(v)
}
