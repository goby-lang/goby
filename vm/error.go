package vm

type ErrorClass struct {
	*BaseClass
}

type Error struct {
	Class   *ErrorClass
	Message string
}

func (e *Error) Type() objectType {
	return errorObj
}

func (e *Error) Inspect() string {
	return "ERROR: " + e.Message
}

func (e *Error) returnClass() Class {
	return e.Class
}
