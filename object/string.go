package object

type StringClass struct {
	*BaseClass
}

type StringObject struct {
	Class *StringClass
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

