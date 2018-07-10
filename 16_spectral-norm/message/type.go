package message

type Times struct {
	V int
}

// GetOp returns the label.
func (*Times) GetOp() string { return "Times" }

type Done struct {
	V int
}

// GetOp returns the label.
func (*Done) GetOp() string { return "Done" }

type Next struct {
	V int
}

// GetOp returns the label.
func (*Next) GetOp() string { return "Next" }

type TimeStr struct {
	V int
}

// GetOp returns the label.
func (*TimeStr) GetOp() string { return "TimeStr" }

type End struct {
	V int
}

// GetOp returns the label.
func (*End) GetOp() string { return "End" }
