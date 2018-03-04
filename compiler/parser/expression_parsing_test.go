package parser

import (
	"github.com/goby-lang/goby/compiler/ast"
	"github.com/goby-lang/goby/compiler/lexer"
	"testing"
)

const (
	Ident = iota
	Const
	Ivar
)

func TestArgumentPairExpressionFail(t *testing.T) {
	tests := []struct {
		input string
		error string
	}{
		{`foo: "bar"`, `unexpected : Line: 0`},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)
		_, err := p.ParseProgram()
		if err.Message != tt.error {
			t.Fatal("Expected hash literal parsing error")
			t.Fatal("expect: ", tt.error)
			t.Fatal("actual: ", err.Message)
		}
	}
}

func TestArrayExpression(t *testing.T) {
	tests := []struct {
		input            string
		expectedElements []int
	}{
		{`[]`, []int{}},
		{`[1]`, []int{1}},
		{`[1,2,4,5]`, []int{1, 2, 4, 5}},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)
		program, err := p.ParseProgram()

		if err != nil {
			t.Fatal(err.Message)
		}

		arrayExp := program.FirstStmt().IsExpression(t).IsArrayExpression(t)

		for i, elem := range arrayExp.TestableElements() {
			elem.IsIntegerLiteral(t).ShouldEqualTo(tt.expectedElements[i])
		}
	}
}

func TestArrayIndexExpression(t *testing.T) {
	tests := []struct {
		input         string
		expectedIndex interface{}
	}{
		{`[][1]`, 1},
		{`[1][0]`, 0},
		{`[1,2,4,5][3]`, 3},
		{`[1,2,4,5][foo]`, "foo"},
		{`test[bar]`, "bar"},
		{`test[1]`, 1},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)
		program, err := p.ParseProgram()

		if err != nil {
			t.Fatal(err.Message)
		}

		arrIndexing := program.FirstStmt().IsExpression(t).IsCallExpression(t)

		switch expected := tt.expectedIndex.(type) {
		case int:
			arrIndexing.NthArgument(1).IsIntegerLiteral(t).ShouldEqualTo(expected)
		case string:
			arrIndexing.NthArgument(1).IsIdentifier(t).ShouldHasName(expected)
		}
	}
}

func TestArrayMultipleIndexExpression(t *testing.T) {
	tests := []struct {
		input         string
		expectedIndex []interface{}
	}{
		{`[][1, 2]`, []interface{}{1, 2}},
		{`[][5, 4, 3, 2, 1]`, []interface{}{5, 4, 3, 2, 1}},
		{`[][foo, bar, baz]`, []interface{}{"foo", "bar", "baz"}},
		{`[][foo, 123, baz]`, []interface{}{"foo", 123, "baz"}},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)
		program, err := p.ParseProgram()

		if err != nil {
			t.Fatal(err.Message)
		}

		arrIndexing := program.FirstStmt().IsExpression(t).IsCallExpression(t)

		for i, value := range tt.expectedIndex {
			arg := arrIndexing.NthArgument(i + 1)
			switch expected := value.(type) {
			case int:
				arg.IsIntegerLiteral(t).ShouldEqualTo(expected)
			case string:
				arg.IsIdentifier(t).ShouldHasName(expected)
			}
		}
	}
}

func TestAssignExpressionWithLiteralValue(t *testing.T) {
	tests := []struct {
		input              string
		expectedIdentifier string
		expectedValue      interface{}
		variableType       int
	}{
		{"x = 5;", "x", 5, Ident},
		{"y = true;", "y", true, Ident},

		{"Foo = '123'", "Foo", "123", Const},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)

		program, err := p.ParseProgram()

		if err != nil {
			t.Fatal(err.Message)
		}

		assignExp := program.FirstStmt().IsExpression(t).IsAssignExpression(t)

		switch tt.variableType {
		case Ident:
			assignExp.NthVariable(1).IsIdentifier(t).ShouldHasName(tt.expectedIdentifier)
		case Const:
			assignExp.NthVariable(1).IsConstant(t).ShouldHasName(tt.expectedIdentifier)
		}

		switch v := tt.expectedValue.(type) {
		case int:
			assignExp.TestableValue().IsIntegerLiteral(t).ShouldEqualTo(v)
		case string:
			assignExp.TestableValue().IsStringLiteral(t).ShouldEqualTo(v)
		case bool:
			assignExp.TestableValue().IsBooleanExpression(t).ShouldEqualTo(v)
		}
	}
}

