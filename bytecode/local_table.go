package bytecode

type localTable struct {
	store map[string]int
	count int
	depth int
	upper *localTable
}

func (lt *localTable) get(v string) (int, bool) {
	i, ok := lt.store[v]

	return i, ok
}

func (lt *localTable) set(val string) int {
	c, ok := lt.store[val]

	if !ok {
		c = lt.count
		lt.store[val] = c
		lt.count++
		return c
	}

	return c
}

func (lt *localTable) setLCL(v string, d int) (index, depth int) {
	index, depth, ok := lt.getLCL(v, d)

	if !ok {
		index = lt.set(v)
		depth = lt.depth
		return index, depth
	}

	return index, depth
}

func (lt *localTable) getLCL(v string, d int) (index, depth int, ok bool) {
	index, ok = lt.get(v)

	if ok {
		return index, d - lt.depth, ok
	}

	if lt.upper != nil {
		index, depth, ok = lt.upper.getLCL(v, d)
		return
	}

	return -1, 0, false
}

func newLocalTable(depth int) *localTable {
	s := make(map[string]int)
	return &localTable{store: s, depth: depth}
}
