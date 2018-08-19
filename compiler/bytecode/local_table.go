package bytecode

type localTable struct {
	store map[string]uint8
	count uint8
	depth uint8
	upper *localTable
}

func (lt *localTable) get(v string) (uint8, bool) {
	i, ok := lt.store[v]

	return i, ok
}

func (lt *localTable) set(val string) uint8 {
	c, ok := lt.store[val]

	if !ok {
		c = lt.count
		lt.store[val] = c
		lt.count++
		return c
	}

	return c
}

func (lt *localTable) setLCL(v string, d uint8) (index, depth uint8) {
	index, depth, ok := lt.getLCL(v, d)

	if !ok {
		index = lt.set(v)
		depth = 0
		return index, depth
	}

	return index, depth
}

func (lt *localTable) getLCL(v string, d uint8) (index, depth uint8, ok bool) {
	index, ok = lt.get(v)

	if ok {
		return index, d - lt.depth, ok
	}

	if lt.upper != nil {
		index, depth, ok = lt.upper.getLCL(v, d)
		return
	}

	return 0, 0, false
}

func newLocalTable(depth uint8) *localTable {
	s := make(map[string]uint8)
	return &localTable{store: s, depth: depth}
}
