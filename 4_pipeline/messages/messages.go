package message

// Val is Val(int) signature.
type Val struct {
	V int
}

// GetOp returns the label.
func (*Val) GetOp() string { return "Val" }
