package vm

import (
	"fmt"
	"github.com/rooby-lang/rooby/bytecode"
	"github.com/rooby-lang/rooby/parser"
	"io/ioutil"
	"path"
	"regexp"
	"strconv"
	"strings"
)

// bytecodeParser is responsible for parsing bytecodes
type bytecodeParser struct {
	line       int
	labelTable map[labelType]map[string][]*instructionSet
	vm         *VM
	filename   filename
	blockTable map[string]*instructionSet
	program    *instructionSet
}

// newBytecodeParser initializes bytecodeParser and its label table then returns it
func newBytecodeParser(file filename) *bytecodeParser {
	p := &bytecodeParser{filename: file}
	p.blockTable = make(map[string]*instructionSet)
	p.labelTable = map[labelType]map[string][]*instructionSet{
		bytecode.LabelDef:      make(map[string][]*instructionSet),
		bytecode.LabelDefClass: make(map[string][]*instructionSet),
	}

	return p
}

// parseBytecode parses given bytecodes and transfer them into a sequence of instruction set.
func (p *bytecodeParser) parseBytecode(bytecodes string) []*instructionSet {
	iss := []*instructionSet{}
	bytecodes = removeEmptyLine(strings.TrimSpace(bytecodes))
	bytecodesByLine := strings.Split(bytecodes, "\n")
	p.parseSection(iss, bytecodesByLine)

	return iss
}

func (p *bytecodeParser) parseSection(iss []*instructionSet, bytecodesByLine []string) {
	is := &instructionSet{filename: p.filename}
	count := 0

	// First line is label
	p.parseLabel(is, bytecodesByLine[0])

	for _, text := range bytecodesByLine[1:] {
		count++
		l := strings.TrimSpace(text)
		if strings.HasPrefix(l, "<") {
			p.parseSection(iss, bytecodesByLine[count:])
			break
		} else {
			p.parseInstruction(is, l)
		}
	}

	iss = append(iss, is)
}

func (p *bytecodeParser) parseLabel(is *instructionSet, line string) {
	line = strings.Trim(line, "<")
	line = strings.Trim(line, ">")
	p.setLabel(is, line)
}

func (p *bytecodeParser) setLabel(is *instructionSet, name string) {
	var l *label
	var ln string
	var lt labelType

	if name == bytecode.Program {
		p.program = is
		return
	} else {
		ln = strings.Split(name, ":")[1]
		lt = labelType(strings.Split(name, ":")[0])
	}

	l = &label{name: name, Type: lt}
	is.label = l

	if lt == bytecode.Block {
		p.blockTable[ln] = is
		return
	}

	p.labelTable[lt][ln] = append(p.labelTable[lt][ln], is)
}

// parseInstruction transfer a line of bytecode into an instruction and append it into given instruction set.
func (p *bytecodeParser) parseInstruction(is *instructionSet, line string) {
	var params []interface{}
	var rawParams []string

	tokens := strings.Split(line, " ")
	lineNum, act := tokens[0], tokens[1]
	ln, _ := strconv.ParseInt(lineNum, 0, 64)
	action := builtInActions[operationType(act)]

	if act == bytecode.PutString {
		text := strings.Split(line, "\"")[1]
		params = append(params, text)
	} else if act == bytecode.RequireRelative {
		filepath := tokens[2]
		filepath = path.Join(p.vm.fileDir, filepath)

		file, err := ioutil.ReadFile(filepath + ".ro")

		if err != nil {
			panic(err)
		}

		program := parser.BuildAST(file)
		g := bytecode.NewGenerator(program)
		bytecodes := g.GenerateByteCode(program)
		p.vm.ExecBytecodes(bytecodes, filepath)
		return
	} else if len(tokens) > 2 {
		rawParams = tokens[2:]

		for _, param := range rawParams {
			params = append(params, p.parseParam(param))
		}
	} else if action == nil {
		panic(fmt.Sprintf("Unknown command: %s. line: %d", act, ln))
	}

	is.define(int(ln), action, params...)
}

func (p *bytecodeParser) parseParam(param string) interface{} {
	integer, e := strconv.ParseInt(param, 0, 64)
	if e != nil {
		return param
	}

	i := int(integer)

	return i
}

func removeEmptyLine(s string) string {
	regex, err := regexp.Compile("\n+")
	if err != nil {
		panic(err)
	}
	s = regex.ReplaceAllString(s, "\n")

	return s
}
