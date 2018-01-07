package precedence

// Constants for denoting precedence
const (
	_ = iota
	Lowest
	Normal
	Assign
	Logic
	Range
	Equals
	Compare
	Sum
	Product
	Prefix
	Index
	Call
)
