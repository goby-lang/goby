package object

var (
	StringClass *RString
)

type RString struct {
	*BaseClass
}

type StringObject struct {
	Class *RString
	Value string
}

func (s *StringObject) Type() ObjectType {
	return STRING_OBJ
}

func (s *StringObject) Inspect() string {
	return s.Value
}

func (s *StringObject) ReturnClass() Class {
	return s.Class
}

var (
	stringTable = make(map[string]*StringObject)
)

func InitializeString(value string) *StringObject {
	addr, ok := stringTable[value]

	if !ok {
		s := &StringObject{Value: value, Class: StringClass}
		stringTable[value] = s
		return s
	}

	return addr
}
