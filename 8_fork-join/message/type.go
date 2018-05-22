package message

// Map is Map(int) signature.
type Map struct {
	V int
}

// GetOp returns the label.
func (*Map) GetOp() string { return "Map" }

// Red is Red(int) signature.
type Red struct {
	V int
}

// GetOp returns the label.
func (*Red) GetOp() string { return "Red" }
