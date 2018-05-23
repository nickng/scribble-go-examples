package message

// Int is a wrapper to an int.
type Int struct {
	V int
}

// GetOp returns the label.
func (*Int) GetOp() string { return "Int" }
