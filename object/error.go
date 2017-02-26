package object

type ErrorClass struct {
	*BaseClass
}

type Error struct {
	Class   *ErrorClass
	Message string
}

func (e *Error) Type() ObjectType {
	return ERROR_OBJ
}

func (e *Error) Inspect() string {
	return "ERROR: " + e.Message
}

func (e *Error) ReturnClass() Class {
	return e.Class
}
