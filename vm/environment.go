package vm

func newEnvironment() *environment {
	s := make(map[string]Object)
	return &environment{store: s, outer: nil}
}

type environment struct {
	store map[string]Object
	outer *environment
}

func (e *environment) get(name string) (Object, bool) {
	obj, ok := e.store[name]
	if !ok && e.outer != nil {
		obj, ok = e.outer.get(name)
	}
	return obj, ok
}

func (e *environment) set(name string, val Object) Object {
	e.store[name] = val
	return val
}
