package message

// Next is Next(int) signature.
type Next struct {
	V int
}

// GetOp returns the label.
func (*Next) GetOp() string { return "Next" }
