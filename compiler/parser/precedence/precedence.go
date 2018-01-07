package precedence

// Constants for denoting precedence
const (
	_ int = iota
	LOWEST
	NORMAL
	ASSIGN
	LOGIC
	RANGE
	EQUALS
	COMPARE
	SUM
	PRODUCT
	PREFIX
	INDEX
	CALL
)
