package parser

import (
	"github.com/goby-lang/goby/compiler/ast"
	"github.com/goby-lang/goby/compiler/lexer"
	"testing"
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
			elem.IsIntegerLiteral(t).ShouldEqualTo(t, tt.expectedElements[i])
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
			arrIndexing.NthArgument(1).IsIntegerLiteral(t).ShouldEqualTo(t, expected)
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
				arg.IsIntegerLiteral(t).ShouldEqualTo(t, expected)
			case string:
				arg.IsIdentifier(t).ShouldHasName(expected)
			}
		}
	}
}

func TestAssignInfixExpressionWithLiteralValue(t *testing.T) {
	tests := []struct {
		input              string
		expectedIdentifier string
		expectedValue      interface{}
		variableMatchFunc  func(*testing.T, ast.Expression, string) bool
	}{
		{"x = 5;", "x", 5, testIdentifier},
		{"y = true;", "y", true, testIdentifier},

		{"Foo = '123'", "Foo", "123", testConstant},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)

		program, err := p.ParseProgram()

		if err != nil {
			t.Fatal(err.Message)
		}

		if program == nil {
			t.Fatal("ParseProgram() returned nil")
		}

		testAssignExpression(t, program.FirstStmt().IsExpression(t).IsAssignExpression(t), tt.expectedIdentifier, tt.variableMatchFunc, tt.expectedValue)
	}
}

