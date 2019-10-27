package vm

import "testing"

func TestObjectClassSuperclass(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{`Object.class.name`, "Class"},
		{`Object.superclass.name`, "Object"},
	}

	for i, tt := range tests {
		v := initTestVM()
		evaluated := v.testEval(t, tt.input, getFilename())
		VerifyExpected(t, i, evaluated, tt.expected)
		v.checkCFP(t, i, 0)
		v.checkSP(t, i, 1)
	}
}

func TestObjectTapMethod(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{
			`
			a = 1
			a.tap do |int|
				int + 1
			end
`, 1},
		{
			`
			a = 1
			b = 2
			a.tap do |int|
				b = int + b
			end
			b
`, 3},
	}

	for i, tt := range tests {
		v := initTestVM()
		evaluated := v.testEval(t, tt.input, getFilename())
		VerifyExpected(t, i, evaluated, tt.expected)
		v.checkCFP(t, i, 0)
		v.checkSP(t, i, 1)
	}
}

func TestObjectTapMethodFail(t *testing.T) {
	testsFail := []errorTestCase{
		{`Object.new.tap`, "InternalError: Can't yield without a block", 1},
	}

	for i, tt := range testsFail {
		v := initTestVM()
		evaluated := v.testEval(t, tt.input, getFilename())
		checkErrorMsg(t, i, evaluated, tt.expected)
		v.checkCFP(t, i, tt.expectedCFP)
		v.checkSP(t, i, 1)
	}
}

func TestObjectDupMethod(t *testing.T) {
	setup := `
class Student
	attr_accessor :name
	def initialize(name)
		@name = name
	end
end

class School
	attr_accessor :name, :students

	def initialize(name, students)
		@name = name
		@students = students
	end

	def inspect
		String.fmt("Name: %s. Students: %s", name, students.map do |s| s.name end.join(", "))
	end
end
`
	tests := []struct {
		input    string
		expected interface{}
	}{
		{
			`
stan = Student.new("Stan")
stan.dup.name
`, "Stan"},
		{
			`
stan = Student.new("Stan")
dup = stan.dup
dup.name = "Jane"

[stan.name, dup.name]
`, []interface{}{"Stan", "Jane"}},
		{
			`
stan = Student.new("Stan")
jane = Student.new("Stan")

s1 = School.new("S1", [stan])
s2 = s1.dup
s2.name = "S2"
s2.inspect
`, "Name: S2. Students: Stan"},
	}

	for i, tt := range tests {
		v := initTestVM()
		evaluated := v.testEval(t, setup+tt.input, getFilename())
		VerifyExpected(t, i, evaluated, tt.expected)
		v.checkCFP(t, i, 0)
		v.checkSP(t, i, 1)
	}
}

func TestObjectId(t *testing.T) {
		tests := []struct {
		input    string
		expected interface{}
	}{
		{`1.object_id == 1.object_id`, false},
		{`"123".object_id == "123".object_id`, false},
		{`a = 10; a.object_id == a.object_id`, true},
		{
			`
class Student; end

stan = Student.new
jane = Student.new

stan.object_id == stan.object_id && jane.object_id != stan.object_id
`, true},
	}

	for i, tt := range tests {
		v := initTestVM()
		evaluated := v.testEval(t, tt.input, getFilename())
		VerifyExpected(t, i, evaluated, tt.expected)
		v.checkCFP(t, i, 0)
		v.checkSP(t, i, 1)
	}
}
