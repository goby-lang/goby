package vm

func newEnvironment() *environment {
	s := make(map[string]Object)
	return &environment{store: s}
}

type environment struct {
	store map[string]Object
}

func (e *environment) get(name string) (Object, bool) {
	obj, ok := e.store[name]
	return obj, ok
}

func (e *environment) set(name string, val Object) Object {
	e.store[name] = val
	return val
}
