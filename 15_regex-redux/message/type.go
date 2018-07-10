package message

type Count struct {
	V string
}

// GetOp returns the label.
func (*Count) GetOp() string { return "Count" }

type Measure struct {
	V int
}

// GetOp returns the label.
func (*Measure) GetOp() string { return "Measure" }

type Donec struct {
	V string
}

// GetOp returns the label.
func (*Donec) GetOp() string { return "Donec" }

type Len struct {
	V int
}

// GetOp returns the label.
func (*Len) GetOp() string { return "Len" }