func TestAssignExpressionWithVariableValue(t *testing.T) {
	tests := []struct {
		input              string
		expectedIdentifier string
		expectedValue      string
		variableType       int
		valueType          int
	}{
		{"x = y", "x", "y", Ident, Ident},
		{"@foo = y", "@foo", "y", Ivar, Ident},
		{"y = @foo", "y", "@foo", Ident, Ivar},
		{"Foo = @bar", "Foo", "@bar", Const, Ivar},
		{"@bar = Foo", "@bar", "Foo", Ivar, Const},
		{"@bar = @foo", "@bar", "@foo", Ivar, Ivar},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)

		program, err := p.ParseProgram()

		if err != nil {
			t.Fatal(err.Message)
		}

		assignExp := program.FirstStmt().IsExpression(t).IsAssignExpression(t)

		switch tt.variableType {
		case Ident:
			assignExp.NthVariable(1).IsIdentifier(t).ShouldHasName(tt.expectedIdentifier)
		case Const:
			assignExp.NthVariable(1).IsConstant(t).ShouldHasName(tt.expectedIdentifier)
		case Ivar:
			assignExp.NthVariable(1).IsInstanceVariable(t).ShouldHasName(tt.expectedIdentifier)
		}

		switch tt.valueType {
		case Ident:
			assignExp.TestableValue().IsIdentifier(t).ShouldHasName(tt.expectedValue)
		case Const:
			assignExp.TestableValue().IsConstant(t).ShouldHasName(tt.expectedValue)
		case Ivar:
			assignExp.TestableValue().IsInstanceVariable(t).ShouldHasName(tt.expectedValue)
		}
	}
}

func TestCallExpression(t *testing.T) {
	input := `
		p.add(1, 2 * 3, 4 + 5)
	`

	l := lexer.New(input)
	p := New(l)
	program, err := p.ParseProgram()

	if err != nil {
		t.Fatal(err.Message)
	}

	callExpression := program.FirstStmt().IsExpression(t).IsCallExpression(t)
	callExpression.TestableReceiver().IsIdentifier(t).ShouldHasName("p")
	callExpression.ShouldHasMethodName("add")
	callExpression.ShouldHasNumbersOfArguments(3)

	callExpression.NthArgument(1).IsIntegerLiteral(t).ShouldEqualTo(1)
	infix1 := callExpression.NthArgument(2).IsInfixExpression(t)
	infix1.ShouldHasOperator("*")
	infix1.TestableLeftExpression().IsIntegerLiteral(t).ShouldEqualTo(2)
	infix1.TestableRightExpression().IsIntegerLiteral(t).ShouldEqualTo(3)

	infix2 := callExpression.NthArgument(3).IsInfixExpression(t)
	infix2.ShouldHasOperator("+")
	infix2.TestableLeftExpression().IsIntegerLiteral(t).ShouldEqualTo(4)
	infix2.TestableRightExpression().IsIntegerLiteral(t).ShouldEqualTo(5)
}

func TestCallExpressionWithBlock(t *testing.T) {
	input := `
	[1, 2, 3, 4].each do |i|
	  puts(i)
	end
	`
	l := lexer.New(input)
	p := New(l)
	program, err := p.ParseProgram()

	if err != nil {
		t.Fatal(err.Message)
	}

	callExpression := program.FirstStmt().IsExpression(t).IsCallExpression(t)
	callExpression.TestableReceiver().IsArrayExpression(t)
	callExpression.ShouldHasMethodName("each")
	callExpression.BlockArguments[0].IsIdentifier(t).ShouldHasName("i")

	block := callExpression.Block
	exp := block.Statements[0].(ast.TestableStatement).IsExpression(t)
	exp.IsCallExpression(t).ShouldHasMethodName("puts")
}

