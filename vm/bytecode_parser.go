package vm

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

type Parser struct {
	Line       int
	LabelCount int
	VM         *VM
}

func NewBytecodeParser() *Parser {
	return &Parser{}
}

func (p *Parser) Parse(bytecodes string) []*instructionSet {
	iss := []*instructionSet{}
	bytecodes = removeEmptyLine(strings.TrimSpace(bytecodes))
	bytecodesByLine := strings.Split(bytecodes, "\n")
	p.parseSection(iss, bytecodesByLine)

	return iss
}

func (p *Parser) parseSection(iss []*instructionSet, bytecodesByLine []string) {
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

func (p *Parser) parseLabel(is *instructionSet, line string) {
	line = strings.Trim(line, "<")
	line = strings.Trim(line, ">")
	p.VM.setLabel(is, line)
}

func (p *Parser) parseInstruction(is *instructionSet, line string) {
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
		panic(fmt.Sprintf("Unknown command: %s. Line: %d", act, ln))
	}

	is.Define(int(ln), action, params...)
}

func (p *Parser) parseParam(param string) interface{} {
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
