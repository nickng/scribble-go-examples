package message

// Next is Next(int) signature.
type Next struct {
	V int
}

// GetOp returns the label.
func (*Next) GetOp() string { return "Next" }

// Done is Done(int) signature.
type Done struct {
	V int
}

// GetOp returns the label.
func (*Done) GetOp() string { return "Done" }