func TestCaseExpression(t *testing.T) {
	input := `
	case 2
	when 0
	  0 + 0
	when 1
	  1 + 1
	else
	  2 + 2
	end
	`

	l := lexer.New(input)
	p := New(l)
	program, err := p.ParseProgram()

	if err != nil {
		t.Fatal(err.Message)
	}

	exp := program.FirstStmt().IsExpression(t).IsIfExpression(t)
	exp.ShouldHasNumberOfConditionals(2)
	cs := exp.TestableConditionals()

	c0 := cs[0].IsConditionalExpression(t)
	condition0 := c0.TestableCondition().IsInfixExpression(t)
	condition0.ShouldHasOperator("==")
	condition0.TestableLeftExpression().IsIntegerLiteral(t).ShouldEqualTo(2)
	condition0.TestableRightExpression().IsIntegerLiteral(t).ShouldEqualTo(0)

	consequence0 := c0.TestableConsequence()
	firstConsequenceExp := consequence0.NthStmt(1).IsExpression(t).IsInfixExpression(t)
	firstConsequenceExp.ShouldHasOperator("+")
	firstConsequenceExp.TestableLeftExpression().IsIntegerLiteral(t).ShouldEqualTo(0)
	firstConsequenceExp.TestableRightExpression().IsIntegerLiteral(t).ShouldEqualTo(0)

	c1 := cs[1].IsConditionalExpression(t)
	condition1 := c1.TestableCondition().IsInfixExpression(t)
	condition1.ShouldHasOperator("==")
	condition1.TestableLeftExpression().IsIntegerLiteral(t).ShouldEqualTo(2)
	condition1.TestableRightExpression().IsIntegerLiteral(t).ShouldEqualTo(1)

	consequence1 := c1.TestableConsequence()
	firstConsequenceExp = consequence1.NthStmt(1).IsExpression(t).IsInfixExpression(t)
	firstConsequenceExp.ShouldHasOperator("+")
	firstConsequenceExp.TestableLeftExpression().IsIntegerLiteral(t).ShouldEqualTo(1)
	firstConsequenceExp.TestableRightExpression().IsIntegerLiteral(t).ShouldEqualTo(1)

	alternative := exp.TestableAlternative()
	alternativeInfix := alternative.NthStmt(1).IsExpression(t).IsInfixExpression(t)
	alternativeInfix.ShouldHasOperator("+")
	alternativeInfix.TestableLeftExpression().IsIntegerLiteral(t).ShouldEqualTo(2)
	alternativeInfix.TestableRightExpression().IsIntegerLiteral(t).ShouldEqualTo(2)
}

func TestConstantExpression(t *testing.T) {
	input := `Person;`

	l := lexer.New(input)
	p := New(l)
	program, err := p.ParseProgram()

	if err != nil {
		t.Fatal(err.Message)
	}

	constant := program.FirstStmt().IsExpression(t).IsConstant(t)
	constant.ShouldHasName("Person")
}

func TestHashExpression(t *testing.T) {
	tests := []struct {
		input            string
		expectedElements map[string]int
	}{
		{`{}`, map[string]int{}},
		{
			`{ test: 123 }`,
			map[string]int{
				"test": 123,
			},
		},
		{

			`{ another_string: 456 }`,
			map[string]int{
				"another_string": 456,
			},
		},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)
		program, err := p.ParseProgram()

		if err != nil {
			t.Fatal(err.Message)
		}

		hash := program.FirstStmt().IsExpression(t).IsHashExpression(t)

		for key := range hash.TestableDataPairs() {
			hash.TestableDataPairs()[key].IsIntegerLiteral(t).ShouldEqualTo(tt.expectedElements[key])
		}
	}
}

func TestHashExpressionFail(t *testing.T) {
	tests := []struct {
		input string
		error string
	}{
		{`{ 1 }`, `could not parse "1" as hash key. Line: 0`},
		{`{ "a" }`, `could not parse "a" as hash key. Line: 0`},
		{`{ nil }`, `could not parse "nil" as hash key. Line: 0`},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)
		_, err := p.ParseProgram()

		if err.Message != tt.error {
			t.Fatal("Expected hash literal parsing error")
			t.Fatal("expect: ", tt.error)
			t.Fatal("actual: ", err.Message)
		}
	}
}

func TestHashAccessExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{`{}["123"]`, "123"},
		{`{ test: 1 }["test"]`, "test"},
		{`{ foo: "123" }["foo"]`, "foo"},
		{`{ bar: true }["bar"]`, "bar"},
		{`{ bar: true }[var]`, "var"},
	}

	for i, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)
		program, err := p.ParseProgram()

		if err != nil {
			t.Fatal(err.Message)
		}

		callExp := program.FirstStmt().IsExpression(t).IsCallExpression(t)
		callExp.TestableReceiver().IsHashExpression(t)
		callExp.ShouldHasNumbersOfArguments(1)

		if i < 4 {
			callExp.NthArgument(1).IsStringLiteral(t)
		} else {
			callExp.NthArgument(1).IsIdentifier(t).ShouldHasName("var")
		}

	}
}

