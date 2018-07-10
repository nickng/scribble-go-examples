package message

type Sort struct {
	V int
}

// GetOp returns the label.
func (*Sort) GetOp() string { return "Sort" }

type Match struct {
	V string
}

// GetOp returns the label.
func (*Match) GetOp() string { return "Match" }

type Done struct {
	V int
}

// GetOp returns the label.
func (*Done) GetOp() string { return "Done" }

type Gather struct {
	V string
}

// GetOp returns the label.
func (*Gather) GetOp() string { return "Gather" }
