package vm

func newEnvironment() *environment {
	s := make(map[string]Object)
	return &environment{store: s, outer: nil}
}

func closedEnvironment(outer *environment) *environment {
	env := newEnvironment()
	env.outer = outer
	return env
}

type environment struct {
	store map[string]Object
	outer *environment
}

type scope struct {
	Env  *environment
	Self Object
}

func (e *environment) getCurrent(name string) (Object, bool) {
	obj, ok := e.store[name]
	return obj, ok
}

func (e *environment) getValueLocation(name string) (*environment, bool) {
	env := e
	_, ok := e.store[name]
	if !ok && e.outer != nil {
		env, ok = e.outer.getValueLocation(name)
	}
	return env, ok
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