func TestIdentifierExpression(t *testing.T) {
	input := `foobar;`

	l := lexer.New(input)
	p := New(l)
	program, err := p.ParseProgram()

	if err != nil {
		t.Fatal(err.Message)
	}

	ident := program.FirstStmt().IsExpression(t).IsIdentifier(t)
	ident.ShouldHasName("foobar")

}

func TestIfExpression(t *testing.T) {
	input := `
	if x < y
	  @x + 5
	elsif x == y
	  y + 5
	elsif x > y
	  y - 1
	else
	  y + 4
	end
	`

	l := lexer.New(input)
	p := New(l)
	program, err := p.ParseProgram()

	if err != nil {
		t.Fatal(err.Message)
	}

	exp := program.FirstStmt().IsExpression(t).IsIfExpression(t)
	exp.ShouldHasNumberOfConditionals(3)

	cs := exp.TestableConditionals()

	c0 := cs[0].IsConditionalExpression(t)
	condition0 := c0.TestableCondition().IsInfixExpression(t)
	condition0.ShouldHasOperator("<")
	condition0.TestableLeftExpression().IsIdentifier(t).ShouldHasName("x")
	condition0.TestableRightExpression().IsIdentifier(t).ShouldHasName("y")
	consequence0 := c0.TestableConsequence().NthStmt(1).IsExpression(t).IsInfixExpression(t)
	consequence0.ShouldHasOperator("+")
	consequence0.TestableLeftExpression().IsInstanceVariable(t).ShouldHasName("@x")
	consequence0.TestableRightExpression().IsIntegerLiteral(t).ShouldEqualTo(5)

	c1 := cs[1].IsConditionalExpression(t)
	condition1 := c1.TestableCondition().IsInfixExpression(t)
	condition1.ShouldHasOperator("==")
	condition1.TestableLeftExpression().IsIdentifier(t).ShouldHasName("x")
	condition1.TestableRightExpression().IsIdentifier(t).ShouldHasName("y")
	consequence1 := c1.TestableConsequence().NthStmt(1).IsExpression(t).IsInfixExpression(t)
	consequence1.ShouldHasOperator("+")
	consequence1.TestableLeftExpression().IsIdentifier(t).ShouldHasName("y")
	consequence1.TestableRightExpression().IsIntegerLiteral(t).ShouldEqualTo(5)

	c2 := cs[2].IsConditionalExpression(t)
	condition2 := c2.TestableCondition().IsInfixExpression(t)
	condition2.ShouldHasOperator(">")
	condition2.TestableLeftExpression().IsIdentifier(t).ShouldHasName("x")
	condition2.TestableRightExpression().IsIdentifier(t).ShouldHasName("y")
	consequence2 := c2.TestableConsequence().NthStmt(1).IsExpression(t).IsInfixExpression(t)
	consequence2.ShouldHasOperator("-")
	consequence2.TestableLeftExpression().IsIdentifier(t).ShouldHasName("y")
	consequence2.TestableRightExpression().IsIntegerLiteral(t).ShouldEqualTo(1)

	alternative := exp.TestableAlternative()
	alternativeExp := alternative.NthStmt(1).IsExpression(t).IsInfixExpression(t)
	alternativeExp.ShouldHasOperator("+")
	alternativeExp.TestableLeftExpression().IsIdentifier(t).ShouldHasName("y")
	alternativeExp.TestableRightExpression().IsIntegerLiteral(t).ShouldEqualTo(4)
}

func TestInfixExpression(t *testing.T) {
	infixTests := []struct {
		input      string
		leftValue  int
		operator   string
		rightValue int
	}{
		{"4 + 1;", 4, "+", 1},
		{"3 - 2;", 3, "-", 2},
	}

	for _, tt := range infixTests {
		l := lexer.New(tt.input)
		p := New(l)
		program, err := p.ParseProgram()

		if err != nil {
			t.Fatal(err.Message)
		}

		exp := program.FirstStmt().IsExpression(t).IsInfixExpression(t)
		exp.ShouldHasOperator(tt.operator)
		exp.TestableLeftExpression().IsIntegerLiteral(t).ShouldEqualTo(tt.leftValue)
		exp.TestableRightExpression().IsIntegerLiteral(t).ShouldEqualTo(tt.rightValue)
	}
}

