package vm

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

// bytecodeParser is responsible for parsing bytecodes
type bytecodeParser struct {
	line int
	VM   *VM
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
	is := &instructionSet{}
	count := 0

	// First line is label
	p.parseLabel(is, bytecodesByLine[0])

	for _, text := range bytecodesByLine[1:] {
		count += 1
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
	p.VM.setLabel(is, line)
}

// parseInstruction transfer a line of bytecode into an instruction and append it into given instruction set.
func (p *bytecodeParser) parseInstruction(is *instructionSet, line string) {
	var params []interface{}
	var rawParams []string

	tokens := strings.Split(line, " ")
	lineNum, act := tokens[0], tokens[1]
	ln, _ := strconv.ParseInt(lineNum, 0, 64)
	action := builtInActions[operationType(act)]

	if act == "putstring" {
		text := strings.Split(line, "\"")[1]
		params = append(params, text)
	} else if len(tokens) > 2 {
		rawParams = tokens[2:]

		for _, param := range rawParams {
			params = append(params, p.parseParam(param))
		}
	} else if action == nil {
		panic(fmt.Sprintf("Unknown command: %s. line: %d", act, ln))
	}

	is.Define(int(ln), action, params...)
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