func TestAssignIndexExpressionWithVariableValue(t *testing.T) {
	tests := []struct {
		input              string
		expectedIdentifier string
		expectedValue      string
		variableMatchFunc  func(*testing.T, ast.Expression, string) bool
		valueMatchFunc     func(*testing.T, ast.Expression, string) bool
	}{
		{"x = y", "x", "y", testIdentifier, testIdentifier},
		{"@foo = y", "@foo", "y", testInstanceVariable, testIdentifier},
		{"y = @foo", "y", "@foo", testIdentifier, testInstanceVariable},
		{"Foo = @bar", "Foo", "@bar", testConstant, testInstanceVariable},
		{"@bar = Foo", "@bar", "Foo", testInstanceVariable, testConstant},
		{"@bar = @foo", "@bar", "@foo", testInstanceVariable, testInstanceVariable},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)

		program, err := p.ParseProgram()

		if err != nil {
			t.Fatal(err.Message)
		}

		assignExp := program.FirstStmt().IsExpression(t).IsAssignExpression(t)
		tt.variableMatchFunc(t, assignExp.Variables[0], tt.expectedIdentifier)
		tt.valueMatchFunc(t, assignExp.Value, tt.expectedValue)
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

	if len(callExpression.Arguments) != 3 {
		t.Fatalf("expect %d arguments. got=%d", 3, len(callExpression.Arguments))
	}

	callExpression.NthArgument(1).IsIntegerLiteral(t).ShouldEqualTo(t, 1)
	testInfixExpression(t, callExpression.NthArgument(2).IsInfixExpression(t), 2, "*", 3)
	testInfixExpression(t, callExpression.NthArgument(3).IsInfixExpression(t), 4, "+", 5)
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
	exp := block.Statements[0].(ast.TestingStatement).IsExpression(t)
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

	if len(program.Statements) != 1 {
		t.Fatalf("expect program's statements to be 1. got=%d", len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)

	if !ok {
		t.Fatalf("expect program.Statements[0] to be *ast.ExpressionStatement. got=%T", program.Statements[0])
	}

	exp, ok := stmt.Expression.(*ast.IfExpression)

	if !ok {
		t.Fatalf("expect statement to be an CaseExpression. got=%T", stmt.Expression)
	}

	cs := exp.Conditionals

	if len(cs) != 2 {
		t.Fatalf("expect the length of conditionals to be 2. got=%d", len(cs))
	}

	c0 := cs[0]

	if !testInfixExpression(t, c0.Condition, 2, "==", 0) {
		return
	}

	if len(c0.Consequence.Statements) != 1 {
		t.Errorf("should be only one consequence. got=%d\n", len(c0.Consequence.Statements))
	}

	consequence0, ok := c0.Consequence.Statements[0].(*ast.ExpressionStatement)

	if !ok {
		t.Errorf("expect consequence should be an expression statement. got=%T", c0.Consequence.Statements[0])
	}

	if !testInfixExpression(t, consequence0.Expression, 0, "+", 0) {
		return
	}

	c1 := cs[1]

	if !testInfixExpression(t, c1.Condition, 2, "==", 1) {
		return
	}

	if len(c1.Consequence.Statements) != 1 {
		t.Errorf("should be only one consequence. got=%d\n", len(c1.Consequence.Statements))
	}

	consequence1, ok := c1.Consequence.Statements[0].(*ast.ExpressionStatement)

	if !ok {
		t.Errorf("expect consequence should be an expression statement. got=%T", c1.Consequence.Statements[0])
	}

	if !testInfixExpression(t, consequence1.Expression, 1, "+", 1) {
		return
	}

	alternative, ok := exp.Alternative.Statements[0].(*ast.ExpressionStatement)

	if !ok {
		t.Errorf("expect alternative should be an expression statement. got=%T", exp.Alternative.Statements[0])
	}

	if !testInfixExpression(t, alternative.Expression, 2, "+", 2) {
		return
	}
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

		for key := range hash.Data {
			testIntegerLiteral(t, hash.Data[key], tt.expectedElements[key])
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
		callExp.NthArgument(1)

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
	  x + 5
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
	exp.ShouldHasNumberOfConditionals(t, 3)

	cs := exp.TestableConditionals()

	c0 := cs[0].IsConditionalExpression(t)
	testInfixExpression(t, c0.Condition, "x", "<", "y")
	consequence0 := c0.TestableConsequence()[0].IsExpression(t).IsInfixExpression(t)
	testInfixExpression(t, consequence0, "x", "+", 5)

	c1 := cs[1].IsConditionalExpression(t)
	testInfixExpression(t, c1.Condition, "x", "==", "y")
	consequence1 := c1.TestableConsequence()[0].IsExpression(t).IsInfixExpression(t)
	testInfixExpression(t, consequence1, "y", "+", 5)

	c2 := cs[2].IsConditionalExpression(t)
	testInfixExpression(t, c2.Condition, "x", ">", "y")
	consequence2 := c2.TestableConsequence()[0].IsExpression(t).IsInfixExpression(t)
	testInfixExpression(t, consequence2, "y", "-", 1)

	alternative := exp.TestableAlternative()
	testInfixExpression(t, alternative[0].IsExpression(t), "y", "+", 4)
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

		if len(program.Statements) != 1 {
			t.Fatalf("expect %d statements. got=%d", 1, len(program.Statements))
		}

		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("program.Statments[0] is not ast.ExpressionStatement. got=%T", program.Statements[0])
		}
		testInfixExpression(t, stmt.Expression, tt.leftValue, tt.operator, tt.rightValue)
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
	integerLiteral.ShouldEqualTo(t, 5)
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

func TestMethodArgumentsParsing(t *testing.T) {
	tests := []struct {
		input          string
		expectedParams []string
	}{
		{input: "def add(x, y); end", expectedParams: []string{"x", "y"}},
		{input: `
		def print(x)
		end
		`, expectedParams: []string{"x"}},
		{input: "def test(x, y, z); end", expectedParams: []string{"x", "y", "z"}},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)
		program, err := p.ParseProgram()

		if err != nil {
			t.Fatal(err.Message)
		}
		methodStatement := program.Statements[0].(*ast.DefStatement)

		if len(methodStatement.Parameters) != len(tt.expectedParams) {
			t.Errorf("expect %d parameters. got=%d", len(tt.expectedParams), len(methodStatement.Parameters))
		}

		for i, expectedParam := range tt.expectedParams {
			testIdentifier(t, methodStatement.Parameters[i], expectedParam)
		}
	}
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

	stmt := program.Statements[0].(*ast.ExpressionStatement)
	callExpression := stmt.Expression.(*ast.CallExpression)

	self, ok := callExpression.Receiver.(*ast.SelfExpression)
	if !ok {
		t.Fatalf("expect receiver to be SelfExpression. got=%T", callExpression.Receiver)
	}

	if self.TokenLiteral() != "self" {
		t.Fatalf("expect SelfExpression's token literal to be 'self'. got=%s", self.TokenLiteral())
	}

	testMethodName(t, callExpression, "add")

	if len(callExpression.Arguments) != 3 {
		t.Fatalf("expect %d arguments. got=%d", 3, len(callExpression.Arguments))
	}

	testIntegerLiteral(t, callExpression.Arguments[0], 1)
	testInfixExpression(t, callExpression.Arguments[1], 2, "*", 3)
	testInfixExpression(t, callExpression.Arguments[2], 4, "+", 5)
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

		if len(program.Statements) != 1 {
			t.Fatalf("program has wrong number of statements. got=%d", len(program.Statements))
		}

		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("first program statement is not ast.ExpressionStatement. got=%T", program.Statements[0])
		}

		literal, ok := stmt.Expression.(*ast.StringLiteral)
		testStringLiteral(t, literal, tt.expected)
	}
}

func testAssignExpression(t *testing.T, exp ast.Expression, expectedIdentifier string, variableMatchFunction func(*testing.T, ast.Expression, string) bool, expected interface{}) {
	assignExp := exp.(*ast.TestableAssignExpression)

	if !variableMatchFunction(t, assignExp.Variables[0], expectedIdentifier) {
		return
	}

	switch expected := expected.(type) {
	case int:
		testIntegerLiteral(t, assignExp.Value, expected)
	case string:
		testStringLiteral(t, assignExp.Value, expected)
	case bool:
		testBoolLiteral(t, assignExp.Value, expected)
	}
}