func TestIntegerLiteralExpression(t *testing.T) {
	input := `5;`

	l := lexer.New(input)
	p := New(l)
	program, err := p.ParseProgram()

	if err != nil {
		t.Fatal(err.Message)
	}

	integerLiteral := program.FirstStmt().IsExpression(t).IsIntegerLiteral(t)
	integerLiteral.ShouldEqualTo(5)
}

func TestIntegerLiteralExpressionFail(t *testing.T) {
	input := `9223372036854775808;`

	l := lexer.New(input)
	p := New(l)
	_, err := p.ParseProgram()

	if err == nil {
		t.Fatal("Expected Integer literal parsing error")
	} else if p.error.Message != "could not parse \"9223372036854775808\" as integer. Line: 0" {
		t.Fatalf("Unexpected parsing error: %s", p.error.Message)
	}

	// "could not parse 9223372036854775808 as integer. Line: 1"
}

func TestNamespaceConstant(t *testing.T) {
	input := `
	Foo::Bar
	`

	l := lexer.New(input)
	p := New(l)
	program, err := p.ParseProgram()

	if err != nil {
		t.Fatal(err.Message)
	}

	infixExp := program.FirstStmt().IsExpression(t).IsInfixExpression(t)
	infixExp.ShouldHasOperator("::")
	infixExp.TestableLeftExpression().IsConstant(t).ShouldHasName("Foo")
	infixExp.TestableRightExpression().IsConstant(t).ShouldHasName("Bar")
}

func TestNilExpression(t *testing.T) {
	input := `
	nil
	`

	l := lexer.New(input)
	p := New(l)
	program, err := p.ParseProgram()

	if err != nil {
		t.Fatal(err.Message)
	}

	stmt := program.Statements[0].(*ast.ExpressionStatement)
	_, ok := stmt.Expression.(*ast.NilExpression)

	if !ok {
		t.Fatalf("Expect expression to be NilExpression. got=%T", stmt.Expression)
	}
}

func TestSelfExpression(t *testing.T) {
	input := `
		self.add(1, 2 * 3, 4 + 5);
	`

	l := lexer.New(input)
	p := New(l)
	program, err := p.ParseProgram()

	if err != nil {
		t.Fatal(err.Message)
	}

	callExpression := program.FirstStmt().IsExpression(t).IsCallExpression(t)
	callExpression.ShouldHasMethodName("add")
	callExpression.ShouldHasNumbersOfArguments(3)
	callExpression.TestableReceiver().IsSelfExpression(t)

	callExpression.NthArgument(1).IsIntegerLiteral(t).ShouldEqualTo(1)

	infix1 := callExpression.NthArgument(2).IsInfixExpression(t)
	infix1.ShouldHasOperator("*")
	infix1.TestableLeftExpression().IsIntegerLiteral(t).ShouldEqualTo(2)
	infix1.TestableRightExpression().IsIntegerLiteral(t).ShouldEqualTo(3)

	infix2 := callExpression.NthArgument(3).IsInfixExpression(t)
	infix2.ShouldHasOperator("+")
	infix2.TestableLeftExpression().IsIntegerLiteral(t).ShouldEqualTo(4)
	infix2.TestableRightExpression().IsIntegerLiteral(t).ShouldEqualTo(5)
}

func TestStringLiteralExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{input: `"testString";`, expected: "testString"},
		{input: `'test_string';`, expected: "test_string"},
		{input: `'!@#!@!$123';`, expected: "!@#!@!$123"},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)
		program, err := p.ParseProgram()

		if err != nil {
			t.Fatal(err.Message)
		}

		literal := program.FirstStmt().IsExpression(t).IsStringLiteral(t)
		literal.ShouldEqualTo(tt.expected)
	}
}

func TestArithmeticExpressionFail(t *testing.T) {
	tests := []struct {
		input string
		error string
	}{
		{`{ 1 ++ 1 }`, `unexpected + Line: 0`},
		{`{ 1 * * 1 }`, `unexpected * Line: 0`},
		{`{ 1 ** [1, 2] }`, `expected next token to be }, got **(**) instead. Line: 0`},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)
		_, err := p.ParseProgram()

		if err.Message != tt.error {
			t.Log("Expected arithmetic parsing error")
			t.Log("expect: ", tt.error)
			t.Fatal("actual: ", err.Message)
		}
	}
}
